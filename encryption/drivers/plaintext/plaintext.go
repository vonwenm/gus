package plaintext

import (
	"github.com/cgentry/gus/encryption"

)

type PwdPlaintext struct {}

// The following string should not be changed once you use it.
var internalSalt string


func init(){
	encryption.Register( "plaintext" , &PwdPlaintext{} )
	internalSalt = ""
}
// EncryptPassword will encrypt the password using the magic number within the record.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The magic number will never alter for the life of the record
func (t *PwdPlaintext) EncryptPassword(pwd , salt string ) string {

	if len( internalSalt ) == 0 {
		return pwd + ";" + salt
	}
	return pwd + ";" + salt + ";" + internalSalt
}

func (t *PwdPlaintext) SetInternalSalt( salt string ){
	internalSalt = salt
}

func ( t * PwdPlaintext ) ComparePasswords( hashedPassword , password , salt string ) bool {
	return hashedPassword == t.EncryptPassword( password , salt )
}

