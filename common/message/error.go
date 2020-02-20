package message

import (
	"encoding/json"
)

type ErrorResponse struct {
	ErrCode int `json:"err_code"`
	ErrMessage string `json:"err_message"`
}

func NewErrorMessage(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrCode: code,
		ErrMessage: message,
	}
}

func (c *ErrorResponse) Code() int {
	return c.ErrCode
}

func (c *ErrorResponse) Message() string {
	return c.ErrMessage
}

func (c *ErrorResponse) Error() string {
	b, _ := json.Marshal(c)

	return string(b)
}


