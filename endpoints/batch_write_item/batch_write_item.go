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

// Support for the DynamoDB BatchWriteItem endpoint.
// This package offers support for request sizes that exceed AWS limits.
//
// example use:
//
// tests/batch_write_item-livestest.go
//
package batch_write_item

import (
	"errors"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME       = "BatchWriteItem"
	BATCHWRITE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// actual limit is 1024kb
	QUERY_LIM_BYTES = 1048576
	QUERY_LIM       = 25
	RECURSE_LIM     = 50
)

type DeleteRequest struct {
	Key ep.Item
}

type PutRequest struct {
	Item ep.Item
}

// BatchWriteItem requests can be puts or deletes. The non-nil member of this struct will be the request type specified.
// Do not specify (non-nil) both in one struct instance.
type RequestInstance struct {
	PutRequest *PutRequest
	DeleteRequest *DeleteRequest
}

type request_putrequest struct {
	PutRequest PutRequest
}

type request_deleterequest struct {
	DeleteRequest DeleteRequest
}

// REQUEST STRUCT TYPES

// Table2Requests maps Table names to list of RequestInstances
type Table2Requests map[string] []RequestInstance

type BatchWriteItem struct {
	RequestItems Table2Requests
	ReturnConsumedCapacity ep.ReturnConsumedCapacity
	ReturnItemCollectionMetrics ep.ReturnItemCollectionMetrics
}

// NewBatchWriteItem will return a pointer to an initialized BatchWriteItem struct.
func NewBatchWriteItem() (*BatchWriteItem) {
	b := new(BatchWriteItem)
	b.RequestItems = make(Table2Requests)
	return b
}

type Request BatchWriteItem

type Response struct {
	ConsumedCapacity []ep.ConsumedCapacity
	ItemCollectionMetrics map [string] []ep.ItemCollectionMetrics
	UnprocessedItems Table2Requests
}

// NewResponse will return a pointer to an initialized Response struct.
func NewResponse() (*Response) {
	r := new(Response)
	r.ConsumedCapacity      = make([]ep.ConsumedCapacity,0)
	r.ItemCollectionMetrics = make(map[string] []ep.ItemCollectionMetrics)
	r.UnprocessedItems      = make(Table2Requests)
	return r
}

func (r RequestInstance) MarshalJSON() ([]byte, error) {
	if r.PutRequest != nil {
		var p request_putrequest
		p.PutRequest = *r.PutRequest
		return json.Marshal(p)
	}
	if r.DeleteRequest != nil {
		var d request_deleterequest
		d.DeleteRequest = *r.DeleteRequest
		return json.Marshal(d)
	}
	e := fmt.Sprintf("batch_write_item.Request.MarshalJSON:" +
		"no valid puts or deletes in %v",r)
	return nil, errors.New(e)
}

// Split supports the ability to have BatchWriteItem structs whose size
// excceds the stated AWS limits. This function splits an arbitrarily-sized
// BatchWriteItems into a list of BatchWriteItem structs that are limited
// to the upper bound stated by AWS.
func Split(b BatchWriteItem) ([]BatchWriteItem,error) {
	bs := make([]BatchWriteItem,0)
	bi := NewBatchWriteItem()
	i := 0
	// for each table name (tn) in b.RequestItems
	for tn,_ := range b.RequestItems {
		// for each request in that table's list
		for _,ri := range b.RequestItems[tn] {
			if i == QUERY_LIM {
				// append value of existing bi, make a new one
				bi.ReturnConsumedCapacity = b.ReturnConsumedCapacity
				bi.ReturnItemCollectionMetrics = b.ReturnItemCollectionMetrics
				bs = append(bs,*bi)
				bi = NewBatchWriteItem()
				i = 0
			}
			// if creating a request in bi for tn for the first time, initialize
			if _,tn_in_bi := bi.RequestItems[tn]; !tn_in_bi {
				bi.RequestItems[tn] = make([]RequestInstance,0)
			}
			// append request to list in bi for this tn
			bi.RequestItems[tn] = append(bi.RequestItems[tn],ri)
			i++
		}
	}
	bi.ReturnConsumedCapacity = b.ReturnConsumedCapacity
	bi.ReturnItemCollectionMetrics = b.ReturnItemCollectionMetrics
	bs = append(bs,*bi)
	return bs,nil
}

