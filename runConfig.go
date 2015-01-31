package main

import (
	"github.com/cgentry/gus/command"
	"os"
	"fmt"

)

var cmdConfig = &command.Command{
	Name: 		"config",
	UsageLine: "gus config [-c configfile]",
	Short:     "Build the initial configuration file",
	Long: `
Config will interactivly help you build or edit a configuration file and
save it for you on disk. It will prompt you for each of options with in
a section (for example User Storage and Client Storage).

Strings can be entered free-form or by enclosing them in quotes. If you
use quotes, the contents will be used as entered, with leading and trailing
blanks.

To enter a blank string, enter two quotes (""). If you enter nothing, the
current value will be used.

The default configuration file is stored in ` + DEFAULT_CONFIG_FILENAME +`
but can be overridden by entering the option "-c configfile". The directory
must exist. If the file doesn't exist, defaults will be used when prompting.`,
}

func init(){
	cmdConfig.Run = runConfig
	addCommonCommandFlags(cmdConfig)
}

// Open up the configuration file (if it exists) and then
// prompt the user for all of the configurations values. When complete
// save the file
func runConfig( cmd *command.Command, args []string ){
	command.Box(os.Stdout, "Gus configuration")
	cfile, dirError,_ := GetConfigFileName()
	if dirError != nil  {
		fmt.Println( dirError.Error() )
		return
	}
	if ! command.PromptYesNoDefault(os.Stdout, os.Stdin, "Configure using file " + cfile, false ) {
		fmt.Println("Aborting configuration")
		return
	}

	c,err := GetConfigFile()
	if err != nil {
		c.Default()
	}
	command.PromptForStructFields( &c.Service , config_service_help )
	command.PromptForStructFields( &c.User,config_user_help )

	fmt.Println( c )
	fmt.Println("Config End.")

}

func configInitConfig() error {
	return nil
}

const config_service_help=`
=================================
    Service Configuration
=================================
This sets the general configuration for the running of the program{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`
const config_user_help=`
=================================
    User Database Connection
=================================
Database Connection for User Information
        This sets the configuration for the connection used. There are
        a number of variables, including the type of database you want
        to use for this connection.{{ range . }}
    {{ .Name   }}:
        {{ .Help}}{{ end }}

`
