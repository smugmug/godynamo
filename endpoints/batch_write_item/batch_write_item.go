// Support for the DynamoDB BatchWriteItem endpoint.
// This package offers support for request sizes that exceed AWS limits.
//
// example use:
//
// tests/batch_write_item-livestest.go
//
package batch_write_item

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/capacity"
	"github.com/smugmug/godynamo/types/item"
	"github.com/smugmug/godynamo/types/itemcollectionmetrics"
	"net/http"
)

const (
	ENDPOINT_NAME       = "BatchWriteItem"
	JSON_ENDPOINT_NAME  = ENDPOINT_NAME + "JSON"
	BATCHWRITE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// actual limit is 1024kb
	QUERY_LIM_BYTES = 1048576
	QUERY_LIM       = 25
	RECURSE_LIM     = 50
)

type DeleteRequest struct {
	Key item.Item
}

type PutRequest struct {
	Item item.Item
}

type PutRequestItemJSON struct {
	Item interface{}
}

// BatchWriteItem requests can be puts or deletes.
// The non-nil member of this struct will be the request type specified.
// Do not specify (non-nil) both in one struct instance.
type RequestInstance struct {
	PutRequest    *PutRequest
	DeleteRequest *DeleteRequest
}

// Similar, but supporting the use of basic json for put requests. Note that
// use of basic json is only supported for Items, whereas delete requests
// use keys.
type RequestInstanceItemJSON struct {
	PutRequest    *PutRequestItemJSON
	DeleteRequest *DeleteRequest
}

type Table2Requests map[string][]RequestInstance

type Table2RequestsItemsJSON map[string][]RequestInstanceItemJSON

type BatchWriteItem struct {
	RequestItems                Table2Requests
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
}

func NewBatchWriteItem() *BatchWriteItem {
	b := new(BatchWriteItem)
	b.RequestItems = make(Table2Requests)
	return b
}

type Request BatchWriteItem

type BatchWriteItemJSON struct {
	RequestItems                Table2RequestsItemsJSON
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
}

func NewBatchWriteItemJSON() *BatchWriteItemJSON {
	b := new(BatchWriteItemJSON)
	b.RequestItems = make(Table2RequestsItemsJSON)
	return b
}

// ToBatchWriteItem will attempt to convert a BatchWriteItemJSON to BatchWriteItem
func (bwij *BatchWriteItemJSON) ToBatchWriteItem() (*BatchWriteItem, error) {
	if bwij == nil {
		return nil, errors.New("receiver is nil")
	}
	b := NewBatchWriteItem()
	for tn, ris := range bwij.RequestItems {
		l := len(ris)
		b.RequestItems[tn] = make([]RequestInstance, l)
		for i, ri := range ris {
			if ri.DeleteRequest != nil {
				b.RequestItems[tn][i].DeleteRequest = ri.DeleteRequest
				b.RequestItems[tn][i].PutRequest = nil
			} else if ri.PutRequest != nil {
				a, cerr := attributevalue.InterfaceToAttributeValueMap(ri.PutRequest.Item)
				if cerr != nil {
					return nil, cerr
				}
				b.RequestItems[tn][i].PutRequest = &PutRequest{Item: item.Item(a)}
				b.RequestItems[tn][i].DeleteRequest = nil
			} else {
				return nil, errors.New("no Put or Delete request found")
			}
		}
	}
	b.ReturnConsumedCapacity = bwij.ReturnConsumedCapacity
	b.ReturnItemCollectionMetrics = bwij.ReturnItemCollectionMetrics
	return b, nil
}

type Response struct {
	ConsumedCapacity      []capacity.ConsumedCapacity
	ItemCollectionMetrics itemcollectionmetrics.ItemCollectionMetricsMap
	UnprocessedItems      Table2Requests
}

func NewResponse() *Response {
	r := new(Response)
	r.ConsumedCapacity = make([]capacity.ConsumedCapacity, 0)
	r.ItemCollectionMetrics = itemcollectionmetrics.NewItemCollectionMetricsMap()
	r.UnprocessedItems = make(Table2Requests)
	return r
}

