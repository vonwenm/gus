package sql

import (
	"testing"
	"fmt"
	"github.com/cgentry/gus/storage"
	_ "github.com/mattn/go-sqlite3"
	//sysql "database/sql"
)

//func startUp(){
//	fmt.Println( " Storage is " + storage.GetDriverName() )
//	fmt.Println( storage.ToString() )
//	drive := storage.GetDriver()
//	drive.Open("sqlite" , ":memory:")
//
//	db := drive.GetRawHandle().( sysql.DB )
//
//	sql := `
//        create table User (id integer not null primary key, Fullname text);
//        delete from foo;
//        `
//	_, err := db.Exec(sql)
//	if err != nil {
//		fmt.Println( err )
//	}
//}
func TestStart( t *testing.T ){
	fmt.Println( " Storage is " + storage.GetDriverName() )
	//startUp()
	fmt.Println("OK")
}
