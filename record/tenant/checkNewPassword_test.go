package tenant

import (
	. "github.com/cgentry/gus/ecode"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPasswords(t *testing.T) {
	Convey("Check with too short passwords", t, func() {
		So(CheckNewPassword(""), ShouldEqual, ErrPasswordTooShort)
		So(CheckNewPassword("12345"), ShouldEqual, ErrPasswordTooShort)
	})
	Convey("Check with password", t, func() {
		So(CheckNewPassword("password"), ShouldEqual, ErrPasswordTooSimple)
	})

	Convey("Check with password OK", t, func() {
		So(CheckNewPassword("Th$s1s0k4Apsswd"), ShouldEqual, nil)
	})
}
