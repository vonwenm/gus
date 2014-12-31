package record

import (
	"github.com/cgentry/gus/record/request"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPackager(t *testing.T) {
	Convey("Check Head", t, func() {
		h := request.NewHead()
		h.Domain = "Domain"
		h.Id = "id"
		h.SetSignature("")

		p := NewPackage()

		Convey("Check IsXxx functions", func() {
			So(p.IsBodySet(), ShouldBeFalse)
			So(p.IsHeadSet(), ShouldBeFalse)

			p.SetHead(h)
			So(p.IsHeadSet(), ShouldBeTrue)
			So(p.IsBodySet(), ShouldBeFalse)
			_, y := p.Head.(request.Head)
			So(y, ShouldBeTrue)
		})

	})
}
