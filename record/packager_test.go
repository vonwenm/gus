package record

import (
	"github.com/cgentry/gus/record/request"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type _testDummyPackageBody struct {
	TestBody string
}

func TestPackager(t *testing.T) {
	Convey("Check Head", t, func() {
		h := request.NewHead()
		h.Domain = "Domain"
		h.Id = "id"
		h.SetSignature("")

		p := NewPackage()

		Convey("Check functions with mainly empty values", func() {
			So(p.IsBodySet(), ShouldBeFalse)
			So(p.IsHeadSet(), ShouldBeFalse)
			So(p.IsPackageComplete(), ShouldBeFalse)

			p.SetHead(h)
			So(p.GetSignature(), ShouldBeBlank)
			So(p.IsHeadSet(), ShouldBeTrue)
			So(p.IsBodySet(), ShouldBeFalse)
			So(p.IsPackageComplete(), ShouldBeFalse)
			_, y := p.Head.(request.Head)
			So(y, ShouldBeTrue)
			So(p.GoodSignature(), ShouldBeFalse)

			p.SetSecret([]byte(`aSecret`))
			So(p.GoodSignature(), ShouldBeFalse)
		})

		Convey("Check functions with completed values", func() {

			p.SetHead(h)
			So(p.GetSignature(), ShouldBeBlank)
			p.SetSecret([]byte(`abcdefSecret`))
			p.SetBodyString("Hello there")
			So(p.Body, ShouldContainSubstring, `Hello there`)

			So(p.GetSignature(), ShouldNotBeBlank)
			So(p.GoodSignature(), ShouldBeTrue)
			p.ClearSecret()
			So(p.GoodSignature(), ShouldBeFalse)

			p.SetSecret([]byte(`anotherSecret`))
			d := _testDummyPackageBody{TestBody: `test body`}
			p.SetBody(d)
			So(p.GoodSignature(), ShouldBeTrue)
			So(p.IsBodySet(), ShouldBeTrue)
			So(p.Body, ShouldContainSubstring, `test body`)

			c := make(chan int)
			p.SetBody(c)
			So(p.GoodSignature(), ShouldBeFalse)
			So(p.IsBodySet(), ShouldBeFalse)
		})

	})
}