// EndpointReq for BatchGetItem which assumes its BatchWriteItem struct instance `b`
// conforms to AWS limits. Use this if you do not employ arbitrarily-sized BatchWriteItems
// and instead choose to conform to the AWS limits.
func (b BatchWriteItem) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	if authreq.AUTH_VERSION != authreq.AUTH_V4 {
		e := fmt.Sprintf("batch_write_item(BatchWriteItem).EndpointReq " +
			"auth must be v4")
		return "",0,errors.New(e)
	}
	return authreq.RetryReq_V4(&b,BATCHWRITE_ENDPOINT)
}

func (req Request) EndpointReq() (string,int,error) {
	return (BatchWriteItem(req)).EndpointReq()
}

// DoBatchWrite is an endpoint request handler for BatchWriteItem that supports arbitrarily-sized
// BatchWriteItem struct instances. These are split in a list of conforming BatchWriteItem instances
// via `Split` and the concurrently dispatched to DynamoDB, with the resulting responses stitched
// together. May break your provisioning.
func (b BatchWriteItem) DoBatchWrite() (string,int,error) {
	var err error
	code := http.StatusOK
	body := ""
	bs,split_err := Split(b)
	if split_err != nil {
		e := fmt.Sprintf("batch_write_item.DoBatchWrite: split failed: %s",split_err.Error())
		return body,code,errors.New(e)
	}
	resps := make(chan ep.Endpoint_Response,len(bs))
	for _,bi := range bs {
		go func(bi_ BatchWriteItem) {
			body,code,err := bi_.RetryBatchWrite(0)
			resps <- ep.Endpoint_Response{Body:body,Code:code,Err:err}
		}(bi)
	}
	combined_resp := NewResponse()
	for i := 0; i < len(bs); i++ {
		resp := <- resps
		if resp.Err != nil {
			err = resp.Err
		} else if resp.Code != http.StatusOK {
			code = resp.Code
		} else {
			var r Response
			um_err := json.Unmarshal([]byte(resp.Body),&r)
			if um_err != nil {
				e := fmt.Sprintf("batch_write_item.DoBatchWrite: %s",um_err.Error())
				err = errors.New(e)
			}
			// merge the responses from this call and the recursive one
			_ = combineResponseMetadata(combined_resp,&r)
		}
	}
	body_bytes,marshal_err := json.Marshal(*combined_resp)
	if marshal_err != nil {
		err = marshal_err
	} else {
		body = string(body_bytes)
	}
	return body,code,err
}

// unprocessedKeys2BatchWriteItems will take a response from DynamoDB that indicates some Keys
// require resubmitting, and turns these into a BatchWriteItem struct instance.
func unprocessedItems2BatchWriteItems(req BatchWriteItem,resp *Response) (*BatchWriteItem,error) {
	b := NewBatchWriteItem()
	for tn,_ := range resp.UnprocessedItems {
		for _,reqinst := range resp.UnprocessedItems[tn] {
			var reqinst_cp RequestInstance
			if reqinst.DeleteRequest != nil {
				reqinst_cp.DeleteRequest = new(DeleteRequest)
				reqinst_cp.DeleteRequest.Key = make(ep.Item)
				for k,v := range reqinst.DeleteRequest.Key {
					reqinst_cp.DeleteRequest.Key[k] = v
				}
				b.RequestItems[tn] = append(b.RequestItems[tn],reqinst_cp)
			} else if reqinst.PutRequest != nil {
				reqinst_cp.PutRequest = new(PutRequest)
				reqinst_cp.PutRequest.Item = make(ep.Item)
				for k,v := range reqinst.PutRequest.Item {
					reqinst_cp.PutRequest.Item[k] = v
				}
				b.RequestItems[tn] = append(b.RequestItems[tn],reqinst_cp)
			}
		}
	}
	b.ReturnConsumedCapacity      = req.ReturnConsumedCapacity
	b.ReturnItemCollectionMetrics = req.ReturnItemCollectionMetrics
	return b,nil
}

