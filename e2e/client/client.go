package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ApiClient struct {
	baseURL string
	client  http.Client
	log     zerolog.Logger
}

func NewApiClient(baseURL string) ApiClient {
	return ApiClient{
		baseURL: baseURL,
		client:  http.Client{},
		log:     zerolog.New(zerolog.ConsoleWriter{}),
	}
}

func (ac *ApiClient) get(endpoint string) (*http.Response, error) {
	log.Info().
		Str("endpoint", endpoint).
		Str("method", http.MethodGet).
		Msg("request sent")
	return ac.client.Get(fmt.Sprintf("%s/%s", ac.baseURL, endpoint))
}

func (ac *ApiClient) post(endpoint string, body any) (*http.Response, error) {
	log.Info().
		Str("endpoint", endpoint).
		Interface("body", body).
		Str("method", http.MethodPost).
		Msg("request sent")
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return ac.client.Post(fmt.Sprintf("%s/%s", ac.baseURL, endpoint), "application/json", bytes.NewBuffer(jsonBody))
}

func (ac *ApiClient) put(endpoint string, body any) (*http.Response, error) {
	log.Info().
		Str("endpoint", endpoint).
		Interface("body", body).
		Str("method", http.MethodPut).
		Msg("request sent")
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", ac.baseURL, endpoint), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	return ac.client.Do(req)
}

func (ac *ApiClient) delete(endpoint string) (*http.Response, error) {
	log.Info().
		Str("endpoint", endpoint).
		Str("method", http.MethodDelete).
		Msg("request sent")
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", ac.baseURL, endpoint), nil)
	if err != nil {
		return nil, err
	}
	return ac.client.Do(req)
}
