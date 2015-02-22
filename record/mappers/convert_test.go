package mappers

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/cgentry/gus/record/configure"
	"testing"
)

func TestConversions(t *testing.T) {
	Convey("Check StrToBool - good values", t, func() {
		So(StrToBool("true", false), ShouldBeTrue)
		So(StrToBool("false", true), ShouldBeFalse)
		So(StrToBool("1", false), ShouldBeTrue)
		So(StrToBool("0", true), ShouldBeFalse)
		So(StrToBool("yes", false), ShouldBeTrue)
		So(StrToBool("no", true), ShouldBeFalse)
		So(StrToBool("ok", false), ShouldBeTrue)
	})
	Convey("Check StrToBool - bad values, use defaults", t, func() {
		So(StrToBool("yep", false), ShouldBeFalse)
		So(StrToBool("nope", true), ShouldBeTrue)
	})
	Convey("Check StrToInt - Good value", t, func() {
		So(StrToInt("0"), ShouldEqual, 0)
		So(StrToInt("1"), ShouldEqual, 1)
		So(StrToInt("1000"), ShouldEqual, 1000)
		So(StrToInt("-10"), ShouldEqual, -10)
		So(StrToInt("+10"), ShouldEqual, 10)
	})
	Convey("Check StrToInt - bad value", t, func() {
		So(StrToInt("x"), ShouldEqual, 0)
		So(StrToInt("a100"), ShouldEqual, 0)
		So(StrToInt("0x100"), ShouldEqual, 0)
	})

	Convey("Check StrToTime - good values", t, func() {
		nowStr := "2014-12-29T12:37:36Z"
		now := StrToTime(nowStr)

		So(now.IsZero(), ShouldBeFalse)
		So(now.Format(configure.USER_TIME_STR), ShouldEqual, `2014-12-29T12:37:36Z`)
	})

	Convey("Check StrToTime - bad values", t, func() {
		nowStr := "Miller time"
		now := StrToTime(nowStr)

		So(now.IsZero(), ShouldBeTrue)

	})
}
