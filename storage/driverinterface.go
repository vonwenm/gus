package storage

import (
	"github.com/cgentry/gus/record"
)

// The driver interface defines very general, high level operations for retrieval and storage of
// data. The back-storage can be a flat file, database or document store.
// The interfaces specify NO sql methods and flatten out operations
type Driver interface {
	Open(connect string) (Conn, error)
}

// This is the minimum call set that every driver is required to implement
type Conn interface {
	RegisterUser(user *record.User) error

	UserLogin(user *record.User) error
	UserAuthenticated(user *record.User) error
	UserLogout(user *record.User) error

	UserUpdate(user *record.User) error

	FetchUserByGuid(guid string) (*record.User, error)
	FetchUserByToken(token string) (*record.User, error)
	FetchUserByEmail(email string) (*record.User, error)
	FetchUserByLogin(loginName string) (*record.User, error)
}

// Option Storge Creation interface
type Creater interface {
	CreateStore() error
}

//Optional Closing interface. If this isn't implemented, no error is reported.
type Closer interface {
	Close() error
}

//Optional Reset interface. This will reset any errors and cleanup any intermediate results
type Reseter interface {
	Reset()
}

// Optional database 'ping' interface. This will check the database connection
type Pinger interface {
	Ping() error
}

//Optional Release interface. This will release any locks/resources that a driver may have set
//For example, the MySQL will do a SELECT...FOR UPDATE for all of the FetchXXX calls. The
//release will cause an explicit commit. This, in the code, will be called by a 'defer' call after
//any fetch/insert operation. For other drivers, it can be ignored or perform any other operation
//required.
// Note that SQLITE doesn't do anything at this stage as it isn't really considered a robust, fully
// hardened storage mechanism. Document-style interfaces will probably not use it either.
type Releaser interface {
	Release() error
}
