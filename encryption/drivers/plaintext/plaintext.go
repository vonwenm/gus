package plaintext

import (
	"github.com/cgentry/gus/encryption"
)

const ENCRYPTION_DRIVER_ID = "plaintext"

type PwdPlaintext struct {
	Name  string
	Salt  string
	Short string
	Long  string
}

func New() *PwdPlaintext {
	c := &PwdPlaintext{
		Name:  ENCRYPTION_DRIVER_ID,
		Short: "For testing only! Do not use in production",
		Long:  const_plain_help_template,
		Salt:  "SALT",
	}
	return c
}

func init() {
	encryption.Register(New())
}
func (t *PwdPlaintext) Id() string        { return t.Name }
func (t *PwdPlaintext) ShortHelp() string { return t.Short }
func (t *PwdPlaintext) LongHelp() string  { return t.Long }

// EncryptPassword will encrypt the password using the magic number within the record.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The magic number will never alter for the life of the record
func (t *PwdPlaintext) EncryptPassword(clearPassword, userSalt string) string {

	return clearPassword + ";" + userSalt + ";" + t.Salt + ";Plaintext"
}

// This should be called only when the driver has been selected for use.
func (t *PwdPlaintext) Setup(json string) encryption.CryptDriver {
	opt, err := encryption.UnmarshalOptions(json)
	if err != nil {
		panic(err.Error())
	}

	if len(opt.Salt) > 0 {
		t.Salt = opt.Salt
	}

	return t
}
func (t *PwdPlaintext) ComparePasswords(hashedPassword, password, salt string) bool {
	return hashedPassword == t.EncryptPassword(password, salt)
}

const const_plain_help_template = `
  This does not encrypt passwords and should never be selected for production use. It
  is only to be used by developers and for testing purposes. The format of the password
  output is:
           [user password];[user salt];[driver's salt];Plaintext
  If a user has a salt of 'kjldoeuifnfl203294fkf' and the password is 'BadPassword', with
  defaults it would become:
           BadPassword;kjldoeuifnfl203294fkf;SALT;Plaintext

  Options: There is one option that can be passed in JSON format: "Salt". The default is "SALT".

  Option format: {"Salt": "Salty" }
`
