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

// A collection of types and methods common to all of the endpoint/* packages.
package endpoint

import (
	"fmt"
	"errors"
	"encoding/json"
	"encoding/base64"
	"net/http"
	"strconv"
)

const (
	HASH              = "HASH"
	RANGE             = "RANGE"
	HASH_KEY_ELEMENT  = "HashKeyElement"
	RANGE_KEY_ELEMENT = "RangeKeyElement"
	S                 = "S"
	N                 = "N"
	B                 = "B"
	SS                = "SS"
	NS                = "NS"
	BS                = "BS"
	RETVAL_NONE       = "NONE"
	ALL               = "ALL"
	SIZE              = "SIZE"
	TOTAL             = "TOTAL"
	KEYS_ONLY         = "KEYS_ONLY"
	INCLUDE           = "INCLUDE"
	SELECT_ALL        = "ALL_ATTRIBUTES"
	SELECT_PROJECTED  = "ALL_PROJECTED_ATTRIBUTES"
	SELECT_SPECIFIC   = "SPECIFIC_ATTRIBUTES"
	SELECT_COUNT      = "COUNT"
)

// AttributeValue from DynamoDB translated to Go.
// Numerics are rendered as strings for portability between numeric systems as detailed by AWS.
type AttributeValue struct {
	N string
	S string
	B string
	SS []string
	NS []string
	BS []string
	Type string
}

type attributeValue AttributeValue

type attr_N  struct { N string }
type attr_S  struct { S string }
type attr_B  struct { B string }     // assume already base64 encoded
type attr_NS struct { NS []string }
type attr_SS struct { SS []string }
type attr_BS struct { BS []string }  // assume already base64 encoded

// Empty determines if an AttributeValue is vacuous.
func (a AttributeValue) Empty() bool {
	return a.N        == "" &&
		a.S       == "" &&
		a.B       == "" &&
		len(a.NS) == 0  &&
		len(a.SS) == 0  &&
		len(a.BS) == 0  &&
		a.Type    == ""
}

type AttributeDefinition struct {
	AttributeName string
	AttributeType string
}

type attributeDefinition AttributeDefinition

// Marshal the AWS Attribute Definition.
func (a AttributeDefinition) MarshalJSON() ([]byte, error) {
	if !(a.AttributeType == S || a.AttributeType == N || a.AttributeType == B ||
		a.AttributeType == SS || a.AttributeType == NS || a.AttributeType == BS) {
		e := fmt.Sprintf("endpoint.AttributeDefinition.MarshalJSON: AttributeType %s is not valid",
			a.AttributeType)
		return nil, errors.New(e)
	}
	var ai attributeDefinition
	ai.AttributeName = a.AttributeName
	ai.AttributeType = a.AttributeType
	return json.Marshal(ai)
}

type AttributeDefinitions []AttributeDefinition

// AWSParseFLoats normalizes numbers-as-strings for transport.
func AWSParseFloat(s string) (string,error) {
	// if it is an uint (no decimal), return it as such
	ui,uierr := strconv.ParseInt(s,0,64)
	if uierr == nil {
		return strconv.FormatInt(ui,10),nil
	}
	// try float, uint failed
	f,ferr := strconv.ParseFloat(s,64)
	if ferr != nil {
		return "",ferr
	}
	// aws accepts 38 decimal points, may have to change -1 to that
	return strconv.FormatFloat(f,'f',-1,64),nil
}

// AWSParseBinary can test if a string has already been encoded.
func AWSParseBinary(s string) (error) {
	_,err := base64.StdEncoding.DecodeString(s)
	return err
}

// UnmarshalJSON will assign the proper field in the AttributeValue type.
func (a *AttributeValue) UnmarshalJSON (data []byte) error {
	var ai attributeValue
	t_err := json.Unmarshal(data,&ai)
	if t_err != nil {
		return t_err
	}
	if ai.N != "" {
		a.N = ai.N
		a.Type = N
		return nil
	}
	if ai.S != "" {
		a.S = ai.S
		a.Type = S
		return nil
	}
	if len(ai.NS) != 0 {
		a.NS = ai.NS
		a.Type = NS
		return nil
	}
	if len(ai.SS) != 0 {
		a.SS = ai.SS
		a.Type = SS
		return nil
	}
	if ai.B != "" {
		a.B = ai.B
		a.Type = B
		return nil
	}
	if len(ai.BS) != 0 {
		a.BS = ai.BS
		a.Type = BS
		return nil
	}
	return nil
}

