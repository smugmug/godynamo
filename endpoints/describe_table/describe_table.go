// Copyright (c) 2013,2014 SmugMug, Inc. All rights reserved.
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

// Support for the DynamoDB DescribeTable endpoint.
//
// example use:
//
// tests/create_table-livestest.go, which contains a DescribeTable invocation
//
package describe_table

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME      = "DescribeTable"
	DESCTABLE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	ACTIVE             = "ACTIVE"
)

type Describe struct {
	TableName string
}

type Request Describe

type Response struct {
	Table struct {
		AttributeDefinitions ep.AttributeDefinitions
		CreationDateTime float64
		ItemCount uint64
		KeySchema ep.KeySchema
		LocalSecondaryIndexes []ep.LocalSecondaryIndex
		ProvisionedThroughput ep.ProvisionedThroughputDesc
		TableName string
		TableSizeBytes uint64
		TableStatus string
	}
}

// NewResponse returns a pointer to an instantiation of the local Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.Table.KeySchema             = make(ep.KeySchema,0)
	r.Table.LocalSecondaryIndexes = make([]ep.LocalSecondaryIndex,0)
	return r
}

type StatusResult struct {
	StatusResult bool
}

// PollTableStatus allows the caller to poll a table for a specific status.
func PollTableStatus(tablename string,status string,tries int) (bool,error) {
	// aws docs informs us to poll the describe endpoint until the table
	// "status" is status for this tablename
	wait := time.Duration(2 * time.Second)

	for i:=0; i<tries; i++ {
		active,err := IsTableStatus(tablename,status)
		if err != nil {
			e := fmt.Sprintf("describe_table.PollStatus:%s",
				err.Error())
			return false,errors.New(e)
		}
		if active {
			return active,nil
		}
		time.Sleep(wait) // wait for table to become ACTIVE
	}
	return false,nil
}

// IsTableStatus will test the equality status of a table.
func IsTableStatus(tablename string,status string) (bool,error) {
	d := ep.Endpoint(Describe{TableName:tablename})
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("describe_table.IsTableStatus " +
			"auth must be v4")
		return false,errors.New(e)
	}
	s_resp,s_code,s_err := authreq.RetryReq_V4(d,DESCTABLE_ENDPOINT)
	if s_err != nil {
		e := fmt.Sprintf("describe_table.IsTableStatus: " +
			"check on %s err %s",
			tablename,s_err.Error())
		// if not a 500 problem, don't retry
		if !ep.ServerErr(s_code) {
			return false,errors.New(e)
		}
	}
	if s_resp != "" && s_code == http.StatusOK {
		var resp_json Response
		um_err := json.Unmarshal([]byte(s_resp), &resp_json)
		if um_err != nil {
			um_msg := fmt.Sprintf("describe_table.IsTableStatus:" +
				"cannot unmarshal %s, err: %s\ncheck " +
				"table creation of %s manually",
				s_resp,um_err.Error(),tablename)
			return false,errors.New(um_msg)
		}
		return (resp_json.Table.TableStatus == status),nil
	}
	e := fmt.Sprintf("describe_table.IsTableStatus:does %s exist?",tablename)
	return false,errors.New(e)
}

// TableExists test for table exists: exploit the fact that aws reports 4xx for tables that don't exist.
func (desc Describe) TableExists() (bool,error) {
	_,code,err := desc.EndpointReq()
	if err != nil {
		e := fmt.Sprintf("describe_table.TableExists " +
			"%s",err.Error())
		return false,errors.New(e)
	}
	return (code == http.StatusOK),nil
}

// EndpointReq implements the Endpoint interface.
func (desc Describe) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("describe_table.EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&desc,DESCTABLE_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (Describe(req)).EndpointReq()
}
