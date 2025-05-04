package handler

import (
	"kondangin-backend/internal/dto"
	"kondangin-backend/internal/service"
	"kondangin-backend/internal/utils"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

type InvitationDashboardHandler struct {
	service service.InvitationDashboardService
}

func NewInvitationDashboardHandler(service service.InvitationDashboardService) *InvitationDashboardHandler {
	return &InvitationDashboardHandler{service}
}

func (h *InvitationDashboardHandler) CreateInvitation(c *gin.Context) {
	var req dto.CreateInvitationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input", gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorizedError(c, "Unauthorized", nil)
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		utils.SendUnauthorizedError(c, "Invalid user ID type", nil)
		return
	}

	if err := h.service.CreateInvitation(userID, req); err != nil {
		utils.SendInternalServerError(c, "Failed to create invitation", gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, "Invitation created", nil)
}

func (h *InvitationDashboardHandler) GetInvitationData(c *gin.Context) {
	var req dto.GetDataJSONRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input", gin.H{"error": err.Error()})
		return
	}

	inv, err := h.service.GetInvitationBySubdomain(req.Subdomain)
	if err != nil {
		utils.SendNotFoundError(c, "Invitation not found", nil)
		return
	}

	// Unmarshal data_json
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &data); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", gin.H{"error": err.Error()})
		return
	}

	// Tambahkan subdomain ke response
	data["subdomain"] = inv.Subdomain

	utils.SendSuccess(c, "Succesfully get invitation data json", data)

}

func (h *InvitationDashboardHandler) AddInvitationPermission(c *gin.Context) {
	var req dto.AddInvitationPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input", gin.H{"error": err.Error()})
		return
	}

	requesterIDRaw, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorizedError(c, "Unauthorized", nil)
		return
	}

	requesterID, ok := requesterIDRaw.(uint)
	if !ok {
		utils.SendUnauthorizedError(c, "Invalid user ID", nil)
		return
	}

	if err := h.service.AddPermission(c, req, requesterID); err != nil {
		utils.SendInternalServerError(c, "Failed to add permission", gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, "Permission added", nil)
}
