package delete_item

import (
	"fmt"
	"testing"
	"encoding/json"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"Expected":{"Replies":{"Exists":false}},"ReturnValues":"ALL_OLD"}`,`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}}}`,
	}
	for _,v := range s {
		d := NewDeleteItem()
		um_err := json.Unmarshal([]byte(v),d)
		if um_err != nil {
			_ = fmt.Sprintf("%v\n",um_err)
			t.Errorf("cannot unmarshal to delete:\n" + v + "\n")
		}
		json,jerr := json.Marshal(d)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n",v,string(json))
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{`{"Attributes":{"LastPostedBy":{"S":"fred@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"201303201023"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"Subject":{"S":"HowdoIupdatemultipleitems?"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"}}}`,
	}
	for _,v := range s {
		var d Response
		um_err := json.Unmarshal([]byte(v),&d)
		if um_err != nil {
			t.Errorf("cannot unmarshal to delete:\n" + v + "\n")
		}
		json,jerr := json.Marshal(d)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		_ = fmt.Sprintf("IN:%v, OUT:%v\n",v,string(json))
	}
}
