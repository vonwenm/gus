package service

import (
"crypto/hmac"
"crypto/sha256"
//"encoding/base64"
//"fmt"
	"net/http"
)

/* What is used in the HMAC?
*	A request comes in that looks like
** 		/register/mydomain/12345?login=we_want_you&pwd=password&email=user@something&name=john_doe
* We use:
*	(1) the shared secret key for a particular user (not transmitted)
*	(2) The contents of the PATH statement (/register/mydomain/12345)
*	(3) The contents of each parm, in alphabetical order.
*			email+user@something+login+we_want_you+name+john_doe+pwd+password
*	(4) The date and time stamp.
*			This should be in the headers as x-gus-date
*			OR in the QUERY paramtere: gdate=YYYYMMDDHHMMSS
*
*/
func CreateRestfulHmac( secret string , r * http.Request, parms []string , srqst *ServiceRequest){

	// create a sorted array of keys

	// Start the request up
	h := hmac.New(sha256.New, []byte(secret))
	h.Write( []byte( r.Url.Path ) )
	// In the header, there SHOULD be a date request.

	//

}
