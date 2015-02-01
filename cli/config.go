package cli

// This set of routines encapsulates the cli structure, which holds
// the help text and definitions for a gus cli.
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

// help implements the 'help' cli.
func Help(usageTemplate, cmd string, args []string, commands []*Command) {
	if len(args) == 1 {
		// General help
		printUsage(os.Stdout, usageTemplate, commands)
		return
	}
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s help command\n\nToo many arguments given.\n", cmd, args[0])
		os.Exit(2) // failed at 'go help'
	}

	// gus help cmd
	for _, cmd := range commands {
		if cmd.Id() == args[1] {
			RenderTemplate(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run '%s %s'.\n", args[1], cmd, args[0])
	os.Exit(2) // failed at 'go help cmd'
}

func printUsage(w io.Writer, usageTemplate string, commands []*Command) {
	bw := bufio.NewWriter(w)
	RenderTemplate(bw, usageTemplate, commands)
	bw.Flush()
}

type Helper interface {
	Id() string
	ShortHelp() string
	LongHelp() string
}

// A cli is any subcommand, such as createstore
// (from go source code)
type Command struct {
	Name string

	// Run runs the cli.
	// The args are the arguments after the cli name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-cli>' output.
	Long string

	// Flag is a set of flags specific to this cli.
	Flag flag.FlagSet

	// CustomFlags indicates that the cli will do its own
	// flag parsing.
	CustomFlags bool
}

func (t *Command) Id() string        { return t.Name }
func (t *Command) ShortHelp() string { return t.Short }
func (t *Command) LongHelp() string  { return t.Long }

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Runnable reports whether the cli can be run; otherwise
// it is a documentation pseudo-cli such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var helpTemplate = `{{if .Runnable}}usage: {{.UsageLine}}

{{end}}{{.Long | trim}}
`

// RenderTemplate executes the given template text on data, writing the result to w.
func RenderTemplate(w io.Writer, text string, data interface{}) {
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

type StructNameList struct {
	Name string
	Help string
}

func PromptForStructFields(substruct interface{}, helpTemplate string) {
	val := reflect.Indirect(reflect.ValueOf(substruct))

	nameList, _ := getStructList(substruct)

	RenderTemplate(os.Stdout, helpTemplate, nameList)
	//val := reflect.ValueOf(substruct).Elem()
	for i := 0; i < val.NumField(); i++ {
		helpEntry, ok := nameList[i]
		if !ok {
			continue
		}
		//valueField := val.Field(i)
		typeField := val.Type().Field(i)

		switch typeField.Type.String() {
		default:

		case `bool`:
			val.Field(i).SetBool(PromptYesNoDefault(os.Stdout, os.Stdin, helpEntry.Name,
				val.Field(i).Bool()))
		case `string`:
			val.Field(i).SetString(PromptStringDefault(os.Stdout, os.Stdin, helpEntry.Name,
				string(val.Field(i).String())))
		case `int`:
			val.Field(i).SetInt(int64(PromptIntegerDefault(os.Stdout, os.Stdin, helpEntry.Name,
				int(val.Field(i).Int()))))
		}
		//fmt.Printf("%2d: Field: %s (%s) Current value: %v\n", i+1, name, valueField.Type().Name(), valueField.Interface())

	}
}

// Get a list of names from the structure and include their help.
// Also return the maximum length of the names
func getStructList(substruct interface{}) (map[int]StructNameList, int) {
	val := reflect.Indirect(reflect.ValueOf(substruct))

	nameList := make(map[int]StructNameList)
	maxLength := 0

	for i := 0; i < val.NumField(); i++ {
		help := val.Type().Field(i).Tag.Get("help")
		if help != "" {
			name := val.Type().Field(i).Tag.Get("name")
			if name != "" {
				nameList[i] = StructNameList{Name: name, Help: help}
				if len(name) > maxLength {
					maxLength = len(name)
				}
			}
		}
	}
	return nameList, maxLength
}
func PrintStructValue(w io.Writer, substruct interface{}) {
	val := reflect.Indirect(reflect.ValueOf(substruct))
	nameList, maxLength := getStructList(substruct)
	fieldId := 0
	for i := 0; i < val.NumField(); i++ {
		entry, ok := nameList[i]
		if ok {
			fieldId = fieldId + 1
			fmt.Fprintf(w, "%2d: %*s: '%v'\n", fieldId, maxLength, entry.Name, val.Field(i).Interface())
		}
	}
}
func promptString(w io.Writer, r io.Reader, prompt string) (txt string, err error) {
	bin := bufio.NewReader(r)

	fmt.Fprintf(w, "%s ", prompt)
	txt,err = bin.ReadString('\n')
	//_, err := fmt.Fscanln(r, &txt)

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
	return
}
func PromptStringDefault(w io.Writer, r io.Reader, prompt, defaultValue string) string {

	prompt = fmt.Sprintf("%s\n  (enter='%s')? ", prompt, defaultValue)
	txt,err := promptString(w,r,prompt)

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

// Prompt for a yes/no value but accept a default value if they don't enter anything.
func PromptYesNoDefault(w io.Writer, r io.Reader, prompt string, defaultValue bool) bool {
	var txt string

	if defaultValue {
		txt = PromptStringDefault(w, r, prompt+" (yes/no)? ", "yes")
	} else {
		txt = PromptStringDefault(w, r, prompt+" (yes/no)? ", "no")
	}
	if val, err := ParseBool(txt); err == nil {
		return val
	}
	return defaultValue
}

// Simple prompt for interger. Caller can pass buffers or os.Stdin/os.Stdout
//
func PromptInteger(w io.Writer, r io.Reader, prompt string) (val int, err error) {
	var txt string

	for true {
		txt, err = promptString(w, r, prompt+" ? ")
		if err != nil {
			return
		}
		if txt != "" {
			if val, err = strconv.Atoi(txt); err == nil {
				return
			}
		}
		fmt.Fprintf(w, "\t%s '%s'\n", err.Error(), txt)
	}
	return
}

// Prompt for a integer value but accept a default value if they don't enter anything.
func PromptIntegerDefault(w io.Writer, r io.Reader, prompt string, defaultValue int) int {
	var txt string
	defaultString := strconv.Itoa(defaultValue)

	fullPrompt := fmt.Sprintf("%s (integer value) ", prompt)
	txt = PromptStringDefault(w, r, fullPrompt, defaultString)

	if val, err := strconv.Atoi(txt); err == nil {
		return val
	}
	return defaultValue
}
func Box(w io.Writer, info string) {
	fmt.Fprintf(w, "****************************\n%s\n************************\n", info)
}
