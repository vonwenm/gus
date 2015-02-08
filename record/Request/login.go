package request

import (
	"github.com/cgentry/gus/ecode"
	"strings"
)

type Login struct {
	Login    string
	Password string
}

func NewLogin() *Login {
	return &Login{}
}

func ( r * Login ) Check() error {
	r.Login = strings.TrimSpace(r.Login)
	r.Password = strings.TrimSpace(r.Password)
	if r.Login == "" {
		return ecode.ErrMissingLogin
	}
	if r.Password == "" {
		return ecode.ErrMissingPassword
	}
	return nil
}

func ( r * Login ) GetLogin() string { return r.Login ; }
func ( r * Login ) GetPassword() string { return r.Password; }
