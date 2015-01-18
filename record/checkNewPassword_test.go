package record

import (
	. "github.com/cgentry/gus/ecode"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPasswords(t *testing.T) {
	user := NewTestUser()
	Convey("Check with too short passwords", t, func() {
		So(user.CheckNewPassword(""), ShouldEqual, ErrPasswordTooShort)
		So(user.CheckNewPassword("12345"), ShouldEqual, ErrPasswordTooShort)
	})
	Convey("Check with password", t, func() {
		So(user.CheckNewPassword("password"), ShouldEqual, ErrPasswordTooSimple)
	})

	Convey("Check with password OK", t, func() {
		So(user.CheckNewPassword("Th$s1s0k4Apsswd"), ShouldEqual, nil)
	})
}
