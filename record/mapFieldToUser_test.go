package record

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestID(t *testing.T) {
	Convey("Check setting values of ID", t, func() {
		user := NewUser()
		user.SetDomain("test2")

		Convey("Check for simple ID set", func() {
			_, err := user.MapFieldToUser("id", "1")
			So(err, ShouldBeNil)
			_, err = user.MapFieldToUser("id", "1000222")
			So(err, ShouldNotBeNil)
		})

		Convey("Check for zero and negative ID set", func() {
			found, err := user.MapFieldToUser("id", "0")
			So(err, ShouldNotBeNil)
			So(found, ShouldBeTrue)
			_, err = user.MapFieldToUser("id", "-1")
			So(err, ShouldNotBeNil)
		})

		Convey("Check values are saved - 1", func() {
			_, err := user.MapFieldToUser("id", "1")
			So(err, ShouldBeNil)
			So(user.GetID(), ShouldEqual, 1)
		})
	})
}

func TestName(t *testing.T) {
	Convey("Check setting values of Name", t, func() {
		user := NewUser()
		user.SetDomain("test2")

		Convey("Check for simple Name set", func() {
			_, err := user.MapFieldToUser("Name", "John Doe")
			So(err, ShouldBeNil)
			So("John Doe", ShouldEqual, user.GetName())
			user.MapFieldToUser("NAME", " Test2 ")
			So("Test2", ShouldEqual, user.GetName())
		})

		Convey("Check for Alternate Name set", func() {
			_, err := user.MapFieldToUser("fullName", "test3")
			So(err, ShouldBeNil)
			So("test3", ShouldEqual, user.GetName())
		})
	})
}

func TestNotFound(t *testing.T) {
	Convey("Pass unknown field", t, func() {
		user := NewUser()
		user.SetDomain("test3")
		found, err := user.MapFieldToUser("noway", "John Doe")
		So(err, ShouldNotBeNil)
		So(found, ShouldBeFalse)

	})
}
