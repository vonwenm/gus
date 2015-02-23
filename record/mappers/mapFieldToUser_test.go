package mappers

import (
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/record/tenant"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestID(t *testing.T) {
	Convey("Check setting values of ID", t, func() {
		user := tenant.NewUser()
		user.SetDomain("test2")

		Convey("Check for simple ID set", func() {
			_, err := UserField(user, "id", "1")
			So(err, ShouldBeNil)
			_, err = UserField(user, "id", "1000222")
			So(err, ShouldNotBeNil)
		})

		Convey("Check for zero and negative ID set", func() {
			found, err := UserField(user, "id", "0")
			So(err, ShouldNotBeNil)
			So(found, ShouldBeTrue)
			_, err = UserField(user, "id", "-1")
			So(err, ShouldNotBeNil)
		})

		Convey("Check values are saved - 1", func() {
			_, err := UserField(user, "id", "1")
			So(err, ShouldBeNil)
			So(user.GetID(), ShouldEqual, 1)
		})
	})
}

func TestName(t *testing.T) {
	Convey("Check setting values of Name", t, func() {
		user := tenant.NewUser()
		user.SetDomain("test2")

		Convey("Check for simple Name set", func() {
			_, err := UserField(user, "Name", "John Doe")
			So(err, ShouldBeNil)
			So("John Doe", ShouldEqual, user.FullName)
			UserField(user, "NAME", " Test2 ")
			So("Test2", ShouldEqual, user.FullName)
		})

		Convey("Check for Alternate Name set", func() {
			_, err := UserField(user, "fullName", "test3")
			So(err, ShouldBeNil)
			So("test3", ShouldEqual, user.FullName)
		})
	})
}

func TestNotFound(t *testing.T) {
	Convey("Pass unknown field", t, func() {
		user := tenant.NewUser()
		user.SetDomain("test3")
		found, err := UserField(user, "noway", "John Doe")
		So(err, ShouldNotBeNil)
		So(found, ShouldBeFalse)

	})
}

func TestGeneralFields(t *testing.T) {
	user := tenant.NewUser()
	nowStr := time.Now().Format(configure.USER_TIME_STR)
	now, terr := time.Parse(configure.USER_TIME_STR, nowStr)

	Convey("Pass field", t, func() {
		So(terr, ShouldBeNil)

		found, err := UserField(user, `email`, `myemail@common.com`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Email, ShouldEqual, `myemail@common.com`)

		found, err = UserField(user, `guid`, `123456789012345678901234567890123456789012345678901234567890`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Guid, ShouldEqual, `123456789012345678901234567890123456789012345678901234567890`)

		found, err = UserField(user, `caller`, `A23456789012345678901234567890123456789012345678901234567890`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Guid, ShouldEqual, `A23456789012345678901234567890123456789012345678901234567890`)

		found, err = UserField(user, `domain`, `MyDomain`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Domain, ShouldEqual, `MyDomain`)

		found, err = UserField(user, `password`, `MyPassword`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Password, ShouldEqual, `MyPassword`)

		found, err = UserField(user, `token`, `MyToken`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Token, ShouldEqual, `MyToken`)

		found, err = UserField(user, `salt`, `saltsaltsaltsaltsaltsaltsalt`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.Salt, ShouldEqual, `saltsaltsaltsaltsaltsaltsalt`)

		found, err = UserField(user, `ISactive`, `true`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsActive, ShouldBeTrue)

		found, err = UserField(user, `isLOGGEDin`, `true`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsLoggedIn, ShouldBeTrue)

		found, err = UserField(user, `issystem`, `true`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsSystem, ShouldBeTrue)

		found, err = UserField(user, `issystem`, `false`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsSystem, ShouldBeFalse)

		found, err = UserField(user, `issystem`, `false`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.IsSystem, ShouldBeFalse)

		found, err = UserField(user, `loginname`, `MyLoginName`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LoginName, ShouldEqual, `MyLoginName`)

		found, err = UserField(user, `login`, `MyNewLoginName`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LoginName, ShouldEqual, `MyNewLoginName`)

		found, err = UserField(user, `loginat`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LoginAt.Equal(now), ShouldBeTrue)
		So(user.GetLoginAtStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `logoutat`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LogoutAt.Equal(now), ShouldBeTrue)
		So(user.GetLogoutAtStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `lastfailedat`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LastFailedAt.Equal(now), ShouldBeTrue)
		So(user.GetLastFailedAtStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `TimeoutAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.TimeoutAt.Equal(now), ShouldBeTrue)
		So(user.GetTimeoutStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `MaxSessionAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.MaxSessionAt.Equal(now), ShouldBeTrue)
		So(user.GetMaxSessionAtStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `UpdatedAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.UpdatedAt.Equal(now), ShouldBeTrue)
		So(user.GetUpdatedAtStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `LastAuthAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.LastAuthAt.Equal(now), ShouldBeTrue)
		So(user.GetLastAuthAtStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `DeletedAt`, nowStr)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.DeletedAt.Equal(now), ShouldBeTrue)
		So(user.GetDeletedAtStr(), ShouldEqual, nowStr)

		found, err = UserField(user, `FailCount`, `10`)
		So(err, ShouldBeNil)
		So(found, ShouldBeTrue)
		So(user.FailCount, ShouldEqual, 10)
		So(user.GetFailCountStr(), ShouldEqual, `10`)
	})
}
