package response

import "encoding/json"

type ResponseWithError struct {
	Data  interface{}
	Error string
}

func NewErrorResponse(err error) []byte {
	if res, err := json.Marshal(ResponseWithError{
		Error: err.Error(),
	}); err != nil {
		return res
	}
	return nil
}

func NewDataResponse(data interface{}) []byte {
	if res, err := json.Marshal(ResponseWithError{
		Data: data,
	}); err != nil {
		return res
	}
	return nil
}
