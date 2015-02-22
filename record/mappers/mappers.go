// Mappers contains functions that will map an input to an output. These are
// convenience functions that are gathered here to reduce the problem of circular references
package mappers

import (
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/record/response"
	"strings"
	"strconv"
	"errors"
)

//ResponseFromUser takes a user record and copies the relevant fields from the
// user's record.
func ResponseFromUser(rtn *response.UserReturn , user *tenant.User ) * response.UserReturn {

	rtn.Guid = user.Guid
	rtn.Token = user.Token

	rtn.LoginAt = user.LoginAt
	rtn.LastAuthAt = user.LastAuthAt
	rtn.CreatedAt = user.CreatedAt
	rtn.TimeoutAt = user.TimeoutAt
	rtn.MaxSessionAt = user.MaxSessionAt

	rtn.FullName = user.FullName
	rtn.Email = user.Email
	rtn.LoginName = user.LoginName

	return rtn
}


// UserField will find map a fieldname to a user record and save the field in the record
func UserField( user * tenant.User, key, value string) (found bool, rtn error) {

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
	case "lastauthat":
		rtn = user.SetLastAuthAt(StrToTime(value))
	case "failcount":
		rtn = user.SetFailCount(StrToInt(value))

	case "maxsessionat":
		rtn = user.SetMaxSessionAt(StrToTime(value))
	case "timeoutat":
		rtn = user.SetTimeoutAt(StrToTime(value))

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


/**
 * Create a record from the user record passed to this routine
 * See:		UserReturn
 */
func UserFromCli( rtn * tenant.User , r * tenant.UserCli ) (ortn *tenant.User, err error) {

	ortn = rtn
	if err = rtn.SetDomain(r.Domain); err != nil {
		return
	}
	if err = rtn.SetName(r.FullName); err != nil {
		return
	}
	if err = rtn.SetEmail(r.Email); err != nil {
		return
	}
	if err = rtn.SetLoginName(r.LoginName); err != nil {
		return
	}
	if err = rtn.SetPassword(r.Password); err != nil {
		return
	}
	rtn.SetIsActive(r.Enable)
	err = rtn.SetIsSystem(r.Level == "client")

	return

}
