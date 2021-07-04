package sql

import (
	"errors"
	"github.com/jinzhu/gorm"
	rrs "github.com/valentinvstoyanov/rating-review-system"
)

type PersistentRatingAlertService struct {
	db            *gorm.DB
	entityService rrs.EntityService
}

func NewRatingAlertService(db *gorm.DB, entityService rrs.EntityService) *PersistentRatingAlertService {
	return &PersistentRatingAlertService{db, entityService}
}

func (rs *PersistentRatingAlertService) Create(ratingAlert *rrs.RatingAlert) (*rrs.RatingAlert, error) {
	_, err := rs.entityService.GetById(ratingAlert.EntityId)
	if err != nil {
		return nil, errors.New("failed to find entity with such id: " + err.Error())
	}

	if err := rs.db.Create(&ratingAlert).Error; err != nil {
		return nil, err
	}

	return ratingAlert, nil
}

func (rs *PersistentRatingAlertService) UpdateById(id uint, percentageChange float32, periodMinutes uint) (*rrs.RatingAlert, error) {
	res := db.Model(&rrs.RatingAlert{}).Where("id = ?", id).UpdateColumns(&rrs.RatingAlert{PercentageChange: percentageChange, PeriodMinutes: periodMinutes})

	if err := res.Error; err != nil {
		return nil, err
	}

	if res.RowsAffected != 1 {
		return nil, errors.New("failed to update rating alert")
	}

	return rs.GetById(id)
}

func (rs *PersistentRatingAlertService) GetById(id uint) (*rrs.RatingAlert, error) {
	var ratingAlert rrs.RatingAlert
	if err := rs.db.First(&ratingAlert, id).Error; err != nil {
		return nil, err
	}
	return &ratingAlert, nil
}

func (rs *PersistentRatingAlertService) GetByEntityId(entityId uint) (*rrs.RatingAlert, error) {
	var ratingAlert rrs.RatingAlert
	if err := rs.db.Where(&rrs.RatingAlert{EntityId: entityId}).First(&ratingAlert).Error; err != nil {
		return nil, err
	}
	return &ratingAlert, nil
}

func (rs *PersistentRatingAlertService) DeleteById(id uint) (*rrs.RatingAlert, error) {
	ratingAlert, err := rs.GetById(id)
	if err != nil {
		return nil, errors.New("failed to find rating alert with such id: " + err.Error())
	}

	rs.db.Delete(&rrs.RatingAlert{}, id)

	return ratingAlert, nil
}
