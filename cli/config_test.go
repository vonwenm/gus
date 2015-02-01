package cli

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestYesNo(t *testing.T) {
	var bufin bytes.Buffer
	var bufout bytes.Buffer

	Convey("Prompt YES input", t, func() {
		bufin.Write([]byte("yes\n"))
		rtn, _ := PromptYesNo(&bufout, &bufin, "Do you like pizza")
		So(rtn, ShouldBeTrue)
	})

	Convey("Prompt NO input", t, func() {
		bufin.Reset()
		bufin.Write([]byte("No\n"))
		rtn, _ := PromptYesNo(&bufout, &bufin, "Do you like pizza")
		So(rtn, ShouldBeFalse)
	})

}

func TestBadYesNoInput(t *testing.T) {
	var bufin bytes.Buffer
	var bufout bytes.Buffer

	bufin.Write([]byte("WHAT?\n"))

	rtn, err := PromptYesNo(&bufout, &bufin, "Do you like pizza")

	if err == nil || err.Error() != "EOF" {
		t.Error("Error was not EOF")
	}
	if rtn != false {
		t.Error("Return was not false")
	}
	str := bufout.String()

	if !strings.Contains(str, "Do you like pizza") {
		t.Error("Does not contain prompt")
	}
	if !strings.Contains(str, "WHAT?") {
		t.Errorf("WHAT is not in the error message '%s'", str)
	}
}
func TestPromptEmptyString(t *testing.T) {
	var bufin bytes.Buffer
	var bufout bytes.Buffer

	bufin.Write([]byte("\n"))

	rtn, err := PromptString(&bufout, &bufin, "Do you like pizza")

	if err == nil {
		t.Error("Error was nil (should be EOF)")
	}
	if err != nil && err.Error() != "EOF" {
		t.Errorf("Error was not EOF (%s)", err.Error())
	}
	if rtn != "" {
		t.Error("Return was not empty")
	}
	str := bufout.String()

	if !strings.Contains(str, "Do you like pizza") {
		t.Error("Does not contain prompt")
	}
	if !strings.Contains(str, "Invalid string") {
		t.Errorf("Does not contain \"Invalid string\" ('%s')", str)
	}
}
func TestPromptReturnString(t *testing.T) {
	var bufin bytes.Buffer
	var bufout bytes.Buffer

	bufin.Write([]byte("Yes\n"))

	rtn, err := PromptString(&bufout, &bufin, "Do you like pizza")

	if err != nil {
		t.Errorf("Error was not EOF (%s)", err.Error())
	}
	if rtn == "" {
		t.Error("Return was empty")
	}
	str := bufout.String()
	if !strings.Contains(str, "Do you like pizza") {
		t.Errorf("Does not contain prompt (%s)", str)
	}

	if !strings.Contains(rtn, "Yes") {
		t.Errorf("Does not contain \"Yes\" ('%s')", rtn)
	}
}
func TestPromptStringDefault(t *testing.T) {
	var bufin bytes.Buffer
	var bufout bytes.Buffer

	bufin.Write([]byte("\n"))

	rtn := PromptStringDefault(&bufout, &bufin, "Do you like pizza", "yes")

	if rtn == "" {
		t.Error("Return was empty")
	}

	if !strings.Contains(bufout.String(), "Do you like pizza") {
		t.Error("Does not contain prompt")
	}
	if rtn != "yes" {
		t.Errorf("Reply does not contain 'yes' (%s)\n", rtn)
	}
}
