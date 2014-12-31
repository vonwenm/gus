package request

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
		Convey( "Check empty header" , func(){
			err := h.Check()
			So( err , ShouldNotBeNil )
			So( err.Error(), ShouldContainSubstring, "No domain")
		})
		Convey( "Check with domain only" , func(){
			h.Domain = "Domain"
			err := h.Check()
			So( err , ShouldNotBeNil )
			So( err.Error(), ShouldContainSubstring, "No Id")
		})
		Convey( "Check with domain, token but no dt" , func(){
			h.Domain = "Domain"
			h.Id = "id"
			err := h.Check()
			So( err , ShouldNotBeNil )
			So( err.Error() , ShouldContainSubstring , "No timestamp" )
		})
		Convey( "Check with domain, token good Time" , func(){
			h.Domain = "Domain"
			h.Id = "id"
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
		h.Domain = "Domain"
		h.Id = "id"
		offset := 2 * time.Minute + 1 * time.Second

		So( now.Equal( h.Timestamp ) , ShouldBeTrue )
		So( h.Check() , ShouldBeNil )

		h.Timestamp = h.Timestamp.Add(-1 * offset)
		So( h.Check() , ShouldNotBeNil )
		So( h.Check().Error(), ShouldContainSubstring, "Request expired")


		h.Timestamp = now.Add( offset )
		So( h.Check() , ShouldNotBeNil )
		So( h.Check().Error(), ShouldContainSubstring, "Request in the future")



	})
}
