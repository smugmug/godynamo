package main

import (
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	batch_write_item "github.com/smugmug/godynamo/endpoints/batch_write_item"
	keepalive "github.com/smugmug/godynamo/keepalive"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/item"
	"log"
	"net/http"
)

// this tests "RetryBatchWrite", which does NOT do intelligent splitting and re-assembling
// of requests and responses
func Test1() {
	tn := "test-godynamo-livetest"
	b := batch_write_item.NewBatchWriteItem()
	b.RequestItems[tn] = make([]batch_write_item.RequestInstance, 0)
	for i := 1; i <= 300; i++ {
		var p batch_write_item.PutRequest
		p.Item = item.NewItem()
		k := fmt.Sprintf("AHashKey%d", i)
		v := fmt.Sprintf("%d", i)
		p.Item["TheHashKey"] = &attributevalue.AttributeValue{S: k}
		p.Item["TheRangeKey"] = &attributevalue.AttributeValue{N: v}
		p.Item["SomeValue"] = &attributevalue.AttributeValue{N: v}
		b.RequestItems[tn] =
			append(b.RequestItems[tn],
				batch_write_item.RequestInstance{PutRequest: &p})
	}
	bs, _ := batch_write_item.Split(b)
	for _, bsi := range bs {
		body, code, err := bsi.RetryBatchWrite(0)
		if err != nil || code != http.StatusOK {
			fmt.Printf("error: %v\n%v\n%v\n", body, code, err)
		} else {
			fmt.Printf("worked!: %v\n%v\n%v\n", body, code, err)
		}
	}
}

// this tests "DoBatchGet", which breaks up requests that are larger than the limit
// and re-assembles responses
func Test2() {
	b := batch_write_item.NewBatchWriteItem()
	tn := "test-godynamo-livetest"
	b.RequestItems[tn] = make([]batch_write_item.RequestInstance, 0)
	for i := 201; i <= 300; i++ {
		var p batch_write_item.PutRequest
		p.Item = item.NewItem()
		k := fmt.Sprintf("AHashKey%d", i)
		v := fmt.Sprintf("%d", i)
		p.Item["TheHashKey"] = &attributevalue.AttributeValue{S: k}
		p.Item["TheRangeKey"] = &attributevalue.AttributeValue{N: v}
		p.Item["SomeValue"] = &attributevalue.AttributeValue{N: v}
		b.RequestItems[tn] =
			append(b.RequestItems[tn],
				batch_write_item.RequestInstance{PutRequest: &p})
	}
	body, code, err := b.DoBatchWrite()
	fmt.Printf("%v\n%v\n%v\n", body, code, err)
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
		_ = <-iam_ready_chan
	}
	conf.Vals.ConfLock.RUnlock()

	Test1()
	Test2()
}
