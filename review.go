package rrs

import "time"

type Review struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	CreatorId uint      `json:"creatorId"`
	EntityId  uint      `json:"entityId"`
	Rating    float32   `json:"rating"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
