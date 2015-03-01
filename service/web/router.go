// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package web

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/service"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// These constants provide the route string used in the http.Handle call.
const (
	SRV_REGISTER = "/register/"
	SRV_LOGIN    = "/login/"
	SRV_LOGOUT   = "/logout/"
	SRV_AUTH     = "/authenticate/"
	SRV_ENABLE   = "/enable/"
	SRV_DISABLE  = "/disable/"
	SRV_PING     = "/ping/"
	SRV_UPDATE   = "/update/"
	SRV_HOME     = "/"
	SRV_TEST     = "/test/"

	GUS_VERSION = "0.1"
)

// Route to service defines what a function needs to look like in order for us to call it when
// the route matches.
type RouteToService func(c *configure.Configure, rhandle RouteService, name string, w http.ResponseWriter, r *http.Request)

// RouteService defines a begining point (Handler) and what service we use (ServiceCreator)
type RouteService struct {
	Handler RouteToService
	Server  service.ServiceCreator
}

// RouteTable contains a route name pointing to a service definition.
type RouteTable map[string]RouteService

var RouteMap = RouteTable{
	SRV_REGISTER: {Handler: httpCallService, Server: service.NewServiceRegister},
	SRV_LOGIN:    {Handler: httpCallService, Server: service.NewServiceLogin},
	SRV_LOGOUT:   {Handler: httpCallService, Server: service.NewServiceLogout},
	SRV_AUTH:     {Handler: httpCallService, Server: service.NewServiceAuthenticate},
	SRV_UPDATE:   {Handler: httpCallService, Server: service.NewServiceUpdate},
	SRV_TEST:     {Handler: httpCallService, Server: service.NewServiceTest},
	//SRV_ENABLE:   {Handler: httpCallService , Server: service.NewServiceEnable } ,
	//SRV_DISABLE:  {Handler: httpCallService , Server: service.NewServiceDisable },
	SRV_PING: {Handler: httpPing, Server: nil},
	SRV_HOME: {Handler: httpHome, Server: nil},
}

type RouteHandler struct {
	config *configure.Configure
}

// New creates a new route handler. Route handlers setup the table used to map requests
// to a service function that will call the http Service routine
func New(c *configure.Configure) *RouteHandler {
	return &RouteHandler{config: c}
}

// CreateHandlerFunc is a private function that creates an http.Handler function for the Go http.Handle function.
// This allows us to pass in extra parameters. The 'RouteService' gives us the linking
// points needed.
func (s *RouteHandler) CreateHandlerFunc(name string, rhandle RouteService) http.Handler {
	config := s.config
	fn := func(w http.ResponseWriter, r *http.Request) {
		rhandle.Handler(config, rhandle, name, w, r)
		return
	}
	return http.HandlerFunc(fn)
}

// Register all of the routes to the go handle. The process takes the RouteTable which
// consists of the http path to match and the RouteService. The route service entry
// holds the http function that will call the service function:
// path => http bundling routine -> (calls) -> ServiceRoute, which contains a run entry
func (s *RouteHandler) Register(rmap RouteTable) *RouteHandler {
	for key, handle := range rmap {
		http.Handle(key, s.CreateHandlerFunc(key, handle))
	}
	return s
}

// Serve starts listening and serving content.
func (s *RouteHandler) Serve() {
	serviceAddress := fmt.Sprintf("%s:%d", s.config.Service.Host, s.config.Service.Port)
	http.ListenAndServe(serviceAddress, nil)
	return
}

// Ping is one of the route routines that will simply return a string back to the user. This does not
// call a service routine
func httpPing(c *configure.Configure, rhandle RouteService, name string, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(name))
	return
}

// httpHome is one of the route routines that will simply return an error if the pattern doesn't match antyhing
func httpHome(c *configure.Configure, rhandle RouteService, name string, w http.ResponseWriter, r *http.Request) {
	httpErrorWrite(w, 404, "Invalid page request '"+r.URL.Path+"'")
	return
}

// httpCallService is the main router for setup and route to the service routines. It instantiates the
// service and then calls the Run() for the service
func httpCallService(c *configure.Configure, rhandle RouteService, name string, w http.ResponseWriter, r *http.Request) {
	var err error
	var srv *service.ServiceProcess

	httpRequestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		httpErrorWrite(w, http.StatusBadRequest, err.Error())
		return
	}

	srv = rhandle.Server()

	returnPackage, err := srv.SetupService(c, string(httpRequestBody))

	if err == nil {
		// Check here for client server match// Check here for client certifcate match
		returnPackage, err = srv.Run(srv)
	}

	srv.Teardown()
	httpResponseWrite(w, returnPackage, err)
	return
}

// httpErrorWrite will pack the code and message into a standard response and call
// httpResponseWrite to return the final result
func httpErrorWrite(w http.ResponseWriter, code int, msg string) {
	responsePackage := record.NewPackage()
	err := ecode.NewGeneralError(msg, code)
	responsePackage.SetBodyMarshal(err)
	responsePackage.SetBodyType(record.PACKAGE_BODYTYPE_ERROR)
	httpResponseWrite(w, responsePackage, err)
}

// httpResponseWrite returns to the caller the body and headers properly set for this
// type of call. Headers for type, md5, code and message will all be set
func httpResponseWrite(w http.ResponseWriter, responsePackage record.Packer, err error) {
	httpResponseBody, convertErr := json.Marshal(responsePackage)
	if convertErr != nil {
		httpResponseBody = []byte("{}")
	}

	if err == nil {
		err = ecode.ErrStatusOk
	}
	w.Header().Set("Message", err.Error())
	errInternal, ok := err.(*ecode.GeneralError)
	if ok {
		w.WriteHeader(errInternal.Code())
	}

	// Set the encoding type (JSON)
	w.Header().Set("Content-Type", "text/json")

	// Set the header MD5 checksum
	rsvpMD5 := md5.New()
	io.WriteString(rsvpMD5, string(httpResponseBody))
	rsvpChecksum := base64.StdEncoding.EncodeToString(rsvpMD5.Sum(nil))
	w.Header().Set("CONTENT-MD5", rsvpChecksum)

	w.Header().Set("ETag", rsvpChecksum)
	w.Header().Set("Date", time.Now().Format(time.RFC1123))
	w.Header().Set("Content-Length", strconv.Itoa(len(httpResponseBody)))
	w.Header().Set("Cache-Control", "no-store") // Don't cache this record
	w.Header().Set("Server", "gus/"+GUS_VERSION)

	w.Write(httpResponseBody)
	return
}
