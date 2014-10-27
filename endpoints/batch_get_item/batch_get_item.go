// Support for the DynamoDB BatchGetItem endpoint.
// This package offers support for request sizes that exceed AWS limits.
//
// example use:
//
// tests/batch_get_item-livestest.go
//
package batch_get_item

import (
	"errors"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/types/attributestoget"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/item"
	"github.com/smugmug/godynamo/types/expressionattributenames"
	"github.com/smugmug/godynamo/types/capacity"
	ep "github.com/smugmug/godynamo/endpoint"
)

const (
	ENDPOINT_NAME     = "BatchGetItem"
	BATCHGET_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// actual limit is 1024kb
	QUERY_LIM_BYTES = 1048576
	QUERY_LIM       = 100
	RECURSE_LIM     = 50
)

// RequestInstance indicates what Keys to retrieve for a Table.
type RequestInstance struct {
	AttributesToGet attributestoget.AttributesToGet `json:",omitempty"`
	ConsistentRead bool `json:",omitempty"`
	ExpressionAttributeNames expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	Keys []item.Item
	ProjectionExpression string `json:",omitempty"`
}

func NewRequestInstance() (*RequestInstance) {
	r := new(RequestInstance)
	r.AttributesToGet = attributestoget.NewAttributesToGet()
	r.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	r.Keys = make([]item.Item,0)
	return r
}

// Table2Requests maps Table names to Key and Attribute data to retrieve.
type Table2Requests map[string] *RequestInstance

type BatchGetItem struct {
	RequestItems Table2Requests
	ReturnConsumedCapacity string `json:",omitempty"`
}

func NewBatchGetItem() (*BatchGetItem) {
	b := new(BatchGetItem)
	b.RequestItems = make(Table2Requests)
	return b
}

type Request BatchGetItem

type Response struct {
	ConsumedCapacity []capacity.ConsumedCapacity
	Responses map[string] []item.Item
	UnprocessedKeys Table2Requests
}

func NewResponse() (*Response) {
	r := new(Response)
	r.ConsumedCapacity      = make([]capacity.ConsumedCapacity,0)
	r.Responses             = make(map[string] []item.Item)
	r.UnprocessedKeys       = make(Table2Requests)
	return r
}

// Split supports the ability to have BatchGetItem structs whose size
// excceds the stated AWS limits. This function splits an arbitrarily-sized
// BatchGetItems into a list of BatchGetItem structs that are limited
// to the upper bound stated by AWS.
func Split(b *BatchGetItem) ([]BatchGetItem,error) {
	bs := make([]BatchGetItem,0)
	bi := NewBatchGetItem()
	i := 0
	for tn,_ := range b.RequestItems {
		for _,ri := range b.RequestItems[tn].Keys {
			if i == QUERY_LIM {
 				bi.ReturnConsumedCapacity = b.ReturnConsumedCapacity
 				bs = append(bs,*bi)
 				bi = NewBatchGetItem()
				i = 0
			}
			if _,tn_in_bi := bi.RequestItems[tn]; !tn_in_bi {
				bi.RequestItems[tn] = NewRequestInstance()
				bi.RequestItems[tn].AttributesToGet =
					make(attributestoget.AttributesToGet,
					len(b.RequestItems[tn].AttributesToGet))
				copy(bi.RequestItems[tn].AttributesToGet,
					b.RequestItems[tn].AttributesToGet)
				bi.RequestItems[tn].ConsistentRead = b.RequestItems[tn].ConsistentRead
			}
			bi.RequestItems[tn].Keys = append(bi.RequestItems[tn].Keys,ri)
			i++
		}
	}
	bi.ReturnConsumedCapacity = b.ReturnConsumedCapacity
	bs = append(bs,*bi)
	return bs,nil
}

func (batch_get_item *BatchGetItem) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	reqJSON,json_err := json.Marshal(batch_get_item);
	if json_err != nil {
		return "",0,json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON,BATCHGET_ENDPOINT)
}

func (req *Request) EndpointReq() (string,int,error) {
	batch_get_item := BatchGetItem(*req)
	return batch_get_item.EndpointReq()
}

