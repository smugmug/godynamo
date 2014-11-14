package create_table

import (
	"fmt"
	"testing"
	"encoding/json"
)

func TestRequestMarshal(t *testing.T) {
	s := []string{`{"AttributeDefinitions":[{"AttributeName":"ForumName","AttributeType":"S"},{"AttributeName":"Subject","AttributeType":"S"},{"AttributeName":"LastPostDateTime","AttributeType":"S"}],"TableName":"Thread","KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"Subject","KeyType":"RANGE"}],"LocalSecondaryIndexes":[{"IndexName":"LastPostIndex","KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType":"RANGE"}],"Projection":{"ProjectionType":"KEYS_ONLY"}}],"ProvisionedThroughput":{"ReadCapacityUnits":5,"WriteCapacityUnits":5}}`,`{"AttributeDefinitions":[{"AttributeName":"ForumName","AttributeType":"S"},{"AttributeName":"Subject","AttributeType":"S"},{"AttributeName":"LastPostDateTime","AttributeType":"S"}],"TableName":"Thread","KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"Subject","KeyType":"RANGE"}],"ProvisionedThroughput":{"ReadCapacityUnits":5,"WriteCapacityUnits":5}}`,
	}
	for _,v := range s {
		var c CreateTable
		um_err := json.Unmarshal([]byte(v),&c)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		json,jerr := json.Marshal(c)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n",v,string(json))
	}
}


func TestResponseMarshal(t *testing.T) {
	s := []string{`{"TableDescription":{"AttributeDefinitions":[{"AttributeName":"ForumName","AttributeType":"S"},{"AttributeName":"LastPostDateTime","AttributeType":"S"},{"AttributeName":"Subject","AttributeType":"S"}],"CreationDateTime":1.36372808007E9,"ItemCount":0,"KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"Subject","KeyType":"RANGE"}],"LocalSecondaryIndexes":[{"IndexName":"LastPostIndex","IndexSizeBytes":0,"ItemCount":0,"KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType":"RANGE"}],"Projection":{"ProjectionType":"KEYS_ONLY"}}],"ProvisionedThroughput":{"NumberOfDecreasesToday":0,"ReadCapacityUnits":5,"WriteCapacityUnits":5},"TableName":"Thread","TableSizeBytes":0,"TableStatus":"CREATING"}}`,
	}
	for _,v := range s {
		var c Response
		um_err := json.Unmarshal([]byte(v),&c)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		json,jerr := json.Marshal(c)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n",v,string(json))
	}
}
