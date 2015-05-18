package itemcollectionmetrics

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExpectedMarshal(t *testing.T) {
	s := []string{
		`{"ItemCollectionKey":{"AttributeValue":{"S":"a string"}},"SizeEstimateRangeGB":[0,10]}`,
	}
	for i, v := range s {
		var a ItemCollectionMetrics
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
