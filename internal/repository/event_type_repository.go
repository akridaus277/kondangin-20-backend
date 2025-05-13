package repository

import (
	dao "kondangin-backend/internal/dao"
	model "kondangin-backend/internal/model"
	"log"

	"gorm.io/gorm"
)

type EventTypeRepository interface {
	FindByCodeAndActive(code string, active bool) (*model.EventType, error)
	Create(eventType *model.EventType) (*model.EventType, error)
	FindByID(id uint) (*model.EventType, error)
	FindAll() ([]model.EventType, error)
	Update(eventType *model.EventType) (*model.EventType, error)
	Delete(id uint) error
	GetLovEventType() ([]dao.LovEventTypeDAO, error)
}

type eventTypeRepository struct {
	db *gorm.DB
}

func NewEventTypeRepository(db *gorm.DB) EventTypeRepository {
	return &eventTypeRepository{db}
}

func (r *eventTypeRepository) FindByCodeAndActive(code string, active bool) (*model.EventType, error) {
	var perm model.EventType
	err := r.db.Where("code = ? AND active = ?", code, active).First(&perm).Error
	return &perm, err
}

func (r *eventTypeRepository) Create(eventType *model.EventType) (*model.EventType, error) {
	if err := r.db.Create(eventType).Error; err != nil {
		return nil, err
	}
	return eventType, nil
}

func (r *eventTypeRepository) FindByID(id uint) (*model.EventType, error) {
	var eventType model.EventType
	err := r.db.First(&eventType, id).Error
	if err != nil {
		return nil, err
	}
	return &eventType, nil
}

func (r *eventTypeRepository) FindAll() ([]model.EventType, error) {
	var eventTypes []model.EventType
	err := r.db.Find(&eventTypes).Error
	return eventTypes, err
}

func (r *eventTypeRepository) Update(eventType *model.EventType) (*model.EventType, error) {
	if err := r.db.Save(eventType).Error; err != nil {
		return nil, err
	}
	return eventType, nil
}

func (r *eventTypeRepository) Delete(id uint) error {
	return r.db.Delete(&model.EventType{}, id).Error
}

func (r *eventTypeRepository) GetLovEventType() ([]dao.LovEventTypeDAO, error) {
	// Struct sementara untuk menampung JSON permissions sebagai string
	log.Println("di repo sebelum call db")

	var resultRows []dao.LovEventTypeDAO

	query := `
		select et.code, et.name from event_types et where et.active is true;

	`

	err := r.db.Raw(query).Scan(&resultRows).Error
	if err != nil {
		return nil, err
	}
	log.Printf("isi invitations repo : %+v", resultRows)

	return resultRows, nil
}
