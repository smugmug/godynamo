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

package update_item

import (
	"testing"
	"encoding/json"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{
     "TableName": "Thread",
     "Key": {
         "ForumName": {
             "S": "Amazon DynamoDB"
         },
         "Subject": {
             "S": "How do I update multiple items?"
         }
     },
     "AttributeUpdates": {
         "LastPostedBy": {
             "Value": {
                 "S": "alice@example.com"
             },
             "Action": "PUT"
         }
     },
     "Expected": {
         "LastPostedBy": {
             "Value": {
                 "S": "fred@example.com"
             },
             "Exists": true
         }
     },
     "ReturnValues": "ALL_NEW"
 }`,
	}
	for _,v := range s {
		var u Update
		um_err := json.Unmarshal([]byte(v),&u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		_,jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{
		`{
     "Attributes": {
         "LastPostedBy": {
             "S": "alice@example.com"
         },
         "ForumName": {
             "S": "Amazon DynamoDB"
         },
         "LastPostDateTime": {
             "S": "20130320010350"
         },
         "Tags": {
             "SS": ["Update","Multiple Items","HelpMe"]
         },
         "Subject": {
             "S": "Maximum number of items?"
         },
         "Views": {
             "N": "5"
         },
         "Message": {
             "S": "I want to put 10 million data items to an Amazon DynamoDB table.  Is there an upper limit?"
         }
     }
 }
`,
	}
	for _,v := range s {
		var u Response
		um_err := json.Unmarshal([]byte(v),&u)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		_,jerr := json.Marshal(u)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}
