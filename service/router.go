// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"net/http"
	"github.com/cgentry/gofig"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/record"
	//"fmt"
	"errors"
	"strings"
)

type ServiceHandler struct {}

func NewService(c * gofig.Configuration) {


	http.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) {ServiceRegister(c, w, r)})
	http.ListenAndServe(":8181", nil)
}

// ServiceRegister will handle the calling of Registration for the user.
func ServiceRegister(c *gofig.Configuration , w http.ResponseWriter, r *http.Request) {
	// var caller string

	// Need the request params. Since we have a standard format, parse by default
	qparam := []string{ KEY_EMAIL, KEY_LOGIN, KEY_NAME, KEY_PWD, KEY_HMAC}
	sr, err := ParseParms(r, StandardPathValues, qparam)
	if err != nil {
		ReturnError(w, CODE_BAD_CALL, err)
		return
	}
	caller, errCode , err := CheckCallerAndHmac(r , &sr)
	if err != nil {
		ReturnError( w , errCode , err )
		return
	}

	user := record.NewUser( caller.GetDomain() )
	for _,key := range qparam {
		val,_ := sr.Get( key )
		user.MapFieldToUser( key , val )
	}
	driver := storage.GetDriver()
	driver.RegisterUser( user )
	userReturn := record.NewReturnFromUser(user)

	ReturnUserJson( w , CODE_OK, &userReturn )
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
// From the query parameters, only those that are included will be split
func ParseParms(r *http.Request , list []string , qparam []string) (ServiceRequest , error) {
	sr := NewServiceRequest()
	sr.SetPathKeys(list )
	sr.SetQueryKeys(qparam)

	parts := strings.Split(r.URL.Path, "/")[1:]
	if len(parts) != len(list) {
		return sr, errors.New("Path was invalid")
	}

	// Match up the keys in list to the Path parts
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

	// Time to add in headers that match the pattern X-Srq-
	for key,value := range r.Header {
		if strings.HasPrefix( value , `X-Srq-`){
			sr.Add( key , value );
		}
	}
	return sr, nil
}




func CheckCallerAndHmac( r *http.Request, sr * ServiceRequest  ) (*record.User , int , error ){
	// We need the calling system's secret. This is the token for the caller
	drive := storage.GetDriver()
	caller_guid,found := sr.Get( KEY_CALLER )
	if ! found {
		return  nil, CODE_BAD_CALL , errors.New("Missing Caller identifier")
	}
	caller, err := drive.FetchUserByGuid( caller_guid )
	if err != nil {
		return  nil , CODE_USER_DOESNT_EXIST , err
	}
	if caller.IsSystem != true {
		return  nil, CODE_INVALID_REQUEST , errors.New("Not a valid caller")
	}

	key, err := CreateRestfulHmac(caller.GetToken(), r, sr)
	if err != nil {
		return caller , CODE_BAD_CALL , err
	}
	if ! CompareHmac(key , sr ){
		return caller , CODE_BAD_CALL , errors.New("HMAC errors")
	}
	return caller , CODE_OK , nil
}

