package memory

import (
	"testing"
	"fmt"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/record"
)


func TestStart( t *testing.T ){
	fmt.Println("OK")
}

func TestRegister( t * testing.T ){
	drive := storage.GetDriver()

	user := record.NewTestUser()
	user.SetDomain("Register")
	user.SetToken("TestToken")
	user.SetName( "Just a test name")
	drive.RegisterUser( user )		// Register new user

	user2, err := drive.FetchUserByGuid(  user.GetGuid() )
	if err != nil {
		t.Errorf( "Could not fetch record by GUID: %s" , err )
	}
	if user2.GetToken() != "TestToken" {
		t.Errorf("Token returned is invalid %s" , user2.GetToken())
	}

	user3,err := drive.FetchUserByToken( "TestToken")
	if err != nil {
		t.Errorf( "Could not fetch record by TestToken: %s" , err )
	}
	if user2 == user3  {
		t.Errorf("Same user pointer returned")
	}
	if user2.GetDomain() != user3.GetDomain() ||
			user2.GetToken() != user3.GetToken() ||
			user2.GetName() != user3.GetName() {
		t.Errorf( "User2 doesn't look like User3")
	}
	fmt.Println( user3.String() )
}
