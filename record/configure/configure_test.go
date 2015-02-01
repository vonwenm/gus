package configure

import (
	_ "encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSimpleString(t *testing.T) {
	Convey("Marsh to string", t, func() {
		c := &Configure{}
		c.Service.Host = "hostname"
		c.Service.Port = 9090
		c.User.Name = "sqlite"
		c.User.Dsn = "dsn"
		c.User.Options = "options"
		x := c.String()
		So(x, ShouldContainSubstring, `"hostname"`)
		So(x, ShouldContainSubstring, `9090`)
		So(x, ShouldContainSubstring, `"sqlite"`)
		So(x, ShouldContainSubstring, `"dsn"`)

	})
}
func TestSimpleEncode(t *testing.T) {
	Convey("Encode String to class", t, func() {
		config, err := NewConfigure(t_config_test_1)
		So(err, ShouldBeNil)
		So(config.Service.Host, ShouldEqual, `hostname`)

	})

}

const t_config_test_1 = `{
  "Service": {
    "Host": "hostname",
    "Port": 9090,
    "SepStore": true
  },
  "User": {
    "Name": "sqlite",
    "Dsn": "dsn",
    "Options": "options"
  }
}`
