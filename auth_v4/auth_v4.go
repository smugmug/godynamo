// Copyright (c) 2013, SmugMug, Inc. All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
//       copyright notice, this list of conditions and the following
//       disclaimer in the documentation and/or other materials provided
//       with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY SMUGMUG, INC. ``AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
// PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL SMUGMUG, INC. BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE
// GOODS OR SERVICES;LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
// IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
// OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
// ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Manages AWS Auth v4 requests to DynamoDB.
// See http://docs.aws.amazon.com/general/latest/gr/signature-version-4.html
// for more information on v4 signed requests. For examples, see any of
// the package in the `endpoints` directory.
package auth_v4

import (
	"net/url"
	"net/http"
	"fmt"
	"strconv"
	"errors"
	"io"
	"io/ioutil"
	"encoding/json"
	"strings"
	"hash"
	"time"
	"crypto/sha256"
	"hash/crc32"
	"encoding/hex"
	"github.com/smugmug/godynamo/auth_v4/tasks"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	IAM_WARN_MESSAGE = "check roles sources and make sure you have run one of the roles " +
		"management functions in package conf_iam, such as GoIAM"
)

// Client for executing requests.
var Client *http.Client

// Initialize package-scoped client.
func init() {
	tr := &http.Transport{ResponseHeaderTimeout: time.Duration(20) * time.Second}
	Client = &http.Client{Transport:tr}
}

// GetRespReqID retrieves the unique identifier from the AWS Response
func GetRespReqID(response http.Response) (string,error) {
	if amz_reqid_list,reqid_ok := response.Header["X-Amzn-Requestid"]; reqid_ok {
		if len(amz_reqid_list) == 1 {
			return amz_reqid_list[0],nil
		}
	}
	return "",errors.New("auth_v4.GetRespReqID: no X-Amzn-Requestid found")
}

// MatchCheckSum will perform a local crc32 on the response body and match it against the aws crc32
// *** WARNING ***
// There seems to be a mismatch between what Go calculates and what AWS (java?) calculates here,
// I believe related to utf8 (go) vs utf16 (java), but I don't know enough about encodings to
// solve it. So until that issue is solved, don't use this.
func MatchCheckSum(response http.Response,respbody []byte) (bool,error) {
	if amz_crc_list,crc_ok := response.Header["X-Amz-Crc32"]; crc_ok {
		if len(amz_crc_list) == 1 {
			amz_crc_int32,amz_crc32_err := strconv.Atoi(amz_crc_list[0])
			if amz_crc32_err == nil {
				client_crc_int32 := int(crc32.ChecksumIEEE(respbody))
				if amz_crc_int32 !=  client_crc_int32 {
					_ = fmt.Sprintf("auth_v4.RawReq: resp crc mismatch: amz %d client %d",
						amz_crc_int32,client_crc_int32)
					return false,nil
				}
			}
		} else {
			return false,errors.New("auth_v4.MatchCheckSum: X-Amz-Crc32 malformed")
		}
	} else {
		return false,errors.New("auth_v4.MatchCheckSum: no X-Amz-Crc32 found")
	}
	return true,nil
}

