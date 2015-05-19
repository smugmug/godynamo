// Support for the DynamoDB Query endpoint.
//
// example use:
//
// tests/query-livestest.go
//
package query

import (
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/types/attributestoget"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/aws_strings"
	"github.com/smugmug/godynamo/types/capacity"
	"github.com/smugmug/godynamo/types/condition"
	"github.com/smugmug/godynamo/types/expressionattributenames"
	"github.com/smugmug/godynamo/types/item"
)

const (
	ENDPOINT_NAME  = "Query"
	QUERY_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	OP_EQ          = aws_strings.OP_EQ
	OP_LE          = aws_strings.OP_LE
	OP_LT          = aws_strings.OP_LT
	OP_GE          = aws_strings.OP_GE
	OP_GT          = aws_strings.OP_GT
	OP_BEGINS_WITH = aws_strings.OP_BEGINS_WITH
	OP_BETWEEN     = aws_strings.OP_BETWEEN
	LIMIT          = 10000 // limit of query unless set
)

type ComparisonOperator string

// These are here for backward compatibility
type KeyConditions condition.Conditions
type KeyCondition condition.Condition

type Query struct {
	AttributesToGet           attributestoget.AttributesToGet                   `json:",omitempty"`
	ConditionalOperator       string                                            `json:",omitempty"`
	ConsistentRead            bool                                              // false is sane default
	ExclusiveStartKey         attributevalue.AttributeValueMap                  `json:",omitempty"`
	ExpressionAttributeNames  expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	ExpressionAttributeValues attributevalue.AttributeValueMap                  `json:",omitempty"`
	FilterExpression          string                                            `json:",omitempty"`
	IndexName                 string                                            `json:",omitempty"`
	KeyConditions             condition.Conditions
	Limit                     uint64               `json:",omitempty"`
	ProjectionExpression      string               `json:",omitempty"`
	QueryFilter               condition.Conditions `json:",omitempty"`
	ReturnConsumedCapacity    string               `json:",omitempty"`
	ScanIndexForward          *bool                `json:",omitempty"`
	Select                    string               `json:",omitempty"`
	TableName                 string
}

// NewQuery returns a pointer to an instantiation of the Query struct.
func NewQuery() *Query {
	q := new(Query)
	q.AttributesToGet = attributestoget.NewAttributesToGet()
	q.ExclusiveStartKey = attributevalue.NewAttributeValueMap()
	q.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	q.ExpressionAttributeValues = attributevalue.NewAttributeValueMap()
	q.KeyConditions = condition.NewConditions()
	q.QueryFilter = condition.NewConditions()
	return q
}

type Request Query

type Response struct {
	ConsumedCapacity *capacity.ConsumedCapacity `json:",omitempty"`
	Count            uint64
	Items            []item.Item                      `json:",omitempty"`
	LastEvaluatedKey attributevalue.AttributeValueMap `json:",omitempty"`
	ScannedCount     uint64                           `json:",omitempty"`
}

func NewResponse() *Response {
	r := new(Response)
	r.ConsumedCapacity = capacity.NewConsumedCapacity()
	r.Items = make([]item.Item, 0)
	r.LastEvaluatedKey = attributevalue.NewAttributeValueMap()
	return r
}

// These implementations of EndpointReq use a parameterized conf.

func (query *Query) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if query == nil {
		return nil, 0, errors.New("query.(Query)EndpointReqWithConf: receiver is nil")
	}
	if !conf.IsValid(c) {
		return nil, 0, errors.New("query.EndpointReqWithConf: c is not valid")
	}
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(query)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4WithConf(reqJSON, QUERY_ENDPOINT, c)
}

func (req *Request) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("query.(Request)EndpointReqWithConf: receiver is nil")
	}
	query := Query(*req)
	return query.EndpointReqWithConf(c)
}

// These implementations of EndpointReq use the global conf.

func (query *Query) EndpointReq() ([]byte, int, error) {
	if query == nil {
		return nil, 0, errors.New("query.(Query)EndpointReq: receiver is nil")
	}
	return query.EndpointReqWithConf(&conf.Vals)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("query.(Request)EndpointReq: receiver is nil")
	}
	query := Query(*req)
	return query.EndpointReqWithConf(&conf.Vals)
}

// ValidOp determines if an operation is in the approved list.
func ValidOp(op string) bool {
	return (op == OP_EQ ||
		op == OP_LE ||
		op == OP_LT ||
		op == OP_GE ||
		op == OP_GT ||
		op == OP_BEGINS_WITH ||
		op == OP_BETWEEN)
}
