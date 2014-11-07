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
	for _, v := range s {
		var a ItemCollectionMetrics
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
