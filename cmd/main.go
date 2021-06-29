package main

import (
	"fmt"
	"github.com/gorilla/mux"
	handlers "github.com/valentinvstoyanov/rating-review-system/http"
	"github.com/valentinvstoyanov/rating-review-system/sql"
	"log"
	"net/http"
	"os"
)

func main() {
	sql.CreateDatabaseConnection()
	sql.AutoMigrateAll()

	r := mux.NewRouter()

	userService := sql.NewUserService(sql.GetDB())
	userHandler := handlers.NewUserHandler(userService)
	registerApiRoutes(r.PathPrefix("/api/").Subrouter(), userHandler)

	http.Handle("/", r)

	formattedPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Listening for requests at %s\n", formattedPort)
	log.Fatal(http.ListenAndServe(formattedPort, r))
}

func registerApiRoutes(r *mux.Router, userHandler *handlers.UserHandler) {
	r.HandleFunc("/users/{id:[0-9]+}", userHandler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/users", userHandler.GetAll).Methods(http.MethodGet)
}
