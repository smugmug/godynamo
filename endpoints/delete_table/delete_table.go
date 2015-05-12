// Support for the DynamoDB DeleteTable endpoint.
//
// example use:
//
// tests/delete_table-livestest.go
//
package delete_table

import (
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
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

// These implementations of EndpointReq use a parameterized conf.

func (delete_table *DeleteTable) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
    	if delete_table == nil {
		return nil, 0, errors.New("delete_table.(DeleteTable)EndpointReqWithConf: receiver is nil")
	}
	if !conf.IsValid(c) {
		return nil, 0, errors.New("delete_table.EndpointReqWithConf: c is not valid")
	}
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(delete_table)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4WithConf(reqJSON, DELETETABLE_ENDPOINT, c)
}

func (delete *Delete) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if delete == nil {
		return nil, 0, errors.New("delete_table.(Delete)EndpointReqWithConf: receiver is nil")
	}
	delete_table := DeleteTable(*delete)
	return delete_table.EndpointReqWithConf(c)
}

func (req *Request) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("delete_table.(Request)EndpointReqWithConf: receiver is nil")
	}
	delete_table := DeleteTable(*req)
	return delete_table.EndpointReqWithConf(c)
}

// These implementations of EndpointReq use the global conf.

func (delete_table *DeleteTable) EndpointReq() ([]byte, int, error) {
	if delete_table == nil {
		return nil, 0, errors.New("delete_table.(DeleteTable)EndpointReq: receiver is nil")
	}
	return delete_table.EndpointReqWithConf(&conf.Vals)
}

func (delete *Delete) EndpointReq() ([]byte, int, error) {
	if delete == nil {
		return nil, 0, errors.New("delete_table.(Delete)EndpointReq: receiver is nil")
	}
	delete_table := DeleteTable(*delete)
	return delete_table.EndpointReqWithConf(&conf.Vals)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("delete_table.(Request)EndpointReq: receiver is nil")
	}
	delete_table := DeleteTable(*req)
	return delete_table.EndpointReqWithConf(&conf.Vals)
}
