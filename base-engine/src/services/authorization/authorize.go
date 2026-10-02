/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"strings"

	"base-engine/auth"
	"base-engine/gen"
)

// AccessMode 描述一次资源访问的语义。
type AccessMode string

const (
	AccessRead     AccessMode = "READ"
	AccessCreate   AccessMode = "CREATE"
	AccessUpdate   AccessMode = "UPDATE"
	AccessDelete   AccessMode = "DELETE"
	AccessRecover  AccessMode = "RECOVER"
	AccessRelation AccessMode = "RELATION"
)

// Intent 是中央授权器接收的完整资源意图。
type Intent struct {
	Action                 string
	Mode                   AccessMode
	ResourceOrganizationID *string
	StoreID                *string
}

// Authorize 按工作台、权限、组织与门店范围执行固定顺序校验。
func Authorize(principal *auth.WorkspacePrincipal, intent Intent) error {
	if principal == nil {
		return auth.NewError(auth.CodeAuthRequired)
	}
	if principal.WorkspaceType == auth.WorkspaceTypeDiscovery {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if !principal.Has(intent.Action) {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return authorizeHeadquarters(principal, intent)
	}
	if principal.WorkspaceType != auth.WorkspaceTypeFranchise {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return authorizeFranchise(principal, intent)
}

func authorizeHeadquarters(principal *auth.WorkspacePrincipal, intent Intent) error {
	if intent.Mode == AccessRead || intent.Mode == AccessRelation {
		return nil
	}
	if isGovernanceAction(intent.Action) {
		return nil
	}
	if strings.HasPrefix(intent.Action, "hq") {
		if sameOrganization(principal.OrganizationID, intent.ResourceOrganizationID) {
			return nil
		}
		return auth.NewError(auth.CodePermissionDenied)
	}
	return auth.NewError(auth.CodeWorkspaceForbidden)
}

func sameOrganization(left, right *string) bool {
	return left != nil && right != nil && *left == *right
}

func authorizeFranchise(principal *auth.WorkspacePrincipal, intent Intent) error {
	if intent.ResourceOrganizationID != nil {
		if principal.OrganizationID == nil || *principal.OrganizationID != *intent.ResourceOrganizationID {
			return auth.NewError(auth.CodePermissionDenied)
		}
	}
	if intent.StoreID != nil && !principal.HasStore(*intent.StoreID) {
		return auth.NewError(auth.CodeStoreScopeDenied)
	}
	return nil
}

func isGovernanceAction(action string) bool {
	switch action {
	case "franchise:provision", "organization:suspend", "organization:restore", "store:approve", "store:reject", "account:update":
		return true
	default:
		return false
	}
}

func authorizeDeletionState(isDelete *int64, mode AccessMode) error {
	active := isDelete == nil || *isDelete == 1
	if mode == AccessRecover && !active {
		return nil
	}
	if mode != AccessRecover && active {
		return nil
	}
	return auth.NewError(auth.CodePermissionDenied)
}

// AuthorizeOperationalStore 拒绝未启用门店承载运营写入。
func AuthorizeOperationalStore(store *gen.Store) error {
	if store == nil || store.Lifecycle != gen.StoreLifecycleActive {
		return auth.NewError(auth.CodeStoreNotActive)
	}
	return nil
}
