package main

import (
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/record/configure"
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/encryption"
	"os"
	"errors"
)

const (
	DEFAULT_CMD_USER_LEVEL = "user"
)

var cmdUser = &cli.Command{
	Name:      "user",
	UsageLine: "gus user [add|enable|show|disable] [-c configfile] [-priv level] [-email mail] [-login name] ",
	Short:     "Manipulate users' information in the store system.",
	Long: `
This has three subcommands:
    add         add a new user to the database
    enable      Enable the user account
    disable     Disable the user account, but don't delete it
    show        Display the record that matches the search criteria
The criteria are:
    priv        Select either a normal "user" (default) or "client" systems
    email       Search for records matching the email address.
    login       Search for records matching the user/client login name
`,
}

var cmdUserAdd = &cli.Command{
	Name:      "useradd",
	UsageLine: "gus useradd [-c configfile] [-priv level] [-enable]",
	Short:     "Add a new user to the system.",
	Long: `
Add a new user to the system, specifying the privledge level. This
allows you to bootstrap users into the system. The levels can be:
	user		Normal user (default)
	client		Clients are allowed to remotely authenticate.

To enable the record, you must add -enable or the record will be added
but not enabled.
`,
}

var cmdUserActive = &cli.Command{
	Name:      "useractive",
	UsageLine: "gus useractive [-c configfile] [-priv level] [-group domain] [-enable] [-login name] [-email email address]",
	Short:     "Enable or disable a user",
	Long: `
This will allow you to enable or disable any user in the system. To enable a user, you
must add '-enable' If not, the user will be disabled.

The lookup can be either by login name or by email. Either criteria may be used to look the
user up. Use "-priv client" if you are using separate client/user store and this is a client
lookup. Lookups require the -group flag.
`,
}

var cmdUserCli *record.UserCli

func init() {
	cmdUserCli = record.NewUserCli()

	cmdUser.Run = runUser
	addCommonCommandFlags(cmdUser)
	cmdUser.Flag.StringVar(&cmdUserCli.Level, "priv", DEFAULT_CONFIG_FILENAME, "")
	cmdUser.Flag.StringVar(&cmdUserCli.LoginName, "login" , "","")
	cmdUser.Flag.StringVar(&cmdUserCli.Email, "email" , "","")
	cmdUser.Flag.StringVar(&cmdUserCli.Domain, "group" , "","")

	cmdUserAdd.Run = runUserAdd
	cmdUserActive.Run = runUserActive
}
func runUser(cmd *cli.Command, args []string ){
	var err error
	if len(args) == 0 {
		fmt.Fprintf( os.Stderr, "%s\n", cmd.UsageLine )
		return
	}
	subCommand := args[0]
	cmd.Flag.Parse(args[1:])
	args = cmd.Flag.Args()

	if subCommand != "add" {
		if cmdUserCli.Domain == "" {
			err = errors.New("Domain is required for " + subCommand)
		}else if cmdUserCli.Email == "" && cmdUserCli.LoginName == "" {
			err = errors.New("Email or login is required for " + subCommand)
		}
		if err != nil{
			runtimeFail("Missing parameters" , err )
		}
	}

	switch {
	case subCommand == "add" :
		runUserAdd( cmd , args)
	case subCommand == 	"show" :
		runUserShow(cmd,args)
	case subCommand == "enable" :
		runUserEnable(cmd,args)
	case subCommand == "disable" :
		runUserDisable(cmd, args)
	}
}

