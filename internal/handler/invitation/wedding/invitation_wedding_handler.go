package handler

import (
	dto "kondangin-backend/internal/dto/invitation/wedding"
	models "kondangin-backend/internal/model"
	service "kondangin-backend/internal/service/invitation/wedding"
	"kondangin-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type InvitationWeddingHandler struct {
	service service.WeddingService
}

func NewInvitationWeddingHandler(service service.WeddingService) *InvitationWeddingHandler {
	return &InvitationWeddingHandler{service}
}

func (h *InvitationWeddingHandler) UpdateMainEvent(c *gin.Context) {

	// Ambil field form
	subdomain := c.PostForm("subdomain")
	title := c.PostForm("payload.title")
	mainEventDate := c.PostForm("payload.mainEventDate")
	mainEventLocation := c.PostForm("payload.mainEventLocation")
	mainEventLocationMap := c.PostForm("payload.mainEventLocationMap")
	file, _ := c.FormFile("payload.coverPhoto")

	// Buat struct DTO
	req := dto.MainEventUpdateRequest{
		Subdomain: subdomain,
		Payload: dto.MainEventUpdatePayload{
			Title:                title,
			MainEventDate:        mainEventDate,
			MainEventLocation:    mainEventLocation,
			MainEventLocationMap: mainEventLocationMap,
			CoverPhoto:           file,
		},
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

	err := h.service.UpdateMainEvent(c, req, user)
	if err != nil {
		return
	}
	utils.SendSuccess(c, "Main Event Successfully Updated", nil)

}

func (h *InvitationWeddingHandler) GetMainEvent(c *gin.Context) {
	var req dto.MainEventGetRequest

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

	mainEvent, err := h.service.GetMainEvent(c, req, user)
	if err != nil {
		return
	}
	utils.SendSuccess(c, "Main Event Successfully Retrieved", mainEvent)

}
