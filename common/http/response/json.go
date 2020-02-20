package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	JSONContentType = "application/json"
	JSONCharset = "utf-8"
)

type JSONResponder interface {
	Write(w http.ResponseWriter, status int, data interface{})
	Data(w http.ResponseWriter, status int, data interface{})
	Error(w http.ResponseWriter, status int, error ErrorContent)
}

type ErrorContent interface {
	Code() int
	Message() string
}

type ErrorResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
}

type jsonResponder struct {
	contentType string
}

func NewDefaultJSONResponder() JSONResponder {
	return &jsonResponder{
		contentType: fmt.Sprintf("%s; charset=%s", JSONContentType, JSONCharset),
	}
}

func (c *jsonResponder) Write(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", c.contentType)
	w.WriteHeader(status)
	if data == nil {
		return
	}

	content, _ := json.Marshal(data)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	_, _ = w.Write(content)
}

func (c *jsonResponder) Data(w http.ResponseWriter, status int, data interface{}) {
	content := DataResponse{Data: data}
	c.Write(w, status, content)
}

func (c *jsonResponder) Error(w http.ResponseWriter, status int, error ErrorContent) {
	content := ErrorResponse{
		Code: error.Code(),
		Message: error.Message(),
	}
	c.Write(w, status, content)
}
