// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"database/sql"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/storage"
	_ "github.com/mattn/go-sqlite3" // Register sqlite3 with the main system
	"net/http"
)

// Register this driver to the main storage driver with a unique name
func init() {
	storage.Register(STORAGE_IDENTITY, &SqliteDriver{})
}

const STORAGE_IDENTITY = "sqlite"
const DRIVER_IDENTITY = "sqlite3"

// These define all of the fields that are in the database, not in the User record.
const (
	FIELD_GUID           = storage.FIELD_GUID
	FIELD_FULLNAME       = storage.FIELD_NAME
	FIELD_EMAIL          = storage.FIELD_EMAIL
	FIELD_DOMAIN         = `Domain`
	FIELD_LOGINNAME      = storage.FIELD_LOGIN
	FIELD_PASSWORD       = `Password`
	FIELD_TOKEN          = storage.FIELD_TOKEN
	FIELD_SALT           = `Salt`
	FIELD_ISACTIVE       = `IsActive`
	FIELD_ISLOGGEDIN     = `IsLoggedIn`
	FIELD_ISSYSTEM       = `IsSystem`
	FIELD_FAILCOUNT      = `FailCount`
	FIELD_LOGIN_DT       = `LoginAt`
	FIELD_LOGOUT_DT      = `LogoutAt`
	FIELD_LASTAUTH_DT    = `LastAuthAt`
	FIELD_LASTFAILED_DT  = `LastFailedAt`
	FIELD_MAX_SESSION_DT = `MaxSessionAt`
	FIELD_TIMEOUT_DT     = `TimeoutAt`
	FIELD_CREATED_DT     = `CreatedAt`
	FIELD_UPDATED_DT     = `UpdatedAt`
	FIELD_DELETED_DT     = `DeletedAt`
)

type SqliteDriver struct{}

type SqliteConn struct {
	db      *sql.DB
	dsn     string
	options string
}

// Fetch a raw database Sqlite driver
func NewSqliteDriver() SqliteDriver {
	return SqliteDriver{}
}

// Return the raw database handle to the caller. This allows more flexible options
func (t *SqliteConn) GetRawHandle() interface{} {
	return t.db
}

// The main driver will call this function to get a connection to the SqlLite db driver.
// it then 'routes' calls through this connection.
func (t SqliteDriver) Open(dsnConnect string, extraDriverOptions string) (storage.Conn, error) {
	var err error
	store := &SqliteConn{
		dsn:     dsnConnect,
		options: extraDriverOptions,
	}
	store.db, err = sql.Open(DRIVER_IDENTITY, dsnConnect)
	return store, NewGeneralFromError(err, http.StatusInternalServerError)
}

// Close the connection to the database (if it is open)
func (t *SqliteConn) Close() error {
	if t.db == nil {
		return nil
	}
	err := t.db.Close()
	t.db = nil
	return NewGeneralFromError(err, http.StatusInternalServerError)
}
