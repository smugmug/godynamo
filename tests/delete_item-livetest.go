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
	"net/http"
	ep "github.com/smugmug/godynamo/endpoint"
	delete_item "github.com/smugmug/godynamo/endpoints/delete_item"
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
	tablename1 := tn
	fmt.Printf("tablename1: %s\n",tablename1)

	var code int
	var err error
	var body string

	// DELETE AN ITEM
	var del_item1 delete_item.Request
        del_item1.TableName = tablename1
	del_item1.Key = make(ep.Item)
	del_item1.Key["TheHashKey"] = ep.AttributeValue{S:"AHashKey1"}
	del_item1.Key["TheRangeKey"] = ep.AttributeValue{N:"1"}

        _,code,err = del_item1.EndpointReq()
        if err != nil || code != http.StatusOK {
               fmt.Printf("fail delete %d %v %s\n",code,err,body)
               os.Exit(1)
        }
}