// RawReq will sign and transmit the request to the AWS DynanoDB endpoint.
// This method is DynamoDB-specific.
func RawReq(reqJSON []byte,amzTarget string) (string,string,int,error) {
	url,url_err := url.Parse(conf.Vals.Network.DynamoDB.URL)
	if url_err != nil {
		e := "auth_v4.RawReq:parse " +
			conf.Vals.Network.DynamoDB.URL +
			" " + url_err.Error()
		return "","",0,errors.New(e)
	}

	// initialize req with body reader
	body := strings.NewReader(string(reqJSON))
	request,req_err := http.NewRequest(aws_const.METHOD,url.String(),body)
	if req_err != nil {
		e := fmt.Sprintf("auth_v4.RawReq:failed init conn %s",req_err.Error())
		return "","",0,errors.New(e)
	}

	// add headers
	// content type
	request.Header.Add(aws_const.CONTENT_TYPE_HDR,aws_const.CTYPE)
	// amz target
	request.Header.Add(aws_const.AMZ_TARGET_HDR,amzTarget)
	// dates
	now := time.Now()
	request.Header.Add(aws_const.X_AMZ_DATE_HDR,
		now.UTC().Format(aws_const.ISO8601FMT_CONDENSED))

	// encode request json payload
	var h256 hash.Hash = sha256.New()
	h256.Write(reqJSON)
	hexPayload := string(hex.EncodeToString([]byte(h256.Sum(nil))))

	// create the various signed formats aws uses for v4 signed reqs
	service := strings.ToLower(aws_const.DYNAMODB)
	canonical_request := tasks.CanonicalRequest(
		conf.Vals.Network.DynamoDB.Host,
		request.Header.Get(aws_const.X_AMZ_DATE_HDR),
		request.Header.Get(aws_const.AMZ_TARGET_HDR),
		hexPayload)
	str2sign := tasks.String2Sign(now,canonical_request,
		conf.Vals.Network.DynamoDB.Zone,
		service)

	// obtain the aws secret credential from the global Auth or from IAM
	var secret string
	if conf.Vals.UseIAM == true {
		conf.Vals.ConfLock.RLock()
		secret = conf.Vals.IAM.Credentials.Secret
		conf.Vals.ConfLock.RUnlock()
	} else {
		secret = conf.Vals.Auth.Secret
	}
	if secret == "" {
		panic("auth_v4.cacheable_hmacs: no Secret defined; " + IAM_WARN_MESSAGE)
	}

	signature := tasks.MakeSignature(str2sign,conf.Vals.Network.DynamoDB.Zone,service,secret)

	// obtain the aws accessKey credential from the global Auth or from IAM
	// if using IAM, read the token while we have the lock
	var accessKey,token string
	if conf.Vals.UseIAM == true {
		conf.Vals.ConfLock.RLock()
		accessKey = conf.Vals.IAM.Credentials.AccessKey
		token = conf.Vals.IAM.Credentials.Token
		conf.Vals.ConfLock.RUnlock()
	} else {
		accessKey = conf.Vals.Auth.AccessKey
	}
	if accessKey == "" {
		panic("auth_v4.RawReq: no Access Key defined; " + IAM_WARN_MESSAGE)
	}

	v4auth := "AWS4-HMAC-SHA256 Credential=" + accessKey +
		"/" + now.UTC().Format(aws_const.ISODATEFMT) + "/" +
		conf.Vals.Network.DynamoDB.Zone + "/" + service + "/aws4_request," +
		"SignedHeaders=content-type;host;x-amz-date;x-amz-target," +
		"Signature=" + signature
	request.Header.Add("Authorization",v4auth)
	if conf.Vals.UseIAM == true {
		if token == "" {
			panic("auth_v4.RawReq: no Token defined;" + IAM_WARN_MESSAGE)
		}
		request.Header.Add(aws_const.X_AMZ_SECURITY_TOKEN_HDR,token)
	}

	// where we finally send req to aws
	response,rsp_err := Client.Do(request)

	if rsp_err != nil {
		return "","",0,rsp_err
	}
	defer response.Body.Close()
	respbody,read_err := ioutil.ReadAll(response.Body)
	if read_err != nil && read_err != io.EOF {
		e := fmt.Sprintf("auth_v4.RawReq:err reading resp body: %s",read_err.Error())
		return "","",0,errors.New(e)
	}

	amz_requestid,amz_requestid_err := GetRespReqID(*response)
	if amz_requestid_err != nil {
		return "","",0,amz_requestid_err
	}

	return string(respbody),amz_requestid,response.StatusCode,nil
}

// Req prepares a RawReq call from either a ep.Endpoint instance or a []byte representation
// serialization of the request payload. DynamoDB-specific.
func Req(v interface{},amzTarget string) (string,string,int,error) {
	// we take two types here, either an ep.Endpoint implementor, or
	// a []byte representing the marshaled json
	_,ep_ok := interface{}(v).(ep.Endpoint)
	if ep_ok {
		reqJSON,json_err := json.Marshal(v);
		if json_err != nil {
			return "","",0,json_err
		}
		return RawReq(reqJSON,amzTarget)
	}
	v_bytes,v_ok := v.([]byte)
	if v_ok {
		return RawReq(v_bytes,amzTarget)
	}
	return "","",0,errors.New("auth_v4.Req:v unknown type")
}
