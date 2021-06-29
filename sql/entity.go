package sql

import (
	"errors"
	"github.com/jinzhu/gorm"
	rrs "github.com/valentinvstoyanov/rating-review-system"
)

type PersistentEntityService struct {
	db          *gorm.DB
	userService rrs.UserService
}

func NewEntityService(db *gorm.DB, userService rrs.UserService) *PersistentEntityService {
	return &PersistentEntityService{db, userService}
}

func (es *PersistentEntityService) Create(entity *rrs.Entity) (*rrs.Entity, error) {
	if _, err := es.userService.GetById(entity.CreatorId); err != nil {
		return nil, errors.New("Failed to find creator with such id: " + err.Error())
	}

	if err := es.db.Create(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (es *PersistentEntityService) GetById(id uint) (*rrs.Entity, error) {
	var entity rrs.Entity
	if err := es.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (es *PersistentEntityService) GetAll() []rrs.Entity {
	var entities []rrs.Entity
	es.db.Find(&entities)
	return entities
}
