package client

import (
	"fmt"
	"net/http"

	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/auth"
)

// Login sends a request to /auth/login with user credentials,
// then sets the cookie for the ApiClient to authenticate further requests.
func (ac *ApiClient) Login(username, password string) error {
	res, err := ac.post("auth/login", auth.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil || res.StatusCode != http.StatusOK {
		return fmt.Errorf("login request error")
	}
	return nil
}

// Register sends a request to /auth/register, and returns the status code
func (ac *ApiClient) Register(username, email, password string) error {
	res, err := ac.post("auth/register", auth.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("register request error: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		var body []byte
		res.Body.Read(body)
		return fmt.Errorf("register request error code: %v, %v", res.StatusCode, res.Body)
	}
	return nil
}

// CheckLogin sends a request to /auth/check, and returns the status code
func (ac *ApiClient) CheckLogin() error {
	res, err := ac.post("auth/check", nil)
	if err != nil {
		return fmt.Errorf("check request error: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		var body []byte
		res.Body.Read(body)
		return fmt.Errorf("register request error code: %v, %v", res.StatusCode, res.Body)
	}
	return nil
}
