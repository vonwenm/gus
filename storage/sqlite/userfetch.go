// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"github.com/cgentry/gus/record"
	"errors"
	"fmt"
	//"database/sql"
)

func (t *StorageMem) fetchUserByField(field, val string) (*record.User, error) {
	cmd := fmt.Sprintf(`SELECT * FROM User WHERE %s = ?`, field)
	rows, err := t.db.Query(cmd, val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := mapColumnsToUser(rows)
	if len(users) == 0 {
		return nil, errors.New("No records found")
	}
	return users[0], err

}

func (t *StorageMem) FetchUserByToken(token string) (*record.User, error) {
	return t.fetchUserByField("Token", token)
}

func (t *StorageMem) FetchUserByGuid(guid string) (*record.User, error) {
	return t.fetchUserByField("Guid", guid)
}

func (t *StorageMem) FetchUserByEmail(email string) (*record.User, error) {
	return t.fetchUserByField("Email", email)
}

func (t *StorageMem) FetchUserByLogin(value string) (*record.User, error) {
	return t.fetchUserByField("LoginName", value)
}
