package dto

import "mime/multipart"

type MainEventUpdateRequest struct {
	Subdomain string
	Payload   MainEventUpdatePayload
}

type MainEventUpdatePayload struct {
	Title                string
	MainEventDate        string
	MainEventLocation    string
	MainEventLocationMap string
	CoverPhoto           *multipart.FileHeader
}

type MainEventGetRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}
