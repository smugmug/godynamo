// Manages AWS Auth v4 requests to DynamoDB.
// See http://docs.aws.amazon.com/general/latest/gr/signature-version-4.html
// for more information on v4 signed requests. For examples, see any of
// the package in the `endpoints` directory.
package auth_v4

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/smugmug/godynamo/auth_v4/tasks"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
	"hash"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	IAM_WARN_MESSAGE = "check roles sources and make sure you have run one of the roles " +
		"management functions in package conf_iam, such as GoIAM"
)

// Client for executing requests.
var Client *http.Client

// Initialize package-scoped client.
func init() {
	// The timeout seems too-long, but it accomodates the exponential decay retry loop.
	// Programs using this can either change this directly or use goroutine timeouts
	// to impose a local minimum.
	tr := &http.Transport{MaxIdleConnsPerHost: 250,
		ResponseHeaderTimeout: time.Duration(20) * time.Second}
	Client = &http.Client{Transport: tr}
}

// GetRespReqID retrieves the unique identifier from the AWS Response
func GetRespReqID(response http.Response) (string, error) {
	if amz_reqid_list, reqid_ok := response.Header["X-Amzn-Requestid"]; reqid_ok {
		if len(amz_reqid_list) == 1 {
			return amz_reqid_list[0], nil
		}
	}
	return "", errors.New("auth_v4.GetRespReqID: no X-Amzn-Requestid found")
}

// MatchCheckSum will perform a local crc32 on the response body and match it against the aws crc32
// *** WARNING ***
// There seems to be a mismatch between what Go calculates and what AWS (java?) calculates here,
// I believe related to utf8 (go) vs utf16 (java), but I don't know enough about encodings to
// solve it. So until that issue is solved, don't use this.
func MatchCheckSum(response http.Response, respbody []byte) (bool, error) {
	if amz_crc_list, crc_ok := response.Header["X-Amz-Crc32"]; crc_ok {
		if len(amz_crc_list) == 1 {
			amz_crc_int32, amz_crc32_err := strconv.Atoi(amz_crc_list[0])
			if amz_crc32_err == nil {
				client_crc_int32 := int(crc32.ChecksumIEEE(respbody))
				if amz_crc_int32 != client_crc_int32 {
					_ = fmt.Sprintf("auth_v4.MatchCheckSum: resp crc mismatch: amz %d client %d",
						amz_crc_int32, client_crc_int32)
					return false, nil
				}
			}
		} else {
			return false, errors.New("auth_v4.MatchCheckSum: X-Amz-Crc32 malformed")
		}
	} else {
		return false, errors.New("auth_v4.MatchCheckSum: no X-Amz-Crc32 found")
	}
	return true, nil
}

