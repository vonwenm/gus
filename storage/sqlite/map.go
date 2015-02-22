// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"database/sql"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/record/mappers"
)

func mapColumnsToUser(rows *sql.Rows) []*tenant.User {

	var allUsers []*tenant.User
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	vpoint := make([]interface{}, count)
	var vstr string

	for rows.Next() {
		for i := range columns {
			vpoint[i] = &values[i]
		}
		user := tenant.NewUser()
		rows.Scan(vpoint...)

		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				vstr = string(b)
				mappers.UserField(user,col, vstr)
			} // End columns

			allUsers = append(allUsers, user)
		}
	}
	return allUsers
}
