package memory

import (
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/service"
	"strconv"
	"database/sql"
	"time"
	//"fmt"
	"errors"
)

func (t *StorageMem) RegisterUser(user *record.User) error {

	stmt, err := t.GetRegisterSql()
	if err == nil {
		_, err = stmt.Exec(
			user.GetGuid(), user.GetFullName(), user.GetEmail(),
			user.GetDomain(), user.GetLoginName(), user.GetPassword(),
			user.GetToken(), user.GetSalt(), strconv.FormatBool(user.IsActive),
			strconv.FormatBool(user.IsLoggedIn), user.GetLoginAtStr(), user.GetLogoutAtStr(),
			user.GetLastAuthAtStr(), user.GetLastFailedAtStr(), user.GetFailCountStr(),
			user.GetMaxSessionAtStr(), user.GetTimeoutStr(), user.GetCreatedAtStr(),
			user.GetUpdatedAtStr(), user.GetDeletedAtStr() , strconv.FormatBool(user.IsSystem ))

	}
	if err == nil {                                // Some error occured...see if there is a duplicate
		return nil
	}
	code, err2 := t.CheckUserExists(user)    // Check to see if user exists

	if code != service.CODE_USER_DOESNT_EXIST {    // Is not an invalid gUID
		return err2
	}
	return err
}

func (t * StorageMem) CheckUserExists(user *record.User) ( int , error ) {
	stmt, err := t.GetRegisterChecksSql()
	if err == nil {
		var guid , domain, email, login sql.NullString
		row := stmt.QueryRow(
			user.GetGuid(),
			user.GetDomain(),
			user.GetEmail(),
			user.GetLoginName())

		err := row.Scan(&guid, &domain, &email, &login)

		if err == nil {
			if guid.Valid && guid.String == user.GetGuid() {
				return service.CODE_DUPLICATE_KEY, errors.New("User GUID exists")
			}
			if email.Valid && email.String == user.GetEmail() {
				return service.CODE_DUPLICATE_EMAIL, errors.New("Email already in use")
			}
			if login.Valid && login.String == user.GetLoginName() {
				return service.CODE_DUPLICATE_LOGIN_NAME, errors.New("Login name already exists")
			}
		}
		return service.CODE_USER_DOESNT_EXIST, err
	}

	return service.CODE_INTERNAL_ERROR, err
}

func StrToTime(t string) time.Time {
	if val, err := time.Parse(record.USER_TIME_STR, t); err == nil {
		return val
	}

	return time.Unix(0, 0)

}

func StrToBool(t string) bool {
	if val, err := strconv.ParseBool(t); err == nil {
		return val
	}
	return false
}

func StrToInt(t string) int {
	if val, err := strconv.ParseInt(t, 10, 32); err == nil {
		if val >= 0 {
			return int(val)
		}
	}
	return 0
}
