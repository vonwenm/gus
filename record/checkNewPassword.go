package record

// This is a 'holding' routine for better password verification logic. Right now it does
// very little.
//

import (
	"strings"
	. "github.com/cgentry/gus/ecode"
)

const (
	PASSWORD_MINIMUM_LENGTH = 6
)

func (user *User) CheckNewPassword(newPassword string) error {
	if len(newPassword) < PASSWORD_MINIMUM_LENGTH {
		return ErrPasswordTooShort
	}

	lpwd := strings.ToLower(newPassword)

	if lpwd == "password" {
		return ErrPasswordTooSimple
	}
	return nil
}
