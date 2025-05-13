package service

import (
	models "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
)

type InvitationPermissionService interface {
	Create(permission *models.InvitationPermission) (*models.InvitationPermission, error)
	Update(permission *models.InvitationPermission) (*models.InvitationPermission, error)
	Delete(id uint) error
	GetByID(id uint) (*models.InvitationPermission, error)
	GetAll() ([]models.InvitationPermission, error)
	GetByInvitationAndUser(invitationID, userID uint) (*models.InvitationPermission, error)
}

type invitationPermissionService struct {
	repo repository.InvitationPermissionRepository
}

func NewInvitationPermissionService(repo repository.InvitationPermissionRepository) InvitationPermissionService {
	return &invitationPermissionService{repo}
}

func (s *invitationPermissionService) Create(permission *models.InvitationPermission) (*models.InvitationPermission, error) {
	return s.repo.Create(permission)
}

func (s *invitationPermissionService) Update(permission *models.InvitationPermission) (*models.InvitationPermission, error) {
	return s.repo.Update(permission)
}

func (s *invitationPermissionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *invitationPermissionService) GetByID(id uint) (*models.InvitationPermission, error) {
	return s.repo.FindByID(id)
}

func (s *invitationPermissionService) GetAll() ([]models.InvitationPermission, error) {
	return s.repo.FindAll()
}

func (s *invitationPermissionService) GetByInvitationAndUser(invitationID, userID uint) (*models.InvitationPermission, error) {
	return s.repo.FindByInvitationAndUser(invitationID, userID)
}
