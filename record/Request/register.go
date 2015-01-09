package request

import (
	"errors"
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
	if r.Login == "" {
		return errors.New("Login cannot be blank")
	}
	if r.Name == "" {
		return errors.New("Name cannot be blank")
	}
	if r.Email == "" {
		return errors.New("Email cannot be blank")
	}
	if r.Password == "" {
		return errors.New("Password cannot be blank")
	}
	if len(r.Password) < 6 {
		return errors.New("Password length is too short")
	}

	return nil
}
