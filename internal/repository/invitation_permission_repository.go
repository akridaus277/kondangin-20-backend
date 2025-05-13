package repository

import (
	models "kondangin-backend/internal/model"

	"gorm.io/gorm"
)

type InvitationPermissionRepository interface {
	Create(permission *models.InvitationPermission) (*models.InvitationPermission, error)
	Update(permission *models.InvitationPermission) (*models.InvitationPermission, error)
	Delete(id uint) error
	FindByID(id uint) (*models.InvitationPermission, error)
	FindAll() ([]models.InvitationPermission, error)
	FindByInvitationAndUser(invitationID, userID uint) (*models.InvitationPermission, error)
}

type invitationPermissionRepository struct {
	db *gorm.DB
}

func NewInvitationPermissionRepository(db *gorm.DB) InvitationPermissionRepository {
	return &invitationPermissionRepository{db}
}

func (r *invitationPermissionRepository) Create(permission *models.InvitationPermission) (*models.InvitationPermission, error) {
	if err := r.db.Create(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *invitationPermissionRepository) Update(permission *models.InvitationPermission) (*models.InvitationPermission, error) {
	if err := r.db.Save(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *invitationPermissionRepository) Delete(id uint) error {
	return r.db.Delete(&models.InvitationPermission{}, id).Error
}

func (r *invitationPermissionRepository) FindByID(id uint) (*models.InvitationPermission, error) {
	var permission models.InvitationPermission
	err := r.db.Preload("Invitation").Preload("User").First(&permission, id).Error
	return &permission, err
}

func (r *invitationPermissionRepository) FindAll() ([]models.InvitationPermission, error) {
	var permissions []models.InvitationPermission
	err := r.db.Preload("Invitation").Preload("User").Find(&permissions).Error
	return permissions, err
}

func (r *invitationPermissionRepository) FindByInvitationAndUser(invitationID, userID uint) (*models.InvitationPermission, error) {
	var permission models.InvitationPermission
	err := r.db.Where("invitation_id = ? AND user_id = ?", invitationID, userID).First(&permission).Error
	return &permission, err
}
