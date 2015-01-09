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
	Convey("Create User", t, func() {
		So(err, ShouldBeNil)

		user := record.NewTestUser()
		user.SetDomain("Register")
		user.SetToken("TestToken")
		user.SetName("Just a test name")
		user.SetEmail("et@home.com")
		user.SetLoginName("justlogin")

		err = db.RegisterUser(user) // Register new user
		So(err, ShouldBeNil)

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
		So(err, ShouldBeNil)
		user, err = db.FetchUserByGuid(user.GetGuid())

		err = db.UserLogout(user)
		So(err, ShouldBeNil)
		err = db.UserLogout(user)
		So(err.Error(), ShouldEqual, storage.ErrUserNotLoggedIn.Error())

		db.Close()
	})

}

/*
func TestBadFetches(t *testing.T) {

	db, err := storage.Open(STORAGE_IDENTITY, STORE_LOCAL)

	Convey("Create User", t, func() {
		So(err, ShouldBeNil)

		// FETCH BY EMAIL
		user2, err := db.FetchUserByGuid(`noguid`)
		So(err, ShouldNotBeNil)
		So(user2, ShouldBeNil)

		// Fetch by TOKEN
		user3, err := db.FetchUserByToken("notken")
		So(err, ShouldNotBeNil)
		So(user3, ShouldBeNil)

		// FETCH BY EMAIL
		user4, err := db.FetchUserByEmail(`nomail@home.com`)
		So(err, ShouldNotBeNil)
		So(user4, ShouldBeNil)

		// FETCH BY EMAIL
		user5, err := db.FetchUserByLogin(`nologin`)
		So(err, ShouldNotBeNil)
		So(user5, ShouldBeNil)
		db.Close()
	})

}

func TestDuplicateKey(t *testing.T) {

	Convey("Re-create User", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, STORE_LOCAL)
		So(err, ShouldBeNil)
		user := record.NewTestUser()
		user.SetDomain("Register")
		user.SetToken("TestToken")
		user.SetName("Just a test name")
		user.SetEmail("et@home.com")
		user.SetLoginName("justlogin")
		err = db.RegisterUser(user) // Register new user
		So(err.Error(), ShouldEqual, storage.ErrDuplicateEmail.Error())

		user.SetEmail("better@home.com")
		err = db.RegisterUser(user)
		So(err.Error(), ShouldEqual, storage.ErrDuplicateLogin.Error())
		db.Close()

	})
}

func TestClose(t *testing.T) {

	Convey("Open and close connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, STORE_LOCAL)
		So(err, ShouldBeNil)
		m := db.GetStorageConnector().(*SqliteConn)
		db.Close()
		So(db.IsOpen(), ShouldBeFalse)

		// Test bad fetch...
		// FETCH BY EMAIL
		user2, err := db.FetchUserByGuid(`noguid`)
		So(err, ShouldNotBeNil)
		So(user2, ShouldBeNil)

		// Try registering
		user := record.NewTestUser()
		err = db.RegisterUser(user) // Register new user
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "not open")

		err = m.Close()
		So(err, ShouldBeNil)

		// Try low-level fetch
		user2, err = m.FetchUserByGuid(`noguid`)
		So(user2, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "not open")

		// Try low-level fetch
		db, err = storage.Open(STORAGE_IDENTITY, STORE_LOCAL)
		So(err, ShouldBeNil)
		m = db.GetStorageConnector().(*SqliteConn)
		clearSqliteTest()
		user2, err = m.FetchUserByGuid(`noguid`)
		So(user2, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "table")
	})
}
func TestGetLowLevelHandle(t *testing.T) {

	Convey("Check low-level handle access", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, STORE_LOCAL)
		So(err, ShouldBeNil)
		m := db.GetStorageConnector().(*SqliteConn)

		dbh := m.GetRawHandle()
		So(dbh, ShouldNotBeNil)

		dbhMem, err := sql.Open(DRIVER_IDENTITY, ":memory:")
		So(err, ShouldBeNil)
		So(dbh, ShouldHaveSameTypeAs, dbhMem)
		dbhMem.Close()

		db.Close()
		So(db.IsOpen(), ShouldBeFalse)
		user := record.NewTestUser()
		err = db.RegisterUser(user) // Register new user
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "not open")

		dbh = m.GetRawHandle()
		So(dbh, ShouldBeNil)
	})
}

func TestSaveRecord(t *testing.T) {

	Convey("Test Saving Record", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, STORE_LOCAL)
		So(err, ShouldBeNil)
		defer clearSqliteTest()
		err = db.CreateStore()
		So(err, ShouldBeNil)

		user := record.NewTestUser()
		err = db.RegisterUser(user)
		So(err, ShouldBeNil)

		user.SetToken(`SaveRecordToken`)
		err = db.UserUpdate(user)
		So(err, ShouldBeNil)
		user, err = db.FetchUserByToken(`SaveRecordToken`)
		So(err, ShouldBeNil)
		So(user, ShouldNotBeNil)
		So(user.GetToken(), ShouldEqual, `SaveRecordToken`)
		db.Close()
	})

}
*/
