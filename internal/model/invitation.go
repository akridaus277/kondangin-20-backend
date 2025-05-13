package models

import "time"

type Invitation struct {
	ID           uint   `gorm:"primaryKey"`
	Subdomain    string `gorm:"unique;not null" json:"subdomain"`
	DataJSON     string `gorm:"type:json" json:"data_json"`
	PropertyJSON string `gorm:"type:json" json:"property_json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    string
	UpdatedBy    string
	Active       bool
	UserID       uint // Foreign key ke User

	User                  User                   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	InvitationPermissions []InvitationPermission `gorm:"foreignKey:InvitationID"`
}
