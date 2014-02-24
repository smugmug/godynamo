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

// Manages reading the configuration file format into our shared conf state.
// The AWS SDKs utilize a conf file format that this package attempts compatibility with while
// supporting extra fields. This is why the type detailing the file and internal formats differ.
// See SAMPLE-aws-config.json in the source repository for a sample.
package conf

import (
	"sync"
)

const (
	CONF_NAME          = "aws-config.json"
	ROLE_PROVIDER_FILE = "file"
)

// SDK_conf_File roughly matches the format as used by recent amazon SDKs, plus some additions.
type SDK_conf_file struct {
	Extends []string
	Services struct {
		Default_settings struct {
			Params struct {
				// Traditional AWS access/secret authentication pair.
				Access_key_id string
				Secret_access_key string
				// If you use syslogd (a linux or *bsd system), you may set this to "true".
				// (currently unused)
				Use_sys_log bool
			}
		}
		Dynamo_db struct {
			// Your dynamo hostname.
			Host string
			// Your aws zone.
			Zone string
			IAM struct {
				// Set to true to use IAM authentication.
				Use_iam bool
				// The role provider is described in the goawsroles package.
				// See: https://github.com/smugmug/goawsroles/
				// Currently the only support is for the "file" provider, whereby
				// roles data is written to local files.
				Role_provider string
				// The identifier (filename, etc) for the IAM Access Key
				Access_key string
				// The identifier (filename, etc) for the IAM Secret Key
				Secret_key string
				// The identifier (filename, etc) for the IAM Token
				Token string
				// If using the "file" role provider, the base dir to read IAM files.
				Base_dir string
				// Set to true if you would like the roles resource watched for changes
				// and automatically (and atomically) updated.
				Watch bool
			}
		}
	}
}

// AWS_Conf is the structure used internally in godynamo.
type AWS_Conf struct {
	// Set to true if this struct is populated correctly.
	Initialized bool
	// Traditional AWS authentication pair.
	Auth struct {
		AccessKey string
		Secret string
	}
	// Dynamo connection data.
	Network struct {
		DynamoDB struct {
			Host string
			IP   string
			Zone string
			URL  string
		}
	}
	// If using syslogd
	UseSysLog bool
	// If using IAM
	UseIAM bool
	// The IAM role provider info
	IAM struct {
		RoleProvider string
		Watch bool
		// Tells you where the credentials can be read from
		File struct {
			AccessKey string
			Secret string
			Token string
			BaseDir string
		}
		// The credentials themselves, once loaded from Files.* above
		// these are kept distinct from the global AccessKey and Secret
		// in the event a caller wants a mixed model
		Credentials struct {
			AccessKey string
			Secret string
			Token string
		}
	}
	// Lock used when accessing IAM values, which will change during execution.
	// other values will persist for program duration so they can be read without locking.
	ConfLock sync.RWMutex
}

// Vals is the global conf vals struct. It is shared throughout the duration of program execution.
// Use the embedded ConfLock mutex to use it safely.
var (
	Vals AWS_Conf
)
