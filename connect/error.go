package connect

import (
	"encoding/json"
	"fmt"
)

// APIError is returned for non-success HTTP responses from connect-state endpoints.
type APIError struct {
	Status  int
	Message string
	Body    string
}

func (e APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("connect: api error (%d): %s", e.Status, e.Message)
	}
	return fmt.Sprintf("connect: api error (%d)", e.Status)
}

func newAPIError(status int, body []byte) error {
	err := APIError{Status: status, Body: string(body)}
	if len(body) == 0 {
		return err
	}
	payload := struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
		Message string `json:"message"`
	}{}
	_ = json.Unmarshal(body, &payload)
	if err.Message = payload.Error.Message; err.Message == "" {
		err.Message = payload.Message
	}
	return err
}
