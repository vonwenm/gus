package record

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

func (user *User) GetID() int {
	return user.Id
}

func (user *User) GetCreatedAt() time.Time {
	return user.CreatedAt
}

func (user *User) GetCreatedAtStr() string {
	return user.CreatedAt.Format(USER_TIME_STR)
}

func (user *User) GetUpdatedAt() time.Time {
	return user.UpdatedAt
}

func (user *User) GetUpdatedAtStr() string {
	return user.UpdatedAt.Format(USER_TIME_STR)
}

func (user *User) GetLastAuthAt() time.Time {
	return user.LastAuthAt
}

func (user *User) GetLastAuthAtStr() string {
	return user.LastAuthAt.Format(USER_TIME_STR)
}

func (user *User) GetDeletedAt() time.Time {
	return user.DeletedAt
}

func (user *User) GetDeletedAtStr() string {
	return user.DeletedAt.Format(USER_TIME_STR)
}

func (user *User) GetLastFailedAt() time.Time {
	return user.LastFailedAt
}

func (user *User) GetLastFailedAtStr() string {
	return user.LastFailedAt.Format(USER_TIME_STR)
}

func (user *User) GetFailCount() int {
	return user.FailCount
}

func (user *User) GetFailCountStr() string {
	return strconv.Itoa(user.FailCount)
}

func (user *User) GetMaxSessionAt() time.Time {
	return user.MaxSessionAt
}

func (user *User) GetTimeoutAt() time.Time {
	return user.TimeoutAt
}

func (user *User) GetTimeoutStr() string {
	return user.TimeoutAt.Format(USER_TIME_STR)
}

func (user *User) GetMaxSessionAtStr() string {
	return user.MaxSessionAt.Format(USER_TIME_STR)
}

func (user *User) GetFullName() string {
	return user.FullName
}

func (user *User) GetName() string {
	return user.FullName
}

func (user *User) GetLoginName() string {
	return user.LoginName
}

// GetGuid will return the unique guid for this user
func (user *User) GetGuid() string {
	if user.Guid == "" {
		guid := md5.New()
		guid.Write([]byte(user.Domain))                  // Add in the user's domain
		guid.Write([]byte(user.GetCreatedAt().String())) // Add in the creation string
		guid.Write([]byte(user.GetSalt()))               // And the user's magic (unique) number
		out := guid.Sum(nil)
		user.Guid = fmt.Sprintf("%x-%x-%x-%x-%x", out[0:4], out[4:6], out[6:8], out[8:10], out[10:])
	}
	return user.Guid
}

// GetDomain will get the domain name for a user
func (user *User) GetDomain() string {
	return user.Domain
}

// GetSalt will get the special account-specific magic number.
// Normally used for salting various other functions, like password
func (user *User) GetSalt() string {
	return user.Salt
}

func (user *User) GetPassword() string {
	return user.Password
}

func (user *User) GetLoginAt() time.Time {
	return user.LoginAt
}

func (user *User) GetLoginAtStr() string {
	return user.LoginAt.Format(USER_TIME_STR)
}

func (user *User) GetLogoutAt() time.Time {
	return user.LogoutAt
}

func (user *User) GetLogoutAtStr() string {
	return user.LogoutAt.Format(USER_TIME_STR)
}

func (user *User) GetEmail() string {
	return user.Email
}

func (user *User) GetToken() string {
	return user.Token
}

// GetToken will check the status of the user and return the token from the record
func (user *User) GetTokenWithExpiration() (string, UserStatusCode) {
	rtn := user.CheckExpirationDates()
	return user.Token, rtn
}
