package main

import (
	"github.com/cgentry/gus/command"
)

var cmdConfig = &command.Command{
	Name: "config",
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

The default configuration file is stored in ` + DEFAULT_CONFIG_LOCATION +`
but can be overridden by entering the option "-c configfile". The directory
must exist. If the file doesn't exist, defaults will be used when prompting.`
}

func init(){
	cmdConfig.Run = runConfig
	cmd.Flag.String(&configFileName,"c","")
}

// Open up the configuration file (if it exists) and then
// prompt the user for all of the configurations values. When complete
// save the file
func runConfig( cmd *Command, args []string ){

}
