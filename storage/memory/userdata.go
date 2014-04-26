package memory

import (
	"github.com/cgentry/gus/record"
)
/*
This package stores data per-user and per-session.
 */

func (t *StorageMem)  GetSessionData(    user * record.User , name string ) ( []byte , error ){

	var byteReturn []byte
	cmd := `SELECT Blob FROM UserData WHERE Guid = ? AND Name = ? AND Type="session"`
	err := t.db.QueryRow(cmd, user.GetGuid() , name ).Scan( &byteReturn)
	return byteReturn , err
}

func (t *StorageMem)  SaveSessionData( user * record.User , name string , data *[]byte ) error {

	cmd := `INSERT OR REPLACE INTO UserData ( Data  , Guid , Name , Type ) VALUES( ? ,?,?, "session")  `
	_,err := t.db.Exec(cmd, *data , user.GetGuid() , name  )
	return  err
}
func (t *StorageMem)  DeleteSessionData( user * record.User , name string )  error {
	return nil
}
