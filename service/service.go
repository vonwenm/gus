// Package service is the main entry point for processing requests. Any external interface should use these calls
// to make changes to the user store.

package service

import (
	"encoding/json"
	"fmt"
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/record/head"
	"github.com/cgentry/gus/record/mappers"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/storage"
	"net/http"
	"strings"
)

// Permissions for updating
const (
	SERVICE_EMPTY_BODY = ""

	PERMIT_ALL      = "permit_all"
	PERMIT_LOGIN    = "permit_login"
	PERMIT_PASSWORD = "permit_password"
	PERMIT_NAME     = "permit_name"
	PERMIT_EMAIL    = "permit_email"
)

// Service creator is any function that returns a pointer to a ServiceProcess.
type ServiceCreator func() *ServiceProcess

type ServiceProcessor interface {
	Setup(*configure.Configure, string) *record.Package
	Teardown() error
	GeneralError(string, error) *record.Package
	Response(string, error) *record.Package
}

// All of the service control requirements are stored in this structure. This points to the
// runtime function that will receive this information.
type ServiceProcess struct {
	// Points to the function that will process the request. It will be passed this structure.
	Run func(*ServiceProcess) *record.Package

	// Point to the configuration class. All configuration is held here.
	Config *configure.Configure

	// The request head is the unpacked header that was decoded from the incoming
	// request package.
	RequestHead *head.Head
	// The request body, unpacked for the service record
	RequestBody record.BodyInterface
	// Client record making the request. This can come from the user or client database
	Client *tenant.User

	// Header for response we are sending back.
	ResponseHead *head.Head
	// Record package - this is what will be returned from any of the calls.
	ResponsePackage record.Packer

	// Datastore for user records. The clientstore, if separate, is not stored as it is
	// only needed once to access the client record.
	UserStore storage.Storer

	// Options can be used to set any options that are desired for the routines.
	Options map[string]string

	SetFlag bool
}

func NewServiceRegister() *ServiceProcess {
	r := &ServiceProcess{
		Run:         register,
		RequestBody: &request.Register{},
		SetFlag:     false,
	}
	return r.Reset()
}

// The Structure that gives us the entry point for user Login
func NewServiceLogin() *ServiceProcess {
	r := &ServiceProcess{
		Run:         login,
		RequestBody: &request.Login{},
	}
	return r.Reset()
}

// The Structure that gives us the entry point for user Logout
func NewServiceLogout() *ServiceProcess {
	r := &ServiceProcess{
		Run:         logout,
		RequestBody: &request.Logout{},
	}
	return r.Reset()
}

// The Structure that gives us the entry point for user record updates
func NewServiceUpdate() *ServiceProcess {
	r := &ServiceProcess{
		Run:         update,
		RequestBody: &request.Update{},
	}
	return r.Reset()
}

// The Structure that gives us the entry point for user record updates
func NewServiceAuthenticate() *ServiceProcess {
	r := &ServiceProcess{
		Run:         authenticate,
		RequestBody: &request.Authenticate{},
	}
	return r.Reset()

}

// NewServiceTest is the entry point for a client checking a connection
func NewServiceTest() *ServiceProcess {
	r := &ServiceProcess{
		Run:         servicetest,
		RequestBody: &request.Test{},
	}
	return r.Reset()

}

// Setup the service structure for common values required.  This will take the request package and
// unpack it into the header and service-specific body.
func (s *ServiceProcess) SetupService(c *configure.Configure, requestPackage string) record.Packer {
	var err error

	// Unpack the incoming request, saving the body and header in our structure.
	// ensure the package has all of the required elements.
	pack := record.NewPackage()
	if err = json.Unmarshal([]byte(requestPackage), pack); err != nil {
		return s.ResponseCode(err.Error(), ecode.ErrBadPackage)
	}

	if !pack.IsPackageComplete() {
		return s.ResponseCode(`"{"Message": "Package is not complete"}`, ecode.ErrBadPackage)
	}
	s.RequestHead, _ = pack.GetHead().(*head.Head)
	s.ResponseHead.Sequence = s.RequestHead.Sequence
	s.ResponseHead.Id = s.RequestHead.Id
	if err = s.RequestHead.Check(); err != nil {
		return s.ResponseCode(`{"Message":"Header is not complete."}`, err)
	}

	// Open up the storage handles. We only need the UserStore as the ClientStore is transitory
	// and only used to read in the client record.
	s.UserStore, err = storage.Open(s.Config.User.Name, s.Config.User.Dsn, s.Config.User.Options)
	if err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}
	s.SetFlag = true
	if s.Config.Service.ClientStore {
		var clientStore *storage.Store
		clientStore, err = storage.Open(s.Config.User.Name, s.Config.User.Dsn, s.Config.User.Options)
		if err == nil {
			s.Client, err = clientStore.FetchUserByLogin(s.RequestHead.Domain, s.RequestHead.Id)
			if err != nil {
				clientStore.Release()
				clientStore.Close()
			}
		}
	} else {
		s.Client, err = s.UserStore.FetchUserByLogin(s.RequestHead.Domain, s.RequestHead.Id)
		s.UserStore.Release()
	}
	if err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}
	s.ResponsePackage.SetSecret([]byte(s.Client.Salt))

	// Confirm that the signature is good. We wait here so we can use the client record.
	if !record.GoodSignature(pack) {
		return s.PackageErr(ecode.ErrInvalidChecksum)
	}
	// Unpack the body. The body is defined as an interface, so we can do a check here.
	if err = json.Unmarshal([]byte(pack.GetBody()), s.RequestBody); err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, ecode.ErrBadBody)
	}
	if err = s.RequestBody.Check(); err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}

	fmt.Printf("USER STORE: %+v\n", s.UserStore)
	s.ResponseOk("")
	return nil
}

