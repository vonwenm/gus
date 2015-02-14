package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/cgentry/gus/record/configure"
	"net/http/httptest"
	"net/http"
)

func TestParseParms(t *testing.T) {

	Convey("Parse Parms", t, func() {
			c := configure.New()
			c.Default()
			So( c.Service.Port, ShouldNotEqual, 0 )
			Convey( "Test ping" , func(){
					route := RouteService{Handler: httpPing , Server: nil }
					req, err := http.NewRequest("GET", "http://example.com/ping", nil)
					So( err, ShouldBeNil)
					w := httptest.NewRecorder()
					httpPing(c , route , "/ping/" , w, req)
					So( w.Code, ShouldEqual, 200)
					So( w.Body.String(), ShouldEqual, "/ping/")
				})
			Convey( "Test Ping route" , func(){

				})
		/*
			config := gofig.NewConfiguration()
			writer := httptest.NewRecorder()
			reader, err := http.NewRequest( "GET" , "http://test.com/register")
			So( err , ShouldBeNil )
			usr,rqst, resp := StdParseParms( config. writer, reader )

			So( user , ShouldBeNil )
			So( rqst , ShouldNotBeNil )
			So( resp , ShouldNotBeNil )
		*/
	})
}
