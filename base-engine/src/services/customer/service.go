package customer

import (
	"base-engine/auth"
	"base-engine/src/services/audit"
	"base-engine/src/services/authentication"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"time"
)

// Service owns headquarters customer records and their business mutations.
type Service struct {
	db            *gorm.DB
	audit         *audit.Service
	now           func() time.Time
	lookupLimiter *authentication.LoginLimiter
}

func NewService(db *gorm.DB, auditService *audit.Service) *Service {
	return &Service{db: db, audit: auditService, now: time.Now, lookupLimiter: authentication.NewLoginLimiter(20, time.Minute)}
}

func headquartersID(principal *auth.WorkspacePrincipal, action string, mode authorization.AccessMode) (string, error) {
	if principal == nil {
		return "", auth.NewError(auth.CodeAuthRequired)
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || principal.OrganizationID == nil {
		return "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := authorization.Authorize(principal, authorization.Intent{
		Action: action, Mode: mode, ResourceOrganizationID: principal.OrganizationID,
	}); err != nil {
		return "", err
	}
	return *principal.OrganizationID, nil
}

func (s *Service) auditMember(tx *gorm.DB, principal *auth.WorkspacePrincipal, action, memberID, evidence string, reasonCodes ...string) error {
	metadata := audit.Metadata{EvidenceReference: evidence}
	if len(reasonCodes) > 0 {
		metadata.ReasonCode = reasonCodes[0]
	}
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: principal.OrganizationID, Action: action,
		ResourceType: "customerMember", ResourceID: memberID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, metadata),
	})
}
