package storage

import (
	"fmt"
)

const (
	BANK_USER_OK            = iota
	BANK_USER_NOT_FOUND     = iota
	BANK_USER_TOKEN_INVALID = iota
	BANK_USER_DATA_NOTFOUND = iota
)

/*
 *			Dynamic interfaces
 */
var driverMap = make(map[string]Driver)
var driverSelect string = "unset"

const driver_name = "Storage"

func Register(name string, driver Driver) {
	if driver == nil {
		panic(driver_name + " driver: Register driver is nil")
	}
	if _, dup := driverMap[name]; dup {
		panic(driver_name + " driver: Register called twice for '" + name + "'")
	}
	driverMap[name] = driver

	// First in...first registered
	if len(driverMap) == 1 {
		SetDriver(name)
	}
}

func ToString() string {
	rtn := fmt.Sprintf("Length is %d\n", len(driverMap))
	for key := range driverMap {
		rtn = rtn + key + "\n"
	}
	return rtn
}
func GetDriverName() string {
	return driverSelect
}

func SetDriver(name string) {
	if _, found := driverMap[name]; !found {
		panic(driver_name + " driver: '" + name + "'. Name not found")
	}
	driverSelect = name
}

// GetDriver will return the driver class associated with the curent driver setup
func GetDriver() Driver {
	if d, found := driverMap[driverSelect]; found {
		return d
	}
	panic(driver_name + " driver: '" + driverSelect + "' name not found")
}


