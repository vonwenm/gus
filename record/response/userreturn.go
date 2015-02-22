package response

import (
	"time"
	"github.com/cgentry/gus/record/stamp"
)

// User will contain the USER's minimum data from any operation operation.
type UserReturn struct {
	stamp.Timestamp

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

func NewUserReturn( ) *UserReturn  {
	rtn := &UserReturn{}
	rtn.SetStamp( time.Now() )

	return rtn
}
func ( u *UserReturn ) Check() error {
	return nil
}
