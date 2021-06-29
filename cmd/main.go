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

	entityService := sql.NewEntityService(sql.GetDB(), userService)
	entityHandler := handlers.NewEntityHandler(entityService)

	reviewService := sql.NewReviewService(sql.GetDB(), userService, entityService)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	registerApiRoutes(r.PathPrefix("/api").Subrouter(), userHandler, entityHandler, reviewHandler)

	http.Handle("/", r)

	formattedPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Listening for requests at %s\n", formattedPort)
	log.Fatal(http.ListenAndServe(formattedPort, r))
}

func registerApiRoutes(r *mux.Router, userHandler *handlers.UserHandler, entityHandler *handlers.EntityHandler, reviewHandler *handlers.ReviewHandler) {
	r.HandleFunc("/users/{id:[0-9]+}", userHandler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/users", userHandler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/entities/{id:[0-9]+}", entityHandler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/entities", entityHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/entities", entityHandler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/reviews/{id:[0-9]+}", reviewHandler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/reviews", reviewHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/reviews", reviewHandler.GetAll).Methods(http.MethodGet)
}
