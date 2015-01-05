// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package mock

import (
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"strings"
)

const STORAGE_IDENTITY = "mock"

type MockDriver struct{}

type MockConn struct {
	internalError      error
	internalUserRecord *record.User

	lastError error
	lastName  string

	errorReturnForCall map[string]error
	userReturnForType  map[string]*record.User
	callList           map[string]int // Called routines

}

// Register this driver to the main storage driver with a unique name and activate it.
func RegisterMockStore() {
	storage.Register(STORAGE_IDENTITY, &MockDriver{})
}

func NewMockConn() *MockConn {
	store := &MockConn{}
	store.Reset()
	return store
}

func (t MockDriver) Open(connect string) (storage.Conn, error) {
	s := NewMockConn()
	s.routineCalled(`open`)
	return s, nil
}
func (t *MockConn) CreateStore() error {
	t.routineCalled(`CreateStore`)
	return t.lookupError(`CreateStore`)
}

func (t *MockConn) Close() error {
	t.routineCalled(`close`)
	return t.lookupError(`close`)
}

func (t *MockConn) RegisterUser(u *record.User) error {
	t.routineCalled(`RegisterUser`)
	t.internalUserRecord = u
	t.routineCalled(`RegisterUser`)
	return t.lookupError(`RegisterUser`)
}

func (t *MockConn) UserLogin(user *record.User) error {
	t.routineCalled(`UserLogin`)
	t.internalUserRecord = user
	return t.lookupError(`UserLogin`)
}

func (t *MockConn) UserAuthenticated(user *record.User) error {
	t.routineCalled(`UserAuthenticated`)
	t.internalUserRecord = user
	return t.lookupError(`UserAuthenticated`)
}

func (t *MockConn) UserLogout(user *record.User) error {
	t.routineCalled(`UserLogout`)
	t.internalUserRecord = user
	return t.lookupError(`UserLogout`)
}

func (t *MockConn) UserUpdate(user *record.User) error {
	t.routineCalled(`UserUpdate`)
	t.internalUserRecord = user
	return t.lookupError(`UserUpdate`)
}

func (t *MockConn) FetchUserByGuid(guid string) (user *record.User, err error) {
	t.routineCalled(`FetchUserByGuid`)
	err = t.lookupError(`FetchUserByGuid`)
	if err == nil {
		user = t.lookupUser(`guid`)
	}
	return
}

func (t *MockConn) FetchUserByToken(token string) (user *record.User, err error) {
	t.routineCalled(`FetchUserByToken`)
	err = t.lookupError(`FetchUserByToken`)
	if err == nil {
		user = t.lookupUser(`token`)
	}
	return
}

func (t *MockConn) FetchUserByEmail(email string) (user *record.User, err error) {
	t.routineCalled(`FetchUserByEmail`)
	err = t.lookupError(`FetchUserByEmail`)
	if err == nil {
		user = t.lookupUser(`email`)
	}
	return
}

func (t *MockConn) FetchUserByLogin(login string) (user *record.User, err error) {
	t.routineCalled(`FetchUserByLogin`)
	err = t.lookupError(`FetchUserByLogin`)
	if err == nil {
		user = t.lookupUser(`login`)
	}
	return
}

/*
 * Internal routines to handle tables and returns
 */

func (t *MockConn) lookupUser(name string) *record.User {
	if val, ok := t.userReturnForType[strings.ToLower(name)]; ok {
		return val
	}
	return nil
}
func (t *MockConn) lookupError(name string) error {
	if val, ok := t.errorReturnForCall[strings.ToLower(name)]; ok {
		t.lastError = val
		return val
	}
	t.lastError = nil
	return nil
}

func (t *MockConn) routineCalled(name string) {
	t.lastName = strings.ToLower(name)
	if v, ok := t.callList[t.lastName]; ok {
		t.callList[t.lastName] = v + 1
	} else {
		t.callList[t.lastName] = 1
	}
}

// TESTING SIDE

// Set what error to set for a message
func (t *MockConn) ForCallReturnError(call string, e error) {
	t.errorReturnForCall[strings.ToLower(call)] = e
}

func (t *MockConn) ForLookupByTypeReturn(call string, u *record.User) {
	t.userReturnForType[strings.ToLower(call)] = u
}

func (t *MockConn) WasCalled(call string) bool {
	if v, ok := t.callList[strings.ToLower(call)]; ok {
		return v > 0
	}
	return false
}

func (t *MockConn) WasCalledOnce(call string) bool {
	if v, ok := t.callList[strings.ToLower(call)]; ok {
		return v == 1
	}
	return false
}

func (t *MockConn) LastUserRecord() *record.User {
	return t.internalUserRecord
}

func (t *MockConn) Reset() {
	t.internalUserRecord = nil
	t.errorReturnForCall = make(map[string]error)
	t.userReturnForType = make(map[string]*record.User)
	t.callList = make(map[string]int)
}
