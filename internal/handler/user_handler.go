package handler

import (
	"kondangin-backend/internal/dto"
	"kondangin-backend/internal/service"
	"kondangin-backend/internal/utils"

	"encoding/base64"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) EncryptPassword(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input", gin.H{"error": err.Error()})
		return
	}

	// Parse public key dan private key
	pubKey, err := utils.ParsePublicKey()
	if err != nil {
		utils.SendInternalServerError(c, "Failed to parse public key", gin.H{"error": err.Error()})
		return
	}
	// Encrypt password
	originalPassword := req.Password
	encryptedPassword, err := utils.EncryptRSA(pubKey, originalPassword)
	if err != nil {
		utils.SendInternalServerError(c, "Error encrypting password", gin.H{"error": err.Error()})
		return
	}

	// Encode hasil enkripsi ke Base64
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(encryptedPassword))

	utils.SendSuccess(c, "Password encrypted successfully", encodedPassword)
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input", gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.RegisterUser(req)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to register user", gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, "User registered successfully", res)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input", gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.LoginUser(req)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid email or password", nil)
		return
	}

	utils.SendSuccess(c, "Login successful", res)
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.SendBadRequestError(c, "Token required", nil)
		return
	}

	if err := h.service.VerifyEmail(token); err != nil {
		utils.SendBadRequestError(c, err.Error(), nil)
		return
	}

	utils.SendSuccess(c, "Email successfully verified", nil)
}

func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResendVerificationEmail(req.Email); err != nil {
		utils.SendInternalServerError(c, "Failed to resend email verification", gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, "Verification email sent", nil)
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ForgotPassword(req.Email); err != nil {
		utils.SendInternalServerError(c, "Failed to send password reset email", gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, "Reset password email sent", nil)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.ValidateResetPasswordToken(req.Token)
	if err != nil {
		utils.SendBadRequestError(c, "Token tidak valid atau kadaluarsa", gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(claims.Subject, req.NewPassword); err != nil {
		utils.SendInternalServerError(c, "Gagal mereset password", gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, "Password berhasil direset", nil)
}
