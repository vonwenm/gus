// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"github.com/cgentry/gofig"
	"github.com/cgentry/gosr"
	"github.com/cgentry/gosr/ghttp"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/storage"
	"net/http"
)

type ServiceHandler struct{}
type ParseParms func(*gofig.Configuration, http.ResponseWriter, *http.Request) (*record.User, *ghttp.Request, *ghttp.Response)

func NewService(c *gofig.Configuration) {

	http.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) { httpRegister(c, w, r) })
	http.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) { ServiceLogin(c, w, r) })
	http.HandleFunc("/logout/", func(w http.ResponseWriter, r *http.Request) { ServiceLogout(c, w, r) })
	http.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) { ServiceUpdate(c, w, r) })
	http.HandleFunc("/enable/", func(w http.ResponseWriter, r *http.Request) { ServiceEnable(c, w, r) })
	http.HandleFunc("/diable/", func(w http.ResponseWriter, r *http.Request) { ServiceDisable(c, w, r) })

	http.ListenAndServe(":8181", nil)
}

func ServiceLogin(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}
func ServiceLogout(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}
func ServiceUpdate(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}
func ServiceEnable(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}
func ServiceDisable(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}

// StdParseParms takes the url path and puts it into a standard map
// the request always looks like:
// /cmd/domain/caller-appid/hmac/identifier....
//			login: identifier(login-name)/password
//			auth:  identifier (token)
//			logout: identifier
//			register: identifier(login-name)/password/email
//			lookup: identifier/type (email|name|guid)
//			save:   identifier/type(session|user)/name(of item)
//			retrieve: identifier/type/name
//
//			inactive: indentifer/type
//			active:
//

/**
 * StdParseParms
 *		Parse the parameters from the request and put them into a standard format
 *		If there are errors, we will return the answer with an error
 */
var StdParseParms ParseParms = func(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) (
	*record.User, *ghttp.Request, *ghttp.Response) {

	var err error
	var subscriber *record.User

	// Decode the request into standard request format
	rqst := ghttp.NewRequest()
	answr := ghttp.NewResponse()

	if err = rqst.Decode(r, ""); err == nil { // Decode the request

		// Need the subscriber's record
		subscriber, err = FindUser(rqst.GetUser(), true) // All requests must have a user...
		if err == nil {                                  // .. got one
			err = rqst.Verify([]byte(subscriber.GetSalt()), 15)
		}
	}
	if err != nil {
		answr.SetError(gosr.NewErrorWithText(CODE_BAD_CALL, err.Error())).Encode(w) // Pack and send
		subscriber = nil
	}
	return subscriber, rqst, answr
}

func sendError(c *gofig.Configuration, w http.ResponseWriter) {

}

/*
 * Lookup the user withing the data store. Simple wrapper.
 */
func FindUser(caller_guid string, needSystem bool) (*record.User, *gosr.Error) {

	// We need the calling system's secret. This is the token for the caller
	drive := storage.GetDriver()
	caller, err := drive.FetchUserByGuid(caller_guid)
	if err != nil {
		return nil, gosr.NewErrorWithText(http.StatusBadRequest, gosr.INVALID_USER_PWD)
	}
	if needSystem && !caller.IsSystem {
		return nil, gosr.NewErrorWithText(http.StatusBadRequest, gosr.INVALID_USER_PWD)
	}

	return caller, nil
}
