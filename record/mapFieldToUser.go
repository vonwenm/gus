package record

/*
 *  MapFieldToUser
 *		This will take a key/value pair and map the input name into
 *		our record. Values here can be strings, integers or booleans.
 *
 *  See:
 *		user.go			- record definition
 *		usersetter.go	- Setters for value
 */
import (
	"errors"
	"strconv"
	"strings"
)

// Map commonly used names to fields within the user record
func (user *User) MapFieldToUser(key, value string) (found bool, rtn error) {

	iValue := 0
	found = true
	value = strings.TrimSpace(value) // No spaces around field

	switch strings.ToLower(key) {

	case "name":
		fallthrough
	case "fullname":
		rtn = user.SetName(value)

	case "email":
		rtn = user.SetEmail(value)

	case "caller":
		fallthrough // the field "caller" is used for guid as an identifier to make less confusion
	case "guid":
		rtn = user.SetGuid(value)

	case "domain":
		rtn = user.SetDomain(value)
	case "password":
		rtn = user.SetPasswordStr(value)
	case "token":
		rtn = user.SetToken(value)

	case "salt":
		rtn = user.SetSalt(value)
	case "isactive":
		rtn = user.SetIsActive(StrToBool(value, user.IsActive))
	case "isloggedin":
		rtn = user.SetIsLoggedIn(StrToBool(value, user.IsLoggedIn))
	case "issystem":
		rtn = user.SetIsSystem(StrToBool(value, user.IsActive))

	case "loginat":
		rtn = user.SetLoginAt(StrToTime(value))
	case "logoutat":
		rtn = user.SetLogoutAt(StrToTime(value))
	case "lastfailedat":
		rtn = user.SetLastFailedAt(StrToTime(value))
	case "failcount":
		rtn = user.SetFailCount(StrToInt(value))

	case "maxsessionat":
		rtn = user.SetMaxSessionAt(StrToTime(value))
	case "timeoutat":
		rtn = user.SetTimeoutAt(StrToTime(value))

	case "createdat":
		rtn = user.SetCreatedAt(StrToTime(value))
	case "updatedat":
		rtn = user.SetUpdatedAt(StrToTime(value))
	case "deletedat":
		rtn = user.SetDeletedAt(StrToTime(value))

	case "login":
		fallthrough
	case "loginname":
		rtn = user.SetLoginName(value)

	case "id":
		iValue, rtn = strconv.Atoi(value)
		if rtn == nil {
			rtn = user.SetID(iValue)
		}

	default:
		found = false
		rtn = errors.New("Invalid field " + key)

	}
	return
}
