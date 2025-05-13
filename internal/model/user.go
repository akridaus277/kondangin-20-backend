package models

import (
	"time"
)

type User struct {
	ID                 uint   `gorm:"primaryKey"`
	Username           string `gorm:"unique;not null" json:"username"`
	Email              string `gorm:"unique;not null"`
	Password           string `gorm:"not null"`
	Name               string
	Active             bool
	VerificationToken  string
	ResetPasswordToken string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CreatedBy          string
	UpdatedBy          string
	Invitations        []Invitation `gorm:"foreignKey:UserID"`
}
