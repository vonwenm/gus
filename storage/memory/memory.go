// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package memory

import (
	"database/sql"
	"github.com/cgentry/gus/storage"
	_ "github.com/mattn/go-sqlite3"
)

const storage_ident = "memory"

type StorageMem struct {
	engine    string
	db        *sql.DB
	lastError error
}

func init() {

	dbh, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	store := &StorageMem{engine: "sqlite3", db: dbh}

	if err := store.CreateStore(); err != nil {
		panic(err)
	}
	storage.Register(storage_ident, store)
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
	t.db, t.lastError = sql.Open( name , connect )
	return t.lastError
}

func (t *StorageMem) Close() error {
	t.lastError = t.db.Close()
	return t.lastError
}
