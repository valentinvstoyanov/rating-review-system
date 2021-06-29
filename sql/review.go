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
		return nil, errors.New("Failed to find creator with such id: " + err.Error())
	}

	if _, err := rs.entityService.GetById(review.EntityId); err != nil {
		return nil, errors.New("Failed to find entity with such id: " + err.Error())
	}

	if err := rs.db.Create(&review).Error; err != nil {
		return nil, err
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
