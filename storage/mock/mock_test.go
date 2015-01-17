package mock

import (
	"errors"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMockRegisterDriver(t *testing.T) {
	Convey("Load the Mock storage", t, func() {
		So(storage.IsRegistered(STORAGE_IDENTITY), ShouldBeFalse)
		RegisterMockStore()
		So(storage.IsRegistered(STORAGE_IDENTITY), ShouldBeTrue)
	})
}
func TestCreateStore(t *testing.T) {
	Convey("Call CreateStore", t, func() {
		storage.ResetRegister()
		So(storage.IsRegistered(STORAGE_IDENTITY), ShouldBeFalse)
		RegisterMockStore()
		So(storage.IsRegistered(STORAGE_IDENTITY), ShouldBeTrue)

		db, err := storage.Open(STORAGE_IDENTITY, "test")
		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		m.Reset()
		So(m.WasCalled(`createstore`), ShouldBeFalse)
		So(m.WasCalledOnce(`createstore`), ShouldBeFalse)
		db.CreateStore()
		So(m.WasCalled(`createstore`), ShouldBeTrue)
		So(m.WasCalledOnce(`createstore`), ShouldBeTrue)
		db.CreateStore()
		So(m.WasCalled(`createstore`), ShouldBeTrue)
		So(m.WasCalledOnce(`createstore`), ShouldBeFalse)

		So(m.WasCalled(`close`), ShouldBeFalse)
	})
}
func TestMockOpenClose(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")
		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)

		So(db.Close(), ShouldBeNil)
		So(db.Close(), ShouldNotBeNil)
		m.Reset()
	})
}

func TestMockRegister(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		So(db.RegisterUser(u), ShouldBeNil)
		So(m.WasCalled(`RegisterUser`), ShouldBeTrue)
		rtnUser := m.LastUserRecord()
		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}

func TestUserLoginLogout(t *testing.T) {
	Convey("TestUserLogin", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		So(db.UserLogin(u), ShouldBeNil)
		So(m.WasCalled(`UserLogin`), ShouldBeTrue)
		rtnUser := m.LastUserRecord()
		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.UserAuthenticated(u), ShouldBeNil)
		So(m.WasCalled(`UserAuthenticated`), ShouldBeTrue)
		rtnUser = m.LastUserRecord()
		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.UserLogout(u), ShouldBeNil)
		So(m.WasCalled(`UserLogout`), ShouldBeTrue)
		rtnUser = m.LastUserRecord()
		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}

func TestFetchUserByGuid(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`guid`, u)
		rtnUser, err := db.FetchUserByGuid(`guid`)
		So(err, ShouldBeNil)
		So(m.WasCalled(`FetchUserByGuid`), ShouldBeTrue)

		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
func TestFetchUserByGuidWithError(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`guid`, u)
		m.ForCallReturnError(`FetchUserByGuid`, errors.New("GUID ERROR"))
		rtnUser, err := db.FetchUserByGuid(`guid`)
		So(m.WasCalled(`FetchUserByGuid`), ShouldBeTrue)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `GUID ERROR`)
		So(rtnUser, ShouldBeNil)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
func TestFetchUserByToken(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`token`, u)
		rtnUser, err := db.FetchUserByToken(`token`)
		So(err, ShouldBeNil)
		So(m.WasCalled(`FetchUserByToken`), ShouldBeTrue)

		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
func TestWithErrorFetchUserByToken(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`token`, u)
		m.ForCallReturnError(`FetchUserByToken`, errors.New("TOKEN ERROR"))
		rtnUser, err := db.FetchUserByToken(`token`)
		So(m.WasCalled(`FetchUserByToken`), ShouldBeTrue)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `TOKEN ERROR`)
		So(rtnUser, ShouldBeNil)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
func TestFetchUserByEmail(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`email`, u)
		rtnUser, err := db.FetchUserByEmail(`email`)
		So(err, ShouldBeNil)
		So(m.WasCalled(`FetchUserByEmail`), ShouldBeTrue)

		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
func TestWithErrorFetchUserByEmail(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`email`, u)
		m.ForCallReturnError(`FetchUserByEmail`, errors.New("EMAIL ERROR"))
		rtnUser, err := db.FetchUserByEmail(`email`)
		So(m.WasCalled(`FetchUserByEmail`), ShouldBeTrue)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `EMAIL ERROR`)
		So(rtnUser, ShouldBeNil)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
func TestFetchUserByLogin(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`login`, u)
		rtnUser, err := db.FetchUserByLogin(`login`)
		So(err, ShouldBeNil)
		So(m.WasCalled(`FetchUserByLogin`), ShouldBeTrue)

		So(rtnUser.Domain, ShouldEqual, u.Domain)
		So(rtnUser.FullName, ShouldEqual, u.FullName)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
func TestWithErrorFetchUserByLogin(t *testing.T) {
	Convey("Open the connection", t, func() {
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		u := record.NewTestUser()
		m.ForLookupByTypeReturn(`login`, u)
		m.ForCallReturnError(`FetchUserByLogin`, errors.New("LOGIN ERROR"))
		rtnUser, err := db.FetchUserByLogin(`login`)
		So(m.WasCalled(`FetchUserByLogin`), ShouldBeTrue)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `LOGIN ERROR`)
		So(rtnUser, ShouldBeNil)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}

func TestWithNoRecordFetchUserByLogin(t *testing.T) {
	Convey("Open the connection", t, func() {
		storage.ResetRegister()
		RegisterMockStore()
		db, err := storage.Open(STORAGE_IDENTITY, "test")

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)
		m, _ := db.GetStorageConnector().(*MockConn)
		So(m.WasCalled(`open`), ShouldBeTrue)

		rtnUser, err := db.FetchUserByLogin(`login`)
		So(m.WasCalled(`FetchUserByLogin`), ShouldBeTrue)

		So(err, ShouldBeNil)
		So(rtnUser, ShouldBeNil)

		So(db.Close(), ShouldBeNil)
		So(m.WasCalled(`close`), ShouldBeTrue)
		m.Reset()
	})
}
