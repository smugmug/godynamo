package describe_table

import (
	"fmt"
	"testing"
	"encoding/json"
)

func TestRequestMarshal(t *testing.T) {
	s := []string{
		`{"TableName":"Thread"}`,
	}
	for _,v := range s {
		var d DescribeTable
		um_err := json.Unmarshal([]byte(v),&d)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		_,jerr := json.Marshal(d)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}

	}
}


func TestResponseMarshal(t *testing.T) {
	s := []string{`{"Table":{"AttributeDefinitions":[{"AttributeName":"ForumName","AttributeType":"S"},{"AttributeName":"LastPostDateTime","AttributeType":"S"},{"AttributeName":"Subject","AttributeType":"S"}],"CreationDateTime":1.36372808007E9,"ItemCount":0,"KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"Subject","KeyType":"RANGE"}],"LocalSecondaryIndexes":[{"IndexName":"LastPostIndex","IndexSizeBytes":0,"ItemCount":0,"KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType":"RANGE"}],"Projection":{"ProjectionType":"KEYS_ONLY"}}],"ProvisionedThroughput":{"NumberOfDecreasesToday":0,"ReadCapacityUnits":5,"WriteCapacityUnits":5},"TableName":"Thread","TableSizeBytes":0,"TableStatus":"CREATING"}}`,
	}
	for _,v := range s {
		var d Response
		um_err := json.Unmarshal([]byte(v),&d)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		json,jerr := json.Marshal(d)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n",v,string(json))
	}
}