// Allocate storage for all of the data in the structure. This will "reset" the storage
// and let the service be re-used.
func (s *ServiceProcess) Reset() *ServiceProcess {
	s.ResponseHead = head.New()
	s.ResponsePackage = record.NewPackage()
	s.Options = make(map[string]string)
	return s
}

// Perform any cleanup needed, closing any connections.
func (s *ServiceProcess) Teardown() error {
	if s.UserStore != nil && s.SetFlag {
		s.UserStore.Release()
		s.UserStore.Close()
		s.SetFlag = false
		s.UserStore = nil
	}
	return nil
}

// test will simply check that a package is correctly formatted, the key is valid and the
// timestamp is correct. It then sends an OK message back. Most of the data checking is
// done already in "Setup"
func servicetest(s *ServiceProcess) *record.Package {
	s.ResponsePackage.SetBodyMarshal( response.NewAck(`test`) )
	return s.PackageOk()
}

// register will register a new user into the main s.UserStore. This will package up the response into a common
// response package after checking the integrity of the request.
func register(s *ServiceProcess) *record.Package {
	var err error
	var eUpdate mappers.ErrSetter

	request, ok := s.RequestBody.(*request.Register)
	if !ok {
		return s.ResponseCode(SERVICE_EMPTY_BODY, ecode.ErrBadBody)
	}
	newUser := tenant.NewUser()
	eUpdate.Set(newUser.SetDomain, s.Client.Domain)
	eUpdate.Set(newUser.SetEmail, request.Email)
	eUpdate.Set(newUser.SetLoginName, request.Login)
	eUpdate.Set(newUser.SetName, request.Name)
	if err = eUpdate.Set(newUser.SetPassword, request.Password); err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}

	if err = s.UserStore.UserInsert(newUser); err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}

	err = s.ResponsePackage.SetBodyMarshal(mappers.ResponseFromUser(response.NewUserReturn(), newUser) )
	if err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}
	return s.PackageOk()

}

// ServiceLogin will Login a user that is registered in the s.UserStore. This will package up the response into a common
// response package after checking the integrity of the request.
func login(s *ServiceProcess) *record.Package {
	login := s.RequestBody.(*request.Login)
	if s.UserStore == nil {
		panic("The userstore is nil")
	}

	defer s.UserStore.Release()
	user, err := s.UserStore.FetchUserByLogin(s.Client.Domain, login.Login)
	if err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}
	// Process the login request. This checks the password that was passed
	if err = user.Login(login.Password); err != nil {
		s.UserStore.UserUpdate(user) // Try and save the error counters
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}

	if err = s.UserStore.UserUpdate(user); err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}
	err = s.ResponsePackage.SetBodyMarshal(mappers.ResponseFromUser(response.NewUserReturn(), user) )
	if err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}
	return s.PackageOk()
}

// ServiceLogout will logout the user that is currently logged in. Only the token is required for this operation.
// If the user is not logged in then an error will be returned. If a user isn't found, a 'NotLoggedIn'
// is returned instead. This is a more precise message for a logout condition
func logout(s *ServiceProcess) *record.Package {
	var err error

	logout, _ := s.RequestBody.(*request.Logout)

	// Find the user - we have to use the TOKEN name for this
	user, err := s.UserStore.UserFetch(s.Client.Domain, storage.FIELD_TOKEN, logout.Token)
	if err != nil {
		if err == ecode.ErrUserNotFound {
			return s.ResponseCode("", ecode.ErrUserNotLoggedIn)
		}

		return s.ResponseCode("", err)
	}
	defer s.UserStore.Release()

	if err = user.Logout(); err != nil {
		return s.ResponseCode("", err)
	}

	if err = s.UserStore.UserUpdate(user); err != nil {
		return s.ResponseCode("", err)
	}
	err = s.ResponsePackage.SetBodyMarshal( response.NewAck(`logout`) )
	if err != nil {
		return s.ResponseCode(SERVICE_EMPTY_BODY, err)
	}
	return s.PackageOk()
}

