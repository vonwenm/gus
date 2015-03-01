// Package bcrypt will encrypt passwords using a strong algorithim, bcrypt,
// and is the recommended driver. The setup option passed in can be:
// { Cost: n, Salt: "string"}
// Where cost is how many iterations bcrypt should perform. The higher the
// number the longer it takes to generate and break. The default is 7.
// Salt is an additional string you want to add in make the hash harder
// to guess. If you don't include one, the internal salt will be used.

// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package

package bcrypt

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/cgentry/gus/encryption"
	"log"
)

type PwdBcrypt struct {
	Name  string
	Salt  string
	Cost  int
	Short string
	Long  string
}

const ENCRYPTION_DRIVER_ID = "bcrypt"

// init will register this driver for use.
func init() {
	encryption.Register(New())
}
func (t *PwdBcrypt) Id() string        { return t.Name }
func (t *PwdBcrypt) ShortHelp() string { return t.Short }
func (t *PwdBcrypt) LongHelp() string  { return t.Long }

// New will create a BCRYPT strucutre. The salt is given a static string but
// can be set up on selection from the driver. This must be the same with every
// load or you won't be able to login anymore.
func New() *PwdBcrypt {
	c := &PwdBcrypt{
		Name:  ENCRYPTION_DRIVER_ID,
		Short: "Standard high-quality encryption using BCRYPT methods",
		Long:  const_bcrypt_help_template,
		Cost:  7,
		Salt:  "vniiO5UD0w5GpJkPijwQCT63MuMjyWnyi5TtUWBGInCq84zaFFsSwGm9DK8UyUeQp{2h&gV,KoQi9ysC",
	}
	return c
}

// EncryptPassword will encrypt the password using the magic number within the record.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The magic number will never alter for the life of the record
func (t *PwdBcrypt) EncryptPassword(clearPassword, userSalt string) string {
	saltyPassword := []byte(clearPassword + t.Salt + userSalt + encryption.GetStaticSalt(0))
	pass1, _ := bcrypt.GenerateFromPassword(saltyPassword, t.Cost)
	return string(pass1)
}

// Setup should be called only when the driver has been selected for use.
func (t *PwdBcrypt) Setup(jsonOptions string) encryption.EncryptDriver {
	if jsonOptions != "" {
		opt, err := encryption.UnmarshalOptions(jsonOptions)
		if err != nil {
			log.Printf("Bcrypt: Could not unmarshal '%s' options: ignored.", jsonOptions)
			return t
		}
		if opt.Cost > 0 {
			t.Cost = opt.Cost
		}
		if len(opt.Salt) > 0 {
			t.Salt = opt.Salt
		}
	}
	return t
}

// ComparePasswords must be called with a bcrypt password.
func (t *PwdBcrypt) ComparePasswords(hashedPassword, clearPassword, userSalt string) bool {
	saltyPassword := []byte(clearPassword + t.Salt + userSalt + encryption.GetStaticSalt(0))
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), saltyPassword)
	return err == nil
}

const const_bcrypt_help_template = `
  The bcrypt function is the default password hash algorithm for BSD and many other systems.
  Besides incorporating a salt to protect against rainbow table attacks, bcrypt is an adaptive
  function: over time, the iteration count can be increased to make it slower, so it remains
  resistant to brute-force search attacks even with increasing computation power.

  Options: There are two options that should be passed by JSON strings. They are:
      "Cost" and "Salt". Cost is the number of iterations you want for the function, making
      it more costly to encrypt (which is a good thing). Salt is an additional bit of
      encryption you want added when it is encrypting the password. The salt should
      be a long, random string of any characters. Do not include quotes.

      The cost defaults to '7' and the salt has a long, random string built in. You must
      not change the salt after you have set it or passwords will never match again.

  Option format: {"Cost" : 7, "Salt": "abcd...........xyz" }

`
