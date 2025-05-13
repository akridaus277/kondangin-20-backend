package service

import (
	"errors"
	"kondangin-backend/internal/dto"
	model "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
	"kondangin-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvitationDashboardService interface {
	GetInvitationBySubdomain(subdomain string) (*model.Invitation, error)
	AddPermission(c *gin.Context, req dto.AddInvitationPermissionRequest, requesterID uint) error
}

type invitationDashboardService struct {
	invitationRepo           repository.InvitationRepository
	invitationPermissionRepo repository.InvitationPermissionRepository
}

func NewInvitationDashboardService(invitationRepo repository.InvitationRepository, invitationPermissionRepo repository.InvitationPermissionRepository) InvitationDashboardService {
	return &invitationDashboardService{invitationRepo, invitationPermissionRepo}
}

func (s *invitationDashboardService) GetInvitationBySubdomain(subdomain string) (*model.Invitation, error) {
	return s.invitationRepo.FindBySubdomain(subdomain)
}

func (s *invitationDashboardService) AddPermission(c *gin.Context, req dto.AddInvitationPermissionRequest, requesterID uint) error {
	inv, err := s.invitationRepo.FindBySubdomain(req.Subdomain)
	if err != nil {
		return err
	}
	// Pastikan hanya pemilik yang bisa memberikan izin
	if inv.UserID != requesterID {
		utils.SendUnauthorizedError(c, "Unauthorized", nil)
		return nil
	}

	existing, err := s.invitationPermissionRepo.FindByInvitationAndUser(inv.ID, req.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Buat baru
		_, err = s.invitationPermissionRepo.Create(&model.InvitationPermission{
			InvitationID: inv.ID,
			UserID:       req.UserID,
			Permissions:  req.Permissions,
		})
		if err != nil {
			utils.SendInternalServerError(c, "Failed to create invitation permisison", nil)
			return errors.New("")
		}
	} else if err != nil {
		return err
	}

	// Update existing
	existing.Permissions = req.Permissions
	_, err = s.invitationPermissionRepo.Update(existing)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update invitation", nil)
		return errors.New("")
	}

	return nil
}
