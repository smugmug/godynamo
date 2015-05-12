package main

import (
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	batch_write_item "github.com/smugmug/godynamo/endpoints/batch_write_item"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/item"
	"net/http"
	"os"
)

// these tests are just like batch_write_item-livestest except they use a parameterized conf

// this tests "RetryBatchWrite", which does NOT do intelligent splitting and re-assembling
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
		body, code, err := bsi.RetryBatchWriteWithConf(0, home_conf)
		if err != nil || code != http.StatusOK {
			fmt.Printf("error: %v\n%v\n%v\n", string(body), code, err)
		} else {
			fmt.Printf("worked!: %v\n%v\n%v\n", string(body), code, err)
		}
	}
}

// this tests "DoBatchWrite", which breaks up requests that are larger than the limit
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
	body, code, err := b.DoBatchWriteWithConf(home_conf)
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)
}

func main() {
	Test1()
	Test2()
}
