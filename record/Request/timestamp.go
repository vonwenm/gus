package request

import (
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/record/stamp"
)

type Test struct {
	*stamp.Timestamp
}

func NewTest() *Test {
	r := &Test{}
	r.Timestamp = stamp.New()
	return r
}

func (r *Test) Check() error {

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
