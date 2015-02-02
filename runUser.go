package main

import (
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/record/configure"
	"fmt"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/encryption"
	"os"
)

const (
	DEFAULT_CMD_USER_LEVEL = "user"
)

var cmdUser = &cli.Command{
	Name:      "user",
	UsageLine: "gus user [add|enable|show|disable] [-c configfile] [-priv level] ",
	Short:     "Manipulate users' information in the store system.",
	Long: `
Add a new user to the system, specifying the privledge level. This
allows you to bootstrap users into the system. The levels can be:
	user		Normal user (default)
	client		Clients are allowed to remotely authenticate.

To enable the record, you must add -enable or the record will be added
but not enabled.
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
	UsageLine: "gus useractive [-c configfile] [-priv level] [-enable] [-login name] [-email email address]",
	Short:     "Enable or disable a user",
	Long: `
This will allow you to enable or disable any user in the system. To enable a user, you
must add '-enable' If not, the user will be disabled.

The lookup can be either by login name or by email. Either criteria may be used to look the
user up. Use "-priv client" if you are using separate client/user store and this is a client
lookup.
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

	cmdUserAdd.Run = runUserAdd
	addCommonCommandFlags(cmdUserAdd)
	cmdUserAdd.Flag.StringVar(&cmdUserCli.Level, "priv", DEFAULT_CONFIG_FILENAME, "")
	cmdUserAdd.Flag.BoolVar(&cmdUserCli.Enable, "enable", false, "")

	cmdUserActive.Run = runUserActive
	addCommonCommandFlags(cmdUserActive)
	cmdUserActive.Flag.StringVar(&cmdUserCli.Level, "priv", DEFAULT_CONFIG_FILENAME, "")
	cmdUserActive.Flag.BoolVar(&cmdUserCli.Enable, "enable", false, "")
	cmdUserActive.Flag.StringVar(&cmdUserCli.LoginName, "login" , "","")
	cmdUserActive.Flag.StringVar(&cmdUserCli.Email, "email" , "","")
}
func runUser(cmd *cli.Command, args []string ){
	if len(args) == 0 {
		fmt.Fprintf( os.Stderr, "%s\n", cmd.UsageLine )
		return
	}
	switch {
	case args[0] == "add" :
		runUserAdd( cmd , args[1:])
	case args[0] == "show" :
		fmt.Println("Show!", args)
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


const cmd_useradd_help = `
=================================
   Add New User
=================================
Add a new user to the system. You can select either a client or user
to add, then you will be prompted for each of the fields.{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`
