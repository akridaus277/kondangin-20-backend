package dto

import (
	"gorm.io/datatypes"
)

type CreateInvitationRequest struct {
	Subdomain string `json:"subdomain"`
	EventType string `json:"eventType"`
	Date      string `json:"date"`
	Title     string `json:"title"`
}

type InvitationListResponse struct {
	Title       string                      `json:"title"`
	Status      string                      `json:"status"`
	Subdomain   string                      `json:"subdomain"`
	Owner       string                      `json:"owner"`
	Permissions datatypes.JSONSlice[string] `json:"permissions"`
	CoverPhoto  string                      `json:"coverPhoto"`
}

type LovEventTypeResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
