package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	delete_item "github.com/smugmug/godynamo/endpoints/delete_item"
	get_item "github.com/smugmug/godynamo/endpoints/get_item"
	put_item "github.com/smugmug/godynamo/endpoints/put_item"
	update_item "github.com/smugmug/godynamo/endpoints/update_item"
	keepalive "github.com/smugmug/godynamo/keepalive"
	"github.com/smugmug/godynamo/types/attributevalue"
	"log"
	"net/http"
	"os"
)

// this test program runs a bunch of item-oriented operations.

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
	tablename1 := tn
	fmt.Printf("tablename1: %s\n", tablename1)

	var code int
	var err error
	var body []byte

	// INSERT SINGLE ITEM
	hk := "a-hash-key"
	rk := "1"
	put1 := put_item.NewPutItem()
	put1.TableName = tablename1
	av1 := attributevalue.NewAttributeValue()
	av2 := attributevalue.NewAttributeValue()
	av3 := attributevalue.NewAttributeValue()
	av4 := attributevalue.NewAttributeValue()
	av5 := attributevalue.NewAttributeValue()
	av6 := attributevalue.NewAttributeValue()
	av7 := attributevalue.NewAttributeValue()
	av1.S = hk
	av2.N = rk
	av3.InsertSS("pk1_a")
	av3.InsertSS("pk1_c")
	av4.InsertNS("1")
	av4.InsertNS_float64(2)
	av4.InsertNS("3")
	av4.InsertNS("-7.2432342")
	av5.N = "1"
	av6.B = base64.StdEncoding.EncodeToString([]byte("hello"))
	av7.InsertBS(base64.StdEncoding.EncodeToString([]byte("hello")))
	av7.InsertBS(base64.StdEncoding.EncodeToString([]byte("there")))
	put1.Item["TheHashKey"] = av1
	put1.Item["TheRangeKey"] = av2
	put1.Item["stringlist"] = av3
	put1.Item["numlist"] = av4
	put1.Item["num"] = av5
	put1.Item["byte"] = av6
	put1.Item["bytelist"] = av7

	body, code, err = put1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("put1 failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)

	// GET THAT ITEM
	get1 := get_item.NewGetItem()
	get1.TableName = tablename1
	get1.Key["TheHashKey"] = av1
	get1.Key["TheRangeKey"] = av2

	get_json, get_json_err := json.Marshal(get1)
	if get_json_err != nil {
		fmt.Printf("%v\n", get_json_err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(get_json))

	body, code, err = get1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("get failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)

	var gr get_item.Response
	um_err := json.Unmarshal([]byte(body), &gr)
	if um_err != nil {
		fmt.Printf("get resp unmarshal failed %s\n", um_err.Error())
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(body))

	// UPDATE THAT ITEM
	up1 := update_item.NewUpdateItem()
	new_attr_val := "new string here"
	up1.TableName = tablename1
	up1.Key["TheHashKey"] = av1
	up1.Key["TheRangeKey"] = av2

	up1.AttributeUpdates = attributevalue.NewAttributeValueUpdateMap()
	up_avu := attributevalue.NewAttributeValueUpdate()
	up_avu.Action = update_item.ACTION_PUT
	up_avu.Value = &attributevalue.AttributeValue{S: new_attr_val}
	up1.AttributeUpdates["new_string"] = up_avu

	del_avu := attributevalue.NewAttributeValueUpdate()
	del_avu.Action = update_item.ACTION_DEL
	del_avu.Value = attributevalue.NewAttributeValue()
	del_avu.Value.InsertSS("pk1_a")
	up1.AttributeUpdates["stringlist"] = del_avu

	del2_avu := attributevalue.NewAttributeValueUpdate()
	del2_avu.Action = update_item.ACTION_DEL
	del2_avu.Value = &attributevalue.AttributeValue{}
	up1.AttributeUpdates["byte"] = del2_avu

	add_avu := attributevalue.NewAttributeValueUpdate()
	add_avu.Action = update_item.ACTION_ADD
	add_avu.Value = &attributevalue.AttributeValue{N: "4"}
	up1.AttributeUpdates["num"] = add_avu

	up1.ReturnValues = update_item.RETVAL_ALL_NEW

	update_item_json, update_item_json_err := json.Marshal(up1)
	if update_item_json_err != nil {
		fmt.Printf("%v\n", update_item_json_err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(update_item_json))

	body, code, err = up1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("update item failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)

	var ur update_item.Response
	um_err = json.Unmarshal([]byte(body), &ur)
	if um_err != nil {
		fmt.Printf("update resp unmarshal failed %s\n", um_err.Error())
		os.Exit(1)
	}

	// GET IT AGAIN
	body, code, err = get1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("get failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)

	// DELETE THE ITEM
	del1 := delete_item.NewDeleteItem()
	del1.TableName = tablename1
	del1.Key["TheHashKey"] = av1
	del1.Key["TheRangeKey"] = av2

	del1.ReturnValues = delete_item.RETVAL_ALL_OLD
	body, code, err = del1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("delete item failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)

	fmt.Printf("PASSED\n")
}
