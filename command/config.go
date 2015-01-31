package command

// This set of routines encapsulates the command structure, which holds
// the help text and definitions for a gus command.
//
// See License for copyright information.
// Portions of this code are (c) 2012 by the Go Authors. This code has been
// been adapted for use in GUS and is not the sole work of Charles Gentry.
//
import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

func Usage(usageTemplate string, commands []*Command) {
	printUsage(os.Stderr, usageTemplate, commands)
	os.Exit(2)
}

// help implements the 'help' command.
func Help(usageTemplate , cmd string, args []string, commands []*Command) {
	if len(args) == 1 {
		// General help
		printUsage(os.Stdout, usageTemplate, commands)
		return
	}
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s help command\n\nToo many arguments given.\n", cmd , args[0])
		os.Exit(2) // failed at 'go help'
	}

	// gus help cmd
	for _, cmd := range commands {
		if cmd.Name == args[1] {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run '%s %s'.\n", args[1] , cmd , args[0])
	os.Exit(2) // failed at 'go help cmd'
}

func printUsage(w io.Writer, usageTemplate string, commands []*Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
}

// A command is any subcommand, such as createstore
// (from go source code)
type Command struct {
	Name string

	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var helpTemplate = `{{if .Runnable}}usage: {{.UsageLine}}

{{end}}{{.Long | trim}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func ParseBool(txt string) (rtn bool, err error) {
	answer := strings.ToLower(txt)
	if answer == `y` || answer == `yes` || answer == `ok` {
		return true, nil
	}
	if answer == `n` || answer == `no` {
		return false, nil
	}
	rtn, err = strconv.ParseBool(txt)
	return
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

type HelpList struct {
	Name	string
	Help	string
}
func PromptForStructFields(substruct interface{}, helpTemplate string ) {
	val := reflect.Indirect(reflect.ValueOf(substruct))

	helpList := make(map[int]HelpList)

	for i := 0; i < val.NumField(); i++ {
		help := val.Type().Field(i).Tag.Get("help")
		if help != "" {
			name := val.Type().Field(i).Tag.Get("name")
			if name != "" {
				helpList[i] = HelpList{Name:name, Help:help}
			}
		}
	}
	tmpl(os.Stdout, helpTemplate, helpList )
	//val := reflect.ValueOf(substruct).Elem()
	for i := 0; i < val.NumField(); i++ {
		hostEntry, ok := helpList[i]
		if ! ok {
			continue
		}
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		switch typeField.Type {
		default:
			fmt.Println("Don't know what this is. Skipping.")
		case bool:
			PromptYesNoDefault( os.Stdout, os.Stdin, helpEntry.Name, valueField )
		case string:
			PromptStringDefault(os.stdout, os.Stdin, helpEntry.Name, valueField )
		}
				//fmt.Printf("%2d: Field: %s (%s) Current value: %v\n", i+1, name, valueField.Type().Name(), valueField.Interface())

	}
}
func promptString(w io.Writer, r io.Reader, prompt string) (string, error) {
	var txt string

	fmt.Fprintf(w, "%s ", prompt)
	_, err := fmt.Fscanln(r, &txt)

	if err != nil && err.Error() == `unexpected newline` {
		return "", nil
	}
	if err == nil {

		txt = strings.TrimSpace(txt)
		if (strings.HasPrefix(txt, `"`) && strings.HasPrefix(txt, `"`)) ||
			(strings.HasPrefix(txt, `'`) && strings.HasPrefix(txt, `'`)) {
			return txt[1 : len(txt)-1], err
		}
	}
	return txt, err
}
func PromptStringDefault(w io.Writer, r io.Reader, prompt, defaultValue string) string {
	var txt string

	fmt.Fprintf(w, "%s? ", prompt)
	_, err := fmt.Fscanln(r, &txt)

	if err == nil {

		txt = strings.TrimSpace(txt)
		if (strings.HasPrefix(txt, `"`) && strings.HasPrefix(txt, `"`)) ||
			(strings.HasPrefix(txt, `'`) && strings.HasPrefix(txt, `'`)) {
			return txt[1 : len(txt)-1]
		}
		return txt
	}

	return defaultValue
}

func PromptString(w io.Writer, r io.Reader, prompt string) (string, error) {

	for true {
		txt, err := promptString(w, r, prompt+"? ")

		if err != nil {
			return "", err
		}

		if txt != "" {
			return txt, nil
		}

		fmt.Fprintf(w, "\tInvalid string (use \"\" for empty string) '%s'\n", txt)

	}
	return "", nil
}

// Simple prompt for yes or no values. Caller can pass buffers or os.Stdin/os.Stdout
//
func PromptYesNo(w io.Writer, r io.Reader, prompt string) (bool, error) {

	for true {
		txt, err := promptString(w, r, prompt+" (yes/no)? ")
		if err != nil {
			return false, err
		}
		if txt != "" {
			if val, err := ParseBool(txt); err == nil {
				return val, nil
			}
		}
		fmt.Fprintf(w, "\tInvalid yes/no response '%s'\n", txt)
	}
	return false, nil
}

func PromptYesNoDefault(w io.Writer, r io.Reader, prompt string, defaultValue bool) bool {
	var txt string

	if defaultValue {
		fullPrompt := fmt.Sprintf("%s\n(yes/no, enter = %s)", prompt, "yes")
		txt = PromptStringDefault(w, r, fullPrompt, "yes")
	} else {
		fullPrompt := fmt.Sprintf("%s\n(yes/no, enter = %s)", prompt, "no")
		txt = PromptStringDefault(w, r, fullPrompt, "no")
	}
	if val, err := ParseBool(txt); err == nil {
		return val
	}
	return defaultValue
}

func Box( w io.Writer , info string ){
	fmt.Fprintf(w , "****************************\n%s\n************************\n", info )
}
