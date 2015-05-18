package globalsecondaryindex

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestGlobalSecondaryIndexMarshal(t *testing.T) {
	s := []string{`{"IndexName":"LastPostIndex","KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType":"RANGE"}],"Projection":{"NonKeyAttributes":null,"ProjectionType":"KEYS_ONLY"},"ProvisionedThroughput":{"ReadCapacityUnits":200,"WriteCapacityUnits":200}}`}
	for i, v := range s {
		var a GlobalSecondaryIndex
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
