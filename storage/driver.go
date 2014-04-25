package storage

import (
	"github.com/cgentry/gus/record"
	"fmt"
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
var driverMap   = make( map[string] Driver )
var driverSelect string = "unset"

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
		SetDriver( name )
	}
}

func ToString() string {
	rtn := fmt.Sprintf("Length is %d\n" , len( driverMap))
	for key := range driverMap {
		rtn = rtn + key + "\n"
	}
	return rtn
}
func GetDriverName() string {
	return driverSelect
}

func SetDriver( name string ){
	if _ , found := driverMap[name] ; ! found {
		panic( driver_name + " driver: '" + name + "'. Name not found")
	}
	driverSelect = name
}

// GetEncryption will return the driver class associated with the curent driver setup
func GetDriver( ) ( Driver  ){
	if d, found := driverMap[driverSelect] ;  found {
		return d
	}
	panic( driver_name + " driver: '" + driverSelect + "' name not found")
}


// The banker interface defines very general, high level operations for retrieval and storage of
// data. The back-storage can be a flat file, database or document store.
// The interfaces specify NO sql methods and flatten out operations
type Driver interface {

	Open( name, connect string ) error
	Close() error

	RegisterUser(     user * record.User )	 error 	// Save initial routine

	FetchUserByGuid(  guid string )(   * record.User , error )
	FetchUserByToken( token string )(  * record.User , error )

	FetchUserByEmail( email string )(  * record.User , error )
	FetchUserByLogin( login string )(  * record.User , error )

	/*
	ExpireSessionRecords()				// The implementer might want to do this in a go routine

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

