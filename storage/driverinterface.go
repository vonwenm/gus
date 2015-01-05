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
