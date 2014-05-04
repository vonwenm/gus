// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package memory

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"os"
)

const STORE_LOCAL="/tmp/store.tst"


func TestSetup( t *testing.T ){
	if STORE_LOCAL != ":memory:" {
		os.Remove( STORE_LOCAL )
	}
}

func TestRegister(t *testing.T) {

	drive := storage.GetDriver()
	drive.Open( "sqlite3" , STORE_LOCAL)
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
	Convey( "Test fetching data" , t , func(){
		So( err , ShouldBeNil)
		So( user, ShouldNotBeNil )
	})
	Convey( "Altering and saving data" , t , func(){
		data := "This is the time for all sessions to end"
		var headers = storage.HeaderMap{"one": "1", "two": "2"}
		So( err , ShouldBeNil )
		err := drive.SaveSessionData(user, "**TEST**", &data, &headers)
		So( err, ShouldBeNil )
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
	})

}


func TestSaveUserData(t *testing.T) {
	drive := storage.GetDriver()
	user, err := drive.FetchUserByEmail("et@home.com")
	Convey( "Check fetching data ", t , func(){
		So( err, ShouldBeNil )
		So( user, ShouldNotBeNil)
	})
	data := "Keep me around! I want to LIVE"
	var headers = storage.HeaderMap{"one": "1", "two": "2"}
	err = drive.SaveUserData(user, "**TEST**", &data, &headers)

	Convey( "Check saved data" , t , func(){
		So( err , ShouldBeNil)
		newdata, newheaders, err := drive.GetUserData(user, "**TEST**")
		So( err, ShouldBeNil)
		
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
	})



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

