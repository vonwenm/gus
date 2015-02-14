// Package logit handles the setup and logging for the program. Like other driver-based packages,
// you select the logger you want and pass it configuration information. Then you simply log
// messages and they will be formatted by the driver

package logit

import (
	"strings"
)

// The interface gives the set of methods that an encryption driver must implement.
type LogitDriver interface {
	Open() LogitDriver
	Write(level int, logval ...interface{} )
	Close()

	Id() string
	ShortHelp() string
	LongHelp() string
}

/*
 *			Dynamic interfaces
 */
var driverMap = make(map[string]LogitDriver)
var driverSelect string

const driver_name = "Logit"

func GetMap() map[string]LogitDriver {
	return driverMap
}
// Determine if a driver is registered or not. This encapsulates the map and simply returns a boolean flag.
func IsRegistered(name string) bool {
	_, ok := driverMap[name]
	return ok
}

// Register a new driver in the system.
func Register(driver LogitDriver) {
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
func Select(name string) LogitDriver {
	if driver, found := driverMap[name]; !found {
		panic(driver_name + " driver: '" + name + "'. Name not found")
	}
	if driverSelect != "" {
		d := GetDriver()
		d.Close()
	}
	driverSelect = name
	return driver

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

