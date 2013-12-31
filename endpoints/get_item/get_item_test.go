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

package get_item

import (
	"testing"
	"encoding/json"
)

func TestRequestMarshal(t *testing.T) {
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
    "AttributesToGet": ["LastPostDateTime","Message","Tags"],
    "ConsistentRead": true,
    "ReturnConsumedCapacity": "TOTAL"
}`,
	}
	for _,v := range s {
		var g Get
		um_err := json.Unmarshal([]byte(v),&g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		_,jerr := json.Marshal(g)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}

func TestResponseMarshal(t *testing.T) {
	s := []string{
		`{
    "ConsumedCapacity": {
        "CapacityUnits": 1,
        "TableName": "Thread"
    },
    "Item": {
        "Tags": {
            "SS": ["Update","Multiple Items","HelpMe"]
        },
        "LastPostDateTime": {
            "S": "201303190436"
        },
        "Message": {
            "S": "I want to update multiple items in a single API call. What's the best way to do that?"
        }
    }
}`,
	}
	for _,v := range s {
		var g Response
		um_err := json.Unmarshal([]byte(v),&g)
		if um_err != nil {
			t.Errorf("cannot unmarshal to create:\n" + v + "\n")
		}
		_,jerr := json.Marshal(g)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
	}
}
