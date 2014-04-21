package storage

import (
	"github.com/CGentry/gus"
)

const (
	BANK_USER_OK 			= iota
	BANK_USER_NOT_FOUND		= iota
	BANK_USER_TOKEN_INVALID = iota
	BANK_USER_DATA_NOTFOUND = iota
)

/*
 *			Dynamic interfaces
 */
var driverMap = make( map[string] EncryptDriver )
var driverSelect string
const driver_name = "Storage"

func Register( name string , driver Driver ){
	if driver == nil {
		panic(driver_name + " driver: Register driver is nil")
	}
	if _,dup := driverMap[name] ; dup {
		panic( driver_name + " driver: Register called twice for '" + name + "'")
	}
	driverMap[name] = driver

	// First in...first registered
	if len( driverMap ) == 1 {
		driverMap = name
	}

}

func SetDriver( name string ){
	if _ , found := driverMap[name] ; ! found {
		panic( driver_name + " driver: '" + name + "'. Name not found")
	}
	driverSelect = name
}

// GetEncryption will return the driver class associated with the curent driver setup
func GetDriver( ) ( Driver  ){
	if d, found := driverMap[name] ;  found {
		return d
	}
	panic( driver_name + " driver: '" + name + "'. Name not found")
}


// The banker interface defines very general, high level operations for retrieval and storage of
// data. The back-storage can be a flat file, database or document store.
// The interfaces specify NO sql methods and flatten out operations
type Driver interface {
	FetchUserByGuid(  guid string )(   * model.User , int )
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

