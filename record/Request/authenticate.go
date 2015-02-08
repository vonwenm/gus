package request

import (
	"github.com/cgentry/gus/ecode"
	"strings"
)

type Authenticate struct {
	Token string
}

func NewAuthenticate() *Authenticate {
	return &Authenticate{}
}


func ( r * Authenticate ) Check() error {
	r.Token = strings.TrimSpace(r.Token)
	if r.Token == "" {
		return ecode.ErrMissingToken
	}
	return nil
}
