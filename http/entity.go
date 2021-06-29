package http

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"log"
	"net/http"
)

type EntityHandler struct {
	entityService rrs.EntityService
}

func NewEntityHandler(entityService rrs.EntityService) *EntityHandler {
	return &EntityHandler{entityService}
}

func (eh *EntityHandler) Create(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	type CreateReq struct {
		Name      string `validate:"required,min=1,max=100" json:"name"`
		CreatorId uint   `validate:"gte=0" json:"creatorId"`
	}

	var createReq CreateReq

	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		log.Printf("Failed to decode entity from request body, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(createReq); err != nil {
		log.Printf("Invalid create entity request, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	entity, err := eh.entityService.Create(&rrs.Entity{
		Name: createReq.Name,
		CreatorId:  createReq.CreatorId,
	})

	if err != nil {
		log.Printf("Failed to create entity, err=%s\n", err)
		handleError(w, err, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(entity)
}

func (eh *EntityHandler) GetById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	id, ok := extractParamId("id", w, req)
	if !ok {
		return
	}

	entity, err := eh.entityService.GetById(id)
	if err != nil {
		log.Printf("Entity not found id=%d, err=%s\n", id, err)
		handleError(w, err, http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(entity)
}

func (eh *EntityHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-type", "application/json")

	entities := eh.entityService.GetAll()

	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(entities)
}
