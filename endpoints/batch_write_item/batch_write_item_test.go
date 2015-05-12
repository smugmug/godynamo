package batch_write_item

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/item"
	"testing"
)

func TestNil(t *testing.T) {
	b := NewBatchWriteItem()
	_,_,err := b.DoBatchWriteWithConf(nil)
	if err == nil {
		t.Errorf("nil conf should result in error")
	}
}

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{"RequestItems":{"Forum":[{"PutRequest":{"Item":{"Name":{"S":"AmazonDynamoDB"},"Category":{"S":"AmazonWebServices"}}}},{"PutRequest":{"Item":{"Name":{"S":"AmazonRDS"},"Category":{"S":"AmazonWebServices"}}}},{"PutRequest":{"Item":{"Name":{"S":"AmazonRedshift"},"Category":{"S":"AmazonWebServices"}}}},{"PutRequest":{"Item":{"Name":{"S":"AmazonElastiCache"},"Category":{"S":"AmazonWebServices"}}}}]},"ReturnConsumedCapacity":"TOTAL"}`,
	}
	for _, v := range s {
		b := NewBatchWriteItem()
		um_err := json.Unmarshal([]byte(v), b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchWriteItem: %v", um_err)
			t.Errorf(e)
		}
		_, jerr := json.Marshal(*b)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{
		`{"UnprocessedItems":{"Forum":[{"PutRequest":{"Item":{"Name":{"S":"AmazonElastiCache"},"Category":{"S":"AmazonWebServices"}}}}]},"ConsumedCapacity":[{"TableName":"Forum","CapacityUnits":3}]}`,
	}
	for _, v := range s {
		var b Response
		um_err := json.Unmarshal([]byte(v), &b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchWriteItem: %v", um_err)
			t.Errorf(e)
		}
		_, jerr := json.Marshal(b)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestUnprocessed(t *testing.T) {
	resp := `{
     "UnprocessedItems": {
         "Forum": [
             {
                 "PutRequest": {
                     "Item": {
                         "Name": {
                             "S": "Amazon ElastiCache"
                         },
                         "Category": {
                             "S": "Amazon Web Services"
                         }
                     }
                 }
             }
         ]
     },
     "ConsumedCapacity": [
         {
             "TableName": "Forum",
             "CapacityUnits": 3
         }
     ]
 }`
	req := `{
     "RequestItems": {
         "Forum": [
             {
                 "PutRequest": {
                     "Item": {
                         "Name": {
                             "S": "Amazon DynamoDB"
                         },
                         "Category": {
                             "S": "Amazon Web Services"
                         }
                     }
                 }
             },
             {
                 "PutRequest": {
                     "Item": {
                         "Name": {
                             "S": "Amazon RDS"
                         },
                         "Category": {
                             "S": "Amazon Web Services"
                         }
                     }
                 }
             },
             {
                 "PutRequest": {
                     "Item": {
                         "Name": {
                             "S": "Amazon Redshift"
                         },
                         "Category": {
                             "S": "Amazon Web Services"
                         }
                     }
                 }
             },
             {
                 "PutRequest": {
                     "Item": {
                         "Name": {
                             "S": "Amazon ElastiCache"
                         },
                         "Category": {
                             "S": "Amazon Web Services"
                         }
                     }
                 }
             }
         ]
     },
     "ReturnConsumedCapacity": "TOTAL"
 }`
	var r_resp Response
	um_err := json.Unmarshal([]byte(resp), &r_resp)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal BatchWriteItem: %v", um_err)
		t.Errorf(e)
	}
	_, jerr := json.Marshal(r_resp)
	if jerr != nil {
		t.Errorf("cannot marshal\n")
	}
	var r_req BatchWriteItem
	um_err = json.Unmarshal([]byte(req), &r_req)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal BatchWriteItem: %v", um_err)
		t.Errorf(e)
	}
	_, jerr2 := json.Marshal(r_req)
	if jerr2 != nil {
		t.Errorf("cannot marshal\n")
	}
	n_req, n_req_err := unprocessedItems2BatchWriteItems(&r_req, &r_resp)
	if n_req_err != nil {
		e := fmt.Sprintf("cannot make new batchwrite:%v", n_req_err)
		t.Errorf(e)
	}
	n_json, n_jerr := json.Marshal(*n_req)
	if n_jerr != nil {
		t.Errorf("cannot marshal\n")
	}
	_ = fmt.Sprintf("NEW:%s\n", string(n_json))
}

func TestSplit1(t *testing.T) {
	b := NewBatchWriteItem()
	b.RequestItems["foo"] = make([]RequestInstance, 0)
	for i := 0; i < 100; i++ {
		var p PutRequest
		p.Item = make(item.Item)
		k := fmt.Sprintf("TheKey%d", i)
		v := fmt.Sprintf("TheVal%d", i)
		p.Item["Key"] = &attributevalue.AttributeValue{S: k}
		p.Item["Val"] = &attributevalue.AttributeValue{S: v}
		b.RequestItems["foo"] = append(b.RequestItems["foo"], RequestInstance{PutRequest: &p})
	}
	bs, _ := Split(b)
	if len(bs) != 4 {
		e := fmt.Sprintf("len should be 4, it is %d\n", len(bs))
		t.Errorf(e)
	}
	i := 0
	for _, bsi := range bs {
		for _, ris := range bsi.RequestItems {
			if len(ris) != 25 {
				t.Errorf("requests len should be 25, it is not")

			}
		}
		i++
	}
}

func TestSplit2(t *testing.T) {
	b := NewBatchWriteItem()
	b.RequestItems["foo"] = make([]RequestInstance, 0)
	for i := 0; i < 23; i++ {
		var p PutRequest
		p.Item = make(item.Item)
		k := fmt.Sprintf("TheKey%d", i)
		v := fmt.Sprintf("TheVal%d", i)
		p.Item["Key"] = &attributevalue.AttributeValue{S: k}
		p.Item["Val"] = &attributevalue.AttributeValue{S: v}
		b.RequestItems["foo"] = append(b.RequestItems["foo"], RequestInstance{PutRequest: &p})
	}
	bs, _ := Split(b)
	if len(bs) != 1 {
		t.Errorf("list should have been split in 4, it is not")
	}
	i := 0
	for _, bsi := range bs {
		for _, ris := range bsi.RequestItems {
			if len(ris) != 23 {
				t.Errorf("requests len should be 25, it is not")
			}
		}
		i++
	}
}
