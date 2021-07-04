package rrs

import "time"

type RatingAlert struct {
	Id               uint      `gorm:"primaryKey" json:"id"`
	EntityId         uint      `gorm:"unique" json:"entityId"`
	PercentageChange float32   `json:"percentageChange"`
	PeriodMinutes    uint      `json:"periodMinutes"`
	LastTriggeredAt  time.Time `json:"lastTriggeredAt"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type RatingAlertService interface {
	Create(ratingAlert *RatingAlert) (*RatingAlert, error)
	UpdateById(id uint, percentageChange float32, periodMinutes uint) (*RatingAlert, error)
	GetById(id uint) (*RatingAlert, error)
	GetByEntityId(entityId uint) (*RatingAlert, error)
	DeleteById(id uint) (*RatingAlert, error)
	UpdateLastTriggeredAtById(id uint, lastTriggeredAt time.Time) error
}

type RatingAlertTriggerService interface {
	Trigger(ratingAlert *RatingAlert) (*RatingAlert, error)
}
