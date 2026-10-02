package ai

import (
	"base-engine/auth"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
)

func (s *ModelConfigStore) auditModelConfig(tx *gorm.DB, principal *auth.WorkspacePrincipal, action string, version uint64, status, result string) error {
	if s.audit == nil {
		return ErrModelConfigUnavailable
	}
	var sessionID *string
	if principal.SessionID != "" {
		sessionID = &principal.SessionID
	}
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: sessionID,
		Action: action, ResourceType: "aiModelConfig", ResultCode: result,
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{
			Source: "headquarters", TargetStatus: status, ModelConfigVersion: version,
		}),
	})
}
