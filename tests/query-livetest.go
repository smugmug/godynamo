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
	query "github.com/smugmug/godynamo/endpoints/query"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
)

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
	q := query.NewQuery()
	q.TableName = tn
	q.Select = ep.SELECT_ALL
	k_v1 := fmt.Sprintf("AHashKey%d",100)
	var kc query.KeyCondition
	kc.AttributeValueList = make([]ep.AttributeValue,1)
	kc.AttributeValueList[0] = ep.AttributeValue{S:k_v1}
	kc.ComparisonOperator = query.OP_EQ
	q.Limit = 10000
	q.KeyConditions["TheHashKey"] = kc
	json,_ := json.Marshal(q)
	fmt.Printf("JSON:%s\n",string(json))
	body,code,err := q.EndpointReq()
	fmt.Printf("%v\n%v\n%v\n",body,code,err)
}
