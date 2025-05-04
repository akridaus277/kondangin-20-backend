package repository

import (
	model "kondangin-backend/internal/model"

	"gorm.io/gorm"
)

type InvitationRepository interface {
	Create(inv *model.Invitation) error
	FindBySubdomain(subdomain string) (*model.Invitation, error)
}

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db}
}

func (r *invitationRepository) Create(inv *model.Invitation) error {
	return r.db.Create(inv).Error
}

func (r *invitationRepository) FindBySubdomain(subdomain string) (*model.Invitation, error) {
	var invitation model.Invitation
	if err := r.db.Where("subdomain = ?", subdomain).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}