// rawReqAll takes each parameter independently, forms and signs the request, and returns the
// result (and error codes).
func rawReqAll(reqJSON []byte, amzTarget string, useIAM bool, url, host, port, zone, IAMSecret, IAMAccessKey, IAMToken, authSecret, authAccessKey string) ([]byte, string, int, error) {

	// initialize req with body reader
	body := strings.NewReader(string(reqJSON))
	request, req_err := http.NewRequest(aws_const.METHOD, url, body)
	if req_err != nil {
		e := fmt.Sprintf("auth_v4.rawReqAll:failed init conn %s", req_err.Error())
		return nil, "", 0, errors.New(e)
	}

	// add headers
	// content type
	request.Header.Add(aws_const.CONTENT_TYPE_HDR, aws_const.CTYPE)
	// amz target
	request.Header.Add(aws_const.AMZ_TARGET_HDR, amzTarget)
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
		host,
		port,
		request.Header.Get(aws_const.X_AMZ_DATE_HDR),
		request.Header.Get(aws_const.AMZ_TARGET_HDR),
		hexPayload)
	str2sign := tasks.String2Sign(now, canonical_request,
		zone,
		service)

	// obtain the aws secret credential from the global Auth or from IAM
	var secret string
	if useIAM == true {
		secret = IAMSecret
	} else {
		secret = authSecret
	}
	if secret == "" {
		panic("auth_v4.rawReqAll: no Secret defined; " + IAM_WARN_MESSAGE)
	}

	signature := tasks.MakeSignature(str2sign, zone, service, secret)

	// obtain the aws accessKey credential from the global Auth or from IAM
	// if using IAM, read the token while we have the lock
	var accessKey, token string
	if useIAM == true {
		accessKey = IAMAccessKey
		token = IAMToken
	} else {
		accessKey = authAccessKey
	}
	if accessKey == "" {
		panic("auth_v4.rawReqAll: no Access Key defined; " + IAM_WARN_MESSAGE)
	}

	v4auth := "AWS4-HMAC-SHA256 Credential=" + accessKey +
		"/" + now.UTC().Format(aws_const.ISODATEFMT) + "/" +
		zone + "/" + service + "/aws4_request," +
		"SignedHeaders=content-type;host;x-amz-date;x-amz-target," +
		"Signature=" + signature

	request.Header.Add("Authorization", v4auth)
	if useIAM == true {
		if token == "" {
			panic("auth_v4.rawReqAll: no Token defined;" + IAM_WARN_MESSAGE)
		}
		request.Header.Add(aws_const.X_AMZ_SECURITY_TOKEN_HDR, token)
	}

	// where we finally send req to aws
	response, rsp_err := Client.Do(request)

	if rsp_err != nil {
		return nil, "", 0, rsp_err
	}

	respbody, read_err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if read_err != nil && read_err != io.EOF {
		e := fmt.Sprintf("auth_v4.rawReqAll:err reading resp body: %s", read_err.Error())
		return nil, "", 0, errors.New(e)
	}

	amz_requestid, amz_requestid_err := GetRespReqID(*response)
	if amz_requestid_err != nil {
		return nil, "", 0, amz_requestid_err
	}

	return respbody, amz_requestid, response.StatusCode, nil
}

// RawReqWithConf will sign and transmit the request to the AWS DynamoDB endpoint.
// reqJSON is the json request
// amzTarget is the dynamoDB endpoint
// c is the configuration struct
// returns []byte respBody, string aws reqID, int http code, error
func RawReqWithConf(reqJSON []byte, amzTarget string, c *conf.AWS_Conf) ([]byte, string, int, error) {
	if !conf.IsValid(c) {
		return nil, "", 0, errors.New("auth_v4.RawReqWithConf: conf not valid")
	}
	// shadow conf vars in a read lock to minimize contention
	var our_c conf.AWS_Conf
	cp_err := our_c.Copy(c)
	if cp_err != nil {
		return nil, "", 0, cp_err
	}
	return rawReqAll(
		reqJSON,
		amzTarget,
		our_c.UseIAM,
		our_c.Network.DynamoDB.URL,
		our_c.Network.DynamoDB.Host,
		our_c.Network.DynamoDB.Port,
		our_c.Network.DynamoDB.Zone,
		our_c.IAM.Credentials.Secret,
		our_c.IAM.Credentials.AccessKey,
		our_c.IAM.Credentials.Token,
		our_c.Auth.Secret,
		our_c.Auth.AccessKey)
}

// RawReq will sign and transmit the request to the AWS DynamoDB endpoint.
// This method uses the global conf.Vals to obtain credential and configuation information.
func RawReq(reqJSON []byte, amzTarget string) ([]byte, string, int, error) {
	return RawReqWithConf(reqJSON, amzTarget, &conf.Vals)
}

// Req  will sign and transmit the request to the AWS DynamoDB endpoint.
// This method uses the global conf.Vals to obtain credential and configuation information.
// At one point, RawReq and Req were different, now RawReq is just an alias.
func Req(reqJSON []byte, amzTarget string) ([]byte, string, int, error) {
	return RawReqWithConf(reqJSON, amzTarget, &conf.Vals)
}

// ReqConf is just a wrapper for RawReq if we need to massage data
// before dispatch. Uses parameterized conf.
func ReqWithConf(reqJSON []byte, amzTarget string, c *conf.AWS_Conf) ([]byte, string, int, error) {
	return RawReqWithConf(reqJSON, amzTarget, c)
}
