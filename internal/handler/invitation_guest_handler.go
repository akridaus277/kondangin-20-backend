package handler

import (
	"kondangin-backend/internal/dto"
	"kondangin-backend/internal/service"
	"kondangin-backend/internal/utils"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

type InvitationGuestHandler struct {
	service service.InvitationGuestService
}

func NewInvitationGuestHandler(service service.InvitationGuestService) *InvitationGuestHandler {
	return &InvitationGuestHandler{service}
}

func (h *InvitationGuestHandler) GetInvitationData(c *gin.Context) {
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
	data["templateCode"] = "template-1"

	utils.SendSuccess(c, "Succesfully get invitation data json", data)

}
