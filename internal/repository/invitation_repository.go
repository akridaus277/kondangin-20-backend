package repository

import (
	"encoding/json"
	dao "kondangin-backend/internal/dao"
	model "kondangin-backend/internal/model"
	"log"

	"gorm.io/gorm"
)

type InvitationRepository interface {
	Create(invitation *model.Invitation) (*model.Invitation, error)
	FindByID(id uint) (*model.Invitation, error)
	FindBySubdomain(subdomain string) (*model.Invitation, error)
	Update(invitation *model.Invitation) (*model.Invitation, error)
	Delete(id uint) error
	FindAll() ([]model.Invitation, error)
	FindInvitationsByUserAccess(userID uint) ([]model.Invitation, error)
	FindByInvitationAndUserNative(userID uint) ([]dao.MemberDashboardInvitationDAO, error)
}

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db}
}

func (r *invitationRepository) Create(invitation *model.Invitation) (*model.Invitation, error) {
	if err := r.db.Create(invitation).Error; err != nil {
		return nil, err
	}
	return invitation, nil
}

func (r *invitationRepository) FindByID(id uint) (*model.Invitation, error) {
	var inv model.Invitation
	err := r.db.First(&inv, id).Error
	return &inv, err
}

func (r *invitationRepository) FindBySubdomain(subdomain string) (*model.Invitation, error) {
	var inv model.Invitation
	err := r.db.Where("subdomain = ?", subdomain).First(&inv).Error
	return &inv, err
}

func (r *invitationRepository) Update(invitation *model.Invitation) (*model.Invitation, error) {
	if err := r.db.Save(invitation).Error; err != nil {
		return nil, err
	}
	return invitation, nil
}

func (r *invitationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Invitation{}, id).Error
}

func (r *invitationRepository) FindAll() ([]model.Invitation, error) {
	var invitations []model.Invitation
	err := r.db.Find(&invitations).Error
	return invitations, err
}

func (r *invitationRepository) FindInvitationsByUserAccess(userID uint) ([]model.Invitation, error) {
	var invitations []model.Invitation

	err := r.db.
		Model(&model.Invitation{}).
		Joins("LEFT JOIN invitation_permissions ON invitations.id = invitation_permissions.invitation_id").
		Where("invitations.user_id = ? OR invitation_permissions.user_id = ?", userID, userID).
		Preload("InvitationPermissions", "user_id = ?", userID). // 🔥 filter hanya untuk user tersebut
		Preload("User").                                         // preload owner info
		Find(&invitations).Error

	return invitations, err
}

func (r *invitationRepository) FindByInvitationAndUserNative(userID uint) ([]dao.MemberDashboardInvitationDAO, error) {
	// Struct sementara untuk menampung JSON permissions sebagai string
	log.Println("di repo sebelum call db")
	type rawPermissionRow struct {
		InvitationTitle       string `gorm:"column:invitation_title"`
		InvitationStatus      string `gorm:"column:invitation_status"`
		InvitationSubdomain   string `gorm:"column:invitation_subdomain"`
		InvitationOwner       string `gorm:"column:invitation_owner"`
		InvitationPermissions string `gorm:"column:permissions"`
		InvitationCoverPhoto  string `gorm:"column:invitation_cover_photo"`
	}

	var rawRows []rawPermissionRow

	query := `
		SELECT 
		JSON_UNQUOTE(JSON_EXTRACT(i.data_json, '$.title')) AS invitation_title,
		CASE WHEN i.active IS TRUE THEN "active" ELSE "inactive" END AS invitation_status,
		i.subdomain AS invitation_subdomain,
		u.username AS invitation_owner,
		ip.permissions,
		JSON_UNQUOTE(JSON_EXTRACT(i.data_json, '$.coverPhoto')) AS invitation_cover_photo
		FROM users u
		JOIN invitation_permissions ip ON ip.user_id = u.id
		JOIN invitations i ON ip.invitation_id = i.id 
		WHERE u.id = ?;

	`

	err := r.db.Raw(query, userID).Scan(&rawRows).Error
	if err != nil {
		return nil, err
	}
	log.Printf("isi rawRows repo : %+v", rawRows)

	var result []dao.MemberDashboardInvitationDAO
	for _, row := range rawRows {
		log.Printf("isi row repo : %+v", row)
		var permissions []string
		if err := json.Unmarshal([]byte(row.InvitationPermissions), &permissions); err != nil {
			return nil, err
		}

		result = append(result, dao.MemberDashboardInvitationDAO{
			InvitationTitle:       row.InvitationTitle,
			InvitationStatus:      row.InvitationStatus,
			InvitationSubdomain:   row.InvitationSubdomain,
			InvitationOwner:       row.InvitationOwner,
			InvitationPermissions: permissions,
			InvitationCoverPhoto:  row.InvitationCoverPhoto,
		})
	}
	log.Printf("isi invitations repo : %+v", result)

	return result, nil
}
