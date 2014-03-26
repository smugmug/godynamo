package main

import (
	"fmt"
	"net/http"
	ep "github.com/smugmug/godynamo/endpoint"
	get "github.com/smugmug/godynamo/endpoints/get_item"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	"log"
	keepalive "github.com/smugmug/godynamo/keepalive"
)

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

	for i := 1; i <= 300; i++ {
		var get1 get.Request
		get1.TableName = "test.dynamo-new-api"
		get1.Key = make(ep.Item)
		k := fmt.Sprintf("AHashKey%d",i)
		v := fmt.Sprintf("%d",i)
		get1.Key["TheHashKey"] = ep.AttributeValue{S:k}
		get1.Key["TheRangeKey"] = ep.AttributeValue{N:v}
		body,code,err := get1.EndpointReq()
		if err != nil || code != http.StatusOK {
			fmt.Printf("get failed %d %v %s\n",code,err,body)
		}
		fmt.Printf("%s\n",string(body))
	}
}
