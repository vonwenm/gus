package storage

import (
	"fmt"
	"github.com/cgentry/gus/record"
	. "github.com/cgentry/gus/ecode"
)

const (
	BANK_USER_OK            = iota
	BANK_USER_NOT_FOUND     = iota
	BANK_USER_TOKEN_INVALID = iota
	BANK_USER_DATA_NOTFOUND = iota
)

var driverMap = make(map[string]Driver,2)
const driver_name = "Storage"

// Find the driver by the 'name' and add it into the map so it can be opened.
// Each driver can only be registered once. To remove all drivers, call
// ResetRegister()
func Register(name string, driver Driver) error {
	if driver == nil {
		panic(driver_name + " driver: Register driver is nil")
	}
	if _, dup := driverMap[name]; dup {
		return ErrAlreadyRegistered
	}
	driverMap[name] = driver
	return nil
}

func String() string {
	rtn := fmt.Sprintf("Length is %d\n", len(driverMap))
	for key := range driverMap {
		rtn = rtn + key + "\n"
	}
	return rtn
}

func ResetRegister() {
	driverMap = make(map[string]Driver)
}

// Store holds the state for any storage driver. It allows you to have
// consistent returns, such as getting the last error, discovering how
// a connection was made (connectString) or the name of the driver (name)
type Store struct {
	name          string
	connectString string
	isOpen        bool
	lastError     error
	driver        Driver
	connection    Conn
}

func IsRegistered(name string) bool {
	_, ok := driverMap[name]
	return ok
}

// Open a connection to the storage mechanism and return both a storage
// structure and an error status of the open
func Open(name string, connect string) (*Store, error) {
	s := &Store{
		name:          name,
		isOpen:        false,
		connectString: connect,
		lastError:     ErrNoDriverFound,
	}
	if driver, ok := driverMap[name]; ok {
		s.driver = driver
		s.connection, s.lastError = driver.Open(connect)
	}
	if s.lastError == nil {
		s.isOpen = true
	}
	return s, s.lastError
}

// Return the actual connection to the database for low-level access.
// This should be avoided unless you are coding for a very non-portable
// function
func (s *Store) GetStorageConnector() Conn {
	return s.connection
}

/*
 * The following functions are provided by this class and are not
 * encapsulated
 */
// Return the last known error condition that was given by a call
func (s *Store) GetLastError() error {
	return s.lastError
}

func (s *Store) IsOpen() bool {
	return s.isOpen
}

// Save the last error for later retrieval, and return the error
func (s *Store) saveLastError(e error) error {
	s.lastError = e
	return e
}

/*
 * Optional interfaces that may be provided by the storage mechanism.
 * If not provided, they should return a 'good' result rather than an error
 */

// Reset any errors or intermediate conditions
func (s *Store) Reset() {
	s.lastError = nil
	if reseter, found := s.connection.(Reseter); found {
		reseter.Reset()
	}
	return
}

// Release any locks or memory
func (s *Store) Release() error {
	s.lastError = nil
	if release, found := s.connection.(Releaser); found {
		s.lastError = release.Release()
	}
	return s.lastError
}

// Close the connection to the storage mechanism. If there is no close routine
// ignore the call
func (s *Store) Close() error {
	if s.isOpen != true {
		return s.saveLastError(ErrNotOpen)
	}
	s.isOpen = false
	s.lastError = nil
	if closer, found := s.connection.(Closer); found {
		s.lastError = closer.Close()
	}
	return s.lastError
}

// If implemented, create the basic storage. If not implemented, an error will be returned.
func (s *Store) CreateStore() error {
	if s.isOpen != true {
		return s.saveLastError(ErrNotOpen)
	}

	if creater, found := s.connection.(Creater); found {
		return s.saveLastError(creater.CreateStore())
	}
	return ErrNoSupport
}

func (s *Store) Ping() error {
	s.lastError = nil
	if pinger, found := s.connection.(Pinger); found {
		s.lastError = pinger.Ping()
	}
	return s.lastError
}

/*
 * Mandatory functions
 */
func (s *Store) RegisterUser(user *record.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	s.lastError = s.connection.RegisterUser(user)
	return s.lastError
}

func (s *Store) UserLogin(user *record.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserLogin(user))
}

func (s *Store) UserAuthenticated(user *record.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserAuthenticated(user))
}

func (s *Store) UserLogout(user *record.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserLogout(user))
}

func (s *Store) UserUpdate(user *record.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserUpdate(user))
}

func (s *Store) FetchUserByGuid(guid string) (*record.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.FetchUserByGuid(guid)
	s.lastError = err
	return rec, err
}

func (s *Store) FetchUserByToken(token string) (*record.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.FetchUserByToken(token)
	s.lastError = err
	return rec, err
}

func (s *Store) FetchUserByEmail(email string) (*record.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.FetchUserByEmail(email)
	s.lastError = err
	return rec, err
}

// Fetch the user record by the login string. Returns the record if found and an error code
func (s *Store) FetchUserByLogin(loginName string) (*record.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.FetchUserByLogin(loginName)
	s.lastError = err
	return rec, err
}
