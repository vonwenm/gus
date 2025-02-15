package response

import (
	"github.com/cgentry/gus/record/stamp"
	"time"
)

type Register struct {
	stamp.Timestamp
	Login   string
	Token   string
	Expires time.Time
}

func NewRegister() *Register {
	rtn := &Register{}
	rtn.SetStamp(time.Now())
	return rtn
}
