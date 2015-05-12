package main

import (
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	ep "github.com/smugmug/godynamo/endpoint"
	query "github.com/smugmug/godynamo/endpoints/query"
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
	body, code, err := q.EndpointReqWithConf(home_conf)
	if err != nil || code != http.StatusOK {
		fmt.Printf("query failed %d %v %s\n", code, err, body)
	}
	fmt.Printf("%v\n%v\n%v\n", string(body), code, err)
}
