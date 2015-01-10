package service

import (
	//"net/http"
	//"github.com/cgentry/gofig"
	"encoding/json"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/storage"
	"net/http"
)

// ServiceRegister will register a new user into the main store. This will package up the response into a common
// response package after checking the integrity of the request.
func ServiceRegister(store *storage.Store, caller *record.User, requestPackage *record.Package) *record.Package {

	requestHead, packError := serviceGetRequestHead(caller, requestPackage)
	if packError != nil {
		return packError
	}

	responseHead := response.NewHead()
	responseHead.Sequence = requestHead.Sequence

	register := request.NewRegister()
	err := json.Unmarshal([]byte(requestPackage.Body), &register)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidBody)
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
		return serviceReturnStorageError(caller, &responseHead, "", serr)
	}

	returnUserJson, err := json.Marshal(record.NewReturnFromUser(newUser))
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", err)
	}
	return serviceReturnStorageError(caller, &responseHead, string(returnUserJson), storage.ErrStatusOk)

}

// ServiceLogin will Login a user that is registered in the store. This will package up the response into a common
// response package after checking the integrity of the request.
func ServiceLogin(store *storage.Store, caller *record.User, requestPackage *record.Package) *record.Package {

	requestHead, packError := serviceGetRequestHead(caller, requestPackage)
	if packError != nil {
		return packError
	}

	responseHead := response.NewHead()
	responseHead.Sequence = requestHead.Sequence

	login := request.NewLogin()
	err := json.Unmarshal([]byte(requestPackage.Body), &login)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidBody)
	}

	if err = requestHead.Check(); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusNotAcceptable, err.Error())
	}

	// Find the user - we have to use the LOGIN name for this
	user, err := store.FetchUserByLogin(login.Login)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", err)
	}
	// Process the login request. This checks the password that was passed
	status, err := user.Login(login.Password)
	if err != nil {
		return serviceReturnResponse(caller, &responseHead, "", int(status), err.Error())
	}

	serr := store.UserLogin(user)
	if serr != storage.ErrStatusOk {
		return serviceReturnStorageError(caller, &responseHead, "", serr)
	}
	user, err = store.FetchUserByLogin(login.Login)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", err)
	}

	returnUserJson, err := json.Marshal(record.NewReturnFromUser(user))

	if err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusInternalServerError, err.Error())
	}

	return serviceReturnStorageError(caller, &responseHead, string(returnUserJson), storage.ErrStatusOk)

}

func ServiceLogout(store *storage.Store, caller *record.User, requestPackage *record.Package) *record.Package {
	var err error

	requestHead, packError := serviceGetRequestHead(caller, requestPackage)
	if packError != nil {
		return packError
	}

	responseHead := response.NewHead()
	responseHead.Sequence = requestHead.Sequence

	logout := request.NewLogout()
	if err = json.Unmarshal([]byte(requestPackage.Body), &logout); err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidBody)
	}

	if err = requestHead.Check(); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusNotAcceptable, err.Error())
	}

	// Find the user - we have to use the LOGIN name for this
	user, err := store.FetchUserByToken(logout.Token)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", err)
	}

	if serr := store.UserLogout(user); serr != storage.ErrStatusOk {
		return serviceReturnStorageError(caller, &responseHead, "", serr)
	}

	return serviceReturnStorageError(caller, &responseHead, "", storage.ErrStatusOk)
}

// ServiceUpdate will update any field for a user, given a token (user must be logged in)
// This performs general update services for fields
func ServiceUpdate(store *storage.Store, caller *record.User, requestPackage *record.Package) *record.Package {
	var err error
	var dirtyRecord bool = false

	update := request.NewUpdate()
	responseHead := response.NewHead()

	requestHead, OK := requestPackage.Head.(request.Head)
	if !OK {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidHeader)
	}
	if err = requestHead.Check(); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusNotAcceptable, err.Error())
	}
	if !requestPackage.GoodSignature() {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidChecksum)
	}

	responseHead.Sequence = requestHead.Sequence

	if err = json.Unmarshal([]byte(requestPackage.Body), &update); err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidBody)
	}
	if err = update.Check(); err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusNotAcceptable, err.Error())
	}

	// Find the user via Token
	user, err := store.FetchUserByToken(update.Token)
	if err != nil {
		return serviceReturnStorageError(caller, &responseHead, "", err)
	}

	if update.Login != "" {
		if err = user.SetLoginName(update.Login); err != nil {
			return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
		}
		dirtyRecord = true
	}
	if update.Name != "" {
		if err = user.SetName(update.Name); err != nil {
			return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
		}
		dirtyRecord = true
	}
	if update.Email != "" {
		if err = user.SetEmail(update.Email); err != nil {
			return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, err.Error())
		}
		dirtyRecord = true
	}
	if update.OldPassword != "" && update.NewPassword != "" {
		if status := user.ChangePassword(update.OldPassword, update.Token, update.NewPassword); status != record.USER_OK {
			return serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidPasswordOrUser)
		}
		dirtyRecord = true
	}

	if !dirtyRecord {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusBadRequest, "No fields included for update")
	}
	if serr := store.UserUpdate(user); serr != storage.ErrStatusOk {
		return serviceReturnStorageError(caller, &responseHead, "", serr)
	}

	returnUserJson, err := json.Marshal(record.NewReturnFromUser(user))
	if err != nil {
		return serviceReturnResponse(caller, &responseHead, "", http.StatusInternalServerError, err.Error())
	}
	return serviceReturnStorageError(caller, &responseHead, string(returnUserJson), storage.ErrStatusOk)

}

// Return a response package based upon the caller, header, body and status code/message
// This will pack up all the data for a simple response that can be sent using http/rpc/queue
func serviceReturnResponse(caller *record.User, responseHead *response.Head, responseBody string, code int, msg string) *record.Package {
	responseHead.Message = msg
	responseHead.Code = code

	responsePackage := record.NewPackage()
	responsePackage.SetSecret([]byte(caller.GetSalt()))
	responsePackage.SetBodyString(responseBody)
	responsePackage.SetHead(responseHead)
	responsePackage.ClearSecret()

	return responsePackage
}

// Return a response package based upon the caller, header, body and storage error (which contains code/error)
// This will pack up all the data for a simple response that can be sent using http/rpc/queue
func serviceReturnStorageError(caller *record.User, responseHead *response.Head, responseBody string, err error) *record.Package {
	var code int
	if serr, ok := err.(*storage.StorageError); ok {
		code = serr.Code()
	} else {
		code = http.StatusInternalServerError
	}
	return serviceReturnResponse(caller, responseHead, responseBody, code, err.Error())

}

func serviceGetRequestHead(caller *record.User, requestPackage *record.Package) (request.Head, *record.Package) {
	responseHead := response.NewHead()

	if requestPackage == nil || requestPackage.Head == nil {
		requestHead := request.Head{}
		return requestHead, serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidHeader)
	}
	requestHead, OK := requestPackage.Head.(request.Head)
	if !OK {
		return requestHead, serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidHeader)
	}
	if err := requestHead.Check(); err != nil {
		return requestHead, serviceReturnResponse(caller, &responseHead, "", storage.ErrInvalidHeader.Code(), err.Error())
	}
	if !requestPackage.GoodSignature() {
		return requestHead, serviceReturnStorageError(caller, &responseHead, "", storage.ErrInvalidChecksum)
	}
	return requestHead, nil
}
