package record

import (
	"time"
	"github.com/cgentry/gus/encryption"
	//"fmt"
)

// Set, or reset, the user's ID. When an ID is set, the GUID is reset.
func (user * User) SetID(id int) * User {
	if id > 0 && user.Id == 0 {
		user.Id = id
	}
	return user
}
// SetName sets the fullname for the user
func (user * User) SetName(name string) * User {
	user.FullName = name
	return user
}

func (user * User) SetEmail(val string) * User {
	user.Email = val
	return user
}


func ( user * User ) SetGuid( val string ) * User {
	user.Guid = val
	return user
}

func (user * User) SetLoginName(name string) * User {
	user.LoginName = name
	return user
}


func (user * User) SetToken(val string) * User {
	user.Token = val
	return user
}

func ( user * User ) SetLoginAtStr( val string ) * User {
	return user
}

func ( user * User ) SetLoginAt( t time.Time ) * User {
	user.LoginAt = t
	return user
}

// SetLoginAt - set the logged in date/time as now
func (user * User) SetLoginAtNow() * User {
	user.LoginAt = time.Now()
	return user
}

// GetCreatedAt to be now. The created at can only be done once
func (user * User) GetCreatedAt() time.Time {
	return user.CreatedAt
}

// SetUpdateAt will set the update time stamp to now
func (user * User) SetUpdatedAtNow() * User {
	user.UpdatedAt = time.Now()
	return user
}

func ( user * User ) SetUpdatedAt( t time.Time ) * User {
	user.UpdatedAt = t
	return user
}

// SetDomain will set the domain name for this record.
func (user * User) SetDomain(val string) * User {
	user.Domain = val
	return user
}

func ( user * User ) SetPasswordStr( pwd string ) * User {
	user.Password = pwd
	return user
}
func (user * User) SetPassword(newPassword string) int {
	if user.Password != "" {
		return USER_INVALID
	}
	user.Password = encryption.GetDriver().EncryptPassword(newPassword, user.Salt)

	return USER_OK
}

func ( user * User ) SetSalt( val string ) * User {
	user.Salt = val
	return user
}

func ( user * User ) SetIsActive( val bool ) * User {
	user.IsActive = val
	return user
}

func ( user * User ) SetIsLoggedIn( val bool ) * User {
	user.IsLoggedIn = val
	return user
}

func ( user * User ) SetLogoutAt( t time.Time ) * User {
	user.LogoutAt = t
	return user
}

func ( user * User ) SetLastFailedAt( t time.Time ) * User {
	user.LastFailedAt = t
	return user
}

func ( user * User ) SetFailCount( i int) * User {
	user.FailCount = i
	return user
}

func ( user * User ) SetMaxSessionAt( t time.Time ) * User {
	user.MaxSessionAt = t
	return user
}

func ( user * User ) SetTimeoutAt( t time.Time ) * User {
	user.TimeoutAt = t
	return user
}

func ( user * User ) SetCreatedAt( t time.Time ) * User {
	user.CreatedAt= t
	return user
}

func ( user * User ) SetDeletedAt( t time.Time ) * User {
	user.DeletedAt = t
	return user
}



