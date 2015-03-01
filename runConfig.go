package main

import (
	"encoding/json"
	"fmt"
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/storage"
	"os"
	"strings"
)

const (
	CONFIG_SETUP_AUTOSALT = "##salt##"
)

var cmdConfig = &cli.Command{
	Name:      "config",
	UsageLine: "gus config [-c configfile] [list]",
	Short:     "Build the initial configuration file",
	Long: `
Config will interactivly help you build, edit or list a configuration file and
save it for you on disk. It will prompt you for each of options with in
a section (for example User Storage and Client Storage).

To see what the values are, use:
    gus config [-c configfile] list

Strings can be entered free-form or by enclosing them in quotes. If you
use quotes, the contents will be used as entered, with leading and trailing
blanks.

To enter a blank string, enter two quotes (""). If you enter nothing, the
current value will be used.

The default configuration file is stored in ` + DEFAULT_CONFIG_FILENAME + `
but can be overridden by entering the option "-c configfile". The directory
must exist. If the file doesn't exist, defaults will be used when prompting.
`,
}
var cmdCreateStore = &cli.Command{
	Name:      "createstore",
	UsageLine: "gus createstore [-c configfile]",
	Short:     "Initialise the store",
	Long: `
Initialise the user and client stores. This will be a non-destructive
operation. If the stores exist already, nothing should occur.
`,
}

func init() {
	cmdConfig.Run = runConfig
	addCommonCommandFlags(cmdConfig)

	cmdCreateStore.Run = runConfigStore
	addCommonCommandFlags(cmdCreateStore)
}

// Open up the configuration file (if it exists) and then
// prompt the user for all of the configurations values. When complete
// save the file
func runConfig(cmd *cli.Command, args []string) {
	var promptForValues bool

	if len(args) > 0 && args[0] == "list" {
		runConfigList()
		return
	}

	cli.Box(os.Stdout, "GUS configuration")
	cfile, dirError, _ := GetConfigFileName()
	if dirError != nil {
		fmt.Println(dirError.Error())
		return
	}
	if !cli.PromptYesNoDefault(os.Stdout, os.Stdin, "Configure using file "+cfile, false) {
		fmt.Println("Aborting configuration")
		return
	}

	c, err := GetConfigFile()
	if err != nil {
		c.Default()
	}
	for promptForValues = true; promptForValues; {
		cli.PromptForStructFields(&c.Service, template_cmd_help_config)
		fmt.Println("\nValues are:")
		cli.PrintStructValue(os.Stdout, &c.Service)
		promptForValues = cli.PromptYesNoDefault(os.Stdout, os.Stdin, "Re-enter values", false)
	}

	for promptForValues = true; promptForValues; {
		cli.PromptForStructFields(&c.Encrypt, template_cmd_help_config_crypt)
		if strings.Contains(c.Encrypt.Options, CONFIG_SETUP_AUTOSALT) {
			c.Encrypt.Options = strings.Replace(c.Encrypt.Options, CONFIG_SETUP_AUTOSALT, tenant.CreateSalt(200), -1)
		}
		fmt.Println("\nValues are:")
		cli.PrintStructValue(os.Stdout, &c.Encrypt)
		promptForValues = cli.PromptYesNoDefault(os.Stdout, os.Stdin, "Re-enter values", false)
	}
	for promptForValues = true; promptForValues; {
		cli.PromptForStructFields(&c.User, template_cmd_help_config_user)
		fmt.Println("\nValues are:")
		cli.PrintStructValue(os.Stdout, &c.User)
		promptForValues = cli.PromptYesNoDefault(os.Stdout, os.Stdin, "Re-enter values", false)
	}
	if c.Service.ClientStore {
		for promptForValues = true; promptForValues; {
			cli.PromptForStructFields(&c.Client, template_cmd_help_config_client)
			fmt.Println("\nValues are:")
			cli.PrintStructValue(os.Stdout, &c.Client)
			promptForValues = cli.PromptYesNoDefault(os.Stdout, os.Stdin, "Re-enter values", false)
		}
	}
	if cli.PromptYesNoDefault(os.Stdout, os.Stdin, "Ok to save configuration values", true) {
		cdata, err := json.MarshalIndent(c, "", "  ")
		if err == nil {
			err = SaveConfigFile(cdata)
		}
		if err == nil {
			fmt.Println("\nConfiguration saved")
		} else {
			fmt.Printf("\nError when saving the configuration: %s\n", err.Error())
		}
	} else {
		fmt.Println("\nConfiguration not saved")
	}
	return
}

// Open the configuration and then call the 'create store' routines
// prompt the user for all of the configurations values. When complete
// save the file
func runConfigStore(cmd *cli.Command, args []string) {
	c, err := GetConfigFile()
	if err != nil {
		runtimeFail("Opening configuration file", err)
	}
	userStore, err := storage.Open(c.User.Name, c.User.Dsn, c.User.Options)
	if err != nil {
		runtimeFail("Opening user store", err)
	}
	defer userStore.Close()
	if err = userStore.CreateStore(); err != nil {
		runtimeFail("Creating user store", err)
	}

	if c.Service.ClientStore {
		clientStore, err := storage.Open(c.Client.Name, c.Client.Dsn, c.Client.Options)
		if err != nil {
			runtimeFail("Opening client store", err)
		}
		defer clientStore.Close()
		if err = clientStore.CreateStore(); err != nil {
			runtimeFail("Creating client store", err)
		}
	}
	fmt.Fprintf(os.Stdout, "Stores created\n")
}

func configInitConfig() error {
	return nil
}

func runConfigList() {
	cli.Box(os.Stdout, "Gus Configuration file")
	c, err := GetConfigFile()
	if err != nil {
		fmt.Println("Configuration error: " + err.Error())
		return
	}
	fmt.Println("\n")
	cli.Box(os.Stdout, "Service Configuration")
	cli.PrintStructValue(os.Stdout, &c.Service)
	fmt.Println("\n")

	fmt.Println("\n")
	cli.Box(os.Stdout, "Service Configuration")
	cli.PrintStructValue(os.Stdout, &c.Encrypt)
	fmt.Println("\n")

	cli.Box(os.Stdout, "User Storage Configuration")
	cli.PrintStructValue(os.Stdout, &c.User)
	fmt.Println("\n")

	if c.Service.ClientStore {
		cli.Box(os.Stdout, "Client Storage Configuration")
		cli.PrintStructValue(os.Stdout, &c.Client)
	}
	fmt.Println("\n")
}

const template_cmd_help_config = `
=================================
    Service Configuration
=================================
This sets the general configuration for the running of the program{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`
const template_cmd_help_config_user = `
=================================
    User Storage Connection
=================================
Connection for User Information
        This sets the configuration for the connection used. There are
        a number of variables, including the type of database you want
        to use for this connection.{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`

const template_cmd_help_config_client = `
=================================
    Client Storage Connection
=================================
Connection for Client Information
        This sets the configuration for the connection used. There are
        a number of variables, including the type of database you want
        to use for this connection.{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`
const template_cmd_help_config_crypt = `
=================================
    Password Encryption Driver
=================================
What encryption technology to use.
        This sets the configuration for the encryption driver to use. The
        standard selection is bcrypt. The options are usually JSON
        encoded and would include additional data for use with randomising
        the password. This would be encoded as:
        { "salt" : "data to include, usually random" }{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`
