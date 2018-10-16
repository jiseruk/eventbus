package app

type APIError struct {
	Message      string   `json:"message"`
	Code        string      `json:"code"`
	Status		int			`json:"status"`
}

func NewAPIError(status int, code string, message string) *APIError {
	return &APIError{
		Status:		status,
		Message:    message,
		Code:       code,
	}
}

