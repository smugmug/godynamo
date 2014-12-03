package main

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	get "github.com/smugmug/godynamo/endpoints/get_item"
	keepalive "github.com/smugmug/godynamo/keepalive"
	"github.com/smugmug/godynamo/types/attributevalue"
	"log"
	"net/http"
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
		_ = <-iam_ready_chan
	}
	conf.Vals.ConfLock.RUnlock()

	get1 := get.NewGetItem()
	get1.TableName = "test-godynamo-livetest"
	// make sure this item has actually been inserted previously
	get1.Key["TheHashKey"] = &attributevalue.AttributeValue{S: "AHashKey264"}
	get1.Key["TheRangeKey"] = &attributevalue.AttributeValue{N: "264"}
	body, code, err := get1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("get failed %d %v %s\n", code, err, body)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)

	resp := get.NewResponse()
	um_err := json.Unmarshal([]byte(body), resp)
	if um_err != nil {
		log.Fatal(um_err)
	}
	j, jerr := json.Marshal(resp)
	if jerr != nil {
		log.Fatal(jerr)
	}
	fmt.Printf("RESP:%s\n", string(j))

	// Try converting the Response to a ResponseItemJSON
	c, cerr := resp.ToResponseItemJSON()
	if cerr != nil {
		log.Fatal(cerr)
	}
	jc, jcerr := json.Marshal(c)
	if jcerr != nil {
		log.Fatal(jcerr)
	}
	fmt.Printf("JSON:%s\n", string(jc))
}
