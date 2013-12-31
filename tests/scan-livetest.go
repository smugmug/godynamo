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
	"encoding/json"
	ep "github.com/smugmug/godynamo/endpoint"
	scan "github.com/smugmug/godynamo/endpoints/scan"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
)

func main() {
	conf_file.Read()
	if conf.Vals.Initialized == false {
		panic("the conf.Vals global conf struct has not been initialized")
	}

	iam_ready_chan := make(chan bool)
	go conf_iam.GoIAM(iam_ready_chan)
	iam_ready := <- iam_ready_chan
	if iam_ready {
		fmt.Printf("using iam\n")
	} else {
		fmt.Printf("not using iam\n")
	}

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
