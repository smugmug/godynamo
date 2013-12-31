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

// Manages writing roles credentials into the global conf.Vals state.
//
// example use:
//
//   iam_ready_chan := make(chan bool)
//   go conf_iam.GoIAM(iam_ready_chan)
//   iam_ready := <- iam_ready_chan
//   if iam_ready {
//	fmt.Printf("using iam\n")
//   } else {
//	fmt.Printf("not using iam\n")
//   }
package conf_iam

import (
	"fmt"
	"time"
	"errors"
	"log/syslog"
	"github.com/bradclawsie/slog"
	conf "github.com/smugmug/godynamo/conf"
	roles_files "github.com/smugmug/goawsroles/roles_files"
)

// AssignCredentials locks the global state and copy over roles data.
func AssignCredentials(rf *roles_files.RolesFiles) (error) {
	accessKey,secret,token,get_err := rf.Get()
	if get_err != nil {
		e := fmt.Sprintf("conf_iam.ReadIAM:cannot get a role file:%s",
			get_err.Error())
		return errors.New(e)
	}
	conf.Vals.ConfLock.Lock()
	conf.Vals.IAM.Credentials.AccessKey = accessKey
	conf.Vals.IAM.Credentials.Secret    = secret
	conf.Vals.IAM.Credentials.Token     = token
	conf.Vals.ConfLock.Unlock()
	e := fmt.Sprintf("IAM credentials assigned at %v",time.Now())
	slog.SLog(syslog.LOG_NOTICE,e,true)
	return nil
}

// ReadIAM explicitly mutates the global shared conf.Vals state by reading in the IAM. Use
// this function at program startup or any time you need to force a refresh of the IAM
// credentials.
func ReadIAM(rf *roles_files.RolesFiles) (error) {
	roles_read_err := rf.RolesRead()
	if roles_read_err != nil {
		e := fmt.Sprintf("conf_iam.ReadIAM:cannot perform initial roles read: %s",
			roles_read_err.Error())
		return errors.New(e)
	}
	return AssignCredentials(rf)
}

// WatchIAM will receive notifications for changes in IAM files and update credentials when a read signal is received.
func WatchIAM(rf *roles_files.RolesFiles,watch_err_chan chan error) {
	err_chan := make(chan error)
	read_signal := make(chan bool)
	go rf.RolesWatch(err_chan,read_signal)
	e := "IAM watching set to true, waiting..."
	slog.SLog(syslog.LOG_NOTICE,e,true)
	for {
		select {
		case roles_watch_err := <- err_chan:
			watch_err_chan <- roles_watch_err
		case <- read_signal:
			e := "WatchIAM received a read signal"
			slog.SLog(syslog.LOG_NOTICE,e,true)
			assign_err := AssignCredentials(rf)
			if assign_err != nil {
				watch_err_chan <- assign_err
			}
		}
	}
}

// GoIAM is a convenience wrapper for callers using roles_files instantiation of the roles interface.
// First there is a blocking read on the roles files to get the initial roles information. Then the
// file notification watcher will run as a goroutine, and resetting the global conf.Vals roles
// values. If IAM Credentials are ready for use, the parameter chan `ready_chan` will receive a true
// value, otherwise false. A false value on this chan should indicate to a caller that another auth
// mechanism (for example, hardocded credentials) should be used.
func GoIAM(ready_chan chan bool) {
	use_iam := false
	conf.Vals.ConfLock.RLock()
	use_iam = conf.Vals.UseIAM
	conf.Vals.ConfLock.RUnlock()
	if use_iam == true {
		rf := roles_files.NewRolesFiles()
		watching := false
		conf.Vals.ConfLock.RLock()
		rf.BaseDir       = conf.Vals.IAM.File.BaseDir
		rf.AccessKeyFile = conf.Vals.IAM.File.AccessKey
		rf.SecretFile    = conf.Vals.IAM.File.Secret
		rf.TokenFile     = conf.Vals.IAM.File.Token
		watching         = conf.Vals.IAM.Watch
		conf.Vals.ConfLock.RUnlock()
		roles_read_err := ReadIAM(rf)
		if roles_read_err != nil {
			e := fmt.Sprintf("conf_iam.GoIAM:cannot perform initial roles read: %s",
				roles_read_err.Error())
			slog.SLog(syslog.LOG_ERR,e,true)
			conf.Vals.UseIAM = false
			ready_chan <- false
		}
		// signal to caller that iam roles are ready to use
		ready_chan <- true
		if watching == true {
			watch_err := make(chan error)
			go WatchIAM(rf,watch_err)
			go func() {
				select {
				case err := <- watch_err:
					if err != nil {
						slog.SLog(syslog.LOG_ERR,err.Error(),true)
						// caller can fall back to hard-coded perms
						// or live with the panic
						conf.Vals.UseIAM = false
					}
				}
			}()
		}
	} else {
		// signal to the caller than iam roles are selected as a auth mechanism
		ready_chan <- false
	}
}