// DoBatchGet is an endpoint request handler for BatchGetItem that supports arbitrarily-sized
// BatchGetItem struct instances. These are split in a list of conforming BatchGetItem instances
// via `Split` and the concurrently dispatched to DynamoDB, with the resulting responses stitched
// together. May break your provisioning.
func (b *BatchGetItem) DoBatchGet() (string,int,error) {
	var err error
	code := http.StatusOK
	body := ""
	bs,split_err := Split(b)
	if split_err != nil {
		e := fmt.Sprintf("batch_get_item.DoBatchGet: split failed: %s",split_err.Error())
		return body,code,errors.New(e)
	}
	resps := make(chan ep.Endpoint_Response,len(bs))
	for _,bi := range bs {
		go func(bi_ BatchGetItem) {
			body,code,err := bi_.RetryBatchGet(0)
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
				e := fmt.Sprintf("batch_get_item.DoBatchGet: %s",um_err.Error())
				err = errors.New(e)
			}
			// merge the responses from this call and the recursive one
			_ = combineResponseMetadata(combined_resp,&r)
			_ = combineResponses(combined_resp,&r)
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

// unprocessedKeys2BatchGetItems will take a response from DynamoDB that indicates some Keys
// require resubmitting, and turns these into a BatchGetItem struct instance.
func unprocessedKeys2BatchGetItems(req *BatchGetItem,resp *Response) (*BatchGetItem,error) {
	b := NewBatchGetItem()
	b.ReturnConsumedCapacity = req.ReturnConsumedCapacity
	for tn,_ := range resp.UnprocessedKeys {
		if _,tn_in_b := b.RequestItems[tn]; !tn_in_b {
			b.RequestItems[tn] = NewRequestInstance()
			b.RequestItems[tn].AttributesToGet = make(
				attributestoget.AttributesToGet,
				len(resp.UnprocessedKeys[tn].AttributesToGet))
			copy(b.RequestItems[tn].AttributesToGet,
				resp.UnprocessedKeys[tn].AttributesToGet)
			b.RequestItems[tn].ConsistentRead =
				resp.UnprocessedKeys[tn].ConsistentRead
			for _,item_src := range resp.UnprocessedKeys[tn].Keys {
				item_cp := item.NewItem()
				for k,v := range item_src {
					v_cp := attributevalue.NewAttributeValue()
					cp_err := v.Copy(v_cp)
					if cp_err != nil {
						return nil,cp_err
					}
					item_cp[k] = v_cp
				}
				b.RequestItems[tn].Keys = append(b.RequestItems[tn].Keys,item_cp)
			}
		}
	}
	return b,nil
}

// Add ConsumedCapacity from "this" Response to "all", the eventual stitched Response.
func combineResponseMetadata(all,this *Response) (error) {
	combinedConsumedCapacity := make([]capacity.ConsumedCapacity,0)
	for _,this_cc := range this.ConsumedCapacity {
		var cc capacity.ConsumedCapacity
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
	return nil
}

// Add actual response data from "this" Response to "all", the eventual stitched Response.
func combineResponses(all,this *Response) (error) {
	for tn,_ := range this.Responses {
		if _,tn_in_all := all.Responses[tn]; !tn_in_all {
			all.Responses[tn] = make([]item.Item,0)
		}
		for _,item_src := range this.Responses[tn] {
			item_cp := item.NewItem()
			for k,v := range item_src {
				v_cp := attributevalue.NewAttributeValue()
				cp_err := v.Copy(v_cp)
				if cp_err != nil {
					return cp_err
				}
				item_cp[k] = v_cp
			}
			all.Responses[tn] = append(all.Responses[tn],item_cp)
		}
	}
	return nil
}

// RetryBatchGet will attempt to fully complete a conforming BatchGetItem request.
// Callers for this method should be of len QUERY_LIM or less (see DoBatchGets()).
// This is different than EndpointReq in that it will extract UnprocessedKeys and
// form new BatchGetItem's based on those, and combine any results.
func (b *BatchGetItem) RetryBatchGet(depth int) (string,int,error) {
	if depth > RECURSE_LIM {
		e := fmt.Sprintf("batch_get_item.RetryBatchGet: recursion depth exceeded")
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
		e := fmt.Sprintf("batch_get_item.RetryBatchGet: %s",um_err.Error())
		return "",0,errors.New(e)
	}
	// if there are unprocessed items remaining from this call...
	if len(resp.UnprocessedKeys) > 0 {
		// make a new BatchGetItem object based on the unprocessed items
		n_req,n_req_err := unprocessedKeys2BatchGetItems(b,&resp)
		if n_req_err != nil {
			e := fmt.Sprintf("batch_get_item.RetryBatchGet: %s",n_req_err.Error())
			return "",0,errors.New(e)
		}
		// call this function on the new object
		n_body,n_code,n_err := n_req.RetryBatchGet(depth+1)
		if n_err != nil || n_code != http.StatusOK {
			return n_body,n_code,n_err
		}
		// get the response as an object
		var n_resp Response
		um_err := json.Unmarshal([]byte(n_body),&n_resp)
		if um_err != nil {
			e := fmt.Sprintf("batch_get_item.RetryBatchGet: %s",um_err.Error())
			return "",0,errors.New(e)
		}
		// merge the responses from this call and the recursive one
		_ = combineResponseMetadata(&resp,&n_resp)
		_ = combineResponses(&resp,&n_resp)
		// make a response string again out of the merged responses
		resp_json,resp_json_err := json.Marshal(resp)
		if resp_json_err != nil {
			e := fmt.Sprintf("batch_get_item.RetryBatchGet: %s",resp_json_err.Error())
			return "",0,errors.New(e)
		}
		body = string(resp_json)
	}
	return body,code,err
}