// Add ConsumedCapacity from "this" Response to "all", the eventual stitched Response.
func combineResponseMetadata(all,this *Response) (error) {
	combinedConsumedCapacity := make([]ep.ConsumedCapacity,0)
	for _,this_cc := range this.ConsumedCapacity {
		var cc ep.ConsumedCapacity
		cc.TableName = this_cc.TableName
		cc.CapacityUnits = this_cc.CapacityUnits
		for _,all_cc := range all.ConsumedCapacity {
			if all_cc.TableName == this_cc.TableName {
				cc.CapacityUnits += all_cc.CapacityUnits
			}
		}
		combinedConsumedCapacity = append(combinedConsumedCapacity,cc)
	}
	all.ConsumedCapacity = combinedConsumedCapacity
	for tn,_ := range this.ItemCollectionMetrics {
		for _,icm := range this.ItemCollectionMetrics[tn] {
			if _,tn_is_all := all.ItemCollectionMetrics[tn]; !tn_is_all {
				all.ItemCollectionMetrics[tn] =
					make([]ep.ItemCollectionMetrics,0)
			}
			all.ItemCollectionMetrics[tn] = append(all.ItemCollectionMetrics[tn],icm)
		}
	}
	return nil
}

// RetryBatchWrite will attempt to fully complete a conforming BatchWriteItem request.
// Callers for this method should be of len QUERY_LIM or less (see DoBatchWrites()).
// This is different than EndpointReq in that it will extract UnprocessedKeys and
// form new BatchWriteItem's based on those, and combine any results.
func (b BatchWriteItem) RetryBatchWrite(depth int) (string,int,error) {
	if depth > RECURSE_LIM {
		e := fmt.Sprintf("batch_write_item.RetryBatchWrite: recursion depth exceeded")
		return "",0,errors.New(e)
	}
	body,code,err := b.EndpointReq()
	if err != nil || code != http.StatusOK {
		return body,code,err
	}
	// we'll need an actual Response object
	var resp Response
	um_err := json.Unmarshal([]byte(body),&resp)
	if um_err != nil {
		e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s",um_err.Error())
		return "",0,errors.New(e)
	}
	// if there are unprocessed items remaining from this call...
	if len(resp.UnprocessedItems) > 0 {
		// make a new BatchWriteItem object based on the unprocessed items
		n_req,n_req_err := unprocessedItems2BatchWriteItems(b,&resp)
		if n_req_err != nil {
			e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s",n_req_err.Error())
			return "",0,errors.New(e)
		}
		// call this function on the new object
		n_body,n_code,n_err := n_req.RetryBatchWrite(depth+1)
		if n_err != nil || n_code != http.StatusOK {
			return n_body,n_code,n_err
		}
		// get the response as an object
		var n_resp Response
		um_err := json.Unmarshal([]byte(n_body),&n_resp)
		if um_err != nil {
			e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s",um_err.Error())
			return "",0,errors.New(e)
		}
		// merge the responses from this call and the recursive one
		_ = combineResponseMetadata(&resp,&n_resp)
		// make a response string again out of the merged responses
		resp_json,resp_json_err := json.Marshal(resp)
		if resp_json_err != nil {
			e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s",resp_json_err.Error())
			return "",0,errors.New(e)
		}
		body = string(resp_json)
	}
	return body,code,err
}
