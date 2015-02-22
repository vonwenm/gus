// The encryption drivers are used to encrypt and decrypt passwords stored for a user.
// Standard encryption drivers are sha512, bcrypt and plaintext (not suitable for
// production systems)
//
// All drivers need to call Register in order to be usable by the system. Failure to select
// an encryption driver will cause a panic during runtime.
//
// Drivers are selected by:
//    crypt := encryption.Select( driverName ).Setup( options )
// The call to Setup() is optional and driver specific.

package encryption

import (
	"encoding/json"
	"strings"
)

// The interface gives the set of methods that an encryption driver must implement.
type EncryptDriver interface {
	EncryptPassword(password string, salt string) string
	ComparePasswords(string, string, string) bool
	Setup(string) EncryptDriver
	Id() string
	ShortHelp() string
	LongHelp() string
}

// These are common parameters used by many drivers. Each driver may use structures that are specific to
// that driver.
type CryptOptions struct {
	StaticSalt bool
	Cost       int
	Salt       string
}

// Unmarshal a json string containing the common options defined in CryptOptions and return
// the option structure
func UnmarshalOptions(jsonOption string) (opt *CryptOptions, err error) {
	opt = &CryptOptions{}
	jsonOption = strings.TrimSpace(jsonOption)
	if jsonOption != "" {
		err = json.Unmarshal([]byte(jsonOption), opt)
	}
	return
}

/*
 *			Dynamic interfaces
 */
var driverMap = make(map[string]EncryptDriver)
var driverSelect string

const driver_name = "Encryption"

func GetMap() map[string]EncryptDriver {
	return driverMap
}

// Determine if a driver is registered or not. This encapsulates the map and simply returns a boolean flag.
func IsRegistered(name string) bool {
	_, ok := driverMap[name]
	return ok
}

// Register a new driver in the system.
func Register(driver EncryptDriver) {
	if driver == nil {
		panic(driver_name + " driver: Registered driver is nil")
	}
	name := driver.Id()
	if _, dup := driverMap[name]; dup {
		panic(driver_name + " driver: Register called twice for '" + name + "'")
	}
	driverMap[name] = driver
}

// Pick a registered driver for use in the system. Only one driver can be selected at a time.
func Select(name string) EncryptDriver {
	if driver, found := driverMap[name]; !found {
		panic(driver_name + " driver: '" + name + "'. Name not found")
	} else {
		driverSelect = name
		return driver
	}
}

// GetEncryption will return the driver class associated with the current driver setup
func GetDriver() (driver EncryptDriver) {
	var found bool

	if driverSelect != "" {
		if driver, found = driverMap[driverSelect]; found {
			return
		}
	} else {
		if len(driverMap) == 1 {
			for driverSelect, driver = range driverMap {
				return
			}
		}
	}
	panic(driver_name + " driver: Nothing registered")

}

func GetStaticSalt(offset int) string {
	modIndex := offset % len(encryption_salts)
	return encryption_salts[modIndex]
}
