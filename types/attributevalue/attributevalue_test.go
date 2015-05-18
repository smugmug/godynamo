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
		`{"BS":["d2VsbCBoZWxsbyB0aGVyZQ==","aGkgdGhlcmUK","aG93ZHk="]}`,
		`{"L":[{"S":"a string"},{"L":[{"S":"another string"}]}]}`,
	}
	for _, v := range s {
		var a AttributeValue
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}

		json, jerr := json.Marshal(a)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		o := string(json)
		if (len(v) != len(o)) {
			e := fmt.Sprintf("%s\n%s\nshould be equal through round-trip",v,o)
			t.Errorf(e)
		}
	}
}

func TestAttributeValueMarshalDeDuplicate(t *testing.T) {
	s := []string{
		`{"SS":["a","b","c","c","c"]}`,
		`{"NS":["42","1","0","0","1","1","1","42"]}`,
		`{"M":{"key1":{"S":"a string"},"key2":{"L":[{"NS":["42","42","1"]},{"S":"a string"},{"L":[{"S":"another string"}]}]}}}`,
	}
	c := []string{
		`{"SS":["a","b","c"]}`,
		`{"NS":["42","1","0"]}`,
		`{"M":{"key1":{"S":"a string"},"key2":{"L":[{"NS":["42","1"]},{"S":"a string"},{"L":[{"S":"another string"}]}]}}}`,
	}
	for i, v := range s {
		var a AttributeValue
		um_err := json.Unmarshal([]byte(v), &a)
		if um_err != nil {
			t.Errorf("cannot unmarshal\n")
		}
		
		json, jerr := json.Marshal(a)
		if jerr != nil {
			t.Errorf("cannot marshal\n")
		}
		o := string(json)
		if (len(c[i]) != len(o)) {
			e := fmt.Sprintf("%s\n%s\nshould be equal through round-trip",c[i],o)
			t.Errorf(e)
		}
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
		} else {
			_ = fmt.Sprintf("%v\n", jerr)
		}
	}
	a = NewAttributeValue()
	a.N = "1"
	a.B = "fsdfa"
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
		}
	}

	a = NewAttributeValue()
	a.N = "1"
	a.InsertSS("a")
	if a.Valid() {
		_, jerr := json.Marshal(a)
		if jerr == nil {
			t.Errorf("should not have been able to marshal\n")
		}
	}

	a = NewAttributeValue()
	if a.Valid() {
		t.Errorf("should not pass valid\n")
	}
}

// Empty AttributeValues should emit null
func TestAttributeValueEmpty(t *testing.T) {
	a := NewAttributeValue()
	_, jerr := json.Marshal(a)
	if jerr != nil {
		t.Errorf("cannot marshal\n")
	}

	var a2 AttributeValue
	_, jerr = json.Marshal(a2)
	if jerr != nil {
		t.Errorf("cannot marshal\n")
	}
}

// Test the Insert funtions
func TestAttributeValueInserts(t *testing.T) {
	a1 := NewAttributeValue()
	a1.InsertSS("hi")
	a1.InsertSS("hi") // duplicate, should be removed
	a1.InsertSS("bye")
	json, jerr := json.Marshal(a1)
	if jerr != nil {
		t.Errorf("cannot marshal\n")
	}
	c := `{"SS":["hi","bye"]}`
	if len(c) != len(string(json)) {
		e := fmt.Sprintf("%s\n%s\nshould be same",c,string(json))
		t.Errorf(e)
	}
}

// Test the Insert functions
func TestAttributeValueInserts2(t *testing.T) {
	a1 := NewAttributeValue()
	a1.InsertSS("hi")
	a1.InsertSS("hi") // duplicate, should be removed
	a1.InsertSS("bye")

	a2 := NewAttributeValue()
	_ = a2.InsertL(a1)
	a1 = nil // should be fine, above line should make a new copy
	json, jerr := json.Marshal(a2)
	if jerr != nil {
		t.Errorf("cannot marshal\n")
	}
	c := `{"L":[{"SS":["hi","bye"]}]}`
	if len(c) != len(string(json)) {
		e := fmt.Sprintf("%s\n%s\nshould be same",c,string(json))
		t.Errorf(e)
	}

	a3 := NewAttributeValue()
	nerr := a3.InsertN("fred")
	if nerr == nil {
		t.Errorf("should have returned error from InsertN\n")
	}
	berr := a3.InsertB("1")
	if berr == nil {
		t.Errorf("should have returned error from InsertB\n")
	}
}

