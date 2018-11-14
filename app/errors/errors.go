package errors

import (
	"fmt"
	"strings"
)

type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
}

func NewAPIError(status int, code string, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
		Code:    code,
	}
}

var (
	ErrorFieldRequired = "The field is required"
	ErrorFieldsInList  = "Must be one of [%s]"
)

func GetErrorMsg(msg string, field string) string {
	return fmt.Sprintf(msg, field)
}

func GetInListError(list ...string) string {
	return fmt.Sprintf(ErrorFieldsInList, strings.Join(list, ", "))
}
