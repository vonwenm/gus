package record

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
			So(user.GetID(), ShouldEqual, 1)
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
			So(user.GetGuid(), ShouldEqual, "3F2504E0-4F89-41D3-9A0C-0305E82C3301")
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

			So(user.CheckPassword("123456"), ShouldEqual, USER_OK)
		})
	})
}
