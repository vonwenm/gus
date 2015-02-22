package stamp

import (
	"time"
)

var unixTimeZero = time.Unix(0, 0)

type Timestamper interface {
	Check() error
	SetStamp(time.Time)
	GetStamp() time.Time
	Age() int
	Window(int) int
}

// A timestamp is a simple Golang time structure
type Timestamp struct {
	Stamp time.Time
}

// Return the current time value
func (t *Timestamp) GetStamp() time.Time {
	return t.Stamp
}

// Set the time stamp to the time passed to the function.
func (t *Timestamp) SetStamp(when time.Time) {
		t.Stamp = when
	return
}

// Return how old, in seconds, this timestamp is
func (t *Timestamp) Age() int {

	diff := t.Stamp.Sub(time.Now())
	secdiff := int(diff.Seconds())
	return secdiff
}

// Determine if the age of this record is within a window of time. Return zero if
// it is, or the age difference if it isn't.
func (t *Timestamp) Window(winSeconds int) int {
	age := t.Age()

	if age > winSeconds {
		return age
	}
	if age < -1*winSeconds {
		return age
	}

	return 0
}

// Time is set when it is non-zero and it isn't equal to the Unix 'epoch'
func (t *Timestamp) IsTimeSet() bool {
	if !t.GetStamp().IsZero() {
		if !t.GetStamp().Equal(unixTimeZero) {
			return true
		}
	}
	return false
}

// Return a new, initialised Timestamp structure. The time is set to 'now'.
func New() *Timestamp {
	return &Timestamp{Stamp: time.Now()}
}
