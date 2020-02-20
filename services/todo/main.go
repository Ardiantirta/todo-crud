package main

import (
	"fmt"
	"github.com/ardiantirta/todo-crud/common/http/request"
	"github.com/ardiantirta/todo-crud/common/http/response"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	todoHttp "github.com/ardiantirta/todo-crud/services/todo/delivery/http"
	_todoRepository "github.com/ardiantirta/todo-crud/services/todo/repository"
	_todoService "github.com/ardiantirta/todo-crud/services/todo/service"
)

func init() {
	viper.SetConfigFile("./config/config.json")
	if err := viper.ReadInConfig(); err != nil {
		panic("err")
	}

	if viper.GetBool("debug") {
		fmt.Println("service run on debug mode")
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
	dbConn, err := gorm.Open("postgres", dsn)
	if err != nil {
		logrus.Error(err)
	}

	if err := dbConn.DB().Ping(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	fmt.Println("ping from db")

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	dbConn.Debug().AutoMigrate(
		&_todoRepository.Todo{},
	)

	r := mux.NewRouter()

	defaultHandler := request.NewDefaultHandler(response.NewDefaultJSONResponder())
	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defaultHandler.Index(w, r)
		return
	}))).Methods(http.MethodGet)

	todoRepository := _todoRepository.NewTodoRepository(dbConn)
	todoService := _todoService.NewTodoService(todoRepository)
	todoHttp.NewTodoHandler(r, todoService)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE"})

	log.Fatal(http.ListenAndServe(viper.GetString("server.address"), handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
