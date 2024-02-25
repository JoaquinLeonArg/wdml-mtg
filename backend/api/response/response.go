package response

import (
	"encoding/json"
)

type ResponseWithError struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func NewErrorResponse(err error) []byte {
	res, err := json.Marshal(ResponseWithError{
		Error: err.Error(),
	})
	if err != nil {
		return nil
	}
	return res
}

func NewDataResponse(data interface{}) []byte {
	res, err := json.Marshal(ResponseWithError{
		Data: data,
	})
	if err != nil {
		return nil
	}
	return res
}
