package tenant

import (
	_ "github.com/cgentry/gus/encryption/drivers/plaintext"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	//"fmt"
	"time"
)

func TestNewUser(t *testing.T) {
	tuser := NewUser()
	tuser.SetDomain("Domainxyz")
	if tuser.Domain != "Domainxyz" {
		t.Errorf("Domain not correct: '%s'", tuser.Domain)
	}
}

func TestAuthenticate(t *testing.T) {
	pwd := "TestingPassvord"
	tuser := NewUser()
	tuser.SetDomain("dom")
	tuser.SetID(1234)
	tuser.SetPassword(pwd)

	Convey("Authenticatin test", t, func() {
		err := tuser.Login(pwd)
		So(err, ShouldBeNil)

	})
}

func TestCreateSalt(t *testing.T) {
	Convey("Multiple salts", t, func() {
		rtn1 := CreateSalt(20)
		So(rtn1, ShouldNotBeBlank)
		So(len(rtn1), ShouldEqual, 40)
		rtn2 := CreateSalt(20)
		So(rtn1, ShouldNotEqual, rtn2)

		rtn3 := CreateSalt(10)
		So(len(rtn3), ShouldEqual, 20)

		rtn4 := CreateSalt(9)
		So(len(rtn4), ShouldEqual, 18)

		rtn5 := CreateSalt(0)
		So(len(rtn5), ShouldEqual, 0)
	})
}

func TestGuid(t *testing.T) {
	tuser := NewUser()

	Convey("Guid test", t, func() {

		So(tuser.Guid, ShouldNotBeBlank)
		So(len(tuser.Guid), ShouldEqual, 36)
		So(tuser.Salt, ShouldNotBeBlank)
		So(len(tuser.Salt), ShouldEqual, 64)
	})
}

func TestLogin(t *testing.T) {
	pwd := "TestingPassvord"
	tuser := NewUser()
	tuser.SetDomain("dom")
	tuser.SetID(1234)
	tuser.SetPassword(pwd)
	Convey("Test good login", t, func() {
		err := tuser.Login(pwd)
		So(err, ShouldBeNil)

		err = tuser.Authenticate(tuser.Token)
		So(err, ShouldEqual, nil)

		err = tuser.Authenticate(``)
		So(err, ShouldNotEqual, nil)
	})
	Convey("Test bad login", t, func() {
		err := tuser.Login(`this isn't going to work`)
		So(err, ShouldNotBeNil)
		So(tuser.Token, ShouldBeBlank)
		So(tuser.FailCount, ShouldEqual, 1)
	})

}
func TestExpiredLogin(t *testing.T) {
	pwd := "TestingPassvord"
	tuser := NewUser()
	tuser.SetDomain("dom")
	tuser.SetID(1234)
	tuser.SetPassword(pwd)
	now := time.Now()
	Convey("Test good login", t, func() {
		err := tuser.Login(pwd)
		So(err, ShouldBeNil)

		err = tuser.Authenticate(tuser.Token)
		So(err, ShouldEqual, nil)
		Convey("Test Timeout", func() {
			tuser.TimeoutAt = now.Add(-1 * time.Second)
			err = tuser.Authenticate(tuser.Token)
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Timeout with MAX session length", func() {
			tuser.MaxSessionAt = now.Add(-1 * time.Second)
			err = tuser.Authenticate(tuser.Token)
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Timeout when TimeoutAt passed", func() {
			tuser.SetTimeoutAt(now.Add(-1 * time.Second))
			err = tuser.Authenticate(tuser.Token)
			So(err, ShouldNotEqual, nil)
		})

	})
}

func TestChangePassword(t *testing.T) {
	pwd := "TestingPassvord"
	now := time.Now()

	Convey("Test good login", t, func() {
		tuser := NewUser()
		tuser.SetDomain("dom")
		tuser.SetID(12345)
		tuser.SetPassword(pwd)
		err := tuser.Login(pwd)

		So(err, ShouldBeNil)

		Convey("Test FailChangePassword", func() {
			err = tuser.ChangePassword(`5544668899`, `abcdefghij`)
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test ChangePassword OK", func() {
			err = tuser.ChangePassword(pwd, `NewPassword`)
			So(err, ShouldEqual, nil)
		})
		Convey("Test too short password", func() {
			err = tuser.ChangePassword(pwd, `x`)
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Change When Not logged in", func() {
			tuser.Logout()
			err = tuser.ChangePassword(pwd, `NewPassword`)
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Change with MAX time passed", func() {
			tuser.MaxSessionAt = now.Add(-1 * time.Second)
			err = tuser.Authenticate(tuser.Token)
			So(err, ShouldNotEqual, nil)
			err = tuser.ChangePassword(pwd, `NewPassword`)
			So(err, ShouldNotEqual, nil)
		})

	})
}
