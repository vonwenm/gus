package storage

import (
	. "github.com/CGentry/gus"
)

const (
	BANK_USER_OK 			= iota
	BANK_USER_NOT_FOUND		= iota
	BANK_USER_TOKEN_INVALID = iota
	BANK_USER_DATA_NOTFOUND = iota
)

var driverStore = make( map[string] Driver )
var driverSelect string


// The banker interface defines very general, high level operations for retrieval and storage of
// data. The back-storage can be a flat file, database or document store.
// The interfaces specify NO sql methods and flatten out operations
type Driver interface {
	FetchUserByGuid(  guid string )(   * User , int )
	/*
	FetchUserByToken( token string )(  * User , int )
	FetchUserByEmail( email string )(  * User , int )
	FetchUserByLogin( login string )(  * User , int )

	ExpireSessionRecords()				// The implementer might want to do this in a go routine

	RegisterUser(   user * User )	 	// Save initial routine

	SaveUserLogin(  user *User ) error	// Save relevant data for being logged in
	SaveUserAuth(   user *User ) error	// User just did an authentication - save info
	SaveUserLogoff( user *User)	error	// User just logged off

	GetSessionData(    user * User , name string ) ( []byte , int )
	SaveSessionData(   user * User , name string , data *[]byte ) int
	DeleteSessionData( user * User , name string )  int


	GetUserData(     user * User ,  name string )( []byte , int )
	SaveUserData(    user * User ,  name string , data *[]byte ) int
	DeleteUserData(  user * User ,  name string ) int
	*/

}

// Register storage function for later use. The first storage mechanism will be the default
// Multiple mechanisms may be stored but only one will be used
func RegisterStorage( name string , driver Driver ){
	if driver == nil {
		panic("Storage driver: Register driver is nil")
	}

	if _,dup := driverStore[name] ; dup {
		panic( "Storage driver: Register called twice for '" + name + "'")
	}

	driverStore[name] = driver
	if driverSelect == "" || len( driverStore ) == 1  {
		driverSelect = name
	}
}

// Fetch the storage mechanism that is active
func GetStorageDriver() Driver {
	if driverSelect == "" || len( driverStore ) == 0  {
		panic( "Storage driver: Nothing selected")
	}
	return driverStore[driverSelect]
}
