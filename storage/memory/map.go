package memory

import (
	"strings"
	"github.com/cgentry/gus/record"
	"database/sql"
)
func mapColumnsToUser(rows * sql.Rows) []*record.User {

	var allUsers [] *record.User
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

		for i , col := range columns {
			val := values[i]
			if b, ok := val.([]byte) ; ok {
				vstr = string(b)


				switch strings.ToLower(col) {
				case "fullname" : user.SetName( vstr )
				case "email"	: user.SetEmail( vstr )
				case "guid"     : user.SetGuid( vstr )

				case "domain"	: user.SetDomain( vstr )
				case "password" : user.SetPasswordStr( vstr )
				case "token"    : user.SetToken( vstr )

				case "salt"	    : user.SetSalt( vstr )
				case "isactive"	: user.SetIsActive( StrToBool(vstr ) )
				case "isloggedin" : user.SetIsLoggedIn( StrToBool( vstr) )

				case "loginat"	: user.SetLoginAt( StrToTime(  vstr ) )
				case "logoutat" : user.SetLogoutAt( StrToTime(  vstr ) )
				case "lastfailedat" : user.SetLastFailedAt( StrToTime(vstr) )
				case "failcount" : user.SetFailCount( StrToInt(vstr) )

				case "maxsessionat": user.SetMaxSessionAt( StrToTime(vstr ) )
				case "timeoutat" : user.SetTimeoutAt( StrToTime(  vstr ) )

				case "createdat" : user.SetCreatedAt( StrToTime(  vstr ) )
				case "updatedat" : user.SetUpdatedAt( StrToTime(  vstr ) )
				case "deletedat" : user.SetDeletedAt( StrToTime(  vstr ) )
				case "loginname" : user.SetLoginName( vstr )

				}
			}
		} // End columns

		allUsers = append(allUsers, user)
	}
	return allUsers
}
