package storage


const CODE_OK = 0
const CODE_INVALID_GUID = 1
// 500 range errors
const (
	CODE_INTERNAL_ERROR = 500+iota
)

const (
	CODE_DUPLICATE_KEY  = 600+iota
	CODE_DUPLICATE_EMAIL = 600+iota
	CODE_DUPLICATE_LOGIN_NAME = 600+iota
)
