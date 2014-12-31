// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	//. "github.com/smartystreets/goconvey/convey"
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


