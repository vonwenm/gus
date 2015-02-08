package ecode

// Common errors for the GUS system are defined here, in one place. When
// adding new drivers or new portions to the system, it is advised that you use these
// messages, where possible. This will make the client code easier to maintain
//
// All messages have a common format: a code and a message string. The codes give
// a crude type of response while the messages give a detailed error.
//
import (
	"net/http"
)

type ErrorCoder interface {
	Error() string
	Code() int
}

/* ------------------------------------------------------------------------------------ */

type GeneralError struct {
	errorString string
	errorCode   int
}

func NewGeneralError(msg string, code int) ErrorCoder {
	return &GeneralError{errorString: msg, errorCode: code}
}

func NewGeneralFromError(e error, code int) ErrorCoder {
	if e == nil {
		return nil
	}
	return &GeneralError{errorString: e.Error(), errorCode: code}

}

func (s *GeneralError) Error() string { return s.errorString }
func (s *GeneralError) Code() int     { return s.errorCode }

var ErrBadPackage         = NewGeneralError("Package: Bad format" , http.StatusBadRequest )
var ErrBadBody            = NewGeneralError("Package: Cannot unarshal body", http.StatusBadRequest )

var ErrHeadNoDomain       = NewGeneralError("Head: No domain", http.StatusBadRequest)
var ErrHeadNoId           = NewGeneralError("Head: No Id", http.StatusBadRequest)
var ErrHeadNoTimestamp    = NewGeneralError("Head: No timestamp set", http.StatusBadRequest)
var ErrHeadFuture         = NewGeneralError("Head: Request in the future", http.StatusBadRequest)
var ErrHeadExpired        = NewGeneralError("Head: Request expired", http.StatusBadRequest)


var ErrRequestNoTimestamp = NewGeneralError( "Request: No timestamp set", http.StatusBadRequest)
var ErrRequestFuture      = NewGeneralError( "Request: Request in the future", http.StatusBadRequest)
var ErrRequestExpired     = NewGeneralError( "Request: Request expired", http.StatusBadRequest)
var ErrMissingLogin       = NewGeneralError( "Request: Missing login" , http.StatusBadRequest )
var ErrMissingName        = NewGeneralError( "Request: Missing Name" , http.StatusBadRequest )
var ErrMissingPassword    = NewGeneralError( "Request: Missing Password" , http.StatusBadRequest )
var ErrMissingToken       = NewGeneralError( "Request: Missing token" , http.StatusBadRequest )
var ErrMissingEmail       = NewGeneralError( "Request: Missing Email" , http.StatusBadRequest )
var ErrMissingPasswordNew = NewGeneralError( "Request: Missing New Password" , http.StatusBadRequest )
var ErrMatchingPassword   = NewGeneralError( "Request: Old and new passwords match" , http.StatusBadRequest )
var ErrPasswordTooShort   = NewGeneralError( "Request: Password is too short", http.StatusBadRequest)

var ErrSessionExpired     = NewGeneralError( "User session expired", http.StatusUnauthorized)
var ErrPasswordTooSimple  = NewGeneralError( "Password is too simple", http.StatusBadRequest)

// Storage Errors
var ErrInvalidHeader      = NewGeneralError("Invalid header in request", http.StatusBadRequest)
var ErrInvalidChecksum    = NewGeneralError("Invalid Checksum", http.StatusBadRequest)
var ErrInvalidBody        = NewGeneralError("Invalid body (mistmatch request?)", http.StatusBadRequest)
var ErrEmptyFieldForLookup = NewGeneralError("Lookup field is empty", http.StatusBadRequest)
var ErrInvalidPasswordOrUser = NewGeneralError("Invalid password or user id", http.StatusBadRequest)
var ErrMatchAnyNotSupported = NewGeneralError("Storage driver does not support 'MATCH_ANY_DOMAIN' for fetch operation", http.StatusInternalServerError)
var ErrNoDriverFound      = NewGeneralError("No storage driver found", http.StatusInternalServerError)
var ErrNoSupport          = NewGeneralError("Storage driver does not support function call", http.StatusNotImplemented)
var ErrNotOpen            = NewGeneralError("Storage driver is not open", http.StatusInternalServerError)
var ErrAlreadyRegistered  = NewGeneralError("Storage driver already registered", http.StatusInternalServerError)
var ErrInternalDatabase   = NewGeneralError("Internal storage error while executing operation", http.StatusInternalServerError)
var ErrCannotSetId        = NewGeneralError("User id cannot be set", http.StatusBadRequest)
var ErrUserNotFound       = NewGeneralError("User not found", http.StatusNotFound)

var ErrShortGuid          = NewGeneralError( "GUID must be at least 32 characters long", http.StatusInternalServerError)
var ErrDuplicateGuid      = NewGeneralError( "User GUID already in use", http.StatusInternalServerError)
var ErrInvalidGuid        = NewGeneralError( "Invalid Guid for lookup", http.StatusNotFound)
var ErrInvalidEmail       = NewGeneralError( "Invalid email for lookup", http.StatusNotFound)
var ErrInvalidToken       = NewGeneralError( "Invalid token for lookup", http.StatusNotFound)


var ErrDuplicateEmail     = NewGeneralError( "Email address already registered", http.StatusConflict)
var ErrDuplicateLogin     = NewGeneralError( "Login name already exists", http.StatusConflict)

var ErrUserNotRegistered  = NewGeneralError("User not registered", http.StatusBadRequest)
var ErrUserNotLoggedIn    = NewGeneralError("User not logged in", http.StatusBadRequest)
var ErrUserLoggedIn       = NewGeneralError("User already logged in", http.StatusBadRequest)
var ErrUserNotActive      = NewGeneralError("User is not yet activated", http.StatusUnauthorized)

var ErrStatusOk           = NewGeneralError("", http.StatusOK)
