// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	"database/sql"
	"github.com/cgentry/gus/storage"
	_ "github.com/mattn/go-sqlite3" // Register sqlite3 with the main system
)

const STORAGE_IDENTITY = "sqlite"
const DRIVER_IDENTITY = "sqlite3"

type SqliteDriver struct{}

type SqliteConn struct {
	db  *sql.DB
	dsn string
}

// Register this driver to the main storage driver with a unique name
func init() {
	storage.Register(STORAGE_IDENTITY, &SqliteDriver{})
}

func (t *SqliteConn) GetRawHandle() interface{} {
	return t.db
}

func (t SqliteDriver) Open(dsnConnect string) (storage.Conn, error) {
	var err error
	store := &SqliteConn{
		dsn: dsnConnect,
	}
	store.db, err = sql.Open(DRIVER_IDENTITY, dsnConnect)
	return store, err
}

func (t *SqliteConn) Close() error {
	if t.db == nil {
		return nil
	}
	err := t.db.Close()
	t.db = nil
	return err
}
