// Support for the DynamoDB UpdateItem endpoint.
//
// example use:
//
// tests/item_operations-livestest.go
//
package update_item

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
	ENDPOINT_NAME       = "UpdateItem"
	UPDATEITEM_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	// the permitted Action flags for this op
	ACTION_PUT = aws_strings.ACTION_PUT
	ACTION_DEL = aws_strings.ACTION_DEL
	ACTION_ADD = aws_strings.ACTION_ADD
	// the permitted ReturnValues flags for this op
	RETVAL_ALL_OLD     = aws_strings.RETVAL_ALL_OLD
	RETVAL_ALL_NEW     = aws_strings.RETVAL_ALL_NEW
	RETVAL_UPDATED_OLD = aws_strings.RETVAL_UPDATED_OLD
	RETVAL_UPDATED_NEW = aws_strings.RETVAL_UPDATED_NEW
	RETVAL_NONE        = aws_strings.RETVAL_NONE
)

type AttributeUpdates attributevalue.AttributeValueMap

type UpdateItem struct {
	AttributeUpdates            attributevalue.AttributeValueUpdateMap            `json:",omitempty"`
	ConditionExpression         string                                            `json:",omitempty"`
	ConditionalOperator         string                                            `json:",omitempty"`
	Expected                    expected.Expected                                 `json:",omitempty"`
	ExpressionAttributeNames    expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	ExpressionAttributeValues   attributevalue.AttributeValueMap                  `json:",omitempty"`
	Key                         item.Key
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	TableName                   string
	UpdateExpression            string `json:",omitempty"`
}

// NewUpdate returns a pointer to an instantiation of the Update struct.
func NewUpdateItem() *UpdateItem {
	u := new(UpdateItem)
	u.AttributeUpdates = attributevalue.NewAttributeValueUpdateMap()
	u.Expected = expected.NewExpected()
	u.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	u.ExpressionAttributeValues = attributevalue.NewAttributeValueMap()
	u.Key = item.NewKey()
	return u
}

// Update is an alias for backwards compatibility
type Update UpdateItem

func NewUpdate() *Update {
	update_item := NewUpdateItem()
	update := Update(*update_item)
	return &update
}

type Request UpdateItem

type Response attributesresponse.AttributesResponse

func NewResponse() *Response {
	a := attributesresponse.NewAttributesResponse()
	r := Response(*a)
	return &r
}

func (update_item *UpdateItem) EndpointReq() (string, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(update_item)
	if json_err != nil {
		return "", 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, UPDATEITEM_ENDPOINT)
}

func (update *Update) EndpointReq() (string, int, error) {
	update_item := UpdateItem(*update)
	return update_item.EndpointReq()
}

func (req *Request) EndpointReq() (string, int, error) {
	update_item := UpdateItem(*req)
	return update_item.EndpointReq()
}
