// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package memory

import (
	"database/sql"
	"github.com/cgentry/gus/record"
)

func mapColumnsToUser(rows *sql.Rows) []*record.User {

	var allUsers []*record.User
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	vpoint := make([]interface{}, count)
	var vstr string

	for rows.Next() {
		for i, _ := range columns {
			vpoint[i] = &values[i]
		}
		user := record.NewUser("")
		rows.Scan(vpoint...)

		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				vstr = string(b)
				user.MapFieldToUser(col,vstr)
		} // End columns

		allUsers = append(allUsers, user)
	}
	}
	return allUsers
}

