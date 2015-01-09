package record

import (
	"time"
)

// UserReturn will contain the USER's minimum data from a login/register operation.
type UserReturn struct {
	Guid  string // Permanent User ID for external linking (Within systems)
	Token string // Send THIS to login with

	FullName  string // FULL user name ("Jane Doe")
	LoginName string // The ID they use to login with
	Email     string // Email address

	LoginAt    string // THIS login time
	LastAuthAt string // Last login time

	TimeoutAt time.Time // Required to authenticate by

	CreatedAt string // When the user was created
}

/**
 * Create a record from the user record passed to this routine
 * See:		UserReturn
 */
func NewReturnFromUser(user *User) UserReturn {

	rtn := UserReturn{}
	rtn.Guid = user.GetGuid()
	if user.Token == "" {
		user.Token = user.CreateToken()
	}
	rtn.Token = user.Token

	rtn.LoginAt = user.LoginAt.Format(time.RFC3339)
	rtn.LastAuthAt = user.LoginAt.Format(time.RFC3339)
	rtn.CreatedAt = user.CreatedAt.Format(time.RFC3339)
	rtn.FullName = user.GetFullName()
	rtn.Email = user.GetEmail()
	rtn.LoginName = user.GetLoginName()

	return rtn

}
