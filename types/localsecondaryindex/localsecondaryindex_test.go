package localsecondaryindex


import (
	"fmt"
	"testing"
	"encoding/json"
)

func TestLocalSecondaryIndexMarshal(t *testing.T) {
	s := []string{
		`{"IndexName":"LastPostIndex","KeySchema":[{"AttributeName":"ForumName","KeyType":"HASH"},{"AttributeName":"LastPostDateTime","KeyType":"RANGE"}],"Projection":{"ProjectionType":"KEYS_ONLY"}}`,
	}
	for _,v := range s {
		var a LocalSecondaryIndex
		um_err := json.Unmarshal([]byte(v),&a)
		if um_err != nil {
			fmt.Printf("%v\n",um_err)
			t.Errorf("cannot unmarshal\n")
		}

		json,jerr := json.Marshal(a)
		if jerr != nil {
			fmt.Printf("%v\n",jerr)
			t.Errorf("cannot marshal\n")
		}
		fmt.Printf("IN:%v, OUT:%v\n",v,string(json))
	}
}
