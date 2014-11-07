package attributevalue

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Round trip some data. Note the sets have repeated elements...make sure they
// are eliminated
func TestAttributeValueMarshal(t *testing.T) {
	s := []string{
		`{"S":"a string"}`,
		`{"B":"aGkgdGhlcmUK"}`,
		`{"N":"5"}`,
		`{"BOOL":true}`,
		`{"NULL":false}`,
		`{"SS":["a","b","c","c","c"]}`,
		`{"BS":["aGkgdGhlcmUK","aG93ZHk=","d2VsbCBoZWxsbyB0aGVyZQ=="]}`,
		`{"NS":["42","1","0","0","1","1","1","42"]}`,
		`{"L":[{"S":"a string"},{"L":[{"S":"another string"}]}]}`,
		`{"M":{"key1":{"S":"a string"},"key2":{"L":[{"NS":["42","42","1"]},{"S":"a string"},{"L":[{"S":"another string"}]}]}}}`,
	}
	for _, v := range s {
		fmt.Printf("--------\n")
		fmt.Printf("IN:%v\n", v)
		var a AttributeValue
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			fmt.Printf("%v\n", um_err)
			t.Errorf("cannot unmarshal\n")
		}

		json, jerr := json.Marshal(a)
		if jerr != nil {
			fmt.Printf("%v\n", jerr)
			t.Errorf("cannot marshal\n")
			return
		}
		fmt.Printf("OUT:%v\n", string(json))
	}
}

// Demonstrate the use of the Valid function
func TestAttributeValueInvalid(t *testing.T) {
	a := NewAttributeValue()
	a.N = "1"
	a.S = "a"
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
			return
		} else {
			fmt.Printf("%v\n", jerr)
		}
	}
	a = NewAttributeValue()
	a.N = "1"
	a.B = "fsdfa"
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
			return
		} else {
			fmt.Printf("%v\n", jerr)
		}
	}

	a = NewAttributeValue()
	a.N = "1"
	a.InsertSS("a")
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
			return
		} else {
			fmt.Printf("%v\n", jerr)
		}
	}
}

// Empty AttributeValues should emit null
func TestAttributeValueEmpty(t *testing.T) {
	a := NewAttributeValue()
	json_bytes, jerr := json.Marshal(a)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	fmt.Printf("OUT:%v\n", string(json_bytes))

	var a2 AttributeValue
	json_bytes, jerr = json.Marshal(a2)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	fmt.Printf("OUT:%v\n", string(json_bytes))
}

// Test the Insert funtions
func TestAttributeValueInserts(t *testing.T) {
	a1 := NewAttributeValue()
	a1.InsertSS("hi")
	a1.InsertSS("hi") // duplicate, should be removed
	a1.InsertSS("bye")
	json, jerr := json.Marshal(a1)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	fmt.Printf("OUT:%v\n", string(json))
}

// Test the Insert functions
func TestAttributeValueInserts2(t *testing.T) {
	a1 := NewAttributeValue()
	_ = a1.InsertSS("hi")
	_ = a1.InsertSS("hi") // duplicate, should be removed
	_ = a1.InsertSS("bye")
	json_bytes, jerr := json.Marshal(a1)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	fmt.Printf("OUT:%v\n", string(json_bytes))

	a2 := NewAttributeValue()
	_ = a2.InsertL(a1)
	a1 = nil // should be fine, above line should make a new copy
	json_bytes, jerr = json.Marshal(a2)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	fmt.Printf("OUT:%v\n", string(json_bytes))

	a3 := NewAttributeValue()
	nerr := a3.InsertN("fred")
	if nerr == nil {
		t.Errorf("should have returned error from InsertN\n")
		return
	} else {
		fmt.Printf("%v\n", nerr)
	}
	berr := a3.InsertB("1")
	if berr == nil {
		t.Errorf("should have returned error from InsertB\n")
		return
	} else {
		fmt.Printf("%v\n", berr)
	}
}

// Should fail, a2 is uninitialized
func TestBadCopy(t *testing.T) {
	a1 := NewAttributeValue()
	_ = a1.InsertSS("hi")
	_ = a1.InsertSS("bye")

	var a2 = new(AttributeValue)

	cp_err := a1.Copy(a2)
	if a2 == nil {
		t.Errorf("should have returned error from Copy\n")
		return
	} else {
		fmt.Printf("%v\n", cp_err)
	}
}

// Make sure Valid emits as null
func TestAttributeValueUpdate(t *testing.T) {
	a := NewAttributeValueUpdate()
	a.Action = "DELETE"
	json_bytes, jerr := json.Marshal(a)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
		t.Errorf("cannot marshal\n")
		return
	}
	fmt.Printf("OUT:%v\n", string(json_bytes))

}
