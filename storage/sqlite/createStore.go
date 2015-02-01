// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package sqlite

import (
	. "github.com/cgentry/gus/ecode"
	"net/http"
)

// CreateStore is a non-destructive storage creation mechanism. It can be called on the cli line
// with the option -C
func (t *SqliteConn) CreateStore() error {

	sql := []string{
		`CREATE TABLE IF NOT EXISTS User (
			Guid         text primary key,
			LoginName    text ,
			Email        text ,
			Token        text UNIQUE,

			Salt         text,

			FullName     text,
			Domain       text,
			Password     text,

			IsActive     integer,
			IsLoggedIn   integer,
			IsSystem     integer,

			LoginAt      text,
			LogoutAt     text,
			LastAuthAt   text,
			LastFailedAt text,
			FailCount    integer ,

			MaxSessionAt text,
			TimeoutAt    text,

			MaxSessionAtSec int8,
			TimeoutAtSec    int8,

			CreatedAt    text,
			UpdatedAt    text,
			DeletedAt    text);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idxlogin      ON User(LoginName,Domain)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idxEmail      ON User(Email,Domain)`,
		`CREATE        INDEX IF NOT EXISTS idxfullname   ON User(FullName);`,
		`CREATE        INDEX IF NOT EXISTS idxMaxSession ON User(MaxSessionAt);`,
		`CREATE        INDEX IF NOT EXISTS idxTimeoutAt  ON User(TimeoutAt);`,
	}

	for _, cmd := range sql {
		if _, err := t.db.Exec(cmd); err != nil {
			return NewGeneralFromError(err, http.StatusInternalServerError)
		}
	}

	return nil
}
