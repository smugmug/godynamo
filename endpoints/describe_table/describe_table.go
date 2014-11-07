// Support for the DynamoDB DescribeTable endpoint.
//
// example use:
//
// tests/create_table-livestest.go, which contains a DescribeTable invocation
//
package describe_table

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smugmug/godynamo/authreq"
	"github.com/smugmug/godynamo/aws_const"
	ep "github.com/smugmug/godynamo/endpoint"
	"github.com/smugmug/godynamo/types/attributedefinition"
	"github.com/smugmug/godynamo/types/globalsecondaryindex"
	"github.com/smugmug/godynamo/types/keydefinition"
	"github.com/smugmug/godynamo/types/localsecondaryindex"
	"github.com/smugmug/godynamo/types/provisionedthroughput"
	"net/http"
	"time"
)

const (
	ENDPOINT_NAME      = "DescribeTable"
	DESCTABLE_ENDPOINT = aws_const.ENDPOINT_PREFIX + ENDPOINT_NAME
	ACTIVE             = "ACTIVE"
)

type DescribeTable struct {
	TableName string
}

// Describe is an alias for backwards compatibility
type Describe DescribeTable

type Request DescribeTable

func NewDescribeTable() *DescribeTable {
	d := new(DescribeTable)
	return d
}

type Response struct {
	Table struct {
		AttributeDefinitions   attributedefinition.AttributeDefinitions
		CreationDateTime       float64
		GlobalSecondaryIndexes []globalsecondaryindex.GlobalSecondaryIndexDesc
		ItemCount              uint64
		KeySchema              keydefinition.KeySchema
		LocalSecondaryIndexes  []localsecondaryindex.LocalSecondaryIndexDesc
		ProvisionedThroughput  provisionedthroughput.ProvisionedThroughputDesc
		TableName              string
		TableSizeBytes         uint64
		TableStatus            string
	}
}

func NewResponse() *Response {
	r := new(Response)
	r.Table.AttributeDefinitions = make(attributedefinition.AttributeDefinitions, 0)
	r.Table.GlobalSecondaryIndexes = make([]globalsecondaryindex.GlobalSecondaryIndexDesc, 0)
	r.Table.KeySchema = make(keydefinition.KeySchema, 0)
	r.Table.LocalSecondaryIndexes = make([]localsecondaryindex.LocalSecondaryIndexDesc, 0)
	return r
}

type StatusResult struct {
	StatusResult bool
}

func (describe_table *DescribeTable) EndpointReq() (string, int, error) {
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(describe_table)
	if json_err != nil {
		return "", 0, json_err
	}
	return authreq.RetryReqJSON_V4(reqJSON, DESCTABLE_ENDPOINT)
}

func (describe *Describe) EndpointReq() (string, int, error) {
	describe_table := DescribeTable(*describe)
	return describe_table.EndpointReq()
}

func (req *Request) EndpointReq() (string, int, error) {
	describe_table := DescribeTable(*req)
	return describe_table.EndpointReq()
}

// PollTableStatus allows the caller to poll a table for a specific status.
func PollTableStatus(tablename string, status string, tries int) (bool, error) {
	// aws docs informs us to poll the describe endpoint until the table
	// "status" is status for this tablename
	wait := time.Duration(2 * time.Second)

	for i := 0; i < tries; i++ {
		active, err := IsTableStatus(tablename, status)
		if err != nil {
			e := fmt.Sprintf("describe_table.PollStatus:%s",
				err.Error())
			return false, errors.New(e)
		}
		if active {
			return active, nil
		}
		time.Sleep(wait) // wait for table to become ACTIVE
	}
	return false, nil
}

// IsTableStatus will test the equality status of a table.
func IsTableStatus(tablename string, status string) (bool, error) {
	d := ep.Endpoint(&DescribeTable{TableName: tablename})
	s_resp, s_code, s_err := authreq.RetryReq_V4(d, DESCTABLE_ENDPOINT)
	if s_err != nil {
		e := fmt.Sprintf("describe_table.IsTableStatus: "+
			"check on %s err %s",
			tablename, s_err.Error())
		// if not a 500 problem, don't retry
		if !ep.ServerErr(s_code) {
			return false, errors.New(e)
		}
	}
	if s_resp != "" && s_code == http.StatusOK {
		var resp_json Response
		um_err := json.Unmarshal([]byte(s_resp), &resp_json)
		if um_err != nil {
			um_msg := fmt.Sprintf("describe_table.IsTableStatus:"+
				"cannot unmarshal %s, err: %s\ncheck "+
				"table creation of %s manually",
				s_resp, um_err.Error(), tablename)
			return false, errors.New(um_msg)
		}
		return (resp_json.Table.TableStatus == status), nil
	}
	e := fmt.Sprintf("describe_table.IsTableStatus:does %s exist?", tablename)
	return false, errors.New(e)
}

// TableExists test for table exists: exploit the fact that aws reports 4xx for tables that don't exist.
func (desc DescribeTable) TableExists() (bool, error) {
	_, code, err := desc.EndpointReq()
	if err != nil {
		e := fmt.Sprintf("describe_table.TableExists "+
			"%s", err.Error())
		return false, errors.New(e)
	}
	return (code == http.StatusOK), nil
}
