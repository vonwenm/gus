package service

import (
	//"net/http"
	//"github.com/cgentry/gofig"
	//"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/storage"
	//"github.com/cgentry/gosr"
	"encoding/json"
	"net/http"
)

/*
 * ServiceRegister will create a new record for the user.
 */

func ServiceRegister(caller *record.User, requestPackage *record.Package) *record.Package {
	var responseBody string

	register := request.NewRegister()
	responseHead := response.NewHead()

	requestHead, OK := requestPackage.Head.(request.Head)
	if !OK {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusInternalServerError, "Could not convert head to proper type")
	}
	responseHead.Sequence = requestHead.Sequence

	if !requestPackage.GoodSignature() {
		return serviceReturnResponse(caller, &responseHead, responseBody, http.StatusUnauthorized, "Invalid checksum")
	}

	err := json.Unmarshal([]byte(requestPackage.Body), &register)
	if err != nil {
		return serviceReturnResponse(caller, &responseHead, responseBody, http.StatusBadRequest, err.Error())
	}

	if err = requestHead.Check(); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusNotAcceptable, err.Error())
	}

	// Good domain - save it.
	newUser := record.NewUser()
	if err = newUser.SetDomain(caller.GetDomain()); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
	}
	if err = newUser.SetName(register.Name); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
	}
	if err = newUser.SetEmail(register.Email); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
	}
	if err = newUser.SetPassword(register.Password); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
	}
	/*
		for _, key := range qparam {
			_,err := newUser.MapFieldToUser(key, requestPackage.Parameters.Get(key))
			if err != nil {
				answr.SetError( gosr.NewErrorWithText( CODE_BAD_CALL , err.Error() )).Encode(w)
				return
			}
		}
	*/
	driver := storage.GetDriver()
	driver.RegisterUser(newUser)

	returnUserJson, err := json.Marshal(record.NewReturnFromUser(newUser))
	if err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusInternalServerError, err.Error())
	}
	return serviceReturnResponse(caller, &responseHead, string(returnUserJson), http.StatusOK, "")

}
func serviceReturnResponse(caller *record.User, responseHead *response.Head, responseBody string, code int, msg string) *record.Package {
	responseHead.Message = msg
	responseHead.Code = code

	responsePackage := record.NewPackage()
	responsePackage.SetSecret([]byte(caller.GetSalt()))
	responsePackage.SetBodyString(responseBody)
	responsePackage.SetHead(responseHead)

	return responsePackage
}

/*
func ServiceRegister(c *gofig.Configuration , w http.ResponseWriter, r *http.Request , parseParms ParseParms) {

	var newUser * record.User
	caller , rqst , answr := parseParms(c, w, r)                        // Standardise format

	if caller != nil {
		qparam := []string{ KEY_EMAIL, KEY_LOGIN, KEY_NAME, KEY_PWD}    // Check normal fields
		if err := requestPackage.Parameters.IsPresent(qparam); err != nil {       // .. if not present
			answr.SetError(err).Encode(w)                               // .. send back and error
		}else {                                                         // All fields preset - register
			user := record.NewUser( )                					// Fetch the user's domain
			user.SetDomain( caller.GetDomain() )
			if err := user.Unmarshall(rqst.Content.GetContent()); err != nil {     // Pass in the body of the request
				answr.SetErrorAndCode(err, http.StatusPreconditionFailed).Encode(w)
			}else {                                                      // Good domain - save it.
				newUser = record.NewUser()
				newUser.SetDomain(caller.GetDomain())
				for _, key := range qparam {
					_,err := newUser.MapFieldToUser(key, rqst.Parameters.Get(key))
					if err != nil {
						answr.SetError( gosr.NewErrorWithText( CODE_BAD_CALL , err.Error() )).Encode(w)
						return
					}
				}

				driver := storage.GetDriver()
				driver.RegisterUser(user)

				userReturn := record.NewReturnFromUser(newUser)

				ReturnUserJson(w, CODE_OK, &userReturn)
			}
		}
	}
	return
}
*/