// Split supports the ability to have BatchWriteItem structs whose size
// excceds the stated AWS limits. This function splits an arbitrarily-sized
// BatchWriteItems into a list of BatchWriteItem structs that are limited
// to the upper bound stated by AWS.
func Split(b *BatchWriteItem) ([]BatchWriteItem, error) {
	bs := make([]BatchWriteItem, 0)
	bi := NewBatchWriteItem()
	i := 0
	// for each table name (tn) in b.RequestItems
	for tn, _ := range b.RequestItems {
		// for each request in that table's list
		for _, ri := range b.RequestItems[tn] {
			if i == QUERY_LIM {
				// append value of existing bi, make a new one
				bi.ReturnConsumedCapacity = b.ReturnConsumedCapacity
				bi.ReturnItemCollectionMetrics = b.ReturnItemCollectionMetrics
				bs = append(bs, *bi)
				bi = NewBatchWriteItem()
				i = 0
			}
			// if creating a request in bi for tn for the first time, initialize
			if _, tn_in_bi := bi.RequestItems[tn]; !tn_in_bi {
				bi.RequestItems[tn] = make([]RequestInstance, 0)
			}
			// append request to list in bi for this tn
			bi.RequestItems[tn] = append(bi.RequestItems[tn], ri)
			i++
		}
	}
	bi.ReturnConsumedCapacity = b.ReturnConsumedCapacity
	bi.ReturnItemCollectionMetrics = b.ReturnItemCollectionMetrics
	bs = append(bs, *bi)
	return bs, nil
}

func (batch_write_item *BatchWriteItem) EndpointReq() (string, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(batch_write_item)
	if json_err != nil {
		return "", 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, BATCHWRITE_ENDPOINT)
}

func (req *Request) EndpointReq() (string, int, error) {
	batch_write_item := BatchWriteItem(*req)
	return batch_write_item.EndpointReq()
}

// DoBatchWrite is an endpoint request handler for BatchWriteItem that supports arbitrarily-sized
// BatchWriteItem struct instances. These are split in a list of conforming BatchWriteItem instances
// via `Split` and the concurrently dispatched to DynamoDB, with the resulting responses stitched
// together. May break your provisioning.
func (b *BatchWriteItem) DoBatchWrite() (string, int, error) {
	var err error
	code := http.StatusOK
	body := ""
	bs, split_err := Split(b)
	if split_err != nil {
		e := fmt.Sprintf("batch_write_item.DoBatchWrite: split failed: %s", split_err.Error())
		return body, code, errors.New(e)
	}
	resps := make(chan ep.Endpoint_Response, len(bs))
	for _, bi := range bs {
		go func(bi_ BatchWriteItem) {
			body, code, err := bi_.RetryBatchWrite(0)
			resps <- ep.Endpoint_Response{Body: body, Code: code, Err: err}
		}(bi)
	}
	combined_resp := NewResponse()
	for i := 0; i < len(bs); i++ {
		resp := <-resps
		if resp.Err != nil {
			err = resp.Err
		} else if resp.Code != http.StatusOK {
			code = resp.Code
		} else {
			var r Response
			um_err := json.Unmarshal([]byte(resp.Body), &r)
			if um_err != nil {
				e := fmt.Sprintf("batch_write_item.DoBatchWrite: %s", um_err.Error())
				err = errors.New(e)
			}
			// merge the responses from this call and the recursive one
			_ = combineResponseMetadata(combined_resp, &r)
		}
	}
	body_bytes, marshal_err := json.Marshal(*combined_resp)
	if marshal_err != nil {
		err = marshal_err
	} else {
		body = string(body_bytes)
	}
	return body, code, err
}

// unprocessedKeys2BatchWriteItems will take a response from DynamoDB that indicates some Keys
// require resubmitting, and turns these into a BatchWriteItem struct instance.
func unprocessedItems2BatchWriteItems(req *BatchWriteItem, resp *Response) (*BatchWriteItem, error) {
	b := NewBatchWriteItem()
	for tn, _ := range resp.UnprocessedItems {
		for _, reqinst := range resp.UnprocessedItems[tn] {
			var reqinst_cp RequestInstance
			if reqinst.DeleteRequest != nil {
				reqinst_cp.DeleteRequest = new(DeleteRequest)
				reqinst_cp.DeleteRequest.Key = make(item.Item)
				for k, v := range reqinst.DeleteRequest.Key {
					v_cp := attributevalue.NewAttributeValue()
					cp_err := v.Copy(v_cp)
					if cp_err != nil {
						return nil, cp_err
					}
					reqinst_cp.DeleteRequest.Key[k] = v_cp
				}
				b.RequestItems[tn] = append(b.RequestItems[tn], reqinst_cp)
			} else if reqinst.PutRequest != nil {
				reqinst_cp.PutRequest = new(PutRequest)
				reqinst_cp.PutRequest.Item = make(item.Item)
				for k, v := range reqinst.PutRequest.Item {
					v_cp := attributevalue.NewAttributeValue()
					cp_err := v.Copy(v_cp)
					if cp_err != nil {
						return nil, cp_err
					}
					reqinst_cp.PutRequest.Item[k] = v_cp
				}
				b.RequestItems[tn] = append(b.RequestItems[tn], reqinst_cp)
			}
		}
	}
	b.ReturnConsumedCapacity = req.ReturnConsumedCapacity
	b.ReturnItemCollectionMetrics = req.ReturnItemCollectionMetrics
	return b, nil
}

