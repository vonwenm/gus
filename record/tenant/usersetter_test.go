package tenant

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/cgentry/gus/record/configure"
	"testing"
	"time"
)

func TestSetID(t *testing.T) {

	Convey("Check setting values of ID", t, func() {
		user := NewUser()
		user.SetDomain("test1")

		Convey("Check for simple ID set", func() {
			err := user.SetID(1)
			So(err, ShouldBeNil)
			err = user.SetID(1000222)
			So(err, ShouldNotBeNil)
		})

		Convey("Check for zero and negative ID set", func() {
			err := user.SetID(0)
			So(err, ShouldNotBeNil)
			err = user.SetID(-1)
			So(err, ShouldNotBeNil)
		})

		Convey("Check values are saved - 1", func() {
			err := user.SetID(1)
			So(err, ShouldBeNil)
			So(user.Id, ShouldEqual, 1)
		})
	})
}

func TestSetGuid(t *testing.T) {

	Convey("Check setting values of Email", t, func() {
		user := NewUser()
		user.SetDomain("test2")

		Convey("Check for simple Guid set", func() {
			err := user.SetGuid("3F2504E0-4F89-41D3-9A0C-0305E82C3301")
			So(err, ShouldBeNil)
			err = user.SetGuid("12")
			So(err, ShouldNotBeNil)
		})

		Convey("Check values are saved ", func() {
			err := user.SetGuid("3F2504E0-4F89-41D3-9A0C-0305E82C3301")
			So(err, ShouldBeNil)
			So(user.Guid, ShouldEqual, "3F2504E0-4F89-41D3-9A0C-0305E82C3301")
		})
	})
}

func TestSetPasswrd(t *testing.T) {

	Convey("Check setting values of password", t, func() {
		user := NewUser()
		user.SetDomain("test3")

		Convey("Check for simple password set", func() {
			err := user.SetPassword("12345678")
			So(err, ShouldBeNil)
			err = user.SetPassword("12")
			So(err, ShouldNotBeNil)
		})

		Convey("Check values are saved ", func() {
			err := user.SetPassword("123456")
			So(err, ShouldBeNil)

			So(user.CheckPassword("123456"), ShouldEqual, nil)
		})
	})
}

func TestSetName(t *testing.T) {
	Convey("Check name setting", t, func() {
		user := NewUser()
		err := user.SetName("Hello")
		So(err, ShouldBeNil)
		So(user.FullName, ShouldEqual, "Hello")
		err = user.SetName(``)
		So(err, ShouldNotBeNil)
		So(user.FullName, ShouldEqual, "Hello")
	})
}

func TestDateSetters(t *testing.T) {
	user := NewUser()
	nowStr := time.Now().Format(configure.USER_TIME_STR)
	now, err := time.Parse(configure.USER_TIME_STR, nowStr)

	Convey("Check setters setting", t, func() {
		So(err, ShouldBeNil)

		err = user.SetCreatedAt(now)
		So(err, ShouldBeNil)
		So(user.CreatedAt.Equal(now), ShouldBeTrue)
		So(user.GetCreatedAtStr(), ShouldEqual, nowStr)

		err = user.SetDeletedAt(now)
		So(err, ShouldBeNil)
		So(user.DeletedAt.Equal(now), ShouldBeTrue)
		So(user.GetDeletedAtStr(), ShouldEqual, nowStr)

		err = user.SetLastAuthAt(now)
		So(err, ShouldBeNil)
		So(user.LastAuthAt.Equal(now), ShouldBeTrue)
		So(user.GetLastAuthAtStr(), ShouldEqual, nowStr)

		err = user.SetLoginAt(now)
		So(err, ShouldBeNil)
		So(user.LoginAt.Equal(now), ShouldBeTrue)
		So(user.GetLoginAtStr(), ShouldEqual, nowStr)
	})
}
