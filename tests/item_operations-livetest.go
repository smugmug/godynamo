// Copyright (c) 2013,2014 SmugMug, Inc. All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
//       copyright notice, this list of conditions and the following
//       disclaimer in the documentation and/or other materials provided
//       with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY SMUGMUG, INC. ``AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
// PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL SMUGMUG, INC. BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE
// GOODS OR SERVICES;LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
// IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
// OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
// ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"fmt"
	"os"
	"encoding/json"
        "net/http"
	"encoding/base64"
	put_item "github.com/smugmug/godynamo/endpoints/put_item"
	get_item "github.com/smugmug/godynamo/endpoints/get_item"
	update_item "github.com/smugmug/godynamo/endpoints/update_item"
	delete_item "github.com/smugmug/godynamo/endpoints/delete_item"
	ep "github.com/smugmug/godynamo/endpoint"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
)

// this test program runs a bunch of item-oriented operations.

func main() {
	// conf file must be read in before anything else, to initialize permissions etc
	conf_file.Read()
	if conf.Vals.Initialized == false {
		panic("the conf.Vals global conf struct has not been initialized")
	}

	// deal with iam, or not
	iam_ready_chan := make(chan bool)
	go conf_iam.GoIAM(iam_ready_chan)
	iam_ready := <- iam_ready_chan
	if iam_ready {
		fmt.Printf("using iam\n")
	} else {
		fmt.Printf("not using iam\n")
	}

	tn := "test-godynamo-livetest"
	tablename1 := tn
	fmt.Printf("tablename1: %s\n",tablename1)

	var code int
	var err error
	var body string

	// INSERT SINGLE ITEM
	hk := "a-hash-key"
	rk := "1"
        var put1 put_item.Request
        put1.TableName = tablename1
        var av1,av2,av3,av4,av5,av6,av7 ep.AttributeValue
        av1.S = hk
        av2.N = rk
        av3.SS = []string{"pk1_a","pk1_b","pk1_c"}
        av4.NS = []string{"1","2","3","-7.234234234234234e+09"}
        av5.N  = "1"
        av6.B  = base64.StdEncoding.EncodeToString([]byte("hello"))
        av7.BS = []string{base64.StdEncoding.EncodeToString([]byte("hello")),
                base64.StdEncoding.EncodeToString([]byte("there"))}
	put1.Item = make(ep.Item)
	put1.Item["TheHashKey"] = av1
	put1.Item["TheRangeKey"] = av2
	put1.Item["stringlist"] = av3
	put1.Item["numlist"] = av4
	put1.Item["num"] = av5
	put1.Item["byte"] = av6
	put1.Item["bytelist"] = av7
        // recommended to make sure pk does not exist yet in table
	// (i.e. this will be a new item)
	put1.Expected = make(ep.Expected)
	put1.Expected["t1.hk"] = ep.Constraints{Exists:false}
	body,code,err = put1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("put1 failed %d %v %s\n",code,err,body)
		os.Exit(1)
	}

	// GET THAT ITEM
        var get1 get_item.Request
        get1.TableName = tablename1
	get1.Key = make(ep.Item)
        get1.Key["TheHashKey"] = av1
        get1.Key["TheRangeKey"] = av2

	get_json,get_json_err := json.Marshal(get1)
	if (get_json_err != nil) {
		fmt.Printf("%v\n",get_json_err)
		os.Exit(1)
	}
	fmt.Printf("%s\n",string(get_json))

        body,code,err = get1.EndpointReq()
        if err != nil || code != http.StatusOK {
               fmt.Printf("get failed %d %v %s\n",code,err,body)
               os.Exit(1)
        }
	var gr get_item.Response
	um_err := json.Unmarshal([]byte(body),&gr)
	if um_err != nil {
		fmt.Printf("get resp unmarshal failed %s\n",um_err.Error())
		os.Exit(1)
	}

	// USE PUT TO REPLACE THAT ITEM CONDITIONALLY
        // **SET IT UP TO FAIL**
        var put2 put_item.Request
        put2.TableName = tablename1
        var av1_2,av2_2,av3_2,av4_2,av5_2,av6_2 ep.AttributeValue
        av1_2.S = hk
        av2_2.N = rk
        av3_2.SS = []string{"pk1_d","pk1_e","pk1_f"}
        av4_2.NS = []string{"4","5","6"}
        av5_2.N  = "2"
        av6_2.S  = "hello there"
        put2.Item = make(ep.Item)
	put2.Item["TheHashKey"] = av1_2
	put2.Item["TheRangeKey"] = av2_2
        put2.Item["stringlist"] = av3_2
        put2.Item["numlist"] = av4_2
        put2.Item["num"] = av5_2
        put2.Item["string"] = av6_2
        put2.Item["byte"] = av6
        put2.Item["bytelist"] = av7
        put2.Expected = make(ep.Expected)
        put2.Expected["TheHashKey"] = ep.Constraints{Exists:false}
        body,code,err = put2.EndpointReq()
        if err == nil && code == http.StatusOK {
                fmt.Printf("put2 should have failed %d %v %s\n",code,err,body)
                os.Exit(1)
        }

	// NOW MAKE IT SO THAT PUT WORKS
        put2.Expected = make(ep.Expected)
        put2.ReturnValues = put_item.RETVAL_ALL_OLD
        put2.Expected["num"] = ep.Constraints{Exists:true,Value:av5}
        body,code,err = put2.EndpointReq()
        if err != nil || code != http.StatusOK {
                fmt.Printf("put2 item failed %d %v %s\n",code,err,body)
                os.Exit(1)
        }

	// UPDATE THAT ITEM
        var up1 update_item.Request
        new_attr_val := "new string here"
        up1.TableName = tablename1
	up1.Key = make(ep.Item)
        up1.Key["TheHashKey"] = av1
        up1.Key["TheRangeKey"] = av2

        up1.AttributeUpdates = make(update_item.AttributeUpdates)
        up1.AttributeUpdates["new_string"] =
                update_item.AttributeAction{Value:ep.AttributeValue{S:new_attr_val},Action:update_item.ACTION_PUT}
        var del_stringlist ep.AttributeValue
        del_stringlist.SS = []string{"pk1_a"}
        up1.AttributeUpdates["stringlist"] =
                update_item.AttributeAction{Value:del_stringlist,Action:update_item.ACTION_DEL}
        up1.AttributeUpdates["byte"] = update_item.AttributeAction{Value:ep.AttributeValue{},Action:update_item.ACTION_DEL}
        up1.AttributeUpdates["num"] = update_item.AttributeAction{Value:ep.AttributeValue{N:"4"},Action:update_item.ACTION_ADD}
        up1.ReturnValues = update_item.RETVAL_ALL_NEW

	update_item_json,update_item_json_err := json.Marshal(up1)
	if (update_item_json_err != nil) {
		fmt.Printf("%v\n",update_item_json_err)
		os.Exit(1)
	}
	fmt.Printf("%s\n",string(update_item_json))

        body,code,err = up1.EndpointReq()
        if err != nil || code != http.StatusOK {
               fmt.Printf("update item failed %d %v %s\n",code,err,body)
               os.Exit(1)
        }

	var ur update_item.Response
	um_err = json.Unmarshal([]byte(body),&ur)
	if um_err != nil {
		fmt.Printf("update resp unmarshal failed %s\n",um_err.Error())
		os.Exit(1)
	}

	// GET IT AGAIN
        body,code,err = get1.EndpointReq()
        if err != nil || code != http.StatusOK {
               fmt.Printf("get failed %d %v %s\n",code,err,body)
               os.Exit(1)
        }

	// DELETE THE ITEM
	var del1 delete_item.Request
	del1.Key = make(ep.Item)
        del1.TableName = tablename1
        del1.Key["TheHashKey"] = av1
        del1.Key["TheRangeKey"] = av2

        del1.Expected = make(ep.Expected)
        del1.Expected["num"] =  ep.Constraints{Exists:true,Value:ep.AttributeValue{N:"6"}}
        del1.ReturnValues = delete_item.RETVAL_ALL_OLD
        body,code,err = del1.EndpointReq()
        if err != nil || code != http.StatusOK {
                fmt.Printf("delete item failed %d %v %s\n",code,err,body)
                os.Exit(1)
        }
}
