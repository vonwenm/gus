package record

import (
	"github.com/cgentry/gus/record/head"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type _testDummyPackageBody struct {
	TestBody string
}

func TestPackager(t *testing.T) {
	Convey("Check Head", t, func() {
		var p Packer
		var hi *head.Head

		h := head.New()
		hi = h
		h.Domain = "Domain"
		h.Id = "id"
		h.SetSignature([]byte(""))
		hi.Check()

		p = NewPackage()

		Convey("Check functions with mainly empty values", func() {
			//So(p.IsBodySet(), ShouldBeTrue)
			//So(p.IsHeadSet(), ShouldBeTrue)
			//So(p.IsPackageComplete(), ShouldBeFalse)

			p.SetHead(h)
			sig, err := p.GetHead().GetSignature()
			So( err, ShouldBeNil)
			So(string(sig), ShouldBeBlank)
			So(p.IsHeadSet(), ShouldBeTrue)
			//So(p.IsBodySet(), ShouldBeFalse)
			//So(p.IsPackageComplete(), ShouldBeFalse)

			So(GoodSignature(p), ShouldBeFalse)

			p.SetSecret([]byte(`aSecret`))
			So(GoodSignature(p), ShouldBeFalse)
		})

		Convey("Check functions with completed values", func() {

			p.SetHead(h)
				sig, err := p.GetHead().GetSignature()
				So( err, ShouldBeNil)
				So(string(sig), ShouldBeBlank)
				p.SetSecret([]byte(`abcdefSecret`))
			p.SetBody("Hello there")
			So(p.GetBody(), ShouldContainSubstring, `Hello there`)
			SignPackage(p)

			So(p.GetHead().IsSignatureSet() , ShouldBeTrue )
			So(GoodSignature(p), ShouldBeTrue)
			p.ClearSecret()
			So( string(p.GetSecret()) , ShouldBeBlank)
			So(GoodSignature(p), ShouldBeFalse)

			p.SetSecret([]byte(`anotherSecret`))
			d := _testDummyPackageBody{TestBody: `test body`}
			p.SetBodyMarshal(d)
			SignPackage(p)
			So(GoodSignature(p), ShouldBeTrue)
			So(p.IsBodySet(), ShouldBeTrue)
			So(p.GetBody(), ShouldContainSubstring, `test body`)
		})
	})
}
