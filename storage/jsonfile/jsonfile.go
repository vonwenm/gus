// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package jsonfile

/*
import (
	"encoding/json"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/storage"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Register this driver to the main storage driver with a unique name
func init() {
	storage.Register(NewJsonFileDriver())
}

const STORAGE_IDENTITY = "jsonfile"
const DRIVER_IDENTITY = "jsonfile"

type JsonFileDriver struct {
	Name  string
	Short string
	Long  string
}

type JsonFileConn struct {
	filename string
	filemod  time.Time
	filesize int
	isdirty  bool
	isMonitor bool
	userlist []tenant.User
}

// Fetch a raw database JsonFile driver
func NewJsonFileDriver() *JsonFileDriver {
	return &JsonFileDriver{
		Name:  STORAGE_IDENTITY,
		Short: "JsonFile storage driver",
		Long:  const_jsonfile_help_template,
	}
}

// The main driver will call this function to get a connection to the SqlLite db driver.
// it then 'routes' calls through this connection.
func (t *JsonFileDriver) Open(jsonfile string, extraDriverOptions string) (storage.Conn, error) {
	store := &JsonFileConn{filename: jsonfile, isdirty: false, isMonitor: bool}
	go store.Monitor() // Start the MONITOR in the background
	return store, nil
}

// Monitor will check the file for changes every minute. This will change to an notify routine once the
// functions are merged from experimental to main. This routine should work with any OS (rather than just one
// or two)
func (t *JsonFileConn) Monitor() {
	var err error
	var buff []byte
	var stat os.FileInfo

	for t.isMonitor {
		// Check to see if we need to write out the data to the file
		if t.isdirty {
			t.WriteLock()
			var fp *os.File
			if fp, err = os.Create(t.filename); err == nil {
				if buff, err := json.MarshalIndent(t.userlist, "", "  "); err == nil {
					fp.Write(buff)
					fp.Close()
					if stat, err = os.Stat(t.filename); err == nil {
						t.filemod = stat.ModTime()
					}
				}
				t.ReleaseLock()
				t.isdirty = false
			}
			//
			// ELSE We need to check the file to see if it is different
		} else {
			if stat, err = os.Stat(t.filename); err != nil {
				if os.IsNotExist(err) {
					t.isdirty = true
					err = nil
				}
			} else if stat.ModTime().After(t.filemod) {
				t.WriteLock()
				if buff, err = ioutil.ReadFile(t.filename); err == nil {
					if err = json.Unmarshal(buff, t.userlist); err != nil {
						log.Printf("JsonFileDriver.Monitor: %s", err.Error())
					}
				}
				t.ReleaseLock()
			}
		}
		if err != nil {
			log.Printf("JsonFileDriver.Monitor: %s", err.Error())
		}
		time.Sleep(60 * time.Second) // One minute sleep period
	}
}

func (t *JsonFileConn) ReadLock() {

}
func (t *JsonFileConn) WriteLock() {

}
func (t *JsonFileConn) ReleaseLock() {

}
func (t *JsonFileDriver) Id() string        { return t.Name }
func (t *JsonFileDriver) ShortHelp() string { return t.Short }
func (t *JsonFileDriver) LongHelp() string  { return t.Long }
func (t *JsonFileDriver) Usage() string     { return t.Short }

// Return the raw database handle to the caller. This allows more flexible options
func (t *JsonFileConn) GetRawHandle() interface{} {
	return t.db
}

// Close the connection to the database (if it is open)
func (t *JsonFileConn) Close() error {
	return nil
}

const const_jsonfile_help_template = `

   This is a store that uses JSON formatted strings and stores the data in a file.

   `

func (t *JsonFileConn) UserUpdate(user *tenant.User) error {
	if err, ok := t.errList[user.Guid]; ok {
		return err
	}
	t.db[user.Guid] = user
	return nil
}
func (t *JsonFileConn) UserInsert(user *tenant.User) error {
	if err, ok := t.errList[user.Guid]; ok {
		return err
	}
	t.db[user.Guid] = user
	return nil
}

func (t *JsonFileConn) UserFetch(domain, key, value string) (*tenant.User, error) {
	found := false
	for _, user := range t.db {

		if domain == storage.MATCH_ANY_DOMAIN || domain == user.Domain {
			switch key {
			case storage.FIELD_GUID:
				found = (value == user.Guid)
			case storage.FIELD_EMAIL:
				found = (value == user.Email)
			case storage.FIELD_LOGIN:
				found = (value == user.LoginName)
			case storage.FIELD_TOKEN:
				found = (value == user.Token)
			}
			if found {
				if err, ok := t.errList[user.Guid]; ok {
					return nil, err
				}
				return user, nil
			}
		}
	}
	return nil, ErrUserNotFound
}
*/
