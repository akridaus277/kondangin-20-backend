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
