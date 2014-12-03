// Support for the DynamoDB ListTables endpoint.
//
// example use:
//
// tests/create_table-livestest.go, which contains a ListTables invocation
//
package list_tables

import (
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
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

func (list_tables *ListTables) EndpointReq() ([]byte, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(list_tables)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, LISTTABLE_ENDPOINT)
}

func (list *List) EndpointReq() ([]byte, int, error) {
	list_tables := ListTables(*list)
	return list_tables.EndpointReq()
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	list_tables := ListTables(*req)
	return list_tables.EndpointReq()
}
