package scan

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{`{"TableName":"Reply","ReturnConsumedCapacity":"TOTAL"}`,`{"TableName":"Reply","ScanFilter":{"PostedBy":{"AttributeValueList":[{"S":"joe@example.com"}],"ComparisonOperator":"EQ"}},"ReturnConsumedCapacity":"TOTAL"}`,
	}
	for _,v := range s {
		var q Scan
		um_err := json.Unmarshal([]byte(v),&q)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal Query: %v",um_err)
			t.Errorf(e)
		}
		json,jerr := json.Marshal(q)
		if jerr != nil {
			t.Errorf("cannot marshal %v\n",jerr)
		}
		fmt.Printf("IN:%v, OUT:%v\n",v,string(json))
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{`{"ConsumedCapacity":{"CapacityUnits":0.5,"TableName":"Reply"},"Count":4,"Items":[{"PostedBy":{"S":"joe@example.com"},"ReplyDateTime":{"S":"20130320115336"},"Id":{"S":"AmazonDynamoDB#HowdoIupdatemultipleitems?"},"Message":{"S":"HaveyoulookedattheBatchWriteItemAPI?"}},{"PostedBy":{"S":"fred@example.com"},"ReplyDateTime":{"S":"20130320115342"},"Id":{"S":"AmazonDynamoDB#HowdoIupdatemultipleitems?"},"Message":{"S":"No,Ididn'tknowaboutthat.WherecanIfindmoreinformation?"}},{"PostedBy":{"S":"joe@example.com"},"ReplyDateTime":{"S":"20130320115347"},"Id":{"S":"AmazonDynamoDB#HowdoIupdatemultipleitems?"},"Message":{"S":"BatchWriteItemisdocumentedintheAmazonDynamoDBAPIReference."}},{"PostedBy":{"S":"fred@example.com"},"ReplyDateTime":{"S":"20130320115352"},"Id":{"S":"AmazonDynamoDB#HowdoIupdatemultipleitems?"},"Message":{"S":"OK,I'lltakealookatthat.Thanks!"}}],"ScannedCount":4}`,`{"ConsumedCapacity":{"CapacityUnits":0.5,"TableName":"Reply"},"Count":2,"Items":[{"PostedBy":{"S":"joe@example.com"},"ReplyDateTime":{"S":"20130320115336"},"Id":{"S":"AmazonDynamoDB#HowdoIupdatemultipleitems?"},"Message":{"S":"HaveyoulookedattheBatchWriteItemAPI?"}},{"PostedBy":{"S":"joe@example.com"},"ReplyDateTime":{"S":"20130320115347"},"Id":{"S":"AmazonDynamoDB#HowdoIupdatemultipleitems?"},"Message":{"S":"BatchWriteItemisdocumentedintheAmazonDynamoDBAPIReference."}}],"ScannedCount":4}`,
	}
	for _,v := range s {
		var q Response
		um_err := json.Unmarshal([]byte(v),&q)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal Response: %v",um_err)
			t.Errorf(e)
		}
		json,jerr := json.Marshal(q)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n",v,string(json))
	}
}
