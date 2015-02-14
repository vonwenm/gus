package response

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestNewHead(t *testing.T) {
	h := NewHead()
	Convey("Generate New Head", t, func() {
		So(h.GetStamp().IsZero(), ShouldBeFalse)
		So(h.IsTimeSet(), ShouldBeTrue)
	})
}

func TestGetSignature(t *testing.T) {
	h := NewHead()
	Convey("Set and test signature values", t, func() {
		So(h.GetSignature(), ShouldBeBlank)
		h.SetSignature("abcdef")
		So(h.GetSignature(), ShouldEqual, "abcdef")
	})
}
func TestCheckheader(t *testing.T) {
	h := NewHead()
	Convey("Check Head", t, func() {
		Convey("No dt", func() {
				h.SetStamp( time.Unix(0,0))
			err := h.Check()

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "No timestamp")
		})

		Convey("Check with domain, token good Time", func() {
			h.SetStamp(time.Now() )
			So(h.IsTimeSet(), ShouldBeTrue)
			err := h.Check()
			So(err, ShouldBeNil)
		})
	})
}

func TestTimeRange(t *testing.T) {
	var now time.Time
	Convey("Check TimeRange", t, func() {
			h := NewHead()
			now = h.GetStamp()

			offset := 8*time.Minute + 1*time.Second

			So(now.Equal(h.GetStamp()), ShouldBeTrue)
			So(h.Check(), ShouldBeNil)

			h.SetStamp(now.Add(offset))
			So(h.Check(), ShouldNotBeNil)
			So(h.Check().Error(), ShouldContainSubstring, "Request in the future")

			h.SetStamp(now.Add(-1 * offset))
			So(h.Check(), ShouldNotBeNil)
			So(h.Check().Error(), ShouldContainSubstring, "Request expired")

		})
}
