// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
// This requires the optional package bcrypt

package bcrypt

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/cgentry/gus/encryption"
)

type PwdBcrypt struct {}

// The following string should not be changed once you use it.
var internalSalt string

var cost int


func init(){
	encryption.Register( "bcrypt" , &PwdBcrypt{} )
	internalSalt = "}o2P@56ha*6T321h£HcQXleH~$JKR1.t6jwqay%van6e9CSo^gtfyUeQp{2h&gV,KoQi9ysC"
	cost = 7
}
// EncryptPassword will encrypt the password using the magic number within the record.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The magic number will never alter for the life of the record
func (t *PwdBcrypt) EncryptPassword(pwd , salt string ) string {

	pass1 , _ := bcrypt.GenerateFromPassword( []byte(pwd + internalSalt) , cost )
	return string( pass1 )
}

func (t *PwdBcrypt) SetInternalSalt( salt string ){
	internalSalt = salt
}

func ( t * PwdBcrypt ) ComparePasswords( hashedPassword , password , salt string ) bool {
	err := bcrypt.CompareHashAndPassword( []byte(hashedPassword) , []byte( password + internalSalt)  )
	return err == nil
}

func ( t * PwdBcrypt ) SetCost ( costValue int ){
	if costValue > 0 {
		cost = costValue
	}
}
