// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package mock

import (
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/storage"
)

// Register this driver to the main storage driver with a unique name
func init() {
	storage.Register(NewMockDriver())
}

const STORAGE_IDENTITY = "mock"
const DRIVER_IDENTITY = "mock"

type MockDriver struct {
	Name  string
	Short string
	Long  string
}

type MockConn struct {
	db      map[string]*tenant.User
	errList map[string]error
}

// Fetch a raw database Mock driver
func NewMockDriver() *MockDriver {
	return &MockDriver{
		Name:  STORAGE_IDENTITY,
		Short: "Mock storage driver",
		Long:  const_mock_help_template,
	}
}

// The main driver will call this function to get a connection to the SqlLite db driver.
// it then 'routes' calls through this connection.
func (t *MockDriver) Open(option1 string, extraDriverOptions string) (storage.Conn, error) {
	store := &MockConn{}
	store.db = make(map[string]*tenant.User)
	store.errList = make(map[string]error)
	return store, nil
}
func (t *MockDriver) Id() string        { return t.Name }
func (t *MockDriver) ShortHelp() string { return t.Short }
func (t *MockDriver) LongHelp() string  { return t.Long }
func (t *MockDriver) Usage() string     { return t.Short }

// Return the raw database handle to the caller. This allows more flexible options
func (t *MockConn) GetRawHandle() interface{} {
	return t.db
}

// Close the connection to the database (if it is open)
func (t *MockConn) Close() error {
	return nil
}

const const_mock_help_template = `

   This is a dummy driver used for testing purposes.

   `

func (t *MockConn) UserUpdate(user *tenant.User) error {
	if err, ok := t.errList[user.Guid]; ok {
		return err
	}
	t.db[user.Guid] = user
	return nil
}
func (t *MockConn) UserInsert(user *tenant.User) error {
	if err, ok := t.errList[user.Guid]; ok {
		return err
	}
	t.db[user.Guid] = user
	return nil
}

func (t *MockConn) UserFetch(domain, key, value string) (*tenant.User, error) {
	found := false
	for _, user := range t.db {

		if domain == storage.MATCH_ANY_DOMAIN || domain == user.Domain {
			switch key {
			case storage.FIELD_GUID:
				found = (value == user.Guid)
			case storage.FIELD_EMAIL:
				found = (value == user.Email)
			case storage.FIELD_LOGIN:
				found = (value == user.LoginName)
			case storage.FIELD_TOKEN:
				found = (value == user.Token)
			}
			if found {
				if err, ok := t.errList[user.Guid]; ok {
					return nil, err
				}
				return user, nil
			}
		}
	}
	return nil, ErrUserNotFound
}
