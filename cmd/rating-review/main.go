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
	entityService := sql.NewEntityService(sql.GetDB(), userService)
	ratingAlertService := sql.NewRatingAlertService(sql.GetDB(), entityService)
	reviewService := sql.NewReviewService(sql.GetDB(), userService, entityService)
	ratingAlertTriggerService := sql.NewRatingAlertTriggerService(entityService, reviewService, ratingAlertService)

	userHandler := handlers.NewUserHandler(userService)
	entityHandler := handlers.NewEntityHandler(entityService)
	ratingAlertHandler := handlers.NewRatingAlertHandler(ratingAlertService)
	reviewHandler := handlers.NewReviewHandler(reviewService, ratingAlertService, ratingAlertTriggerService)

	registerApiRoutes(r.PathPrefix("/api").Subrouter(), userHandler, entityHandler, reviewHandler, ratingAlertHandler)

	http.Handle("/", r)

	formattedPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Listening for requests at %s\n", formattedPort)
	log.Fatal(http.ListenAndServe(formattedPort, r))
}

func registerApiRoutes(r *mux.Router, userHandler *handlers.UserHandler, entityHandler *handlers.EntityHandler, reviewHandler *handlers.ReviewHandler, ratingAlertHandler *handlers.RatingAlertHandler) {
	r.HandleFunc("/users/{id:[0-9]+}/reviews", reviewHandler.GetByCreatorId).Methods(http.MethodGet)
	r.HandleFunc("/users/{id:[0-9]+}", userHandler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/users", userHandler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/entities/{id:[0-9]+}/rating-alerts", ratingAlertHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/entities/{id:[0-9]+}/reviews", reviewHandler.GetByEntityId).Methods(http.MethodGet)
	r.HandleFunc("/entities/{id:[0-9]+}", entityHandler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/entities", entityHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/entities", entityHandler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/reviews/{id:[0-9]+}", reviewHandler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/reviews", reviewHandler.Create).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.HandleFunc("/reviews", reviewHandler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/rating-alerts", ratingAlertHandler.GetByEntityId).Queries("entityId", "{entityId:[0-9]+}").Methods(http.MethodGet)
	r.HandleFunc("/rating-alerts/{id:[0-9]+}", ratingAlertHandler.UpdateById).Methods(http.MethodPut).Headers("Content-Type", "application/json")
	r.HandleFunc("/rating-alerts/{id:[0-9]+}", ratingAlertHandler.DeleteById).Methods(http.MethodDelete)
	r.HandleFunc("/rating-alerts/{id:[0-9]+}", ratingAlertHandler.GetById).Methods(http.MethodGet)
}
