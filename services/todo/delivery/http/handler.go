package http

import (
	"encoding/json"
	"github.com/ardiantirta/todo-crud/common/http/request"
	"github.com/ardiantirta/todo-crud/common/http/response"
	"github.com/ardiantirta/todo-crud/common/message"
	"github.com/ardiantirta/todo-crud/services/todo/service"
	"github.com/ardiantirta/todo-crud/services/todo/service/todo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

type TodoHandler struct {
	TodoService service.Service
	jsonResponder response.JSONResponder
}

func NewTodoHandler(r *mux.Router, todoService service.Service) {
	handler := &TodoHandler{
		TodoService: todoService,
		jsonResponder: response.NewDefaultJSONResponder(),
	}

	v1 := r.PathPrefix("/todo").Subrouter()
	v1.Handle("", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.Create))).Methods(http.MethodPost)
	v1.Handle("", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.GetAll))).Methods(http.MethodGet)
	v1.Handle("/search", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.GetByTitle))).Methods(http.MethodGet)
	v1.Handle("/{id}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.GetByID))).Methods(http.MethodGet)
	v1.Handle("/done/{id}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.MarkAsDone))).Methods(http.MethodPut)
	v1.Handle("/favorite/{id}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.MarkAsFavorite))).Methods(http.MethodPut)
	v1.Handle("/{id}",handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.UpdateData))).Methods(http.MethodPut)
	v1.Handle("/{id}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.DeleteByID))).Methods(http.MethodDelete)
}

func (c *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	formData := new(todo.CreateRequest)

	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "invalid json body"))
		return
	}

	resp, err := c.TodoService.Create(r.Context(), formData)
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	c.jsonResponder.Data(w, http.StatusOK, resp)
	return
}

func (c *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	isDone := ""
	isDoneSlice, ok := r.URL.Query()["is_done"]
	if ok && len(isDoneSlice[0]) > 0 {
		isDone = isDoneSlice[0]
	}

	isFavorite := ""
	isFavoriteSlice, ok := r.URL.Query()["is_favorite"]
	if ok && len(isFavoriteSlice[0]) > 0 {
		isFavorite = isFavoriteSlice[0]
	}

	params := map[string]interface{}{
		"is_done": isDone,
		"is_favorite": isFavorite,
	}

	resp, err := c.TodoService.GetAll(r.Context(), params)
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	c.jsonResponder.Data(w, http.StatusOK, resp)
	return
}

func (c *TodoHandler) GetByTitle(w http.ResponseWriter, r *http.Request) {

	title := ""
	titleSlice, ok := r.URL.Query()["title"]
	if ok && len(titleSlice[0]) > 0 {
		title = titleSlice[0]
	}

	params := map[string]interface{}{
		"title": title,
	}

	resp, err := c.TodoService.GetByTitle(r.Context(), params)
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	c.jsonResponder.Data(w, http.StatusOK, resp)
	return
}

func (c *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id is not a valid number"))
		return
	}

	if id <= 0 {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id should be a positive number"))
		return
	}

	resp, err := c.TodoService.GetByID(r.Context(), id)
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	c.jsonResponder.Data(w, http.StatusOK, resp)
	return
}

func (c *TodoHandler) UpdateData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id is not a valid number"))
		return
	}

	if id <= 0 {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id should be a positive number"))
		return
	}

	formData := new(todo.UpdateRequest)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "invalid json body"))
		return
	}

	resp, err := c.TodoService.UpdateData(r.Context(), id, formData)
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	c.jsonResponder.Data(w, http.StatusOK, resp)
	return
}

func (c *TodoHandler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id is not a valid number"))
		return
	}

	if id <= 0 {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id should be a positive number"))
		return
	}

	formData := new(todo.DoneRequest)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "invalid json body"))
		return
	}

	if err := c.TodoService.MarkAsDone(r.Context(), id, formData); err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	defaultHandler := request.NewDefaultHandler(c.jsonResponder)
	defaultHandler.Index(w, r)
	return
}

func (c *TodoHandler) MarkAsFavorite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id is not a valid number"))
		return
	}

	if id <= 0 {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id shoud be a positive number"))
		return
	}

	formData := new(todo.FavoriteRequest)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "invalid json body"))
		return
	}

	if err := c.TodoService.MarkAsFavorite(r.Context(), id, formData); err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	defaultHandler := request.NewDefaultHandler(c.jsonResponder)
	defaultHandler.Index(w, r)
	return
}

func (c *TodoHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id is not a valid number"))
		return
	}

	if id <= 0 {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, "id should be a positive number"))
		return
	}

	if err := c.TodoService.DeleteByID(r.Context(), id); err != nil {
		c.jsonResponder.Error(w, http.StatusBadRequest, message.NewErrorMessage(0, err.Error()))
		return
	}

	defaultHandler := request.NewDefaultHandler(c.jsonResponder)
	defaultHandler.Index(w, r)
	return
}
