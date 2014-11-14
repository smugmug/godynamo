package delete_table

import (
	"fmt"
	"testing"
	"encoding/json"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{"TableName": "Reply"}`,
	}
	for _,v := range s {
		var u DeleteTable
		um_err := json.Unmarshal([]byte(v),&u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		_,jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{
		`{"TableDescription":{"ItemCount":0,"ProvisionedThroughput":{"NumberOfDecreasesToday":0,"ReadCapacityUnits":5,"WriteCapacityUnits":5},"TableName":"Reply","TableSizeBytes":0,"TableStatus":"DELETING"}}`,
	}
	for _,v := range s {
		var u Response
		um_err := json.Unmarshal([]byte(v),&u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		json,jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n",v,string(json))
	}
}
