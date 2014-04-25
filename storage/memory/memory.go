package memory

import (
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"strconv"
	"strings"
	"time"
	"errors"
)

const storage_ident = "memory"

type StorageMem struct{
	db            *sql.DB
	lastError    error
}

func loadSql(dbh * sql.DB) {
	sql := `
        create table User (
        	Guid text primary key,
        	FullName text,
        	Email    text,


        	Domain    text,
        	LoginName text,
        	Password  text,
        	Token     text,

        	Salt text,

        	IsActive   integer,
        	IsLoggedIn integer,

        	LoginAt      text,
        	LogoutAt     text,
        	LastAuthAt   text,
        	LastFailedAt text,
        	FailCount    integer ,

        	MaxSessionAt text,
        	TimeoutAt  text,

        	CreatedAt text,
        	UpdatedAt text,
        	DeletedAt text);
        create unique index uemail on User(Email);
        create unique index ulname on User(LoginName);
        create index ufullname on User(FullName);
        create index utoken on User(Token);
        `
	_, err := dbh.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
}


func init() {
	dbh, err := sql.Open("sqlite3", "/tmp/junk")
	if err != nil {
		panic(err)
	}
	loadSql(dbh)
	storage.Register(storage_ident, &StorageMem{ db: dbh })
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

// Convert db data to user structure. For this, we expect a 1:1 mapping
func (t * StorageMem) mapToUser(rows *sql.Rows) []record.User {

	var users []record.User

	for rows.Next() {
		cols, _ := rows.Columns()
		//user := new(record.User)

		for name := range cols {
			fmt.Printf("Name is %s\n", name)
		}
		// Copy data over

	}
	return users
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

func mapColumnsToUser(rows * sql.Rows) []*record.User {

	var allUsers [] *record.User
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	vpoint := make([]interface{}, count)
	var vstr string

	for rows.Next() {
		for i, _ := range columns {
			vpoint[i] = &values[i]
		}
		user := record.NewUser("")
		rows.Scan(vpoint...)

		for i , col := range columns {
			val := values[i]
			if b, ok := val.([]byte) ; ok {
				vstr = string(b)


			switch strings.ToLower(col) {
				case "fullname" : user.SetName( vstr )
				case "email"	: user.SetEmail( vstr )
				case "guid"     : user.SetGuid( vstr )

				case "domain"	: user.SetDomain( vstr )
				case "password" : user.SetPasswordStr( vstr )
				case "token"    : user.SetToken( vstr )

				case "salt"	    : user.SetSalt( vstr )
				case "isactive"	: user.SetIsActive( StrToBool(vstr ) )
				case "isloggedin" : user.SetIsLoggedIn( StrToBool( vstr) )

				case "loginat"	: user.SetLoginAt( StrToTime(  vstr ) )
				case "logoutat" : user.SetLogoutAt( StrToTime(  vstr ) )
				case "lastfailedat" : user.SetLastFailedAt( StrToTime(vstr) )
				case "failcount" : user.SetFailCount( StrToInt(vstr) )

				case "maxsessionat": user.SetMaxSessionAt( StrToTime(vstr ) )
				case "timeoutat" : user.SetTimeoutAt( StrToTime(  vstr ) )

				case "createdat" : user.SetCreatedAt( StrToTime(  vstr ) )
				case "updatedat" : user.SetUpdatedAt( StrToTime(  vstr ) )
				case "deletedat" : user.SetDeletedAt( StrToTime(  vstr ) )
				case "loginname" : user.SetLoginName( vstr )

			}
			}
		} // End columns

		allUsers = append(allUsers, user)
	}
	return allUsers
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

