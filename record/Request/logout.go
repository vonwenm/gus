package request

import (
	"github.com/cgentry/gus/ecode"
	"strings"
)

type Logout struct {
	Token string
}

func NewLogout() *Logout {
	return &Logout{}
}

func ( r * Logout ) Check() error {
	r.Token = strings.TrimSpace(r.Token)
	if r.Token == "" {
		return ecode.ErrMissingToken
	}
	return nil
}
