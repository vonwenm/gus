// This package simply defines the interfaces that any password generator must define

package gus

type Passworder interface {

	EncryptPassword( user * User , password string ) string
}


