package gus

import (
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"time"
)

const (
	USER_OK                  = iota
	USER_EXPIRED             = iota
	USER_INVALID             = iota
	USER_PASSWD_TOO_SHORT    = iota
	USER_PASSWORD_TOO_SIMPLE = iota
)

func init() {
	userControl = new(UserControl)
	userControl.SetMaxDuration("24h")
	userControl.SetTimeout("20m")
}

// A control variable holds all of the values we need to control logging in
var userControl *UserControl

type User struct {
	Id       int64  // our internal ID
	FullName string // Simple full name
	Email    string // User's primary email

	Guid string // Simple, unique user ID for external use

	Domain    string // What system is using this user
	LoginName string // What they login with
	Password  string // Password (encrypted) for user
	Token     string // Generated at login-time

	MagicNumber string // Magic number used to hash values for user

	IsActive   bool // Is this an active user
	IsLoggedIn bool // Is this user currently logged in

	LoginAt      time.Time // Last login time
	LogoutAt     time.Time // Last logout time
	LastAuthAt   time.Time // Last successful Authorisation
	LastFailedAt time.Time // Last failed login
	FailCount    int       // Current number of failed logins

	MaxSession time.Time // When they MUST logout by
	TimeoutAt  time.Time // Required to authenticate by

	CreatedAt time.Time // Creation date (immutable)
	UpdatedAt time.Time // Last updated
	DeletedAt time.Time // When deleted

}

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

type UserControl struct {
	MaximumSessionDuration  time.Duration
	TimeSinceAuthentication time.Duration
}

/*
 *				INTERNAL ROUTINES
 */

