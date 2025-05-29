package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png" // <- WAJIB untuk decode PNG
	weddingDto "kondangin-backend/internal/dto/invitation/wedding"
	models "kondangin-backend/internal/model"
	modelService "kondangin-backend/internal/service"
	validator "kondangin-backend/internal/service/invitation/wedding/validator"
	"kondangin-backend/internal/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type WeddingService interface {
	UpdateMainEvent(c *gin.Context, input weddingDto.MainEventUpdateRequest, user *models.User) error
	GetMainEvent(c *gin.Context, input weddingDto.MainEventGetRequest, user *models.User) (interface{}, error)
	UpdateCouples(c *gin.Context, input weddingDto.CouplesUpdateRequest, user *models.User) error
	GetCouples(c *gin.Context, input weddingDto.CouplesGetRequest, user *models.User) (interface{}, error)
	UpdateEventSchedules(c *gin.Context, input weddingDto.EventSchedulesUpdateRequest, user *models.User) error
	GetEventSchedules(c *gin.Context, input weddingDto.EventSchedulesGetRequest, user *models.User) (interface{}, error)
	UpdatePhotoTemplates(c *gin.Context, input weddingDto.PhotoTemplatesUpdateRequest, user *models.User) error
	GetPhotoTemplates(c *gin.Context, input weddingDto.PhotoTemplatesGetRequest, user *models.User) (interface{}, error)
	UpdatePhotoGalleries(c *gin.Context, input weddingDto.PhotoGalleriesUpdateRequest, user *models.User) error
	GetPhotoGalleries(c *gin.Context, input weddingDto.PhotoGalleriesGetRequest, user *models.User) (interface{}, error)
	UpdateVideoGalleries(c *gin.Context, input weddingDto.VideoGalleriesUpdateRequest, user *models.User) error
	GetVideoGalleries(c *gin.Context, input weddingDto.VideoGalleriesGetRequest, user *models.User) (interface{}, error)
}

type weddingService struct {
	validator         validator.WeddingValidator
	invitationService modelService.InvitationService
}

func NewWeddingService(validator validator.WeddingValidator, invitationService modelService.InvitationService) WeddingService {
	return &weddingService{validator, invitationService}
}

func (s *weddingService) UpdateMainEvent(c *gin.Context, input weddingDto.MainEventUpdateRequest, user *models.User) error {
	errValidation := s.validator.ValidateUpdateMainEvent(c, input, user.ID)
	if errValidation != nil {
		return errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return errors.New("")
	}

	weddingMainEvent := invitationDataJson["mainEvent"]
	weddingMainEvent.(map[string]interface{})["title"] = input.Payload.Title
	weddingMainEvent.(map[string]interface{})["mainEventDate"] = input.Payload.MainEventDate
	weddingMainEvent.(map[string]interface{})["mainEventLocation"] = input.Payload.MainEventLocation
	weddingMainEvent.(map[string]interface{})["mainEventLocationMap"] = input.Payload.MainEventLocationMap

	if input.Payload.CoverPhoto != nil {
		_, url, err := SaveUploadedFile(input.Payload.CoverPhoto)
		if err != nil {
			utils.SendInternalServerError(c, "Failed to save image", nil)
			return errors.New("")
		}
		err = DeleteUploadedFile(weddingMainEvent.(map[string]interface{})["coverPhoto"].(string))
		if err != nil {
			utils.SendInternalServerError(c, "Failed to delete old image", nil)
			return errors.New("")
		}
		weddingMainEvent.(map[string]interface{})["coverPhoto"] = url
	}

	invitationDataJson["mainEvent"] = weddingMainEvent

	// Marshal kembali ke string JSON
	updatedDataJSONBytes, err := json.Marshal(invitationDataJson)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated data to JSON", nil)
		return errors.New("")
	}
	// Konversi ke string
	updatedDataJSONString := string(updatedDataJSONBytes)

	// update model
	inv.DataJSON = updatedDataJSONString
	_, err = s.invitationService.UpdateInvitation(c, inv)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update invitation", nil)
		return errors.New("")
	}

	return nil
}

