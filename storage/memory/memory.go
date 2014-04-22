package memory

import (
	"fmt"
	"github.com/cgentry/gus"
	"github.com/cgentry/gus/storage"
	_ "github.com/mattn/go-sqlite3"
	 "database/sql"
	"strconv"
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

        	MaxSession text,
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
	dbh, err := sql.Open("sqlite3" , ":memory:")
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
func (t * StorageMem) mapToUser(rows *sql.Rows) []gus.User {

	var users []gus.User

	for rows.Next() {
		cols, _ := rows.Columns()
		//user := new(gus.User)

		for name := range cols {
			fmt.Printf("Name is %s\n", name)
		}
		// Copy data over

	}
	return users
}

func (t *StorageMem) FetchUserByGuid(guid string) (   * gus.User , int ) {
	return nil, storage.BANK_USER_NOT_FOUND
}

func ( t *StorageMem) RegisterUser(   user * gus.User )	{
	token, _ := user.GetToken()
	cmd := `INSERT OR IGNORE INTO User
			(Guid,FullName,Email,Domain,LoginName,Password,Token,Salt,
        	IsActive,IsLoggedIn,LoginAt,LogoutAt,LastAuthAt,LastFailedAt,
        	FailCount,MaxSession,TimeoutAt,CreatedAt,UpdatedAt,DeletedAt
        	)
           VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?,?, ?, ?, ?, ?,)`

	t.db.Exec( cmd , user.GetGuid() , user.GetFullName() ,
		user.GetEmail() , user.GetDomain() , user.GetLoginName(),
		user.GetPassword() , token , user.GetSalt() ,
		strconv.FormatBool( user.IsActive ) , strconv.FormatBool( user.IsLoggedIn ),
		user.GetLoginAtStr() , user.GetLoginAtStr() , user.GetLastAuthAt() ,
		user.GetLastFailedAtStr() , user.GetFailCountStr() , user.GetMaxSessionStr() ,
		user.GetTimeoutStr() , user.GetCreatedAtStr() , user.GetDeletedAtStr() )
}

func FetchUserByGuid(  guid string )(   * gus.User , error ){
	cmd := `SELECT * FROM User WHERE Guild = ?`
	q,err := t.db.Query( cmd , guid )
	return nil , err
}

