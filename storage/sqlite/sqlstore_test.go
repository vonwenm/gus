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
	//"fmt"
)

const STORE_LOCAL = "/tmp/test_store.sqlite3"

func clearSqliteTest() {
	if STORE_LOCAL != ":memory:" {
		os.Remove(STORE_LOCAL)
	}
}

func TestSimpleRegisterCycle(t *testing.T) {
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
		So(serr.Error(), ShouldBeBlank)

		// FETCH BY EMAIL
		user2, err := db.FetchUserByGuid(user.GetGuid())
		So(err, ShouldBeNil)
		So(user2.GetDomain(), ShouldEqual, user.GetDomain())
		So(user2.GetToken(), ShouldEqual, user.GetToken())
		So(user2.GetName(), ShouldEqual, user.GetName())

		// Fetch by TOKEN
		user3, err := db.FetchUserByToken("TestToken")
		So(err, ShouldBeNil)
		So(user3.GetDomain(), ShouldEqual, user.GetDomain())
		So(user3.GetToken(), ShouldEqual, user.GetToken())
		So(user3.GetName(), ShouldEqual, user.GetName())

		// FETCH BY EMAIL
		user4, err := db.FetchUserByEmail(`et@home.com`)
		So(err, ShouldBeNil)
		So(user4.GetDomain(), ShouldEqual, user.GetDomain())
		So(user4.GetToken(), ShouldEqual, user.GetToken())
		So(user4.GetName(), ShouldEqual, user.GetName())

		// FETCH BY EMAIL
		user5, err := db.FetchUserByLogin(`justlogin`)
		So(err, ShouldBeNil)
		So(user5.GetDomain(), ShouldEqual, user.GetDomain())
		So(user5.GetToken(), ShouldEqual, user.GetToken())
		So(user5.GetName(), ShouldEqual, user.GetName())

		// By default, a registered user is NOT logged in...
		err = db.UserLogin(user)
		So(err.Error(), ShouldBeBlank)
		user, err = db.FetchUserByGuid(user.GetGuid())

		err = db.UserLogout(user)
		So(err.Error(), ShouldBeBlank )
		err = db.UserLogout(user)
		So(err.Error(), ShouldEqual, storage.ErrUserNotLoggedIn.Error())

		db.Close()
	})

}

