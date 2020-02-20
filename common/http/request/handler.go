package request

import (
	"github.com/ardiantirta/todo-crud/common/http/response"
	"github.com/ardiantirta/todo-crud/common/message"
	"net/http"
)

type handler struct {
	jsonResponder response.JSONResponder
}

func NewDefaultHandler(jsonResponder response.JSONResponder) *handler {
	return &handler{
		jsonResponder: jsonResponder,
	}
}

func (c *handler) Index(w http.ResponseWriter, r *http.Request) {
	content := map[string]bool {
		"status": true,
	}

	c.jsonResponder.Write(w, http.StatusOK, content)
}

func (c *handler) NotFound(w http.ResponseWriter, r *http.Request) {
	c.jsonResponder.Error(w, http.StatusNotFound, message.NewErrorMessage(http.StatusNotFound, "url not found"))
}

func (c *handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	c.jsonResponder.Error(w, http.StatusMethodNotAllowed, message.NewErrorMessage(http.StatusMethodNotAllowed, "method not allowed"))
}


