package main

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	scan "github.com/smugmug/godynamo/endpoints/scan"
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

	s := scan.NewScan()
	tn := "test-godynamo-livetest"
	s.TableName = tn
	k_v1 := fmt.Sprintf("AHashKey%d", 100)

	kc := condition.NewCondition()
	kc.AttributeValueList = make([]*attributevalue.AttributeValue, 1)
	kc.AttributeValueList[0] = &attributevalue.AttributeValue{S: k_v1}
	kc.ComparisonOperator = scan.OP_EQ

	s.ScanFilter["TheHashKey"] = kc
	jsonstr, _ := json.Marshal(s)
	fmt.Printf("JSON:%s\n", string(jsonstr))
	body, code, err := s.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("scan failed %d %v %s\n", code, err, body)
	}
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)

	var r scan.Response
	um_err := json.Unmarshal([]byte(body), &r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v", um_err)
		fmt.Printf("%s\n", e)
	}

	s = scan.NewScan()
	s.TableName = tn
	jsonstr, _ = json.Marshal(s)
	fmt.Printf("JSON:%s\n", string(jsonstr))
	body, code, err = s.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("scan failed %d %v %s\n", code, err, body)
	}
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)

	um_err = json.Unmarshal([]byte(body), &r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v", um_err)
		fmt.Printf("%s\n", e)
	}

	kc = condition.NewCondition()
	kc.AttributeValueList = make([]*attributevalue.AttributeValue, 2)
	kc.AttributeValueList[0] = &attributevalue.AttributeValue{N: "270"}
	kc.AttributeValueList[1] = &attributevalue.AttributeValue{N: "290"}
	kc.ComparisonOperator = scan.OP_BETWEEN
	s.ScanFilter["SomeValue"] = kc

	jsonstr, _ = json.Marshal(s)
	fmt.Printf("JSON:%s\n", string(jsonstr))
	body, code, err = s.EndpointReq()
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)
	um_err = json.Unmarshal([]byte(body), &r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v", um_err)
		fmt.Printf("%s\n", e)
	}
	k_v2 := fmt.Sprintf("AHashKey%d", 290)
	r_v2 := fmt.Sprintf("%d", 290)
	s.ExclusiveStartKey["TheHashKey"] = &attributevalue.AttributeValue{S: k_v2}
	s.ExclusiveStartKey["TheRangeKey"] = &attributevalue.AttributeValue{N: r_v2}
	jsonstr, _ = json.Marshal(s)
	fmt.Printf("JSON:%s\n", string(jsonstr))
	body, code, err = s.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("scan failed %d %v %s\n", code, err, body)
	}
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)

	um_err = json.Unmarshal([]byte(body), &r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v", um_err)
		fmt.Printf("%s\n", e)
	}
}
