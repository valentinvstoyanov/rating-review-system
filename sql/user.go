package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/valentinvstoyanov/rating-review-system"
)

type PersistentUserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *PersistentUserService {
	return &PersistentUserService{db}
}

func (us *PersistentUserService) Create(user *rrs.User) (*rrs.User, error) {
	if err := us.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (us *PersistentUserService) GetById(id uint) (*rrs.User, error) {
	var user rrs.User
	if err := us.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *PersistentUserService) GetAll() []rrs.User {
	var users []rrs.User
	us.db.Find(&users)
	return users
}