// Add ConsumedCapacity from "this" Response to "all", the eventual stitched Response.
func combineResponseMetadata(all, this *Response) error {
	combinedConsumedCapacity := make([]capacity.ConsumedCapacity, 0)
	for _, this_cc := range this.ConsumedCapacity {
		var cc capacity.ConsumedCapacity
		cc.TableName = this_cc.TableName
		cc.CapacityUnits = this_cc.CapacityUnits
		for _, all_cc := range all.ConsumedCapacity {
			if all_cc.TableName == this_cc.TableName {
				cc.CapacityUnits += all_cc.CapacityUnits
			}
		}
		combinedConsumedCapacity = append(combinedConsumedCapacity, cc)
	}
	all.ConsumedCapacity = combinedConsumedCapacity
	for tn, _ := range this.ItemCollectionMetrics {
		for _, icm := range this.ItemCollectionMetrics[tn] {
			if _, tn_is_all := all.ItemCollectionMetrics[tn]; !tn_is_all {
				all.ItemCollectionMetrics[tn] =
					make([]*itemcollectionmetrics.ItemCollectionMetrics, 0)
			}
			all.ItemCollectionMetrics[tn] = append(all.ItemCollectionMetrics[tn], icm)
		}
	}
	return nil
}

// RetryBatchWrite will attempt to fully complete a conforming BatchWriteItem request.
// Callers for this method should be of len QUERY_LIM or less (see DoBatchWrites()).
// This is different than EndpointReq in that it will extract UnprocessedKeys and
// form new BatchWriteItem's based on those, and combine any results.
func (b *BatchWriteItem) RetryBatchWrite(depth int) (string, int, error) {
	if depth > RECURSE_LIM {
		e := fmt.Sprintf("batch_write_item.RetryBatchWrite: recursion depth exceeded")
		return "", 0, errors.New(e)
	}
	body, code, err := b.EndpointReq()
	if err != nil || code != http.StatusOK {
		return body, code, err
	}
	// we'll need an actual Response object
	var resp Response
	um_err := json.Unmarshal([]byte(body), &resp)
	if um_err != nil {
		e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s", um_err.Error())
		return "", 0, errors.New(e)
	}
	// if there are unprocessed items remaining from this call...
	if len(resp.UnprocessedItems) > 0 {
		// make a new BatchWriteItem object based on the unprocessed items
		n_req, n_req_err := unprocessedItems2BatchWriteItems(b, &resp)
		if n_req_err != nil {
			e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s", n_req_err.Error())
			return "", 0, errors.New(e)
		}
		// call this function on the new object
		n_body, n_code, n_err := n_req.RetryBatchWrite(depth + 1)
		if n_err != nil || n_code != http.StatusOK {
			return n_body, n_code, n_err
		}
		// get the response as an object
		var n_resp Response
		um_err := json.Unmarshal([]byte(n_body), &n_resp)
		if um_err != nil {
			e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s", um_err.Error())
			return "", 0, errors.New(e)
		}
		// merge the responses from this call and the recursive one
		_ = combineResponseMetadata(&resp, &n_resp)
		// make a response string again out of the merged responses
		resp_json, resp_json_err := json.Marshal(resp)
		if resp_json_err != nil {
			e := fmt.Sprintf("batch_write_item.RetryBatchWrite: %s", resp_json_err.Error())
			return "", 0, errors.New(e)
		}
		body = string(resp_json)
	}
	return body, code, err
}
