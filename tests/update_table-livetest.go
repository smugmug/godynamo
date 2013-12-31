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
	update_table "github.com/smugmug/godynamo/endpoints/update_table"
	list "github.com/smugmug/godynamo/endpoints/list_tables"
	desc "github.com/smugmug/godynamo/endpoints/describe_table"
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

        var update_table1 update_table.Request
        update_table1.TableName = tablename1
        update_table1.ProvisionedThroughput.ReadCapacityUnits = 200
        update_table1.ProvisionedThroughput.WriteCapacityUnits = 200
        body,code,err = update_table1.EndpointReq()
        if err != nil || code != http.StatusOK {
                fmt.Printf("update table failed %d %v %s\n",code,err,body)
                os.Exit(1)
        }

        // WAIT FOR THE PROVISIONING TO FINISH
        fmt.Printf("checking for ACTIVE status for update....\n")
        active,poll_err := desc.PollTableStatus(tablename1,
                desc.ACTIVE,100)
	if poll_err != nil {
		fmt.Printf("poll1:%v\n",poll_err)
		os.Exit(1)
	}
	fmt.Printf("ACTIVE:%v\n",active)

	var desc1 desc.Describe
	desc1.TableName = tablename1
	body,code,err = desc1.EndpointReq()
	fmt.Printf("desc:%v\n%v\n,%v\n",body,code,err)

	// WAIT FOR IT TO BE ACTIVE
	fmt.Printf("checking for ACTIVE status for table....\n")
	active,poll_err = desc.PollTableStatus(tablename1,desc.ACTIVE,100)
	if poll_err != nil {
		fmt.Printf("poll1:%v\n",poll_err)
		os.Exit(1)
	}
	fmt.Printf("ACTIVE:%v\n",active)

	// List TABLES
	var l list.List
	l.ExclusiveStartTableName = ""
	l.Limit = 100
	lbody,lcode,lerr := l.EndpointReq()
	fmt.Printf("%v\n%v\n,%v\n",lbody,lcode,lerr)

}
