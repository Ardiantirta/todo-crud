package http

import (
	"encoding/json"
	"fmt"
	"github.com/ardiantirta/todo-crud/models"
	"github.com/ardiantirta/todo-crud/todo/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"os"
	"strconv"
)

type TodoHandler struct {
	TodoService service.Service
}

func NewTodoHandler(r *mux.Router, ts service.Service) {
	handler := &TodoHandler{
		TodoService: ts,
	}

	r.Handle("/todos", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.FetchTodo))).Methods(http.MethodGet)
	r.Handle("/todoschannel", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.FetchWChannel))).Methods(http.MethodGet)
	r.Handle("/todo/{id}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.GetTodoById))).Methods(http.MethodGet)
	r.Handle("/todo", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.CreateTodo))).Methods(http.MethodPost)
	r.Handle("/todobulk", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.CreateBulkTodo))).Methods(http.MethodPost)
	r.Handle("/todo/{id}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.UpdateTodo))).Methods(http.MethodPut)
	r.Handle("/todo/{id}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.DeleteTodo))).Methods(http.MethodDelete)
}

func (s *TodoHandler) FetchTodo(w http.ResponseWriter, r *http.Request) {
	qsPage := r.URL.Query()["page"][0]
	page, err := strconv.Atoi(qsPage)
	if err != nil {
		logrus.Error(err)
	}

	qsLimit := r.URL.Query()["limit"][0]
	limit, err := strconv.Atoi(qsLimit)
	if err != nil {
		logrus.Error(err)
	}

	listTodo, totalCount, err := s.TodoService.Fetch(page, limit)

	data := make(map[string]interface{})
	data["data"] = listTodo
	data["count"] = len(listTodo)
	data["limit"] = limit
	data["page"] = page
	data["min_page"] = 1
	data["max_page"] = math.Ceil(float64(totalCount) / float64(limit))

	if err != nil {
		logrus.Error(err)

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.ResponseHttp{
		Code:    strconv.Itoa(http.StatusOK),
		Data:    data,
		Message: "",
	})
}

func (s *TodoHandler) FetchWChannel(w http.ResponseWriter, r *http.Request) {
	qsPage := r.URL.Query()["page"][0]
	qsLimit := r.URL.Query()["limit"][0]
	qsChannel := r.URL.Query()["channel"][0]

	page, err := strconv.Atoi(qsPage)
	if err != nil {
		logrus.Error(err)
	}

	limit, err := strconv.Atoi(qsLimit)
	if err != nil {
		logrus.Error(err)
	}

	channel, err := strconv.Atoi(qsChannel)
	if err != nil {
		logrus.Error(err)
	}

	listTodo, totalCount, err := s.TodoService.FetchWChannel(page, limit, channel)

	data := make(map[string]interface{})
	data["data"] = listTodo
	data["count"] = totalCount
	data["limit"] = limit
	data["page"] = page
	data["min_page"] = 1
	data["max_page"] = math.Ceil(float64(totalCount) / float64(limit))

	if err != nil {
		logrus.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusOK, models.ResponseHttp{
		Data:    data,
		Code:    strconv.Itoa(http.StatusOK),
		Message: "",
	})

}

func (s *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.Error(err)
	}

	res, err := s.TodoService.GetById(int64(id))
	if err != nil {
		logrus.Error(err)

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["data"] = res

	respondWithJSON(w, http.StatusOK, models.ResponseHttp{
		Code:    strconv.Itoa(http.StatusOK),
		Data:    data,
		Message: "",
	})
}

func (s *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		logrus.Error(err)
	}
	err = s.TodoService.Create(todo)
	if err != nil {
		logrus.Error(err)

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.ResponseHttp{
		Data:    nil,
		Code:    strconv.Itoa(http.StatusOK),
		Message: "Insert Todo: Success",
	})
}

func (s *TodoHandler) CreateBulkTodo(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logrus.Error(err)
	}

	bTodo, err := json.Marshal(body["todo"])
	if err != nil {
		logrus.Error(err)
	}

	var todo models.Todo
	err = json.Unmarshal(bTodo, &todo)
	if err != nil {
		logrus.Error(err)
	}

	countBulk := int(body["countBulk"].(float64))

	err = s.TodoService.CreateBulk(todo, countBulk)
	if err != nil {
		logrus.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	data := make(map[string]interface{})
	data["todo"] = todo
	data["countBulk"] = countBulk

	message := fmt.Sprintf("insert %d data success", data["countBulk"])

	respondWithJSON(w, http.StatusOK, models.ResponseHttp{
		Data:    nil,
		Code:    strconv.Itoa(http.StatusOK),
		Message: message,
	})
}

func (s *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.Error(err)
	}

	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		logrus.Error(err)
	}

	err = s.TodoService.Update(int64(id), todo)
	if err != nil {
		logrus.Error(err)

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.ResponseHttp{
		Data:    nil,
		Code:    strconv.Itoa(http.StatusOK),
		Message: "Update todo: success",
	})
}

func (s *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.Error(err)
	}

	err = s.TodoService.Delete(int64(id))
	if err != nil {
		logrus.Error(err)

		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.ResponseHttp{
		Data:    nil,
		Code:    strconv.Itoa(http.StatusOK),
		Message: "Delete Todo: Success",
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
