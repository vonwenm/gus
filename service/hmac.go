// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	//"fmt"
	"net/http"
	"errors"
	"bytes"
)

const HEADER_DATE = "X-Gus-Date"
const QUERY_DATE = "date"

/* What is used in the HMAC?
*	A request comes in that looks like
** 		/register/mydomain/12345?login=we_want_you&pwd=password&email=user@something&name=john_doe
* We use:
*	(1) the shared secret key for a particular user (not transmitted)
*	(2) The Command, hashed with the secret (sha256)
*	(3) The contents of each parm, in alphabetical order.
*			email+user@something+login+we_want_you+name+john_doe+pwd+password
*	(4) The date and time stamp.
*			This should be in the headers as "X-Gus-Date"
*
* So, this is:  hmac(  base64(sha256( secret + cmd ) ) + query-cmd + query-value + date )
* The HMAC sent should be encoded in base64 also.
*
*/
func CreateRestfulHmac(secret string , r *http.Request , srqst *ServiceRequest) ( string, error) {

	var found bool
	var date  []string

	cmds := srqst.GetPathKeys()                        // Ones we should skip...
	cmd, _ := srqst.Get("cmd")                        // Standard command

	// First, a simple hash rather than the hmac for the command string
	s := sha256.New()                                // Quick hash of the CMD
	s.Write([]byte(secret))                        // Add in the secret...
	s.Write([]byte(cmd))                            // .. then the command
	cmd = base64.StdEncoding.EncodeToString(s.Sum(nil))

	h := hmac.New(sha256.New, []byte(secret))        // Start the hmac up
	h.Write([]byte(cmd))                        // Adding in the fresh command hash


	date = r.Header[ HEADER_DATE ]                    // Find the date (should be in header)
	if len(date) == 0 {                                    // ... and if it isn't...
		// Oh for heavens sake...
		query := r.URL.Query()                        // Fetch the query list
		date, found = query[QUERY_DATE]                // See if they tucked it in there
		if !found {                                // ... NOPE - error
			return "", errors.New("No date/time specified for key check")
		}
	}

	keys := srqst.SortedKeys()                        // get a list of the keys in sorted order

	for _, key := range keys {                        // for each key (in order)
		found = false
		for _, path := range cmds {
			found = (bytes.Compare([]byte(path), []byte(key)) == 0 )
			if found {
				break
			}
		}
		if !found {
			if val, found := srqst.Get(key); found {
				h.Write([]byte(key))                        // Add in the key and ...
				h.Write([]byte(val))                // ... the key value
			}
		}
	}

	h.Write([]byte(date[0]))                        // Add the date exactly as specified
	return base64.StdEncoding.EncodeToString(h.Sum(nil)) , nil
}

func CompareHmac( hmacComputed string , srqst * ServiceRequest) bool {

	if sent, found := srqst.Get("hmac") ; found {
		return hmac.Equal([]byte(sent), []byte(hmacComputed))
	}
	return false
}

