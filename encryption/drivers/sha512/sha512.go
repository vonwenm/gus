package sha512

import (
	"crypto/sha512"
	"encoding/base64"
	"github.com/cgentry/gus/encryption"
	//"time"
)

type PwdSha512 struct {}

// The following string should not be changed once you use it.
var pwd_magic string


func init(){
	encryption.Register( "sha512" , &PwdSha512{} )
	pwd_magic = "}o2P@56ha*6T321h£HcQXleH~$JKR1.t6jwqay%van6e9CSo^gtfyUeQp{2h&gV,KoQi9ysC"
}
// EncryptPassword will encrypt the password using the magic number within the record.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The magic number will never alter for the life of the record
func (t *PwdSha512) EncryptPassword(pwd , salt string ) string {

	var pass1 []byte = []byte("")
	crypt := sha512.New()
	for i := 0; i < 4; i++ {
		crypt.Write([]byte(pass1))
		crypt.Write([]byte(salt))
		crypt.Write([]byte(pwd))
		crypt.Write([]byte(pwd_magic))
		pass1 = crypt.Sum(nil)
		crypt.Reset()
	}

	return base64.StdEncoding.EncodeToString( pass1 )
}

func (t *PwdSha512) SetInternalSalt( salt string ){
	pwd_magic = salt
}

