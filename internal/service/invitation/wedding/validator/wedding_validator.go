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
	ValidateUpdateCouples(c *gin.Context, req weddingDto.CouplesUpdateRequest, userId uint) error
	ValidateGetCouples(c *gin.Context, req weddingDto.CouplesGetRequest, userId uint) error
	ValidateUpdateEventSchedules(c *gin.Context, req weddingDto.EventSchedulesUpdateRequest, userId uint) error
	ValidateGetEventSchedules(c *gin.Context, req weddingDto.EventSchedulesGetRequest, userId uint) error
	ValidateUpdatePhotoTemplates(c *gin.Context, req weddingDto.PhotoTemplatesUpdateRequest, userId uint) error
	ValidateGetPhotoTemplates(c *gin.Context, req weddingDto.PhotoTemplatesGetRequest, userId uint) error
	ValidateUpdatePhotoGalleries(c *gin.Context, req weddingDto.PhotoGalleriesUpdateRequest, userId uint) error
	ValidateGetPhotoGalleries(c *gin.Context, req weddingDto.PhotoGalleriesGetRequest, userId uint) error
	ValidateUpdateVideoGalleries(c *gin.Context, req weddingDto.VideoGalleriesUpdateRequest, userId uint) error
	ValidateGetVideoGalleries(c *gin.Context, req weddingDto.VideoGalleriesGetRequest, userId uint) error
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

func (s *memberDashboardValidator) ValidateUpdateCouples(c *gin.Context, req weddingDto.CouplesUpdateRequest, userId uint) error {
	errNull := s.NullCheckUpdateCouples(c, req)
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

	return nil

}

