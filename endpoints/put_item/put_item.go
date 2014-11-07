// Support for the DynamoDB PutItem endpoint.
//
// example use:
//
// tests/put_item-livestest.go
//
package put_item

import (
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/types/attributesresponse"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/aws_strings"
	"github.com/smugmug/godynamo/types/expected"
	"github.com/smugmug/godynamo/types/expressionattributenames"
	"github.com/smugmug/godynamo/types/item"
)

const (
	ENDPOINT_NAME    = "PutItem"
	PUTITEM_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// the permitted ReturnValues flags for this op
	RETVAL_ALL_OLD = aws_strings.RETVAL_ALL_OLD
	RETVAL_NONE    = aws_strings.RETVAL_NONE
)

type PutItem struct {
	ConditionExpression         string                                            `json:",omitempty"`
	ConditionalOperator         string                                            `json:",omitempty"`
	Expected                    expected.Expected                                 `json:",omitempty"`
	ExpressionAttributeNames    expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	ExpressionAttributeValues   attributevalue.AttributeValueMap                  `json:",omitempty"`
	Item                        item.Item
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	TableName                   string
}

// NewPut will return a pointer to an initialized Put struct.
func NewPutItem() *PutItem {
	p := new(PutItem)
	p.Expected = expected.NewExpected()
	p.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	p.ExpressionAttributeValues = attributevalue.NewAttributeValueMap()
	p.Item = item.NewItem()
	return p
}

type Request PutItem

// Put is an alias for backwards compatibility
type Put PutItem

func NewPut() *Put {
	put_item := NewPutItem()
	put := Put(*put_item)
	return &put
}

type Response attributesresponse.AttributesResponse

func NewResponse() *Response {
	a := attributesresponse.NewAttributesResponse()
	r := Response(*a)
	return &r
}

func (put_item *PutItem) EndpointReq() (string, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(put_item)
	if json_err != nil {
		return "", 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, PUTITEM_ENDPOINT)
}

func (put *Put) EndpointReq() (string, int, error) {
	put_item := PutItem(*put)
	return put_item.EndpointReq()
}

func (req *Request) EndpointReq() (string, int, error) {
	put_item := PutItem(*req)
	return put_item.EndpointReq()
}

// ValidItem validates the size of a json serialization of an Item.
// AWS says items can only be 400k bytes binary
func ValidItem(i string) bool {
	return !(len([]byte(i)) > 400000)
}
