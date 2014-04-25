package record

import (
	"strings"
)
func (user * User) CheckNewPassword(newPassword string) int {
	if len(newPassword) < 6 {
		return USER_PASSWD_TOO_SHORT
	}
	lpwd := strings.ToLower(newPassword)
	if strings.Contains(lpwd,user.FullName) ||
		strings.ToLower(user.Email) == lpwd ||
		strings.Contains(lpwd , "password")  {
		return USER_PASSWORD_TOO_SIMPLE
	}
	return USER_OK
}
