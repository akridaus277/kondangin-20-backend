package service

import (
	"encoding/base64"
	"fmt"
	"kondangin-backend/internal/dto"
	model "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
	utils "kondangin-backend/internal/utils"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(req dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
	LoginUser(req dto.LoginUserRequest) (*dto.LoginUserResponse, error)
	VerifyEmail(token string) error
	ResendVerificationEmail(email string) error
	ForgotPassword(email string) error
	ResetPassword(email string, newPassword string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) RegisterUser(req dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	// Decrypt password dari frontend
	privateKey, err := utils.ParsePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to decode password: %w", err)
	}

	decryptedPassword, err := utils.DecryptRSA(privateKey, string(decodedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate email verification token (pakai JWT / UUID)
	verificationToken, err := utils.GenerateVerificationToken(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Simpan user ke database
	user := model.User{
		Email:             req.Email,
		Password:          string(hashedPassword),
		Name:              req.Name,
		IsActive:          false, // belum aktif
		VerificationToken: verificationToken,
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Kirim email verifikasi

	appUrl := os.Getenv("APP_URL")
	verificationLink := fmt.Sprintf(appUrl+"/verify?token=%s", verificationToken)

	if err := utils.SendVerificationEmail(user.Email, "Verifikasi Email Anda", verificationLink); err != nil {
		return nil, fmt.Errorf("failed to send verification email: %w", err)
	}

	// Response
	res := &dto.RegisterUserResponse{
		Email: user.Email,
		Name:  user.Name,
	}

	return res, nil
}

func (s *userService) LoginUser(req dto.LoginUserRequest) (*dto.LoginUserResponse, error) {

	// Decrypt password
	privateKey, err := utils.ParsePrivateKey()
	if err != nil {
		return nil, err
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(req.Password)
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := utils.DecryptRSA(privateKey, string(decodedPassword))
	if err != nil {
		return nil, err
	}

	// Cari user di database
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, err
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(decryptedPassword)); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return nil, err
	}

	res := &dto.LoginUserResponse{
		Token: token,
	}

	return res, nil
}

func (s *userService) VerifyEmail(token string) error {
	email, err := utils.ParseVerificationToken(token)
	if err != nil {
		return err
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	//Check Token is Used
	if user.VerificationToken == "" {
		// Jika user tidak ditemukan
		return fmt.Errorf("token not found")
	}

	if user.IsActive {
		return fmt.Errorf("account already activated")
	}

	user.IsActive = true
	user.VerificationToken = ""
	return s.repo.Update(user)
}

func (s *userService) ResendVerificationEmail(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.IsActive {
		return fmt.Errorf("user already verified")
	}

	// Generate new token
	token, err := utils.GenerateVerificationToken(email)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	user.VerificationToken = token

	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update verification token: %w", err)
	}

	// Kirim ulang email verifikasi
	return utils.SendVerificationEmail(user.Email, "Verifikasi Email Anda", token)
}

func (s *userService) ForgotPassword(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	resetToken, err := utils.GenerateResetPasswordToken(user.Email)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	user.ResetPasswordToken = resetToken // Pastikan field ini ada di model
	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	return utils.SendResetPasswordEmail(user.Email, "Reset Password", resetToken)
}

// ResetPassword mengubah password user berdasarkan userID
func (s *userService) ResetPassword(email string, newPassword string) error {
	// Verifikasi apakah user ada di database
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// Jika user tidak ditemukan
		return fmt.Errorf("user not found")
	}

	//Check Token is Used
	if user.ResetPasswordToken == "" {
		// Jika user tidak ditemukan
		return fmt.Errorf("token not found")
	}

	// Decrypt password dari frontend
	privateKey, err := utils.ParsePrivateKey()
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(newPassword)
	if err != nil {
		return fmt.Errorf("failed to decode password: %w", err)
	}

	decryptedPassword, err := utils.DecryptRSA(privateKey, string(decodedPassword))
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Hash password baru menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password di database
	user.ResetPasswordToken = ""
	user.Password = string(hashedPassword)
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Password berhasil diperbarui
	return nil
}
