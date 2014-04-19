package gus

import (
	"testing"
	//"fmt"
)

func TestNewUser( t *testing.T ){
	tuser := NewUser( "Domainxyz")
	if tuser.GetDomain() != "Domainxyz" {
		t.Errorf( "Domain not correct: '%s'" , tuser.Domain)
	}
}

func TestGuid( t * testing.T ){
	tuser := NewUser( "Domain")
	guid := tuser.GetGuid()

	if len( guid ) < 32 {
		t.Errorf("Guid is too short: %s" , guid )
	}
}

func TestAuthenticate( t * testing.T ){
	pwd := "TestingPassvord"
	tuser := NewUser("dom")
	tuser.SetID( 1234)
	tuser.SetPassword( pwd )

	rtn , err := tuser.Login( pwd )

	if err != nil {
		t.Error( err )
	}

	if rtn != USER_OK {
		t.Errorf( "Invalid password %d\n" , rtn )
	}
}
