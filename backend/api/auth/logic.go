package auth

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameInvalid = fmt.Errorf("USERNAME_INVALID")
	ErrPasswordWeak    = fmt.Errorf("PASSWORD_WEAK")
	ErrPasswordTooLong = fmt.Errorf("PASSWORD_LONG")
	ErrEmailInvalid    = fmt.Errorf("EMAIL_INVALID")
	ErrInvalidAuth     = fmt.Errorf("INVALID_AUTH")
	ErrInternal        = fmt.Errorf("INTERNAL_ERROR")
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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 10)
	if err != nil {
		return ErrInternal
	}

	err = db.CreateUser(domain.User{
		Username:          registerRequest.Username,
		Email:             registerRequest.Email,
		Password:          passwordHash,
		Description:       "New user!",
		CreatedAt:         primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:         primitive.NewDateTimeFromTime(time.Now()),
		ProfilePictureURL: "",
	})
	return err
}

func LoginUser(loginRequest LoginRequest) (string, error) {
	user, err := db.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return "", ErrInvalidAuth
	}

	if bcrypt.CompareHashAndPassword(user.Password, []byte(loginRequest.Password)) != nil {
		return "", ErrInvalidAuth
	}

	token, err := CreateToken(user.ID.String(), user.Username)
	if err != nil {
		return "", ErrInternal
	}

	return token, nil
}
