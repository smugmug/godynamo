// Support for the DynamoDB UpdateTable endpoint.
//
// example use:
//
// tests/update_table-livestest.go
//
package update_table

import (
	"encoding/json"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/types/provisionedthroughput"
	"github.com/smugmug/godynamo/types/globalsecondaryindex"
	create_table "github.com/smugmug/godynamo/endpoints/create_table"
)

const (
	ENDPOINT_NAME        = "UpdateTable"
	UPDATETABLE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
)

type UpdateTable struct {
	GlobalSecondaryIndexUpdates *globalsecondaryindex.GlobalSecondaryIndexUpdates `json:",omitempty"`
	TableName string
	ProvisionedThroughput *provisionedthroughput.ProvisionedThroughput `json:",omitempty"`
}

func NewUpdateTable() (*UpdateTable) {
	update_table := new(UpdateTable)
	update_table.GlobalSecondaryIndexUpdates =
		globalsecondaryindex.NewGlobalSecondaryIndexUpdates()
	update_table.ProvisionedThroughput =
		provisionedthroughput.NewProvisionedThroughput()
	return update_table
}

// Update is an alias for backwards compatibility
type Update UpdateTable

func NewUpdate() (*Update) {
	update_table := NewUpdateTable()
	update := Update(*update_table)
	return &update
}

type Request UpdateTable

type Response create_table.Response

func NewResponse() (*Response) {
	cr := create_table.NewResponse()
	r := Response(*cr)
	return &r
}

func (update_table *UpdateTable) EndpointReq() (string,int,error) {
	// returns resp_body,code,err
	reqJSON,json_err := json.Marshal(update_table);
	if json_err != nil {
		return "",0,json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON,UPDATETABLE_ENDPOINT)
}

func (update *Update) EndpointReq() (string,int,error) {
	update_table := UpdateTable(*update)
	return update_table.EndpointReq()
}

func (req *Request) EndpointReq() (string,int,error) {
	update_table := UpdateTable(*req)
	return update_table.EndpointReq()
}
