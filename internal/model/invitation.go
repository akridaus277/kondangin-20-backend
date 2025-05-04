package models

type Invitation struct {
	ID           uint   `gorm:"primaryKey"`
	Subdomain    string `gorm:"unique;not null" json:"subdomain"`
	DataJSON     string `gorm:"type:json" json:"data_json"`
	PropertyJSON string `gorm:"type:json" json:"property_json"`
	UserID       uint   // Foreign key ke User

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
