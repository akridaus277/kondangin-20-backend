package service

import (
	model "kondangin-backend/internal/model"

	"gorm.io/gorm"
)

func HasPermission(db *gorm.DB, userID, invitationID uint, required string) bool {
	var perm model.InvitationPermission
	if err := db.
		Where("user_id = ? AND invitation_id = ?", userID, invitationID).
		First(&perm).Error; err != nil {
		return false
	}

	for _, p := range perm.Permissions {
		if p == required {
			return true
		}
	}
	return false
}
