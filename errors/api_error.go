package errors

import (
	"encoding/json"
)

type APIError struct {
	Status  int         `json:"-"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (this *APIError) Error() string {
	bs, err := json.Marshal(this)
	if err != nil {
		return this.Message
	}
	return string(bs)
}
