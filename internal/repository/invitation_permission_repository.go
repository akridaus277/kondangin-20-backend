package repository

import (
	model "kondangin-backend/internal/model"

	"gorm.io/gorm"
)

type InvitationPermissionRepository interface {
	FindPermission(invitationID uint, userID uint) (*model.InvitationPermission, error)
	CreatePermission(permission *model.InvitationPermission) error
	UpdatePermission(permission *model.InvitationPermission) error
}

type invitationPermissionRepository struct {
	db *gorm.DB
}

func NewInvitationPermissionRepository(db *gorm.DB) InvitationPermissionRepository {
	return &invitationPermissionRepository{db}
}

func (r *invitationPermissionRepository) FindPermission(invitationID uint, userID uint) (*model.InvitationPermission, error) {
	var perm model.InvitationPermission
	err := r.db.Where("invitation_id = ? AND user_id = ?", invitationID, userID).First(&perm).Error
	return &perm, err
}

func (r *invitationPermissionRepository) CreatePermission(permission *model.InvitationPermission) error {
	return r.db.Create(permission).Error
}

func (r *invitationPermissionRepository) UpdatePermission(permission *model.InvitationPermission) error {
	return r.db.Save(permission).Error
}
