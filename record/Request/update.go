package request

import (
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/record/stamp"
	"strings"
)

type Update struct {
	*stamp.Timestamp
	Token       string
	Login       string
	Name        string
	Email       string
	NewPassword string
	OldPassword string
}

func NewUpdate() *Update {
	r := &Update{}
	r.Timestamp = stamp.New()
	return r
}

func (r *Update) Check() error {
	r.Token = strings.TrimSpace(r.Token)
	r.Login = strings.TrimSpace(r.Login)
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.NewPassword = strings.TrimSpace(r.NewPassword)
	r.OldPassword = strings.TrimSpace(r.OldPassword)

	if r.NewPassword == "" {
		return ecode.ErrMissingPasswordNew
	}
	if len(r.NewPassword) < 6 {
		return ecode.ErrPasswordTooShort
	}
	if r.OldPassword == "" {
		return ecode.ErrMissingPassword
	}
	if r.NewPassword == r.OldPassword {
		return ecode.ErrMatchingPassword
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
