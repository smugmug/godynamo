package conf_file

import (
	"testing"
)

func TestReadConfFile(t *testing.T) {
	_, err := ReadConfFile("./test_aws-config.json")
	if err != nil {
		t.Errorf(err.Error())
	}
}
