// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	//"database/sql"
	//. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	//"time"
)

const STORE_LOCAL = "/tmp/test_store.sqlite3"

func clearSqliteTest() {
	if STORE_LOCAL != ":memory:" {
		os.Remove(STORE_LOCAL)
	}
}

func TestSimpleRegisterCycle(t *testing.T) {
	//var compareTime1, compareTime2 time.Time
	clearSqliteTest()
	/*
		dbConn, err := storage.Open(STORAGE_IDENTITY, STORE_LOCAL)
		defer clearSqliteTest()
		defer dbConn.Close()
		dbConn.CreateStore()
	*/

	dbGeneralCon, err := NewSqliteDriver().Open(STORE_LOCAL, ``)

	Convey("Create User", t, func() {
		So(err, ShouldBeNil)
		defer clearSqliteTest()

		dbConn, ok := dbGeneralCon.(*SqliteConn) // To force getting at the raw calls...
		So(ok, ShouldBeTrue)
		dbConn.CreateStore()

		user := record.NewTestUser()
		user.SetDomain("Register")
		user.SetToken("TestToken")
		user.SetName("Just a test name")
		user.SetEmail("et@home.com")
		user.SetLoginName("justlogin")

		serr := dbConn.UserInsert(user) // Register new user
		So(serr, ShouldBeNil)

		// FETCH BY EMAIL
		user2, err := dbConn.UserFetch(storage.MATCH_ANY_DOMAIN, storage.FIELD_GUID, user.Guid)
		So(err, ShouldBeNil)
		So(user2.Domain, ShouldEqual, user.Domain)
		So(user2.Token, ShouldEqual, user.Token)
		So(user2.FullName, ShouldEqual, user.FullName)

		// Fetch by TOKEN
		user3, err := dbConn.UserFetch(storage.MATCH_ANY_DOMAIN, storage.FIELD_TOKEN, user.Token)
		So(err, ShouldBeNil)
		So(user3.Domain, ShouldEqual, user.Domain)
		So(user3.Token, ShouldEqual, user.Token)
		So(user3.FullName, ShouldEqual, user.FullName)

		// FETCH BY EMAIL
		user4, err := dbConn.UserFetch(user.Domain, storage.FIELD_EMAIL, user.Email)
		So(err, ShouldBeNil)
		So(user4.Domain, ShouldEqual, user.Domain)
		So(user4.Token, ShouldEqual, user.Token)
		So(user4.FullName, ShouldEqual, user.FullName)

		// FETCH BY Login name
		user5, err := dbConn.UserFetch(user.Domain, storage.FIELD_LOGIN, user.LoginName)
		So(err, ShouldBeNil)
		So(user5.Domain, ShouldEqual, user.Domain)
		So(user5.Token, ShouldEqual, user.Token)
		So(user5.FullName, ShouldEqual, user.FullName)
		/*
			// By default, a registered user is NOT logged in...
			compareTime1 = user.LoginAt
			compareTime2= user.UpdatedAt
			err = dbConn.UserLogin(user)
			So(err, ShouldBeNil )
			user, err = dbConn.FetchUserByGuid(user.Guid)
			So( user.IsLoggedIn, ShouldBeTrue)
			So(user.LoginAt.Equal(compareTime1), ShouldBeFalse)
			So(user.UpdatedAt.Equal(compareTime2), ShouldBeFalse)
			So(user.LoginAt.Equal(user.UpdatedAt), ShouldBeTrue)

			compareTime1 = user.LogoutAt
			compareTime2= user.UpdatedAt

			err = dbConn.UserLogout(user)
			So(err, ShouldBeNil)
			So(user.IsLoggedIn, ShouldBeFalse)
			So(user.LogoutAt.Equal(compareTime1), ShouldBeFalse)
			So(user.UpdatedAt.Equal(compareTime2), ShouldBeFalse)
			So(user.LogoutAt.Equal(user.UpdatedAt), ShouldBeTrue)
			err = dbConn.UserLogout(user)
			So(err, ShouldNotBeNil )
			So(err.Error(), ShouldEqual, ErrUserNotLoggedIn.Error())

			dbConn.Close()
		*/
	})

}
