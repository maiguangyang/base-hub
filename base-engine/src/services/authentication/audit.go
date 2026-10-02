/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authentication

import (
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
)

func (s *Service) writeAuthenticationAudit(tx *gorm.DB, account *gen.Account, sessionID *string, resultCode string) error {
	if s.audit == nil {
		return nil
	}
	accountID := ""
	if account != nil {
		accountID = account.ID
	}
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: accountID, SessionID: sessionID,
		Action: "authentication:login", ResourceType: "account", ResourceID: accountID,
		ResultCode: resultCode, Metadata: audit.Metadata{Source: "password"},
	})
}

func (s *Service) writePasswordChangeAudit(tx *gorm.DB, principalSessionID string, account *gen.Account) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: account.ID, SessionID: &principalSessionID,
		Action: "account:password_change", ResourceType: "account", ResourceID: account.ID,
		ResultCode: "SUCCESS", Metadata: audit.Metadata{Source: "self_service"},
	})
}
