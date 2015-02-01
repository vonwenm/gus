package main

import (
	"fmt"
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/encryption"
	"github.com/cgentry/gus/storage"
	"os"
)

var helpStore = &cli.Command{
	Name:      "store",
	UsageLine: "gus store [driver-name]",
	Short:     "Display a list of what drivers are available",
	Long: `
Display all of the drivers that are compiled into this runtime. If
you add in the 'driver-name', it will list specific help for that driver.

Each driver may require different paramters. The driver will give you some
details, but you should refer to the documentation
`,
}
var helpEncrypt = &cli.Command{
	Name:      "encrypt",
	UsageLine: "gus encrypt [driver-name]",
	Short:     "Display a list of what drivers are available",
	Long: `
Display all of the drivers that are compiled into this runtime. If
you add in the 'driver-name', it will list specific help for that driver.

Each driver may require different paramters. The driver will give you some
details, but you should refer to the documentation
`,
}

func init() {
	helpStore.Run = runStore
	helpEncrypt.Run = runEncrypt
}

// Output any help that is required
func runStore(cmd *cli.Command, args []string) {
	listStore := storage.GetMap()

	if len(args) == 0 {
		cli.RenderTemplate(os.Stdout, template_storage_list, listStore)
		return
	}
	if len(args) == 1 {
		if entry, ok := listStore[args[0]]; ok {
			cli.RenderTemplate(os.Stdout, template_storage_entry, entry)
			return
		}
		fmt.Fprintf(os.Stderr, "'%s' is not a valid storage driver\n", args[0])
	} else {
		fmt.Fprintf(os.Stderr, "Only one parameter for store command\nUse 'gus help store' for more information\n")
	}

	return
}

const template_storage_list = `
List of storage drivers available:{{ range . }}
  {{ .Id }}: {{ .ShortHelp }}{{ end }}

`
const template_storage_entry = `
{{ .Id }}: {{ .ShortHelp }}
{{ .LongHelp }}
`

func runEncrypt(cmd *cli.Command, args []string) {
	listStore := encryption.GetMap()

	if len(args) == 0 {
		cli.RenderTemplate(os.Stdout, template_encryption_list, listStore)
		return
	}
	if len(args) == 1 {
		if entry, ok := listStore[args[0]]; ok {
			cli.RenderTemplate(os.Stdout, template_encryption_entry, entry)
			return
		}
		fmt.Fprintf(os.Stderr, "'%s' is not a valid encryption driver\n", args[0])
	} else {
		fmt.Fprintf(os.Stderr, "Only one parameter for encrypt command\nUse 'gus help encrypt' for more information\n")
	}

	return
}

const template_encryption_list = `
List of encryption drivers available:{{ range . }}
  {{ .Id }}: {{ .ShortHelp }}{{ end }}

`
const template_encryption_entry = `
{{ .Id }}: {{ .ShortHelp }}
{{ .LongHelp }}
`
