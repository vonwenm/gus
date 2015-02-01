package main

import (
	"encoding/json"
	"flag"
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/record/configure"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"runtime/debug"
)

var configFileName string
var commands = []*cli.Command{
	cmdConfig,
	cmdCreateStore,
	cmdUserAdd,
	cmdUserActive,
	helpStore,
	helpEncrypt,
}

var help_template = `Usage:

          go command [arguments]

{{range .}}
    {{.Id | printf "%-11s"}} {{ .ShortHelp }}{{end}}

Use "gus help [command]" for more information about a cli.

`

func Usage() {
	cli.Usage(help_template, commands)
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
		cli.Help(help_template, "gus", args, commands)
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

func GetConfigFileName() (Filename string, DirExists error, FileExists error) {
	if configFileName == "" {
		configFileName = DEFAULT_CONFIG_FILENAME
	}
	_, DirExists = os.Stat(filepath.Dir(configFileName))
	_, FileExists = os.Stat(configFileName)
	Filename = configFileName
	return
}

func GetConfigFile() (*configure.Configure, error) {
	var err error
	c := configure.New()

	fname, dirError, fileError := GetConfigFileName()
	if dirError != nil {
		return c, dirError
	}
	if fileError != nil {
		return c, fileError
	}
	fdata, err := ioutil.ReadFile(fname)
	if err == nil {
		err = json.Unmarshal(fdata, c)
	}

	return c, err
}

func SaveConfigFile(jsonString []byte) error {
	file, direrror, _ := GetConfigFileName()
	if direrror != nil {
		return direrror
	}
	return ioutil.WriteFile(file, jsonString, DEFAULT_CONFIG_PERMISSIONS)
}

// addCommonCommandFlags will add in flags that are system-wide.
func addCommonCommandFlags(cmd *cli.Command) {
	cmd.Flag.StringVar(&configFileName, "c", DEFAULT_CONFIG_FILENAME, "")
}

func runtimeFail(msg string,err error) {
	var rpt int
	emsg := err.Error()
	if  len( emsg) > len(msg ){
		rpt = len(emsg)
	}else{
		rpt = len(msg)
	}

	stars := strings.Repeat("*", rpt+4)
	fmt.Fprintf(os.Stderr, "%s\n* %-*s *\n* %-*s *\n%s\n\n", stars, rpt,msg, rpt, emsg, stars)
	fmt.Fprintln(os.Stderr, "STACK TRACE:")
	debug.PrintStack()
	fmt.Fprintln(os.Stderr,"\n")
	os.Exit(1)
}
