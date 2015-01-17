package sqlite

// The fields that we save. Defined as constants so I don't have to retype them

const (
	DB_FIELD_COUNT_UPDATE = 20
	DB_FIELD_LIST_UPDATE  = `
            FullName,     Email,
			Domain,       LoginName,    Password,
			Token,        Salt,         IsActive,
			IsLoggedIn,   LoginAt,      LogoutAt,
			LastAuthAt,   LastFailedAt, FailCount,
			MaxSessionAt, TimeoutAt,    CreatedAt,
			UpdatedAt,    DeletedAt,    IsSystem`

	DB_FIELD_COUNT_STATUS = 13
	DB_FIELD_LIST_STATUS  = `
			IsActive,	  IsLoggedIn,   LoginAt,
			LogoutAt,	  LastAuthAt,   LastFailedAt,
			FailCount,    MaxSessionAt, TimeoutAt,
			CreatedAt,    UpdatedAt,    DeletedAt,
			IsSystem`

	DB_FIELD_COUNT_ALL = DB_FIELD_COUNT_UPDATE + 1
	DB_FIELD_LIST_ALL  = `Guid, ` + DB_FIELD_LIST_UPDATE
)

const (
	FIELD_GUID           = `Guid`
	FIELD_FULLNAME       = `FullName`
	FIELD_EMAIL          = `Email`
	FIELD_DOMAIN         = `Domain`
	FIELD_LOGINNAME      = `LoginName`
	FIELD_PASSWORD       = `Password`
	FIELD_TOKEN          = `Token`
	FIELD_SALT           = `Salt`
	FIELD_ISACTIVE       = `IsActive`
	FIELD_ISLOGGEDIN     = `IsLoggedIn`
	FIELD_LOGIN_DT       = `LoginAt`
	FIELD_LOGOUT_DT      = `LogoutAt`
	FIELD_LASTAUTH_DT    = `LastAuthAt`
	FIELD_LASTFAILED_DT  = `LastFailedAt`
	FIELD_FAILCOUNT      = `FailCount`
	FIELD_MAX_SESSION_DT = `MaxSessionAt`
	FIELD_TIMEOUT_DT     = `TimeoutAt`
	FIELD_CREATED_DT     = `CreatedAt`
	FIELD_UPDATED_DT     = `UpdatedAt`
	FIELD_DELETED_DT     = `DeletedAt`
	FIELD_ISSYSTEM       = `IsSystem`
)

const (
	MAP_REGISTER_INSERT = iota
	MAP_REGISTER_CHECK  = iota
)

/*

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
*/
