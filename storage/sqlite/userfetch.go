// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"strings"
	. "github.com/cgentry/gus/ecode"
	"net/http"
)

func (t *SqliteConn) fetchUserByField(field, val string) (*record.User, error) {
	field = strings.TrimSpace(field)
	if field == `` {
		return nil, ErrEmptyFieldForLookup
	}
	if t.db == nil {
		return nil, ErrNotOpen
	}
	cmd := fmt.Sprintf(`SELECT * FROM User WHERE %s = ?`, field)
	rows, err := t.db.Query(cmd, val)
	if err != nil {
		return nil, NewGeneralFromError(err, http.StatusInternalServerError)
	}
	defer rows.Close()

	users := mapColumnsToUser(rows)
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil

}

func (t *SqliteConn) FetchUserByToken(token string) (*record.User, error) {
	return t.fetchUserByField(FIELD_TOKEN, token)
}

func (t *SqliteConn) FetchUserByGuid(guid string) (*record.User, error) {
	return t.fetchUserByField(FIELD_GUID, guid)
}

func (t *SqliteConn) FetchUserByEmail(email string) (*record.User, error) {
	return t.fetchUserByField(FIELD_EMAIL, email)
}

func (t *SqliteConn) FetchUserByLogin(value string) (*record.User, error) {
	return t.fetchUserByField(FIELD_LOGINNAME, value)
}
