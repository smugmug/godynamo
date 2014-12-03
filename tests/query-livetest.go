package main

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	ep "github.com/smugmug/godynamo/endpoint"
	query "github.com/smugmug/godynamo/endpoints/query"
	keepalive "github.com/smugmug/godynamo/keepalive"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/condition"
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

	tn := "test-godynamo-livetest"
	q := query.NewQuery()
	q.TableName = tn
	q.Select = ep.SELECT_ALL
	k_v1 := fmt.Sprintf("AHashKey%d", 100)
	kc := condition.NewCondition()
	kc.AttributeValueList = make([]*attributevalue.AttributeValue, 1)
	kc.AttributeValueList[0] = &attributevalue.AttributeValue{S: k_v1}
	kc.ComparisonOperator = query.OP_EQ
	q.Limit = 10000
	q.KeyConditions["TheHashKey"] = kc
	json, _ := json.Marshal(q)
	fmt.Printf("JSON:%s\n", string(json))
	body, code, err := q.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("query failed %d %v %s\n", code, err, body)
	}
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)
}
