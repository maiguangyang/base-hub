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

type CustomerMemberResultTypeResolver struct {
	*gen.GeneratedCustomerMemberResultTypeResolver
}

func (r *Resolver) CustomerMemberResultType() gen.CustomerMemberResultTypeResolver {
	return &CustomerMemberResultTypeResolver{&gen.GeneratedCustomerMemberResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerMemberResolver struct {
	*gen.GeneratedCustomerMemberResolver
}

func (r *Resolver) CustomerMember() gen.CustomerMemberResolver {
	return &CustomerMemberResolver{&gen.GeneratedCustomerMemberResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerBenefitPolicyResultTypeResolver struct {
	*gen.GeneratedCustomerBenefitPolicyResultTypeResolver
}

func (r *Resolver) CustomerBenefitPolicyResultType() gen.CustomerBenefitPolicyResultTypeResolver {
	return &CustomerBenefitPolicyResultTypeResolver{&gen.GeneratedCustomerBenefitPolicyResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerBenefitPolicyResolver struct {
	*gen.GeneratedCustomerBenefitPolicyResolver
}

func (r *Resolver) CustomerBenefitPolicy() gen.CustomerBenefitPolicyResolver {
	return &CustomerBenefitPolicyResolver{&gen.GeneratedCustomerBenefitPolicyResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerDailyPointGrantBudgetResultTypeResolver struct {
	*gen.GeneratedCustomerDailyPointGrantBudgetResultTypeResolver
}

func (r *Resolver) CustomerDailyPointGrantBudgetResultType() gen.CustomerDailyPointGrantBudgetResultTypeResolver {
	return &CustomerDailyPointGrantBudgetResultTypeResolver{&gen.GeneratedCustomerDailyPointGrantBudgetResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerDailyPointGrantBudgetResolver struct {
	*gen.GeneratedCustomerDailyPointGrantBudgetResolver
}

func (r *Resolver) CustomerDailyPointGrantBudget() gen.CustomerDailyPointGrantBudgetResolver {
	return &CustomerDailyPointGrantBudgetResolver{&gen.GeneratedCustomerDailyPointGrantBudgetResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerPointEntryResultTypeResolver struct {
	*gen.GeneratedCustomerPointEntryResultTypeResolver
}

func (r *Resolver) CustomerPointEntryResultType() gen.CustomerPointEntryResultTypeResolver {
	return &CustomerPointEntryResultTypeResolver{&gen.GeneratedCustomerPointEntryResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerPointEntryResolver struct {
	*gen.GeneratedCustomerPointEntryResolver
}

func (r *Resolver) CustomerPointEntry() gen.CustomerPointEntryResolver {
	return &CustomerPointEntryResolver{&gen.GeneratedCustomerPointEntryResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerCouponTemplateResultTypeResolver struct {
	*gen.GeneratedCustomerCouponTemplateResultTypeResolver
}

func (r *Resolver) CustomerCouponTemplateResultType() gen.CustomerCouponTemplateResultTypeResolver {
	return &CustomerCouponTemplateResultTypeResolver{&gen.GeneratedCustomerCouponTemplateResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerCouponTemplateResolver struct {
	*gen.GeneratedCustomerCouponTemplateResolver
}

func (r *Resolver) CustomerCouponTemplate() gen.CustomerCouponTemplateResolver {
	return &CustomerCouponTemplateResolver{&gen.GeneratedCustomerCouponTemplateResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductCategoryResultTypeResolver struct {
	*gen.GeneratedProductCategoryResultTypeResolver
}

func (r *Resolver) ProductCategoryResultType() gen.ProductCategoryResultTypeResolver {
	return &ProductCategoryResultTypeResolver{&gen.GeneratedProductCategoryResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductCategoryResolver struct {
	*gen.GeneratedProductCategoryResolver
}

func (r *Resolver) ProductCategory() gen.ProductCategoryResolver {
	return &ProductCategoryResolver{&gen.GeneratedProductCategoryResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductBrandResultTypeResolver struct {
	*gen.GeneratedProductBrandResultTypeResolver
}

func (r *Resolver) ProductBrandResultType() gen.ProductBrandResultTypeResolver {
	return &ProductBrandResultTypeResolver{&gen.GeneratedProductBrandResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductBrandResolver struct {
	*gen.GeneratedProductBrandResolver
}

func (r *Resolver) ProductBrand() gen.ProductBrandResolver {
	return &ProductBrandResolver{&gen.GeneratedProductBrandResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductResultTypeResolver struct {
	*gen.GeneratedProductResultTypeResolver
}

func (r *Resolver) ProductResultType() gen.ProductResultTypeResolver {
	return &ProductResultTypeResolver{&gen.GeneratedProductResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductResolver struct {
	*gen.GeneratedProductResolver
}

func (r *Resolver) Product() gen.ProductResolver {
	return &ProductResolver{&gen.GeneratedProductResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductSkuResultTypeResolver struct {
	*gen.GeneratedProductSkuResultTypeResolver
}

func (r *Resolver) ProductSkuResultType() gen.ProductSkuResultTypeResolver {
	return &ProductSkuResultTypeResolver{&gen.GeneratedProductSkuResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductSkuResolver struct {
	*gen.GeneratedProductSkuResolver
}

func (r *Resolver) ProductSku() gen.ProductSkuResolver {
	return &ProductSkuResolver{&gen.GeneratedProductSkuResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductPackageResultTypeResolver struct {
	*gen.GeneratedProductPackageResultTypeResolver
}

func (r *Resolver) ProductPackageResultType() gen.ProductPackageResultTypeResolver {
	return &ProductPackageResultTypeResolver{&gen.GeneratedProductPackageResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductPackageResolver struct {
	*gen.GeneratedProductPackageResolver
}

func (r *Resolver) ProductPackage() gen.ProductPackageResolver {
	return &ProductPackageResolver{&gen.GeneratedProductPackageResolver{GeneratedResolver: r.GeneratedResolver}}
}

type SpecificationDefinitionResultTypeResolver struct {
	*gen.GeneratedSpecificationDefinitionResultTypeResolver
}

func (r *Resolver) SpecificationDefinitionResultType() gen.SpecificationDefinitionResultTypeResolver {
	return &SpecificationDefinitionResultTypeResolver{&gen.GeneratedSpecificationDefinitionResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type SpecificationDefinitionResolver struct {
	*gen.GeneratedSpecificationDefinitionResolver
}

func (r *Resolver) SpecificationDefinition() gen.SpecificationDefinitionResolver {
	return &SpecificationDefinitionResolver{&gen.GeneratedSpecificationDefinitionResolver{GeneratedResolver: r.GeneratedResolver}}
}

type SpecificationValueResultTypeResolver struct {
	*gen.GeneratedSpecificationValueResultTypeResolver
}

func (r *Resolver) SpecificationValueResultType() gen.SpecificationValueResultTypeResolver {
	return &SpecificationValueResultTypeResolver{&gen.GeneratedSpecificationValueResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type SpecificationValueResolver struct {
	*gen.GeneratedSpecificationValueResolver
}

func (r *Resolver) SpecificationValue() gen.SpecificationValueResolver {
	return &SpecificationValueResolver{&gen.GeneratedSpecificationValueResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductSpecificationChoiceResultTypeResolver struct {
	*gen.GeneratedProductSpecificationChoiceResultTypeResolver
}

func (r *Resolver) ProductSpecificationChoiceResultType() gen.ProductSpecificationChoiceResultTypeResolver {
	return &ProductSpecificationChoiceResultTypeResolver{&gen.GeneratedProductSpecificationChoiceResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductSpecificationChoiceResolver struct {
	*gen.GeneratedProductSpecificationChoiceResolver
}

func (r *Resolver) ProductSpecificationChoice() gen.ProductSpecificationChoiceResolver {
	return &ProductSpecificationChoiceResolver{&gen.GeneratedProductSpecificationChoiceResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductSkuSpecificationValueResultTypeResolver struct {
	*gen.GeneratedProductSkuSpecificationValueResultTypeResolver
}

func (r *Resolver) ProductSkuSpecificationValueResultType() gen.ProductSkuSpecificationValueResultTypeResolver {
	return &ProductSkuSpecificationValueResultTypeResolver{&gen.GeneratedProductSkuSpecificationValueResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductSkuSpecificationValueResolver struct {
	*gen.GeneratedProductSkuSpecificationValueResolver
}

func (r *Resolver) ProductSkuSpecificationValue() gen.ProductSkuSpecificationValueResolver {
	return &ProductSkuSpecificationValueResolver{&gen.GeneratedProductSkuSpecificationValueResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductPackageTemplateResultTypeResolver struct {
	*gen.GeneratedProductPackageTemplateResultTypeResolver
}

func (r *Resolver) ProductPackageTemplateResultType() gen.ProductPackageTemplateResultTypeResolver {
	return &ProductPackageTemplateResultTypeResolver{&gen.GeneratedProductPackageTemplateResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type ProductPackageTemplateResolver struct {
	*gen.GeneratedProductPackageTemplateResolver
}

func (r *Resolver) ProductPackageTemplate() gen.ProductPackageTemplateResolver {
	return &ProductPackageTemplateResolver{&gen.GeneratedProductPackageTemplateResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreListingResultTypeResolver struct {
	*gen.GeneratedStoreListingResultTypeResolver
}

func (r *Resolver) StoreListingResultType() gen.StoreListingResultTypeResolver {
	return &StoreListingResultTypeResolver{&gen.GeneratedStoreListingResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreListingResolver struct {
	*gen.GeneratedStoreListingResolver
}

func (r *Resolver) StoreListing() gen.StoreListingResolver {
	return &StoreListingResolver{&gen.GeneratedStoreListingResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePackageOfferResultTypeResolver struct {
	*gen.GeneratedStorePackageOfferResultTypeResolver
}

func (r *Resolver) StorePackageOfferResultType() gen.StorePackageOfferResultTypeResolver {
	return &StorePackageOfferResultTypeResolver{&gen.GeneratedStorePackageOfferResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePackageOfferResolver struct {
	*gen.GeneratedStorePackageOfferResolver
}

func (r *Resolver) StorePackageOffer() gen.StorePackageOfferResolver {
	return &StorePackageOfferResolver{&gen.GeneratedStorePackageOfferResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePriceRevisionResultTypeResolver struct {
	*gen.GeneratedStorePriceRevisionResultTypeResolver
}

func (r *Resolver) StorePriceRevisionResultType() gen.StorePriceRevisionResultTypeResolver {
	return &StorePriceRevisionResultTypeResolver{&gen.GeneratedStorePriceRevisionResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePriceRevisionResolver struct {
	*gen.GeneratedStorePriceRevisionResolver
}

func (r *Resolver) StorePriceRevision() gen.StorePriceRevisionResolver {
	return &StorePriceRevisionResolver{&gen.GeneratedStorePriceRevisionResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreInventoryBatchResultTypeResolver struct {
	*gen.GeneratedStoreInventoryBatchResultTypeResolver
}

func (r *Resolver) StoreInventoryBatchResultType() gen.StoreInventoryBatchResultTypeResolver {
	return &StoreInventoryBatchResultTypeResolver{&gen.GeneratedStoreInventoryBatchResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreInventoryBatchResolver struct {
	*gen.GeneratedStoreInventoryBatchResolver
}

func (r *Resolver) StoreInventoryBatch() gen.StoreInventoryBatchResolver {
	return &StoreInventoryBatchResolver{&gen.GeneratedStoreInventoryBatchResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStockBalanceResultTypeResolver struct {
	*gen.GeneratedStoreStockBalanceResultTypeResolver
}

func (r *Resolver) StoreStockBalanceResultType() gen.StoreStockBalanceResultTypeResolver {
	return &StoreStockBalanceResultTypeResolver{&gen.GeneratedStoreStockBalanceResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStockBalanceResolver struct {
	*gen.GeneratedStoreStockBalanceResolver
}

func (r *Resolver) StoreStockBalance() gen.StoreStockBalanceResolver {
	return &StoreStockBalanceResolver{&gen.GeneratedStoreStockBalanceResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStocktakeResultTypeResolver struct {
	*gen.GeneratedStoreStocktakeResultTypeResolver
}

func (r *Resolver) StoreStocktakeResultType() gen.StoreStocktakeResultTypeResolver {
	return &StoreStocktakeResultTypeResolver{&gen.GeneratedStoreStocktakeResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStocktakeResolver struct {
	*gen.GeneratedStoreStocktakeResolver
}

func (r *Resolver) StoreStocktake() gen.StoreStocktakeResolver {
	return &StoreStocktakeResolver{&gen.GeneratedStoreStocktakeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStocktakeLineResultTypeResolver struct {
	*gen.GeneratedStoreStocktakeLineResultTypeResolver
}

func (r *Resolver) StoreStocktakeLineResultType() gen.StoreStocktakeLineResultTypeResolver {
	return &StoreStocktakeLineResultTypeResolver{&gen.GeneratedStoreStocktakeLineResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStocktakeLineResolver struct {
	*gen.GeneratedStoreStocktakeLineResolver
}

func (r *Resolver) StoreStocktakeLine() gen.StoreStocktakeLineResolver {
	return &StoreStocktakeLineResolver{&gen.GeneratedStoreStocktakeLineResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStockMovementResultTypeResolver struct {
	*gen.GeneratedStoreStockMovementResultTypeResolver
}

func (r *Resolver) StoreStockMovementResultType() gen.StoreStockMovementResultTypeResolver {
	return &StoreStockMovementResultTypeResolver{&gen.GeneratedStoreStockMovementResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StoreStockMovementResolver struct {
	*gen.GeneratedStoreStockMovementResolver
}

func (r *Resolver) StoreStockMovement() gen.StoreStockMovementResolver {
	return &StoreStockMovementResolver{&gen.GeneratedStoreStockMovementResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePromotionResultTypeResolver struct {
	*gen.GeneratedStorePromotionResultTypeResolver
}

func (r *Resolver) StorePromotionResultType() gen.StorePromotionResultTypeResolver {
	return &StorePromotionResultTypeResolver{&gen.GeneratedStorePromotionResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePromotionResolver struct {
	*gen.GeneratedStorePromotionResolver
}

func (r *Resolver) StorePromotion() gen.StorePromotionResolver {
	return &StorePromotionResolver{&gen.GeneratedStorePromotionResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePromotionTargetResultTypeResolver struct {
	*gen.GeneratedStorePromotionTargetResultTypeResolver
}

func (r *Resolver) StorePromotionTargetResultType() gen.StorePromotionTargetResultTypeResolver {
	return &StorePromotionTargetResultTypeResolver{&gen.GeneratedStorePromotionTargetResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type StorePromotionTargetResolver struct {
	*gen.GeneratedStorePromotionTargetResolver
}

func (r *Resolver) StorePromotionTarget() gen.StorePromotionTargetResolver {
	return &StorePromotionTargetResolver{&gen.GeneratedStorePromotionTargetResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerCouponGrantResultTypeResolver struct {
	*gen.GeneratedCustomerCouponGrantResultTypeResolver
}

func (r *Resolver) CustomerCouponGrantResultType() gen.CustomerCouponGrantResultTypeResolver {
	return &CustomerCouponGrantResultTypeResolver{&gen.GeneratedCustomerCouponGrantResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerCouponGrantResolver struct {
	*gen.GeneratedCustomerCouponGrantResolver
}

func (r *Resolver) CustomerCouponGrant() gen.CustomerCouponGrantResolver {
	return &CustomerCouponGrantResolver{&gen.GeneratedCustomerCouponGrantResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerCouponDistributionJobResultTypeResolver struct {
	*gen.GeneratedCustomerCouponDistributionJobResultTypeResolver
}

func (r *Resolver) CustomerCouponDistributionJobResultType() gen.CustomerCouponDistributionJobResultTypeResolver {
	return &CustomerCouponDistributionJobResultTypeResolver{&gen.GeneratedCustomerCouponDistributionJobResultTypeResolver{GeneratedResolver: r.GeneratedResolver}}
}

type CustomerCouponDistributionJobResolver struct {
	*gen.GeneratedCustomerCouponDistributionJobResolver
}

func (r *Resolver) CustomerCouponDistributionJob() gen.CustomerCouponDistributionJobResolver {
	return &CustomerCouponDistributionJobResolver{&gen.GeneratedCustomerCouponDistributionJobResolver{GeneratedResolver: r.GeneratedResolver}}
}