func (s *weddingService) UpdateCouples(c *gin.Context, input weddingDto.CouplesUpdateRequest, user *models.User) error {
	errValidation := s.validator.ValidateUpdateCouples(c, input, user.ID)
	if errValidation != nil {
		return errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return errors.New("")
	}

	updatedWeddingCouples := []interface{}{}
	for i := 0; i < len(input.Payload); i++ {
		updatedWeddingCouples[i].(map[string]interface{})["name"] = input.Payload[i].Name
		updatedWeddingCouples[i].(map[string]interface{})["nickname"] = input.Payload[i].NickName
		updatedWeddingCouples[i].(map[string]interface{})["gender"] = input.Payload[i].Gender
		updatedWeddingCouples[i].(map[string]interface{})["additionalInfo"] = input.Payload[i].AdditionalInfo
		updatedWeddingCouples[i].(map[string]interface{})["tiktok"] = input.Payload[i].Tiktok
		updatedWeddingCouples[i].(map[string]interface{})["instagram"] = input.Payload[i].Instagram
		updatedWeddingCouples[i].(map[string]interface{})["facebook"] = input.Payload[i].Facebook
		if input.Payload[i].Photo != nil {
			_, url, err := SaveUploadedFile(input.Payload[i].Photo)
			if err != nil {
				utils.SendInternalServerError(c, "Failed to save image", nil)
				return errors.New("")
			}
			err = DeleteUploadedFile(updatedWeddingCouples[i].(map[string]interface{})["linkPhoto"].(string))
			if err != nil {
				utils.SendInternalServerError(c, "Failed to delete old image", nil)
				return errors.New("")
			}
			updatedWeddingCouples[i].(map[string]interface{})["linkPhoto"] = url
		}
	}

	invitationDataJson["couples"] = updatedWeddingCouples

	// Marshal kembali ke string JSON
	updatedDataJSONBytes, err := json.Marshal(invitationDataJson)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated data to JSON", nil)
		return errors.New("")
	}
	// Konversi ke string
	updatedDataJSONString := string(updatedDataJSONBytes)

	// update model
	inv.DataJSON = updatedDataJSONString
	_, err = s.invitationService.UpdateInvitation(c, inv)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update invitation", nil)
		return errors.New("")
	}

	return nil
}

func (s *weddingService) UpdateEventSchedules(c *gin.Context, input weddingDto.EventSchedulesUpdateRequest, user *models.User) error {
	errValidation := s.validator.ValidateUpdateEventSchedules(c, input, user.ID)
	if errValidation != nil {
		return errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return errors.New("")
	}

	updatedWeddingEventSchedules := []interface{}{}
	for i := 0; i < len(input.Payload); i++ {
		updatedWeddingEventSchedules[i].(map[string]interface{})["title"] = input.Payload[i].Title
		updatedWeddingEventSchedules[i].(map[string]interface{})["addressLocation"] = input.Payload[i].AddressLocation
		updatedWeddingEventSchedules[i].(map[string]interface{})["date"] = input.Payload[i].Date
		updatedWeddingEventSchedules[i].(map[string]interface{})["startHour"] = input.Payload[i].StartHour
		updatedWeddingEventSchedules[i].(map[string]interface{})["endHour"] = input.Payload[i].EndHour
		updatedWeddingEventSchedules[i].(map[string]interface{})["timezone"] = input.Payload[i].Timezone
		updatedWeddingEventSchedules[i].(map[string]interface{})["customHour"] = input.Payload[i].CustomHour
		updatedWeddingEventSchedules[i].(map[string]interface{})["locationLink"] = input.Payload[i].LocationLink
		updatedWeddingEventSchedules[i].(map[string]interface{})["description"] = input.Payload[i].Description
		updatedWeddingEventSchedules[i].(map[string]interface{})["icon"] = input.Payload[i].Icon
	}

	invitationDataJson["couples"] = updatedWeddingEventSchedules

	// Marshal kembali ke string JSON
	updatedDataJSONBytes, err := json.Marshal(invitationDataJson)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated data to JSON", nil)
		return errors.New("")
	}
	// Konversi ke string
	updatedDataJSONString := string(updatedDataJSONBytes)

	// update model
	inv.DataJSON = updatedDataJSONString
	_, err = s.invitationService.UpdateInvitation(c, inv)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update invitation", nil)
		return errors.New("")
	}

	return nil
}

