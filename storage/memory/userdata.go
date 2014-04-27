// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package memory

import (
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
)

const (
	GROUP_SESSION = "session" // Transient, session data
	GROUP_USER    = "user"    // Permanent, user-level data
)

// This file handles all of the user data storage, both session and long-term
// It performs simple upserts/selects from the database.
//
// Note that in SQLITE3 an INSERT OR REPLACE will cause a record to be deleted, then inserted
//
// Get...Data retrieves the current named permanent record for a user. The validity of the
// session is not checked.

/* ---------------------------------------------------------
 *                 ACCESS SESSION (TEMPORARY) DATA
 * ---------------------------------------------------------
 */
func (t *StorageMem) GetSessionData(user *record.User, name string) (string, storage.HeaderMap, error) {
	return t.getData(user, name, GROUP_SESSION)
}

func (t *StorageMem) SaveSessionData(user *record.User, name string, data *string, headers *storage.HeaderMap) error {
	return t.saveData(user, name, data, headers, GROUP_SESSION)
}

// DeleteSessionData removes the named session record from the database
func (t *StorageMem) DeleteSessionData(user *record.User, name string) error {
	return t.deleteData(user, name, GROUP_SESSION)
}

/* ---------------------------------------------------------
 *                 ACCESS USER (PERMANENT) DATA
 * ---------------------------------------------------------
 */

func (t *StorageMem) GetUserData(user *record.User, name string) (string, storage.HeaderMap, error) {
	return t.getData(user, name, GROUP_USER)
}

func (t *StorageMem) SaveUserData(user *record.User, name string, data *string, headers *storage.HeaderMap) error {
	return t.saveData(user, name, data, headers, GROUP_USER)
}

func (t *StorageMem) DeleteUserData(user *record.User, name string) error {
	return t.deleteData(user, name, GROUP_USER)
}

/* ---------------------------------------------------------
 * Low level worker functions that handle the real delete/read/inserts.
 * This allows specific 'typing' but simply maps to a different record
 * ---------------------------------------------------------
 */

func (t *StorageMem) getData(user *record.User, name, storageType string) (string, storage.HeaderMap, error) {

	var dataReturn string
	var tempString string

	cmd := `SELECT Data , Headers FROM UserData WHERE Guid = ? AND Name = ? AND Type=?`
	row := t.db.QueryRow(cmd, user.GetGuid(), name, storageType)

	err := row.Scan(&dataReturn, &tempString)

	return dataReturn, storage.NewHeaderMap(tempString), err
}

func (t *StorageMem) saveData(user *record.User, name string, data *string, headers *storage.HeaderMap, storageType string) error {

	cmd := `INSERT OR REPLACE INTO UserData ( Data  , Headers, Guid , Name , Type ) VALUES( ? ,?,?,?,? )  `
	_, err := t.db.Exec(cmd, *data, headers.ToString(), user.GetGuid(), name, storageType)
	return err
}

func (t *StorageMem) deleteData(user *record.User, name, storageType string) error {
	cmd := `DELETE FROM UserData WHERE Guid = ? AND Name = ? AND Type = ?`
	_, err := t.db.Exec(cmd, user.GetGuid(), name, storageType)
	return err
}