// MarshalJSON will emit JSON for an AttributeValue while skipping empty values.
func (a AttributeValue) MarshalJSON() ([]byte, error) {
	if a.N != "" {
		fs,ferr := AWSParseFloat(a.N)
		if ferr != nil {
			return nil,ferr
		}
		return json.Marshal(attr_N{N:fs})
	}
	if a.S != "" {
		return json.Marshal(attr_S{S:a.S})
	}
	if len(a.B) != 0 {
		b64err := AWSParseBinary(a.B)
		if b64err != nil {
			return nil,b64err
		}
		return json.Marshal(attr_B{B:a.B})
	}
	if len(a.NS) != 0 {
		var ans attr_NS
		for _,v := range a.NS {
			fs,ferr := AWSParseFloat(v)
			if ferr != nil {
				return nil,ferr
			}
			ans.NS = append(ans.NS,fs)
		}
		return json.Marshal(ans)
	}
	if len(a.SS) != 0 {
		var ass attr_SS
		ass.SS = a.SS
		return json.Marshal(ass)
	}
	if len(a.BS) != 0 {
		for _,v := range a.BS {
			b64err := AWSParseBinary(v)
			if b64err != nil {
				return nil,b64err
			}
		}
		return json.Marshal(attr_BS{BS:a.BS})
	}
	e := fmt.Sprintf("AttributeValue.MarshalJSON: no fields in %v",a)
	return nil, errors.New(e)
}

// Items are the data transport mechanism for Dynamo, mapping keys to AttributeValues
type Item map[string] AttributeValue

type item Item

func (i Item) MarshalJSON() ([]byte, error) {
	var ii item
	if len(i) == 0 {
		ii = nil
	} else {
		ii = item(i)
	}
	return json.Marshal(ii)
}

type ConsumedCapacityUnit float32

type ConsumedCapacity struct {
	CapacityUnits ConsumedCapacityUnit
	TableName string
}

// AttributeResponse is response for PutItem,UpdateItem,DeleteItem,etc
type AttributesResponse struct {
	Attributes Item
	ConsumedCapacity ConsumedCapacity
	ItemCollectionMetrics ItemCollectionMetrics
}

// AttributesToGet must be of len 1 or greater
type AttributesToGet []string

type attributesToGet AttributesToGet

func (a AttributesToGet) MarshalJSON() ([]byte, error) {
	var ai attributesToGet
	if len(a) == 0 {
		ai = nil
	} else {
		ai = attributesToGet(a)
	}
	return json.Marshal(ai)
}

// NullableString is a string that when empty is marshaled as null
// use this when there is an *optional* string parameter.
type NullableString string

func (n NullableString) MarshalJSON() ([]byte, error) {
	sn := string(n)
	if sn == "" {
		var i interface{}
		return json.Marshal(i)
	}
	return json.Marshal(sn)
}

// NullableUInt64 is a uint64 that when empty (0) is marshaled as null
// use this when there is an *optional* uint64 parameter
type NullableUInt64 uint64

func (n NullableUInt64) MarshalJSON() ([]byte, error) {
	in := uint64(n)
	if in == 0 {
		var i interface{}
		return json.Marshal(i)
	}
	return json.Marshal(in)
}

type ReturnValues string

type returnValues ReturnValues

func (r ReturnValues) MarshalJSON() ([]byte, error) {
	var ri returnValues
	if string(r) == "" {
		ri = returnValues(RETVAL_NONE)
	} else {
		ri = returnValues(r)
	}
	return json.Marshal(ri)
}

type ReturnConsumedCapacity string

type returnConsumedCapacity ReturnConsumedCapacity

func (r ReturnConsumedCapacity) MarshalJSON() ([]byte, error) {
	var ri returnConsumedCapacity
	if string(r) == "" {
		ri = returnConsumedCapacity(RETVAL_NONE)
	} else if (string(r) == RETVAL_NONE) || (string(r) == TOTAL) {
		ri = returnConsumedCapacity(r)
	} else {
		e := fmt.Sprintf("ReturnConsumedCapacity.MarshalJSON:%s bad",r)
		return nil,errors.New(e)
	}
	return json.Marshal(ri)
}

type ReturnItemCollectionMetrics string

type returnItemCollectionMetrics ReturnItemCollectionMetrics

func (r ReturnItemCollectionMetrics) MarshalJSON() ([]byte, error) {
	var ri returnItemCollectionMetrics
	if string(r) == "" {
		ri = returnItemCollectionMetrics(RETVAL_NONE)
	} else if (string(r) == RETVAL_NONE) || (string(r) == SIZE) {
		ri = returnItemCollectionMetrics(r)
	} else {
		e := fmt.Sprintf("ReturnItemCollectionMetrics.MarshalJSON:%s bad",r)
		return nil,errors.New(e)
	}
	return json.Marshal(ri)
}

type ItemCollectionMetrics struct {
	ItemCollectionKey Item
	SizeEstimateRangeGB [2]uint64
}

// Select is used by Query and Scan
type Select string

type s_elect Select // select is a go keyword

func (s Select) MarshalJSON() ([]byte, error) {
	var si s_elect
	ss := string(s)
	if ss == SELECT_ALL || ss == SELECT_PROJECTED ||
		ss == SELECT_SPECIFIC || ss == SELECT_COUNT {
		si = s_elect(ss)
		return json.Marshal(si)
	} else if ss == "" {
		var in interface{}
		return json.Marshal(in)
	} else {
		e := fmt.Sprintf("Select.MarshalJSON:%s bad",ss)
		return nil,errors.New(e)
	}
}

