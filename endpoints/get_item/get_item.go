// Support for the DynamoDB GetItem endpoint.
//
// example use:
//
// tests/get_item-livestest.go
//
package get_item

import (
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/types/attributestoget"
	"github.com/smugmug/godynamo/types/capacity"
	"github.com/smugmug/godynamo/types/expressionattributenames"
	"github.com/smugmug/godynamo/types/item"
)

const (
	ENDPOINT_NAME    = "GetItem"
	GETITEM_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
)

type GetItem struct {
	AttributesToGet          attributestoget.AttributesToGet                   `json:",omitempty"`
	ConsistentRead           bool                                              // false is sane default
	ExpressionAttributeNames expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	Key                      item.Key
	ProjectionExpression     string `json:",omitempty"`
	ReturnConsumedCapacity   string `json:",omitempty"`
	TableName                string
}

func NewGetItem() *GetItem {
	g := new(GetItem)
	g.Key = item.NewKey()
	g.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	g.AttributesToGet = make(attributestoget.AttributesToGet, 0)
	return g
}

// Get is an alias for backwards compatibility
type Get GetItem

func NewGet() *Get {
	get_item := NewGetItem()
	get := Get(*get_item)
	return &get
}

type Request GetItem

type Response struct {
	Item             item.Item
	ConsumedCapacity *capacity.ConsumedCapacity `json:",omitempty"`
}

func NewResponse() *Response {
	r := new(Response)
	r.Item = item.NewItem()
	r.ConsumedCapacity = capacity.NewConsumedCapacity()
	return r
}

func (get_item *GetItem) EndpointReq() (string, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(get_item)
	if json_err != nil {
		return "", 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, GETITEM_ENDPOINT)
}

func (get *Get) EndpointReq() (string, int, error) {
	get_item := GetItem(*get)
	return get_item.EndpointReq()
}

func (req *Request) EndpointReq() (string, int, error) {
	get_item := GetItem(*req)
	return get_item.EndpointReq()
}
