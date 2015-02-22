package tenant

// This is a 'holding' routine for better password verification logic. Right now it does
// very little.
//

import (
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/configure"
	"strings"
)

func CheckNewPassword(newPassword string) error {
	if len(newPassword) < configure.PASSWORD_MINIMUM_LENGTH {
		return ErrPasswordTooShort
	}

	lpwd := strings.ToLower(newPassword)

	if lpwd == "password" {
		return ErrPasswordTooSimple
	}
	return nil
}
