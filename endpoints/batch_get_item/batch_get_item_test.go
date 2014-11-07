// Tests JSON formats as described on the AWS docs site. For live tests, see ../../tests
package batch_get_item

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/item"
	"testing"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{`{"RequestItems":{"Forum":{"Keys":[{"Name":{"S":"AmazonDynamoDB"}},{"Name":{"S":"AmazonRDS"}},{"Name":{"S":"AmazonRedshift"}}],"AttributesToGet":["Name","Threads","Messages","Views"]},"Thread":{"Keys":[{"ForumName":{"S":"AmazonDynamoDB"}},{"Subject":{"S":"Concurrentreads"}}],"AttributesToGet":["Tags","Message"]}},"ReturnConsumedCapacity":"TOTAL"}`}
	for _, v := range s {
		b := NewBatchGetItem()
		um_err := json.Unmarshal([]byte(v), b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchGetItem: %v", um_err)
			t.Errorf(e)
		}
		_, jerr := json.Marshal(*b)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{`{"Responses":{"Forum":[{"Name":{"S":"AmazonDynamoDB"},"Threads":{"N":"5"},"Messages":{"N":"19"},"Views":{"N":"35"}},{"Name":{"S":"AmazonRDS"},"Threads":{"N":"8"},"Messages":{"N":"32"},"Views":{"N":"38"}},{"Name":{"S":"AmazonRedshift"},"Threads":{"N":"12"},"Messages":{"N":"55"},"Views":{"N":"47"}}],"Thread":[{"Tags":{"SS":["Reads","MultipleUsers"]},"Message":{"S":"Howmanyuserscanreadasingledataitematatime?Arethereanylimits?"}}]},"UnprocessedKeys":{"Forum":{"Keys":[{"Name":{"S":"AmazonDynamoDB"}},{"Name":{"S":"AmazonRDS"}},{"Name":{"S":"AmazonRedshift"}}],"AttributesToGet":["Name","Threads","Messages","Views"]}},"ConsumedCapacity":[{"TableName":"Forum","CapacityUnits":3},{"TableName":"Thread","CapacityUnits":1}]}`}
	for _, v := range s {
		b := NewResponse()
		um_err := json.Unmarshal([]byte(v), b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchGetItem: %v", um_err)
			t.Errorf(e)
		}
		_, jerr := json.Marshal(*b)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestUnprocessedToBatchGet(t *testing.T) {
	s := []string{`{"Responses":{"Forum":[{"Name":{"S":"AmazonDynamoDB"},"Threads":{"N":"5"},"Messages":{"N":"19"},"Views":{"N":"35"}},{"Name":{"S":"AmazonRDS"},"Threads":{"N":"8"},"Messages":{"N":"32"},"Views":{"N":"38"}},{"Name":{"S":"AmazonRedshift"},"Threads":{"N":"12"},"Messages":{"N":"55"},"Views":{"N":"47"}}],"Thread":[{"Tags":{"SS":["Reads","MultipleUsers"]},"Message":{"S":"Howmanyuserscanreadasingledataitematatime?Arethereanylimits?"}}]},"UnprocessedKeys":{"Forum":{"Keys":[{"Name":{"S":"AmazonDynamoDB"}},{"Name":{"S":"AmazonRDS"}},{"Name":{"S":"AmazonRedshift"}}],"AttributesToGet":["Name","Threads","Messages","Views"]}},"ConsumedCapacity":[{"TableName":"Forum","CapacityUnits":3},{"TableName":"Thread","CapacityUnits":1}]}`}
	for _, v := range s {
		bg := NewBatchGetItem()
		r := NewResponse()
		um_err := json.Unmarshal([]byte(v), r)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchGetItem: %v", um_err)
			t.Errorf(e)
		}
		nbg, _ := unprocessedKeys2BatchGetItems(bg, r)
		json, jerr := json.Marshal(*nbg)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("%s\n", string(json))
	}
}

func TestSplit1(t *testing.T) {
	b := NewBatchGetItem()
	b.RequestItems["foo"] = NewRequestInstance()
	for i := 0; i < 400; i++ {
		key := make(item.Item)
		k := fmt.Sprintf("TheKey%d", i)
		key["KeyName"] = &attributevalue.AttributeValue{S: k}
		b.RequestItems["foo"].Keys = append(b.RequestItems["foo"].Keys, key)
	}
	bs, _ := Split(b)
	if len(bs) != 4 {
		e := fmt.Sprintf("len should be 4, it is %d\n", len(bs))
		t.Errorf(e)
	}
	i := 0
	for _, bsi := range bs {
		json, _ := json.Marshal(bsi)
		fmt.Printf("\n\n%s\n\n", string(json))
		i++
	}
}
