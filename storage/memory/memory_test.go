// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package memory

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"testing"
)

func TestRegister(t *testing.T) {

	drive := storage.GetDriver()
	drive.Open( "sqlite3" , ":memory:")
	drive.CreateStore()

	user := record.NewTestUser()
	user.SetDomain("Register")
	user.SetToken("TestToken")
	user.SetName("Just a test name")
	user.SetEmail("et@home.com")
	drive.RegisterUser(user) // Register new user

	user2, err := drive.FetchUserByGuid(user.GetGuid())
	if err != nil {
		t.Fatalf("Could not fetch record by GUID '%s' : %s", user.GetGuid(), err)
	}
	if user2.GetToken() != "TestToken" {
		t.Errorf("Token returned is invalid %s", user2.GetToken())
	}

	user3, err := drive.FetchUserByToken("TestToken")
	if err != nil {
		t.Errorf("Could not fetch record by TestToken: %s", err)
	}
	if user2 == user3 {
		t.Errorf("Same user pointer returned")
	}
	if user2.GetDomain() != user3.GetDomain() ||
		user2.GetToken() != user3.GetToken() ||
		user2.GetName() != user3.GetName() {
		t.Errorf("User2 doesn't look like User3")
	}
	fmt.Println(user3.String())
}

func TestDuplicateKey(t *testing.T) {

	drive := storage.GetDriver()
	user := record.NewTestUser()
	user.SetDomain("Register")
	user.SetToken("TestToken")
	user.SetName("Just a test name")
	user.SetEmail("et@home.com")
	err := drive.RegisterUser(user) // Register new user
	if err == nil {
		t.Errorf( "Did not get the error (601) of duplicate mail")
	}
}

func TestSaveSessionData(t *testing.T) {
	drive := storage.GetDriver()
	user, err := drive.FetchUserByEmail("et@home.com")
	if err != nil {
		t.Errorf("Could not fetch by email: %s", err)
	}
	data := "This is the time for all sessions to end"
	var headers = storage.HeaderMap{"one": "1", "two": "2"}
	if err = drive.SaveSessionData(user, "**TEST**", &data, &headers); err != nil {

		t.Errorf("Could not save user: %s", err)
	}

	newdata, newheaders, err := drive.GetSessionData(user, "**TEST**")
	if data != newdata {
		t.Errorf("Data does not match: '%s' != '%s'", data, newdata)
	}
	for key, value := range headers {
		if cval, found := newheaders[key]; found {
			if cval != value {
				t.Errorf("Value return '%s' does not match expected '%s", cval, value)
			}
		} else {
			t.Errorf("The key '%s' (value) not found in return map", key)
		}
	}
}


func TestSaveUserData(t *testing.T) {
	drive := storage.GetDriver()
	user, err := drive.FetchUserByEmail("et@home.com")
	if err != nil {
		t.Errorf("Could not fetch by email: %s", err)
	}
	data := "Keep me around! I want to LIVE"
	var headers = storage.HeaderMap{"one": "1", "two": "2"}
	if err = drive.SaveUserData(user, "**TEST**", &data, &headers); err != nil {

		t.Errorf("Could not save user: %s", err)
	}

	newdata, newheaders, err := drive.GetUserData(user, "**TEST**")
	if err != nil {
		t.Errorf("GetUserData error: %s" , err )
	}
	if data != newdata {
		t.Errorf("Data does not match: '%s' != '%s'", data, newdata)
	}
	for key, value := range headers {
		if cval, found := newheaders[key]; found {
			if cval != value {
				t.Errorf("Value return '%s' does not match expected '%s", cval, value)
			}
		} else {
			t.Errorf("The key '%s' (value) not found in return map", key)
		}
	}
}

func TestDeleteSessionData( t *testing.T ){
	drive := storage.GetDriver()
	user, err := drive.FetchUserByEmail("et@home.com")
	if err != nil {
		t.Errorf("Could not fetch by email: %s", err)
	}
	drive.DeleteSessionData( user , "**TEST**" )
	newdata, newheaders, err := drive.GetSessionData(user, "**TEST**")

	if newdata != "" {
		t.Errorf( "Data should be blank: %s" , newdata )
	}
	if newheaders != nil {
		t.Errorf( "Headers should be emtpy %s" , newheaders)
	}
	if err == nil {
		t.Errorf( "GetSessionData: no error occured!" )
	}
}

func TestDeleteUserData( t *testing.T ){
	drive := storage.GetDriver()
	user, err := drive.FetchUserByEmail("et@home.com")
	if err != nil {
		t.Errorf("Could not fetch by email: %s", err)
	}
	drive.DeleteUserData( user , "**TEST**" )
	newdata, newheaders, err := drive.GetUserData(user, "**TEST**")

	if newdata != "" {
		t.Errorf( "Data should be blank: %s" , newdata )
	}
	if newheaders != nil {
		t.Errorf( "Headers should be emtpy %s" , newheaders)
	}
	if err == nil {
		t.Errorf( "GetUserData: no error occured!" )
	}
}

