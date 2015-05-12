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
	"errors"
	"fmt"
	conf "github.com/smugmug/godynamo/conf"
	roles_files "github.com/smugmug/goawsroles/roles_files"
	"log"
	"time"
)

// AssignCredentialsToConf will safely copy the credentials data from rf to the conf c.
func AssignCredentialsToConf(rf *roles_files.RolesFiles, c *conf.AWS_Conf) error {
	if rf == nil || c == nil {
		return errors.New("conf_iam.AssignCredentialsToConf: rf or c is nil")
	}
	accessKey, secret, token, get_err := rf.Get()
	if get_err != nil {
		e := fmt.Sprintf("conf_iam.AssignCredentialsToConf:cannot get a role file:%s",
			get_err.Error())
		return errors.New(e)
	}
	c.ConfLock.Lock()
	c.IAM.Credentials.AccessKey = accessKey
	c.IAM.Credentials.Secret = secret
	c.IAM.Credentials.Token = token
	c.ConfLock.Unlock()
	e := fmt.Sprintf("IAM credentials assigned at %v", time.Now())
	log.Printf(e)
	return nil
}

// AssignCredentials will safely copy the credentials data from rf to the global conf.Vals.
func AssignCredentials(rf *roles_files.RolesFiles) error {
	return AssignCredentialsToConf(rf, &conf.Vals)
}

// ReadIAMToConf will read the credentials data from rf and safely assign it to conf c.
func ReadIAMToConf(rf *roles_files.RolesFiles, c *conf.AWS_Conf) error {
	if rf == nil || c == nil {
		return errors.New("conf_iam.ReadIAMToConf: rf or c is nil")
	}
	roles_read_err := rf.RolesRead()
	if roles_read_err != nil {
		e := fmt.Sprintf("conf_iam.ReadIAM:cannot perform initial roles read: %s",
			roles_read_err.Error())
		return errors.New(e)
	}
	return AssignCredentialsToConf(rf, c)
}

// ReadIAM will read the credentials data from rf and safely assign it to the global conf.Vals.
func ReadIAM(rf *roles_files.RolesFiles) error {
	return ReadIAMToConf(rf, &conf.Vals)
}

// WatchIAMToConf will begin rf's RolesWatch method and wait to receive signals that new credentials
// are available to be assigned, which will then be safely copied to conf c.
func WatchIAMToConf(rf *roles_files.RolesFiles, c *conf.AWS_Conf, watch_err_chan chan error) {
	if rf == nil || c == nil {
		watch_err_chan <- errors.New("conf_iam.WatchIAMToConf: rf or c is nil")
		return
	}
	err_chan := make(chan error)
	read_signal := make(chan bool)
	go rf.RolesWatch(err_chan, read_signal)
	e := "IAM watching set to true, waiting..."
	log.Printf(e)
	for {
		select {
		case roles_watch_err := <-err_chan:
			watch_err_chan <- roles_watch_err
		case <-read_signal:
			e := "WatchIAM received a read signal"
			log.Printf(e)
			assign_err := AssignCredentialsToConf(rf, c)
			if assign_err != nil {
				watch_err_chan <- assign_err
			}
		}
	}
}

// WatchIAMToConf will begin rf's RolesWatch method and wait to receive signals that new credentials
// are available to be assigned, which will then be safely copied to the global conf.Vals.
func WatchIAM(rf *roles_files.RolesFiles, watch_err_chan chan error) {
	WatchIAMToConf(rf, &conf.Vals, watch_err_chan)
}

// GoIAMToConf is a convenience wrapper for callers using roles_files instantiation of the roles interface.
// First there is a blocking read on the roles files to get the initial roles information. Then the
// file notification watcher will run as a goroutine, and resetting conf c's roles
// values. If IAM Credentials are ready for use, the parameter chan `ready_chan` will receive a true
// value, otherwise false. A false value on this chan should indicate to a caller that another auth
// mechanism (for example, hardocded credentials) should be used.
func GoIAMToConf(c *conf.AWS_Conf, ready_chan chan bool) {
	if c == nil {
		log.Printf("conf_iam.GoIAMToConf: c is nil")
		ready_chan <- false
		return
	}
	use_iam := false
	c.ConfLock.RLock()
	use_iam = c.UseIAM
	c.ConfLock.RUnlock()
	if use_iam == true {
		rf := roles_files.NewRolesFiles()
		watching := false
		c.ConfLock.RLock()
		rf.BaseDir = c.IAM.File.BaseDir
		rf.AccessKeyFile = c.IAM.File.AccessKey
		rf.SecretFile = c.IAM.File.Secret
		rf.TokenFile = c.IAM.File.Token
		watching = c.IAM.Watch
		c.ConfLock.RUnlock()
		roles_read_err := ReadIAMToConf(rf, c)
		if roles_read_err != nil {
			e := fmt.Sprintf("conf_iam.GoIAMToConf:cannot perform initial roles read: %s",
				roles_read_err.Error())
			log.Printf(e)
			c.ConfLock.Lock()
			c.UseIAM = false
			c.ConfLock.Unlock()
			ready_chan <- false
			return
		}
		// signal to caller that iam roles are ready to use
		ready_chan <- true
		if watching == true {
			watch_err := make(chan error)
			go WatchIAMToConf(rf, c, watch_err)
			go func() {
				select {
				case err := <-watch_err:
					if err != nil {
						log.Printf(err.Error())
						// caller can fall back to hard-coded perms
						// or live with the panic
						c.ConfLock.Lock()
						c.UseIAM = false
						c.ConfLock.Unlock()
					}
				}
			}()
		}
	} else {
		// signal to the caller than iam roles are not selected as a auth mechanism
		e := fmt.Sprintf("conf_iam.GoIAMToConf: not using IAM")
		log.Printf(e)
		ready_chan <- false
	}
}

// GoIAM calls the GoIAMToConf credential manager, assigning to the global conf.Vals.
func GoIAM(ready_chan chan bool) {
	GoIAMToConf(&conf.Vals, ready_chan)
}
