// Manages reading the conf file into the global var as described in the `conf` package.
package conf_file

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

// ReadConfFile will attempt to read in the conf file path passed as a parameter
// and convert it into a conf.AWS_Conf struct pointer. You can use this to read in
// a conf for a file of your own choosing.
func ReadConfFile(conf_file string) (*conf.AWS_Conf, error) {
	conf_bytes, conf_err := ioutil.ReadFile(conf_file)
	if conf_err != nil {
		e := fmt.Sprintf("conf_file.ReadConfFile: cannot read conf file %s", conf_file)
		return nil, errors.New(e)
	}

	var cf conf.SDK_conf_file
	um_err := json.Unmarshal(conf_bytes, &cf)
	if um_err != nil {
		e := fmt.Sprintf("conf_file.ReadConfFile: cannot unmarshal %s. json err %s",
			conf_file, um_err.Error())
		return nil, errors.New(e)
	}

	var c conf.AWS_Conf
	// make sure the dynamo endpoint is available
	addrs, addrs_err := net.LookupIP(cf.Services.Dynamo_db.Host)
	if addrs_err != nil {
		e := fmt.Sprintf("conf_file.ReadConfFile: cannot lookup hostname %s",
			cf.Services.Dynamo_db.Host)
		return nil, errors.New(e)
	}
	dynamo_ip := (addrs[0]).String()

	// assign the values to our globally-available c struct instance
	c.Auth.AccessKey = cf.Services.Default_settings.Params.Access_key_id
	c.Auth.Secret = cf.Services.Default_settings.Params.Secret_access_key
	c.UseSysLog = cf.Services.Default_settings.Params.Use_sys_log
	c.Network.DynamoDB.Host = cf.Services.Dynamo_db.Host
	c.Network.DynamoDB.IP = dynamo_ip
	c.Network.DynamoDB.Zone = cf.Services.Dynamo_db.Zone
	scheme := "http"
	port := aws_const.PORT // already a string
	if cf.Services.Dynamo_db.Scheme != "" {
		scheme = cf.Services.Dynamo_db.Scheme
	}
	if cf.Services.Dynamo_db.Port != 0 {
		port = strconv.Itoa(cf.Services.Dynamo_db.Port)
	}
	c.Network.DynamoDB.Port = port
	c.Network.DynamoDB.Scheme = scheme
	c.Network.DynamoDB.URL = scheme + "://" + c.Network.DynamoDB.Host +
		":" + port
	_, url_err := url.Parse(c.Network.DynamoDB.URL)
	if url_err != nil {
		return nil, errors.New("conf_file.ReadConfFile: conf.Network.DynamoDB.URL malformed")
	}

	// If set to true, programs that are written with godynamo may
	// opt to launch the keepalive goroutine to keep conns open.
	c.Network.DynamoDB.KeepAlive = cf.Services.Dynamo_db.KeepAlive

	// read in flags for IAM support
	if cf.Services.Dynamo_db.IAM.Use_iam == true {
		// caller will have to check the RoleProvider to dispatch further Roles features
		c.IAM.RoleProvider = cf.Services.Dynamo_db.IAM.Role_provider
		c.IAM.File.BaseDir = cf.Services.Dynamo_db.IAM.Base_dir
		c.IAM.File.AccessKey = cf.Services.Dynamo_db.IAM.Access_key
		c.IAM.File.Secret = cf.Services.Dynamo_db.IAM.Secret_key
		c.IAM.File.Token = cf.Services.Dynamo_db.IAM.Token
		if cf.Services.Dynamo_db.IAM.Watch == true {
			c.IAM.Watch = true
		} else {
			c.IAM.Watch = false
		}
		c.UseIAM = true
	}
	c.Initialized = true
	return &c, nil
}

// ReadDefaultConfs will check the preset standard locations for conf files and
// attempt to create a conf.AWS_Conf struct pointer with the first one found.
// The order of precedence is:
// 1. file in $GODYNAMO_CONF_FILE if you wish to set it.
// 2. $HOME/.aws-config.json
// 3. /etc/aws-config.json (note lack of prefix '.')
// This function can be useful when moving back and forth between environments
// where you wish to have a conf file in $HOME that precludes one in /etc.
func ReadDefaultConfs() (*conf.AWS_Conf, error) {
	local_conf := os.Getenv("HOME") + string(filepath.Separator) + "." + conf.CONF_NAME
	etc_conf := string(filepath.Separator) + "etc" + string(filepath.Separator) + conf.CONF_NAME
	conf_files := make([]string, 0)
	// assumes that if set, this is a fully-qualified file path
	const env_conf = "GODYNAMO_CONF_FILE"
	// assumes that if set, this is a fully-qualified file path
	if os.Getenv(env_conf) != "" {
		conf_files = append(conf_files, os.Getenv(env_conf))
	}
	conf_files = append(conf_files, local_conf)
	conf_files = append(conf_files, etc_conf)

CONF_LOCATIONS:
	for _, conf_file := range conf_files {
		c, c_err := ReadConfFile(conf_file)
		if c_err != nil {
			e := fmt.Sprintf("conf_file.ReadDefaultConfs: problem with conf %s: %s",
				conf_file, c_err.Error())
			log.Printf(e)
			continue CONF_LOCATIONS
		} else {
			return c, nil
		}
	}
	e := fmt.Sprintf("conf_File.ReadDefaultConfs: no conf struct found in %v", conf_files)
	return nil, errors.New(e)
}

// ReadGlobal will take a configuration and assign it to the global conf.Vals struct instance.
func ReadGlobal() {
	local_conf, conf_err := ReadDefaultConfs()
	if conf_err != nil {
		e := fmt.Sprintf("cannot locate a valid global configuration: %s", conf_err.Error())
		panic(e)
	}
	conf.Vals.Copy(local_conf)
}

// Read assigns a configuration to the global conf.Vals struct instance.
// This function exists as backwards compatibility only. Do not use it in new code,
// instead use ReadLocal.
func Read() {
	ReadGlobal()
}
