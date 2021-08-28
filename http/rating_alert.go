package http

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"log"
	"net/http"
)

type RatingAlertHandler struct {
	ratingAlertService rrs.RatingAlertService
}

func NewRatingAlertHandler(ratingAlertService rrs.RatingAlertService) *RatingAlertHandler {
	return &RatingAlertHandler{ratingAlertService}
}

func (rah *RatingAlertHandler) Create(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	entityId, ok := extractUintPathVar("id", w, req)
	if !ok {
		return
	}

	type CreateReq struct {
		PercentageChange float32 `validate:"gt=0,lte=100" json:"percentageChange"`
		PeriodMinutes    uint    `validate:"gt=0,lte=1000000" json:"periodMinutes"`
	}

	var createReq CreateReq

	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		log.Printf("Failed to decode create rating alert from request body, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(createReq); err != nil {
		log.Printf("Invalid create rating alert request, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	ratingAlert, err := rah.ratingAlertService.Create(&rrs.RatingAlert{
		EntityId:         entityId,
		PercentageChange: createReq.PercentageChange,
		PeriodMinutes:    createReq.PeriodMinutes,
	})

	if err != nil {
		log.Printf("Failed to create rating alert, err=%s\n", err)
		handleError(w, err, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ratingAlert)
}

func (rah *RatingAlertHandler) UpdateById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	id, ok := extractUintPathVar("id", w, req)
	if !ok {
		return
	}

	type UpdateReq struct {
		PercentageChange float32 `validate:"gt=0,lte=100" json:"percentageChange"`
		PeriodMinutes    uint    `validate:"gt=1,lte=1000000" json:"periodMinutes"`
	}

	var updateReq UpdateReq

	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		log.Printf("Failed to decode update rating alert from request body, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(updateReq); err != nil {
		log.Printf("Invalid update rating alert request, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	ratingAlert, err := rah.ratingAlertService.UpdateById(id, updateReq.PercentageChange, updateReq.PeriodMinutes)

	if err != nil {
		log.Printf("Failed to update rating alert, err=%s\n", err)
		handleError(w, err, http.StatusConflict)
		return
	}

	_ = json.NewEncoder(w).Encode(ratingAlert)
}

func (rah *RatingAlertHandler) GetById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	id, ok := extractUintPathVar("id", w, req)
	if !ok {
		return
	}

	ratingAlert, err := rah.ratingAlertService.GetById(id)
	if err != nil {
		log.Printf("Rating alert not found id=%d, err=%s\n", id, err)
		handleError(w, err, http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(ratingAlert)
}

func (rah *RatingAlertHandler) GetByEntityId(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	entityId, ok := extractUintRequestParam("entityId", w, req)
	if !ok {
		return
	}

	ratingAlert, err := rah.ratingAlertService.GetByEntityId(entityId)
	if err != nil {
		log.Printf("Rating alert not found entityId=%d, err=%s\n", entityId, err)
		handleError(w, err, http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(ratingAlert)
}

func (rah *RatingAlertHandler) DeleteById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	id, ok := extractUintPathVar("id", w, req)
	if !ok {
		return
	}

	ratingAlert, err := rah.ratingAlertService.DeleteById(id)
	if err != nil {
		log.Printf("Rating alert not found id=%d, err=%s\n", id, err)
		handleError(w, err, http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(ratingAlert)
}
