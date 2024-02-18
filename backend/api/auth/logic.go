package auth

import (
	"fmt"
	"net/mail"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var (
	ErrUsernameInvalid = fmt.Errorf("USERNAME_INVALID")
	ErrPasswordWeak    = fmt.Errorf("PASSWORD_WEAK")
	ErrPasswordTooLong = fmt.Errorf("PASSWORD_LONG")
	ErrEmailInvalid    = fmt.Errorf("EMAIL_INVALID")
)

func CreateUser(registerRequest RegisterRequest) error {
	if len(registerRequest.Username) < 3 || len(registerRequest.Username) > 12 {
		return ErrUsernameInvalid
	}
	if passwordvalidator.GetEntropy(registerRequest.Password) < 50 {
		return ErrPasswordWeak
	}
	if len(registerRequest.Password) > 64 {
		return ErrPasswordTooLong
	}
	if _, err := mail.ParseAddress(registerRequest.Email); err != nil {
		return ErrEmailInvalid
	}
	err := db.CreateUser(domain.User{
		Username:    registerRequest.Username,
		Email:       registerRequest.Email,
		Password:    registerRequest.Password,
		Description: "New user!",
		// CreatedAt:
		// UpdatedAt:
		// ProfilePictureURL: "",
	})
	return err
}
