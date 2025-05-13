package service

import (
	dao "kondangin-backend/internal/dao"
	models "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
	"log"
)

type EventTypeService interface {
	Create(eventType *models.EventType) (*models.EventType, error)
	GetByID(id uint) (*models.EventType, error)
	GetAll() ([]models.EventType, error)
	Update(eventType *models.EventType) (*models.EventType, error)
	Delete(id uint) error
	GetLovEventType() ([]dao.LovEventTypeDAO, error)
}

type eventTypeService struct {
	repo repository.EventTypeRepository
}

func NewEventTypeService(repo repository.EventTypeRepository) EventTypeService {
	return &eventTypeService{repo}
}

func (s *eventTypeService) Create(eventType *models.EventType) (*models.EventType, error) {
	return s.repo.Create(eventType)
}

func (s *eventTypeService) GetByID(id uint) (*models.EventType, error) {
	return s.repo.FindByID(id)
}

func (s *eventTypeService) GetAll() ([]models.EventType, error) {
	return s.repo.FindAll()
}

func (s *eventTypeService) Update(eventType *models.EventType) (*models.EventType, error) {
	return s.repo.Update(eventType)
}

func (s *eventTypeService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *eventTypeService) GetLovEventType() ([]dao.LovEventTypeDAO, error) {
	log.Println("di invi service sebelum call repo.FindByInvitationAndUserNative")
	return s.repo.GetLovEventType()
}
