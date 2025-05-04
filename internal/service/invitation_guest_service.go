package service

import (
	model "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
)

type InvitationGuestService interface {
	GetInvitationBySubdomain(subdomain string) (*model.Invitation, error)
}

type invitationGuestService struct {
	repo repository.InvitationRepository
}

func NewInvitationGuestService(repo repository.InvitationRepository) InvitationGuestService {
	return &invitationGuestService{repo}
}

func (s *invitationGuestService) GetInvitationBySubdomain(subdomain string) (*model.Invitation, error) {
	return s.repo.FindBySubdomain(subdomain)
}
