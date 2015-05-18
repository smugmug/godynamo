// Support for the DynamoDB UpdateTable endpoint.
//
// example use:
//
// tests/update_table-livestest.go
//
package update_table

import (
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
	create_table "github.com/smugmug/godynamo/endpoints/create_table"
	"github.com/smugmug/godynamo/types/globalsecondaryindex"
	"github.com/smugmug/godynamo/types/provisionedthroughput"
)

const (
	ENDPOINT_NAME        = "UpdateTable"
	UPDATETABLE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
)

type UpdateTable struct {
	GlobalSecondaryIndexUpdates *globalsecondaryindex.GlobalSecondaryIndexUpdates `json:",omitempty"`
	TableName                   string
	ProvisionedThroughput       *provisionedthroughput.ProvisionedThroughput `json:",omitempty"`
}

func NewUpdateTable() *UpdateTable {
	update_table := new(UpdateTable)
	update_table.GlobalSecondaryIndexUpdates =
		globalsecondaryindex.NewGlobalSecondaryIndexUpdates()
	update_table.ProvisionedThroughput =
		provisionedthroughput.NewProvisionedThroughput()
	return update_table
}

// Update is an alias for backwards compatibility
type Update UpdateTable

func NewUpdate() *Update {
	update_table := NewUpdateTable()
	update := Update(*update_table)
	return &update
}

type Request UpdateTable

type Response create_table.Response

func NewResponse() *Response {
	cr := create_table.NewResponse()
	r := Response(*cr)
	return &r
}

// These implementations of EndpointReq use a parameterized conf.

func (update_table *UpdateTable) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if update_table == nil {
		return nil, 0, errors.New("update_table.(UpdateTable)EndpointReqWithConf: receiver is nil")
	}
	if !conf.IsValid(c) {
		return nil, 0, errors.New("update_table.EndpointReqWithConf: c is not valid")
	}
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(update_table)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4WithConf(reqJSON, UPDATETABLE_ENDPOINT, c)
}

func (update *Update) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if update == nil {
		return nil, 0, errors.New("update_table.(Update)EndpointReqWithConf: receiver is nil")
	}
	update_table := UpdateTable(*update)
	return update_table.EndpointReqWithConf(c)
}

func (req *Request) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("update_table.(Request)EndpointReqWithConf: receiver is nil")
	}
	update_table := UpdateTable(*req)
	return update_table.EndpointReqWithConf(c)
}

// These implementations of EndpointReq use the global conf.

func (update_table *UpdateTable) EndpointReq() ([]byte, int, error) {
	if update_table == nil {
		return nil, 0, errors.New("update_table.(UpdateTable)EndpointReq: receiver is nil")
	}
	return update_table.EndpointReqWithConf(&conf.Vals)
}

func (update *Update) EndpointReq() ([]byte, int, error) {
	if update == nil {
		return nil, 0, errors.New("update_table.(Update)EndpointReq: receiver is nil")
	}
	update_table := UpdateTable(*update)
	return update_table.EndpointReqWithConf(&conf.Vals)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("update_table.(Request)EndpointReq: receiver is nil")
	}
	update_table := UpdateTable(*req)
	return update_table.EndpointReqWithConf(&conf.Vals)
}
