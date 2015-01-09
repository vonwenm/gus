package request

import (
	"errors"
	"strings"
)

type Update struct {
	Token    string
	Login    string
	Name     string
	Email    string
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

	if r.NewPassword != "" {
		return errors.New("Password cannot be blank")
	}
	if len(r.NewPassword) < 6 {
		return errors.New("Password length is too short")
	}

	return nil
}
