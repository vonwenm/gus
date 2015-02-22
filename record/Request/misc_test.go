package request

import (
	"github.com/cgentry/gus/ecode"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	Convey("Test check and create", t, func() {
		entity := NewAuthenticate()
		err := entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingToken)

		entity.Token = "HI"
		err = entity.Check()
		So(err, ShouldBeNil)

		entity.SetStamp(time.Unix(0, 0))
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrRequestNoTimestamp)

	})
}

func TestLogin(t *testing.T) {
	Convey("Test check and create", t, func() {
		entity := NewLogin()
		err := entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingLogin)

		entity.Login = "login"
		So(entity.GetLogin(), ShouldEqual, "login")
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingPassword)

		entity.Password = "pwd"
		So(entity.GetPassword(), ShouldEqual, "pwd")
		err = entity.Check()
		So(err, ShouldBeNil)

		entity.SetStamp(time.Unix(0, 0))
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrRequestNoTimestamp)

	})
}
func TestLogout(t *testing.T) {
	Convey("Test check and create", t, func() {
		entity := NewLogout()
		err := entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingToken)

		entity.Token = "HI"
		err = entity.Check()
		So(err, ShouldBeNil)

		entity.SetStamp(time.Unix(0, 0))
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrRequestNoTimestamp)

	})
}

func TestRegister(t *testing.T) {
	Convey("Test check and create", t, func() {
		entity := NewRegister()
		err := entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingEmail)

		entity.Email = "e@mail.com"
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingLogin)

		entity.Login = "login"
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingName)

		entity.Name = "Name"
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingPassword)

		entity.Password = "pwd"
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrPasswordTooShort)

		entity.Password = "this is a long password"
		err = entity.Check()
		So(err, ShouldBeNil)

		entity.SetStamp(time.Unix(0, 0))
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrRequestNoTimestamp)

	})
}

func TestUpdate(t *testing.T) {
	Convey("Test check and create", t, func() {
		entity := NewUpdate()
		err := entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingPasswordNew)

		entity.NewPassword = "pwd"
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrPasswordTooShort)

		entity.NewPassword = "Long password"
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMissingPassword)

		entity.OldPassword = "Long password"
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrMatchingPassword)

		entity.NewPassword = "New Long password"
		err = entity.Check()
		So(err, ShouldBeNil)

		entity.SetStamp(time.Unix(0, 0))
		err = entity.Check()
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.ErrRequestNoTimestamp)

	})
}
