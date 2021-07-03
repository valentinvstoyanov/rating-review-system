package rrs

import "time"

type Entity struct {
	Id           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"unique" json:"name"`
	CreatorId    uint      `json:"creatorId"`
	AvgRating    float32   `json:"avgRating"`
	ReviewsCount uint      `json:"reviewsCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type EntityService interface {
	Create(entity *Entity) (*Entity, error)
	GetById(id uint) (*Entity, error)
	GetAll() []Entity
	UpdateRating(id uint, rating float32, reviewsCount uint) (float32, error)
}
