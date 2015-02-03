// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package

package record

import (
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/encryption"
	"time"
)

const (
	USER_OK                  = iota
	USER_EXPIRED             = iota
	USER_INVALID             = iota
	USER_PASSWD_TOO_SHORT    = iota
	USER_PASSWORD_TOO_SIMPLE = iota
)

const USER_TIME_STR = time.RFC3339

// Standard name for the user store.
const USER_STORE_NAME = "User"

func initialiseControl() *UserControl {
	u := UserControl{}

	u.SetMaxDuration("24h")
	u.SetTimeout("20m")
	return &u
}

/*
 * 		UserControl routines. These set the fields for timeout and control
 */
// A control variable holds all of the values we need to control logging in
var userControl = initialiseControl()

type UserControl struct {
	MaximumSessionDuration  time.Duration
	TimeSinceAuthentication time.Duration
}

func (uc *UserControl) SetMaxDuration(interval string) (err error) {
	uc.MaximumSessionDuration, err = time.ParseDuration(interval)
	return err
}

func (uc *UserControl) SetTimeout(interval string) (err error) {
	uc.TimeSinceAuthentication, err = time.ParseDuration(interval)

	return err
}

type UserInterface interface {
	Login(string) error
	Logout()
	ChangePassword(oldPassword, newPassword string) error
	Authenticate(token string) error
}

// User is the internal record used to store all of the data that is held
// for a single user. The database routines need to take care of serialising/mapping
// the data out to long-term storage (DB, File, etc.)
type User struct {
	Id       int
	FullName string 	`name:"User's fullname" help:"User's full name (title, first, surname)"`
	Email    string 	`name:"User's email address" help:"User's email address."`
	IsSystem bool   	`name:"System user" help:"True if the this is a client otherwise a standard user"`

	Guid string 		`name:"User's GUID" help:"How the user is identified by this system. A unique key"`

	Domain    string 	`name:"Domain" help:"The group that the user belongs to"`
	LoginName string 	`name:"Login name" help:"The name the user uses to login with"`
	Password  string 	`name:"Encrypted password" help:"This is the user's encrypted password."`
	Token     string // Generated at login-time

	Salt string // Magic number used to hash values for user

	IsActive   bool 	`name:"User is enabled"   help:"If disabled, the user will not be able to login"`
	IsLoggedIn bool // Is this user currently logged in

	LoginAt      time.Time // Last login time
	LogoutAt     time.Time // Last logout time
	LastAuthAt   time.Time // Last successful Authorisation
	LastFailedAt time.Time // Last failed login
	FailCount    int       // Current number of failed logins

	MaxSessionAt time.Time // When they MUST logout by
	TimeoutAt    time.Time // Required to authenticate by

	CreatedAt time.Time // Creation date (immutable)
	UpdatedAt time.Time // Last updated
	DeletedAt time.Time // When deleted

}

// This is the minimum data needed for a user's record. It is NOT used
// for anything other than a minimum set.
type UserJson struct {
	FullName  string
	Email     string
	LoginName string
	Password  string
}

func (u *User) String() string {
	return fmt.Sprintf(`User Record
	Guid: 		'%s'
	Domain:		'%s'
	LoginName:	'%s'
	Token:		'%s'
	IsLoggedIn:	'%t'
	IsActive:	'%t'
	TimeoutAt:  '%s'
	NOW:        '%s'
	`,
		u.Guid,
		u.Domain,
		u.LoginName,
		u.Token,
		u.IsLoggedIn,
		u.IsActive,
		u.GetTimeoutStr(),
		time.Now().Format(USER_TIME_STR))

}

/*
 *				INTERNAL ROUTINES
 */

// CreateSalt will create a magic number for use with other functions,
// like creating a GUID or a token.
func CreateSalt(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil { // This should never happen
		panic(err.Error()) // ...and won't be covered in coverage report
	} // ...if it does - we can't run the system

	return fmt.Sprintf("%x", b)
}

/*
 *			PUBLIC ROUTINES
 *			-- User --
 */
// NewUser creates a new, empty user record. The domain is set to blank and
// the "Salt" field is a crypto-random number in order to produce
// unique values
func NewUser() *User {
	now := time.Now()
	user := new(User)
	user.Id = 0 // Flag for 'not created'
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsSystem = false

	user.SetDomain("")
	user.Salt = CreateSalt(20)
	user.Token = CreateSalt(20)		// Tokens need to be unique
	user.IsActive = true
	user.IsLoggedIn = false
	user.GenerateGuid()
	return user
}

// NewTestUser will generate a nonsense test user
func NewTestUser() *User {
	user := NewUser()
	user.SetDomain("_test")
	user.SetName("test")
	user.SetEmail("test@nowhere.com")

	return user
}

