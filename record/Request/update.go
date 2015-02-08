package request

import (
	"github.com/cgentry/gus/ecode"
	"strings"
)

type Update struct {
	Token       string
	Login       string
	Name        string
	Email       string
	NewPassword string
	OldPassword string
}

func NewUpdate() *Update {
	return &Update{}
}

func (r *Update) Check() error {
	r.Login = strings.TrimSpace(r.Login)
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.NewPassword = strings.TrimSpace(r.NewPassword)
	r.OldPassword = strings.TrimSpace(r.OldPassword)

	if r.NewPassword == "" {
		return ecode.ErrMissingPasswordNew
	}
	if r.OldPassword == "" {
		return ecode.ErrMissingPassword
	}
	if r.NewPassword == r.OldPassword {
		return ecode.ErrMatchingPassword
	}
	if len( r.OldPassword ) < 6 {
		return ecode.ErrPasswordTooShort
	}

	return nil
}
