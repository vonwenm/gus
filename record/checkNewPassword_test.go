package record

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPasswords(t *testing.T) {
	user := NewTestUser()
	Convey("Check with too short passwords", t, func() {
		So(user.CheckNewPassword(""), ShouldEqual, USER_PASSWD_TOO_SHORT)
		So(user.CheckNewPassword("12345"), ShouldEqual, USER_PASSWD_TOO_SHORT)
	})
	Convey("Check with password", t, func() {
		So(user.CheckNewPassword("password"), ShouldEqual, USER_PASSWORD_TOO_SIMPLE)
	})

	Convey("Check with password OK", t, func() {
		So(user.CheckNewPassword("Th$s1s0k4Apsswd"), ShouldEqual, USER_OK)
	})
}
