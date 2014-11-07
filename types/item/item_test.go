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
		`{"ItemName":{"BS":["aGkgdGhlcmUK","aGkgdGhlcmUK","aGkgdGhlcmUK"]}}`,
		`{"ItemName":{"NS":["42","1","0"]}}`,
		`{"ItemName":{"L":[{"S":"a string"},{"L":[{"S":"another string"}]}]}}`,
		`{"ItemName":{"M":{"key1":{"S":"a string"},"key2":{"L":[{"S":"a string"},{"L":[{"S":"another string"}]}]}}}}`,
	}
	for _, v := range s {
		var a Item
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
