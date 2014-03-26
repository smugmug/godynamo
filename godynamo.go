// A dummy package for installers who wish to get all GoDynamo depencies in one "go get" invocation.
package godynamo

import (
	"fmt"
	"log"
	put_item "github.com/smugmug/godynamo/endpoints/put_item"
	get_item "github.com/smugmug/godynamo/endpoints/get_item"
	update_item "github.com/smugmug/godynamo/endpoints/update_item"
	update_table "github.com/smugmug/godynamo/endpoints/update_table"
	delete_item "github.com/smugmug/godynamo/endpoints/delete_item"
	delete_table "github.com/smugmug/godynamo/endpoints/delete_table"
	describe_table "github.com/smugmug/godynamo/endpoints/describe_table"
	list_tables "github.com/smugmug/godynamo/endpoints/list_tables"
	batch_write_item "github.com/smugmug/godynamo/endpoints/batch_write_item"
	batch_get_item "github.com/smugmug/godynamo/endpoints/batch_get_item"
	create "github.com/smugmug/godynamo/endpoints/create_table"
	query "github.com/smugmug/godynamo/endpoints/query"
	scan "github.com/smugmug/godynamo/endpoints/scan"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	keepalive "github.com/smugmug/godynamo/keepalive"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
)

// This program serves only to include all of the libraries in GoDynamo so that you can
// just run one instance of 'go get' and get the totality of the library.

func installAll() {
	// conf file must be read in before anything else, to initialize permissions etc
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
		_ = <- iam_ready_chan
	}
	conf.Vals.ConfLock.RUnlock()

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
	fmt.Printf("%v%v%v%v%v%v%v%v%v%v%v%v%v",get1,put1,up1,upt1,del1,batchw1,batchg1,create1,delt1,query1,scan1,desc1,list1)
}
