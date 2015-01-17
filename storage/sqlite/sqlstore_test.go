// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	//"database/sql"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	. "github.com/cgentry/gus/ecode"
	"time"
)

const STORE_LOCAL = "/tmp/test_store.sqlite3"

func clearSqliteTest() {
	if STORE_LOCAL != ":memory:" {
		os.Remove(STORE_LOCAL)
	}
}

func TestSimpleRegisterCycle(t *testing.T) {
	var compareTime1, compareTime2 time.Time
	clearSqliteTest()
	db, err := storage.Open(STORAGE_IDENTITY, STORE_LOCAL)
	db.CreateStore()
	defer clearSqliteTest()

	Convey("Create User", t, func() {
		So(err, ShouldBeNil)

		user := record.NewTestUser()
		user.SetDomain("Register")
		user.SetToken("TestToken")
		user.SetName("Just a test name")
		user.SetEmail("et@home.com")
		user.SetLoginName("justlogin")

		serr := db.RegisterUser(user) // Register new user
		So(serr,ShouldBeNil)

		// FETCH BY EMAIL
		user2, err := db.FetchUserByGuid(user.Guid)
		So(err, ShouldBeNil)
		So(user2.Domain, ShouldEqual, user.Domain)
		So(user2.Token, ShouldEqual, user.Token)
		So(user2.FullName, ShouldEqual, user.FullName)

		// Fetch by TOKEN
		user3, err := db.FetchUserByToken("TestToken")
		So(err, ShouldBeNil)
		So(user3.Domain, ShouldEqual, user.Domain)
		So(user3.Token, ShouldEqual, user.Token)
		So(user3.FullName, ShouldEqual, user.FullName)

		// FETCH BY EMAIL
		user4, err := db.FetchUserByEmail(`et@home.com`)
		So(err, ShouldBeNil)
		So(user4.Domain, ShouldEqual, user.Domain)
		So(user4.Token, ShouldEqual, user.Token)
		So(user4.FullName, ShouldEqual, user.FullName)

		// FETCH BY EMAIL
		user5, err := db.FetchUserByLogin(`justlogin`)
		So(err, ShouldBeNil)
		So(user5.Domain, ShouldEqual, user.Domain)
		So(user5.Token, ShouldEqual, user.Token)
		So(user5.FullName, ShouldEqual, user.FullName)

		// By default, a registered user is NOT logged in...
		compareTime1 = user.LoginAt
		compareTime2= user.UpdatedAt
		err = db.UserLogin(user)
		So(err, ShouldBeNil )
		user, err = db.FetchUserByGuid(user.Guid)
		So( user.IsLoggedIn, ShouldBeTrue)
		So(user.LoginAt.Equal(compareTime1), ShouldBeFalse)
		So(user.UpdatedAt.Equal(compareTime2), ShouldBeFalse)
		So(user.LoginAt.Equal(user.UpdatedAt), ShouldBeTrue)

		compareTime1 = user.LogoutAt
		compareTime2= user.UpdatedAt

		err = db.UserLogout(user)
		So(err, ShouldBeNil)
		So(user.IsLoggedIn, ShouldBeFalse)
		So(user.LogoutAt.Equal(compareTime1), ShouldBeFalse)
		So(user.UpdatedAt.Equal(compareTime2), ShouldBeFalse)
		So(user.LogoutAt.Equal(user.UpdatedAt), ShouldBeTrue)
		err = db.UserLogout(user)
		So(err, ShouldNotBeNil )
		So(err.Error(), ShouldEqual, ErrUserNotLoggedIn.Error())

		db.Close()
	})

}
