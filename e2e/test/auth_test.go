package test

import (
	"testing"

	"github.com/joaquinleonarg/wdml-mtg/e2e/client"
	"github.com/joaquinleonarg/wdml-mtg/e2e/config"
	"github.com/joaquinleonarg/wdml-mtg/e2e/pkg/mongo"
	"github.com/stretchr/testify/require"
)

// TestRegisterAndLoginOk registers a valid user, then logs it in
func TestRegisterAndLoginOk(t *testing.T) {
	c := client.NewApiClient(config.Config.APIBaseURL)
	mongo.Cleanup()

	username := "test_user"
	email := "valid@email.com"
	password := "th!s_1s_@_s3cure_pASSw0rd"

	// Register the user
	err := c.Register(username, email, password)
	require.NoError(t, err)

	// Login with that user
	err = c.Login(username, password)
	require.NoError(t, err)

	// Check if the cookie was set
	err = c.CheckLogin()
	require.NoError(t, err)
}

// TestRegisterAndLoginWrongCredentials registers a valid user, then logs in with the wrong credentials
func TestRegisterAndLoginWrongCredentials(t *testing.T) {
	c := client.NewApiClient(config.Config.APIBaseURL)
	mongo.Cleanup()

	username := "test_user"
	wrong_username := "miguel"
	email := "valid@email.com"
	password := "th!s_1s_@_s3cure_pASSw0rd"
	wrong_password := "th!s_1s_n0t_th3_pASSw0rd"

	// Register the user
	err := c.Register(username, email, password)
	require.NoError(t, err)

	// Login with the wrong password
	err = c.Login(username, wrong_password)
	require.Error(t, err)

	// Login with the wrong username
	err = c.Login(wrong_username, password)
	require.Error(t, err)

	// Login with the email instead of username
	err = c.Login(email, password)
	require.Error(t, err)
}

// TestRegisterDuplicatedUser registers a valid user, then tries to reuse username/email
func TestRegisterDuplicatedUser(t *testing.T) {
	c := client.NewApiClient(config.Config.APIBaseURL)
	mongo.Cleanup()

	username := "test_user"
	other_username := "miguel"
	email := "valid@email.com"
	other_email := "miguel@email.com"
	password := "th!s_1s_@_s3cure_pASSw0rd"

	// Register the user
	err := c.Register(username, email, password)
	require.NoError(t, err)

	// Register with the same username
	err = c.Register(username, other_email, password)
	require.Error(t, err)

	// Register with the same email
	err = c.Register(other_username, email, password)
	require.Error(t, err)
}
