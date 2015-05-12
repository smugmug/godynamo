package conf

import (
	roles_master "github.com/smugmug/goawsroles/roles_master"
	roles_simple "github.com/smugmug/goawsroles/roles_simple"
	"testing"
)

func TestRolesSimple(t *testing.T) {
	r := roles_simple.NewRolesSimple("myAccessKey", "mySecret", "myToken")
	var c AWS_Conf
	cp_err := c.CredentialsFromRoles(r)
	if cp_err != nil {
		t.Errorf(cp_err.Error())
	}
	if !c.UseIAM {
		t.Errorf("UseIAM should be true")
	}
}

func TestRolesMaster(t *testing.T) {
	r := roles_master.NewRolesMaster("myAccessKey", "mySecret")
	var c AWS_Conf
	cp_err := c.CredentialsFromRoles(r)
	if cp_err != nil {
		t.Errorf(cp_err.Error())
	}
	if c.UseIAM {
		t.Errorf("UseIAM should be false")
	}
}
