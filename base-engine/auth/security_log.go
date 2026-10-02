/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package auth

import "log"

// LogAuthorizationDenied 记录不含敏感信息的总部授权拒绝事件。
func LogAuthorizationDenied(principal *WorkspacePrincipal, action, resourceType, resourceID string, err error) {
	if principal == nil || principal.WorkspaceType != WorkspaceTypeHeadquarters {
		return
	}
	code := ErrorCode(err)
	if code == "" {
		return
	}
	log.Printf(
		"security_event=authorization_denied request_id=%s actor_account_id=%s actor_membership_id=%s organization_id=%s action=%s resource_type=%s resource_id=%s result_code=%s",
		principal.RequestID, principal.AccountID, optionalID(principal.MembershipID), optionalID(principal.OrganizationID),
		action, resourceType, resourceID, code,
	)
}

func optionalID(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
