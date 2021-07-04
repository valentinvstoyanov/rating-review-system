package http

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"log"
	"net/http"
)

type ReviewHandler struct {
	reviewService rrs.ReviewService
}

func NewReviewHandler(reviewService rrs.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService}
}

func (rh *ReviewHandler) Create(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	type CreateReq struct {
		Rating    float32 `validate:"gte=1,lte=5" json:"rating"`
		Content   string  `validate:"required,min=1,max=1000" json:"content"`
		EntityId  uint    `validate:"gte=0" json:"entityId"`
		CreatorId uint    `validate:"gte=0" json:"creatorId"`
	}

	var createReq CreateReq

	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		log.Printf("Failed to decode review from request body, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(createReq); err != nil {
		log.Printf("Invalid create review request, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return
	}

	review, err := rh.reviewService.Create(&rrs.Review{
		Rating:    createReq.Rating,
		Content:   createReq.Content,
		EntityId:  createReq.EntityId,
		CreatorId: createReq.CreatorId,
	})

	if err != nil {
		log.Printf("Failed to create review, err=%s\n", err)
		handleError(w, err, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(review)
}

func (rh *ReviewHandler) GetById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	id, ok := extractParamId("id", w, req)
	if !ok {
		return
	}

	review, err := rh.reviewService.GetById(id)
	if err != nil {
		log.Printf("Review not found id=%d, err=%s\n", id, err)
		handleError(w, err, http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(review)
}

func (rh *ReviewHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-type", "application/json")

	reviews := rh.reviewService.GetAll()

	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(reviews)
}

func (rh *ReviewHandler) GetByEntityId(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	entityId, ok := extractParamId("id", w, req)
	if !ok {
		return
	}

	reviews := rh.reviewService.GetByEntityId(entityId)

	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(reviews)
}

func (rh *ReviewHandler) GetByCreatorId(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "application/json")

	creatorId, ok := extractParamId("id", w, req)
	if !ok {
		return
	}

	reviews := rh.reviewService.GetByCreatorId(creatorId)

	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(reviews)
}
