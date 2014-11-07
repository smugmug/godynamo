package update_item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{`{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"HowdoIupdatemultipleitems?"}},"AttributeUpdates":{"LastPostedBy":{"Value":{"S":"alice@example.com"},"Action":"PUT"}},"Expected":{"LastPostedBy":{"Value":{"S":"fred@example.com"},"Exists":true}},"ReturnValues":"ALL_NEW"}`, `{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"Maximumnumberofitems?"}},"UpdateExpression":"setLastPostedBy=:val1","ConditionExpression":"LastPostedBy=:val2","ExpressionAttributeValues":{":val1":{"S":"alice@example.com"},":val2":{"S":"fred@example.com"}},"ReturnValues":"ALL_NEW"}`, `{"TableName":"Thread","Key":{"ForumName":{"S":"AmazonDynamoDB"},"Subject":{"S":"Aquestionaboutupdates"}},"UpdateExpression":"setReplies=Replies+:num","ExpressionAttributeValues":{":num":{"N":"1"}},"ReturnValues":"NONE"}`}
	for _, v := range s {
		var u UpdateItem
		um_err := json.Unmarshal([]byte(v), &u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		json, jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n", v, string(json))
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{`{"Attributes":{"LastPostedBy":{"S":"alice@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"20130320010350"},"Tags":{"SS":["Update","MultipleItems","HelpMe"]},"Subject":{"S":"Maximumnumberofitems?"},"Views":{"N":"5"},"Message":{"S":"Iwanttoput10milliondataitemstoanAmazonDynamoDBtable.Isthereanupperlimit?"}}}`}
	for _, v := range s {
		var u Response
		um_err := json.Unmarshal([]byte(v), &u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		json, jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n", v, string(json))
	}
}
