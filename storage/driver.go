package storage

import (
	"fmt"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/tenant"
)

// These are the names of fields we expect to occur in the database and will
// pass to database functions when performing UserFetch operations. You may
// map them in the driver-level routines in order to provide names that are
// more appropriate to the driver mechanism.
const (
	FIELD_EMAIL = `Email`
	FIELD_NAME  = `FullName`
	FIELD_GUID  = `Guid`
	FIELD_LOGIN = `LoginName`
	FIELD_TOKEN = `Token`
)
const driver_name = "Storage"

// For domains, you need to handle the MATCH_ANY_DOMAIN as a parameter. This will be a placeholder
// when there are UNIQUE keys (e.g. GUID and TOKEN)
const MATCH_ANY_DOMAIN = `*`

// All drivers that are registered are stored here.
var driverMap = make(map[string]StorageDriver, 2)

// Find the driver by the 'name' and add it into the map so it can be opened.
// Each driver can only be registered once. To remove all drivers, call
// ResetRegister()
func Register(driver StorageDriver) error {
	if driver == nil {
		panic("Register driver is nil")
	}
	id := driver.Id()
	if _, dup := driverMap[id]; dup {
		return nil
	}
	driverMap[id] = driver

	return nil
}

// Return a simple string respresentation of the drivers that are registered
func String() string {
	rtn := fmt.Sprintf("Length is %d\n", len(driverMap))
	for key := range driverMap {
		rtn = rtn + key + "\n"
	}
	return rtn
}

func GetMap() map[string]StorageDriver {
	return driverMap
}

// Remove all the drivers that have been registered
func ResetRegister() {
	driverMap = make(map[string]StorageDriver, 2)
}

func IsRegistered(name string) bool {
	_, ok := driverMap[name]
	return ok
}

type Storer interface {
	Close() error
	GetStorageConnector() Conn
	LastError() error
	IsOpen() bool
	Ping() error
	Release() error
	Reset()

	FetchUserByEmail(domain, email string)(*tenant.User, error)
	FetchUserByGuid(guid string) (*tenant.User, error)
	FetchUserByLogin(domain, loginName string) (*tenant.User, error)
	FetchUserByToken(token string) (*tenant.User, error)
	UserFetch(domain, lookupKey, lookkupValue string) (*tenant.User, error)
	UserInsert(user *tenant.User) error
	UserUpdate(user *tenant.User) error
}
// Store holds the state for any storage driver. It allows you to have
// consistent returns, such as getting the last error, discovering how
// a connection was made (connectString) or the name of the driver (name)
type Store struct {
	name          string
	connectString string
	isOpen        bool
	lastError     error
	driver        StorageDriver
	connection    Conn
}

// Open a connection to the storage mechanism and return both a storage
// structure and an error status of the open
func Open(name string, connect string, extraDriverOptions string) (*Store, error) {
	s := &Store{
		name:          name,
		isOpen:        false,
		connectString: connect,
		lastError:     ErrNoDriverFound,
	}
	if driver, ok := driverMap[name]; ok {
		s.driver = driver
		s.connection, s.lastError = driver.Open(connect, extraDriverOptions)
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
func (s *Store) LastError() error {
	return s.lastError
}

func ( s *Store )SetLastError( err error ) *Store {
	s.lastError = err
	return s
}

// IsOpen will return the the open status of the connection
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
	s.SetLastError( nil )
	if reseter, found := s.connection.(Reseter); found {
		reseter.Reset()
	}
	return
}

// Release any locks or memory
func (s *Store) Release() error {
	s.SetLastError( nil )
	if release, found := s.connection.(Releaser); found {
		s.SetLastError(  release.Release() )
	}
	return s.LastError()
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

func (s *Store) UserUpdate(user *tenant.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserUpdate(user))
}

func (s *Store) UserInsert(user *tenant.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserInsert(user))
}

// Fetch a user's record using the domain, a field name and the field value. There will only be one record
// returned. If you pass MATCH_ANY_DOMAIN as the domain, this will only be valid for a small number of
// key-types (e.g. enforced unique keys.)
func (s *Store) UserFetch(domain, lookupKey, lookkupValue string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	if domain == MATCH_ANY_DOMAIN {
		if lookupKey != FIELD_GUID || lookupKey != FIELD_TOKEN {
			return nil, ErrMatchAnyNotSupported
		}
	}
	rec, err := s.connection.UserFetch(domain, lookupKey, lookkupValue)
	s.lastError = err
	return rec, err
}

/* ------------------------ THE FOLLOWING ARE 'CONVENIENCE' FUNCTIONS ***********************/

// Fetch a user by the GUID. No domains are required as this is the primary (or unique) key
func (s *Store) FetchUserByGuid(guid string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(MATCH_ANY_DOMAIN, FIELD_GUID, guid)
	s.lastError = err
	return rec, err
}

// Fetch a user by the logged-in token. If the user is not logged in, a 'User not found' error is returned.
func (s *Store) FetchUserByToken(token string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(MATCH_ANY_DOMAIN, FIELD_TOKEN, token)
	s.lastError = err
	return rec, err
}

// Fetch a user by the email. Emails are not unique, except within a domain.
func (s *Store) FetchUserByEmail(domain, email string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(domain, FIELD_EMAIL, email)
	s.lastError = err
	return rec, err
}

// Fetch the user record by the login string. Login names are only unique within the domain
func (s *Store) FetchUserByLogin(domain, loginName string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(domain, FIELD_LOGIN, loginName)
	s.lastError = err
	return rec, err
}
