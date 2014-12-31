package record

// This is a 'holding' routine for better password verification logic. Right now it does
// very little.
//

import (
	"strings"
)

const (
	PASSWORD_MINIMUM_LENGTH = 6
)

func (user *User) CheckNewPassword(newPassword string) UserStatusCode {
	if len(newPassword) < PASSWORD_MINIMUM_LENGTH {
		return USER_PASSWD_TOO_SHORT
	}

	lpwd := strings.ToLower(newPassword)

	if lpwd == "password" {
		return USER_PASSWORD_TOO_SIMPLE
	}
	return USER_OK
}
