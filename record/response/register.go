package response

import (
	"time"
	"github.com/cgentry/gus/record/stamp"
)

type Register struct {
	stamp.Timestamp
	Login   string
	Token   string
	Expires time.Time
}

func NewRegister() *Register {
	rtn := &Register{}
	rtn.SetStamp( time.Now() )
	return rtn
}
