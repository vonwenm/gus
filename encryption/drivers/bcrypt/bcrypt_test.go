package bcrypt

import (
	"github.com/cgentry/gus/encryption"
	"github.com/cgentry/gus/record"
	"testing"
)

func TestGenerate(t *testing.T) {

	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("hello", user.Salt)
	if pwd == "hello" {
		t.Errorf("pwd didn't get encrypted")
	}
}

func TestCompare(t *testing.T) {
	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("123456", user.Salt)

	if !encryption.GetDriver().ComparePasswords(pwd, "123456", user.Salt) {
		t.Errorf("Passwords didn't match")
	}

}

func TestIsLongEnough(t *testing.T) {
	user := record.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword("hello", user.Salt)
	pwdLen := len(pwd)
	if pwdLen < 60 {
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
