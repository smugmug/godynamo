// Copyright (c) 2013,2014 SmugMug, Inc. All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
//       copyright notice, this list of conditions and the following
//       disclaimer in the documentation and/or other materials provided
//       with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY SMUGMUG, INC. ``AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
// PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL SMUGMUG, INC. BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE
// GOODS OR SERVICES;LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
// IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
// OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
// ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package batch_write_item

import (
	"testing"
	"fmt"
	"encoding/json"
	ep "github.com/smugmug/godynamo/endpoint"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		` {
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
 }`,
	}
	for _,v := range s {
		b := NewBatchWriteItem()
		um_err := json.Unmarshal([]byte(v),b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchWriteItem: %v",um_err)
			t.Errorf(e)
		}
		_,jerr := json.Marshal(*b)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{
		`{
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
 }`,
	}
	for _,v := range s {
		var b Response
		um_err := json.Unmarshal([]byte(v),&b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchWriteItem: %v",um_err)
			t.Errorf(e)
		}
		_,jerr := json.Marshal(b)
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
	um_err := json.Unmarshal([]byte(resp),&r_resp)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal BatchWriteItem: %v",um_err)
		t.Errorf(e)
	}
	_,jerr := json.Marshal(r_resp)
	if jerr != nil {
		t.Errorf("cannot marshal\n")
	}
	var r_req BatchWriteItem
	um_err = json.Unmarshal([]byte(req),&r_req)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal BatchWriteItem: %v",um_err)
		t.Errorf(e)
	}
	_,jerr2 := json.Marshal(r_req)
	if jerr2 != nil {
		t.Errorf("cannot marshal\n")
	}
	n_req,n_req_err := unprocessedItems2BatchWriteItems(r_req,&r_resp)
	if n_req_err != nil {
		e := fmt.Sprintf("cannot make new batchwrite:%v",n_req_err)
		t.Errorf(e)
	}
	n_json,n_jerr := json.Marshal(*n_req)
	if n_jerr != nil {
		t.Errorf("cannot marshal\n")
	}
	fmt.Printf("NEW:%s\n",string(n_json))
}

func TestSplit1(t *testing.T) {
	b := NewBatchWriteItem()
	b.RequestItems["foo"] = make([]RequestInstance,0)
	for i := 0; i< 100; i++ {
		var p PutRequest
		p.Item = make(ep.Item)
		k := fmt.Sprintf("TheKey%d",i)
		v := fmt.Sprintf("TheVal%d",i)
		p.Item["Key"] = ep.AttributeValue{S:k}
		p.Item["Val"] = ep.AttributeValue{S:v}
		b.RequestItems["foo"] = append(b.RequestItems["foo"],RequestInstance{PutRequest:&p})
	}
	bs,_ := Split(*b)
	if len(bs) != 4 {
		e := fmt.Sprintf("len should be 4, it is %d\n",len(bs))
		t.Errorf(e)
	}
	i := 0
	for _,bsi := range bs {
		for _,ris := range bsi.RequestItems {
			if len(ris) != 25 {
				t.Errorf("requests len should be 25, it is not")

			}
		}
		i++
	}
}

func TestSplit2(t *testing.T) {
	b := NewBatchWriteItem()
	b.RequestItems["foo"] = make([]RequestInstance,0)
	for i := 0; i< 23; i++ {
		var p PutRequest
		p.Item = make(ep.Item)
		k := fmt.Sprintf("TheKey%d",i)
		v := fmt.Sprintf("TheVal%d",i)
		p.Item["Key"] = ep.AttributeValue{S:k}
		p.Item["Val"] = ep.AttributeValue{S:v}
		b.RequestItems["foo"] = append(b.RequestItems["foo"],RequestInstance{PutRequest:&p})
	}
	bs,_ := Split(*b)
	if len(bs) != 1 {
		t.Errorf("list should have been split in 4, it is not")
	}
	i := 0
	for _,bsi := range bs {
		for _,ris := range bsi.RequestItems {
			if len(ris) != 23 {
				t.Errorf("requests len should be 25, it is not")
			}
		}
		i++
	}
}
