package request

import (
	"errors"
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

	if (r.NewPassword == "" && r.OldPassword != "") || (r.NewPassword != "" && r.OldPassword == "") {
		return errors.New("New and old password must be set")
	}
	if r.NewPassword != "" && r.OldPassword != "" {
		if len(r.NewPassword) < 6 {
			return errors.New("Password length is too short")
		}
	}

	return nil
}
