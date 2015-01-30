package configure

// Configure holds the structures used to drive the program and setup connections.
// These are JSON encoded and stored, somewhere, and loaded at program startup.
// There is a CLI interface to interactively build the options
import (
	"encoding/json"
)

type Configure struct {
	Service Service
	User  Store		`help:"The storage for the user data"`
	Client Store	`help:"The storage for the client can be different than for the user store"`
}


type Store struct {
	Name    string `help:"The storage driver you want to use." name:"Driver Name"`
	Dsn     string `help:"The specific driver data-name. Usually how to connect to the database." name:"DSN"`
	Options string `help:"Options passed to the driver. Check the driver for what options are availble." name:"Driver options"`
}

type Service struct {
	Host string  `name:"Hostname" help:"The IP address or host name to listen on. Leave empty to listen to all."`
	Port int	`name:"Port" help:"The port that the service should listen on."`
	ClientId int `name:"Header ID for Client" help:"If you are using an HTTPS load balancer, what header is set for the client id. (Must match the Email address.)`
	SepStore bool `name:"Separate User/Client store" help:"Do you want separate client and user storage?"`
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

const DEFAULT_CONFIG=`{
  "Service": {
    "Port": 9090,
    "SepStore": false
  },
  "User": {
    "Name": "mongo",
    "Dsn": "dsn",
    "Options": "User"
  }
}`
