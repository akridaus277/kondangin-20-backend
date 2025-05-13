package models

import (
	"time"

	"gorm.io/datatypes"
)

type InvitationPermission struct {
	ID           uint `gorm:"primaryKey"`
	InvitationID uint `gorm:"index" json:"invitation_id"`
	UserID       uint `gorm:"index" json:"user_id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    string
	UpdatedBy    string

	Permissions datatypes.JSONSlice[string] `gorm:"type:json" json:"permissions"`

	Invitation Invitation `gorm:"foreignKey:InvitationID;constraint:OnDelete:CASCADE"`
	User       User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
