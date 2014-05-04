// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"net/http"
	"github.com/cgentry/gofig"
	"fmt"
	"encoding/json"
	"errors"
	"strings"
)

type ServiceHandler struct {}

type StatusReturn struct {
	Status     int
	Message    string
}

var stdPathParam = []string{"cmd", "domain", "caller", "hmac"}

func NewService(c * gofig.Configuration) {

	http.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) {ServiceRegister(c, w, r)})
	/*

	*/

	http.ListenAndServe(":8181", nil)
}

// ServiceRegister will handle the calling of Registration for the user.
func ServiceRegister(c *gofig.Configuration , w http.ResponseWriter, r *http.Request) {
	// Need the request params. Since we have a standard format, parse by default
	qparam := []string{ "email", "login", "name", "password"}
	srequest, err := ParseParms(r, stdPathParam, qparam)
	if err != nil {
		ReturnError(w, CODE_BAD_CALL, err)
	}else {
		fmt.Fprintf(w, "Good boy!<br>%s<br>", srequest)
	}

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
func ParseParms(r *http.Request , list []string , qparam []string) (ServiceRequest , error) {
	sr := NewServiceRequest()
	sr.SetPathKeys( list )

	parts := strings.Split(r.URL.Path, "/")[1:]
	plen := len(parts)
	if plen != len(list) {
		return sr, errors.New("Path was invalid")
	}

	// Match up all the keys
	for i, key := range list {
		sr.Add( key , parts[i] )
	}

	// AND now for the parameters (follows the ? portion)
	query := r.URL.Query()
	for _, key := range qparam {
		if _, found := query[key]; !found {
			return sr, errors.New("Missing query parameter '"+key+"'")
		}else {
			sr.Add( key , query.Get(key) )
		}
	}
	return sr, nil
}


func ReturnError(w http.ResponseWriter , code int , err error) {

	msg := StatusReturn{ Status : code , Message : err.Error() }
	rtn, _ := json.Marshal(msg)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(rtn)
}

