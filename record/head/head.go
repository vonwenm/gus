package head

import (
	"fmt"
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/record/signature"
	"github.com/cgentry/gus/record/stamp"
	"time"
)

// HeaderInterface defines the minimum set of calls required for a header
//
type HeaderInterface interface {
	IsTimeSet() bool
	GetSignature() []byte
	SetSignature([]byte)
	Check() error
}

// Head implements the record.HeaderInterface
type Head struct {
	Code     int
	Domain   string
	Id       string
	Sequence int
	BodyType string
	*stamp.Timestamp
	*signature.Signature
}

// New will create a new header and fill in the basic information
func New() *Head {
	h := new(Head)
	h.Timestamp = stamp.New()
	h.SetStamp( time.Now() )
	h.Signature = signature.New()
	return h
}

// Check to see if the package has the required information
func (h *Head) Check() error {

	if h.Domain == "" {
		return ecode.ErrHeadNoDomain
	}
	if h.Id == "" {
		return ecode.ErrHeadNoId
	}
	if !h.IsTimeSet() {
		return ecode.ErrHeadNoTimestamp
	}

	// Note: stale time is always 2 minutes old. You can check for earlier times...
	window := h.Window(configure.TIMESTAMP_EXPIRATION)
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

// Convert the Head value to a string
func (h *Head) String() string {
	return fmt.Sprintf("Domain: '%s', Id: '%s', Time: '%v', Signature: '%s', Set? %v",
		h.Domain, h.Id, h.Timestamp, string(h.Signature.GetSignature()), h.IsTimeSet())
}
