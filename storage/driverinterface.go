package storage

import (
	"github.com/cgentry/gus/record"
)

// The driver interface defines very general, high level operations for retrieval and storage of
// data. The back-storage can be a flat file, database or document store.
// The interfaces specify NO sql methods and flatten out operations
type Driver interface {
	Open(name, connect string) error
	CreateStore() error
	Close() error

	RegisterUser(user *record.User) error // Save initial routine

	SaveUserLogin(user *record.User) error  // Save relevant data for being logged in
	SaveUserAuth(user *record.User) error   // User just did an authentication - save info
	SaveUserLogoff(user *record.User) error // User just logged off

	FetchUserByGuid(guid string) (*record.User, error)
	FetchUserByToken(token string) (*record.User, error)
	FetchUserByEmail(email string) (*record.User, error)
	FetchUserByLogin(login string) (*record.User, error)

	GetLastError() error

}
