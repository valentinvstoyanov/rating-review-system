package rrs

import "time"

type RatingAlert struct {
	Id               uint      `gorm:"primaryKey" json:"id"`
	EntityId         uint      `json:"entityId"`
	PercentageChange float32   `json:"percentageChange"`
	PeriodMinutes    uint      `json:"periodMinutes"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type RatingAlertService interface {
	Create(ratingAlert *RatingAlert) (*RatingAlert, error)
	UpdateById(id uint, percentageChange float32, periodMinutes uint) (*RatingAlert, error)
	GetById(id uint) (*RatingAlert, error)
	GetByEntityId(entityId uint) (*RatingAlert, error)
	DeleteById(id uint) (*RatingAlert, error)
}
