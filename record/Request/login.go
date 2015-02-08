package request

import (
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/stamp"
	"strings"
)

// The login structure is passed for every login request. It is usually encoded into
// json in the request package. All fields are required and are tested by the
// check() routine.
type Login struct {
	*stamp.Timestamp
	Login    string
	Password string
}

// Create a new login request with the time set to 'now'
func NewLogin() *Login {
	r := &Login{}
	r.Timestamp = stamp.New()
	return r
}

// Check to see that all login elements are correct. Time must be within a
// window in order for it to pass (within 2 minutes of play)
func (r *Login) Check() error {
	r.Login = strings.TrimSpace(r.Login)
	r.Password = strings.TrimSpace(r.Password)
	if r.Login == "" {
		return ecode.ErrMissingLogin
	}
	if r.Password == "" {
		return ecode.ErrMissingPassword
	}
	if !r.IsTimeSet() {
		return ecode.ErrRequestNoTimestamp
	}
	window := r.Window(TIMESTAMP_EXPIRATION)
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

func (r *Login) GetLogin() string    { return r.Login }
func (r *Login) GetPassword() string { return r.Password }
