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

type CouplesUpdateRequest struct {
	Subdomain string
	Payload   []CouplesUpdatePayload
}

type CouplesUpdatePayload struct {
	Name           string
	NickName       string
	Gender         string
	AdditionalInfo string
	Tiktok         string
	Instagram      string
	Facebook       string
	Photo          *multipart.FileHeader
}

type CouplesGetRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}

type EventSchedulesUpdateRequest struct {
	Subdomain string
	Payload   []EventSchedulesUpdatePayload
}

type EventSchedulesUpdatePayload struct {
	Date            string
	Title           string
	Description     string
	AddressLocation string
	LocationLink    string
	Timezone        string
	StartHour       string
	EndHour         string
	CustomHour      string
	Icon            string
}

type EventSchedulesGetRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}

type PhotoTemplatesUpdateRequest struct {
	Subdomain string
	Payload   []PhotoTemplatesUpdatePayload
}

type PhotoTemplatesUpdatePayload struct {
	LinkPhoto string
	Photo     *multipart.FileHeader
}

type PhotoTemplatesGetRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}

type PhotoGalleriesUpdateRequest struct {
	Subdomain string
	Payload   []PhotoGalleriesUpdatePayload
}

type PhotoGalleriesUpdatePayload struct {
	LinkPhoto string
	Photo     *multipart.FileHeader
}

type PhotoGalleriesGetRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}

type VideoGalleriesUpdateRequest struct {
	Subdomain string
	Payload   []VideoGalleriesUpdatePayload
}

type VideoGalleriesUpdatePayload struct {
	Platform  string
	LinkVideo string
}

type VideoGalleriesGetRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}

type UtterancesUpdateRequest struct {
	Subdomain string
	Payload   []VideoGalleriesUpdatePayload
}

type UtterancesUpdatePayload struct {
	Platform  string
	LinkVideo string
}

type UtterancesGetRequest struct {
	Subdomain string `json:"subdomain" binding:"required"`
}
