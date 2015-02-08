package request

import (
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/stamp"
	"strings"
)

type Authenticate struct {
	*stamp.Timestamp
	Token string
}

func NewAuthenticate() *Authenticate {
	r := &Authenticate{}
	r.Timestamp = stamp.New()
	return r
}

func (r *Authenticate) Check() error {
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
