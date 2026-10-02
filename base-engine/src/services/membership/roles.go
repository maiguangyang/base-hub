/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

import (
	"context"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

// SaveRoleInput 描述自定义角色的名称与 Allow 权限集合。
type SaveRoleInput struct {
	ID            *string
	Name          string
	PermissionIDs []string
}

// RoleService 管理自定义角色及权限 scope 不变量。
type RoleService struct {
	db    *gorm.DB
	audit *audit.Service
}

// NewRoleService 创建角色服务。
func NewRoleService(db *gorm.DB, auditService *audit.Service) *RoleService {
	return &RoleService{db: db, audit: auditService}
}

// SaveCustomRole 创建或更新当前组织的 CUSTOM 角色。
func (s *RoleService) SaveCustomRole(ctx context.Context, principal *auth.WorkspacePrincipal, input SaveRoleInput) (*gen.OperatorRole, error) {
	organizationID, action, scope, err := roleWorkspace(principal, input.ID == nil)
	if err != nil {
		return nil, err
	}
	mode := authorization.AccessUpdate
	if input.ID == nil {
		mode = authorization.AccessCreate
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: action, Mode: mode, ResourceOrganizationID: &organizationID}); err != nil {
		return nil, err
	}
	var result *gen.OperatorRole
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		role, txErr := s.saveRoleTransaction(tx, input, organizationID, scope)
		result = role
		return txErr
	})
	return result, err
}

func (s *RoleService) saveRoleTransaction(tx *gorm.DB, input SaveRoleInput, organizationID string, scope gen.PermissionScope) (*gen.OperatorRole, error) {
	permissions, err := loadScopedPermissions(tx, input.PermissionIDs, scope)
	if err != nil {
		return nil, err
	}
	role := &gen.OperatorRole{}
	if input.ID == nil {
		role = &gen.OperatorRole{ID: uuid.Must(uuid.NewV4()).String(), Name: input.Name, Kind: gen.RoleKindCustom, OrganizationID: organizationID}
		if err := tx.Create(role).Error; err != nil {
			return nil, err
		}
	} else {
		if err := tx.Where("is_delete IS NULL OR is_delete = ?", 1).
			First(role, "id = ? AND organization_id = ?", *input.ID, organizationID).Error; err != nil {
			return nil, auth.NewError(auth.CodePermissionDenied)
		}
		if role.Kind != gen.RoleKindCustom {
			return nil, auth.NewError(auth.CodePermissionDenied)
		}
		if err := tx.Model(role).Update("name", input.Name).Error; err != nil {
			return nil, err
		}
	}
	if err := tx.Model(role).Association("Permissions").Replace(permissions); err != nil {
		return nil, err
	}
	return role, nil
}

func roleWorkspace(principal *auth.WorkspacePrincipal, creating bool) (string, string, gen.PermissionScope, error) {
	if principal == nil || principal.OrganizationID == nil {
		return "", "", "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	verb := "update"
	if creating {
		verb = "create"
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return *principal.OrganizationID, "hqRole:" + verb, gen.PermissionScopeSystem, nil
	}
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		return *principal.OrganizationID, "operatorRole:" + verb, gen.PermissionScopeTenant, nil
	}
	return "", "", "", auth.NewError(auth.CodeWorkspaceForbidden)
}

func loadScopedPermissions(tx *gorm.DB, ids []string, scope gen.PermissionScope) ([]*gen.Permission, error) {
	permissions := make([]*gen.Permission, 0, len(ids))
	if len(ids) == 0 {
		return permissions, nil
	}
	if err := tx.Where("id IN ? AND scope = ?", ids, scope).
		Where("is_delete IS NULL OR is_delete = ?", 1).Find(&permissions).Error; err != nil {
		return nil, err
	}
	if len(permissions) != len(ids) {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return permissions, nil
}

// ResolvePermissionUnion 返回多角色 Allow 权限并集。
func ResolvePermissionUnion(roles []*gen.OperatorRole) map[string]struct{} {
	result := make(map[string]struct{})
	for _, role := range roles {
		for _, permission := range role.Permissions {
			result[permission.Action] = struct{}{}
		}
	}
	return result
}
