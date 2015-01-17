package record

import (
	"strconv"
)

func (user *User) GetID() int {
	return user.Id
}

func (user *User) GetCreatedAtStr() string {
	return user.CreatedAt.Format(USER_TIME_STR)
}

func (user *User) GetUpdatedAtStr() string {
	return user.UpdatedAt.Format(USER_TIME_STR)
}

func (user *User) GetLastAuthAtStr() string {
	return user.LastAuthAt.Format(USER_TIME_STR)
}

func (user *User) GetDeletedAtStr() string {
	return user.DeletedAt.Format(USER_TIME_STR)
}

func (user *User) GetLastFailedAtStr() string {
	return user.LastFailedAt.Format(USER_TIME_STR)
}

func (user *User) GetFailCountStr() string {
	return strconv.Itoa(user.FailCount)
}

func (user *User) GetTimeoutStr() string {
	return user.TimeoutAt.Format(USER_TIME_STR)
}

func (user *User) GetMaxSessionAtStr() string {
	return user.MaxSessionAt.Format(USER_TIME_STR)
}

func (user *User) GetLoginAtStr() string {
	return user.LoginAt.Format(USER_TIME_STR)
}

func (user *User) GetLogoutAtStr() string {
	return user.LogoutAt.Format(USER_TIME_STR)
}
