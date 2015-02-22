package tenant

// UserCli contains tags used for prompting
type UserCli struct {
	FullName  string `name:"User's full name"     help:"The full user's name (title, first and last) of the user."`
	LoginName string `name:"User's login id"      help:"This is what the user would use to identify themselves to the system."`
	Email     string `name:"User's email address" help:"The user's real email address, if available."`
	Domain    string `name:"User's group"         help:"What group, or domain, does this user belong to."`
	Password  string `name:"Password"             help:"Password for user"`

	Level  string // Level for the user. From the flags set.
	Enable bool   `name:"Enable"               help:"Enable user record"`
}

func NewUserCli() *UserCli {
	return &UserCli{}
}
