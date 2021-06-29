package sql

import (
	"github.com/jinzhu/gorm"
	rrs "github.com/valentinvstoyanov/rating-review-system"
)

type PersistentEntityService struct {
	db *gorm.DB
}

func NewEntityService(db *gorm.DB) *PersistentEntityService {
	return &PersistentEntityService{db}
}

func (us *PersistentEntityService) Create(entity *rrs.Entity) (*rrs.Entity, error) {
	if err := us.db.Create(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (us *PersistentEntityService) GetById(id uint) (*rrs.Entity, error) {
	var entity rrs.Entity
	if err := us.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (us *PersistentEntityService) GetAll() []rrs.Entity {
	var entities []rrs.Entity
	us.db.Find(&entities)
	return entities
}
