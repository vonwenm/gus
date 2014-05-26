// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

// General call errors and returns
const (
	CODE_OK                 = 200

)

// Caller errors are in the 400 range
const (
	CODE_BAD_CALL = 400 + iota
	CODE_USER_DOESNT_EXIST  = iota
	CODE_INVALID_GUID       = iota
	CODE_INVALID_HMAC       = iota
	CODE_INVALID_REQUEST    = iota
	CODE_INVALID_PARAMETERS = iota
)

// Internal errors are in the 500 errors
const (
	CODE_INTERNAL_ERROR = 500 + iota
)

// Database returns
const (
	CODE_DUPLICATE_KEY        = 600 + iota
	CODE_DUPLICATE_EMAIL      = 600 + iota
	CODE_DUPLICATE_LOGIN_NAME = 600 + iota
)
