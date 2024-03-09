package errors

import "fmt"

var (
	ErrInternal = fmt.Errorf("INTERNAL_ERROR")

	ErrInvalidAuth        = fmt.Errorf("INVALID_AUTH")
	ErrDuplicatedResource = fmt.Errorf("DUPLICATED_RESOURCE")
	ErrNotFound           = fmt.Errorf("NOT_FOUND")
	ErrUnauthorized       = fmt.Errorf("UNAUTHORIZED")
	ErrNoData             = fmt.Errorf("NO_DATA")

	// Auth
	ErrUsernameInvalid = fmt.Errorf("USERNAME_INVALID")
	ErrPasswordWeak    = fmt.Errorf("PASSWORD_WEAK")
	ErrPasswordTooLong = fmt.Errorf("PASSWORD_LONG")
	ErrEmailInvalid    = fmt.Errorf("EMAIL_INVALID")
)