// createMagicNumber will create a magic number for use with other functions,
// like creating a GUID or a token.
func createMagicNumber(len int) string {
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

func (uc *UserControl) SetMaxDuration(interval string) (err error) {
	uc.MaximumSessionDuration, err = time.ParseDuration(interval)
	return err
}

func (uc *UserControl) SetTimeout(interval string) (err error) {
	uc.TimeSinceAuthentication, err = time.ParseDuration(interval)
	return err
}

/*
 *			PUBLIC ROUTINES
 *			-- User --
 */
// NewUser creates a new user record with only the Domain filled in.
// The "MagicNumber" field is a crypto-random number in order to produce
// unique values
func NewUser(domain string) User {
	user := User{}
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	user.SetDomain(domain)
	user.MagicNumber = createMagicNumber(20)
	user.Token = ""
	return user
}

// Set, or reset, the user's ID. When an ID is set, the GUID is reset.
func (user *User) SetID(id int64) *User {
	if id > 0 && user.Id == 0 {
		user.Id = id
	}
	return user
}

// SetName sets the fullname for the user
func (user *User) SetName(name string) *User {
	user.FullName = name
	return user
}

// SetLoginAt - set the logged in date/time as now
func (user *User) SetLoginAt() *User {
	user.LoginAt = time.Now()
	return user
}

// GetCreatedAt to be now. The created at can only be done once
func (user *User) GetCreatedAt() time.Time {
	return user.CreatedAt
}

// SetUpdateAt will set the update time stamp to now
func (user *User) SetUpdatedAt() *User {
	user.UpdatedAt = time.Now()
	return user
}

// GetGuid will return the unique guid for this user
func (user *User) GetGuid() string {
	if user.Guid == "" {
		guid := md5.New()
		guid.Write([]byte(user.Domain))						// Add in the user's domain
		guid.Write([]byte( user.GetCreatedAt().String()))	// Add in the creation string
		guid.Write([]byte(user.GetMagicNumber()))			// And the user's magic (unique) number
		out := guid.Sum(nil)
		user.Guid = fmt.Sprintf("%x-%x-%x-%x-%x", out[0:4], out[4:6], out[6:8], out[8:10], out[10:len(out)])
	}
	return user.Guid
}

// SetDomain will set the domain name for this record.
func (user *User) SetDomain(name string) *User {
	user.Domain = name
	return user
}

// GetDomain will get the domain name for a user
func (user *User) GetDomain() string {
	return user.Domain
}

// GetMagicNumber will get the special account-specific magic number.
// Normally used for salting various other functions, like password
func (user *User) GetMagicNumber() string {
	return user.MagicNumber
}

func (user *User) GetPassword() string {
	return user.Password
}

func (user *User ) SetPassword( newPassword string ) int {
	if user.Password != "" {
		return USER_INVALID
	}
	t,found := pwdDrivers[driverName]
	if ! found {
		panic( "No encryption library found")
	}
	user.Password = t.EncryptPassword( user , newPassword )

	return USER_OK
}

// GetToken will check the status of the user and return the token from the record
func (user *User) GetToken() (string, int) {
	rtn := user.CheckExpirationDates()
	return user.Token, rtn
}

// CreateToken will generate a short-use token for confirmation with authentication.
// The token can be used as a ticket until it expires. Any program can gain access
// to user information with it. Tokens can be saves
func (user *User) CreateToken() string {
	guid := md5.New()
	guid.Write([]byte(user.GetGuid()))			// Always based on user's GUID
	guid.Write([]byte(createMagicNumber(20)))	// And a non-repeatable magic number
	out := guid.Sum(nil)
	return fmt.Sprintf("%x-%x-%x-%x-%x", out[0:4], out[4:6], out[6:8], out[8:10], out[10:len(out)])
}


// CheckExpirationDates will see if the token is valid or expired. If it
// is expired, the token will be cleared and the proper status will be set
func (user *User) CheckExpirationDates() int {
	if user.Token != "" && user.IsLoggedIn {
		if user.LastAuthAt.Before(user.MaxSession) || user.LastAuthAt.Before(user.TimeoutAt) {
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
func (user *User) Authenticate(token string) int {
	if token != "" && user.IsLoggedIn {
		if checkToken, err := user.GetToken(); err == USER_OK && token == checkToken {
			return USER_OK
		}
	}
	return USER_INVALID
}

// Login will authenticate the user and create the tokens required later
func (user *User) Login(password string) (int, error) {

	now := time.Now()                     // Get time marker all the times

	if status := user.CheckPassword( password ); status != USER_OK {
		user.LastFailedAt = now // Save failure date/time
		user.IsLoggedIn = false // Mark as not logged in
		user.Token = ""         // Clear the token
		user.FailCount++        // Increment failure count
		return status , errors.New( "Invalid login/password")
	}
	t,_ := pwdDrivers[driverName]
	pwd := t.EncryptPassword(user , password) // Encrypt password


	if pwd != user.Password { // Password differ?
		user.LastFailedAt = now // Save failure date/time
		user.IsLoggedIn = false // Mark as not logged in
		user.Token = ""         // Clear the token
		user.FailCount++        // Increment failure count
		return USER_INVALID, errors.New("Invalid login/password for user")
	}

	user.Token = user.CreateToken() // Give him a ticket

	user.MaxSession = now.Add(userControl.MaximumSessionDuration)
	user.TimeoutAt = now.Add(userControl.TimeSinceAuthentication)
	user.LastAuthAt = now
	user.IsLoggedIn = true
	user.FailCount = 0

	return USER_OK, nil
}

// Logout will mark the record as 'logged out' and the user will be removed from the system
func (user *User) Logout() {
	user.Token = ""
	user.IsLoggedIn = false
	user.LogoutAt = time.Now()
}



// ChangePassword to the new password. The user must be logged in for this
func ( user *User ) ChangePassword( oldPassword, token, newPassword string) int {
	if user.Authenticate(token) == USER_OK {
		t,found := pwdDrivers[driverName]
		if  ! found {
			panic("No password encryption set")
		}
		if t.EncryptPassword(user , oldPassword) == user.Password {
			if check := user.CheckNewPassword(newPassword); check != USER_OK {
				return check
			}
			user.Password = t.EncryptPassword( user , newPassword)
			return USER_OK
		}
	}
	return USER_INVALID
}

func ( user *User ) CheckPassword(  testPassword string ) int {
	t,found := pwdDrivers[driverName]
	if ! found {
		panic("No password encryption set")
	}
	pwd := t.EncryptPassword(user , testPassword) // Encrypt password

	if pwd != user.GetPassword() {
		return USER_INVALID
	}
	return USER_OK
}


/*
 *			PUBLIC ROUTINES
 *			-- UserReturn --
 */
// GetUserReturn
func (user *User) GetUserReturn() UserReturn {

	user.Token = user.CreateToken()
	rtn := UserReturn{}
	rtn.Guid = user.GetGuid()
	rtn.Token = user.Token
	rtn.LoginAt = user.LoginAt.Format(time.RFC3339)
	rtn.LastAuthAt = user.LoginAt.Format(time.RFC3339)
	rtn.CreatedAt = user.CreatedAt.Format(time.RFC3339)

	return rtn

}

func (user *User) String() string {
	return fmt.Sprintf( "Name: %s\nPassword: %s\nEmail: %s\n" ,
		user.FullName , user.Password , user.Email )
}

/*
 *			Dynamic interfaces
 *			-- UserReturn --
 */
var pwdDrivers = make( map[string] Passworder )
var driverName = "sha512"

func RegisterPassword( name string , driver Passworder ){
	if driver == nil {
		panic("Password driver: Register driver is nil")
	}
	if _,dup := pwdDrivers[name] ; dup {
		panic( "Password driver: Register called twice for '" + name + "'")
	}
	pwdDrivers[name] = driver
}

func SetPasswordEncryption( name string ){
	driverName = name
}
