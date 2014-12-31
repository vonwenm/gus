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

type StorageMem struct {
	engine    string
	db        *sql.DB
	lastError error
	openName  string
	openConn  string
}

// Register this driver to the main storage driver with a unique name
func init() {
	store := &StorageMem{engine: "sqlite3"}
	storage.Register(STORAGE_IDENTITY, store)
}

func (t *StorageMem) GetRawHandle() interface{} {
	return t.db
}

func (t *StorageMem) GetLastError() error {
	return t.lastError
}

func (t *StorageMem) WasLastOk() bool {
	return t.lastError == nil
}

func (t *StorageMem) Open(name, connect string) error {
	t.openName = name
	t.openConn = connect
	t.db, t.lastError = sql.Open( name , connect )
	return t.lastError
}

func (t *StorageMem) Close() error {
	if t.db == nil {
		t.lastError = nil
	}else{
		t.lastError = t.db.Close()
	}
	return t.lastError
}
