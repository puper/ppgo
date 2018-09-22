package errors

import (
	"fmt"
	"net/http"
	"strings"
)

var templates map[int]*errorTemplate

type errorTemplate struct {
	Status  int
	Message string
}

func AddTemplate(code int, status int, message string) {
	templates[code] = &errorTemplate{
		Status:  status,
		Message: message,
	}
}

func NewAPIError(code int, params map[string]interface{}) *APIError {
	err := &APIError{
		Status:  http.StatusOK,
		Code:    code,
		Message: "",
	}

	if template, ok := templates[code]; ok {
		err.Status = template.Status
		err.Message = replacePlaceholders(template.Message, params)
	}

	return err
}

func replacePlaceholders(message string, params map[string]interface{}) string {
	if len(message) == 0 || params == nil {
		return ""
	}
	for key, value := range params {
		message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
	}
	return message
}
