package record

import (
	"time"
)

// Record passed back to the caller
type UserReturn struct {
	Guid  string // External ID
	Token string // Internal user access token

	FullName string
	Email    string

	LoginAt    string
	LastAuthAt string

	TimeoutAt time.Time // Required to authenticate by

	CreatedAt string
}

func NewReturnFromUser( user * User) UserReturn {

		rtn := UserReturn{ Token: user.CreateToken() }
		rtn.Guid = user.GetGuid()
		rtn.Token = user.Token

		rtn.LoginAt = user.LoginAt.Format(time.RFC3339)
		rtn.LastAuthAt = user.LoginAt.Format(time.RFC3339)
		rtn.CreatedAt = user.CreatedAt.Format(time.RFC3339)

		return rtn

}
