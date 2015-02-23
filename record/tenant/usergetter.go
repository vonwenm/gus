package tenant

import (
	"github.com/cgentry/gus/record/configure"
	"strconv"
)

func (user *User) GetID() int {
	return user.Id
}

func (user *User) GetCreatedAtStr() string {
	return user.CreatedAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetUpdatedAtStr() string {
	return user.UpdatedAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetLastAuthAtStr() string {
	return user.LastAuthAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetDeletedAtStr() string {
	return user.DeletedAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetLastFailedAtStr() string {
	return user.LastFailedAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetFailCountStr() string {
	return strconv.Itoa(user.FailCount)
}

func (user *User) GetTimeoutStr() string {
	return user.TimeoutAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetMaxSessionAtStr() string {
	return user.MaxSessionAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetLoginAtStr() string {
	return user.LoginAt.Format(configure.USER_TIME_STR)
}

func (user *User) GetLogoutAtStr() string {
	return user.LogoutAt.Format(configure.USER_TIME_STR)
}
