// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"net/http"
	"github.com/cgentry/gofig"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gosr"
	gosrhttp "github.com/cgentry/gosr/http"
	//"fmt"
	"errors"
	"strings"
)

type ServiceHandler struct {}

func NewService(c * gofig.Configuration) {

	http.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) {ServiceRegister(c, w, r)})
	http.ListenAndServe(":8181", nil)
}

// ServiceRegister The body of a request must contain the registration details for a
// new user.
func ServiceRegister(c *gofig.Configuration , w http.ResponseWriter, r *http.Request) {

	caller , rqst , answr := parseParms(c, w, r)
	if caller != nil {

		// Need the request params. Since we have a standard format, parse by default

		qparam := []string{ KEY_EMAIL, KEY_LOGIN, KEY_NAME, KEY_PWD}
		if err := rqst.Parameters.IsPresent(qparam); err != nil {
			answ.SetError(err)
			answ.Encode(w)
		}else {
			user := record.NewUser(caller.GetDomain())
			err := user.Unmarshall( rqst.GetBody() )					// Pass in the body of the request
			if err != nil {
				answr.Status = http.PreconditionFailed
				answr.StatusText = err.Error()
			}else{
				// OK...now the real work
			}
			for _, key := range qparam {
				user.MapFieldToUser(key, rqst.Parameters.Get(key))
			}

			driver := storage.GetDriver()
			driver.RegisterUser(user)

			userReturn := record.NewReturnFromUser(user)

			ReturnUserJson(w, CODE_OK, &userReturn)
		}
	}
	return
}

// ParseParms takes the url path and puts it into a standard map
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

// ParseParms expects a list of path parameters and a list of query parameters that are required.
// From the query parameters, only thos that are included will be split
func parseParms(c *gofig.Configuration , w http.ResponseWriter, r *http.Request) (
	 *record.User,*gosrhttp.Request, *gosrhttp.Response) {

	var err error
	var subscriber * record.User

// Decode the request into standard request format
	rqst := gosrhttp.NewRequest()
	answr := gosr.NewResponse()

	if err = rqst.Decode(r); err == nil {

		// Need the subscriber's record
		subscriber, err = FindUser( rqst.GetUser() )
		if err == nil && subscriber.IsSystem {
			err = rqst.Verify( []byte( subscriber.GetSalt()) , 15 )
		}else{
			err.StatusText = gosr.INVALID_SUBSCRIBER
		}
	}
	if err != nil {
		answr.SetError(err).Encode(w)					// Pack and send
		subscriber = nil
	}
	return subscriber , rqst , answr
}




func FindUser( caller_guid string  ) (*record.User , * gosr.Error ){

	// We need the calling system's secret. This is the token for the caller
	drive := storage.GetDriver()
	caller, err := drive.FetchUserByGuid( caller_guid )
	if err != nil {
		return  nil , gosr.NewErrorWithText( http.StatusBadRequest, gosr.INVALID_USER_PWD )
	}

	return caller ,  nil
}

func CheckParameters( *)

