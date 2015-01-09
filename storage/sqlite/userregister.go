package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"net/http"
	"strconv"
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
func registerUser(db *sql.DB, user *record.User) *storage.StorageError {

	if err := checkUserExists(db, user); err != storage.ErrUserNotFound {
		return err
	}

	if cmd_user_register == "" {
		cmd_user_register = fmt.Sprintf(
			`INSERT INTO %s
			(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
		    VALUES (%s %s)`,
			record.USER_STORE_NAME,
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
			FIELD_ISSYSTEM,
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
			`?`)

	}
	fmt.Println( 60 )
	stmt, err := db.Prepare(cmd_user_register)
	fmt.Println( 62 )
	if err != nil {
		fmt.Println( 64 )
		return storage.NewStorageFromError( err, http.StatusInternalServerError)
	}
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
		fmtTime /*Created_DT */, user.GetDeletedAtStr(),
		fmtTime, /*Updated_DT*/
	)
	if err != nil {
		fmt.Println( 92 )
		return storage.NewStorageFromError( err, http.StatusInternalServerError)
	}

	if count, err := result.RowsAffected(); err != nil {
		fmt.Println( 97 )
		return storage.NewStorageFromError( err, http.StatusInternalServerError)
	} else {
		if count == 0 {
			fmt.Println( 101 )
			return storage.ErrUserNotRegistered
		}
	}
	fmt.Println( 105 )
	return storage.ErrStatusOk
}

func checkUserExists(db *sql.DB, user *record.User) * storage.StorageError {
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
