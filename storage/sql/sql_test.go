package sql

import (
	"testing"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/cgentry/gus/storage"
	sysql "database/sql"
)

func startUp(){
	fmt.Println( " Storage is " + storage.GetDriverName() )
	drive := storage.GetDriver()
	drive.Open("sqlite" , ":memory:")

	db := drive.GetRawHandle().( sysql.DB )

	sql := `
        create table User (id integer not null primary key, Fullname text);
        delete from foo;
        `
	_, err := db.Exec(sql)
	if err != nil {
		fmt.Println( err )
	}
}
func TestStart( t *testing.T ){
	fmt.Println( " Storage is " + storage.GetDriverName() )
	startUp()
	fmt.Println("OK")
}
