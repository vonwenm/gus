// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/record"
	"strings"
)

const STORAGE_IDENTITY = "mock"

type MockStore struct {
	internalError error
	internalUserRecord  * record.User					// Record we were passed

	lastError error
	lastName  string

	errorReturnForCall map[string]error
	userReturnForType  map[string]*record.User
	callList map[string]int								// Called routines

	
}

// Register this driver to the main storage driver with a unique name and activate it.
func RegisterMockStore() * MockStore {
	store := &MockStore{}

	store.Reset()

	storage.Register(STORAGE_IDENTITY, store)
	storage.SetDriver(STORAGE_IDENTITY)
	return store
}

func (t *MockStore) GetLastError() error {
	return t.lastError
}

func (t *MockStore) WasLastOk() bool {
	return t.lastError == nil
}

func (t *MockStore) Open(name, connect string) error {
	t.routineCalled(`open`)
	return t.lookupError(`open`)
}
func (t *MockStore) CreateStore() error {
	t.routineCalled(`CreateStore`)
	return t.lookupError( `CreateStore`)
}

func (t *MockStore) Close() error {
	t.routineCalled(`close`)
	return t.lookupError( `close`)
}

func ( t *MockStore) RegisterUser(u *record.User) error {
	t.routineCalled(`RegisterUser`)
	t.internalUserRecord = u
	t.routineCalled(`RegisterUser`)
	return t.lookupError( `RegisterUser`)
}


func ( t *MockStore) SaveUserLogin(user *record.User) error  {
	t.routineCalled(`SaveUserLogin`)
	t.internalUserRecord = user
	return t.lookupError( `SaveUserLogin`)
}


func ( t *MockStore) SaveUserAuth(user *record.User) error {
	t.routineCalled(`SaveUserAuth`)
	t.internalUserRecord = user
	return t.lookupError( `SaveUserAuth`)
}


func ( t *MockStore) SaveUserLogoff(user *record.User) error {
	t.routineCalled(`SaveUserLogoff`)
	t.internalUserRecord = user
	return t.lookupError(`SaveUserLogoff`)
}



func ( t *MockStore) FetchUserByGuid(guid string) ( user *record.User,err error){
	t.routineCalled(`FetchUserByGuid`)
	err = t.lookupError(`FetchUserByGuid`)
	if err == nil {
		user = t.lookupUser(`guid`)
	}
	return
}


func ( t *MockStore) FetchUserByToken(token string) ( user *record.User,err error){
	t.routineCalled(`FetchUserByToken`)
	err = t.lookupError(`FetchUserByToken`)
	if err == nil {
		user = t.lookupUser(`token`)
	}
	return
}


func ( t *MockStore) FetchUserByEmail(email string) ( user *record.User,err error){
	t.routineCalled(`FetchUserByEmail`)
	err = t.lookupError(`FetchUserByEmail`)
	if err == nil {
		user = t.lookupUser(`email`)
	}
	return
}


func ( t *MockStore) FetchUserByLogin(login string)( user *record.User,err error){
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

func (t *MockStore) lookupUser( name string ) * record.User {
	if val,ok := t.userReturnForType[ strings.ToLower(name)]; ok {
		return val
	}
	return nil
}
func (t *MockStore) lookupError( name string ) error {
	if val,ok := t.errorReturnForCall[ strings.ToLower(name) ]; ok {
		t.lastError = val
		return val
	}
	t.lastError = nil
	return nil
}

func ( t * MockStore ) routineCalled(name string ){
	t.lastName = strings.ToLower(name)
	if v,ok := t.callList[t.lastName]; ok {
		t.callList[t.lastName] = v+1
	}else{
		t.callList[t.lastName] = 1
	}
}


// TESTING SIDE

// Set what error to set for a message
func ( t * MockStore ) ForCallReturnError( call string , e error ) {
	t.errorReturnForCall[ strings.ToLower( call )] = e
}

func ( t * MockStore ) ForLookupByTypeReturn( call string , u *record.User ){
	t.userReturnForType[ strings.ToLower(call)] = u
}

func ( t * MockStore ) WasCalled( call string ) bool {
	if v,ok := t.callList[strings.ToLower(call)]; ok {
		return v > 0
	}
	return false
}

func ( t * MockStore )WasCalledOnce( call string ) bool {
	if v,ok := t.callList[strings.ToLower(call)]; ok {
		return v == 1
	}
	return false
}

func ( t * MockStore)LastUserRecord() * record.User {
	return t.internalUserRecord
}

func ( t * MockStore ) Reset(){
	t.internalUserRecord = nil
	t.errorReturnForCall = make( map[string]error)
	t.userReturnForType  = make( map[string]*record.User)
	t.callList = make( map[string]int )
}
