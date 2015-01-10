// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

// The SAVE routines come in two flavours: the very specific functions and a generic 'save all fields'
// functions. The specific only update the relevant fields using conditions to ensure that records
// are only updated when appropriate (acting as business logic) while the generic routine will save almost
// all fields that can be updated.
// Not all storage routines will be as specific, but it is best to do when ever posible

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"net/http"
	"strconv"
	"time"
)

var cmd_user_login string
var cmd_user_authenticated string
var cmd_user_logout string
var cmd_user_update string

// Save the user's record when they login. This will update the token, status and dates
// Condition: The user must not be logged in and must be active
func (t *SqliteConn) UserLogin(user *record.User) error {

	if !user.IsActive {
		return storage.ErrUserNotActive
	}
	if cmd_user_login == "" {
		cmd_user_login = fmt.Sprintf(`UPDATE %s
			 SET %s = ?,  %s = ?,
			     %s = ?,  %s = ?,
			     %s = ?
           WHERE %s = ? AND %s = ?`,
			record.USER_STORE_NAME,
			FIELD_TOKEN,
			FIELD_ISLOGGEDIN,
			FIELD_LASTAUTH_DT,
			FIELD_LOGIN_DT,
			FIELD_UPDATED_DT,
			FIELD_GUID,
			FIELD_ISACTIVE)

	}
	now := time.Now()
	fmtTime := now.Format(record.USER_TIME_STR)
	result, err := t.db.Exec(cmd_user_login,
		user.GetToken(),
		strconv.FormatBool(true /* IS LOGGED IN */),
		fmtTime,
		fmtTime,
		fmtTime,
		user.GetGuid(),
		strconv.FormatBool(true))
	if err != nil {
		return err
	}
	if numRows, err := result.RowsAffected(); err != nil {
		return err
	} else if numRows == 0 {
		return storage.ErrUserLoggedIn

	}

	return storage.ErrStatusOk

}

// The user has been authenticated; this will only update the authenticated and updated dates
// Condition: This is done by GUID, the record must be active and logged in
func (t *SqliteConn) UserAuthenticated(user *record.User) error {
	if cmd_user_authenticated == "" {
		cmd_user_authenticated = fmt.Sprintf(`UPDATE %s
			 SET %s = ?,  %s = ?
           WHERE %s = ? AND %s = ? AND %s = ?`,
			record.USER_STORE_NAME,
			FIELD_LASTAUTH_DT, FIELD_UPDATED_DT,
			FIELD_GUID, FIELD_ISLOGGEDIN, FIELD_ISACTIVE)

	}
	now := time.Now()
	fmtTime := now.Format(record.USER_TIME_STR)
	_, err := t.db.Exec(cmd_user_authenticated,
		fmtTime, fmtTime,
		user.GetGuid(), strconv.FormatBool(true), strconv.FormatBool(true))
	return err
}

// The user has logged out; this will update the token, status and dates
// Condition: The user must be logged in and must be active
func (t *SqliteConn) UserLogout(user *record.User) error {
	if cmd_user_logout == "" {
		cmd_user_logout = fmt.Sprintf(`UPDATE %s
			 SET %s = ?,  %s = ?,
			     %s = ?
           WHERE %s = ? AND %s = ?`,
			record.USER_STORE_NAME,
			FIELD_ISLOGGEDIN,
			FIELD_LOGOUT_DT,
			FIELD_UPDATED_DT,
			FIELD_GUID,
			FIELD_ISLOGGEDIN)
	}
	now := time.Now()
	fmtTime := now.Format(record.USER_TIME_STR)
	result, err := t.db.Exec(cmd_user_logout,
		strconv.FormatBool(false /* FIELD_ISLOGGEDIN */),
		fmtTime,                  // FIELD_LOGOUT_DT
		fmtTime,                  // FIELD_UPDATED_DT
		user.GetGuid(),           // FIELD_GUID
		strconv.FormatBool(true)) // FIELD_ISLOGGEDIN
	if err != nil {
		return storage.NewStorageFromError(err, http.StatusInternalServerError)
	}
	if numRows, err := result.RowsAffected(); err != nil {
		return storage.NewStorageFromError(err, http.StatusNotFound)
	} else {
		if numRows == 0 {
			return storage.ErrUserNotLoggedIn
		}
	}
	return storage.ErrStatusOk

}

// Save most of the user record. This is used to perform most updates, including
// password. The only fields that will NOT be updated are Domain, Salt, and CreatedAt.
// UpdatedAt will always be set in this routine from the current time, not from the record.
// Higher level routines can use this to completely update portions of a record. This is NOT
// atomic as the read/update routines do not lock records. This shouldn't be a problem for
// most cases.
func (t *SqliteConn) UserUpdate(user *record.User) error {
	if cmd_user_update == "" {
		cmd_user_update = fmt.Sprintf(`UPDATE %s
			 SET %s = ?,  %s = ?,
			     %s = ?,  %s = ?,
			     %s = ?,  %s = ?,
			     %s = ?,  %s = ?,
			     %s = ?,  %s = ?,
			     %s = ?,  %s = ?,
			     %s = ?,  %s = ?,
			     %s = ?,  %s = ?
           WHERE %s = ? `,
			record.USER_STORE_NAME,
			FIELD_FULLNAME,
			FIELD_EMAIL,
			FIELD_LOGINNAME,
			FIELD_PASSWORD,
			FIELD_TOKEN,
			FIELD_LOGIN_DT,
			FIELD_LOGOUT_DT,
			FIELD_LASTAUTH_DT,
			FIELD_LASTFAILED_DT,
			FIELD_FAILCOUNT,
			FIELD_MAX_SESSION_DT,
			FIELD_TIMEOUT_DT,
			FIELD_UPDATED_DT,
			FIELD_DELETED_DT,
			FIELD_ISACTIVE,
			FIELD_ISLOGGEDIN,
			FIELD_GUID,
		)

	}
	now := time.Now()
	fmtTime := now.Format(record.USER_TIME_STR)
	_, err := t.db.Exec(cmd_user_update,
		user.GetFullName(),        /* FIELD_FULLNAME		*/
		user.GetEmail(),           /* FIELD_EMAIL	  		*/
		user.GetLoginName(),       /* FIELD_LOGINNAME 		*/
		user.GetPassword(),        /* FIELD_PASSWORD 		*/
		user.GetToken(),           /* FIELD_TOKEN 			*/
		user.GetLoginAtStr(),      /* FIELD_LOGIN_DT 		*/
		user.GetLogoutAtStr(),     /* FIELD_LOGOUT_DT 		*/
		user.GetLastAuthAtStr(),   /* FIELD_LASTAUTH_DT 	*/
		user.GetLastFailedAtStr(), /* FIELD_LASTFAILED_DT 	*/
		user.GetFailCountStr(),    /* FIELD_FAILCOUNT 		*/
		user.GetMaxSessionAtStr(), /* FIELD_MAX_SESSION_DT */
		user.GetTimeoutStr(),      /* FIELD_TIMEOUT_DT 	*/
		fmtTime,                   /* FIELD_UPDATED_DT 	*/
		user.GetDeletedAtStr(),    /* FIELD_DELETED_DT 	*/
		strconv.FormatBool(user.IsActive),
		strconv.FormatBool(user.IsLoggedIn),
		user.GetGuid()) /* FIELD_GUID */

	return err
}