// Generate a unique GUID for the user record. This GUID will be based upon random numbers
// and the creation string.
func (user *User) GenerateGuid() {
	if user.Guid == "" {
		guid := md5.New()
		guid.Write([]byte(user.GetCreatedAtStr())) // Add in the creation string
		guid.Write([]byte(user.Salt))              // And the user's magic (unique) number
		guid.Write([]byte(CreateSalt(60)))         // A bit more entropy (so it isn't repeatable)
		out := guid.Sum(nil)
		user.Guid = fmt.Sprintf("%x-%x-%x-%x-%x", out[0:4], out[4:6], out[6:8], out[8:10], out[10:])
	}
}

// CreateToken will generate a short-use token for confirmation with authentication.
// The token can be used as a ticket until it expires. Any program can gain access
// to user information with it.
func (user *User) CreateToken() string {
	guid := md5.New()
	guid.Write([]byte(user.Guid))      // Always based on user's GUID
	guid.Write([]byte(CreateSalt(20))) // And a non-repeatable magic number
	out := guid.Sum(nil)
	return fmt.Sprintf("%x-%x-%x-%x-%x", out[0:4], out[4:6], out[6:8], out[8:10], out[10:])
}

// CheckExpirationDates will see if the token is valid or expired. If it
// is expired, the token will be cleared and the proper status will be set
func (user *User) CheckExpirationDates() error {

	if user.LastAuthAt.Before(user.MaxSessionAt) && user.LastAuthAt.Before(user.TimeoutAt) {
		user.LastAuthAt = time.Now()
		user.TimeoutAt = user.LastAuthAt.Add(userControl.TimeSinceAuthentication)
		user.UpdatedAt = user.LastAuthAt
		return nil
	}
	user.Logout()
	return ErrSessionExpired

}

// Authenticate checks the user's token to see if it is valid. This is a post-login process
// The user's record should be saved after this operation
func (user *User) Authenticate(token string) error {
	user.UpdatedAt = time.Now()
	if token != "" && user.IsLoggedIn && token == user.Token {
		if err := user.CheckExpirationDates(); err == nil {
			user.LastAuthAt = user.UpdatedAt
			return nil
		}
	}
	return ErrSessionExpired
}

// Login will authenticate the user and create the tokens required later
func (user *User) Login(password string) error {

	now := time.Now() // Get time marker all the times
	user.UpdatedAt = now

	if err := user.CheckPassword(password); err != nil {
		user.LastFailedAt = now // Save failure date/time
		user.IsLoggedIn = false // Mark as not logged in
		user.Token = ""         // Clear the token
		user.FailCount++        // Increment failure count
		return err
	}

	user.Token = user.CreateToken() // Give him a ticket

	user.MaxSessionAt = now.Add(userControl.MaximumSessionDuration)
	user.TimeoutAt = now.Add(userControl.TimeSinceAuthentication)
	user.IsLoggedIn = true
	user.LastAuthAt = now
	user.LoginAt = now

	user.FailCount = 0

	return nil
}

// Logout will mark the record as 'logged out' and the user will be removed from the system
func (user *User) Logout() error {
	if !user.IsLoggedIn {
		return ErrUserNotLoggedIn
	}
	user.Token = ""
	user.IsLoggedIn = false
	user.LogoutAt = time.Now()
	user.UpdatedAt = user.LogoutAt
	user.IsLoggedIn = false
	return nil
}

// ChangePassword to the new password. The user must be logged in for this
func (user *User) ChangePassword(oldPassword, newPassword string) error {
	if user.Authenticate(user.Token) == nil {
		t := encryption.GetDriver()
		if t.EncryptPassword(oldPassword, user.Salt) == user.Password {
			if err := user.CheckNewPassword(newPassword); err != nil {
				return err
			}
			user.Password = t.EncryptPassword(newPassword, user.Salt)
			user.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrInvalidPasswordOrUser
}

func (user *User) CheckPassword(testPassword string) error {
	pwd := encryption.GetDriver().EncryptPassword(testPassword, user.Salt) // Encrypt password

	if pwd != user.Password {
		return ErrInvalidPasswordOrUser
	}
	return nil
}

// This is used when a user loses their password. They request a password reset
// based upon their email address. If they are logged in, they will not be able to
// reset the password ( 'User still logged in'). When they are NOT logged in, a token
// is generated and set. The client program later calls this with the email address
// AND the token. When this is confirmed, the password is then set to the value of the
// token and the user can go ahead and login as normal.
func (user *User) GenerateLostPassword() (newPassword string, err error) {

	return
}

func (user *User) ConfirmLostPassword(lostPwdToken string) (err error) {
	if lostPwdToken == user.Token {
		user.Password = encryption.GetDriver().EncryptPassword(lostPwdToken, user.Salt)
	} else {
		err = errors.New(`Invalid password or user id`)
	}
	return
}
