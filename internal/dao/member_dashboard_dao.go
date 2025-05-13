package dto

type MemberDashboardInvitationDAO struct {
	InvitationTitle       string
	InvitationStatus      string
	InvitationSubdomain   string
	InvitationOwner       string
	InvitationPermissions []string
	InvitationCoverPhoto  string
}

type LovEventTypeDAO struct {
	Code string
	Name string
}
