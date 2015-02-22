// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"fmt"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/storage"
	"net/http"
	"strings"
)

func (t *SqliteConn) UserFetch(domain, field, val string) (*tenant.User, error) {
	if domain == storage.MATCH_ANY_DOMAIN {
		return t.fetchUserByFieldAny(field, val)
	}
	return t.fetchUserByField(domain, field, val)
}
func (t *SqliteConn) fetchUserByField(domain, field, val string) (*tenant.User, error) {
	field = strings.TrimSpace(field)
	if field == `` {
		return nil, ErrEmptyFieldForLookup
	}
	if t.db == nil {
		return nil, ErrNotOpen
	}
	cmd := fmt.Sprintf(`SELECT *
			 FROM %s
			WHERE %s = ?
			  AND %s = ?`,
		tenant.USER_STORE_NAME,
		FIELD_DOMAIN,
		field)
	rows, err := t.db.Query(cmd, domain, val)
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
func (t *SqliteConn) fetchUserByFieldAny(field, val string) (*tenant.User, error) {
	field = strings.TrimSpace(field)
	if field == `` {
		return nil, ErrEmptyFieldForLookup
	}
	if t.db == nil {
		return nil, ErrNotOpen
	}
	cmd := fmt.Sprintf(`SELECT *
			 FROM %s
			WHERE  %s = ?`,
		tenant.USER_STORE_NAME,
		field)

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
