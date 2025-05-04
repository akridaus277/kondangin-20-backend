package models

import (
	"time"
)

type User struct {
	ID                 uint   `gorm:"primaryKey"`
	Email              string `gorm:"unique"`
	Password           string
	Name               string
	IsActive           bool
	VerificationToken  string
	ResetPasswordToken string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Invitations        []Invitation `gorm:"foreignKey:UserID"`
}
