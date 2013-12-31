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

package delete_item

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
    "Expected": {
        "Replies": {
            "Exists": false
        }
    },
    "ReturnValues": "ALL_OLD"
}`,
	}
	for _,v := range s {
		var d Delete
		um_err := json.Unmarshal([]byte(v),&d)
		if um_err != nil {
			t.Errorf("cannot unmarshal to delete:\n" + v + "\n")
		}
		_,jerr := json.Marshal(d)
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
            "S": "fred@example.com"
        },
        "ForumName": {
            "S": "Amazon DynamoDB"
        },
        "LastPostDateTime": {
            "S": "201303201023"
        },
        "Tags": {
            "SS": ["Update","Multiple Items","HelpMe"]
        },
        "Subject": {
            "S": "How do I update multiple items?"
        },
        "Message": {
            "S": "I want to update multiple items in a single API call. What's the best way to do that?"
        }
    }
}`,
	}
	for _,v := range s {
		var d Response
		um_err := json.Unmarshal([]byte(v),&d)
		if um_err != nil {
			t.Errorf("cannot unmarshal to delete:\n" + v + "\n")
		}
		_,jerr := json.Marshal(d)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}

	}
}
