package memory

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"strconv"
	//"strings"
	"time"
	"errors"
)

const storage_ident = "memory"

type StorageMem struct{
	db            *sql.DB
	lastError    error
}




func init() {
	dbh, err := sql.Open("sqlite3", "/tmp/junk")
	if err != nil {
		panic(err)
	}
	store := &StorageMem{ db: dbh }
	if err := store.CreateStore(); err != nil {
		panic( err )
	}
	storage.Register(storage_ident, store )
}

func (t *StorageMem) GetRawHandle() interface{} {
	return t.db
}

func (t *StorageMem) GetLastError() error {
	return t.lastError
}

func (t * StorageMem) WasLastOk() bool {
	return t.lastError == nil
}

func (t * StorageMem) Open(name, connect string) error {
	t.lastError = nil
	return t.lastError
}


func (t *StorageMem) Close() error {
	t.lastError = t.db.Close()
	return t.lastError
}

func (t *StorageMem) RegisterUser(user * record.User) error {

	cmd := `INSERT OR IGNORE INTO User
			(Guid,     FullName,     Email,
			Domain,    LoginName,    Password,
			Token,     Salt,         IsActive,
			IsLoggedIn,LoginAt,      LogoutAt,
			LastAuthAt,LastFailedAt, FailCount,
			MaxSessionAt,TimeoutAt,    CreatedAt,
			UpdatedAt, DeletedAt
        	)
           VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?,?, ?, ?, ?, ?)`


	_ , err := t.db.Exec(cmd,
		user.GetGuid(), user.GetFullName(), user.GetEmail(),
		user.GetDomain(), user.GetLoginName(), user.GetPassword(),
		user.GetToken(), user.GetSalt(), strconv.FormatBool(user.IsActive),
		strconv.FormatBool(user.IsLoggedIn), user.GetLoginAtStr(), user.GetLogoutAtStr(),
		user.GetLastAuthAtStr(), user.GetLastFailedAtStr(), user.GetFailCountStr(),
		user.GetMaxSessionAtStr(), user.GetTimeoutStr(), user.GetCreatedAtStr(),
		user.GetUpdatedAtStr(), user.GetDeletedAtStr())

	return err
}


func (t *StorageMem) fetchUserByField( field, val string ) ( * record.User, error ){
	cmd := fmt.Sprintf( `SELECT * FROM User WHERE %s = ?` , field )
	rows, err := t.db.Query(cmd, val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := mapColumnsToUser(rows)
	if len( users ) == 0 {
		return nil , errors.New("No records found")
	}
	return users[0], err

}

func (t *StorageMem) FetchUserByToken(token string) (   * record.User , error ) {
	return t.fetchUserByField( "Token" , token )
}

func (t *StorageMem) FetchUserByGuid(guid string) (   * record.User , error ) {
	return t.fetchUserByField( "Guid" , guid )
}

func (t *StorageMem) FetchUserByEmail( email string )( *record.User , error ){
	return t.fetchUserByField( "Email" , email )
}

func (t *StorageMem) FetchUserByLogin( value string )( *record.User , error ){
	return t.fetchUserByField( "LoginName" , value )
}

func StrToTime( t string ) time.Time {
	if val, err := time.Parse( record.USER_TIME_STR , t ); err == nil {
		return val
	}

	return time.Unix(0,0)

}

func StrToBool( t string ) bool {
	if val,err := strconv.ParseBool( t ) ; err == nil {
		return val
	}
	return false
}

func StrToInt( t string ) int {
	if val,err := strconv.ParseInt( t , 10 , 32 ); err == nil {
		if val >= 0 {
			return int(val)
		}
	}
	return 0
}

