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

// Support for the DynamoDB PutItem endpoint.
//
// example use:
//
// tests/put_item-livestest.go
//
package put_item

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME           = "PutItem"
	PUTITEM_ENDPOINT        = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// the permitted ReturnValues flags for this op
	RETVAL_ALL_OLD		= "ALL_OLD"
	RETVAL_NONE             = ep.RETVAL_NONE
)

type Put struct {
	TableName string
	Item ep.Item
	Expected ep.Expected
	ReturnValues ep.ReturnValues
}

// NewPut will return a pointer to an initialized Put struct.
func NewPut() (*Put) {
	p := new(Put)
	p.Item = make(ep.Item)
	p.Expected = make(ep.Expected)
	return p
}

type Request Put

type put Put

type Response ep.AttributesResponse

// NewResponse will return a pointer to an initialized Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.Attributes = make(ep.Item)
	return r
}

func (p Put) MarshalJSON() ([]byte, error) {
	var pi put
	if p.Expected == nil || len(p.Expected) == 0 {
		pi.Expected = nil
	} else {
		pi.Expected = p.Expected
	}
	pi.TableName = p.TableName
	pi.ReturnValues = p.ReturnValues
	pi.Item = p.Item
	return json.Marshal(pi)
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(Put(r))
}

// EndpointReq implements the Endpoint interface.
func (put Put) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("put_item(Put).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&put,PUTITEM_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (Put(req)).EndpointReq()
}

// ValidItem validates the size of a json serialization of an Item.
// AWS says items can only be 64k bytes binary
// potential utf8 (utf8 chars *can* occupy 4 bytes)
func ValidItem(i string) bool {
	return !(len([]byte(i)) > 65536)
}
