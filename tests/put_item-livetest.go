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
	"time"
	"net/http"
	ep "github.com/smugmug/godynamo/endpoint"
	put "github.com/smugmug/godynamo/endpoints/put_item"
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

	var put1 put.Request
	put1.TableName = "test-godynamo-livetest"
	k := fmt.Sprintf("hk1")
	v := fmt.Sprintf("%v",time.Now().Unix())
	put1.Item = make(ep.Item)
	put1.Item["TheHashKey"] = ep.AttributeValue{S:k}
	put1.Item["TheRangeKey"] = ep.AttributeValue{N:v}
	n := fmt.Sprintf("%v",time.Now().Unix())
	put1.Item["Mtime"] = ep.AttributeValue{N:n}
	put1.Item["SomeJunk"] = ep.AttributeValue{S:"some junk"}
	put1.Item["SomeJunks"] = ep.AttributeValue{SS:[]string{"some junk1","some junk2"}}
	body,code,err := put1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("put failed %d %v %s\n",code,err,body)
	}
	fmt.Printf("%s\n",string(body))
}
