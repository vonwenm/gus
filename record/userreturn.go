package record

import (
	"time"
)

// UserReturn will contain the USER's minimum data from any operation operation.
type UserReturn struct {
	Guid  string // Permanent User ID for external linking (Within systems)
	Token string // Send THIS to login with

	FullName  string // FULL user name ("Jane Doe")
	LoginName string // The ID they use to login with
	Email     string // Email address

	LoginAt      time.Time // THIS login time
	LastAuthAt   time.Time // Last login time
	TimeoutAt    time.Time // Required to authenticate by
	MaxSessionAt time.Time // This is when the user will be forced off

	CreatedAt time.Time // When the user was created
}

/**
 * Create a record from the user record passed to this routine
 * See:		UserReturn
 */
func NewReturnFromUser(user *User) UserReturn {

	rtn := UserReturn{}
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
