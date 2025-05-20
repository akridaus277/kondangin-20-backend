package service

import (
	"errors"
	"fmt"
	weddingDto "kondangin-backend/internal/dto/invitation/wedding"
	"kondangin-backend/internal/repository"
	modelService "kondangin-backend/internal/service"
	"kondangin-backend/internal/utils"
	globalValidator "kondangin-backend/internal/validator"

	"github.com/gin-gonic/gin"
)

type WeddingValidator interface {
	ValidateUpdateMainEvent(c *gin.Context, req weddingDto.MainEventUpdateRequest, userId uint) error
	ValidateGetMainEvent(c *gin.Context, req weddingDto.MainEventGetRequest, userId uint) error
}

type memberDashboardValidator struct {
	invitationRepo              repository.InvitationRepository
	invitationPermissionService modelService.InvitationPermissionService
	eventTypeRepo               repository.EventTypeRepository
}

func NewWeddingValidator(
	invitationRepo repository.InvitationRepository,
	invitationPermissionService modelService.InvitationPermissionService,
	eventTypeRepo repository.EventTypeRepository) WeddingValidator {
	return &memberDashboardValidator{invitationRepo, invitationPermissionService, eventTypeRepo}
}

func (s *memberDashboardValidator) ValidateUpdateMainEvent(c *gin.Context, req weddingDto.MainEventUpdateRequest, userId uint) error {
	errNull := s.NullCheckUpdateMainEvent(c, req)
	if errNull != nil {
		return errors.New("")

	}

	inv, err := s.invitationRepo.FindBySubdomain(req.Subdomain)
	if err != nil {
		utils.SendBadRequestError(c, "Subdomain does not exist", nil)
		return errors.New("")
	}

	//nanti tambahin disini validasi permission
	hasPermission, err := s.invitationPermissionService.HasPermission(inv.ID, userId, "update")
	if !hasPermission || err != nil {
		utils.SendForbiddenError(c, "You have no permission to update", nil)
		return errors.New("")
	}

	const (
		DateFormatYYYYMMDD = "2006-01-02" // Harus begini meskipun tidak ingin hardcoded di tempat lain
	)
	if !utils.IsValidDateFormat(req.Payload.MainEventDate, DateFormatYYYYMMDD) {
		utils.SendBadRequestError(c, "Invalid date format, expected YYYY-mm-dd", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardValidator) ValidateGetMainEvent(c *gin.Context, req weddingDto.MainEventGetRequest, userId uint) error {
	errNull := s.NullCheckGetMainEvent(c, req)
	if errNull != nil {
		return errors.New("")

	}

	inv, err := s.invitationRepo.FindBySubdomain(req.Subdomain)
	if err != nil {
		utils.SendBadRequestError(c, "Subdomain does not exist", nil)
		return errors.New("")
	}

	//nanti tambahin disini validasi permission
	hasPermission, err := s.invitationPermissionService.HasPermission(inv.ID, userId, "view")
	if !hasPermission || err != nil {
		utils.SendForbiddenError(c, "You have no permission to update", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardValidator) NullCheckUpdateMainEvent(c *gin.Context, req weddingDto.MainEventUpdateRequest) error {
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload) {
		utils.SendBadRequestError(c, "Payload is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload.Title) {
		utils.SendBadRequestError(c, "Title is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload.MainEventDate) {
		utils.SendBadRequestError(c, "Date is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload.MainEventLocation) {
		utils.SendBadRequestError(c, "Location is empty", nil)
		return errors.New("")
	}
	// if globalValidator.IsNilOrEmpty(req.Payload.CoverPhoto) {
	// 	utils.SendBadRequestError(c, "Cover Photo is empty", nil)
	// 	return errors.New("")
	// }

	return nil

}

func (s *memberDashboardValidator) NullCheckGetMainEvent(c *gin.Context, req weddingDto.MainEventGetRequest) error {
	fmt.Println("NullCheckGetMainEvent : ", req)
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}

	return nil

}
