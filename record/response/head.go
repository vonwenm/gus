package response

import (
	"errors"
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

type Head struct {
	Code      int
	Message   string
	Timestamp time.Time
	Sequence  int
	*Signature
}

func NewHead() Head {
	h := Head{}
	h.Signature = new(Signature)
	h.Signature.SetSignature("")
	h.Timestamp = time.Now()
	return h
}

func (h Head) Check() error {
	if h.Timestamp.IsZero() {
		return errors.New("Head: No timestamp set")
	}
	// Note: stale time is always 2 minutes old. You can check for earlier times...
	diff := h.Timestamp.Sub(time.Now())
	if diff != 0 {
		if diff > 0 {
			if diff > TIMESTAMP_EXPIRATION {
				return errors.New("Head: Request in the future")
			}
		} else {
			if diff < -1*TIMESTAMP_EXPIRATION {
				return errors.New("Head: Request expired")
			}
		}
	}
	return nil
}

func (h Head) IsTimeSet() bool {
	return !h.Timestamp.IsZero()
}
