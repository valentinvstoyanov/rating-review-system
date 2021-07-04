package http

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"log"
	"net/http"
)

type UserHandler struct {
	userService rrs.UserService
}

func NewUserHandler(userService rrs.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (uh *UserHandler) Create(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	type CreateReq struct {
		FirstName string `validate:"required,min=1,max=100" json:"firstName"`
		LastName  string `validate:"required,min=1,max=100" json:"lastName"`
		Email     string `validate:"email,required" json:"email"`
	}

	var createReq CreateReq

	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		log.Printf("Failed to decode user from request body, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(createReq); err != nil {
		log.Printf("Invalid create user request, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	user, err := uh.userService.Create(&rrs.User{
		FirstName: createReq.FirstName,
		LastName:  createReq.LastName,
		Email:     createReq.Email,
	})

	if err != nil {
		log.Printf("Failed to create user, err=%s\n", err)
		handleError(w, err, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) GetById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	id, ok := extractUintPathVar("id", w, req)
	if !ok {
		return
	}

	user, err := uh.userService.GetById(id)
	if err != nil {
		log.Printf("User not found id=%d, err=%s\n", id, err)
		handleError(w, err, http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-type", "application/json")

	users := uh.userService.GetAll()

	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(users)
}
