package main

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	batch_get_item "github.com/smugmug/godynamo/endpoints/batch_get_item"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/item"
	"net/http"
	"os"
)

// this tests "RetryBatchGet", which does NOT do intelligent splitting and re-assembling
// of requests and responses
func Test1() {
	home := os.Getenv("HOME")
	home_conf_file := home + string(os.PathSeparator) + "." + conf.CONF_NAME
	home_conf, home_conf_err := conf_file.ReadConfFile(home_conf_file)
	if home_conf_err != nil {
		panic("cannot read conf from " + home_conf_file)
	}
	home_conf.ConfLock.RLock()
	if home_conf.Initialized == false {
		panic("conf struct has not been initialized")
	}

	b := batch_get_item.NewBatchGetItem()
	tn := "test-godynamo-livetest"
	b.RequestItems[tn] = batch_get_item.NewRequestInstance()
	for i := 1; i <= 200; i++ {
		item := item.NewItem()
		k := fmt.Sprintf("AHashKey%d", i)
		v := fmt.Sprintf("%d", i)
		item["TheHashKey"] = &attributevalue.AttributeValue{S: k}
		item["TheRangeKey"] = &attributevalue.AttributeValue{N: v}
		b.RequestItems[tn].Keys =
			append(b.RequestItems[tn].Keys, item)

	}
	_, jerr := json.Marshal(b)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
	} else {
		//fmt.Printf("%s\n",string(json))
	}
	bs, _ := batch_get_item.Split(b)
	for _, bsi := range bs {
		body, code, err := bsi.RetryBatchGetWithConf(0, home_conf)
		if err != nil || code != http.StatusOK {
			fmt.Printf("error: %v\n%v\n%v\n", string(body), code, err)
		} else {
			fmt.Printf("worked!: %v\n%v\n%v\n", string(body), code, err)
		}
	}
}

// this tests "DoBatchGet", which breaks up requests that are larger than the limit
// and re-assembles responses
func Test2() {
	home := os.Getenv("HOME")
	home_conf_file := home + string(os.PathSeparator) + "." + conf.CONF_NAME
	home_conf, home_conf_err := conf_file.ReadConfFile(home_conf_file)
	if home_conf_err != nil {
		panic("cannot read conf from " + home_conf_file)
	}
	home_conf.ConfLock.RLock()
	if home_conf.Initialized == false {
		panic("conf struct has not been initialized")
	}

	b := batch_get_item.NewBatchGetItem()
	tn := "test-godynamo-livetest"
	b.RequestItems[tn] = batch_get_item.NewRequestInstance()
	for i := 1; i <= 300; i++ {
		item := item.NewItem()
		k := fmt.Sprintf("AHashKey%d", i)
		v := fmt.Sprintf("%d", i)
		item["TheHashKey"] = &attributevalue.AttributeValue{S: k}
		item["TheRangeKey"] = &attributevalue.AttributeValue{N: v}
		b.RequestItems[tn].Keys =
			append(b.RequestItems[tn].Keys, item)

	}
	_, jerr := json.Marshal(*b)
	if jerr != nil {
		fmt.Printf("%v\n", jerr)
	} else {
		//fmt.Printf("%s\n",string(json))
	}
	body, code, err := b.DoBatchGetWithConf(home_conf)
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)
}

func main() {
	Test1()
	Test2()
}
