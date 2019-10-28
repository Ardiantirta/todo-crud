package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ardiantirta/todo-crud/models"
	"github.com/ardiantirta/todo-crud/todo/repository"
	"github.com/ardiantirta/todo-crud/todo/service"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"net/url"
	"os"

	httpTransport "github.com/ardiantirta/todo-crud/todo/delivery/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode.")
	}
}

func main() {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connStr, val.Encode())
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.Error(err)
	}

	err = dbConn.Ping()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := mux.NewRouter()

	todoRepository := repository.NewTodoRepository(dbConn)

	todoService := service.NewTodoService(todoRepository)

	httpTransport.NewTodoHandler(r, todoService)

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(models.ResponseHttp{
			Code: "200",
			Data: "pong",
		})
	})

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH"})

	log.Fatal(http.ListenAndServe(viper.GetString("server.address"), handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
