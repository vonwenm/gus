// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package

package record

import (
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"time"
	"github.com/cgentry/gus/encryption"
)

const (
	USER_OK                  = iota
	USER_EXPIRED             = iota
	USER_INVALID             = iota
	USER_PASSWD_TOO_SHORT    = iota
	USER_PASSWORD_TOO_SIMPLE = iota
)

const USER_TIME_STR = time.RFC3339


func init() {
	userControl = new(UserControl)
	userControl.SetMaxDuration("24h")
	userControl.SetTimeout("20m")
}

// A control variable holds all of the values we need to control logging in
var userControl * UserControl

type User struct {
	Id       int  // our internal ID
	FullName string // Simple full name
	Email    string // User's primary email
	IsSystem bool   // User is INTERNAL or EXTERNAL

	Guid string // Simple, unique user ID for external use

	Domain    string // What system is using this user
	LoginName string // What they login with
	Password  string // Password (encrypted) for user
	Token     string // Generated at login-time

	Salt string // Magic number used to hash values for user

	IsActive   bool // Is this an active user
	IsLoggedIn bool // Is this user currently logged in

	LoginAt      time.Time // Last login time
	LogoutAt     time.Time // Last logout time
	LastAuthAt   time.Time // Last successful Authorisation
	LastFailedAt time.Time // Last failed login
	FailCount    int       // Current number of failed logins

	MaxSessionAt time.Time // When they MUST logout by
	TimeoutAt  time.Time // Required to authenticate by

	CreatedAt time.Time // Creation date (immutable)
	UpdatedAt time.Time // Last updated
	DeletedAt time.Time // When deleted

}

func ( u * User ) String() string {
	return fmt.Sprintf( `User Record
	Guid: 		'%s'
	Domain:		'%s'
	LoginName:	'%s'
	Token:		'%s'
	` ,
	u.GetGuid() ,
	u.GetDomain() ,
	u.GetLoginName(),
	u.Token )

}



type UserControl struct {
	MaximumSessionDuration  time.Duration
	TimeSinceAuthentication time.Duration
}

/*
 *				INTERNAL ROUTINES
 */

// CreateSalt will create a magic number for use with other functions,
// like creating a GUID or a token.
func CreateSalt(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%x", b)
}

/*
 * 		UserControl routines. These set the fields for timeout and control
 */

func (uc * UserControl) SetMaxDuration(interval string) (err error) {
	uc.MaximumSessionDuration, err = time.ParseDuration(interval)
	return err
}

func (uc * UserControl) SetTimeout(interval string) (err error) {
	uc.TimeSinceAuthentication, err = time.ParseDuration(interval)
	return err
}

/*
 *			PUBLIC ROUTINES
 *			-- User --
 */
// NewUser creates a new user record with only the Domain filled in.
// The "Salt" field is a crypto-random number in order to produce
// unique values
func NewUser(domain string) * User {
	user := new( User )
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.IsSystem  = false

	user.SetDomain(domain)
	user.Salt = CreateSalt(20)
	user.Token = ""
	return user
}

// NewTestUser will generate a nonsense test user
func NewTestUser() * User {
	user := NewUser("_test")
	user.SetName( "test")
	user.SetEmail( "test@nowhere.com" )
	return user
}

// CreateToken will generate a short-use token for confirmation with authentication.
// The token can be used as a ticket until it expires. Any program can gain access
// to user information with it. Tokens can be saves
func (user * User) CreateToken() string {
	guid := md5.New()
	guid.Write([]byte(user.GetGuid()))			// Always based on user's GUID
	guid.Write([]byte(CreateSalt(20)))	// And a non-repeatable magic number
	out := guid.Sum(nil)
	return fmt.Sprintf("%x-%x-%x-%x-%x", out[0:4], out[4:6], out[6:8], out[8:10], out[10:len(out)])
}


// CheckExpirationDates will see if the token is valid or expired. If it
// is expired, the token will be cleared and the proper status will be set
func (user * User) CheckExpirationDates() int {
	if user.Token != "" && user.IsLoggedIn {
		if user.LastAuthAt.Before(user.MaxSessionAt) || user.LastAuthAt.Before(user.TimeoutAt) {
			user.LastAuthAt = time.Now()
			return USER_OK
		}
		user.Logout()
		return USER_EXPIRED
	}
	user.Logout()
	return USER_INVALID
}

// Authenticate checks the user's token to see if it is valid. This is a post-login process
// The user's record should be saved after this operation
func (user * User) Authenticate(token string) int {
	if token != "" && user.IsLoggedIn {
		if checkToken, err := user.GetTokenWithExpiration(); err == USER_OK && token == checkToken {
			return USER_OK
		}
	}
	return USER_INVALID
}


// Login will authenticate the user and create the tokens required later
func (user * User) Login(password string) (int, error) {

	now := time.Now()                     // Get time marker all the times

	if status := user.CheckPassword( password ); status != USER_OK {
		user.LastFailedAt = now // Save failure date/time
		user.IsLoggedIn = false // Mark as not logged in
		user.Token = ""         // Clear the token
		user.FailCount++        // Increment failure count
		return status , errors.New( "Invalid login/password")
	}
	pwd := encryption.GetDriver().EncryptPassword(password , user.Salt) // Encrypt password


	if pwd != user.Password { // Password differ?
		user.LastFailedAt = now // Save failure date/time
		user.IsLoggedIn = false // Mark as not logged in
		user.Token = ""         // Clear the token
		user.FailCount++        // Increment failure count
		return USER_INVALID, errors.New("Invalid login/password for user")
	}

	user.Token = user.CreateToken() // Give him a ticket

	user.MaxSessionAt = now.Add(userControl.MaximumSessionDuration)
	user.TimeoutAt = now.Add(userControl.TimeSinceAuthentication)
	user.LastAuthAt = now
	user.IsLoggedIn = true
	user.FailCount = 0

	return USER_OK, nil
}

// Logout will mark the record as 'logged out' and the user will be removed from the system
func (user * User) Logout() {
	user.Token = ""
	user.IsLoggedIn = false
	user.LogoutAt = time.Now()
}



// ChangePassword to the new password. The user must be logged in for this
func ( user * User ) ChangePassword( oldPassword, token, newPassword string) int {
	if user.Authenticate(token) == USER_OK {
		t := encryption.GetDriver()
		if t.EncryptPassword( oldPassword , user.Salt) == user.Password {
			if check := user.CheckNewPassword(newPassword); check != USER_OK {
				return check
			}
			user.Password = t.EncryptPassword( newPassword , user.Salt)
			return USER_OK
		}
	}
	return USER_INVALID
}

func ( user * User ) CheckPassword(  testPassword string ) int {
	pwd := encryption.GetDriver().EncryptPassword(testPassword , user.Salt) // Encrypt password

	if pwd != user.GetPassword() {
		return USER_INVALID
	}
	return USER_OK
}





