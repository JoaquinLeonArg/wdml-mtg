package errors

import "fmt"

var (
	ErrInternal = fmt.Errorf("INTERNAL_ERROR")

	ErrInvalidAuth        = fmt.Errorf("INVALID_AUTH")
	ErrDuplicatedResource = fmt.Errorf("DUPLICATED_RESOURCE")
	ErrNotFound           = fmt.Errorf("NOT_FOUND")

	// Auth
	ErrUsernameInvalid = fmt.Errorf("USERNAME_INVALID")
	ErrPasswordWeak    = fmt.Errorf("PASSWORD_WEAK")
	ErrPasswordTooLong = fmt.Errorf("PASSWORD_LONG")
	ErrEmailInvalid    = fmt.Errorf("EMAIL_INVALID")
)
