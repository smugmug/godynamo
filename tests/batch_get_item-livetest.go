package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	batch_get_item "github.com/smugmug/godynamo/endpoints/batch_get_item"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	"github.com/smugmug/godynamo/types/item"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	"log"
	keepalive "github.com/smugmug/godynamo/keepalive"
)

// this tests "RetryBatchGet", which does NOT do intelligent splitting and re-assembling
// of requests and responses
func Test1() {
	b := batch_get_item.NewBatchGetItem()
	tn := "test-godynamo-livetest"
	b.RequestItems[tn] = batch_get_item.NewRequestInstance()
	for i := 1; i <= 200; i++ {
		item := item.NewItem()
		k := fmt.Sprintf("AHashKey%d",i)
		v := fmt.Sprintf("%d",i)
		item["TheHashKey"] = &attributevalue.AttributeValue{S:k}
		item["TheRangeKey"] = &attributevalue.AttributeValue{N:v}
		b.RequestItems[tn].Keys =
			append(b.RequestItems[tn].Keys,item)

	}
	_,jerr := json.Marshal(b)
	if jerr != nil {
		fmt.Printf("%v\n",jerr)
	} else {
		//fmt.Printf("%s\n",string(json))
	}
	bs,_ := batch_get_item.Split(b)
	for _,bsi := range bs {
	 	body,code,err := bsi.RetryBatchGet(0)
	 	if err != nil || code != http.StatusOK {
	 		fmt.Printf("error: %v\n%v\n%v\n",body,code,err)
	 	} else {
	 		fmt.Printf("worked!: %v\n%v\n%v\n",body,code,err)
		}
	}
}

// this tests "DoBatchGet", which breaks up requests that are larger than the limit
// and re-assembles responses
func Test2() {
	b := batch_get_item.NewBatchGetItem()
	tn := "test-godynamo-livetest"
	b.RequestItems[tn] = batch_get_item.NewRequestInstance()
	for i := 1; i <= 200; i++ {
		item := item.NewItem()
		k := fmt.Sprintf("AHashKey%d",i)
		v := fmt.Sprintf("%d",i)
		item["TheHashKey"] = &attributevalue.AttributeValue{S:k}
		item["TheRangeKey"] = &attributevalue.AttributeValue{N:v}
		b.RequestItems[tn].Keys =
			append(b.RequestItems[tn].Keys,item)

	}
	_,jerr := json.Marshal(*b)
	if jerr != nil {
		fmt.Printf("%v\n",jerr)
	} else {
		//fmt.Printf("%s\n",string(json))
	}
	body,code,err := b.DoBatchGet()
	fmt.Printf("%v\n%v\n%v\n",body,code,err)
}

func main() {
	conf_file.Read()
	conf.Vals.ConfLock.RLock()
	if conf.Vals.Initialized == false {
		panic("the conf.Vals global conf struct has not been initialized")
	}

	// launch a background poller to keep conns to aws alive
	if conf.Vals.Network.DynamoDB.KeepAlive {
		log.Printf("launching background keepalive")
		go keepalive.KeepAlive([]string{})
	}

	// deal with iam, or not
	if conf.Vals.UseIAM {
		iam_ready_chan := make(chan bool)
		go conf_iam.GoIAM(iam_ready_chan)
		_ = <- iam_ready_chan
	}
	conf.Vals.ConfLock.RUnlock()

	Test1()
	Test2()
}
