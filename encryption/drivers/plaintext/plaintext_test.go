package plaintext

import (
	"github.com/cgentry/gus/encryption"
	"github.com/cgentry/gus/record"
	"testing"
)

func TestGenerate(t *testing.T) {

	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("hello", user.Salt)
	if pwd != "hello;"+user.Salt {
		t.Errorf("Passwords don't match encrypted")
	}
}

func TestRepeatable(t *testing.T) {
	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("123456", user.Salt)
	pwd2 := encryption.GetDriver().EncryptPassword("123456", user.Salt)
	if pwd != pwd2 {
		t.Errorf("Passwords didn't match: '%s' and '%s'", pwd, pwd2)
	}

}

func TestIsLongEnough(t *testing.T) {
	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("hello", user.Salt)
	pwdLen := len(pwd)
	sbLen := len("hello;" + user.Salt)
	if pwdLen != sbLen {
		t.Errorf("PWD isn't long enough %d", pwdLen)
	}
}

func TestSimilarUserDifferntPwd(t *testing.T) {
	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("123456", user.Salt)
	user2 := record.NewTestUser()
	pwd2 := encryption.GetDriver().EncryptPassword("123456", user2.Salt)
	if pwd == pwd2 {
		t.Errorf("Passwords for different users should not match: '%s' and '%s'", pwd, pwd2)
	}
}

func TestAfterChangingSalt(t *testing.T) {
	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("123456", user.Salt)
	encryption.GetDriver().SetInternalSalt("hello - this should screw up password")
	pwd2 := encryption.GetDriver().EncryptPassword("123456", user.Salt)

	if pwd == pwd2 {
		t.Errorf("Passwords with different salts should not match: '%s' and '%s'", pwd, pwd2)
	}
}
