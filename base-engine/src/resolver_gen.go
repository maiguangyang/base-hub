package src

import (
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

func NewResolver(db *gen.DB, ec *gen.EventController, services Dependencies) *Resolver {
	defaults := gen.DefaultResolutionHandlers()
	handlers := authorization.RegisterHandlers(defaults, authorization.HandlerDependencies{Sessions: services.Sessions})
	if err := authorization.ValidateHandlerCoverage(handlers, defaults); err != nil {
		panic(err)
	}
	return &Resolver{GeneratedResolver: &gen.GeneratedResolver{Handlers: handlers, DB: db, EventController: ec}, Services: services}
}

type Resolver struct {
	*gen.GeneratedResolver
	Services Dependencies
}

type MutationResolver struct {
	*gen.GeneratedMutationResolver
	Services Dependencies
}

type QueryResolver struct {
	*gen.GeneratedQueryResolver
	Services Dependencies
}

type SubscriptionResolver struct {
	*gen.GeneratedSubscriptionResolver
	Services Dependencies
}

func (r *Resolver) Mutation() gen.MutationResolver {
	return &MutationResolver{GeneratedMutationResolver: &gen.GeneratedMutationResolver{GeneratedResolver: r.GeneratedResolver}, Services: r.Services}
}
func (r *Resolver) Query() gen.QueryResolver {
	return &QueryResolver{GeneratedQueryResolver: &gen.GeneratedQueryResolver{GeneratedResolver: r.GeneratedResolver}, Services: r.Services}
}

func (r *Resolver) Subscription() gen.SubscriptionResolver {
	return &SubscriptionResolver{GeneratedSubscriptionResolver: &gen.GeneratedSubscriptionResolver{GeneratedResolver: r.GeneratedResolver}, Services: r.Services}
}

type AccountResultTypeResolver struct {
	*gen.GeneratedAccountResultTypeResolver
}

func (r *Resolver) AccountResultType() gen.AccountResultTypeResolver {
	return &AccountResultTypeResolver{&gen.GeneratedAccountResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type AccountResolver struct {
	*gen.GeneratedAccountResolver
}

func (r *Resolver) Account() gen.AccountResolver {
	return &AccountResolver{&gen.GeneratedAccountResolver{GeneratedResolver: r.GeneratedResolver}}
}

type OrganizationResultTypeResolver struct {
	*gen.GeneratedOrganizationResultTypeResolver
}

func (r *Resolver) OrganizationResultType() gen.OrganizationResultTypeResolver {
	return &OrganizationResultTypeResolver{&gen.GeneratedOrganizationResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type OrganizationResolver struct {
	*gen.GeneratedOrganizationResolver
}

func (r *Resolver) Organization() gen.OrganizationResolver {
	return &OrganizationResolver{&gen.GeneratedOrganizationResolver{GeneratedResolver: r.GeneratedResolver}}
}

type OperatorMembershipResultTypeResolver struct {
	*gen.GeneratedOperatorMembershipResultTypeResolver
}

func (r *Resolver) OperatorMembershipResultType() gen.OperatorMembershipResultTypeResolver {
	return &OperatorMembershipResultTypeResolver{&gen.GeneratedOperatorMembershipResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type OperatorMembershipResolver struct {
	*gen.GeneratedOperatorMembershipResolver
}

func (r *Resolver) OperatorMembership() gen.OperatorMembershipResolver {
	return &OperatorMembershipResolver{&gen.GeneratedOperatorMembershipResolver{GeneratedResolver: r.GeneratedResolver}}
}

type PermissionResultTypeResolver struct {
	*gen.GeneratedPermissionResultTypeResolver
}

func (r *Resolver) PermissionResultType() gen.PermissionResultTypeResolver {
	return &PermissionResultTypeResolver{&gen.GeneratedPermissionResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type PermissionResolver struct {
	*gen.GeneratedPermissionResolver
}

func (r *Resolver) Permission() gen.PermissionResolver {
	return &PermissionResolver{&gen.GeneratedPermissionResolver{GeneratedResolver: r.GeneratedResolver}}
}

type OperatorRoleResultTypeResolver struct {
	*gen.GeneratedOperatorRoleResultTypeResolver
}

func (r *Resolver) OperatorRoleResultType() gen.OperatorRoleResultTypeResolver {
	return &OperatorRoleResultTypeResolver{&gen.GeneratedOperatorRoleResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type OperatorRoleResolver struct {
	*gen.GeneratedOperatorRoleResolver
}

func (r *Resolver) OperatorRole() gen.OperatorRoleResolver {
	return &OperatorRoleResolver{&gen.GeneratedOperatorRoleResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreResultTypeResolver struct {
	*gen.GeneratedStoreResultTypeResolver
}

func (r *Resolver) StoreResultType() gen.StoreResultTypeResolver {
	return &StoreResultTypeResolver{&gen.GeneratedStoreResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreResolver struct {
	*gen.GeneratedStoreResolver
}

func (r *Resolver) Store() gen.StoreResolver {
	return &StoreResolver{&gen.GeneratedStoreResolver{GeneratedResolver: r.GeneratedResolver}}
}

type SessionResultTypeResolver struct {
	*gen.GeneratedSessionResultTypeResolver
}

func (r *Resolver) SessionResultType() gen.SessionResultTypeResolver {
	return &SessionResultTypeResolver{&gen.GeneratedSessionResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type SessionResolver struct {
	*gen.GeneratedSessionResolver
}

func (r *Resolver) Session() gen.SessionResolver {
	return &SessionResolver{&gen.GeneratedSessionResolver{GeneratedResolver: r.GeneratedResolver}}
}

type MembershipInvitationResultTypeResolver struct {
	*gen.GeneratedMembershipInvitationResultTypeResolver
}

func (r *Resolver) MembershipInvitationResultType() gen.MembershipInvitationResultTypeResolver {
	return &MembershipInvitationResultTypeResolver{&gen.GeneratedMembershipInvitationResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type MembershipInvitationResolver struct {
	*gen.GeneratedMembershipInvitationResolver
}

func (r *Resolver) MembershipInvitation() gen.MembershipInvitationResolver {
	return &MembershipInvitationResolver{&gen.GeneratedMembershipInvitationResolver{GeneratedResolver: r.GeneratedResolver}}
}

type AuditLogResultTypeResolver struct {
	*gen.GeneratedAuditLogResultTypeResolver
}

func (r *Resolver) AuditLogResultType() gen.AuditLogResultTypeResolver {
	return &AuditLogResultTypeResolver{&gen.GeneratedAuditLogResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type AuditLogResolver struct {
	*gen.GeneratedAuditLogResolver
}

func (r *Resolver) AuditLog() gen.AuditLogResolver {
	return &AuditLogResolver{&gen.GeneratedAuditLogResolver{GeneratedResolver: r.GeneratedResolver}}
}

type FranchiseOpeningRecordResultTypeResolver struct {
	*gen.GeneratedFranchiseOpeningRecordResultTypeResolver
}

func (r *Resolver) FranchiseOpeningRecordResultType() gen.FranchiseOpeningRecordResultTypeResolver {
	return &FranchiseOpeningRecordResultTypeResolver{&gen.GeneratedFranchiseOpeningRecordResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type FranchiseOpeningRecordResolver struct {
	*gen.GeneratedFranchiseOpeningRecordResolver
}

func (r *Resolver) FranchiseOpeningRecord() gen.FranchiseOpeningRecordResolver {
	return &FranchiseOpeningRecordResolver{&gen.GeneratedFranchiseOpeningRecordResolver{GeneratedResolver: r.GeneratedResolver}}
}

type GlobalPaymentConfigResultTypeResolver struct {
	*gen.GeneratedGlobalPaymentConfigResultTypeResolver
}

func (r *Resolver) GlobalPaymentConfigResultType() gen.GlobalPaymentConfigResultTypeResolver {
	return &GlobalPaymentConfigResultTypeResolver{&gen.GeneratedGlobalPaymentConfigResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type FranchisePaymentConfigResultTypeResolver struct {
	*gen.GeneratedFranchisePaymentConfigResultTypeResolver
}

func (r *Resolver) FranchisePaymentConfigResultType() gen.FranchisePaymentConfigResultTypeResolver {
	return &FranchisePaymentConfigResultTypeResolver{&gen.GeneratedFranchisePaymentConfigResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type FranchisePaymentConfigResolver struct {
	*gen.GeneratedFranchisePaymentConfigResolver
}

func (r *Resolver) FranchisePaymentConfig() gen.FranchisePaymentConfigResolver {
	return &FranchisePaymentConfigResolver{&gen.GeneratedFranchisePaymentConfigResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePaymentConfigResultTypeResolver struct {
	*gen.GeneratedStorePaymentConfigResultTypeResolver
}

func (r *Resolver) StorePaymentConfigResultType() gen.StorePaymentConfigResultTypeResolver {
	return &StorePaymentConfigResultTypeResolver{&gen.GeneratedStorePaymentConfigResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePaymentConfigResolver struct {
	*gen.GeneratedStorePaymentConfigResolver
}

func (r *Resolver) StorePaymentConfig() gen.StorePaymentConfigResolver {
	return &StorePaymentConfigResolver{&gen.GeneratedStorePaymentConfigResolver{GeneratedResolver: r.GeneratedResolver}}
}