func runUserAdd(cmd *cli.Command, args []string) {
	var promptForValues bool
	var configStore configure.Store

	c, err := GetConfigFile()
	if err != nil {
		runtimeFail("Opening configuration file", err)
	}
	encryption.Select( c.Encrypt.Name).Setup(c.Encrypt.Options)

	// We've got the config file. Now we need to prompt for the user information
	for promptForValues = true; promptForValues; {
		cli.PromptForStructFields(cmdUserCli, cmd_useradd_help)
		fmt.Println("\nValues are:")
		cli.PrintStructValue(os.Stdout, cmdUserCli)
		promptForValues = cli.PromptYesNoDefault(os.Stdout, os.Stdin, "Re-enter values", false)
	}
	urec, err := cmdUserCli.NewUser()
	if err != nil {
		runtimeFail("Creating user record from input", err )
	}

	if c.Service.ClientStore && urec.IsSystem {
		configStore =c.Client
	}else{
		configStore = c.User
	}
	store,err := storage.Open(configStore.Name, configStore.Dsn, configStore.Options)
	if err != nil {
		runtimeFail("Opening database", err )
	}

	if err := store.UserInsert(urec) ; err != nil {
		runtimeFail("Writing user record" , err )
	}
	fmt.Fprintf( os.Stdout, "User record created for %s\n", urec.FullName )
	return
}

func runUserActive(cmd *cli.Command, args []string) {
	return
}
func runUserEnable(cmd *cli.Command, args[]string ){
	setUserEnableFlag(true)
	return
}

func runUserDisable(cmd *cli.Command, args[]string){
	setUserEnableFlag(false)
	return
}

func runUserShow( cmd *cli.Command , args[]string ){
	var configStore configure.Store

	c, err := GetConfigFile()
	if err != nil {
		runtimeFail("Opening configuration file", err)
	}
	if c.Service.ClientStore && cmdUserCli.Level == "client" {
		configStore =c.Client
	}else{
		configStore = c.User
	}
	store,err := storage.Open(configStore.Name, configStore.Dsn, configStore.Options)
	defer store.Close()
	if err != nil {
		runtimeFail("Opening database", err )
	}
	userRecord := getUserRecordByCli(store, cmdUserCli)
	cli.RenderTemplate(os.Stdout, const_user_show_template, userRecord )
}

func setUserEnableFlag( newFlag bool ){
	var configStore configure.Store

	c, err := GetConfigFile()
	if err != nil {
		runtimeFail("Opening configuration file", err)
	}
	if c.Service.ClientStore && cmdUserCli.Level == "client" {
		configStore =c.Client
	}else{
		configStore = c.User
	}
	store,err := storage.Open(configStore.Name, configStore.Dsn, configStore.Options)
	defer store.Close()
	if err != nil {
		runtimeFail("Opening database", err )
	}
	userRecord := getUserRecordByCli(store, cmdUserCli)
	if userRecord.IsActive != newFlag {
		userRecord.IsActive = newFlag
		err := store.UserUpdate( userRecord )
		if err != nil {
			runtimeFail("Saving user record", err)
		}
		fmt.Fprintf( os.Stdout, "Record saved for user\n")
	}else{
		fmt.Fprintf( os.Stdout, "No change required for user.")
	}
	fmt.Fprintf( os.Stdout, "Done.")
}

func getUserRecordByCli( store * storage.Store , rec *record.UserCli) ( userRec *record.User) {
	var err error
	if rec.Email != "" {
		userRec,err = store.FetchUserByEmail( rec.Domain, rec.Email )
	}else {
		userRec,err = store.FetchUserByLogin( rec.Domain, rec.LoginName )
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n" , err.Error() )
		os.Exit(1)
	}
	return
}

const cmd_useradd_help = `
=================================
   Add New User
=================================
Add a new user to the system. You can select either a client or user
to add, then you will be prompted for each of the fields.{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`
const const_user_show_template = `

User Record for: {{ .FullName }}

Login ID:        {{ .LoginName }}
Email:           {{ .Email     }}

Is Enabled:      {{ .IsActive  }}
Is Logged In:    {{ .IsLoggedIn}}

Last Login:      {{ .LoginAt }}
Last Auth:       {{ .LastAuthAt }}
Last Logout:     {{ .LogoutAt }}

Created At:      {{ .CreatedAt }}
Updated At:      {{ .UpdatedAt }}
`
