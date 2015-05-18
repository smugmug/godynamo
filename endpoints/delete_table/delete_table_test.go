package delete_table

import (
	"encoding/json"
	"testing"
)

func TestNil(t *testing.T) {
	d := NewDeleteTable()
	_,_,err := d.EndpointReqWithConf(nil)
	if err == nil {
		t.Errorf("nil conf should result in error")
	}
}

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{"TableName": "Reply"}`,
	}
	for _, v := range s {
		var u DeleteTable
		um_err := json.Unmarshal([]byte(v), &u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		_, jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{
		`{"TableDescription":{"ItemCount":0,"ProvisionedThroughput":{"NumberOfDecreasesToday":0,"ReadCapacityUnits":5,"WriteCapacityUnits":5},"TableName":"Reply","TableSizeBytes":0,"TableStatus":"DELETING"}}`,
	}
	for _, v := range s {
		var u Response
		um_err := json.Unmarshal([]byte(v), &u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		_, jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}
