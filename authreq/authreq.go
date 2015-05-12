// Implements the wrapper for versioned retryable DynamoDB requests.
// See the init() function below for details about initial conf file processing.
package authreq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smugmug/godynamo/auth_v4"
	"github.com/smugmug/godynamo/aws_const"
	"github.com/smugmug/godynamo/conf"
	ep "github.com/smugmug/godynamo/endpoint"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	// auth version numbers
	AUTH_V2 = 2
	AUTH_V4 = 4
)

// Stipulate the current authorization version.
const AUTH_VERSION = AUTH_V4

var (
	exceeded_msg_bytes, unrecognized_client_msg_bytes, throttling_msg_bytes []byte
)

func init() {
	if AUTH_VERSION != AUTH_V4 {
		panic("authreq: only v4 authentication is enabled")
	}
	// convert these to []byte so we can search within responses
	exceeded_msg_bytes = []byte(aws_const.EXCEEDED_MSG)
	unrecognized_client_msg_bytes = []byte(aws_const.UNRECOGNIZED_CLIENT_MSG)
	throttling_msg_bytes = []byte(aws_const.THROTTLING_MSG)
}

// RetryReq_V4 sends a retry-able request using an ep.Endpoint structure and v4 auth.
// Uses the global conf.
func RetryReq_V4(v ep.Endpoint, amzTarget string) ([]byte, int, error) {
	if !conf.IsValid(&conf.Vals) {
		return nil, 0, errors.New("authreq.RetryReq_V4: conf not valid")
	}
	reqJSON, json_err := json.Marshal(v)
	if json_err != nil {
		return nil, 0, json_err
	}
	return retryReq(reqJSON, amzTarget, &conf.Vals)
}

// RetryReq_V4 sends a retry-able request using a JSON serialized request and v4 auth.
// Uses the global conf.
func RetryReqJSON_V4(reqJSON []byte, amzTarget string) ([]byte, int, error) {
	if !conf.IsValid(&conf.Vals) {
		return nil, 0, errors.New("authreq.RetryReqJSON_V4: conf not valid")
	}
	return retryReq(reqJSON, amzTarget, &conf.Vals)
}

// RetryReq_V4 sends a retry-able request using an ep.Endpoint structure and v4 auth.
// Uses a parameterized conf.
func RetryReq_V4WithConf(v ep.Endpoint, amzTarget string, c *conf.AWS_Conf) ([]byte, int, error) {
	if !conf.IsValid(c) {
		return nil, 0, errors.New("authreq.RetryReqV4WithConf: conf not valid")
	}
	reqJSON, json_err := json.Marshal(v)
	if json_err != nil {
		return nil, 0, json_err
	}
	return retryReq(reqJSON, amzTarget, c)
}

// RetryReq_V4 sends a retry-able request using a JSON serialized request and v4 auth.
// Uses a parameterized conf.
func RetryReqJSON_V4WithConf(reqJSON []byte, amzTarget string, c *conf.AWS_Conf) ([]byte, int, error) {
	if !conf.IsValid(c) {
		return nil, 0, errors.New("authreq.RetryReqJSON_V4WithConf: conf not valid")
	}
	return retryReq(reqJSON, amzTarget, c)
}

// Implement exponential backoff for the req above in the case of 5xx errors
// from aws. Algorithm is lifted from AWS docs.
// returns []byte respBody, int httpcode, error
func retryReq(reqJSON []byte, amzTarget string, c *conf.AWS_Conf) ([]byte, int, error) {
	// conf.IsValid has already been established by caller
	resp_body, amz_requestid, code, resp_err := auth_v4.ReqWithConf(reqJSON, amzTarget, c)
	shouldRetry := false
	if resp_err != nil {
		e := fmt.Sprintf("authreq.retryReq:0 "+
			" try AuthReq Fail:%s (reqid:%s)", resp_err.Error(), amz_requestid)
		log.Printf("authreq.retryReq: call err %s\n", e)
		shouldRetry = true
	}
	// see:
	// http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ErrorHandling.html
	if code >= http.StatusInternalServerError {
		shouldRetry = true // all 5xx codes are deemed retryable by amazon
	}
	if code == http.StatusBadRequest {
		if bytes.Contains(resp_body, exceeded_msg_bytes) {
			log.Printf("authreq.retryReq THROUGHPUT WARNING RETRY\n")
			shouldRetry = true
		} else if bytes.Contains(resp_body, unrecognized_client_msg_bytes) {
			log.Printf("authreq.retryReq CLIENT WARNING RETRY\n")
			shouldRetry = true
		} else if bytes.Contains(resp_body, throttling_msg_bytes) {
			log.Printf("authreq.retryReq THROUGHPUT WARNING RETRY\n")
			shouldRetry = true
		} else {
			log.Printf("authreq.retryReq un-retryable err: %s\n%s (reqid:%s)\n",
				string(resp_body), string(reqJSON), amz_requestid)
			shouldRetry = false
		}
	}
	if !shouldRetry {
		// not retryable
		return resp_body, code, resp_err
	} else {
		// retry the request RETRIES time in the case of a 5xx
		// response, with an exponentially decayed sleep interval

		// seed our rand number generator g
		g := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 1; i < aws_const.RETRIES; i++ {
			// get random delay from range
			// [0..4**i*100 ms)
			log.Printf("authreq.retryReq: BEGIN SLEEP %v (code:%v) (REQ:%s) (reqid:%s)",
				time.Now(), code, string(reqJSON), amz_requestid)
			r := time.Millisecond *
				time.Duration(g.Int63n(int64(
					math.Pow(4, float64(i)))*
					100))
			time.Sleep(r)
			log.Printf("authreq.retryReq END SLEEP %v\n", time.Now())
			shouldRetry = false
			resp_body, amz_requestid, code, resp_err := auth_v4.ReqWithConf(reqJSON, amzTarget, c)
			if resp_err != nil {
				_ = fmt.Sprintf("authreq.retryReq:1 "+
					" try AuthReq Fail:%s (reqid:%s)", resp_err.Error(), amz_requestid)
				shouldRetry = true
			}
			if code >= http.StatusInternalServerError {
				shouldRetry = true
			}
			if code == http.StatusBadRequest {
				if bytes.Contains(resp_body, exceeded_msg_bytes) {
					log.Printf("authreq.retryReq THROUGHPUT WARNING RETRY\n")
					shouldRetry = true
				}
			}
			if !shouldRetry {
				// worked! no need to retry
				log.Printf("authreq.retryReq RETRY LOOP SUCCESS")
				return resp_body, code, resp_err
			}
		}
		e := fmt.Sprintf("authreq.retryReq: failed retries on %s:%s",
			amzTarget, string(reqJSON))
		return nil, 0, errors.New(e)
	}
}
