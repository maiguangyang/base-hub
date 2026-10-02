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
	"base-engine/src/services/membership"
	storeservice "base-engine/src/services/store"
)

// InviteOperator 邀请当前加盟组织的操作员。
func (r *MutationResolver) InviteOperator(ctx context.Context, input gen.InviteOperatorInput) (*gen.InviteOperatorPayload, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	result, err := r.Services.Memberships.InviteOperator(ctx, principal, membership.InviteInput{
		Phone: input.Phone, DisplayName: input.DisplayName, Email: input.Email,
		RoleIDs: input.RoleIds, StoreAccessMode: input.StoreAccessMode, StoreIDs: input.StoreIds,
	})
	if err != nil {
		return nil, err
	}
	return &gen.InviteOperatorPayload{Membership: result.Membership, TemporaryPassword: result.TemporaryPassword, InvitationPending: result.InvitationPending}, nil
}

// AcceptMembershipInvitation 接受属于当前账号的有效邀请。
func (r *MutationResolver) AcceptMembershipInvitation(ctx context.Context, id string) (*gen.OperatorMembership, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return r.Services.Memberships.AcceptInvitation(ctx, principal, id)
}

// ChangeMembershipStatus 修改当前加盟组织成员状态并保护最后一个老板。
func (r *MutationResolver) ChangeMembershipStatus(ctx context.Context, input gen.ChangeMembershipStatusInput) (*gen.OperatorMembership, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if err := r.Services.Memberships.ChangeMembershipStatus(ctx, principal, input.MembershipID, input.Status); err != nil {
		return nil, err
	}
	membership := &gen.OperatorMembership{}
	if err := r.Services.DB.WithContext(ctx).First(membership, "id = ?", input.MembershipID).Error; err != nil {
		return nil, err
	}
	return membership, nil
}

// SubmitStore 提交当前组织门店审核。
func (r *MutationResolver) SubmitStore(ctx context.Context, id string) (*gen.Store, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return r.Services.Stores.Submit(ctx, principal, id)
}

// ReviewStore 由总部使用独立批准或退回权限审核门店。
func (r *MutationResolver) ReviewStore(ctx context.Context, input gen.ReviewStoreInput) (*gen.Store, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	reason := ""
	if input.RejectionReason != nil {
		reason = *input.RejectionReason
	}
	return r.Services.Stores.Review(ctx, principal, storeservice.ReviewInput{StoreID: input.StoreID, Approved: input.Approved, RejectionReason: reason})
}

func requireResolverPrincipal(ctx context.Context) (*auth.WorkspacePrincipal, error) {
	return auth.RequirePrincipal(ctx)
}
