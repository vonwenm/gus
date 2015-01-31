package main

import (
	"github.com/cgentry/gus/command"
	"github.com/cgentry/gus/record/configure"
	"path/filepath"
	"io/ioutil"
	"flag"
	"log"
	"os"
	"encoding/json"

)


var configFileName string
var commands = []*command.Command{
	cmdConfig,
}

var usage_template=`Hello. it is working
`
var help_template=`Usage:

          go command [arguments]

{{range .}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "gus help [command]" for more information about a command.

`
func Usage(){
	command.Usage( usage_template, commands )
	os.Exit(0)
}
func main() {
	flag.Usage = Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		Usage()
	}

	if args[0] == "help" {
		command.Help(help_template , "gus", args , commands )
		return
	}

	// Try and run the command given...
	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }

				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()

			cmd.Run(cmd, args)
			return
		}
	}
}

func GetConfigFileName()( Filename string , DirExists error, FileExists error ){
	if configFileName == "" {
		configFileName = DEFAULT_CONFIG_FILENAME
	}
	_,DirExists = os.Stat( filepath.Dir( configFileName ) )
	_,FileExists = os.Stat( configFileName )
	Filename = configFileName
	return
}

func GetConfigFile() ( * configure.Configure , error ){
	var err error
	c := configure.New()

	fname, dirError,fileError := GetConfigFileName()
	if dirError != nil {
		return c, dirError
	}
	if fileError != nil {
		return c, fileError
	}
	fdata,err := ioutil.ReadFile( fname )
	if err == nil {
		err = json.Unmarshal( fdata, c )
	}
	return c, err
}

func SaveConfigFile( jsonString []byte ) error {
	file, direrror, _ := GetConfigFileName()
	if direrror!= nil {
		return direrror
	}
	return ioutil.WriteFile( file, jsonString, DEFAULT_CONFIG_PERMISSIONS )
}



// addCommonCommandFlags will add in flags that are system-wide.
func addCommonCommandFlags( cmd * command.Command ){
	cmd.Flag.StringVar(&configFileName,"c",DEFAULT_CONFIG_FILENAME , "" )
}

