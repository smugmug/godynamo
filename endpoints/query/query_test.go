package query

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{`{"TableName":"Thread","IndexName":"LastPostIndex","Select":"ALL_ATTRIBUTES","Limit":3,"ConsistentRead":true,"KeyConditions":{"LastPostDateTime":{"AttributeValueList":[{"S":"20130101"},{"S":"20130115"}],"ComparisonOperator":"BETWEEN"},"ForumName":{"AttributeValueList":[{"S":"AmazonDynamoDB"}],"ComparisonOperator":"EQ"}}}`, `{"TableName":"Thread","Select":"COUNT","ConsistentRead":true,"KeyConditions":{"ForumName":{"AttributeValueList":[{"S":"AmazonDynamoDB"}],"ComparisonOperator":"EQ"}}}`}
	for _, v := range s {
		var q Query
		um_err := json.Unmarshal([]byte(v), &q)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal Query: %v", um_err)
			t.Errorf(e)
		}
		json, jerr := json.Marshal(q)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n", v, string(json))
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{`{"Count":3,"Items":[{"LastPostedBy":{"S":"fred@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"20130102054211"},"Tags":{"SS":["Problem","Question"]}},{"LastPostedBy":{"S":"alice@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"20130105111307"},"Tags":{"SS":["Idea"]}},{"LastPostedBy":{"S":"bob@example.com"},"ForumName":{"S":"AmazonDynamoDB"},"LastPostDateTime":{"S":"20130108094417"},"Tags":{"SS":["AppDesign","HelpMe"]}}]}`, `{"Count":17}`}
	for _, v := range s {
		var q Response
		um_err := json.Unmarshal([]byte(v), &q)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal Response: %v", um_err)
			t.Errorf(e)
		}
		json, jerr := json.Marshal(q)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n", v, string(json))
	}
}
