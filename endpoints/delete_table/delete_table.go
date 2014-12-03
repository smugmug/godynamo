// Support for the DynamoDB DeleteTable endpoint.
//
// example use:
//
// tests/delete_table-livestest.go
//
package delete_table

import (
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	create_table "github.com/smugmug/godynamo/endpoints/create_table"
)

const (
	ENDPOINT_NAME        = "DeleteTable"
	DELETETABLE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
)

type DeleteTable struct {
	TableName string
}

// Delete is an alias for backwards compatibility
type Delete DeleteTable

type Request DeleteTable

func NewDeleteTable() *DeleteTable {
	d := new(DeleteTable)
	return d
}

// DeleteTable and CreateTable use the same Response format
type Response create_table.Response

func NewResponse() *Response {
	cr := create_table.NewResponse()
	r := Response(*cr)
	return &r
}

func (delete_table *DeleteTable) EndpointReq() ([]byte, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(delete_table)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, DELETETABLE_ENDPOINT)
}

func (delete *Delete) EndpointReq() ([]byte, int, error) {
	delete_table := DeleteTable(*delete)
	return delete_table.EndpointReq()
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	delete_table := DeleteTable(*req)
	return delete_table.EndpointReq()
}
