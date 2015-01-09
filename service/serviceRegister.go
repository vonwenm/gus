package service

import (
	//"net/http"
	//"github.com/cgentry/gofig"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/storage"
	//"github.com/cgentry/gosr"
	"encoding/json"
	"net/http"
	"fmt"
)

// ServiceRegister will register a new user into the main store. This will package up the response into a common
// response package after checking the integrity of the request.
func ServiceRegister(store *storage.Store, caller *record.User, requestPackage *record.Package) *record.Package {

	register := request.NewRegister()
	responseHead := response.NewHead()

	requestHead, OK := requestPackage.Head.(request.Head)
	if !OK {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidHeader )
	}
	responseHead.Sequence = requestHead.Sequence

	if !requestPackage.GoodSignature() {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidChecksum )
	}

	err := json.Unmarshal([]byte(requestPackage.Body), &register)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidBody )
	}

	if err = requestHead.Check(); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusNotAcceptable, err.Error())
	}

	// Good domain - save it.
	newUser := record.NewUser()
	if err = newUser.SetDomain(caller.GetDomain()); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
	}
	if err = newUser.SetLoginName(register.Login); err != nil {
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

	serr := store.RegisterUser(newUser)
	if serr != storage.ErrStatusOk {
		fmt.Println( err )
		return serviceReturnStorageError(caller, &responseHead, "", serr)
	}

	returnUserJson, err := json.Marshal(record.NewReturnFromUser(newUser))
	if err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusInternalServerError, err.Error())
	}
	return serviceReturnResponse(caller, &responseHead, string(returnUserJson), http.StatusOK, "")

}

// ServiceLogin will Login a user that is registered in the store. This will package up the response into a common
// response package after checking the integrity of the request.
func ServiceLogin(store *storage.Store, caller *record.User, requestPackage *record.Package) *record.Package {

	register := request.NewRegister()
	responseHead := response.NewHead()

	requestHead, OK := requestPackage.Head.(request.Head)
	if !OK {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidHeader )
	}
	responseHead.Sequence = requestHead.Sequence

	if !requestPackage.GoodSignature() {
		return  serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidChecksum )
	}

	err := json.Unmarshal([]byte(requestPackage.Body), &register)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidBody )
	}

	if err = requestHead.Check(); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusNotAcceptable, err.Error())
	}

	// Good domain - save it.
	user,err := store.FetchUserByLogin( register.Login )
	if err != nil {
		return  serviceReturnStorageError(caller, &responseHead, "", err )
	}
	status, err := user.Login(register.Password)
	if err != nil {
		return serviceReturnResponse(caller, &responseHead, "", int(status), err.Error())
	}

	serr := store.UserLogin(user)
	if serr != storage.ErrStatusOk {
		return serviceReturnStorageError(caller, &responseHead, "", serr)
	}
	user,err = store.FetchUserByLogin( register.Login )
	if err != nil {
		return  serviceReturnStorageError(caller, &responseHead, "", err )
	}

	returnUserJson, err := json.Marshal(record.NewReturnFromUser(user))

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

func serviceReturnStorageError(caller *record.User, responseHead *response.Head, responseBody string, err error) *record.Package {
	var code int
	if serr,ok := err.(*storage.StorageError); ok {
		code = serr.Code()
	}else{
		code = http.StatusInternalServerError
	}
	return serviceReturnResponse(caller,responseHead,responseBody,code,err.Error())

}

