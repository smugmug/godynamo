package expected

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestExpectedMarshal(t *testing.T) {
	s := []string{
		`{"MyConstraint1":{"AttributeValueList":[{"S":"a string"}],"ComparisonOperator":"BEGINS_WITH","Value":{"N":"4"},"Exists":true}}`,
		`{"MyConstraint2":{"AttributeValueList":[{"S":"a string"}],"ComparisonOperator":"BEGINS_WITH","Value":{"N":"4"},"Exists":false}}`,
	}
	for i, v := range s {
		var a Expected
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
			return
		}
	}
}
