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
	"github.com/smugmug/godynamo/conf"
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

// These implementations of EndpointReq use a parameterized conf.

func (describe_table *DescribeTable) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if describe_table == nil {
		return nil, 0, errors.New("describe_table.(DescribeTable)EndpointReqWithConf: receiver is nil")
	}
	if !conf.IsValid(c) {
		return nil, 0, errors.New("describe_table.EndpointReqWithConf: c is not valid")
	}
	// returns resp_body,code,err
	reqJSON, json_err := json.Marshal(describe_table)
	if json_err != nil {
		return nil, 0, json_err
	}
	return authreq.RetryReqJSON_V4WithConf(reqJSON, DESCTABLE_ENDPOINT, c)
}

func (describe *Describe) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if describe == nil {
		return nil, 0, errors.New("describe_table.(Describe)EndpointReqWithConf: receiver is nil")
	}
	describe_table := DescribeTable(*describe)
	return describe_table.EndpointReqWithConf(c)
}

func (req *Request) EndpointReqWithConf(c *conf.AWS_Conf) ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("describe_table.(Request)EndpointReqWithConf: receiver is nil")
	}
	describe_table := DescribeTable(*req)
	return describe_table.EndpointReqWithConf(c)
}

// These implementations of EndpointReq use the global conf.

func (describe_table *DescribeTable) EndpointReq() ([]byte, int, error) {
	if describe_table == nil {
		return nil, 0, errors.New("describe_table.(DescribeTable)EndpointReq: receiver is nil")
	}
	return describe_table.EndpointReqWithConf(&conf.Vals)
}

func (describe *Describe) EndpointReq() ([]byte, int, error) {
	if describe == nil {
		return nil, 0, errors.New("describe_table.(Describe)EndpointReq: receiver is nil")
	}
	describe_table := DescribeTable(*describe)
	return describe_table.EndpointReqWithConf(&conf.Vals)
}

func (req *Request) EndpointReq() ([]byte, int, error) {
	if req == nil {
		return nil, 0, errors.New("describe_table.(Request)EndpointReq: receiver is nil")
	}
	describe_table := DescribeTable(*req)
	return describe_table.EndpointReqWithConf(&conf.Vals)
}

// PollTableStatusWithConf allows the caller to poll a table for a specific status.
func PollTableStatusWithConf(tablename string, status string, tries int, c *conf.AWS_Conf) (bool, error) {
	if !conf.IsValid(c) {
		return false, errors.New("describe_table.PollTableStatusWithConf: c is not valid")
	}
	// aws docs informs us to poll the describe endpoint until the table
	// "status" is status for this tablename
	wait := time.Duration(2 * time.Second)

	for i := 0; i < tries; i++ {
		active, err := IsTableStatusWithConf(tablename, status, c)
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

// PollTableStatus is the same as PollTableStatusWithConf but uses the global conf.Vals.
func PollTableStatus(tablename string, status string, tries int) (bool, error) {
	return PollTableStatusWithConf(tablename, status, tries, &conf.Vals)
}

// IsTableStatusWithConf will test the equality status of a table.
func IsTableStatusWithConf(tablename string, status string, c *conf.AWS_Conf) (bool, error) {
	if !conf.IsValid(c) {
		return false, errors.New("describe_table.IsTableStatusWithConf: c is not valid")
	}
	d := ep.Endpoint(&DescribeTable{TableName: tablename})
	s_resp, s_code, s_err := authreq.RetryReq_V4WithConf(d, DESCTABLE_ENDPOINT, c)
	if s_err != nil {
		e := fmt.Sprintf("describe_table.IsTableStatus: "+
			"check on %s err %s",
			tablename, s_err.Error())
		// if not a 500 problem, don't retry
		if !ep.ServerErr(s_code) {
			return false, errors.New(e)
		}
	}
	if s_resp != nil && s_code == http.StatusOK {
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

// IsTableStatus is the same as IsTableStatusWithConf but uses the global conf.Vals.
func IsTableStatus(tablename string, status string) (bool, error) {
	return IsTableStatusWithConf(tablename, status, &conf.Vals)
}

// TableExistsWithconf test for table exists: exploit the fact that aws reports 4xx for tables that don't exist.
func (desc DescribeTable) TableExistsWithConf(c *conf.AWS_Conf) (bool, error) {
	if !conf.IsValid(c) {
		return false, errors.New("describe_table.TableExistsWithConf: c is not valid")
	}
	_, code, err := desc.EndpointReqWithConf(c)
	if err != nil {
		e := fmt.Sprintf("describe_table.TableExistsWithConf "+
			"%s", err.Error())
		return false, errors.New(e)
	}
	return (code == http.StatusOK), nil
}

// TableExists is the same as TableExistsWithConf but uses the global conf.Vals.
func (desc DescribeTable) TableExists() (bool, error) {
	return desc.TableExistsWithConf(&conf.Vals)
}