// Constraints models query constraints.
// Exists as an interface is no accident. This type-above-bool is used to
// marshal and unmarshal into the correct bool defaults for aws, which
// happen to be the opposite of Go bool defaults (aws defaults to true,
// go defaults to false)
type Constraints struct {
	Value AttributeValue
	Exists interface{}
}

type constraints_exists struct {
	Exists bool
}

type constraints Constraints

func (c *Constraints) UnmarshalJSON (data []byte) error {
	var ci constraints
	t_err := json.Unmarshal(data,&ci)
	if t_err != nil {
		return t_err
	}
	if ci.Exists == nil {
		c.Exists = true
	} else {
		_,b_ok := ci.Exists.(bool)
		if b_ok {
			c.Exists = ci.Exists
		} else {
			c.Exists = true
		}
	}
	c.Value = ci.Value
	return nil
}

func (c Constraints) MarshalJSON() ([]byte, error) {
	var ci constraints
	b,b_ok := c.Exists.(bool)
	if !b_ok {
		ci.Exists = true
		b = true
	} else {
		ci.Exists = b
	}
	if b == false {
		return json.Marshal(constraints_exists{Exists:false})
	}
	// Exists == true, needs a comparison val
	if c.Value.Empty() {
		return nil, errors.New("endpoint.Constraints.MarshalJSON: 'true' Exists Expected constraint needs nonempty comparison AttributeValue")
	}
	ci.Value = c.Value
	return json.Marshal(ci)
}

// Expected maps attribute names to Constraints.
type Expected map[string] Constraints

type ProvisionedThroughput struct {
	ReadCapacityUnits uint64
	WriteCapacityUnits uint64
}

type ProvisionedThroughputDesc struct {
	LastIncreaseDateTime float64
	LastDecreaseDateTime float64
	ReadCapacityUnits uint64
	WriteCapacityUnits uint64
	NumberOfDecreasesToday uint64
}

// Packages implementing the Endpoint interface should return the
// string output from the authorized request (or ""), the http code,
// and an error (or nil). This is the fundamental endpoint interface of
// GoDynamo.
type Endpoint interface {
	EndpointReq() (string,int,error)
}

type Endpoint_Response struct {
	Body string
	Code int
	Err  error
}

// ReqErr is a convenience function to see if the request was bad
func ReqErr(code int) bool {
	return code >= http.StatusBadRequest &&
		code < http.StatusInternalServerError
}

// ServerErr is a convenience function to see if the remote server had an internal error
func ServerErr(code int) bool {
	return code >= http.StatusInternalServerError
}

// HttpErr is a convenience function to see determine if the code is an error code.
func HttpErr(code int) bool {
	return ReqErr(code) || ServerErr(code)
}

// KeyDefinition is how table keys are described.
type KeyDefinition struct {
	AttributeName string
	KeyType string
}

type keyDefinition KeyDefinition

func (k KeyDefinition) MarshalJSON() ([]byte, error) {
	if !(k.KeyType == HASH || k.KeyType == RANGE) {
		e := fmt.Sprintf("endpoint.KeyDefinition.MarshalJSON: KeyType %s is not valid",
			k.KeyType)
		return nil, errors.New(e)
	}
	var ki keyDefinition
	ki.AttributeName = k.AttributeName
	ki.KeyType       = k.KeyType
	return json.Marshal(ki)
}

type KeySchema []KeyDefinition

type LocalSecondaryIndex struct {
	IndexName string
	KeySchema KeySchema
	Projection struct {
		NonKeyAttributes []string
		ProjectionType string
	}
}

func NewLocalSecondaryIndex() (*LocalSecondaryIndex) {
	l := new(LocalSecondaryIndex)
	l.KeySchema = make(KeySchema,0)
	l.Projection.NonKeyAttributes = make([]string,0)
	return l
}

type localSecondaryIndex LocalSecondaryIndex

type LocalSecondaryIndexes []LocalSecondaryIndex

func (l LocalSecondaryIndex) MarshalJSON() ([]byte, error) {
	if !(l.Projection.ProjectionType == ALL ||
		l.Projection.ProjectionType == KEYS_ONLY ||
		l.Projection.ProjectionType == INCLUDE) {
		e := fmt.Sprintf("endpoint.LocalSecondaryIndex.MarshalJSON: " +
			"ProjectionType %s is not valid",l.Projection.ProjectionType)
		return nil, errors.New(e)
	}
	if len(l.Projection.NonKeyAttributes) > 20 {
		e := fmt.Sprintf("endpoint.LocalSecondaryIndex.MarshalJSON: " +
			"NonKeyAttributes > 20")
		return nil, errors.New(e)
	}
	var li localSecondaryIndex
	li.IndexName = l.IndexName
	li.KeySchema = l.KeySchema
	li.Projection = l.Projection
	// if present, must have length between 1 and 20
	if len(l.Projection.NonKeyAttributes) == 0 {
		li.Projection.NonKeyAttributes = nil
	}
	li.Projection.ProjectionType = l.Projection.ProjectionType
	return json.Marshal(li)
}
