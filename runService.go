package main

import (
	"encoding/json"
	"fmt"
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"os"
	"strings"
)

var cmdService = &cli.Command{
	Name:      "service",
	UsageLine: "gus service [-c configfile]",
	Short:     "Run the program in service mode.",
	Long: `
Service will listen in on a port and wait for requests for user activity from
clients. Clients will call to register, authenticate, login and logout from
the system. Each request is made over HTTP but must use a PUT instead of a GET.

The single option, "-c" allows you to specify where to load the configuration
file from. The default configuration file is ` + DEFAULT_CONFIG_FILENAME + `.

`,
}

func init() {
	cmdService.Run = runService
	addCommonCommandFlags(cmdService)

	c, err := GetConfigFile()
	if err != nil {
		runtimeFail( "Service failed" , err )
	}

	encryption.Select( c.Encrypt.Name).Setup(c.Encrypt.Options)
}
