package productcatalog

import (
	"base-engine/auth"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
)

func (s *Service) templateAudit(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, resourceType, resourceID string) error {
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: resourceType, ResourceID: resourceID, ResultCode: "SUCCESS"})
}
