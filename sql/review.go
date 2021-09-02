package sql

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"github.com/valentinvstoyanov/rating-review-system/email"
	"github.com/valentinvstoyanov/rating-review-system/env"
	"log"
	"time"
)

type PersistentReviewService struct {
	db            *gorm.DB
	userService   rrs.UserService
	entityService rrs.EntityService
}

func NewReviewService(db *gorm.DB, userService rrs.UserService, entityService rrs.EntityService) *PersistentReviewService {
	return &PersistentReviewService{db, userService, entityService}
}

func (rs *PersistentReviewService) Create(review *rrs.Review) (*rrs.Review, error) {
	reviewer, err := rs.userService.GetById(review.CreatorId)
	if err != nil {
		return nil, errors.New("failed to find creator with such id: " + err.Error())
	}

	entity, err := rs.entityService.GetById(review.EntityId)
	if err != nil {
		return nil, errors.New("failed to find entity with such id: " + err.Error())
	}

	if entity.CreatorId == review.CreatorId {
		return nil, errors.New("cannot rate your own entity")
	}

	entityCreator, err := rs.userService.GetById(entity.CreatorId)
	if err != nil {
		return nil, errors.New("failed to find creator with such id: " + err.Error())
	}

	if err := rs.db.Create(&review).Error; err != nil {
		return nil, err
	}

	newReviewsCount := entity.ReviewsCount + 1
	newAvgRating := (review.Rating + (entity.AvgRating * float32(entity.ReviewsCount))) / float32(newReviewsCount)
	_, err = rs.entityService.UpdateRating(entity.Id, newAvgRating, newReviewsCount)
	if err != nil {
		return nil, errors.New("failed to rate the entity: " + err.Error())
	}

	message := email.Message{
		Sender:   env.GetEnvVar("EMAIL_SENDER"),
		Receiver: entityCreator.Email,
		Subject:  "New entity review",
		Content:  fmt.Sprintf("%s just got reviewed by %s with %.2f stars.\n\n\nEnjoy!\n", entity.Name, reviewer.FirstName+" "+reviewer.LastName, review.Rating),
	}

	if err := email.Send(message); err != nil {
		log.Printf("Failed to notify the creator %d of entity %d for the new review %d\n", entity.CreatorId, entity.Id, review.Id)
	} else {
		log.Printf("New review from %d for entity with id=%d, name=%s and creatorId=%d", review.CreatorId, entity.Id, entity.Name, entity.CreatorId)
	}

	return review, nil
}

func (rs *PersistentReviewService) GetById(id uint) (*rrs.Review, error) {
	var review rrs.Review
	if err := rs.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (rs *PersistentReviewService) GetAll() []rrs.Review {
	var reviews []rrs.Review
	rs.db.Find(&reviews)
	return reviews
}

func (rs *PersistentReviewService) GetByEntityId(entityId uint) []rrs.Review {
	var reviews []rrs.Review
	rs.db.Where(&rrs.Review{EntityId: entityId}).Find(&reviews)
	return reviews
}

func (rs *PersistentReviewService) GetByCreatorId(creatorId uint) []rrs.Review {
	var reviews []rrs.Review
	rs.db.Where(&rrs.Review{CreatorId: creatorId}).Find(&reviews)
	return reviews
}

func (rs *PersistentReviewService) GetByEntityIdInPeriod(entityId uint, startTime time.Time, endTime time.Time) []rrs.Review {
	var reviews []rrs.Review
	rs.db.Where("entity_id = ? AND created_at >= ? AND created_at <= ?", entityId, startTime, endTime).Order("id").Find(&reviews)
	return reviews
}
