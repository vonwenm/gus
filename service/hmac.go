// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"crypto/hmac"
	//"crypto/sha256"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
<<<<<<< HEAD
)

const HEADER_DATE = "X-Hmac-Date"
const QUERY_DATE  = "date"

/* What is used in the HMAC?
*	A request comes in that looks like
** 		/register/mydomain/12345?login=we_want_you&pwd=password&email=user@something&name=john_doe&hmac=xxxxxxx
* We use:
*	(1) the shared secret key for a particular user (not transmitted)
*	(2) The path ( /register/mydomain/12345 )
*	(3) The contents of each parm, in alphabetical order.
*			email+user@something+login+we_want_you+name+john_doe+pwd+password
*			HMAC IS NEVER ADDED
*	(4) The date and time stamp.
*			This should be in the headers as "X-Srq-Date"
=======
	"errors"
	"time"
	//"strings"
	"io/ioutil"
)

const HEADER_TIMESTAMP = "Timestamp"
const HEADER_DATE      = "Date"

/* What is used in the HMAC?
*	A request comes in that looks like
* 	GET	/register/mydomain?login=we_want_you&pwd=password&email=user@something&name=john_doe
>>>>>>> FETCH_HEAD
*
*
*	The authorisation information must occur in the header:
		Authorization: abcd-efg-1234-456:qnR8UCqJggD55PohusaBNviGoOJ67HC6Btry4qXLVZc=

	The split occurs at the first colon (:). The left is the caller's ID. The right is the hmac
	generated by hashing the contents of the URI path (/register/mydomain...) and the contents
	of the body of the request.

	The contents of the headers required for authorisation are:
		* The Client ID (always)
		* The Content-MD5 header (if there is a body)
		* The Content-Type in the header (if present)
		* Either header contents of Timestamp the standard HTTP header 'Date:' in the request
		* The complete request (/register/mydomain?login....)
		* The client secret (not sent in request)

	The following should be generated per request:
		Rqst	MD5		Type
		GET		No		No
		PUT		Yes		Yes
		DELETE	No		No
		POST	Yes		Yes
*/

func GetHmacDate( r *http.Request)( string ,error ){

	requestDate := r.Header.Get( HEADER_TIMESTAMP )		// Header has "Timestamp:"
	if len(requestDate) == 0 {					// Umm..NO
		requestDate = r.Header.Get( HEADER_DATE )		// Header has "Date:" ?

<<<<<<< HEAD
	if len(date) == 0 {                               // ... and if it isn't...
		// Oh for heavens sake...
		query := r.URL.Query()                        // Fetch the query list
		date, found = query[QUERY_DATE]               // See if they tucked it in there
		if !found {                                   // ... NOPE - error
			return "", errors.New("No date/time specified for key check")
		}
	}

	h := hmac.New(sha256.New, []byte(secret))         // Start the hmac up
	h.Write([]byte(r.URL.Path))                       // Adding in the full path (command and ID)
	fmt.Printf( "HASH: Add in %s\n" , r.URL.Path )

	for _, key := range keys {                        // for each key (in order)
		if key != KEY_HMAC {                          // The hash can't be part of the hash
			fmt.Printf( "HASH: PROCESS %s\n" , key)
			for _, queryName := range params {
				fmt.Printf("HASH: QueryName is %s\n" , queryName)
				if queryName  == key {
					if val, found := srqst.Get(key); found {
						h.Write([]byte(key))          // Add in the key and ...
						h.Write([]byte(val))          // ... the key value
						fmt.Printf( "HASH: Add in %s%s\n" , key, val )
					}
				}
			}
		}
	}

	h.Write([]byte(date[0]))                        // Add the date exactly as specified
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func CreateHmac(secret string , r *http.Request , srqst *ServiceRequest) ( string, error) {
	var date  []string
	var found bool

	params := srqst.GetQueryKeys()
	keys := srqst.SortedKeys()                        // get a list of the keys in sorted order

	date = r.Header[ HEADER_DATE ]                    // Find the date (should be in header)

	if len(date) == 0 {                               // ... and if it isn't...
		// Oh for heavens sake...
		query := r.URL.Query()                        // Fetch the query list
		date, found = query[QUERY_DATE]               // See if they tucked it in there
		if !found {                                   // ... NOPE - error
			return "", errors.New("No date/time specified for key check")
=======
		if len(requestDate ) == 0 {
			return "" , errors.New( "No date/time specified for key check" )
>>>>>>> FETCH_HEAD
		}
	}
	// Check to see if timestamp is older than 15min. If so, reject request
	// First, parse this into a time object...
	tstamp , err := time.Parse( time.RFC3339 , requestDate )
	if err != nil {
		return "" , err
	}
	now := time.Now()								// Current time...
	diff := now.Sub( tstamp )						// We want how far in the past it is...
	if diff*time.Minute > 15 || diff*time.Minute < -15 {
		return "",errors.New("Time is outside of 15 minutes")
	}
	return requestDate , nil						// Passed all tests...
}

<<<<<<< HEAD
	h := hmac.New(sha256.New, []byte(secret))         // Start the hmac up
	h.Write([]byte(r.URL.Path))                       // Adding in the full path (command and ID)
	fmt.Printf( "HASH: Add in %s\n" , r.URL.Path )

	for _, key := range keys {                        // for each key (in order)
		if key != KEY_HMAC {                          // The hash can't be part of the hash
			fmt.Printf( "HASH: PROCESS %s\n" , key)
			for _, queryName := range params {
				fmt.Printf("HASH: QueryName is %s\n" , queryName)
				if queryName  == key {
					if val, found := srqst.Get(key); found {
						h.Write([]byte(key))          // Add in the key and ...
						h.Write([]byte(val))          // ... the key value
						fmt.Printf( "HASH: Add in %s%s\n" , key, val )
					}
				}
			}
		}
	}


	h.Write([]byte(date[0]))                        // Add the date exactly as specified

	body,err := ioutil.ReadAll( r.Body )
=======
/**
 * 	Return the base64 of the MD5 of the body.
 *  if the body is empty, you will receive
 */
func ComputeBodyMD5( r *http.Request ) string {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}

	fmt.Println( "BODY IS " + r.FormValue("body") )
	if len( body ) == 0 {
		return ""
	}

	d := md5.New()
	d.Write( body )
	m5 := d.Sum(nil)
	return base64.StdEncoding.EncodeToString( m5 )
}

func GenerateHmac( secret string , r *http.Request )( string , error ){

	requestDate, err := GetHmacDate( r )
>>>>>>> FETCH_HEAD
	if err != nil {
		return "" , err
	}

<<<<<<< HEAD
	fmt.Printf("Add body %s\n" , body )
	h.Write( body )
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
=======
	// Now start to generate the rest of the message
	return requestDate, nil
}
func CreateRestfulHmac(secret string , r *http.Request , srqst *ServiceRequest) ( string, error) {
	return "",nil
>>>>>>> FETCH_HEAD
}

func CompareHmac(hmacComputed string , srqst * ServiceRequest) bool {

	if sent, found := srqst.Get("hmac") ; found {
		return hmac.Equal([]byte(sent), []byte(hmacComputed))
	}
	return false
}



