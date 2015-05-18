package tasks

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"testing"
	"time"
)

func TestCanonicalRequest(t *testing.T) {

	body := []byte(`{"Test":"Stuff"}`)
	var h256 hash.Hash = sha256.New()
	h256.Write(body)
	hexPayload := string(hex.EncodeToString([]byte(h256.Sum(nil))))

	r := CanonicalRequest("dynamodb.us-east-1.amazonaws.com", "80", "x-amz-date: 20110909T233600Z", "", hexPayload)
	h256.Write([]byte(r))
	hashedCanonicalRequest := string(hex.EncodeToString([]byte(h256.Sum(nil))))
	const expected = "0d336e0c2a7878efc381e1fba2606c885afa6c5ad26480d39b56c1c2503271ba"
	if expected != hashedCanonicalRequest {
		t.Errorf("canonical request unexpected")
	}
}

func TestStringToSign(t *testing.T) {
	body := []byte(`{"Test":"Stuff"}`)
	var h256 hash.Hash = sha256.New()
	h256.Write(body)
	hexPayload := string(hex.EncodeToString([]byte(h256.Sum(nil))))

	r := CanonicalRequest("dynamodb.us-east-1.amazonaws.com", "80", "x-amz-date: 20110909T233600Z", "", hexPayload)
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	test_time, t_err := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
	if t_err != nil {
		t.Errorf(t_err.Error())
	}
	str2sign := String2Sign(test_time, r, "us-east-1", "dynamodb")
	h256.Write([]byte(str2sign))
	hashedString2Sign := string(hex.EncodeToString([]byte(h256.Sum(nil))))

	const expected = "3a958ac6ec0702c4d9b05f2b360762893bdb47440f133ffcdfde1d263555479e"
	if expected != hashedString2Sign {
		t.Errorf("string 2 sign unexpected")
	}
}
