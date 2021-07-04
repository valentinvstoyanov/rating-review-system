package sql

import (
	"errors"
	"github.com/jinzhu/gorm"
	rrs "github.com/valentinvstoyanov/rating-review-system"
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
	if _, err := rs.userService.GetById(review.CreatorId); err != nil {
		return nil, errors.New("failed to find creator with such id: " + err.Error())
	}

	entity, err := rs.entityService.GetById(review.EntityId)
	if err != nil {
		return nil, errors.New("failed to find entity with such id: " + err.Error())
	}

	if entity.CreatorId == review.CreatorId {
		return nil, errors.New("cannot rate your own entity")
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

	//TODO: Notify entity.CreatorId that there is new review

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
