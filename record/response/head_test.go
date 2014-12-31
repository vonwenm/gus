package response

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"time"
)

func TestNewHead( t * testing.T ){
	Convey("Generate New Head", t, func() {
		h := NewHead()
		So( h.Timestamp.IsZero() , ShouldBeFalse )
		So( h.IsTimeSet() , ShouldBeTrue)
	})
}

func TestGetSignature( t * testing.T ){
	Convey("Set and test signature values", t, func() {
		h := NewHead()
		So( h.GetSignature() , ShouldBeBlank )
		h.SetSignature("abcdef")
		So( h.GetSignature(), ShouldEqual, "abcdef")
	})
}
func TestCheckheader( t * testing.T ){
	Convey("Check Head", t, func() {
		h := Head{}
		So( h.IsTimeSet() , ShouldBeFalse )

		Convey( "No dt" , func(){

			err := h.Check()
			So( err , ShouldNotBeNil )
			So( err.Error() , ShouldContainSubstring , "No timestamp" )
		})
		Convey( "Check with domain, token good Time" , func(){

			h.Timestamp = time.Now()
			So( h.IsTimeSet() , ShouldBeTrue)
			err := h.Check()
			So( err , ShouldBeNil )
		})

	})
}
func TestTimeRange( t * testing.T ){
	var now time.Time
	Convey("Check TimeRange", t, func() {
		h := NewHead()
		now = h.Timestamp

		So( now.Equal( h.Timestamp ) , ShouldBeTrue )
		So( h.Check() , ShouldBeNil )


		h.Timestamp = now.Add( 2 * time.Minute + 1 * time.Second )
		So( h.Check() , ShouldNotBeNil )
		So( h.Check().Error(), ShouldContainSubstring, "Request in the future")

		h.Timestamp = now.Add( -2 * time.Minute + -1 * time.Second )
		So( h.Check() , ShouldNotBeNil )
		So( h.Check().Error(), ShouldContainSubstring, "Request expired")

	})
}
