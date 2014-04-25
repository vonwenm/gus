package plaintext

import (
	"testing"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/encryption"

)

func TestGenerate( t *testing.T ){

	user := gus.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword(  "hello" , user.GetSalt() )
	if pwd != "hello;" + user.GetSalt() {
		t.Errorf("Passwords don't match encrypted")
	}
}

func TestRepeatable( t *testing.T ){
	user := gus.NewTestUser()
	pwd  := encryption.GetDriver().EncryptPassword( "123456" , user.GetSalt() )
	pwd2 := encryption.GetDriver().EncryptPassword( "123456" , user.GetSalt() )
	if( pwd != pwd2 ){
		t.Errorf( "Passwords didn't match: '%s' and '%s'" , pwd , pwd2 )
	}

}

func TestIsLongEnough( t *testing.T ){
	user := gus.NewTestUser()
	pwd := encryption.GetDriver().EncryptPassword(  "hello" , user.GetSalt() )
	pwdLen := len( pwd )
	sbLen := len ( "hello;" + user.GetSalt() )
	if pwdLen != sbLen {
		t.Errorf("PWD isn't long enough %d" , pwdLen )
	}
}

func TestSimilarUserDifferntPwd( t *testing.T ){
	user := gus.NewTestUser()
	pwd  := encryption.GetDriver().EncryptPassword( "123456" , user.GetSalt() )
	user2 := gus.NewTestUser()
	pwd2 := encryption.GetDriver().EncryptPassword( "123456" , user2.GetSalt() )
	if( pwd == pwd2 ){
		t.Errorf( "Passwords for different users should not match: '%s' and '%s'" , pwd , pwd2 )
	}
}

func TestAfterChangingSalt( t *testing.T ){
	user := gus.NewTestUser()
	pwd  := encryption.GetDriver().EncryptPassword( "123456" , user.GetSalt() )
	encryption.GetDriver().SetInternalSalt( "hello - this should screw up password" )
	pwd2 := encryption.GetDriver().EncryptPassword( "123456" , user.GetSalt() )

	if( pwd == pwd2 ){
		t.Errorf( "Passwords with different salts should not match: '%s' and '%s'" , pwd , pwd2 )
	}
}
