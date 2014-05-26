// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"errors"
	//"bytes"
)

const HEADER_DATE = "X-Gus-Date"
const QUERY_DATE = "date"

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
*			This should be in the headers as "X-Gus-Date"
*
* So, this is:  hmac(  base64(sha256( secret + cmd ) ) + query-cmd + query-value + date )
* The HMAC sent should be encoded in base64 also.
*
*/
func CreateRestfulHmac(secret string , r *http.Request , srqst *ServiceRequest) ( string, error) {
	var date  []string
	var found bool

	params := srqst.GetQueryKeys()
	keys := srqst.SortedKeys()                        // get a list of the keys in sorted order

	date = r.Header[ HEADER_DATE ]                    // Find the date (should be in header)

	if len(date) == 0 {                                    // ... and if it isn't...
		// Oh for heavens sake...
		query := r.URL.Query()                        // Fetch the query list
		date, found = query[QUERY_DATE]                // See if they tucked it in there
		if !found {                                // ... NOPE - error
			return "", errors.New("No date/time specified for key check")
		}
	}

	h := hmac.New(sha256.New, []byte(secret))        // Start the hmac up
	h.Write([]byte(r.URL.Path))                      // Adding in the full path (command and ID)
	fmt.Printf( "HASH: Add in %s\n" , r.URL.Path )

	for _, key := range keys {                        // for each key (in order)
		if key != KEY_HMAC {                          // The hash can't be part of the hash
			fmt.Printf( "HASH: PROCESS %s\n" , key)
			for _, queryName := range params {
				fmt.Printf("HASH: QueryName is %s\n" , queryName)
				if queryName  == key {
					if val, found := srqst.Get(key); found {
						h.Write([]byte(key))                        // Add in the key and ...
						h.Write([]byte(val))                // ... the key value
						fmt.Printf( "HASH: Add in %s%s\n" , key, val )
					}
				}
			}
		}
	}

	h.Write([]byte(date[0]))                        // Add the date exactly as specified
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func CompareHmac(hmacComputed string , srqst * ServiceRequest) bool {

	if sent, found := srqst.Get("hmac") ; found {
		return hmac.Equal([]byte(sent), []byte(hmacComputed))
	}
	return false
}

