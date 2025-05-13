package handler

import (
	"kondangin-backend/internal/dto"
	models "kondangin-backend/internal/model"
	"kondangin-backend/internal/service"
	"kondangin-backend/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type MemberDashboardHandler struct {
	memberDashboardService service.MemberDashboardService
}

func NewMemberDashboardHandler(memberDashboardService service.MemberDashboardService) *MemberDashboardHandler {
	return &MemberDashboardHandler{memberDashboardService}
}

func (h *MemberDashboardHandler) CreateInvitation(c *gin.Context) {
	var req dto.CreateInvitationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input", gin.H{"error": err.Error()})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		utils.SendUnauthorizedError(c, "Unauthorized", nil)
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		utils.SendUnauthorizedError(c, "Invalid user type", nil)
		return
	}

	err := h.memberDashboardService.CreateInvitation(c, user, req)
	if err != nil {
		return
	}

	utils.SendSuccess(c, "Invitation created", nil)
}

func (h *MemberDashboardHandler) GetInvitation(c *gin.Context) {

	userInterface, exists := c.Get("user")
	if !exists {
		utils.SendUnauthorizedError(c, "Unauthorized", nil)
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		utils.SendUnauthorizedError(c, "Invalid user type", nil)
		return
	}

	log.Println("di handler sebelum call memberDashboardService.GetInvitation")
	result, err := h.memberDashboardService.GetInvitation(c, user)
	if err != nil {
		return
	}

	utils.SendSuccess(c, "Get List Invitation Success", result)
}

func (h *MemberDashboardHandler) GetLovEventType(c *gin.Context) {

	userInterface, exists := c.Get("user")
	if !exists {
		utils.SendUnauthorizedError(c, "Unauthorized", nil)
		return
	}

	_, ok := userInterface.(*models.User)
	if !ok {
		utils.SendUnauthorizedError(c, "Invalid user type", nil)
		return
	}

	result, err := h.memberDashboardService.GetLovEventType(c)
	if err != nil {
		return
	}

	utils.SendSuccess(c, "Get LOV Event Type Success", result)
}
