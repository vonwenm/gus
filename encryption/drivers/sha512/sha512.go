package sha512

import (
	"crypto/sha512"
	"encoding/base64"
	"github.com/cgentry/gus/encryption"

	//"time"
)

type PwdSha512 struct {
	Name  string
	Salt  string
	Cost  int
	Short string
	Long  string
}

const ENCRYPTION_DRIVER_ID = "sha512"

// Create a new SHA512 encryption. The salt is given a static string but
// can be set up on selection from the driver. This must be the same with every
// load or you won't be able to login anymore.
func New() *PwdSha512 {
	c := &PwdSha512{
		Name:  ENCRYPTION_DRIVER_ID,
		Short: "Standard quality encryption using SHA512 methods",
		Long:  const_sha512_help_template,
		Cost:  4,
		Salt:  "9u4K6f6pKmpUqF%Cgo9$c2rJfZEPut//ziRbrda8A2KQctVxWYKrUCX28GDww.t6jwqay%van6e9CSo^gtfyUeQp{2h&gV,KoQi9ysC",
	}
	return c
}

func init() {
	encryption.Register(New())
}

func (t *PwdSha512) Id() string        { return t.Name }
func (t *PwdSha512) ShortHelp() string { return t.Short }
func (t *PwdSha512) LongHelp() string  { return t.Long }

// EncryptPassword will encrypt the password using the user's salt and our salt.
// This will be re-iterated for 'cost' number of times.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The internal salt must never alter for the life of the record)
func (t *PwdSha512) EncryptPassword(clearPassword, userSalt string) string {

	previousPass := []byte("")
	crypt := sha512.New()
	for i := 0; i < t.Cost; i++ {
		crypt.Write([]byte(previousPass))
		crypt.Write([]byte(userSalt))
		crypt.Write([]byte(clearPassword))
		crypt.Write([]byte(t.Salt))
		crypt.Write([]byte(encryption.ENCRYPTION_SALT1))
		previousPass = crypt.Sum(nil)
		crypt.Reset()
	}

	return base64.StdEncoding.EncodeToString(previousPass)
}

// This should be called  when the driver has been selected for use.
func (t *PwdSha512) Setup(jsonOption string) encryption.CryptDriver {

	opt, err := encryption.UnmarshalOptions(jsonOption)
	if err != nil {
		panic(err.Error())
	}

	if opt.Cost > 0 {
		t.Cost = opt.Cost
	}
	if len(opt.Salt) > 0 {
		t.Salt = opt.Salt
	}
	return t
}

func (t *PwdSha512) ComparePasswords(hashedPassword, clearPassword, salt string) bool {
	return hashedPassword == t.EncryptPassword(clearPassword, salt)
}

const const_sha512_help_template = `
  The SHA512 encryption driver is a reasonable hash method that attempts to balance cost
  and speed. It will take a password, the users salt, a system salt and hash them into
  a string.

  If you have a cost value > 0, then this is iterated with the previous results making it
  slightly more difficult to crack should the database be comprimised. The system salt should
  also be a long, random string of characters stored in the configuration file. This can increase
  the security by separating the salts into 2 parts: the internally compiled salt (which can be
  altered and the code recompiled) and the external salt stored in a configuration file.


  Options: There are two options that should be passed by JSON strings. They are:
      "Cost" and "Salt". Cost is the number of iterations you want for the function, making
      it more costly to encrypt (which is a good thing). Salt is an additional bit of
      encryption you want added when it is encrypting the password. The salt should
      be a long, random string of any characters. Do not include quotes.

      The cost defaults to '4' and the salt is a very long, random string coded in. You must
      not change these values after you have selected them or passwords will never match
      again. You should include a salt in your configuration to increase the security.

  Option format: {"Cost" : 7, "Salt": "abcd...........xyz" }
                 { "Salt": "abc...........xyz" }

`
