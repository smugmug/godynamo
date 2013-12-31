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

// Support for the DynamoDB DeleteItem endpoint.
//
// example use:
//
// tests/delete_item-livestest.go
//
package delete_item

import (
	"fmt"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME           = "DeleteItem"
	DELETEITEM_ENDPOINT     = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// the permitted ReturnValues flags for this op
	RETVAL_ALL_OLD		= "ALL_OLD"
	RETVAL_NONE             = ep.RETVAL_NONE
)

type Delete struct {
	TableName string
	Key ep.Item
	Expected ep.Expected
	ReturnValues ep.ReturnValues
}

// NewDelete returns a pointer to an instantiation of the Delete struct.
func NewDelete() (*Delete) {
	d := new(Delete)
	d.Key = make(ep.Item)
	d.Expected = make(ep.Expected)
	return d
}

type Request Delete

type Response ep.AttributesResponse

// NewResponse returns a pointer to an instantiation of the local Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.Attributes = make(ep.Item)
	return r
}

// EndpointReq implements the Endpoint interface.
func (del Delete) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("delete_item(Delete).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&del,DELETEITEM_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (Delete(req)).EndpointReq()
}
