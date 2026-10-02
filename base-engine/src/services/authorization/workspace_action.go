/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package authorization

import (
	"base-engine/auth"
)

// WorkspaceAction maps shared GraphQL command permissions to the active workspace policy.
func WorkspaceAction(principal *auth.WorkspacePrincipal, declared string) (string, error) {
	if principal == nil {
		return "", auth.NewError(auth.CodeAuthRequired)
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		switch declared {
		case "operatorMembership:create":
			return "hqMembership:create", nil
		case "operatorMembership:update":
			return "hqMembership:update", nil
		}
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters && principal.WorkspaceType != auth.WorkspaceTypeFranchise {
		return "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return declared, nil
}
