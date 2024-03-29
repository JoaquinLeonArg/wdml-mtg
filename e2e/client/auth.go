package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/auth"
)

// Register sends a request to auth/login with user credentials, and returns the status code
func (ac *ApiClient) Login(username, password string) (string, error) {
	res, err := ac.get("auth/login")
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("login request error")
	}
	cookie, err := res.Request.Cookie("jwt")
	if err != nil {
		return "", fmt.Errorf("jwt not found on login response: %w", err)
	}
	return cookie.Value, nil
}

// Login sends a request to auth/login with user credentials, and returns the status code
func (ac *ApiClient) Register(username, email, password string) error {
	body, err := json.Marshal(auth.RegisterRequest)
	if err != nil {
		return fmt.Errorf("register request error")
	}
	res, err := ac.post("auth/register", body)
	if err != nil || res.StatusCode != http.StatusOK {
		return fmt.Errorf("register request error")
	}
	return nil
}
