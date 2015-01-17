package record

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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
			So("John Doe", ShouldEqual, user.FullName)
			user.MapFieldToUser("NAME", " Test2 ")
			So("Test2", ShouldEqual, user.FullName)
		})

		Convey("Check for Alternate Name set", func() {
			_, err := user.MapFieldToUser("fullName", "test3")
			So(err, ShouldBeNil)
			So("test3", ShouldEqual, user.FullName)
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

func TestGeneralFields(t *testing.T) {
	user := NewUser()
	nowStr := time.Now().Format(USER_TIME_STR)
	now, terr := time.Parse(USER_TIME_STR, nowStr)

	Convey("Pass field", t, func() {
		So(terr, ShouldBeNil)

		found, err := user.MapFieldToUser(`email`, `myemail@common.com`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Email, ShouldEqual, `myemail@common.com`)

		found, err = user.MapFieldToUser(`guid`, `123456789012345678901234567890123456789012345678901234567890`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Guid, ShouldEqual, `123456789012345678901234567890123456789012345678901234567890`)

		found, err = user.MapFieldToUser(`caller`, `A23456789012345678901234567890123456789012345678901234567890`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Guid, ShouldEqual, `A23456789012345678901234567890123456789012345678901234567890`)

		found, err = user.MapFieldToUser(`domain`, `MyDomain`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Domain, ShouldEqual, `MyDomain`)

		found, err = user.MapFieldToUser(`password`, `MyPassword`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Password, ShouldEqual, `MyPassword`)

		found, err = user.MapFieldToUser(`token`, `MyToken`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Token, ShouldEqual, `MyToken`)

		found, err = user.MapFieldToUser(`salt`, `saltsaltsaltsaltsaltsaltsalt`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Salt, ShouldEqual, `saltsaltsaltsaltsaltsaltsalt`)

		found, err = user.MapFieldToUser(`ISactive`, `true`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsActive, ShouldBeTrue)

		found, err = user.MapFieldToUser(`isLOGGEDin`, `true`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsLoggedIn, ShouldBeTrue)

		found, err = user.MapFieldToUser(`issystem`, `true`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsSystem, ShouldBeTrue)

		found, err = user.MapFieldToUser(`issystem`, `false`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsSystem, ShouldBeFalse)

		found, err = user.MapFieldToUser(`issystem`, `false`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsSystem, ShouldBeFalse)

		found, err = user.MapFieldToUser(`loginname`, `MyLoginName`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LoginName, ShouldEqual, `MyLoginName`)

		found, err = user.MapFieldToUser(`login`, `MyNewLoginName`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LoginName, ShouldEqual, `MyNewLoginName`)

		found, err = user.MapFieldToUser(`loginat`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LoginAt.Equal(now), ShouldBeTrue)
		So(user.GetLoginAtStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`logoutat`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LogoutAt.Equal(now), ShouldBeTrue)
		So(user.GetLogoutAtStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`lastfailedat`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LastFailedAt.Equal(now), ShouldBeTrue)
		So(user.GetLastFailedAtStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`TimeoutAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.TimeoutAt.Equal(now), ShouldBeTrue)
		So(user.GetTimeoutStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`MaxSessionAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.MaxSessionAt.Equal(now), ShouldBeTrue)
		So(user.GetMaxSessionAtStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`UpdatedAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.UpdatedAt.Equal(now), ShouldBeTrue)
		So(user.GetUpdatedAtStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`LastAuthAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LastAuthAt.Equal(now), ShouldBeTrue)
		So(user.GetLastAuthAtStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`DeletedAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.DeletedAt.Equal(now), ShouldBeTrue)
		So(user.GetDeletedAtStr(), ShouldEqual, nowStr)

		found, err = user.MapFieldToUser(`FailCount`, `10`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.FailCount, ShouldEqual, 10)
		So(user.GetFailCountStr(), ShouldEqual, `10`)
	})
}
