package request

import (
	"github.com/cgentry/gus/ecode"
	"time"
)

const (
	TIMESTAMP_EXPIRATION = 2 * time.Minute
)

type Signature struct {
	sum string
}

func (s *Signature) SetSignature(newSum string) { s.sum = newSum }
func (s *Signature) GetSignature() string       { return s.sum }

// Head implements the record.HeaderInterface
type Head struct {
	Domain    string
	Id        string
	Timestamp time.Time
	Sequence  int
	*Signature
}

func NewHead() Head {
	h := new(Head)
	h.Timestamp = time.Now()
	h.Signature = new(Signature)
	h.Signature.SetSignature("")
	return *h
}

func (h Head) Check() error {

	if h.Domain == "" {
		return ecode.ErrHeadNoDomain
	}
	if h.Id == "" {
		return ecode.ErrHeadNoId
	}
	if h.Timestamp.IsZero() {
		return ecode.ErrHeadNoTimestamp
	}

	// Note: stale time is always 2 minutes old. You can check for earlier times...
	diff := h.Timestamp.Sub(time.Now())
	if diff != 0 {
		if diff > 0 {
			if diff > TIMESTAMP_EXPIRATION {
				return ecode.ErrHeadFuture
			}
		} else {
			if diff < -1*TIMESTAMP_EXPIRATION {
				return ecode.ErrHeadExpired
			}
		}
	}
	return nil
}

func (h Head) IsTimeSet() bool {
	return !h.Timestamp.IsZero()
}