func (s *memberDashboardValidator) ValidateGetCouples(c *gin.Context, req weddingDto.CouplesGetRequest, userId uint) error {
	errNull := s.NullCheckGetCouples(c, req)
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

func (s *memberDashboardValidator) NullCheckUpdateCouples(c *gin.Context, req weddingDto.CouplesUpdateRequest) error {
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload) {
		utils.SendBadRequestError(c, "Payload is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(len(req.Payload) < 2) {
		utils.SendBadRequestError(c, "Couples must be 2 people", nil)
		return errors.New("")
	}
	for i := 0; i < len(req.Payload); i++ {
		if globalValidator.IsNilOrEmpty(req.Payload[i].Name) {
			utils.SendBadRequestError(c, fmt.Sprintf("Couple %d Name is empty", i+1), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].NickName) {
			utils.SendBadRequestError(c, fmt.Sprintf("Couple %d Nickname is empty", i+1), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].AdditionalInfo) {
			utils.SendBadRequestError(c, fmt.Sprintf("Couple %d Additional Info is empty", i+1), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].Photo) {
			utils.SendBadRequestError(c, fmt.Sprintf("Couple %d Photo is empty", i+1), nil)
			return errors.New("")
		}
	}

	return nil

}

func (s *memberDashboardValidator) NullCheckGetCouples(c *gin.Context, req weddingDto.CouplesGetRequest) error {
	fmt.Println("NullCheckGetMainEvent : ", req)
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardValidator) ValidateUpdateEventSchedules(c *gin.Context, req weddingDto.EventSchedulesUpdateRequest, userId uint) error {
	errNull := s.NullCheckUpdateEventSchedules(c, req)
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

	return nil

}

func (s *memberDashboardValidator) ValidateGetEventSchedules(c *gin.Context, req weddingDto.EventSchedulesGetRequest, userId uint) error {
	errNull := s.NullCheckGetEventSchedules(c, req)
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

func (s *memberDashboardValidator) NullCheckUpdateEventSchedules(c *gin.Context, req weddingDto.EventSchedulesUpdateRequest) error {
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload) {
		utils.SendBadRequestError(c, "Payload is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(len(req.Payload) < 1) {
		utils.SendBadRequestError(c, "Event Schedules must be at least 1 event", nil)
		return errors.New("")
	}
	for i := 0; i < len(req.Payload); i++ {
		if globalValidator.IsNilOrEmpty(req.Payload[i].Title) {
			utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Location is empty", i+1), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].AddressLocation) {
			utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Address Location is empty", i), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].Date) {
			utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Date is empty", i), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].CustomHour) {
			if globalValidator.IsNilOrEmpty(req.Payload[i].StartHour) {
				utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Start hour is empty", i), nil)
				return errors.New("")
			}
			if globalValidator.IsNilOrEmpty(req.Payload[i].EndHour) {
				utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d End hour is empty", i), nil)
				return errors.New("")
			}
			if globalValidator.IsNilOrEmpty(req.Payload[i].Timezone) {
				utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Timezone is empty", i), nil)
				return errors.New("")
			}
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].LocationLink) {
			utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Location Link is empty", i), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].Description) {
			utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Description is empty", i), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].Icon) {
			utils.SendBadRequestError(c, fmt.Sprintf("Event Schedule %d Icon is empty", i), nil)
			return errors.New("")
		}
	}

	return nil

}

func (s *memberDashboardValidator) NullCheckGetEventSchedules(c *gin.Context, req weddingDto.EventSchedulesGetRequest) error {
	fmt.Println("NullCheckGetMainEvent : ", req)
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardValidator) ValidateUpdatePhotoTemplates(c *gin.Context, req weddingDto.PhotoTemplatesUpdateRequest, userId uint) error {
	errNull := s.NullCheckUpdatePhotoTemplates(c, req)
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

	return nil

}

func (s *memberDashboardValidator) ValidateGetPhotoTemplates(c *gin.Context, req weddingDto.PhotoTemplatesGetRequest, userId uint) error {
	errNull := s.NullCheckGetPhotoTemplates(c, req)
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

func (s *memberDashboardValidator) NullCheckUpdatePhotoTemplates(c *gin.Context, req weddingDto.PhotoTemplatesUpdateRequest) error {
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload) {
		utils.SendBadRequestError(c, "Payload is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(len(req.Payload) < 1) {
		utils.SendBadRequestError(c, "Photo Templates must be at least 1 photo", nil)
		return errors.New("")
	}
	for i := 0; i < len(req.Payload); i++ {
		if globalValidator.IsNilOrEmpty(req.Payload[i].LinkPhoto) {
			if globalValidator.IsNilOrEmpty(req.Payload[i].Photo) {
				utils.SendBadRequestError(c, fmt.Sprintf("Photo Templates %d Photo is empty", i), nil)
				return errors.New("")
			}
		}
	}

	return nil

}

func (s *memberDashboardValidator) NullCheckGetPhotoTemplates(c *gin.Context, req weddingDto.PhotoTemplatesGetRequest) error {
	fmt.Println("NullCheckGetMainEvent : ", req)
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardValidator) ValidateUpdatePhotoGalleries(c *gin.Context, req weddingDto.PhotoGalleriesUpdateRequest, userId uint) error {
	errNull := s.NullCheckUpdatePhotoGalleries(c, req)
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

	return nil

}

func (s *memberDashboardValidator) ValidateGetPhotoGalleries(c *gin.Context, req weddingDto.PhotoGalleriesGetRequest, userId uint) error {
	errNull := s.NullCheckGetPhotoGalleries(c, req)
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

func (s *memberDashboardValidator) NullCheckUpdatePhotoGalleries(c *gin.Context, req weddingDto.PhotoGalleriesUpdateRequest) error {
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload) {
		utils.SendBadRequestError(c, "Payload is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(len(req.Payload) < 1) {
		utils.SendBadRequestError(c, "Photo Galleries must be at least 1 photo", nil)
		return errors.New("")
	}
	for i := 0; i < len(req.Payload); i++ {
		if globalValidator.IsNilOrEmpty(req.Payload[i].LinkPhoto) {
			if globalValidator.IsNilOrEmpty(req.Payload[i].Photo) {
				utils.SendBadRequestError(c, fmt.Sprintf("Photo Galleries %d Photo is empty", i), nil)
				return errors.New("")
			}
		}
	}

	return nil

}

func (s *memberDashboardValidator) NullCheckGetPhotoGalleries(c *gin.Context, req weddingDto.PhotoGalleriesGetRequest) error {
	fmt.Println("NullCheckGetMainEvent : ", req)
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardValidator) ValidateUpdateVideoGalleries(c *gin.Context, req weddingDto.VideoGalleriesUpdateRequest, userId uint) error {
	errNull := s.NullCheckUpdateVideoGalleries(c, req)
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

	return nil

}

func (s *memberDashboardValidator) ValidateGetVideoGalleries(c *gin.Context, req weddingDto.VideoGalleriesGetRequest, userId uint) error {
	errNull := s.NullCheckGetVideoGalleries(c, req)
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

func (s *memberDashboardValidator) NullCheckUpdateVideoGalleries(c *gin.Context, req weddingDto.VideoGalleriesUpdateRequest) error {
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(req.Payload) {
		utils.SendBadRequestError(c, "Payload is empty", nil)
		return errors.New("")
	}
	if globalValidator.IsNilOrEmpty(len(req.Payload) < 1) {
		utils.SendBadRequestError(c, "Video Galleries must be at least 1 photo", nil)
		return errors.New("")
	}
	for i := 0; i < len(req.Payload); i++ {
		if globalValidator.IsNilOrEmpty(req.Payload[i].Platform) {
			utils.SendBadRequestError(c, fmt.Sprintf("Video Gallery %d Platform is empty", i+1), nil)
			return errors.New("")
		}
		if globalValidator.IsNilOrEmpty(req.Payload[i].LinkVideo) {
			utils.SendBadRequestError(c, fmt.Sprintf("Video Gallery %d Link Video is empty", i+1), nil)
			return errors.New("")
		}
	}

	return nil

}

func (s *memberDashboardValidator) NullCheckGetVideoGalleries(c *gin.Context, req weddingDto.VideoGalleriesGetRequest) error {
	fmt.Println("NullCheckGetMainEvent : ", req)
	if globalValidator.IsNilOrEmpty(req.Subdomain) {
		utils.SendBadRequestError(c, "Subdomain is empty", nil)
		return errors.New("")
	}

	return nil

}
