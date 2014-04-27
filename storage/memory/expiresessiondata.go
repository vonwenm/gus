package memory
/*
import (
	"time"
)
func (t *StorageMem) ExpireSessionData() {

	go func (){
		tx, err := t.db.Begin()
		if err == nil {
			now := time.Now().Unix()
			cmd := `DELETE FROM UserData WHERE Type = GROUP_SESSION AND Guid IN
					(
						SELECT Guid
						FROM User
						WHERE IsLogged = "true"
						  AND MaxSessionAtSec <= ?
						  AND MaxTimeoutAtSec <= ?
					)`
			clearLogin := `UPDATE User Set IsLoggedIn="false"
						WHERE IsLoggedIn = "true"
						  AND MaxSessionAtSec <= ?
						  AND MaxTimeoutAtSec <= ?`
			tx.Commit()
		}
	}()
}
*/
