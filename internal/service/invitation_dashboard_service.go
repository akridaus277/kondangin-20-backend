package service

import (
	"errors"
	"io/ioutil"
	"kondangin-backend/internal/dto"
	model "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
	"kondangin-backend/internal/utils"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvitationDashboardService interface {
	CreateInvitation(userID uint, req dto.CreateInvitationRequest) error
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

func (s *invitationDashboardService) CreateInvitation(userID uint, req dto.CreateInvitationRequest) error {
	// Baca file default data JSON
	defaultDataJSONPath := filepath.Join("default_json", "invitation_data_json.json")
	defaultDataJSONFile, err := os.Open(defaultDataJSONPath)
	if err != nil {
		return err
	}
	defer defaultDataJSONFile.Close()

	defaultDataJSONByteValue, err := ioutil.ReadAll(defaultDataJSONFile)
	if err != nil {
		return err
	}

	// Baca file default data JSON
	defaultPropertyJSONPath := filepath.Join("default_json", "invitation_property_json.json")
	defaultPropertyJSONFile, err := os.Open(defaultPropertyJSONPath)
	if err != nil {
		return err
	}
	defer defaultPropertyJSONFile.Close()

	defaultPropertyJSONByteValue, err := ioutil.ReadAll(defaultPropertyJSONFile)
	if err != nil {
		return err
	}

	// Gunakan default JSON sebagai isi data_json
	defaultDataJSON := string(defaultDataJSONByteValue)
	defaultPropertyJSON := string(defaultPropertyJSONByteValue)

	inv := &model.Invitation{
		Subdomain:    req.Subdomain,
		DataJSON:     defaultDataJSON,
		PropertyJSON: defaultPropertyJSON,
		UserID:       userID,
	}

	return s.invitationRepo.Create(inv)
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

	existing, err := s.invitationPermissionRepo.FindPermission(inv.ID, req.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Buat baru
		return s.invitationPermissionRepo.CreatePermission(&model.InvitationPermission{
			InvitationID: inv.ID,
			UserID:       req.UserID,
			Permissions:  req.Permissions,
		})
	} else if err != nil {
		return err
	}

	// Update existing
	existing.Permissions = req.Permissions
	return s.invitationPermissionRepo.UpdatePermission(existing)
}