func (s *weddingService) UpdatePhotoTemplates(c *gin.Context, input weddingDto.PhotoTemplatesUpdateRequest, user *models.User) error {
	errValidation := s.validator.ValidateUpdatePhotoTemplates(c, input, user.ID)
	if errValidation != nil {
		return errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return errors.New("")
	}

	updatedPhotoTemplates := []interface{}{}
	for i := 0; i < len(input.Payload); i++ {
		if input.Payload[i].Photo != nil {
			_, url, err := SaveUploadedFile(input.Payload[i].Photo)
			if err != nil {
				utils.SendInternalServerError(c, "Failed to save image", nil)
				return errors.New("")
			}
			err = DeleteUploadedFile(updatedPhotoTemplates[i].(map[string]interface{})["linkPhoto"].(string))
			if err != nil {
				utils.SendInternalServerError(c, "Failed to delete old image", nil)
				return errors.New("")
			}
			updatedPhotoTemplates[i].(map[string]interface{})["linkPhoto"] = url
		}
	}

	invitationDataJson["photoTemplates"] = updatedPhotoTemplates

	// Marshal kembali ke string JSON
	updatedDataJSONBytes, err := json.Marshal(invitationDataJson)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated data to JSON", nil)
		return errors.New("")
	}
	// Konversi ke string
	updatedDataJSONString := string(updatedDataJSONBytes)

	// update model
	inv.DataJSON = updatedDataJSONString
	_, err = s.invitationService.UpdateInvitation(c, inv)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update invitation", nil)
		return errors.New("")
	}

	return nil
}

func (s *weddingService) UpdatePhotoGalleries(c *gin.Context, input weddingDto.PhotoGalleriesUpdateRequest, user *models.User) error {
	errValidation := s.validator.ValidateUpdatePhotoGalleries(c, input, user.ID)
	if errValidation != nil {
		return errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return errors.New("")
	}

	updatedPhotoGalleries := []interface{}{}
	for i := 0; i < len(input.Payload); i++ {
		if input.Payload[i].Photo != nil {
			_, url, err := SaveUploadedFile(input.Payload[i].Photo)
			if err != nil {
				utils.SendInternalServerError(c, "Failed to save image", nil)
				return errors.New("")
			}
			err = DeleteUploadedFile(updatedPhotoGalleries[i].(map[string]interface{})["linkPhoto"].(string))
			if err != nil {
				utils.SendInternalServerError(c, "Failed to delete old image", nil)
				return errors.New("")
			}
			updatedPhotoGalleries[i].(map[string]interface{})["linkPhoto"] = url
		}
	}

	invitationDataJson["photoGalleries"] = updatedPhotoGalleries

	// Marshal kembali ke string JSON
	updatedDataJSONBytes, err := json.Marshal(invitationDataJson)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated data to JSON", nil)
		return errors.New("")
	}
	// Konversi ke string
	updatedDataJSONString := string(updatedDataJSONBytes)

	// update model
	inv.DataJSON = updatedDataJSONString
	_, err = s.invitationService.UpdateInvitation(c, inv)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update invitation", nil)
		return errors.New("")
	}

	return nil
}

func (s *weddingService) UpdateVideoGalleries(c *gin.Context, input weddingDto.VideoGalleriesUpdateRequest, user *models.User) error {
	errValidation := s.validator.ValidateUpdateVideoGalleries(c, input, user.ID)
	if errValidation != nil {
		return errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return errors.New("")
	}

	updatedVideoGalleries := []interface{}{}
	for i := 0; i < len(input.Payload); i++ {
		updatedVideoGalleries[i].(map[string]interface{})["platform"] = input.Payload[i].Platform
		updatedVideoGalleries[i].(map[string]interface{})["linkVideo"] = input.Payload[i].LinkVideo

	}

	invitationDataJson["videoGalleries"] = updatedVideoGalleries

	// Marshal kembali ke string JSON
	updatedDataJSONBytes, err := json.Marshal(invitationDataJson)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to convert updated data to JSON", nil)
		return errors.New("")
	}
	// Konversi ke string
	updatedDataJSONString := string(updatedDataJSONBytes)

	// update model
	inv.DataJSON = updatedDataJSONString
	_, err = s.invitationService.UpdateInvitation(c, inv)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update invitation", nil)
		return errors.New("")
	}

	return nil
}

func (s *weddingService) GetMainEvent(c *gin.Context, input weddingDto.MainEventGetRequest, user *models.User) (interface{}, error) {
	errValidation := s.validator.ValidateGetMainEvent(c, input, user.ID)
	if errValidation != nil {
		return nil, errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return nil, errors.New("")
	}

	weddingMainEvent := invitationDataJson["mainEvent"]

	return weddingMainEvent, nil
}

