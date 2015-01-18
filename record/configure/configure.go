package configure

// Configure holds the structures used to drive the program and setup connections.
// These are JSON encoded and stored, somewhere, and loaded at program startup.
// There is a CLI interface to interactively build the options
import (
	"encoding/json"
)

type Configure struct {
	Service Service
	Driver  Driver
}

type Driver struct {
	Name    string
	Dsn     string
	Options string
}

type Service struct {
	Host string
	Port int
}

type Configurer interface {
	DecodeString(string) error
	String() string
}

func NewConfigure(encodedConfig string) (*Configure, error) {
	c := &Configure{}
	return c, c.DecodeString(encodedConfig)
}

func (c *Configure) DecodeString(encodedConfig string) error {
	err := json.Unmarshal([]byte(encodedConfig), c)
	return err
}
func (c *Configure) String() string {
	s, _ := json.MarshalIndent(c, ``, `  `)
	return string(s)
}
