package http

import (
	"encoding/json"
	"github.com/ardiantirta/todo-crud/models"
	"github.com/ardiantirta/todo-crud/todo/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"strconv"
)

type TodoHandler struct {
	TodoService service.Service
}

func NewTodoHandler(r *mux.Router, ts service.Service) {
	handler := &TodoHandler{
		TodoService: ts,
	}

	r.HandleFunc("/todos", handler.FetchTodo).Methods("GET")
	r.HandleFunc("/todo/{id}", handler.GetTodoById).Methods("GET")
	r.HandleFunc("/todo", handler.CreateTodo).Methods("POST")
	r.HandleFunc("/todo/{id}", handler.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todo/{id}", handler.DeleteTodo).Methods("DELETE")
}

func (s TodoHandler) FetchTodo(w http.ResponseWriter, r *http.Request) {
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

func (s TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {
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

func (s TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
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

func (s TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
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

func (s TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
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
