// Support for the DynamoDB ListTables endpoint.
//
// example use:
//
// tests/create_table-livestest.go, which contains a ListTables invocation
//
package list_tables

import (
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
)

const (
	ENDPOINT_NAME              = "ListTables"
	LISTTABLE_ENDPOINT         = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	EXCLUSIVE_START_TABLE_NAME = "ExclusiveStartTableName"
	LIMIT                      = "Limit"
	LAST_EVALUATED_TABLE_NAME  = "LastEvaluatedTableName"
	AWS_LIMIT                  = 100
)

type ListTables struct {
	ExclusiveStartTableName string `json:",omitempty"`
	Limit                   uint64 `json:",omitempty"`
}

// List is an alias for backwards compatibility
type List ListTables

type Request ListTables

type Response struct {
	TableNames             []string
	LastEvaluatedTableName string `json:",omitempty"`
}

func NewResponse() *Response {
	r := new(Response)
	r.TableNames = make([]string, 0)
	return r
}

// These implementations of EndpointReq use a parameterized conf.

func (list_tables *ListTables) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if list_tables == nil {
		return nil, 0, errors.New("list_tables.(ListTables)EndpointReqWithConf: receiver is nil")
	}
	if !conf.IsValid(c) {
		return nil, 0, errors.New("list_tables.EndpointReqWithConf: c is not valid")
	}
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(list_tables)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4WithConf(reqJSON, LISTTABLE_ENDPOINT, c)
}

func (list *List) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if list == nil {
		return nil, 0, errors.New("list_tables.(List)EndpointReqWithConf: receiver is nil")
	}
	list_tables := ListTables(*list)
	return list_tables.EndpointReqWithConf(c)
}

func (req *Request) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("list_tables.(Request)EndpointReqWithConf: receiver is nil")
	}
	list_tables := ListTables(*req)
	return list_tables.EndpointReqWithConf(c)
}

// These implementations of EndpointReq use the global conf.

func (list_table *ListTables) EndpointReq() ([]byte, int, error) {
	if list_table == nil {
		return nil, 0, errors.New("list_tables.(ListTables)EndpointReq: receiver is nil")
	}
	return list_table.EndpointReqWithConf(&conf.Vals)
}

func (list *List) EndpointReq() ([]byte, int, error) {
	if list == nil {
		return nil, 0, errors.New("list_tables.(List)EndpointReq: receiver is nil")
	}
	list_table := ListTables(*list)
	return list_table.EndpointReqWithConf(&conf.Vals)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("list_tables.(Request)EndpointReq: receiver is nil")
	}
	list_table := ListTables(*req)
	return list_table.EndpointReqWithConf(&conf.Vals)
}
