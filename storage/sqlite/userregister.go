package sqlite

import (
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/service"
	"strconv"
	"database/sql"
	//"fmt"
	"errors"
	"strings"
)

// Register a new user with the system. If the record cannot be inserted, see what kind of error there is
func (t *StorageMem) RegisterUser(user *record.User) error {
	return registerUser( t.db , user )
}

// This is the local function for testing purposes.
func registerUser( db *sql.DB , user *record.User) error {

	sql := `INSERT INTO User (` + DB_FIELD_LIST_ALL + `)
		    VALUES (` + strings.Repeat( "?, ",  DB_FIELD_COUNT_ALL - 1) + `? )`
	stmt, err := db.Prepare(sql)

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
	code, err2 := checkUserExists(db , user)    // Check to see if user exists

	if code != service.CODE_USER_DOESNT_EXIST {    // Is not an invalid gUID
		return err2
	}
	return err
}

func checkUserExists( db *sql.DB , user *record.User) ( int , error ) {
	var guid , domain, email, login sql.NullString

	sql := `SELECT Guid, Domain , Email, LoginName
			FROM User
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
				return service.CODE_DUPLICATE_KEY, errors.New("User GUID exists")
			}
			if email.Valid && email.String == user.GetEmail() {
				return service.CODE_DUPLICATE_EMAIL, errors.New("Email already registered")
			}
			if login.Valid && login.String == user.GetLoginName() {
				return service.CODE_DUPLICATE_LOGIN_NAME, errors.New("Login name already exists")
			}
		}
		return service.CODE_USER_DOESNT_EXIST, err
	}

	return service.CODE_INTERNAL_ERROR, err
}
