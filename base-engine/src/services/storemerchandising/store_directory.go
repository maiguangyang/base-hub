package storemerchandising

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

// ScopedStores lists only stores the operator can use for merchandising work.
func (s *Service) ScopedStores(ctx context.Context, principal *auth.WorkspacePrincipal, page, perPage int) (*Page[*gen.Store], error) {
	orgID, err := directoryOrganization(principal)
	if err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.Store{}).Where("organization_id = ? AND lifecycle = ? AND (is_delete IS NULL OR is_delete = ?)", orgID, gen.StoreLifecycleActive, 1)
	if !principal.AllStores {
		ids := directoryStoreIDs(principal)
		if len(ids) == 0 {
			return &Page[*gen.Store]{Page: page, PerPage: perPage, Data: []*gen.Store{}}, nil
		}
		query = query.Where("id IN ?", ids)
	}
	result := &Page[*gen.Store]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("name, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func directoryOrganization(principal *auth.WorkspacePrincipal) (string, error) {
	if principal == nil {
		return "", auth.NewError(auth.CodeAuthRequired)
	}
	if principal.WorkspaceType != auth.WorkspaceTypeFranchise || principal.OrganizationID == nil {
		return "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	action := scopedReadAction(principal, "franchiseProduct:read", "franchiseStock:read", "franchiseStocktake:read", "franchiseStocktake:record", "franchiseStocktake:post", "franchisePromotion:read", "franchiseCoupon:read")
	intent := authorization.Intent{Action: action, Mode: authorization.AccessRead, ResourceOrganizationID: principal.OrganizationID}
	if err := authorization.Authorize(principal, intent); err != nil {
		return "", err
	}
	return *principal.OrganizationID, nil
}

func directoryStoreIDs(principal *auth.WorkspacePrincipal) []string {
	ids := make([]string, 0, len(principal.StoreIDs))
	for id := range principal.StoreIDs {
		ids = append(ids, id)
	}
	return ids
}
