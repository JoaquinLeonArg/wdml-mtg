package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ApiClient struct {
	baseURL string
	client  http.Client
	log     zerolog.Logger
}

func NewApiClient(baseURL string) *ApiClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic("failed to create cookie jar")
	}
	return &ApiClient{
		baseURL: baseURL,
		client: http.Client{
			Jar: jar,
		},
		log: zerolog.New(zerolog.ConsoleWriter{}),
	}
}

func (ac *ApiClient) get(endpoint string) (*http.Response, error) {
	return ac.doRequest(http.MethodGet, endpoint, nil)
}

func (ac *ApiClient) post(endpoint string, body any) (*http.Response, error) {
	return ac.doRequest(http.MethodPost, endpoint, body)
}

func (ac *ApiClient) put(endpoint string, body any) (*http.Response, error) {
	return ac.doRequest(http.MethodPut, endpoint, body)
}

func (ac *ApiClient) delete(endpoint string) (*http.Response, error) {
	return ac.doRequest(http.MethodDelete, endpoint, nil)
}

func (ac *ApiClient) doRequest(method, endpoint string, body any) (*http.Response, error) {
	log.Info().
		Str("endpoint", endpoint).
		Interface("body", body).
		Str("method", method).
		Msg("sending request")
	jsonBody := []byte{}
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			log.Error().Err(err).Interface("body", body).Msg("failed to marshal request body")
			return nil, err
		}
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", ac.baseURL, endpoint), bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return nil, err
	}
	res, err := ac.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to send request")
		return nil, err
	}
	log.Info().Int("status_code", res.StatusCode).Msg("request completed")
	return res, err
}
