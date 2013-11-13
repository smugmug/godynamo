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

package scan

import (
	"testing"
	"encoding/json"
	"fmt"
 	ep "github.com/smugmug/godynamo/endpoint"
)

func TestSCanFilterOut(t *testing.T) {
	s := NewScan()
	s.TableName = "test-table"
	s.Limit = 30000
	var sf ScanFilter
	var av ep.AttributeValue
	av.N = "123456"
	sf.AttributeValueList = []ep.AttributeValue{av}
	sf.ComparisonOperator = OP_LT
	s.ScanFilter["rndtime"] = sf
	json,json_err := json.Marshal(s)
	if json_err != nil {
		e := fmt.Sprintf("cannot marshal %s",json_err.Error())
		t.Errorf(e)
	}
	fmt.Printf("%s\n",string(json))
}

func TestSCanFilterOut2(t *testing.T) {
	s := NewScan()
	s.TableName = "test-table"
	s.Limit = 30000
	var sf ScanFilter
	var av ep.AttributeValue
	av.N = "123456"
	sf.AttributeValueList = []ep.AttributeValue{av}
	sf.ComparisonOperator = OP_LT
	json,json_err := json.Marshal(s)
	if json_err != nil {
		e := fmt.Sprintf("cannot marshal %s",json_err.Error())
		t.Errorf(e)
	}
	fmt.Printf("%s\n",string(json))
}


func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{
     "TableName": "Reply",
     "ReturnConsumedCapacity": "TOTAL"
 }`,
 `{
     "TableName": "Reply",
     "ScanFilter": {
         "PostedBy": {
             "AttributeValueList": [
                 {
                     "S": "joe@example.com"
                 }
             ],
             "ComparisonOperator": "EQ"
         }
     },
     "ReturnConsumedCapacity": "TOTAL"
 }`,
	}
	for _,v := range s {
		var q Scan
		um_err := json.Unmarshal([]byte(v),&q)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal Query: %v",um_err)
			t.Errorf(e)
		}
		_,jerr := json.Marshal(q)
		if jerr != nil {
			t.Errorf("cannot marshal %v\n",jerr)
		}

	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{
		`{
     "ConsumedCapacity": {
         "CapacityUnits": 0.5,
         "TableName": "Reply"
     },
     "Count": 4,
     "Items": [
         {
             "PostedBy": {
                 "S": "joe@example.com"
             },
             "ReplyDateTime": {
                 "S": "20130320115336"
             },
             "Id": {
                 "S": "Amazon DynamoDB#How do I update multiple items?"
             },
             "Message": {
                 "S": "Have you looked at the BatchWriteItem API?"
             }
         },
         {
             "PostedBy": {
                 "S": "fred@example.com"
             },
             "ReplyDateTime": {
                 "S": "20130320115342"
             },
             "Id": {
                 "S": "Amazon DynamoDB#How do I update multiple items?"},
             "Message": {
                 "S": "No, I didn't know about that.  Where can I find more information?"
             }
         },
         {
             "PostedBy": {
                 "S": "joe@example.com"
             },
             "ReplyDateTime": {
                 "S": "20130320115347"
             },
             "Id": {
                 "S": "Amazon DynamoDB#How do I update multiple items?"
             },
             "Message": {
                 "S": "BatchWriteItem is documented in the Amazon DynamoDB API Reference."
             }
         },
         {
             "PostedBy": {
                 "S": "fred@example.com"
             },
             "ReplyDateTime": {
                 "S": "20130320115352"
             },
             "Id": {
                 "S": "Amazon DynamoDB#How do I update multiple items?"
             },
             "Message": {                 "S": "OK, I'll take a look at that.  Thanks!"
             }
         }
     ],
     "ScannedCount": 4
 }`,
 `{
     "ConsumedCapacity": {
         "CapacityUnits": 0.5,
         "TableName": "Reply"
     },
     "Count": 2,
     "Items": [
         {
             "PostedBy": {
                 "S": "joe@example.com"
             },
             "ReplyDateTime": {
                 "S": "20130320115336"
             },
             "Id": {
                 "S": "Amazon DynamoDB#How do I update multiple items?"
             },
             "Message": {
                 "S": "Have you looked at the BatchWriteItem API?"
             }
         },
         {
             "PostedBy": {
                 "S": "joe@example.com"
             },
             "ReplyDateTime": {
                 "S": "20130320115347"
             },
             "Id": {
                 "S": "Amazon DynamoDB#How do I update multiple items?"},
             "Message": {
                 "S": "BatchWriteItem is documented in the Amazon DynamoDB API Reference."
             }
         }
     ],
     "ScannedCount": 4
 }`,
	}
	for _,v := range s {
		var q Response
		um_err := json.Unmarshal([]byte(v),&q)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal Response: %v",um_err)
			t.Errorf(e)
		}
		_,jerr := json.Marshal(q)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}

	}
}
