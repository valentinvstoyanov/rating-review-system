package rrs

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `gorm:"unique" json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserService interface {
	Create(user *User) (*User, error)
	GetById(id uint) (*User, error)
	GetAll() []User
}
