package request

import (
	"github.com/cgentry/gus/ecode"
	"strings"
)

type Register struct {
	Login    string
	Name     string
	Email    string
	Password string
}

func NewRegister() *Register {
	return &Register{}
}

func (r *Register) Check() error {
	r.Login = strings.TrimSpace(r.Login)
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Email == "" {
		return ecode.ErrMissingEmail
	}
	if r.Login == "" {
		return ecode.ErrMissingLogin
	}
	if r.Name == "" {
		return ecode.ErrMissingName
	}
	if r.Password == "" {
		return ecode.ErrMissingPassword
	}
	if len(r.Password) < 6 {
		return ecode.ErrPasswordTooShort
	}

	return nil
}
