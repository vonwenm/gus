package sql

import (
	"database/sql"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"fmt"
)

type StorageSql struct{
	db 			*sql.DB
	lastError	error
}

const storage_ident = "sql"

func init(){
	storage.Register( storage_ident , &StorageSql{})
}


func ( t *StorageSql ) GetRawHandle() interface{}{
	return t.db
}

func ( t *StorageSql ) GetLastError() error {
	return t.lastError
}

func ( t * StorageSql ) WasLastOk() bool {
	return t.lastError == nil
}

func ( t * StorageSql ) Open( name, connect string ) error {
	t.db , t.lastError = sql.Open( name , connect )
	return t.lastError
}


func ( t *StorageSql ) Close() error {
	t.lastError = t.db.Close()
	return t.lastError
}

// Convert db data to user structure. For this, we expect a 1:1 mapping
func ( t * StorageSql ) mapToUser( rows *sql.Rows )  []record.User {

	var users []record.User

	for rows.Next() {
		cols,_ := rows.Columns()
		//user := new(record.User)

		for name := range cols {
			fmt.Printf( "Name is %s\n" , name)
		}
		// Copy data over

	}
	return users
}

func (t *StorageSql ) FetchUserByGuid(  guid string )(   * record.User , int ){
	return nil , storage.BANK_USER_NOT_FOUND
}
