package request

import (
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/stamp"
	"strings"
)

type Logout struct {
	*stamp.Timestamp
	Token string
}

func NewLogout() *Logout {
	r := &Logout{}
	r.Timestamp = stamp.New()
	return r
}

func (r *Logout) Check() error {
	r.Token = strings.TrimSpace(r.Token)
	if r.Token == "" {
		return ecode.ErrMissingToken
	}
	if !r.IsTimeSet() {
		return ecode.ErrRequestNoTimestamp
	}
	// Note: stale time is always 2 minutes old. You can check for earlier times...
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
