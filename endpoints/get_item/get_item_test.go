package get_item

import (
	"fmt"
	"testing"
	"encoding/json"
)

func TestRequestMarshal(t *testing.T) {
	s := []string{
		`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"AttributesToGet":["LastPostDateTime","Message","Tags"],"ConsistentRead":true,"ReturnConsumedCapacity":"TOTAL"}`,`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"ProjectionExpression":"LastPostDateTime,Message,Tags","ConsistentRead":true,"ReturnConsumedCapacity":"TOTAL"}`,
	}
	for _,v := range s {
		var g GetItem
		um_err := json.Unmarshal([]byte(v),&g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		json,jerr := json.Marshal(g)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n",v,string(json))
	}
}

func TestResponseMarshal(t *testing.T) {
	s := []string{`{"ConsumedCapacity":{"CapacityUnits":1,"TableName":"Thread"},"Item":{"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"LastPostDateTime":{"S":"201303190436"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"}}}`,
	}
	for _,v := range s {
		var g Response
		um_err := json.Unmarshal([]byte(v),&g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		json,jerr := json.Marshal(g)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n",v,string(json))
	}
}
