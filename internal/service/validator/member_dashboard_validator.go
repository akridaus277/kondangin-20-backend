package service

import (
	"errors"
	"kondangin-backend/internal/dto"
	"kondangin-backend/internal/repository"
	"kondangin-backend/internal/utils"
	globalValidator "kondangin-backend/internal/validator"

	"github.com/gin-gonic/gin"
)

type MemberDashboardValidator interface {
	ValidateCreateInvitation(c *gin.Context, req dto.CreateInvitationRequest) error
}

type memberDashboardValidator struct {
	invitationRepo           repository.InvitationRepository
	invitationPermissionRepo repository.InvitationPermissionRepository
	eventTypeRepo            repository.EventTypeRepository
}

func NewMemberDashboardValidator(
	invitationRepo repository.InvitationRepository,
	invitationPermissionRepo repository.InvitationPermissionRepository,
	eventTypeRepo repository.EventTypeRepository) MemberDashboardValidator {
	return &memberDashboardValidator{invitationRepo, invitationPermissionRepo, eventTypeRepo}
}

func (s *memberDashboardValidator) ValidateCreateInvitation(c *gin.Context, req dto.CreateInvitationRequest) error {
	errNull := s.NullCheckCreateInvitation(c, req)
	if errNull != nil {
		return errors.New("")

	}
	_, err := s.eventTypeRepo.FindByCodeAndActive(req.EventType, true)
	if err != nil {
		utils.SendNotFoundError(c, "Event Type is not found", nil)
		return errors.New("")
	}

	_, err = s.invitationRepo.FindBySubdomain(req.Subdomain)
	if err == nil {
		utils.SendBadRequestError(c, "Subdomain is already exists", nil)
		return errors.New("")
	}

	const (
		DateFormatYYYYMMDD = "2006-01-02" // Harus begini meskipun tidak ingin hardcoded di tempat lain
	)
	if !utils.IsValidDateFormat(req.Date, DateFormatYYYYMMDD) {
		utils.SendBadRequestError(c, "Invalid date format, expected YYYY-mm-dd", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardValidator) NullCheckCreateInvitation(c *gin.Context, req dto.CreateInvitationRequest) error {
	if globalValidator.IsNilOrEmpty(req.Date) {
		utils.SendBadRequestError(c, "Date is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.EventType) {
		utils.SendBadRequestError(c, "Event Type is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Title) {
		utils.SendBadRequestError(c, "Title is empty", nil)
		return errors.New("")
	}

	return nil

}