func (s *weddingService) GetCouples(c *gin.Context, input weddingDto.CouplesGetRequest, user *models.User) (interface{}, error) {
	errValidation := s.validator.ValidateGetCouples(c, input, user.ID)
	if errValidation != nil {
		return nil, errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return nil, errors.New("")
	}

	weddingCouples := invitationDataJson["couples"]

	return weddingCouples, nil
}

func (s *weddingService) GetEventSchedules(c *gin.Context, input weddingDto.EventSchedulesGetRequest, user *models.User) (interface{}, error) {
	errValidation := s.validator.ValidateGetEventSchedules(c, input, user.ID)
	if errValidation != nil {
		return nil, errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return nil, errors.New("")
	}

	weddingCouples := invitationDataJson["couples"]

	return weddingCouples, nil
}

func (s *weddingService) GetPhotoTemplates(c *gin.Context, input weddingDto.PhotoTemplatesGetRequest, user *models.User) (interface{}, error) {
	errValidation := s.validator.ValidateGetPhotoTemplates(c, input, user.ID)
	if errValidation != nil {
		return nil, errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return nil, errors.New("")
	}

	weddingCouples := invitationDataJson["couples"]

	return weddingCouples, nil
}

func (s *weddingService) GetPhotoGalleries(c *gin.Context, input weddingDto.PhotoGalleriesGetRequest, user *models.User) (interface{}, error) {
	errValidation := s.validator.ValidateGetPhotoGalleries(c, input, user.ID)
	if errValidation != nil {
		return nil, errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return nil, errors.New("")
	}

	weddingCouples := invitationDataJson["couples"]

	return weddingCouples, nil
}

func (s *weddingService) GetVideoGalleries(c *gin.Context, input weddingDto.VideoGalleriesGetRequest, user *models.User) (interface{}, error) {
	errValidation := s.validator.ValidateGetVideoGalleries(c, input, user.ID)
	if errValidation != nil {
		return nil, errors.New("")
	}

	inv, _ := s.invitationService.GetInvitationBySubdomain(c, input.Subdomain)

	var invitationDataJson map[string]interface{}
	if err := json.Unmarshal([]byte(inv.DataJSON), &invitationDataJson); err != nil {
		utils.SendInternalServerError(c, "Failed to parse invitation data", nil)
		return nil, errors.New("")
	}

	weddingCouples := invitationDataJson["couples"]

	return weddingCouples, nil
}

func SaveUploadedFile(file *multipart.FileHeader) (string, string, error) {
	// Buat folder jika belum ada
	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", "", err
	}

	// Rename file pakai timestamp agar unik
	timestamp := time.Now().UnixNano()
	ext := ".jpg" // paksa .jpg karena kita akan encode ulang ke JPEG
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	fullPath := filepath.Join(uploadPath, filename)

	// Simpan file
	if err := saveFile(file, fullPath); err != nil {
		return "", "", err
	}

	// Return filename dan URL relatif (bisa disesuaikan base URL-nya)
	return filename, "/uploads/" + filename, nil
}

func saveFile(file *multipart.FileHeader, path string) error {
	return saveWithGin(file, path)
}

func saveWithGin(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Decode image
	img, format, err := image.Decode(src)
	if err != nil {
		fmt.Println("Error image decode:", err)
		return err
	}

	// Resize (optional – hanya jika kamu ingin membatasi max dimension)
	resizedImg := resize.Resize(0, 1080, img, resize.Lanczos3) // Tinggi maksimal 1080px

	out, err := os.Create(path)
	if err != nil {
		fmt.Println("Error create path")
		return err
	}
	defer out.Close()

	// Kompresi JPEG dengan quality tinggi tapi ukuran kecil
	opts := &jpeg.Options{Quality: 75} // Ubah ke 60–80 sesuai trial
	// Konversi semua gambar ke JPEG (lebih kecil & kompatibel)
	switch format {
	case "jpeg", "jpg", "png", "gif":
		return jpeg.Encode(out, resizedImg, opts)
	default:
		return errors.New("unsupported image format")
	}

	// _, err = out.ReadFrom(src)
	// return err
}

func DeleteUploadedFile(fileURL string) error {
	// Hilangkan awalan slash jika ada
	cleanURL := strings.TrimPrefix(fileURL, "/")

	// Bangun path relatif ke file
	fullPath := filepath.Join(".", cleanURL)

	// Hapus file
	if err := os.Remove(fullPath); err != nil {
		return err
	}

	return nil
}
