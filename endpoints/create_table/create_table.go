// Support for the DynamoDB CreateTable endpoint.
//
// example use:
//
// see tests/create_table-livestest.go
//
package create_table

import (
	"encoding/json"
	"errors"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/types/attributedefinition"
	"github.com/smugmug/godynamo/types/globalsecondaryindex"
	"github.com/smugmug/godynamo/types/keydefinition"
	"github.com/smugmug/godynamo/types/localsecondaryindex"
	"github.com/smugmug/godynamo/types/provisionedthroughput"
)

const (
	ENDPOINT_NAME        = "CreateTable"
	CREATETABLE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
)

type CreateTable struct {
	AttributeDefinitions   attributedefinition.AttributeDefinitions
	GlobalSecondaryIndexes []globalsecondaryindex.GlobalSecondaryIndex `json:",omitempty"`
	KeySchema              keydefinition.KeySchema
	LocalSecondaryIndexes  []localsecondaryindex.LocalSecondaryIndex `json:",omitempty"`
	ProvisionedThroughput  provisionedthroughput.ProvisionedThroughput
	TableName              string
}

func NewCreateTable() *CreateTable {
	c := new(CreateTable)
	c.AttributeDefinitions = make(attributedefinition.AttributeDefinitions, 0)
	c.GlobalSecondaryIndexes = make([]globalsecondaryindex.GlobalSecondaryIndex, 0)
	c.KeySchema = make(keydefinition.KeySchema, 0)
	c.LocalSecondaryIndexes = make([]localsecondaryindex.LocalSecondaryIndex, 0)
	return c
}

// Create is an alias for backwards compatibility
type Create CreateTable

func NewCreate() *Create {
	create_table := NewCreateTable()
	create := Create(*create_table)
	return &create
}

type Request CreateTable

type Response struct {
	TableDescription struct {
		AttributeDefinitions   attributedefinition.AttributeDefinitions        `json:",omitempty"`
		CreationDateTime       float64                                         `json:",omitempty"`
		GlobalSecondaryIndexes []globalsecondaryindex.GlobalSecondaryIndexDesc `json:",omitempty"`
		ItemCount              uint64                                          `json:",omitempty"`
		KeySchema              keydefinition.KeySchema                         `json:",omitempty"`
		LocalSecondaryIndexes  []localsecondaryindex.LocalSecondaryIndexDesc   `json:",omitempty"`
		ProvisionedThroughput  provisionedthroughput.ProvisionedThroughputDesc `json:",omitempty"`
		TableName              string
		TableSizeBytes         uint64 `json:",omitempty"`
		TableStatus            string
	}
}

func NewResponse() *Response {
	r := new(Response)
	r.TableDescription.AttributeDefinitions = make(attributedefinition.AttributeDefinitions, 0)
	r.TableDescription.GlobalSecondaryIndexes = make([]globalsecondaryindex.GlobalSecondaryIndexDesc, 0)
	r.TableDescription.KeySchema = make(keydefinition.KeySchema, 0)
	r.TableDescription.LocalSecondaryIndexes = make([]localsecondaryindex.LocalSecondaryIndexDesc, 0)

	return r
}

// These implementations of EndpointReq use a parameterized conf.

func (create_table *CreateTable) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if create_table == nil {
		return nil, 0, errors.New("create_table.(CreateTable)EndpointReqWithConf: receiver is nil")
	}
	if !conf.IsValid(c) {
		return nil, 0, errors.New("create_table.EndpointReqWithConf: c is not valid")
	}
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(create_table)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4WithConf(reqJSON, CREATETABLE_ENDPOINT, c)
}

func (create *Create) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if create == nil {
		return nil, 0, errors.New("create_table.(Create)EndpointReqWithConf: receiver is nil")
	}
	create_table := CreateTable(*create)
	return create_table.EndpointReqWithConf(c)
}

func (req *Request) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("create_table.(Request)EndpointReqWithConf: receiver is nil")
	}
	create_table := CreateTable(*req)
	return create_table.EndpointReqWithConf(c)
}

// These implementations of EndpointReq use the global conf.

func (create_table *CreateTable) EndpointReq() ([]byte, int, error) {
	if create_table == nil {
		return nil, 0, errors.New("create_table.(CreateTable)EndpointReq: receiver is nil")
	}
	return create_table.EndpointReqWithConf(&conf.Vals)
}

func (create *Create) EndpointReq() ([]byte, int, error) {
	if create == nil {
		return nil, 0, errors.New("create_table.(Create)EndpointReq: receiver is nil")
	}
	create_table := CreateTable(*create)
	return create_table.EndpointReqWithConf(&conf.Vals)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("create_table.(Request)EndpointReq: receiver is nil")
	}
	create_table := CreateTable(*req)
	return create_table.EndpointReqWithConf(&conf.Vals)
}

// ValidTable is a local validator that helps callers determine if a table name is too long.
func ValidTableName(t string) bool {
	l := len([]byte(t))
	return (l > 3) && (l < 256)
}
