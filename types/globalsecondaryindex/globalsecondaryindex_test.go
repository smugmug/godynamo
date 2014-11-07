package globalsecondaryindex

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestGlobalSecondaryIndexMarshal(t *testing.T) {
	s := []string{`{"IndexName":"LastPostIndex","KeySchema":[{"AttributeName": "ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType": "RANGE"}],"Projection": {"ProjectionType": "KEYS_ONLY"},"ProvisionedThroughput": {"ReadCapacityUnits":200,"WriteCapacityUnits":200}}`}
	for _, v := range s {
		var a GlobalSecondaryIndex
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			fmt.Printf("%v\n", um_err)
			t.Errorf("cannot unmarshal\n")
		}

		json, jerr := json.Marshal(a)
		if jerr != nil {
			fmt.Printf("%v\n", jerr)
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n", v, string(json))
	}
}
