package stamp

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	tim := New()

	Convey("Set and test time values", t, func() {
		So( tim.Age(), ShouldBeLessThan, 1 )

		newTime, _ := time.Parse("2006 Jan 02 15:04:05", "2015 Jan 01 12:15:30.918273645")
		tim.SetStamp( newTime )
		So( newTime.Equal( tim.GetStamp()), ShouldBeTrue )

		age := int( time.Now().Sub( newTime).Seconds() )
		So( age, ShouldBeGreaterThan, 1000 )

		So( age, ShouldEqual, tim.Age())
	})
}
