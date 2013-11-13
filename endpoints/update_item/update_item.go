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

// Support for the DynamoDB UpdateItem endpoint.
//
// example use:
//
// tests/item_operations-livestest.go
//
package update_item

import (
	"fmt"
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME           = "UpdateItem"
	UPDATEITEM_ENDPOINT     = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// the permitted Action flags for this op
	ACTION_PUT		= "PUT"
	ACTION_DEL		= "DELETE"
	ACTION_ADD		= "ADD"
	// the permitted ReturnValues flags for this op
	RETVAL_ALL_OLD		= "ALL_OLD"
	RETVAL_ALL_NEW		= "ALL_NEW"
	RETVAL_UPDATED_OLD	= "UPDATED_OLD"
	RETVAL_UPDATED_NEW	= "UPDATED_NEW"
	RETVAL_NONE             = ep.RETVAL_NONE
)

type AttributeAction struct {
	Value ep.AttributeValue
	Action string
}

type attributeAction AttributeAction

type attributeAction_delete struct {
	Action string
}

func (a AttributeAction) MarshalJSON() ([]byte,error) {
	// deletes can have empty values - remove the attribute
	if a.Action == ACTION_DEL && a.Value.Empty() {
		var ad attributeAction_delete
		ad.Action = ACTION_DEL
		return json.Marshal(ad)
	}
	var ai attributeAction
	ai.Value = a.Value
	if a.Action == "" {
		ai.Action = ACTION_PUT
	} else {
		ai.Action = a.Action
	}
	return json.Marshal(ai)
}

func (a *AttributeAction) UnmarshalJSON (data []byte) error {
	var ai attributeAction
	t_err := json.Unmarshal(data,&ai)
	if t_err != nil {
		return t_err
	}
	if ai.Action == "" {
		a.Action = ACTION_PUT
	} else {
		a.Action = ai.Action
	}
	a.Value = ai.Value
	return nil
}

type AttributeUpdates map[string] AttributeAction

type Update struct {
	TableName string
	Key ep.Item
	AttributeUpdates AttributeUpdates
	Expected ep.Expected
	ReturnValues ep.ReturnValues
	ReturnItemCollectionMetrics ep.ReturnItemCollectionMetrics
}

// NewUpdate returns a pointer to an instantiation of the Update struct.
func NewUpdate() (*Update) {
	u := new(Update)
	u.Key			= make(ep.Item)
	u.AttributeUpdates	= make(AttributeUpdates)
	u.Expected		= make(ep.Expected)
	return u
}

type update Update

type Request Update

type Response ep.AttributesResponse

// NewResponse will return a pointer to an initialized Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.Attributes = make(ep.Item)
	return r
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(Update(r))
}

func (u Update) MarshalJSON() ([]byte, error) {
	var ui update
	if u.Expected == nil || len(u.Expected) == 0 {
		ui.Expected = nil
	} else {
		ui.Expected = u.Expected
	}
	ui.TableName = u.TableName
	ui.Key = u.Key
	ui.AttributeUpdates = u.AttributeUpdates
	ui.ReturnValues = u.ReturnValues
	return json.Marshal(ui)
}

// EndpointReq implements the Endpoint interface.
func (u Update) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("update_item(Update).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&u,UPDATEITEM_ENDPOINT)
}

// EndpointReq implements the Endpoint interface on the local Request type.
func (req Request) EndpointReq() (string,int,error) {
	return (Update(req)).EndpointReq()
}
