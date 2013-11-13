// Copyright (c) 2013, SmugMug, Inc. All rights reserved.
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

// Tests JSON formats as described on the AWS docs site. For live tests, see ../../tests
package batch_get_item

import (
	"testing"
	"fmt"
	"encoding/json"
	ep "github.com/smugmug/godynamo/endpoint"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{
     "RequestItems": {
         "Forum": {
             "Keys": [
                 {
                     "Name":{"S":"Amazon DynamoDB"}
                 },
                 {
                     "Name":{"S":"Amazon RDS"}
                 },
                 {
                     "Name":{"S":"Amazon Redshift"}
                 }
             ],
             "AttributesToGet": [
                 "Name","Threads","Messages","Views"
             ]
         },
         "Thread": {
             "Keys": [
                 {
                     "ForumName":{"S":"Amazon DynamoDB"}
                 },
                 {
                     "Subject":{"S":"Concurrent reads"}
                 }
             ],
             "AttributesToGet": [
                 "Tags","Message"
             ]
         }
     },
     "ReturnConsumedCapacity": "TOTAL"
 }`,
	}
	for _,v := range s {
		b := NewBatchGetItem()
		um_err := json.Unmarshal([]byte(v),b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchGetItem: %v",um_err)
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
    "Responses": {
        "Forum": [
            {
                "Name":{
                    "S":"Amazon DynamoDB"
                },
                "Threads":{
                    "N":"5"
                },
                "Messages":{
                    "N":"19"
                },
                "Views":{
                    "N":"35"
                }
            },
            {
                "Name":{
                    "S":"Amazon RDS"
                },
                "Threads":{
                    "N":"8"
                },
                "Messages":{
                    "N":"32"
                },
                "Views":{
                    "N":"38"
                }
            },
            {
                "Name":{
                    "S":"Amazon Redshift"
                },
                "Threads":{
                    "N":"12"
                },
                "Messages":{
                    "N":"55"
                },
                "Views":{
                    "N":"47"
                }
            }
        ],
        "Thread": [
            {
                "Tags":{
                    "SS":["Reads","MultipleUsers"]
                },
                "Message":{
                    "S":"How many users can read a single data item at a time? Are there any limits?"
                }
            }
        ]
    },
    "UnprocessedKeys": {
        "Forum": {
             "Keys": [
                 {
                     "Name":{"S":"Amazon DynamoDB"}
                 },
                 {
                     "Name":{"S":"Amazon RDS"}
                 },
                 {
                     "Name":{"S":"Amazon Redshift"}
                 }
             ],
             "AttributesToGet": [
                 "Name","Threads","Messages","Views"
             ]
         }
    },
    "ConsumedCapacity": [
        {
            "TableName": "Forum",
            "CapacityUnits": 3
        },
        {
            "TableName": "Thread",
            "CapacityUnits": 1
        }
    ]
}`,

	}
	for _,v := range s {
		b := NewResponse()
		um_err := json.Unmarshal([]byte(v),b)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchGetItem: %v",um_err)
			t.Errorf(e)
		}
		_,jerr := json.Marshal(*b)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestUnprocessedToBatchGet(t *testing.T) {
	s := []string{
		`{
    "Responses": {
        "Forum": [
            {
                "Name":{
                    "S":"Amazon DynamoDB"
                },
                "Threads":{
                    "N":"5"
                },
                "Messages":{
                    "N":"19"
                },
                "Views":{
                    "N":"35"
                }
            },
            {
                "Name":{
                    "S":"Amazon RDS"
                },
                "Threads":{
                    "N":"8"
                },
                "Messages":{
                    "N":"32"
                },
                "Views":{
                    "N":"38"
                }
            },
            {
                "Name":{
                    "S":"Amazon Redshift"
                },
                "Threads":{
                    "N":"12"
                },
                "Messages":{
                    "N":"55"
                },
                "Views":{
                    "N":"47"
                }
            }
        ],
        "Thread": [
            {
                "Tags":{
                    "SS":["Reads","MultipleUsers"]
                },
                "Message":{
                    "S":"How many users can read a single data item at a time? Are there any limits?"
                }
            }
        ]
    },
    "UnprocessedKeys": {
        "Forum": {
             "Keys": [
                 {
                     "Name":{"S":"Amazon DynamoDB"}
                 },
                 {
                     "Name":{"S":"Amazon RDS"}
                 },
                 {
                     "Name":{"S":"Amazon Redshift"}
                 }
             ],
             "AttributesToGet": [
                 "Name","Threads","Messages","Views"
             ]
         }
    },
    "ConsumedCapacity": [
        {
            "TableName": "Forum",
            "CapacityUnits": 3
        },
        {
            "TableName": "Thread",
            "CapacityUnits": 1
        }
    ]
}`,

	}
	for _,v := range s {
		bg := NewBatchGetItem()
		r := NewResponse()
		um_err := json.Unmarshal([]byte(v),r)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal BatchGetItem: %v",um_err)
			t.Errorf(e)
		}
		nbg,_ := unprocessedKeys2BatchGetItems(*bg,r)
		r.UnprocessedKeys["Forum"].Keys[0]["Name"] = ep.AttributeValue{S:"foo"}
		json,jerr := json.Marshal(*nbg)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("%s\n",string(json))
	}
}

func TestSplit1(t *testing.T) {
	b := NewBatchGetItem()
	b.RequestItems["foo"] = NewRequestInstance()
	for i := 0; i< 400; i++ {
		key := make(ep.Item)
		k := fmt.Sprintf("TheKey%d",i)
		key["KeyName"] = ep.AttributeValue{S:k}
		b.RequestItems["foo"].Keys = append(b.RequestItems["foo"].Keys,key)
	}
	bs,_ := Split(*b)
	if len(bs) != 4 {
		e := fmt.Sprintf("len should be 4, it is %d\n",len(bs))
		t.Errorf(e)
	}
	i := 0
	for _,bsi := range bs {
		json,_ := json.Marshal(bsi)
		fmt.Printf("\n\n%s\n\n",string(json))
		i++
	}
}
