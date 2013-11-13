// Copyright (c) 2013, SmugMug, Inc. All rights reserved.
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
	create "github.com/smugmug/godynamo/endpoints/create_table"
	list "github.com/smugmug/godynamo/endpoints/list_tables"
	desc "github.com/smugmug/godynamo/endpoints/describe_table"
	ep "github.com/smugmug/godynamo/endpoint"
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

	tablename1 := "test-godynamo-livetest"
	fmt.Printf("tablename: %s\n",tablename1)

	var code int
	var err error
	var body string

	// CREATE TABLE
	var create1 create.Create
	create1.TableName = tablename1
	create1.ProvisionedThroughput.ReadCapacityUnits = 100
	create1.ProvisionedThroughput.WriteCapacityUnits = 100
	create1.AttributeDefinitions = append(create1.AttributeDefinitions,
		ep.AttributeDefinition{AttributeName:"TheHashKey",AttributeType:ep.S})
	create1.AttributeDefinitions = append(create1.AttributeDefinitions,
		ep.AttributeDefinition{AttributeName:"TheRangeKey",AttributeType:ep.N})
	create1.AttributeDefinitions = append(create1.AttributeDefinitions,
		ep.AttributeDefinition{AttributeName:"AnAttrName",AttributeType:ep.S})
	create1.KeySchema = append(create1.KeySchema,
		ep.KeyDefinition{AttributeName:"TheHashKey",KeyType:ep.HASH})
	create1.KeySchema = append(create1.KeySchema,
		ep.KeyDefinition{AttributeName:"TheRangeKey",KeyType:ep.RANGE})
	lsi := ep.NewLocalSecondaryIndex()
	lsi.IndexName = "AnAttrIndex"
	lsi.Projection.ProjectionType = ep.KEYS_ONLY
	lsi.KeySchema = append(lsi.KeySchema,
		ep.KeyDefinition{AttributeName:"TheHashKey",KeyType:ep.HASH})
	lsi.KeySchema = append(lsi.KeySchema,
		ep.KeyDefinition{AttributeName:"AnAttrName",KeyType:ep.RANGE})
	create1.LocalSecondaryIndexes = append(create1.LocalSecondaryIndexes,*lsi)

	create_json,create_json_err := json.Marshal(create1)
	if (create_json_err != nil) {
		fmt.Printf("%v\n",create_json_err)
		os.Exit(1)
	}
	fmt.Printf("%s\n",string(create_json))
	fmt.Printf("%v\n",create1)

	cbody,ccode,cerr := create1.EndpointReq()
	fmt.Printf("%v\n%v\n,%v\n",cbody,ccode,cerr)

	var desc1 desc.Describe
	desc1.TableName = tablename1
	body,code,err = desc1.EndpointReq()
	fmt.Printf("desc:%v\n%v\n,%v\n",body,code,err)

	// WAIT FOR IT TO BE ACTIVE
	fmt.Printf("checking for ACTIVE status for table....\n")
	active,poll_err := desc.PollTableStatus(tablename1,desc.ACTIVE,100)
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
