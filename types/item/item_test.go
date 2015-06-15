package item

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Roundtrip some examples
func TestItemMarshal(t *testing.T) {
	s := []string{
		`{"ItemName":{"S":"a string"}}`,
		`{"ItemName":{"B":"aGkgdGhlcmUK"}}`,
		`{"ItemName":{"N":"5"}}`,
		`{"ItemName":{"BOOL":true}}`,
		`{"ItemName":{"NULL":true}}`,
		`{"ItemName":{"SS":["a","b","c"]}}`,
		`{"ItemName":{"BS":["aGkgdGhlcmUK"]}}`,
		`{"ItemName":{"NS":["42","1","0"]}}`,
		`{"ItemName":{"L":[{"S":"a string"},{"L":[{"S":"another string"}]}]}}`,
		`{"ItemName":{"M":{"key1":{"S":"a string"},"key2":{"L":[{"S":"a string"},{"L":[{"S":"another string"}]}]}}}}`,
	}
	for i, v := range s {
		var a Item
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
		l := len(a)
		ac := NewItem()
		a.Copy(ac)
		lc := len(ac)
		if l != lc {
			e := fmt.Sprintf("lengths differ: %d %d",l,lc)
			t.Errorf(e)
		}
		fmt.Printf("%d %d\n",l,lc)
		delete(a,"ItemName")

		lc = len(ac)
		if l != lc {
			e := fmt.Sprintf("lengths differ: %d %d",l,lc)
			t.Errorf(e)
		}

	}
}
