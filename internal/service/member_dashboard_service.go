package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"kondangin-backend/internal/dto"
	model "kondangin-backend/internal/model"
	validator "kondangin-backend/internal/service/validator"
	"kondangin-backend/internal/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type MemberDashboardService interface {
	CreateInvitation(c *gin.Context, user *model.User, req dto.CreateInvitationRequest) error
	GetInvitation(c *gin.Context, user *model.User) ([]dto.InvitationListResponse, error)
	GetLovEventType(c *gin.Context) ([]dto.LovEventTypeResponse, error)
}

type memberDashboardService struct {
	invitationService           InvitationService
	invitationPermissionService InvitationPermissionService
	eventTypeService            EventTypeService
	memberDashboardValidator    validator.MemberDashboardValidator
}

func NewMemberDashboardService(
	invitationService InvitationService,
	invitationPermissionService InvitationPermissionService,
	eventTypeService EventTypeService,
	memberDashboardValidator validator.MemberDashboardValidator) MemberDashboardService {
	return &memberDashboardService{invitationService, invitationPermissionService, eventTypeService, memberDashboardValidator}
}

func (s *memberDashboardService) CreateInvitation(c *gin.Context, user *model.User, req dto.CreateInvitationRequest) error {
	// Baca file default data JSON
	defaultDataJSONPath := filepath.Join("default_json", "invitation_data_json_wedding.json")
	defaultDataJSONFile, err := os.Open(defaultDataJSONPath)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to get default data json", nil)
		return errors.New("")
	}
	defer defaultDataJSONFile.Close()

	defaultDataJSONByteValue, err := ioutil.ReadAll(defaultDataJSONFile)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to read default data json", nil)
		return errors.New("")
	}

	// Baca file default data JSON
	defaultPropertyJSONPath := filepath.Join("default_json", "invitation_property_json_wedding.json")
	defaultPropertyJSONFile, err := os.Open(defaultPropertyJSONPath)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to get property json", nil)
		return errors.New("")
	}
	defer defaultPropertyJSONFile.Close()

	defaultPropertyJSONByteValue, err := ioutil.ReadAll(defaultPropertyJSONFile)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to read property json", nil)
		return errors.New("")
	}

	errValidation := s.memberDashboardValidator.ValidateCreateInvitation(c, req)
	if errValidation != nil {
		return errors.New("")
	}

	// Gunakan default JSON sebagai isi data_json
	defaultDataJSON := string(defaultDataJSONByteValue)
	defaultPropertyJSON := string(defaultPropertyJSONByteValue)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(defaultDataJSON), &data); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return errors.New("")
	}
	var property map[string]interface{}
	if err := json.Unmarshal([]byte(defaultPropertyJSON), &property); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation property", nil)
		return errors.New("")
	}

	// Ubah title
	data["mainEvent"].(map[string]interface{})["title"] = req.Title
	data["mainEvent"].(map[string]interface{})["mainEventDate"] = req.Date
	property["eventType"] = req.EventType

	// Marshal kembali ke string JSON
	updatedDataJSONBytes, err := json.Marshal(data)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated data to JSON", nil)
		return errors.New("")
	}
	// Marshal kembali ke string JSON
	updatedPropertyJSONBytes, err := json.Marshal(property)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated property to JSON", nil)
		return errors.New("")
	}

	// Konversi ke string
	updatedDataJSONString := string(updatedDataJSONBytes)
	updatedPropertyJSONString := string(updatedPropertyJSONBytes)

	inv := &model.Invitation{
		Subdomain:    req.Subdomain,
		DataJSON:     updatedDataJSONString,
		PropertyJSON: updatedPropertyJSONString,
		UserID:       user.ID,
		CreatedBy:    user.Username,
	}

	newInvitation, err := s.invitationService.CreateInvitation(c, inv)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create invitation", nil)
		return errors.New("")
	}

	invitationPermission := &model.InvitationPermission{
		InvitationID: newInvitation.ID,
		Permissions:  datatypes.JSONSlice[string]{"view", "update", "delete"},
		UserID:       user.ID,
		CreatedBy:    user.Username,
	}
	_, err = s.invitationPermissionService.Create(invitationPermission)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create invitation", nil)
		return errors.New("")
	}

	return nil

}

func (s *memberDashboardService) GetInvitation(c *gin.Context, user *model.User) ([]dto.InvitationListResponse, error) {
	log.Println("di member dashboard service sebelum call invitationService.GetInvitationsForUser")
	invitations, err := s.invitationService.GetInvitationsForUser(user.ID)
	if err != nil {
		return nil, err
	}
	log.Printf("isi invitations : %+v", invitations)

	var result []dto.InvitationListResponse
	for _, inv := range invitations {

		resp := dto.InvitationListResponse{
			Subdomain:   inv.InvitationSubdomain,
			Title:       inv.InvitationTitle,
			Owner:       inv.InvitationOwner,
			Status:      inv.InvitationStatus,
			CoverPhoto:  inv.InvitationCoverPhoto,
			Permissions: inv.InvitationPermissions,
		}

		result = append(result, resp)
	}

	return result, nil
}

func (s *memberDashboardService) GetLovEventType(c *gin.Context) ([]dto.LovEventTypeResponse, error) {

	lovEventTypes, err := s.eventTypeService.GetLovEventType()
	if err != nil {
		return nil, err
	}
	log.Printf("isi lovEventTypes : %+v", lovEventTypes)

	var result []dto.LovEventTypeResponse
	for _, lov := range lovEventTypes {

		resp := dto.LovEventTypeResponse{
			Code: lov.Code,
			Name: lov.Name,
		}

		result = append(result, resp)
	}

	return result, nil
}
