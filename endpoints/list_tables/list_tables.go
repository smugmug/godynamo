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

// Support for the DynamoDB ListTables endpoint.
//
// example use:
//
// tests/create_table-livestest.go, which contains a ListTables invocation
//
package list_tables

import (
	"fmt"
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME              = "ListTables"
	LISTTABLE_ENDPOINT         = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	EXCLUSIVE_START_TABLE_NAME = "ExclusiveStartTableName"
	LIMIT                      = "Limit"
	LAST_EVALUATED_TABLE_NAME  = "LastEvaluatedTableName"
	AWS_LIMIT                  = 100
)

type List struct {
	ExclusiveStartTableName ep.NullableString
	Limit uint64
}

type list List

type Request List

type Response struct {
	TableNames []string
	LastEvaluatedTableName string
}

// NewResponse returns a pointer to an instantiation of the local Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.TableNames = make([]string,0)
	return r
}

func limit(l uint64) uint64 {
	// no need to error here, just override if need be
	if l == 0 || l > AWS_LIMIT {
		return AWS_LIMIT
	} else {
		return l
	}
}

// custom json marshal
func (l List) MarshalJSON() ([]byte, error) {
	var li list
	li.ExclusiveStartTableName = l.ExclusiveStartTableName
	li.Limit = limit(l.Limit)
	return json.Marshal(li)
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(List(r))
}

// EndpointReq implements the Endpoint interface.
func (list List) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("list_table(List).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&list,LISTTABLE_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (List(req)).EndpointReq()
}
