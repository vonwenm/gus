package memory

import (
	"database/sql"
	"sync"
	"strings"
	//"fmt"
)

/* The fields that we save. Defined as constants so I don't have to retype them */

const (

	DB_FIELD_COUNT_UPDATE = 20
	DB_FIELD_LIST_UPDATE = `
            FullName,     Email,
			Domain,       LoginName,    Password,
			Token,        Salt,         IsActive,
			IsLoggedIn,   LoginAt,      LogoutAt,
			LastAuthAt,   LastFailedAt, FailCount,
			MaxSessionAt, TimeoutAt,    CreatedAt,
			UpdatedAt,    DeletedAt,    IsSystem`

	DB_FIELD_COUNT_STATUS = 13
	DB_FIELD_LIST_STATUS = `
			IsActive,	  IsLoggedIn,   LoginAt,
			LogoutAt,	  LastAuthAt,   LastFailedAt,
			FailCount,    MaxSessionAt, TimeoutAt,
			CreatedAt,    UpdatedAt,    DeletedAt,
			IsSystem`

	DB_FIELD_COUNT_ALL = DB_FIELD_COUNT_UPDATE + 1
	DB_FIELD_LIST_ALL = `Guid, ` + DB_FIELD_LIST_UPDATE

)



const (
	MAP_REGISTER_INSERT = iota
	MAP_REGISTER_CHECK  = iota
)



var stmtMap = make( map[int] * sql.Stmt	,10 )		// Where we store prepared statments
var mutex sync.Mutex								 // For concurrency

func (t *StorageMem) GetRegisterSql() (* sql.Stmt , error) {
	sql := `INSERT INTO User (` + DB_FIELD_LIST_ALL + `)
		    VALUES (` + strings.Repeat( "?, ",  DB_FIELD_COUNT_ALL - 1) + `? )`
	return t.lockAndLoad( &mutex, MAP_REGISTER_INSERT , sql )
}

func (t *StorageMem ) GetRegisterChecksSql() (  * sql.Stmt , error ){
	sql := `SELECT Guid, Domain , Email, LoginName
			FROM User
			WHERE Guid = ? OR ( Domain = ? AND ( Email = ? OR LoginName = ?))`
	return t.lockAndLoad( &mutex, MAP_REGISTER_CHECK , sql )
}

func (t *StorageMem ) lockAndLoad( mutex *sync.Mutex, tag int , sql string) ( * sql.Stmt , error ){

	mutex.Lock()
	defer mutex.Unlock()
	stmt, found := stmtMap[ tag  ]
	if  found {
		return stmt, nil
	}

	stmt, err := t.db.Prepare(sql)
	if err == nil {
		stmtMap[tag] = stmt
	}
	return stmt, err
}
