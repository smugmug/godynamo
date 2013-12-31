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

// Support for the DynamoDB CreateTable endpoint.
//
// example use:
//
// see tests/create_table-livestest.go
//
package create_table

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME        = "CreateTable"
	CREATETABLE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
)

type Create struct {
	TableName string
	AttributeDefinitions ep.AttributeDefinitions
	KeySchema ep.KeySchema
	LocalSecondaryIndexes ep.LocalSecondaryIndexes
	ProvisionedThroughput ep.ProvisionedThroughput
}

type create Create

type Request Create

// NewCreate returns a pointer to an instantiation of the Create struct.
func NewCreate() (*Create) {
	c := new(Create)
	c.AttributeDefinitions  = make(ep.AttributeDefinitions,0)
	c.KeySchema             = make(ep.KeySchema,0)
	c.LocalSecondaryIndexes = make(ep.LocalSecondaryIndexes,0)
	return c
}

type Response struct {
	TableDescription struct {
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
	r.TableDescription.KeySchema             = make(ep.KeySchema,0)
	r.TableDescription.LocalSecondaryIndexes = make([]ep.LocalSecondaryIndex,0)
	return r
}

func (c Create) MarshalJSON() ([]byte, error) {
	if len(c.LocalSecondaryIndexes) > 5 {
		e := fmt.Sprintf("endpoint.Create.MarshalJSON: LocalSecondaryIndexes > 5")
		return nil, errors.New(e)
	}
	if (!ValidTableName(c.TableName)) {
		e := fmt.Sprintf("endpoint.Create.MarshalJSON: TableName %s bad len",c.TableName)
		return nil, errors.New(e)
	}
	var ci create
	ci.TableName = c.TableName
	ci.KeySchema = c.KeySchema
	ci.AttributeDefinitions  = c.AttributeDefinitions
	ci.LocalSecondaryIndexes = c.LocalSecondaryIndexes
	ci.ProvisionedThroughput = c.ProvisionedThroughput
	return json.Marshal(ci)
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(Create(r))
}

// EndpointReq implements the Endpoint interface.
func (c Create) EndpointReq() (string,int,error) {
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("create_table(Create).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&c,CREATETABLE_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (Create(req)).EndpointReq()
}

// ValidTable is a local validator that helps callers determine if a table name is too long.
func ValidTableName(t string) bool {
	l := len([]byte(t))
	return (l > 3) && (l < 256)
}
