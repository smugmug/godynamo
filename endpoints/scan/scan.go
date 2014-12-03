// Support for the DynamoDB Scan endpoint.
//
// example use:
//
// tests/scan-livestest.go
//
package scan

import (
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/types/attributestoget"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/aws_strings"
	"github.com/smugmug/godynamo/types/capacity"
	"github.com/smugmug/godynamo/types/condition"
	"github.com/smugmug/godynamo/types/expressionattributenames"
	"github.com/smugmug/godynamo/types/item"
)

const (
	ENDPOINT_NAME   = "Scan"
	SCAN_ENDPOINT   = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	OP_EQ           = aws_strings.OP_EQ
	OP_NE           = aws_strings.OP_NE
	OP_LE           = aws_strings.OP_LE
	OP_LT           = aws_strings.OP_LT
	OP_GE           = aws_strings.OP_GE
	OP_GT           = aws_strings.OP_GT
	OP_NULL         = aws_strings.OP_NULL
	OP_NOT_NULL     = aws_strings.OP_NOT_NULL
	OP_CONTAINS     = aws_strings.OP_CONTAINS
	OP_NOT_CONTAINS = aws_strings.OP_NOT_CONTAINS
	OP_BEGINS_WITH  = aws_strings.OP_BEGINS_WITH
	OP_IN           = aws_strings.OP_IN
	OP_BETWEEN      = aws_strings.OP_BETWEEN
	LIMIT           = 10000 // limit of scan unless set
)

type ComparisonOperator string

type Scan struct {
	AttributesToGet           attributestoget.AttributesToGet                   `json:",omitempty"`
	ConditionalOperator       string                                            `json:",omitempty"`
	ExclusiveStartKey         attributevalue.AttributeValueMap                  `json:",omitempty"`
	ExpressionAttributeNames  expressionattributenames.ExpressionAttributeNames `json:",omitempty"`
	ExpressionAttributeValues attributevalue.AttributeValueMap                  `json:",omitempty"`
	FilterExpression          string                                            `json:",omitempty"`
	Limit                     uint64                                            `json:",omitempty"`
	ProjectionExpression      string                                            `json:",omitempty"`
	ReturnConsumedCapacity    string                                            `json:",omitempty"`
	ScanFilter                condition.Conditions
	Segment                   uint64 `json:",omitempty"`
	Select                    string `json:",omitempty"`
	TableName                 string
	TotalSegments             uint64 `json:",omitempty"`
}

func NewScan() *Scan {
	s := new(Scan)
	s.AttributesToGet = attributestoget.NewAttributesToGet()
	s.ExclusiveStartKey = attributevalue.NewAttributeValueMap()
	s.ExpressionAttributeNames = expressionattributenames.NewExpressionAttributeNames()
	s.ExpressionAttributeValues = attributevalue.NewAttributeValueMap()
	s.ScanFilter = condition.NewConditions()
	return s
}

type Request Scan

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

func (scan *Scan) EndpointReq() ([]byte, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(scan)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, SCAN_ENDPOINT)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	scan := Scan(*req)
	return scan.EndpointReq()
}

// ValidOp determines if an operation is in the approved list.
func ValidOp(op string) bool {
	return (op == OP_EQ ||
		op == OP_NE ||
		op == OP_LE ||
		op == OP_LT ||
		op == OP_GE ||
		op == OP_GT ||
		op == OP_NULL ||
		op == OP_NOT_NULL ||
		op == OP_CONTAINS ||
		op == OP_NOT_CONTAINS ||
		op == OP_BEGINS_WITH ||
		op == OP_IN ||
		op == OP_BETWEEN)
}
