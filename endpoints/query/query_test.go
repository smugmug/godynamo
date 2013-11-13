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

package query

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestRequestUnmarshal(t *testing.T) {
	s := []string{
		`{
    "TableName": "Thread",
    "IndexName": "LastPostIndex",
    "Select": "ALL_ATTRIBUTES",
    "Limit":3,
    "ConsistentRead": true,
    "KeyConditions": {
        "LastPostDateTime": {
            "AttributeValueList": [
                {
                    "S": "20130101"
                },
                {
                    "S": "20130115"
                }
            ],
            "ComparisonOperator": "BETWEEN"
        },
        "ForumName": {
            "AttributeValueList": [
                {
                    "S": "Amazon DynamoDB"
                }
            ],
            "ComparisonOperator": "EQ"
        }
    }
}`,
`
{
    "TableName": "Thread",
    "Select": "COUNT",
    "ConsistentRead": true,
    "KeyConditions": {
        "ForumName": {
            "AttributeValueList": [
                {
                    "S": "Amazon DynamoDB"
                }
            ],
            "ComparisonOperator": "EQ"
        }
    }
}
  `,
	}
	for _,v := range s {
		var q Query
		um_err := json.Unmarshal([]byte(v),&q)
		if um_err != nil {
			e := fmt.Sprintf("unmarshal Query: %v",um_err)
			t.Errorf(e)
		}
		_,jerr := json.Marshal(q)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}

	}
}

func TestResponseUnmarshal(t *testing.T) {
	s := []string{
		`{
    "Count": 3,
    "Items": [
        {
            "LastPostedBy": {
                "S": "fred@example.com"
            },
            "ForumName": {
                "S": "Amazon DynamoDB"
            },
            "LastPostDateTime": {
                "S": "20130102054211"
            },
            "Tags": {
                "SS": ["Problem","Question"]
            }
        },
        {
            "LastPostedBy": {
                "S": "alice@example.com"
            },
            "ForumName": {
                "S": "Amazon DynamoDB"
            },
            "LastPostDateTime": {
                    "S": "20130105111307"
            },
            "Tags": {
                "SS": ["Idea"]
            }
        },
        {
            "LastPostedBy": {
                "S": "bob@example.com"
            },
            "ForumName": {
                "S": "Amazon DynamoDB"
            },
            "LastPostDateTime": {
                "S": "20130108094417"
            },
            "Tags": {
                "SS": ["AppDesign", "HelpMe"]
            }
        }
    ]
}`,
`{
    "Count":17
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
