// Support for the DynamoDB UpdateItem endpoint.
//
// example use:
//
// tests/item_operations-livestest.go
//
package update_item

import (
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
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

// These implementations of EndpointReq use a parameterized conf.

func (update_item *UpdateItem) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if update_item == nil {
		return nil, 0, errors.New("update_item.(UpdateItem)EndpointReqWithConf: receiver is nil")
	}
	if !conf.IsValid(c) {
		return nil, 0, errors.New("update_item.EndpointReqWithConf: c is not valid")
	}
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(update_item)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4WithConf(reqJSON, UPDATEITEM_ENDPOINT, c)
}

func (update *Update) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if update == nil {
		return nil, 0, errors.New("update_item.(Update)EndpointReqWithConf: receiver is nil")
	}
	update_item := UpdateItem(*update)
	return update_item.EndpointReqWithConf(c)
}

func (req *Request) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("update_item.(Request)EndpointReqWithConf: receiver is nil")
	}
	update_item := UpdateItem(*req)
	return update_item.EndpointReqWithConf(c)
}

// These implementations of EndpointReq use the global conf.

func (update_item *UpdateItem) EndpointReq() ([]byte, int, error) {
	if update_item == nil {
		return nil, 0, errors.New("update_item.(UpdateItem)EndpointReq: receiver is nil")
	}
	return update_item.EndpointReqWithConf(&conf.Vals)
}

func (update *Update) EndpointReq() ([]byte, int, error) {
	if update == nil {
		return nil, 0, errors.New("update_item.(Update)EndpointReq: receiver is nil")
	}
	update_item := UpdateItem(*update)
	return update_item.EndpointReqWithConf(&conf.Vals)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("update_item.(Request)EndpointReq: receiver is nil")
	}
	update_item := UpdateItem(*req)
	return update_item.EndpointReqWithConf(&conf.Vals)
}
