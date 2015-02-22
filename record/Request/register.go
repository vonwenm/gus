package request

import (
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/record/stamp"
	"strings"
)

type Register struct {
	*stamp.Timestamp
	Login    string
	Name     string
	Email    string
	Password string
}

func NewRegister() *Register {
	r := &Register{}
	r.Timestamp = stamp.New()
	return r
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
	if !r.IsTimeSet() {
		return ecode.ErrRequestNoTimestamp
	}
	// Note: stale time is always 2 minutes old. You can check for earlier times...
	window := r.Window(configure.TIMESTAMP_EXPIRATION)
	if window != 0 {
		if window > 0 {
			return ecode.ErrRequestFuture
		}
		if window < 0 {
			return ecode.ErrRequestExpired
		}
	}
	return nil
}
