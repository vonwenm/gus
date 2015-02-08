package response

import (
	"errors"
	"time"
	"github.com/cgentry/gus/record/stamp"
)

const (
	TIMESTAMP_EXPIRATION = 2 * time.Minute
)

type Signature struct {
	sum string
}

func (s *Signature) SetSignature(newSum string) { s.sum = newSum }
func (s *Signature) GetSignature() string       { return s.sum }

type Head struct {
	*stamp.Timestamp
	Code      int
	Message   string
	Id		  string
	Sequence  int
	*Signature
}

func NewHead() Head {
	h := new(Head)
	h.Timestamp = stamp.New()
	h.Signature = new(Signature)
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
