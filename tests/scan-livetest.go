package main

import (
	"fmt"
	"encoding/json"
	ep "github.com/smugmug/godynamo/endpoint"
	scan "github.com/smugmug/godynamo/endpoints/scan"
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

	s := scan.NewScan()
	tn := "test-godynamo-livetest"
	s.TableName = tn
	k_v1 := fmt.Sprintf("AHashKey%d",100)
	s.ScanFilter["TheHashKey"] =
		scan.ScanFilter{AttributeValueList:[]ep.AttributeValue{ep.AttributeValue{S:k_v1}},
		ComparisonOperator:scan.OP_EQ}
	jsonstr,_ := json.Marshal(s)
	fmt.Printf("JSON:%s\n",string(jsonstr))
	body,code,err := s.EndpointReq()
	var r scan.Response
	um_err := json.Unmarshal([]byte(body),&r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v",um_err)
		fmt.Printf("%s\n",e)
	}

	fmt.Printf("%v\n%v\n%v\n",body,code,err)
	s = scan.NewScan()
	s.TableName = tn
	jsonstr,_ = json.Marshal(s)
	fmt.Printf("JSON:%s\n",string(jsonstr))
	body,code,err = s.EndpointReq()
	fmt.Printf("%v\n%v\n%v\n",body,code,err)
	um_err = json.Unmarshal([]byte(body),&r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v",um_err)
		fmt.Printf("%s\n",e)
	}

	s.ScanFilter["SomeValue"] =
		scan.ScanFilter{AttributeValueList:[]ep.AttributeValue{
		ep.AttributeValue{N:"270"},ep.AttributeValue{N:"290"}},
		ComparisonOperator:scan.OP_BETWEEN}
	jsonstr,_ = json.Marshal(s)
	fmt.Printf("JSON:%s\n",string(jsonstr))
	body,code,err = s.EndpointReq()
	fmt.Printf("%v\n%v\n%v\n",body,code,err)
	um_err = json.Unmarshal([]byte(body),&r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v",um_err)
		fmt.Printf("%s\n",e)
	}
	k_v2 := fmt.Sprintf("AHashKey%d",290)
	r_v2 := fmt.Sprintf("%d",290)
	s.ExclusiveStartKey["TheHashKey"] = ep.AttributeValue{S:k_v2}
	s.ExclusiveStartKey["TheRangeKey"] = ep.AttributeValue{N:r_v2}
	jsonstr,_ = json.Marshal(s)
	fmt.Printf("JSON:%s\n",string(jsonstr))
	body,code,err = s.EndpointReq()
	fmt.Printf("%v\n%v\n%v\n",body,code,err)
	um_err = json.Unmarshal([]byte(body),&r)
	if um_err != nil {
		e := fmt.Sprintf("unmarshal Response: %v",um_err)
		fmt.Printf("%s\n",e)
	}
}
