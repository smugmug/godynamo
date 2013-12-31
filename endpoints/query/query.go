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

// Support for the DynamoDB Query endpoint.
//
// example use:
//
// tests/query-livestest.go
//
package query

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
 	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME    = "Query"
	QUERY_ENDPOINT   = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	OP_EQ            = "EQ"
	OP_LE            = "LE"
	OP_LT            = "LT"
	OP_GE            = "GE"
	OP_GT            = "GT"
	OP_BEGINS_WITH   = "BEGINS_WITH"
	OP_BETWEEN       = "BETWEEN"
	LIMIT            = 10000 // limit of query unless set
)

type ComparisonOperator string

type comparisonOperator ComparisonOperator

func (c ComparisonOperator) MarshalJSON() ([]byte, error) {
	var ci comparisonOperator
	var cs = string(c)
	if ValidOp(cs) {
		ci = comparisonOperator(cs)
		return json.Marshal(ci)
	} else {
		e := fmt.Sprintf("ComparisonOperator.MarshalJSON: op %s is not valid",cs)
		return nil, errors.New(e)
	}
}

type KeyConditions map[string] KeyCondition

type KeyCondition struct {
	AttributeValueList []ep.AttributeValue
	ComparisonOperator ComparisonOperator
}

// ScanIndexForward as an interface is no accident. I use this
// type-above-bool to marshal and unmarshal into the correct bool
// defaults for AWS, which happen to be the opposite of Go bool
// defaults (aws defaults to true, go defaults to false)
type Query struct {
	AttributesToGet ep.AttributesToGet
	ConsistentRead bool
	ExclusiveStartKey ep.Item
	IndexName ep.NullableString
	KeyConditions KeyConditions
	Limit ep.NullableUInt64
	ReturnConsumedCapacity ep.ReturnConsumedCapacity
	ScanIndexForward interface{}
	TableName string
	Select ep.Select
}

// NewQuery returns a pointer to an instantiation of the Query struct.
func NewQuery() (*Query) {
	q := new(Query)
	q.AttributesToGet	= make(ep.AttributesToGet,0)
	q.ExclusiveStartKey	= make(ep.Item)
	q.KeyConditions		= make(KeyConditions)
	return q
}

type Request Query

type query Query

type Response struct {
	Count uint64
	Items []ep.Item
	LastEvaluatedKey ep.Item
	ConsumedCapacity ep.ConsumedCapacity
}

// NewResponse will return a pointer to an initialized Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.Items = make([]ep.Item,0)
	r.LastEvaluatedKey = make(ep.Item)
	return r
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(Query(r))
}

func (q Query) MarshalJSON() ([]byte, error) {
	qi := query(q)
	_,b_ok := q.ScanIndexForward.(bool)
	if b_ok {
		qi.ScanIndexForward = q.ScanIndexForward
	} else {
		qi.ScanIndexForward = true
	}
	return json.Marshal(qi)
}

// ValidOp determines if an operation is in the approved list.
func ValidOp(op string) bool {
	return (op == OP_EQ          ||
		op == OP_LE          ||
		op == OP_LT          ||
		op == OP_GE          ||
		op == OP_GT          ||
		op == OP_BEGINS_WITH ||
		op == OP_BETWEEN)
}

// EndpointReq implements the Endpoint interface.
func (q Query) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("query(Query).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&q,QUERY_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (Query(req)).EndpointReq()
}
