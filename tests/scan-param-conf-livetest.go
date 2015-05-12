package main

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	scan "github.com/smugmug/godynamo/endpoints/scan"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/condition"
	"net/http"
	"os"
)

func main() {

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
	body, code, err := s.EndpointReqWithConf(home_conf)
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
	body, code, err = s.EndpointReqWithConf(home_conf)
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
	body, code, err = s.EndpointReqWithConf(home_conf)
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
	body, code, err = s.EndpointReqWithConf(home_conf)
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
