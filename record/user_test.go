package record

import (
	"testing"
	_ "github.com/cgentry/gus/encryption/drivers/plaintext"
	. "github.com/smartystreets/goconvey/convey"
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


	Convey( "Authenticatin test" , t , func(){
		rtn , err := tuser.Login( pwd )
		So( err , ShouldBeNil )
		So( rtn, ShouldEqual , USER_OK )
	})
}
