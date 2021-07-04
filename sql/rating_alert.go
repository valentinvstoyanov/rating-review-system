package sql

import (
	"errors"
	"github.com/jinzhu/gorm"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"time"
)

type PersistentRatingAlertService struct {
	db            *gorm.DB
	entityService rrs.EntityService
}

func NewRatingAlertService(db *gorm.DB, entityService rrs.EntityService) *PersistentRatingAlertService {
	return &PersistentRatingAlertService{db, entityService}
}

func (ras *PersistentRatingAlertService) Create(ratingAlert *rrs.RatingAlert) (*rrs.RatingAlert, error) {
	_, err := ras.entityService.GetById(ratingAlert.EntityId)
	if err != nil {
		return nil, errors.New("failed to find entity with such id: " + err.Error())
	}

	if err := ras.db.Create(&ratingAlert).Error; err != nil {
		return nil, err
	}

	return ratingAlert, nil
}

func (ras *PersistentRatingAlertService) UpdateById(id uint, percentageChange float32, periodMinutes uint) (*rrs.RatingAlert, error) {
	res := db.Model(&rrs.RatingAlert{}).Where("id = ?", id).UpdateColumns(&rrs.RatingAlert{PercentageChange: percentageChange, PeriodMinutes: periodMinutes})

	if err := res.Error; err != nil {
		return nil, err
	}

	if res.RowsAffected != 1 {
		return nil, errors.New("failed to update rating alert")
	}

	return ras.GetById(id)
}

func (ras *PersistentRatingAlertService) GetById(id uint) (*rrs.RatingAlert, error) {
	var ratingAlert rrs.RatingAlert
	if err := ras.db.First(&ratingAlert, id).Error; err != nil {
		return nil, err
	}
	return &ratingAlert, nil
}

func (ras *PersistentRatingAlertService) GetByEntityId(entityId uint) (*rrs.RatingAlert, error) {
	var ratingAlert rrs.RatingAlert
	if err := ras.db.Where(&rrs.RatingAlert{EntityId: entityId}).First(&ratingAlert).Error; err != nil {
		return nil, err
	}
	return &ratingAlert, nil
}

func (ras *PersistentRatingAlertService) DeleteById(id uint) (*rrs.RatingAlert, error) {
	ratingAlert, err := ras.GetById(id)
	if err != nil {
		return nil, errors.New("failed to find rating alert with such id: " + err.Error())
	}

	ras.db.Delete(&rrs.RatingAlert{}, id)

	return ratingAlert, nil
}

func (ras *PersistentRatingAlertService) UpdateLastTriggeredAtById(id uint, lastTriggeredAt time.Time) error {
	res := db.Model(&rrs.RatingAlert{}).Where("id = ?", id).UpdateColumns(&rrs.RatingAlert{LastTriggeredAt: lastTriggeredAt})

	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		return errors.New("failed to update rating alert")
	}

	return nil
}
