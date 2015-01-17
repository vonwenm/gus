package record

/*
 * All of the setters will return either nil or ERROR. This is atypical of most
 * other routines, but this is where rules are enforced for the record's fields
 */

import (
	"errors"
	"github.com/cgentry/gus/encryption"
	"strings"
	"time"
	//"fmt"
)

// Set, or reset, the user's ID. When an ID is set, the GUID is reset.
func (user *User) SetID(id int) error {

	if id > 0 && user.Id == 0 {
		user.Id = id
		return nil
	}
	return errors.New("User id cannot be set")

}

// SetName sets the fullname for the user
func (user *User) SetName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) > 0 {
		user.FullName = name
		return nil
	}
	return errors.New("Name cannot be empty")
}

func (user *User) SetEmail(val string) error {
	user.Email = val
	return nil
}

func (user *User) SetGuid(val string) error {
	if len(val) < 32 {
		return errors.New("GUID must be at least 32 characters long")
	}
	user.Guid = val
	return nil
}

func (user *User) SetLoginName(name string) error {
	user.LoginName = name
	return nil
}

func (user *User) SetToken(val string) error {
	user.Token = val
	return nil
}

func (user *User) SetLoginAt(t time.Time) error {
	user.LoginAt = t
	return nil
}
func (user *User) SetUpdatedAt(t time.Time) error {
	user.UpdatedAt = t
	return nil
}

// SetDomain will set the domain name for this record.
func (user *User) SetDomain(val string) error {
	user.Domain = val
	return nil
}

func (user *User) SetPasswordStr(pwd string) error {
	user.Password = pwd
	return nil
}
func (user *User) SetPassword(newPassword string) error {
	newPassword = strings.TrimSpace(newPassword)
	if len(newPassword) < 6 {
		return errors.New("Password must be at least 6 characters")
	}
	user.Password = encryption.GetDriver().EncryptPassword(newPassword, user.Salt)
	return nil
}

func (user *User) SetSalt(val string) error {
	user.Salt = val
	return nil
}

func (user *User) SetIsActive(val bool) error {
	user.IsActive = val
	return nil
}

func (user *User) SetIsLoggedIn(val bool) error {
	user.IsLoggedIn = val
	return nil
}

func (user *User) SetIsSystem(val bool) error {
	user.IsSystem = val
	return nil
}

func (user *User) SetLogoutAt(t time.Time) error {
	user.LogoutAt = t
	return nil
}

func (user *User) SetLastFailedAt(t time.Time) error {
	user.LastFailedAt = t
	return nil
}

func (user *User) SetLastAuthAt(t time.Time) error {
	user.LastAuthAt = t
	return nil
}

func (user *User) SetFailCount(i int) error {
	user.FailCount = i
	return nil
}

func (user *User) SetMaxSessionAt(t time.Time) error {
	user.MaxSessionAt = t
	return nil
}

func (user *User) SetTimeoutAt(t time.Time) error {
	user.TimeoutAt = t
	return nil
}

func (user *User) SetCreatedAt(t time.Time) error {
	user.CreatedAt = t
	return nil
}

func (user *User) SetDeletedAt(t time.Time) error {
	user.DeletedAt = t
	return nil
}
