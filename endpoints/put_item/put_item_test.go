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
package put_item

import (
	"testing"
	"encoding/json"
)


func TestRequestMarshal(t *testing.T) {
	s := []string{
		`{
    "TableName": "Thread",
    "Item": {
        "LastPostDateTime": {
            "S": "201303190422"
        },
        "Tags": {
            "SS": ["Update","Multiple Items","HelpMe"]
        },
        "ForumName": {
            "S": "Amazon DynamoDB"
        },
        "Message": {
            "S": "I want to update multiple items in a single API call. What's the best way to do that?"
        },
        "Subject": {
            "S": "How do I update multiple items?"
        },
        "LastPostedBy": {
            "S": "fred@example.com"
        }
    },
    "Expected": {
        "ForumName": {
            "Exists": false
        },
        "Subject": {
            "Exists": false
        }
    }
}`,
	}
	for _,v := range s {
		var p Put
		um_err := json.Unmarshal([]byte(v),&p)
		if um_err != nil {
			t.Errorf("cannot unmarshal RequestItems:\n" + v + "\n")
		}
		_,jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}

	}
}

func TestResponseMarshal(t *testing.T) {
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
		}`,
	}
	for _,v := range s {
		var p Response
		um_err := json.Unmarshal([]byte(v),&p)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n" + v + "\n")
		}
		_,jerr := json.Marshal(p)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}
