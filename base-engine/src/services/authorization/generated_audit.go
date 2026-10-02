/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
)

// writeGeneratedAudit 在 generated mutation 的同一事务内写入结构化审计。
func writeGeneratedAudit(ctx context.Context, resolver *gen.GeneratedResolver, principal *auth.WorkspacePrincipal, organizationID, action, resourceType, resourceID string) error {
	return writeGeneratedAuditMetadata(ctx, resolver, principal, organizationID, action, resourceType, resourceID, audit.Metadata{})
}

func writeGeneratedAuditMetadata(ctx context.Context, resolver *gen.GeneratedResolver, principal *auth.WorkspacePrincipal, organizationID, action, resourceType, resourceID string, metadata audit.Metadata) error {
	return audit.NewService().Write(resolverDB(ctx, resolver), audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &organizationID, Action: action,
		ResourceType: resourceType, ResourceID: resourceID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, metadata),
	})
}
