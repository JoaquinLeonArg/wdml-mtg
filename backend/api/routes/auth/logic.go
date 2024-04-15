package auth

import (
	"errors"
	"net/mail"
	"time"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var ()

func CreateUser(registerRequest RegisterRequest) error {
	if len(registerRequest.Username) < 3 || len(registerRequest.Username) > 32 {
		return apiErrors.ErrUsernameInvalid
	}
	if passwordvalidator.GetEntropy(registerRequest.Password) < 50 {
		return apiErrors.ErrPasswordWeak
	}
	if len(registerRequest.Password) > 64 {
		return apiErrors.ErrPasswordTooLong
	}
	if _, err := mail.ParseAddress(registerRequest.Email); err != nil {
		return apiErrors.ErrEmailInvalid
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 10)
	if err != nil {
		return apiErrors.ErrInternal
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
	if errors.Is(err, db.ErrAlreadyExists) {
		return apiErrors.ErrDuplicatedResource
	}
	return nil
}

func LoginUser(loginRequest LoginRequest) (string, error) {
	user, err := db.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return "", apiErrors.ErrInvalidAuth
	}

	if bcrypt.CompareHashAndPassword(user.Password, []byte(loginRequest.Password)) != nil {
		return "", apiErrors.ErrInvalidAuth
	}

	token, err := CreateToken(user.ID.Hex(), user.Username)
	if err != nil {
		return "", apiErrors.ErrInternal
	}

	return token, nil
}
