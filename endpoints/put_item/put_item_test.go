// Tests JSON formats as described on the AWS docs site. For live tests, see ../../tests
package put_item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNil(t *testing.T) {
	g := NewPutItem()
	_,_,err := g.EndpointReqWithConf(nil)
	if err == nil {
		t.Errorf("nil conf should result in error")
	}
}

func TestRequestMarshal(t *testing.T) {
	s := []string{
		`{"Expected":{"ForumName":{"AttributeValueList":null,"ComparisonOperator":"","Value":null,"Exists":false},"Subject":{"AttributeValueList":null,"ComparisonOperator":"","Value":null,"Exists":false},"Value":null},"Item":{"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"201303190422"},"LastPostedBy":{"S":"fred@example.com"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"},"Subject":{"S":"HowdoIupdatemultipleitems?"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]}},"TableName":"Thread"}`,
		`{"TableName":"Thread","Item":{"LastPostDateTime":{"S":"201303190422"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"ForumName":{"S":"AmazonDynamoDB"},"Message":{"S":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?"},"Subject":{"S":"HowdoIupdatemultipleitems?"},"LastPostedBy":{"S":"fred@example.com"}},"ConditionExpression":"ForumName<>:fandSubject<>:s","ExpressionAttributeValues":{":f":{"S":"AmazonDynamoDB"},":s":{"S":"HowdoIupdatemultipleitems?"}}}`}
	for i, v := range s {
		var p PutItem
		um_err := json.Unmarshal([]byte(v), &p)
		if um_err != nil {
			t.Errorf("cannot unmarshal RequestItems:\n" + v + "\n")
		}
		json, jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		if i == 0 {
			if len(s[i]) != len(string(json)) {
				e := fmt.Sprintf("\n%s\n%s\nshould be same",s[i],string(json))
				t.Errorf(e)
			}
		}
	}
}

func TestRequestJSONMarshal(t *testing.T) {
	s := []string{`{"TableName":"Thread","Item":{"LastPostDateTime":"201303190422","Tags":["Update","MultipleItems","HelpMe"],"ForumName":"AmazonDynamoDB","Message":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?","Subject":"HowdoIupdatemultipleitems?","LastPostedBy":"fred@example.com"},"Expected":{"ForumName":{"Exists":false},"Subject":{"Exists":false}}}`, `{"TableName":"Thread","Item":{"LastPostDateTime":"201303190422","Tags":["Update","MultipleItems","HelpMe"],"ForumName":"AmazonDynamoDB","Message":"IwanttoupdatemultipleitemsinasingleAPIcall.What'sthebestwaytodothat?","Subject":"HowdoIupdatemultipleitems?","LastPostedBy":"fred@example.com"},"ConditionExpression":"ForumName<>:fandSubject<>:s","ExpressionAttributeValues":{":f":{"S":"AmazonDynamoDB"},":s":{"S":"HowdoIupdatemultipleitems?"}}}`}
	for _, v := range s {
		var p PutItemJSON
		um_err := json.Unmarshal([]byte(v), &p)
		if um_err != nil {
			t.Errorf("cannot unmarshal RequestItems:\n" + v + "\n")
		}
		json, jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		// there is some unicode encoding here, so if you want, uncomment this line and
		// do visual qa
		_ = fmt.Sprintf("JSON IN:%v, OUT:%v\n", v, string(json))
	}
}

func TestResponseMarshal(t *testing.T) {
	s := []string{`{"Attributes":{"LastPostedBy":{"S":"alice@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"20130320010350"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"Subject":{"S":"Maximumnumberofitems?"},"Views":{"N":"5"},"Message":{"S":"Iwanttoput10milliondataitemstoanAmazonDynamoDBtable.Isthereanupperlimit?"}}}`}
	for i, v := range s {
		var p Response
		um_err := json.Unmarshal([]byte(v), &p)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n" + v + "\n")
		}
		json, jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		if len(s[i]) != len(string(json)) {
			e := fmt.Sprintf("\n%s\n%s\nshould be same",s[i],string(json))
			t.Errorf(e)
		}
	}
}
