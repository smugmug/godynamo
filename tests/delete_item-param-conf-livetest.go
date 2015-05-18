package main

import (
	"fmt"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	delete_item "github.com/smugmug/godynamo/endpoints/delete_item"
	"github.com/smugmug/godynamo/types/attributevalue"
	"net/http"
	"os"
)

func main() {

	home := os.Getenv("HOME")
	home_conf_file := home + string(os.PathSeparator) + "." + conf.CONF_NAME
	home_conf, home_conf_err := conf_file.ReadConfFile(home_conf_file)
	if home_conf_err != nil {
		panic("cannot read conf from " + home_conf_file)
	}
	home_conf.ConfLock.RLock()
	if home_conf.Initialized == false {
		panic("conf struct has not been initialized")
	}

	tn := "test-godynamo-livetest"
	tablename1 := tn
	fmt.Printf("tablename1: %s\n", tablename1)

	// DELETE AN ITEM
	del_item1 := delete_item.NewDeleteItem()
	del_item1.TableName = tablename1
	del_item1.Key["TheHashKey"] = &attributevalue.AttributeValue{S: "AHashKey1"}
	del_item1.Key["TheRangeKey"] = &attributevalue.AttributeValue{N: "1"}

	body, code, err := del_item1.EndpointReqWithConf(home_conf)
	if err != nil || code != http.StatusOK {
		fmt.Printf("fail delete %d %v %s\n", code, err, string(body))
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n,%v\n", string(body), code, err)
}
