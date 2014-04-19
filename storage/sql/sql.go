package sql

import (
	"database/sql"
	"github.com/CGentry/gus"
	"github.com/CGentry/gus/storage"
)

type StorageSql struct{
	db 			*sql.DB
	lastError	error
}


func init(){
	storage.RegisterStorage("sql" , &StorageSql{})
}

func ( t *StorageSql ) GetLastError() error {
	return t.lastError
}

func ( t * StorageSql ) WasLastOk() bool {
	return t.lastError == nil
}

func  OpenSql( name, connect string ) ( * StorageSql , error ) {
	s := new( StorageSql )
	s.db , s.lastError = sql.Open( name , connect )
	return s, s.lastError
}


func ( t *StorageSql ) Close() error {
	t.lastError = t.db.Close()
	return t.lastError
}

// Convert db data to user structure. For this, we expect a 1:1 mapping
func ( t * StorageSql ) mapToUser( rows *sql.Rows )  []gus.User {
	var users []gus.User
	for rows.Next() {
		user := new(gus.User)
		
		users = append( users, *user )

		// Copy data over

	}
	return users
}

func (t *StorageSql ) FetchUserByGuid(  guid string )(   * gus.User , int ){
	return nil , storage.BANK_USER_NOT_FOUND
}
