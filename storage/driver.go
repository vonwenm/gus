package storage

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"net/http"
)

const (
	BANK_USER_OK            = iota
	BANK_USER_NOT_FOUND     = iota
	BANK_USER_TOKEN_INVALID = iota
	BANK_USER_DATA_NOTFOUND = iota
)

type StorageError struct {
	errorString string
	errorCode   int
}

func NewStorageError(msg string, code int) *StorageError {
	return &StorageError{errorString: msg, errorCode: code}
}

func NewStorageFromError(e error, code int) *StorageError {
	if e == nil {
		return nil
	}
	return &StorageError{errorString: e.Error(), errorCode: code}

}

func (s *StorageError) Error() string {
	return s.errorString
}
func (s *StorageError) Code() int {
	return s.errorCode
}

/*
 *			Dynamic interfaces
 */
var driverMap = make(map[string]Driver)

var ErrInvalidHeader = NewStorageError("Invalid header in request", http.StatusBadRequest)
var ErrInvalidChecksum = NewStorageError("Invalid Checksum", http.StatusBadRequest)
var ErrInvalidBody = NewStorageError("Invalid body (mistmatch request?)", http.StatusBadRequest)
var ErrEmptyFieldForLookup = NewStorageError("Lookup field is empty", http.StatusBadRequest)

var ErrNoDriverFound = NewStorageError("No storage driver found", http.StatusInternalServerError)
var ErrNoSupport = NewStorageError("Storage driver does not support function call", http.StatusNotImplemented)
var ErrNotOpen = NewStorageError("Storage driver is not open", http.StatusInternalServerError)
var ErrAlreadyRegistered = NewStorageError("Storage driver already registered", http.StatusInternalServerError)
var ErrInternalDatabase = NewStorageError("Internal storage error while executing operation", http.StatusInternalServerError)

var ErrUserNotFound = NewStorageError("User not found", http.StatusNotFound)

var ErrInvalidGuid = NewStorageError("Invalid Guid for lookup", http.StatusNotFound)
var ErrInvalidEmail = NewStorageError("Invalid email for lookup", http.StatusNotFound)
var ErrInvalidToken = NewStorageError("Invalid token for lookup", http.StatusNotFound)

var ErrDuplicateGuid = NewStorageError("User GUID already in use", http.StatusInternalServerError)
var ErrDuplicateEmail = NewStorageError("Email registered", http.StatusConflict)
var ErrDuplicateLogin = NewStorageError("Login name already exists", http.StatusConflict)

var ErrUserNotRegistered = NewStorageError("User not registered", http.StatusBadRequest)
var ErrUserNotLoggedIn = NewStorageError("User not logged in", http.StatusBadRequest)
var ErrUserLoggedIn = NewStorageError("User already logged in", http.StatusBadRequest)
var ErrUserNotActive = NewStorageError("User is not yet activated", http.StatusUnauthorized)

var ErrStatusOk = NewStorageError("", http.StatusOK)

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

func ToString() string {
	rtn := fmt.Sprintf("Length is %d\n", len(driverMap))
	for key := range driverMap {
		rtn = rtn + key + "\n"
	}
	return rtn
}

func ResetRegister() {
	driverMap = make(map[string]Driver)
}

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
