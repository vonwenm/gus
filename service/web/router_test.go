package web

import (
	"bytes"
	"encoding/json"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/service"
	"github.com/cgentry/gus/storage"
	_ "github.com/cgentry/gus/storage/mock"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseParms(t *testing.T) {
	c := configure.New()
	c.Default()

	Convey("Parse Parms", t, func() {

		So(c.Service.Port, ShouldNotEqual, 0)
		Convey("Test ping", func() {
			route := RouteService{Handler: httpPing, Server: nil}
			req, err := http.NewRequest("GET", "http://example.com/ping", nil)
			So(err, ShouldBeNil)
			w := httptest.NewRecorder()
			httpPing(c, route, "/ping/", w, req)
			So(w.Code, ShouldEqual, 200)
			So(w.Body.String(), ShouldEqual, "/ping/")
		})
		Convey("Test Ping route", func() {
			var testMap = RouteTable{SRV_PING: {Handler: httpPing, Server: nil}}

			w := New(c)
			serve := httptest.NewServer(w.CreateHandlerFunc(SRV_PING, testMap[SRV_PING]))
			defer serve.Close()
			res, err := http.Get(serve.URL + "/ping")
			So(err, ShouldBeNil)
			pingTxt, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			So(err, ShouldBeNil)
			So(string(pingTxt), ShouldEqual, "/ping/")

		})
		Convey("Test Home route", func() {
			var testMap = RouteTable{SRV_HOME: {Handler: httpHome, Server: nil}}

			w := New(c)
			serve := httptest.NewServer(w.CreateHandlerFunc(SRV_HOME, testMap[SRV_HOME]))
			defer serve.Close()
			res, err := http.Get(serve.URL + "/nothere")
			So(err, ShouldBeNil)
			pingTxt, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			So(err, ShouldBeNil)
			So(string(pingTxt), ShouldEndWith, `,"Code":404,"Message":"Invalid page request '/nothere'","Id":"","Sequence":0},"Body":""}`)

			res, err = http.Get(serve.URL + "/")
			So(err, ShouldBeNil)
			pingTxt, err = ioutil.ReadAll(res.Body)
			res.Body.Close()
			So(err, ShouldBeNil)
			So(string(pingTxt), ShouldEndWith, `,"Code":404,"Message":"Invalid page request '/'","Id":"","Sequence":0},"Body":""}`)

		})
		Convey("Test 'test' route", func() {

			var testMap = RouteTable{SRV_TEST: {Handler: httpCallService, Server: service.NewServiceTest}}
			c.User.Name = "mock"
			w := New(c)

			// Fill in a client record
			store, _ := storage.Open("mock", "", "")

			user := record.NewTestUser()
			user.IsSystem = true
			store.UserInsert(user)

			// Create a "test" request. This requires a few fields
			pack := record.NewPackage()
			head := request.NewHead()
			head.Domain = user.Domain
			head.Id = user.LoginName
			body := request.NewTest()
			pack.SetHead(head)
			pack.SetBody(body)
			rqstBody, err := json.Marshal(pack)
			So(err, ShouldBeNil)

			serve := httptest.NewServer(w.CreateHandlerFunc(SRV_TEST, testMap[SRV_TEST]))
			defer serve.Close()
			buff := bytes.NewBuffer(rqstBody)
			res, err := http.Post(serve.URL+"/test/", "text/json", buff)

			So(err, ShouldBeNil)
			pingTxt, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			So(err, ShouldBeNil)
			So(string(pingTxt), ShouldEndWith, `,"Code":404,"Message":"Invalid page request '/nothere'","Id":"","Sequence":0},"Body":""}`)

		})
	})
}
