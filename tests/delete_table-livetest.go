package main

import (
	"fmt"
	"os"
	"net/http"
	delete_table "github.com/smugmug/godynamo/endpoints/delete_table"
	list "github.com/smugmug/godynamo/endpoints/list_tables"
	conf_iam "github.com/smugmug/godynamo/conf_iam"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	"log"
	keepalive "github.com/smugmug/godynamo/keepalive"
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
		_ = <- iam_ready_chan
	}
	conf.Vals.ConfLock.RUnlock()

	tn := "test-godynamo-livetest"
	tablename1 := tn
	fmt.Printf("tablename1: %s\n",tablename1)

	var code int
	var err error
	var body string

	// DELETE THE TABLE
	var del_table1 delete_table.Request
        del_table1.TableName = tablename1
        _,code,err = del_table1.EndpointReq()
        if err != nil || code != http.StatusOK {
               fmt.Printf("fail delete %d %v %s\n",code,err,body)
               os.Exit(1)
        }

	// List TABLES
	var l list.List
	l.ExclusiveStartTableName = ""
	l.Limit = 100
	lbody,lcode,lerr := l.EndpointReq()
	fmt.Printf("%v\n%v\n,%v\n",lbody,lcode,lerr)
}
