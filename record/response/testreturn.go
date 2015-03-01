package response

import (
	"github.com/cgentry/gus/record/stamp"
	"time"
)

type Ack struct {
	stamp.Timestamp
	Request string
}

func NewAck(op string) *Ack {
	rtn := &Ack{}
	rtn.SetStamp(time.Now())
	rtn.Request = op
	return rtn
}
