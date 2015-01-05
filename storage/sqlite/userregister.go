package sqlite

import (
	"database/sql"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"strconv"
	"fmt"
	"strings"
	"time"
)

var cmd_user_register string
var cmd_user_check string
// Register a new user with the system. If the record cannot be inserted, see what kind of error there is
func (t *SqliteConn) RegisterUser(user *record.User) error {
	return registerUser(t.db, user)
}

// This is the local function for testing purposes.
func registerUser(db *sql.DB, user *record.User) error {

	err := checkUserExists(db, user) // Check to see if user exists

	if err.Error() != storage.ErrUserNotFound.Error() {
		return err
	}

	if cmd_user_register == "" {
		cmd_user_register = fmt.Sprintf(
			`INSERT INTO %s
			(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
		    VALUES (%s %s)`,
			record.USER_STORE_NAME ,

			FIELD_DOMAIN,
			FIELD_EMAIL,
			FIELD_FULLNAME,
			FIELD_GUID,
			FIELD_LOGINNAME,
			FIELD_PASSWORD,
			FIELD_SALT,
			FIELD_TOKEN,

			FIELD_ISACTIVE,
			FIELD_ISLOGGEDIN,
			FIELD_ISSYSTEM ,

			FIELD_FAILCOUNT,

			FIELD_LOGIN_DT,
			FIELD_LOGOUT_DT,
			FIELD_LASTAUTH_DT,
			FIELD_LASTFAILED_DT,
			FIELD_MAX_SESSION_DT,
			FIELD_TIMEOUT_DT,

			FIELD_CREATED_DT,
			FIELD_DELETED_DT,
			FIELD_UPDATED_DT,

			strings.Repeat(`?, `, 20),
			`?` )

	}
	stmt, err := db.Prepare(cmd_user_register)

	if err == nil {
		now := time.Now()
		fmtTime := now.Format(record.USER_TIME_STR)
		result, err := stmt.Exec(
			user.GetDomain(),
			user.GetEmail(),
			user.GetFullName(),
			user.GetGuid(),
			user.GetLoginName(),
			user.GetPassword(),
			user.GetSalt(),
			user.GetToken(),

			strconv.FormatBool(user.IsActive),
			strconv.FormatBool(user.IsLoggedIn),
			strconv.FormatBool(user.IsSystem),

			user.GetFailCountStr(),

			user.GetLoginAtStr(),
			user.GetLogoutAtStr(),
			user.GetLastAuthAtStr(),
			user.GetLastFailedAtStr(),
			user.GetMaxSessionAtStr(),
			user.GetTimeoutStr(),
			fmtTime /*Created_DT */,
			user.GetDeletedAtStr(),
			fmtTime /*Updated_DT*/,

		)

		if err == nil { // Some error occured...see if there is a duplicate
			if count, err := result.RowsAffected(); err != nil {
				return nil
			} else {
				if count == 0 {
					return storage.ErrUserNotRegistered
				}
			}
		}
	}

	return err
}

func checkUserExists(db *sql.DB, user *record.User) error {
	var guid, domain, email, login sql.NullString

	sql := `SELECT Guid, Domain , Email, LoginName
			FROM ` + record.USER_STORE_NAME + `
			WHERE Guid = ? OR ( Domain = ? AND ( Email = ? OR LoginName = ?))`

	stmt, err := db.Prepare(sql)
	if err == nil {

		row := stmt.QueryRow(
			user.GetGuid(),
			user.GetDomain(),
			user.GetEmail(),
			user.GetLoginName())

		err := row.Scan(&guid, &domain, &email, &login)

		if err == nil {
			if guid.Valid && guid.String == user.GetGuid() {
				return storage.ErrDuplicateGuid
			}
			if email.Valid && email.String == user.GetEmail() {
				return storage.ErrDuplicateEmail
			}
			if login.Valid && login.String == user.GetLoginName() {
				return storage.ErrDuplicateLogin
			}
		}
		return storage.ErrUserNotFound
	}

	return storage.ErrInternalDatabase
}
