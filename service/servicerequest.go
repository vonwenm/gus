// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//

package service

import (

	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/cgentry/gav"
)

const (

	KEY_CMD		= "cmd"					// The requested action
	KEY_DOMAIN	= "domain"				// Logical group. If not present, "default"
	KEY_CALLER	= "caller"				// What is the ID of the caller
	KEY_HMAC	= "hmac"				// Checksum

	KEY_EMAIL	= "email"
	KEY_PWD		= "pwd"
	KEY_TOKEN	= "token"
	KEY_NAME	= "name"
	KEY_LOGIN	= "login"
	KEY_GUID	= "guid"
)

var StandardPathValues = []string{ KEY_CMD , KEY_DOMAIN, KEY_CALLER }

/*
 *	A ServiceRequest is a simple map that contains all the request parameters as key/value pairs
 *
 */
type ServiceRequest struct {
	request		* http.Request
	body		string						// Body of the message
	Parameters	map[string]string			// Encoded as either a=b&c=d or in header

	Hmac		string						// Hmac that was detected/set
	Date		string						// Date that was detected/set

	PathKeys    []string
	QueryKeys	[]string
	HeaderKeys  []string


	Security	* gav.Secure


}

func NewServiceRequest( r *http.Request ) * ServiceRequest {
	s.request = r
	s := &ServiceRequest{ Parameters : make(map[string]string) }
	s.setBody( r )
	s.Security = gav.NewSecure()
	return s
}

func ( sr * ServiceRequest ) GetUser() ( string , error ){
	return sr.Security.GetUser( request )
}

func ( sr * ServiceRequest ) ValidateSignature( userSecret string )( err errors ){
	err = sr.Security.VerifySignature( sr.request , userSecret ,[]byte( sr.body ) )
	return
}

func ( sr * ServiceRequest ) SetBody( r * http.Request )( err errors ){

	defer r.Body.Close()
	sr.body, err = ioutil.Readall( r.Body )
	return
}

func ( sr * ServiceRequest ) GetBody( ) string {
	return sr.body;
}

func( sr * ServiceRequest ) SetIfHmac( key , value string ) bool {
	if key == `hmac` || strings.ToLower( key ) == `x-srq-hmac` {
		sr.Hmac = value
		return true
	}
	return false
}

func ()
func( sr * ServiceRequest ) SetIfDate( key , value string ) bool {
	if key == `date` || strings.ToLower( key ) == `x-srq-date` {
		sr.Date = value
		return true
	}
	return false
}

func( sr * ServiceRequest ) GetHmac( ) ( string , bool ){
	return sr.Hmac , ( sr.Hmac == `` )
}

func (sr * ServiceRequest) SetPathKeys( paths []string ) * ServiceRequest {
	sr.PathKeys = paths
	return sr
}

func (sr * ServiceRequest) GetPathKeys( ) []string  {
	return sr.PathKeys
}

func (sr * ServiceRequest) SetQueryKeys( queryNames []string ) * ServiceRequest {
	sr.QueryKeys = queryNames
	return sr
}

func (sr * ServiceRequest) GetQueryKeys( ) []string  {
	return sr.QueryKeys
}

// Return all of the keys from the ServiceRequest in a sorted array.
func (sr * ServiceRequest) SortedKeys() []string {
	keys := make([]string, len(sr.Parameters))
	i := 0
	for key, _ := range sr.Parameters {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// Add a new key and value to the service map. Use of this will protect against structure changes
func (sr * ServiceRequest) Add(key, value string) * ServiceRequest {
	if ! sr.SetIfHmac( key , value ){
		sr.Parameters[key] = value
	}
	return sr
}

// Find a value in the service map. Use of this will protect against structure changes
func (sr * ServiceRequest) Get(key string) ( string , bool ) {
	val,found := sr.Parameters[key]
	return val,found
}



