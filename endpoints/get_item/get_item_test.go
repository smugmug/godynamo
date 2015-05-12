package get_item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNil(t *testing.T) {
	g := NewGetItem()
	_,_,err := g.EndpointReqWithConf(nil)
	if err == nil {
		t.Errorf("nil conf should result in error")
	}
}

func TestRequestMarshal(t *testing.T) {
	s := []string{
		`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"AttributesToGet":["LastPostDateTime","Message","Tags"],"ConsistentRead":true,"ReturnConsumedCapacity":"TOTAL"}`, `{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"ProjectionExpression":"LastPostDateTime,Message,Tags","ConsistentRead":true,"ReturnConsumedCapacity":"TOTAL"}`,
	}
	for i, v := range s {
		var g GetItem
		um_err := json.Unmarshal([]byte(v), &g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		json, jerr := json.Marshal(g)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		if len(s[i]) != len(string(json)) {
			e := fmt.Sprintf("\n%s\n%s\nshould be same",s[i],string(json))
			t.Errorf(e)
		}
	}
}

func TestResponseMarshal(t *testing.T) {
	s := []string{`{"ConsumedCapacity":{"CapacityUnits":1,"TableName":"Thread"},"Item":{"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"LastPostDateTime":{"S":"201303190436"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"}}}`}
	j := []string{`{"Item":{"LastPostDateTime":"201303190436","Message":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?","Tags":["MultipleItems","HelpMe","Update"]},"ConsumedCapacity":{"CapacityU\nits":1,"TableName":"Thread"}`}
	for i, v := range s {
		var g Response
		um_err := json.Unmarshal([]byte(v), &g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		json1, jerr1 := json.Marshal(g)
		if jerr1 != nil {
			t.Errorf("cannot marshal\n")
		}
		if len(s[i]) != len(string(json1)) {
			e := fmt.Sprintf("\n%s\n%s\nshould be same",s[i],string(json1))
			t.Errorf(e)
		}
		c, cerr := g.ToResponseItemJSON()
		if cerr != nil {
			e := fmt.Sprintf("cannot convert %v\n", cerr)
			t.Errorf(e)
		}
		json2, jerr2 := json.Marshal(c)
		if jerr2 != nil {
			t.Errorf("cannot marshal\n")
		}
		if len(j[i]) != len(string(json2)) {
			e := fmt.Sprintf("\n%s\n%s\nshould be same",s[i],string(json2))
			t.Errorf(e)
		}
	}
}
