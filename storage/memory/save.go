package memory

import (
	"strconv"
	"github.com/cgentry/gus/record"
	//"database/sql"
)

func (t *StorageMem ) UpdateUserRecordWithCondition( user * record.User , condition , condvalue string ) error {
	cmd := `UPDATE User
			(     FullName,     Email,
			Domain,    LoginName,    Password,
			Token,     Salt,         IsActive,
			IsLoggedIn,LoginAt,      LogoutAt,
			LastAuthAt,LastFailedAt, FailCount,
			MaxSessionAt,TimeoutAt,    CreatedAt,
			UpdatedAt, DeletedAt
        	)
           VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?,?, ?, ?, ?, ?)
           WHERE Guid = ? `
		if condition != "" {
			cmd = cmd + " AND " + condition + " ? "
		}else{
			cmd = cmd + " AND Guid = "
			condvalue = user.GetGuid()			// Duplicate key request should be ignored
		}


	_ , err := t.db.Exec(cmd,
		user.GetFullName(), user.GetEmail(),
		user.GetDomain(), user.GetLoginName(), user.GetPassword(),
		user.GetToken(), user.GetSalt(), strconv.FormatBool(user.IsActive),
		strconv.FormatBool(user.IsLoggedIn), user.GetLoginAtStr(), user.GetLogoutAtStr(),
		user.GetLastAuthAtStr(), user.GetLastFailedAtStr(), user.GetFailCountStr(),
		user.GetMaxSessionAtStr(), user.GetTimeoutStr(), user.GetCreatedAtStr(),
		user.GetUpdatedAtStr(), user.GetDeletedAtStr(), user.GetGuid() , condvalue )

	return err
}
// Save the user's record when they login. (They must be logged off)
func (t *StorageMem) SaveUserLogin(  user * record.User ) error	{
	return t.UpdateUserRecordWithCondition( user , "IsLoggedIn = " , "false")
	return nil
}

// SaveUserAuth will save fields (with conditions)
func (t *StorageMem) SaveUserAuth(   user * record.User ) error	{
	return t.UpdateUserRecordWithCondition( user , "IsLoggedIn = " , "true")
}
// SaveUserLogoff will save the record (with conditions)
func (t *StorageMem) SaveUserLogoff( user * record.User)	error	{
	return t.UpdateUserRecordWithCondition( user , "IsLoggedIn = " , "true" )
	return nil
}
