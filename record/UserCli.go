package record

// UserReturn will contain the USER's minimum data from any operation operation.
type UserCli struct {
	FullName  string `name:"User's full name"     help:"The full user's name (title, first and last) of the user."`
	LoginName string `name:"User's login id"      help:"This is what the user would use to identify themselves to the system."`
	Email     string `name:"User's email address" help:"The user's real email address, if available."`
	Domain    string `name:"User's group"         help:"What group, or domain, does this user belong to."`
	Password  string `name:"Password"             help:"Password for user"`

	Level string // Level for the user. From the flags set.
	Enable	  bool	 `name:"Enable"               help:"Enable user record"`

}

func NewUserCli() *UserCli {
	return &UserCli{}
}

/**
 * Create a record from the user record passed to this routine
 * See:		UserReturn
 */
func (r *UserCli) NewUser() (rtn *User, err error) {

	rtn = NewUser()
	if err = rtn.SetDomain(r.Domain); err != nil {
		return
	}
	if err = rtn.SetName(r.FullName); err != nil {
		return
	}
	if err = rtn.SetEmail(r.Email); err != nil {
		return
	}
	if err = rtn.SetLoginName(r.LoginName); err != nil {
		return
	}
	if err = rtn.SetPassword(r.Password); err != nil {
		return
	}
	rtn.SetIsActive( r.Enable )
	err = rtn.SetIsSystem(r.Level == "client")

	return

}
