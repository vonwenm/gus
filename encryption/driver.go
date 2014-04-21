// This package simply defines the interfaces that any password generator must define

package encryption

type Driver interface {

	EncryptPassword( password string , salt string) string
	SetInternalSalt( salt string )
	ComparePasswords( hashedPassword , password , salt string ) bool
}

/*
 *			Dynamic interfaces
 */
var driverMap = make( map[string] Driver )
var driverSelect string
const driver_name = "Encryption"

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
		driverSelect = name
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
	if d, found := driverMap[driverSelect] ;  found {
		return d
	}
	panic( driver_name + " driver: '" + driverSelect + "'. Name not found")
}

