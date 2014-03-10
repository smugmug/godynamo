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

// Support for the DynamoDB Scan endpoint.
//
// example use:
//
// tests/scan-livestest.go
//
package scan

import (
	"fmt"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"encoding/json"
 	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME   = "Scan"
	SCAN_ENDPOINT   = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	OP_EQ           = "EQ"
	OP_NE           = "NE"
	OP_LE           = "LE"
	OP_LT           = "LT"
	OP_GE           = "GE"
	OP_GT           = "GT"
	OP_NULL         = "NULL"
	OP_NOT_NULL     = "NOT_NULL"
	OP_CONTAINS     = "CONTAINS"
	OP_NOT_CONTAINS = "NOT_CONTAINS"
	OP_BEGINS_WITH  = "BEGINS_WITH"
	OP_IN           = "IN"
	OP_BETWEEN      = "BETWEEN"
	LIMIT           = 10000 // limit of scan unless set
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

type ScanFilters map[string] ScanFilter

type scanFilters ScanFilters

func (s ScanFilters) MarshalJSON() ([]byte,error) {
	if len(s) == 0 {
		var i interface{} = nil
		return json.Marshal(i)
	} else {
		return json.Marshal(scanFilters(s))
	}
}

type ScanFilter struct {
	AttributeValueList []ep.AttributeValue
	ComparisonOperator ComparisonOperator
}

type Scan struct {
	AttributesToGet ep.AttributesToGet
	ExclusiveStartKey ep.Item
	ReturnConsumedCapacity ep.ReturnConsumedCapacity
	Limit ep.NullableUInt64
	ScanFilter ScanFilters
	Select ep.Select
	Segment ep.NullableUInt64
	TableName string
	TotalSegments ep.NullableUInt64
}

type scan Scan

// scanParallel is identical to Scan with one exception. The "Segment" field will
// not evaluate the null in the case of the int value of the TotalSegments being nonzero,
// it will be treated numerically
type scanParallel struct {
	AttributesToGet ep.AttributesToGet
	ExclusiveStartKey ep.Item
	ReturnConsumedCapacity ep.ReturnConsumedCapacity
	Limit ep.NullableUInt64
	ScanFilter ScanFilters
	Select ep.Select
	Segment uint64
	TableName string
	TotalSegments uint64
}

// MarshalJSON for Scan types will test the TotalSegments and if they are nonzero, allow the Segment
// field to be set to zero instead of being null'd, which is the default behavior of both
// Segment and TotalSegments are defaulted to zero
func (s Scan) MarshalJSON() ([]byte, error) {
	if s.TotalSegments != 0 {
		var sp scanParallel
		sp.AttributesToGet = s.AttributesToGet
		sp.ExclusiveStartKey = s.ExclusiveStartKey
		sp.ReturnConsumedCapacity = s.ReturnConsumedCapacity
		sp.Limit = s.Limit
		sp.ScanFilter = s.ScanFilter
		sp.Select = s.Select
		sp.Segment = uint64(s.Segment)
		sp.TableName = s.TableName
		sp.TotalSegments = uint64(s.TotalSegments)
		return json.Marshal(sp)
	} else {
		var si scan
		si = scan(s)
		return json.Marshal(si)
	}
}

// NewScan returns a pointer to an instantiation of the Scan struct.
func NewScan() (*Scan) {
	s := new(Scan)
	s.AttributesToGet	= make(ep.AttributesToGet,0)
	s.ExclusiveStartKey	= make(ep.Item)
	s.ScanFilter		= make(ScanFilters)
	return s
}

type Request Scan

type Response struct {
	Count uint64
	Items []ep.Item
	LastEvaluatedKey ep.Item
	ConsumedCapacity ep.ConsumedCapacity
	ScannedCount uint64
}

// NewResponse will return a pointer to an initialized Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.Items = make([]ep.Item,0)
	r.LastEvaluatedKey = make(ep.Item)
	return r
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(Scan(r))
}

// ValidOp determines if an operation is in the approved list.
func ValidOp(op string) bool {
	return (op == OP_EQ           ||
		op == OP_NE           ||
		op == OP_LE           ||
		op == OP_LT           ||
		op == OP_GE           ||
		op == OP_GT           ||
		op == OP_NULL         ||
		op == OP_NOT_NULL     ||
		op == OP_CONTAINS     ||
		op == OP_NOT_CONTAINS ||
		op == OP_BEGINS_WITH  ||
		op == OP_IN           ||
		op == OP_BETWEEN)
}

// EndpointReq implements the Endpoint interface.
func (s Scan) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("scan(Scan).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&s,SCAN_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (Scan(req)).EndpointReq()
}
