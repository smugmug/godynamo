package capacity

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestCapacityMarshal(t *testing.T) {
	s := []string{
		`{"CapacityUnits":1,"TableName":"mytable"}`,
		`{"CapacityUnits":1.01,"TableName":"mytable","Table":{"CapacityUnits":2.01}}`,
		`{"CapacityUnits":1.01,"TableName":"mytable","Table":{"CapacityUnits":2.01},"LocalSecondaryIndexes":{"mylsi":{"CapacityUnits":10.1}},"GlobalSecondaryIndexes":{"mygsi0":{"CapacityUnits":11.11},"mygsi1":{"CapacityUnits":10.1}}}`,
	}
	for i, v := range s {
		var a ConsumedCapacity
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}

		json, jerr := json.Marshal(a)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		if len(s[i]) != len(string(json)) {
			e := fmt.Sprintf("%s\n%s\nshould be same",s[i],string(json))
			t.Errorf(e)
		}
	}
}
