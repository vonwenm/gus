package record

import (
	"strings"
)

// Map commonly used names to fields within the user record
func ( user * User ) MapFieldToUser( key, value string ) (found bool) {
	
	found = true
	
	switch strings.ToLower(key) {

	case "name":
		fallthrough
	case "fullname":
		user.SetName(value)

	case "email":
		user.SetEmail(value)

	case "caller":
		fallthrough				// the field "caller" is used for guid as an identifier to make less confusion
	case "guid":
		user.SetGuid(value)

	case "domain":
		user.SetDomain(value)
	case "password":
		user.SetPasswordStr(value)
	case "token":
		user.SetToken(value)

	case "salt":
		user.SetSalt(value)
	case "isactive":
		user.SetIsActive(StrToBool(value,user.IsActive))
	case "isloggedin":
		user.SetIsLoggedIn(StrToBool(value,user.IsLoggedIn))
	case "issystem":
		user.SetIsSystem( StrToBool(value,user.IsActive))

	case "loginat":
		user.SetLoginAt(StrToTime(value))
	case "logoutat":
		user.SetLogoutAt(StrToTime(value))
	case "lastfailedat":
		user.SetLastFailedAt(StrToTime(value))
	case "failcount":
		user.SetFailCount(StrToInt(value))

	case "maxsessionat":
		user.SetMaxSessionAt(StrToTime(value))
	case "timeoutat":
		user.SetTimeoutAt(StrToTime(value))

	case "createdat":
		user.SetCreatedAt(StrToTime(value))
	case "updatedat":
		user.SetUpdatedAt(StrToTime(value))
	case "deletedat":
		user.SetDeletedAt(StrToTime(value))

	case "login":
		fallthrough
	case "loginname":
		user.SetLoginName(value)
	default:
		found = false

	}
	return found
}
