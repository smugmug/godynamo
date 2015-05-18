package localsecondaryindex

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLocalSecondaryIndexMarshal(t *testing.T) {
	s := []string{
		`{"IndexName":"LastPostIndex","KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType":"RANGE"}],"Projection":{"ProjectionType":"KEYS_ONLY"}}`,
	}
	for i, v := range s {
		var a LocalSecondaryIndex
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		json, jerr := json.Marshal(a)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		if len(s[i]) != len(string(json)) {
			e := fmt.Sprintf("\n%s\n%s\nshould be same",s[i],string(json))
			t.Errorf(e)
		}
	}
}
