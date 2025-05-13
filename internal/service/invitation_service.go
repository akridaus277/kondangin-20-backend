package service

import (
	dao "kondangin-backend/internal/dao"
	model "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
)

type InvitationService interface {
	CreateInvitation(c *gin.Context, input *model.Invitation) (*model.Invitation, error)
	GetInvitationByID(c *gin.Context, id uint) (*model.Invitation, error)
	GetInvitationBySubdomain(c *gin.Context, subdomain string) (*model.Invitation, error)
	UpdateInvitation(c *gin.Context, inv *model.Invitation) (*model.Invitation, error)
	DeleteInvitation(c *gin.Context, id uint) error
	ListInvitations(c *gin.Context) ([]model.Invitation, error)
	GetInvitationsForUser(userID uint) ([]dao.MemberDashboardInvitationDAO, error)
}

type invitationService struct {
	repo repository.InvitationRepository
}

func NewInvitationService(repo repository.InvitationRepository) InvitationService {
	return &invitationService{repo}
}

func (s *invitationService) CreateInvitation(c *gin.Context, input *model.Invitation) (*model.Invitation, error) {
	return s.repo.Create(input)
}

func (s *invitationService) GetInvitationByID(c *gin.Context, id uint) (*model.Invitation, error) {
	return s.repo.FindByID(id)
}

func (s *invitationService) GetInvitationBySubdomain(c *gin.Context, subdomain string) (*model.Invitation, error) {
	return s.repo.FindBySubdomain(subdomain)
}

func (s *invitationService) UpdateInvitation(c *gin.Context, inv *model.Invitation) (*model.Invitation, error) {
	return s.repo.Update(inv)
}

func (s *invitationService) DeleteInvitation(c *gin.Context, id uint) error {
	return s.repo.Delete(id)
}

func (s *invitationService) ListInvitations(c *gin.Context) ([]model.Invitation, error) {
	return s.repo.FindAll()
}

func (s *invitationService) GetInvitationsForUser(userID uint) ([]dao.MemberDashboardInvitationDAO, error) {
	log.Println("di invi service sebelum call repo.FindByInvitationAndUserNative")
	return s.repo.FindByInvitationAndUserNative(userID)
}