// Authenticate will check to see if the user is logged in and then mark the record as updated. This should
// only be called about once a minute by the client so they
func authenticate(s *ServiceProcess) *record.Package {
	var err error

	auth, _ := s.RequestBody.(*request.Authenticate)

	// Find the user - we have to use the TOKEN name for this
	user, err := s.UserStore.UserFetch(s.Client.Domain, storage.FIELD_TOKEN, auth.Token)
	if err != nil {
		if err == ecode.ErrUserNotFound {
			return s.ResponseCode("", ecode.ErrUserNotLoggedIn)
		}
		return s.ResponseCode("", err)
	}
	defer s.UserStore.Release()

	if err = user.Authenticate(auth.Token); err != nil {
		return s.ResponseCode("", err)
	}

	if err = s.UserStore.UserUpdate(user); err != nil {
		return s.ResponseCode("", err)
	}
	s.ResponsePackage.SetBodyMarshal( response.NewAck(`logout`) )
	return s.PackageOk()
}

// Update is the catch-all for updating the record. The fields that can be updated through THIS call
// are: LoginName, FullName, Email and Password. This limited set allows most front-end applications to
// alter key fields that the user will want to affect. It is only accessible by the users' token, so they
// must be logged in currently.
//
// If a field is blank, the field will not be updated. This allows the front-end to control what is being altered.
//
// If a front-end wants to create multiple interfaces (change password only, for example) it can include options
// in the call which will stop updates from occurring.
func update(s *ServiceProcess) *record.Package {
	var err error
	var eSetter mappers.ErrSetter
	var updatedFields []string

	update := s.RequestBody.(*request.Update)

	if s.Options == nil || len(s.Options) == 0 {
		return s.Response("", http.StatusInternalServerError, "No updates in options")
	}

	// Find the user via Token
	user, err := s.UserStore.UserFetch(s.Client.Domain, storage.FIELD_TOKEN, update.Token)
	if err != nil {
		return s.ResponseCode("", err)
	}
	defer s.UserStore.Release()

	if update.Login != "" && (s.boolOption(PERMIT_ALL) || s.boolOption(PERMIT_LOGIN)) {
		eSetter.Set(user.SetLoginName, update.Login)
		updatedFields = append(updatedFields, "Login")
	}
	if update.Name != "" && (s.boolOption(PERMIT_ALL) || s.boolOption(PERMIT_NAME)) {
		eSetter.Set(user.SetName, update.Name)
		updatedFields = append(updatedFields, "Name")
	}
	if update.Email != "" && (s.boolOption(PERMIT_ALL) || s.boolOption(PERMIT_EMAIL)) {
		eSetter.Set(user.SetEmail, update.Email)
		updatedFields = append(updatedFields, "Email")
	}
	if eSetter.Err != nil {
		return s.ResponseCode("", eSetter.Err)
	}
	if update.OldPassword != "" && update.NewPassword != "" && (s.boolOption(PERMIT_ALL) || s.boolOption(PERMIT_PASSWORD)) {
		if err = user.ChangePassword(update.OldPassword, update.NewPassword); err != nil {
			return s.ResponseCode("", err)
		}
		updatedFields = append(updatedFields, "Password")
	}
	if len(updatedFields) == 0 {
		return s.Response("", http.StatusBadRequest, "No fields included for update")
	}
	if err = s.UserStore.UserUpdate(user); err != nil {
		return s.ResponseCode("", err)
	}

	returnUserJson, err := json.Marshal(mappers.ResponseFromUser(response.NewUserReturn(), user))
	if err != nil {
		return s.Response("", http.StatusInternalServerError, err.Error())
	}
	return s.Response(string(returnUserJson), ecode.ErrStatusOk.Code(), "Fields updated: "+strings.Join(updatedFields, `, `))
}
func (s *ServiceProcess) boolOption(key string) bool {
	_, ok := s.Options[key]
	return ok
}

/*
 *            RESPONSE ROUTINES:
 *				package up any responses and return to caller
 */

func ( s *ServiceProcess) PackageOk() record.Packer {
	record.SignPackage( s.ResponsePackage )
	return s.ResponsePackage
}


func ( s *ServiceProcess ) PackageErr( err Error ) record.Packer {
	s.ResponsePackage.SetBodyMarshal( err )
	return s.ResponsePackage
}
// Return a response package based upon the caller, header, body and status code/message
// This will pack up all the data for a simple response that can be sent using http/rpc/queue
func (s *ServiceProcess) Response(jsonResponse string, returnCode int, statusMsg string) *record.Package {
	s.ResponsePackage.SetBodyString(jsonResponse)
	s.ResponsePackage.SetHead(s.ResponseHead)
	record.SignPackage(s.ResponsePackage)
	s.ResponsePackage.ClearSecret()

	return s.ResponsePackage
}

// Return a response package based upon the caller, header, body and storage error (which contains code/error)
// This will pack up all the data for a simple response that can be sent using http/rpc/queue
func (s *ServiceProcess) ResponseCode(jsonResponseBody string, err error) *record.Package {
	if err == nil {
		return s.ResponseCode(jsonResponseBody, ecode.ErrStatusOk)
	}
	var code int
	if serr, ok := err.(*ecode.GeneralError); ok {
		code = serr.Code()
	} else {
		code = http.StatusInternalServerError
	}
	return s.Response(jsonResponseBody, code, err.Error())
}

// A convenience function to simply return an OK response with data.
func (s *ServiceProcess) ResponseOk(jsonReponseBody string) *record.Package {
	return s.ResponseCode(jsonReponseBody, ecode.ErrStatusOk)
}
