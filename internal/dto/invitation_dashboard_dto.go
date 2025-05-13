package dto

type GetDataJSONRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}

type AddInvitationPermissionRequest struct {
	Subdomain   string   `json:"subdomain" binding:"required"`
	UserID      uint     `json:"user_id" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}