func TestCopy(t *testing.T) {
	a1 := NewAttributeValue()
	_ = a1.InsertSS("hi")
	_ = a1.InsertSS("bye")

	var a2 = new(AttributeValue)

	cp_err := a1.Copy(a2)
	if cp_err != nil {
		t.Errorf("Copy error\n")
	}

	a2 = nil
	cp_err = a1.Copy(a2)
	if cp_err == nil {
		t.Errorf("should have returned error from Copy\n")
	}

	a1 = nil
	a2 = new(AttributeValue)
	cp_err = a1.Copy(a2)
	if cp_err == nil {
		t.Errorf("should have returned error from Copy\n")
	}

	a2 = nil
	cp_err = a1.Copy(a2)
	if cp_err == nil {
		t.Errorf("should have returned error from Copy\n")
	}
}

// Make sure Valid emits as null
func TestAttributeValueUpdate(t *testing.T) {
	a := NewAttributeValueUpdate()
	a.Action = "DELETE"
	_, jerr := json.Marshal(a)
	if jerr != nil {
		t.Errorf("cannot marshal\n")
	}
}

func TestCoerceAttributeValueBasicJSON(t *testing.T) {
	js := []string{`{"a":{"b":"c"},"d":[{"e":"f"},"g","h"],"i":[1.0,2.0,3.0],"j":["x","y"]}`,
		`"a"`, `true`,
		`[1,2,3,2,3]`}
	c := []string{`{"a":{"b":"c"},"d":[{"e":"f"},"g","h"],"i":[1,2,3],"j":["x","y"]}`,
		`"a"`, `true`,
		`[1,2,3]`}
	for i, v := range js {
		j := []byte(v)
		av, av_err := BasicJSONToAttributeValue(j)
		if av_err != nil {
			t.Errorf("cannot coerce")
		}
		_, av_json_err := json.Marshal(av)
		if av_json_err != nil {
			_ = fmt.Sprintf("%v\n", av_json_err)
			t.Errorf("cannot marshal")
		}
		b, cerr := av.ToBasicJSON()
		if cerr != nil {
			t.Errorf("cannot coerce")
		}
		if len(c[i]) != len(string(b)) {
			e := fmt.Sprintf("%s\n%s\nshould be same",c[i],string(b))
			t.Errorf(e)
		}
	}
}

func TestCoerceAttributeValueMapBasicJSON(t *testing.T) {
	js := []string{`{"AS":"1234string","AN":3,"ANS":[1,2,1,2,3],"ASS":["a","a","b"],"ABOOL":true,"AL":["1234string",3,[1,2,3],["a","b"]],"AM":{"AMS":"1234string","AMN":3,"AMNS":[1,2,3],"AMSS":["a","b"],"AMBOOL":true,"AL":["1234string",3,[1,2,3],["a","b"]]}}`}
	c := []string{`{"AS":"1234string","AN":3,"ANS":[1,2,3],"ASS":["a","b"],"ABOOL":true,"AL":["1234string",3,[1,2,3],["a","b"]],"AM":{"AMS":"1234string","AMN":3,"AMNS":[1,2,3],"AMSS":["a","b"],"AMBOOL":true,"AL":["1234string",3,[1,2,3],["a","b"]]}}`}
	for i, v := range js {
		j := []byte(v)
		av, av_err := BasicJSONToAttributeValueMap(j)
		if av_err != nil {
			t.Errorf("cannot coerce")
		}
		_, av_json_err := json.Marshal(av)
		if av_json_err != nil {
			t.Errorf("cannot marshal")
		}
		b, cerr := av.ToBasicJSON()
		if cerr != nil {
			t.Errorf("cannot coerce")
		}
		if len(c[i]) != len(string(b)) {
			e := fmt.Sprintf("%s\n%s\nshould be same",c[i],string(b))
			t.Errorf(e)
		}
	}
}
