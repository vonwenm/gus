package response

import (
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/stamp"
	"github.com/cgentry/gus/record/signature"
)

const (
	TIMESTAMP_EXPIRATION = 2 * 60
)


type Head struct {
	*stamp.Timestamp
	Code      int
	Message   string
	Id		  string
	Sequence  int
	*signature.Signature
}

func NewHead() Head {
	h := new(Head)
	h.Timestamp = stamp.New()
	h.Signature = signature.New()
	h.Signature.SetSignature("")
	return *h
}

func (h Head) Check() error {
	if !h.IsTimeSet() {
		return ecode.ErrHeadNoTimestamp
	}

	// Note: stale time is always 2 minutes old. You can check for earlier times...
	window := h.Window(TIMESTAMP_EXPIRATION)
	if window != 0 {
		if window > 0 {
			return ecode.ErrHeadFuture
		}
		if window < 0 {
			return ecode.ErrHeadExpired
		}
	}
	return nil
}
