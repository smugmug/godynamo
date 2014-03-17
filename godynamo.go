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

// A dummy package for installers who wish to get all GoDynamo depencies in one "go get" invocation.
package godynamo

import (
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	batch_get_item "github.com/smugmug/godynamo/endpoints/batch_get_item"
	batch_write_item "github.com/smugmug/godynamo/endpoints/batch_write_item"
	create "github.com/smugmug/godynamo/endpoints/create_table"
	delete_item "github.com/smugmug/godynamo/endpoints/delete_item"
	delete_table "github.com/smugmug/godynamo/endpoints/delete_table"
	describe_table "github.com/smugmug/godynamo/endpoints/describe_table"
	get_item "github.com/smugmug/godynamo/endpoints/get_item"
	list_tables "github.com/smugmug/godynamo/endpoints/list_tables"
	put_item "github.com/smugmug/godynamo/endpoints/put_item"
	query "github.com/smugmug/godynamo/endpoints/query"
	scan "github.com/smugmug/godynamo/endpoints/scan"
	update_item "github.com/smugmug/godynamo/endpoints/update_item"
	update_table "github.com/smugmug/godynamo/endpoints/update_table"
)

// This program serves only to include all of the libraries in GoDynamo so that you can
// just run one instance of 'go get' and get the totality of the library.

func installAll() {
	// conf file must be read in before anything else, to initialize permissions etc
	conf_file.Read()
	if conf.Vals.Initialized == false {
		panic("the conf.Vals global conf struct has not been initialized")
	}

	// deal with iam, or not
	iam_ready_chan := make(chan bool)
	go conf_iam.GoIAM(iam_ready_chan)
	_ = <-iam_ready_chan

	var get1 get_item.Request
	var put1 put_item.Request
	var up1 update_item.Request
	var upt1 update_table.Request
	var del1 delete_item.Request
	var delt1 delete_table.Request
	var batchw1 batch_write_item.Request
	var batchg1 batch_get_item.Request
	var create1 create.Request
	var query1 query.Request
	var scan1 scan.Request
	var desc1 describe_table.Request
	var list1 list_tables.Request
	fmt.Printf("%v%v%v%v%v%v%v%v%v%v%v%v%v", get1, put1, up1, upt1, del1, batchw1, batchg1, create1, delt1, query1, scan1, desc1, list1)

}
