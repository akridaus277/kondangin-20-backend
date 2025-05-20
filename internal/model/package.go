package models

import "time"

type Package struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Code         string `gorm:"unique;not null" json:"code"`
	Name         string `gorm:"not null" json:"name"`
	DataJSON     string `gorm:"type:json" json:"data_json"`
	PropertyJSON string `gorm:"type:json" json:"property_json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    string
	UpdatedBy    string
	Active       bool `gorm:"default:true" json:"active"`
}
