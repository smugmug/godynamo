package main

import (
	"fmt"
	"time"
	"net/http"
	ep "github.com/smugmug/godynamo/endpoint"
	put "github.com/smugmug/godynamo/endpoints/put_item"
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

	var put1 put.Request
	put1.TableName = "test-godynamo-livetest"
	k := fmt.Sprintf("hk1")
	v := fmt.Sprintf("%v",time.Now().Unix())
	put1.Item = make(ep.Item)
	put1.Item["TheHashKey"] = ep.AttributeValue{S:k}
	put1.Item["TheRangeKey"] = ep.AttributeValue{N:v}
	n := fmt.Sprintf("%v",time.Now().Unix())
	put1.Item["Mtime"] = ep.AttributeValue{N:n}
	put1.Item["SomeJunk"] = ep.AttributeValue{S:"some junk"}
	put1.Item["SomeJunks"] = ep.AttributeValue{SS:[]string{"some junk1","some junk2"}}
	body,code,err := put1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("put failed %d %v %s\n",code,err,body)
	}
	fmt.Printf("%s\n",string(body))
}
