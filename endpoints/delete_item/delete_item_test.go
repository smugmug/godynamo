package delete_item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNil(t *testing.T) {
	d := NewDeleteItem()
	_,_,err := d.EndpointReqWithConf(nil)
	if err == nil {
		t.Errorf("nil conf should result in error")
	}
}

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"Expected":{"Replies":{"AttributeValueList":null,"Exists":false}},"ReturnValues":"ALL_OLD"}`, `{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}}}`,
	}
	for _, v := range s {
		d := NewDeleteItem()
		um_err := json.Unmarshal([]byte(v), d)
		if um_err != nil {
			_ = fmt.Sprintf("%v\n", um_err)
			t.Errorf("cannot unmarshal to delete:\n" + v + "\n")
		}
		_, jerr := json.Marshal(d)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{`{"Attributes":{"LastPostedBy":{"S":"fred@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"201303201023"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"Subject":{"S":"HowdoIupdatemultipleitems?"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"}}}`}
	for i, v := range s {
		var d Response
		um_err := json.Unmarshal([]byte(v), &d)
		if um_err != nil {
			t.Errorf("cannot unmarshal to delete:\n" + v + "\n")
		}
		json, jerr := json.Marshal(d)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		if len(s[i]) != len(string(json)) {
			e := fmt.Sprintf("\n%s\n%s\nshould be same",s[i],string(json))
			t.Errorf(e)
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n", v, string(json))
	}
}
