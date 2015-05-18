package main

import (
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	put "github.com/smugmug/godynamo/endpoints/put_item"
	keepalive "github.com/smugmug/godynamo/keepalive"
	"github.com/smugmug/godynamo/types/attributevalue"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// this is the same as put-item-livetest except here we demonstrate using a parameterized conf
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

	// launch a background poller to keep conns to aws alive
	if home_conf.Network.DynamoDB.KeepAlive {
		log.Printf("launching background keepalive")
		go keepalive.KeepAlive([]string{})
	}

	// deal with iam, or not
	if home_conf.UseIAM {
		iam_ready_chan := make(chan bool)
		go conf_iam.GoIAM(iam_ready_chan)
		_ = <-iam_ready_chan
	}
	home_conf.ConfLock.RUnlock()
	
	put1 := put.NewPutItem()
	put1.TableName = "test-godynamo-livetest"
	
	k := fmt.Sprintf("hk1000")
	v := fmt.Sprintf("%v", time.Now().Unix())
	put1.Item["TheHashKey"] = &attributevalue.AttributeValue{S: k}
	put1.Item["TheRangeKey"] = &attributevalue.AttributeValue{N: v}

	i := fmt.Sprintf("%v",1)
	t := fmt.Sprintf("%v",time.Now().Unix())
	put1.Item["UserID"] = &attributevalue.AttributeValue{N:i}
	put1.Item["Timestamp"] = &attributevalue.AttributeValue{N:t}

	// the Token field is a simple string
	put1.Item["Token"] = &attributevalue.AttributeValue{S:"a token"}

	// the Location must be created as a "map"
	location := attributevalue.NewAttributeValue()
	location.InsertM("Latitude",&attributevalue.AttributeValue{N:"120.01"})
	location.InsertM("Longitude",&attributevalue.AttributeValue{N:"50.99"})
	put1.Item["Location"] = location

	body, code, err := put1.EndpointReqWithConf(home_conf)
	if err != nil || code != http.StatusOK {
		fmt.Printf("put failed %d %v %s\n", code, err, body)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)
}
