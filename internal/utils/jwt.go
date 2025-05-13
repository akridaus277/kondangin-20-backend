package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	models "kondangin-backend/internal/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// Load environment variables from .env file
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

// Fungsi untuk generate JWT token
func GenerateJWT(user models.User) (string, error) {
	// Ambil secret key dari environment variable
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	// Ambil durasi token dalam jam (bisa float)
	tokenDurationHours, _ := strconv.ParseFloat(os.Getenv("JWT_LOGIN_DURATION"), 64)

	// Ubah ke durasi dalam detik, lalu ubah jadi time.Duration
	tokenDuration := time.Duration(tokenDurationHours * float64(time.Hour))

	log.Printf("Generate token untuk user: ID=%d, Username=%s", user.ID, user.Username)

	// Buat klaim (claims) baru
	claims := &jwt.StandardClaims{
		Subject:   user.Username,                                    // ID user di-encode menjadi string
		ExpiresAt: time.Now().Add(tokenDuration * time.Hour).Unix(), // Token akan expired dalam 24 jam
		IssuedAt:  time.Now().Unix(),
	}

	// Membuat token dengan klaim dan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateVerificationToken(email string) (string, error) {
	// Ambil durasi token dalam jam (bisa float)
	tokenDurationHours, _ := strconv.ParseFloat(os.Getenv("JWT_REGISTER_VERIF_DURATION"), 64)

	// Ubah ke durasi dalam detik, lalu ubah jadi time.Duration
	tokenDuration := time.Duration(tokenDurationHours * float64(time.Hour))

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(tokenDuration * time.Hour).Unix(), // expired 1 hari
		"type":  "email_verification",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET") // ambil dari .env

	return token.SignedString([]byte(secret))
}

func ParseVerificationToken(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "email_verification" {
		return "", fmt.Errorf("invalid token type")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token payload")
	}

	return email, nil
}

func GenerateResetPasswordToken(email string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   email,
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateResetPasswordToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

// ValidateJWTToken memvalidasi JWT dan mengembalikan user username
func ValidateJWTToken(tokenString string) (*jwt.StandardClaims, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		log.Printf("Token subject: %s", claims.Subject)
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
