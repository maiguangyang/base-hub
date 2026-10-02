/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"base-engine/src/services/organization"
)

// ProvisionFranchise 开通加盟组织与老板身份。
func (r *MutationResolver) ProvisionFranchise(ctx context.Context, input gen.ProvisionFranchiseInput) (*gen.ProvisionFranchisePayload, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	result, err := r.Services.Organizations.ProvisionFranchise(ctx, principal, organization.ProvisionInput{
		Code: input.Code, Name: input.Name, OwnerPhone: input.OwnerPhone,
		OwnerDisplayName: input.OwnerDisplayName, OwnerEmail: input.OwnerEmail,
	})
	if err != nil {
		return nil, err
	}
	return &gen.ProvisionFranchisePayload{
		Organization: result.Organization, Membership: result.Membership,
		TemporaryPassword: result.TemporaryPassword, InvitationPending: result.InvitationPending,
	}, nil
}

// SuspendOrganization 暂停组织并强制退出全部在线人员。
func (r *MutationResolver) SuspendOrganization(ctx context.Context, input gen.SuspendOrganizationInput) (*gen.Organization, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if err := r.Services.Organizations.Suspend(ctx, principal, input.OrganizationID, input.ReasonCode); err != nil {
		return nil, err
	}
	return loadOrganization(ctx, r.Services, input.OrganizationID)
}

// RestoreOrganization 恢复组织但不复活旧会话。
func (r *MutationResolver) RestoreOrganization(ctx context.Context, id string) (*gen.Organization, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if err := r.Services.Organizations.Restore(ctx, principal, id); err != nil {
		return nil, err
	}
	return loadOrganization(ctx, r.Services, id)
}

// ResetTemporaryPassword 轮换账号密码与全部会话。
func (r *MutationResolver) ResetTemporaryPassword(ctx context.Context, accountID string) (*gen.TemporaryPasswordPayload, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: "account:update", Mode: authorization.AccessUpdate}); err != nil {
		return nil, err
	}
	password, err := resetTemporaryPassword(ctx, r.Services, principal, accountID)
	if err != nil {
		return nil, err
	}
	return &gen.TemporaryPasswordPayload{AccountID: accountID, TemporaryPassword: password}, nil
}

func loadOrganization(ctx context.Context, services Dependencies, id string) (*gen.Organization, error) {
	organization := &gen.Organization{}
	return organization, services.DB.WithContext(ctx).First(organization, "id = ?", id).Error
}
