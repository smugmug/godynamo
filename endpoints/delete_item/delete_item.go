// Support for the DynamoDB DeleteItem endpoint.
//
// example use:
//
// tests/delete_item-livestest.go
//
package delete_item

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
	ENDPOINT_NAME       = "DeleteItem"
	DELETEITEM_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// the permitted ReturnValues flags for this op
	RETVAL_ALL_OLD = aws_strings.RETVAL_ALL_OLD
	RETVAL_NONE    = aws_strings.RETVAL_NONE
)

type DeleteItem struct {
	ConditionExpression         string                                            `json:",omitempty"`
	ConditionalOperator         string                                            `json:",omitempty"`
	Expected                    expected.Expected                                 `json:",omitempty"`
	ExpressionAttributeNames    expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	ExpressionAttributeValues   attributevalue.AttributeValueMap                  `json:",omitempty"`
	Key                         item.Key
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	TableName                   string
}

func NewDeleteItem() *DeleteItem {
	u := new(DeleteItem)
	u.Expected = expected.NewExpected()
	u.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	u.ExpressionAttributeValues = attributevalue.NewAttributeValueMap()
	u.Key = item.NewKey()
	return u
}

// Delete is an alias for backwards compatibility
type Delete DeleteItem

func NewDelete() *Delete {
	delete_item := NewDeleteItem()
	delete := Delete(*delete_item)
	return &delete
}

type Request DeleteItem

type Response attributesresponse.AttributesResponse

func NewResponse() *Response {
	a := attributesresponse.NewAttributesResponse()
	r := Response(*a)
	return &r
}

func (delete_item *DeleteItem) EndpointReq() (string, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(delete_item)
	if json_err != nil {
		return "", 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, DELETEITEM_ENDPOINT)
}

func (delete *Delete) EndpointReq() (string, int, error) {
	delete_item := DeleteItem(*delete)
	return delete_item.EndpointReq()
}

func (req *Request) EndpointReq() (string, int, error) {
	delete_item := DeleteItem(*req)
	return delete_item.EndpointReq()
}
