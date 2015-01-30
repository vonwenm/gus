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

var ErrHeadNoDomain = NewGeneralError("Head: No domain", http.StatusBadRequest)
var ErrHeadNoId = NewGeneralError("Head: No Id", http.StatusBadRequest)
var ErrHeadNoTimestamp = NewGeneralError("Head: No timestamp set", http.StatusBadRequest)
var ErrHeadFuture = NewGeneralError("Head: Request in the future", http.StatusBadRequest)
var ErrHeadExpired = NewGeneralError("Head: Request expired", http.StatusBadRequest)
var ErrSessionExpired = NewGeneralError("User session expired", http.StatusUnauthorized)
var ErrPasswordTooShort = NewGeneralError("Password is too short", http.StatusBadRequest)
var ErrPasswordTooSimple = NewGeneralError("Password is too simple", http.StatusBadRequest)

// Storage Errors
var ErrInvalidHeader = NewGeneralError("Invalid header in request", http.StatusBadRequest)
var ErrInvalidChecksum = NewGeneralError("Invalid Checksum", http.StatusBadRequest)
var ErrInvalidBody = NewGeneralError("Invalid body (mistmatch request?)", http.StatusBadRequest)
var ErrEmptyFieldForLookup = NewGeneralError("Lookup field is empty", http.StatusBadRequest)
var ErrInvalidPasswordOrUser = NewGeneralError("Invalid password or user id", http.StatusBadRequest)
var ErrMatchAnyNotSupported = NewGeneralError("Storage driver does not support 'MATCH_ANY_DOMAIN' for fetch operation", http.StatusInternalServerError)
var ErrNoDriverFound = NewGeneralError("No storage driver found", http.StatusInternalServerError)
var ErrNoSupport = NewGeneralError("Storage driver does not support function call", http.StatusNotImplemented)
var ErrNotOpen = NewGeneralError("Storage driver is not open", http.StatusInternalServerError)
var ErrAlreadyRegistered = NewGeneralError("Storage driver already registered", http.StatusInternalServerError)
var ErrInternalDatabase = NewGeneralError("Internal storage error while executing operation", http.StatusInternalServerError)

var ErrUserNotFound = NewGeneralError("User not found", http.StatusNotFound)

var ErrInvalidGuid = NewGeneralError("Invalid Guid for lookup", http.StatusNotFound)
var ErrInvalidEmail = NewGeneralError("Invalid email for lookup", http.StatusNotFound)
var ErrInvalidToken = NewGeneralError("Invalid token for lookup", http.StatusNotFound)

var ErrDuplicateGuid = NewGeneralError("User GUID already in use", http.StatusInternalServerError)
var ErrDuplicateEmail = NewGeneralError("Email registered", http.StatusConflict)
var ErrDuplicateLogin = NewGeneralError("Login name already exists", http.StatusConflict)

var ErrUserNotRegistered = NewGeneralError("User not registered", http.StatusBadRequest)
var ErrUserNotLoggedIn = NewGeneralError("User not logged in", http.StatusBadRequest)
var ErrUserLoggedIn = NewGeneralError("User already logged in", http.StatusBadRequest)
var ErrUserNotActive = NewGeneralError("User is not yet activated", http.StatusUnauthorized)

var ErrStatusOk = NewGeneralError("", http.StatusOK)