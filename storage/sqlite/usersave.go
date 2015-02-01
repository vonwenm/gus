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
	"database/sql"
	"fmt"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record"
	"net/http"
	"strconv"
	"strings"
)

var cmd_user_update string
var cmd_user_insert string

// Update the database from the user record passed. The only fields that are not updated
// are the CreatedAt and GUID fields. These are only set on an UserInsert call
func (t *SqliteConn) UserUpdate(user *record.User) error {
	if cmd_user_update == "" {
		cmd_user_update = fmt.Sprintf(`UPDATE %s
			 SET %s = ?,
			 	 %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?,
			     %s = ?
           WHERE %s = ? `,
			record.USER_STORE_NAME,

			FIELD_DOMAIN,
			FIELD_EMAIL,
			FIELD_FAILCOUNT,
			FIELD_FULLNAME,
			FIELD_LOGINNAME,
			FIELD_PASSWORD,
			FIELD_SALT,
			FIELD_TOKEN,

			FIELD_ISACTIVE,
			FIELD_ISLOGGEDIN,
			FIELD_ISSYSTEM,

			FIELD_LASTAUTH_DT,
			FIELD_LASTFAILED_DT,
			FIELD_LOGIN_DT,
			FIELD_LOGOUT_DT,
			FIELD_MAX_SESSION_DT,
			FIELD_TIMEOUT_DT,

			FIELD_UPDATED_DT,
			FIELD_DELETED_DT,

			FIELD_GUID,
		)

	}

	_, err := t.db.Exec(cmd_user_update,
		user.Domain,
		user.Email,
		user.GetFailCountStr(),
		user.FullName,
		user.LoginName,
		user.Password,
		user.Salt,
		user.Token,

		strconv.FormatBool(user.IsActive),
		strconv.FormatBool(user.IsLoggedIn),
		strconv.FormatBool(user.IsSystem),

		user.GetLastAuthAtStr(),
		user.GetLastFailedAtStr(),
		user.GetLoginAtStr(),
		user.GetLogoutAtStr(),
		user.GetMaxSessionAtStr(),
		user.GetTimeoutStr(),

		user.GetUpdatedAtStr(),
		user.GetDeletedAtStr(),

		user.Guid) /* FIELD_GUID - KEY*/
	if err != nil {
		return NewGeneralFromError(err, http.StatusInternalServerError)
	}

	return nil
}

// Save most of the user record. This is used to perform general updates, including
// password. The only fields that will NOT be updated are Domain, Salt, and CreatedAt.
// UpdatedAt will always be set in this routine from the current time, not from the record.
// Higher level routines can use this to completely update portions of a record. This is NOT
// atomic as the read/update routines do not lock records. This shouldn't be a problem for
// most cases.
func (t *SqliteConn) UserInsert(user *record.User) error {

	if cmd_user_insert == "" {
		cmd_user_insert = fmt.Sprintf(
			`INSERT INTO %s
			(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
		    VALUES (%s %s)`,
			record.USER_STORE_NAME,

			FIELD_DOMAIN,
			FIELD_EMAIL,
			FIELD_FAILCOUNT,
			FIELD_FULLNAME,
			FIELD_GUID,
			FIELD_LOGINNAME,
			FIELD_PASSWORD,
			FIELD_SALT,
			FIELD_TOKEN,

			FIELD_ISACTIVE,
			FIELD_ISLOGGEDIN,
			FIELD_ISSYSTEM,

			FIELD_LASTAUTH_DT,
			FIELD_LASTFAILED_DT,
			FIELD_LOGIN_DT,
			FIELD_LOGOUT_DT,
			FIELD_MAX_SESSION_DT,
			FIELD_TIMEOUT_DT,

			FIELD_CREATED_DT,
			FIELD_UPDATED_DT,
			FIELD_DELETED_DT,

			strings.Repeat(`?, `, 20), `?`)

	}

	stmt, err := t.db.Prepare(cmd_user_insert)

	if err != nil {
		return NewGeneralFromError(err, http.StatusInternalServerError)
	}

	result, err := stmt.Exec(
		user.Domain,
		user.Email,
		user.GetFailCountStr(),
		user.FullName,
		user.Guid,
		user.LoginName,
		user.Password,
		user.Salt,
		user.Token,

		strconv.FormatBool(user.IsActive),
		strconv.FormatBool(user.IsLoggedIn),
		strconv.FormatBool(user.IsSystem),

		user.GetLastAuthAtStr(),
		user.GetLastFailedAtStr(),
		user.GetLoginAtStr(),
		user.GetLogoutAtStr(),
		user.GetMaxSessionAtStr(),
		user.GetTimeoutStr(),

		user.GetCreatedAtStr(),
		user.GetUpdatedAtStr(),
		user.GetDeletedAtStr(),
	)
	if err != nil {
		if err := t.checkUserExists(user); err != ErrUserNotFound {
			return err
		}
		return NewGeneralFromError(err, http.StatusInternalServerError)
	}

	if count, err := result.RowsAffected(); err != nil {
		return NewGeneralFromError(err, http.StatusInternalServerError)
	} else {
		if count == 0 {
			return ErrUserNotRegistered
		}
	}

	return nil
}

// Release is used to release any locks/resources that may have been created. In SQLITE we
// aren't using any locks, so we don't have to do anything.
func (t *SqliteConn) Release() error {
	return nil
}

func (t *SqliteConn) checkUserExists(user *record.User) error {
	var guid, domain, email, login sql.NullString

	s := fmt.Sprintf(`SELECT
				%s,
				%s ,
				%s,
				%s
			FROM %s
			WHERE %s = ?
			   OR ( %s = ? AND ( %s = ? OR %s = ?))`,
		FIELD_GUID, /* SELECT ... */
		FIELD_DOMAIN,
		FIELD_EMAIL,
		FIELD_LOGINNAME,

		record.USER_STORE_NAME, /* FROM ... */

		FIELD_GUID, /* WHERE ... */
		FIELD_DOMAIN,
		FIELD_EMAIL,
		FIELD_LOGINNAME,
	)

	stmt, err := t.db.Prepare(s)
	if err == nil {

		row := stmt.QueryRow(
			user.Guid,
			user.Domain,
			user.Email,
			user.LoginName)

		if row.Scan(&guid, &domain, &email, &login) == nil {
			if guid.Valid && guid.String == user.Guid {
				return ErrDuplicateGuid
			}
			if email.Valid && email.String == user.Email {
				return ErrDuplicateEmail
			}
			if login.Valid && login.String == user.LoginName {
				return ErrDuplicateLogin
			}
		}
		return ErrUserNotFound
	}

	return ErrInternalDatabase
}
