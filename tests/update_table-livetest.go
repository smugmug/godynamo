package main

import (
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	desc "github.com/smugmug/godynamo/endpoints/describe_table"
	list "github.com/smugmug/godynamo/endpoints/list_tables"
	update_table "github.com/smugmug/godynamo/endpoints/update_table"
	keepalive "github.com/smugmug/godynamo/keepalive"
	"log"
	"net/http"
	"os"
)

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

	update_table1 := update_table.NewUpdateTable()
	update_table1.TableName = tablename1
	update_table1.ProvisionedThroughput.ReadCapacityUnits = 200
	update_table1.ProvisionedThroughput.WriteCapacityUnits = 200
	body, code, err = update_table1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("update table failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}

	// WAIT FOR THE PROVISIONING TO FINISH
	fmt.Printf("checking for ACTIVE status for update....\n")
	active, poll_err := desc.PollTableStatus(tablename1,
		desc.ACTIVE, 100)
	if poll_err != nil {
		fmt.Printf("poll1:%v\n", poll_err)
		os.Exit(1)
	}
	fmt.Printf("ACTIVE:%v\n", active)

	var desc1 desc.DescribeTable
	desc1.TableName = tablename1
	body, code, err = desc1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("desc failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("desc:%v\n%v\n,%v\n", string(body), code, err)

	// WAIT FOR IT TO BE ACTIVE
	fmt.Printf("checking for ACTIVE status for table....\n")
	active, poll_err = desc.PollTableStatus(tablename1, desc.ACTIVE, 100)
	if poll_err != nil {
		fmt.Printf("poll1:%v\n", poll_err)
		os.Exit(1)
	}
	fmt.Printf("ACTIVE:%v\n", active)

	// List TABLES
	var l list.List
	l.ExclusiveStartTableName = ""
	l.Limit = 100
	body, code, err = l.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("list failed %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)

}
