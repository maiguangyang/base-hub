//go:generate go run github.com/99designs/gqlgen generate
package gen

import (
	"context"
)

type ResolutionHandlers struct {
	OnEvent   func(ctx context.Context, r *GeneratedResolver, e *Event) error
	WebSocket func(ctx context.Context, r *GeneratedResolver) (<-chan any, error)

	CreateAccount    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Account, err error)
	UpdateAccount    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Account, err error)
	DeleteAccounts   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryAccounts func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryAccount     func(ctx context.Context, r *GeneratedResolver, opts QueryAccountHandlerOptions) (*Account, error)
	QueryAccounts    func(ctx context.Context, r *GeneratedResolver, opts QueryAccountsHandlerOptions) (*AccountResultType, error)

	AccountMemberships func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*OperatorMembership, err error)

	AccountInitializedOrganizations func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*Organization, err error)

	AccountOpeningRecords func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*FranchiseOpeningRecord, err error)

	AccountRecordedOpeningRecords func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*FranchiseOpeningRecord, err error)

	AccountSessions func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*Session, err error)

	AccountReviewedStores func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*Store, err error)

	AccountSentMembershipInvitations func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*MembershipInvitation, err error)

	AccountAuditLogs func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*AuditLog, err error)

	AccountCreatedStocktakes func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*StoreStocktake, err error)

	AccountPostedStocktakes func(ctx context.Context, r *GeneratedResolver, obj *Account) (res []*StoreStocktake, err error)

	CreateOrganization    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Organization, err error)
	UpdateOrganization    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Organization, err error)
	DeleteOrganizations   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryOrganizations func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryOrganization     func(ctx context.Context, r *GeneratedResolver, opts QueryOrganizationHandlerOptions) (*Organization, error)
	QueryOrganizations    func(ctx context.Context, r *GeneratedResolver, opts QueryOrganizationsHandlerOptions) (*OrganizationResultType, error)

	OrganizationMemberships func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*OperatorMembership, err error)

	OrganizationInitialAccount func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res *Account, err error)

	OrganizationOpeningRecords func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*FranchiseOpeningRecord, err error)

	OrganizationStores func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*Store, err error)

	OrganizationRoles func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*OperatorRole, err error)

	OrganizationSessions func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*Session, err error)

	OrganizationAuditLogs func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*AuditLog, err error)

	OrganizationPaymentConfigs func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*FranchisePaymentConfig, err error)

	OrganizationCustomerMembers func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*CustomerMember, err error)

	OrganizationCustomerBenefitPolicies func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*CustomerBenefitPolicy, err error)

	OrganizationCustomerDailyPointGrantBudgets func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*CustomerDailyPointGrantBudget, err error)

	OrganizationCustomerPointEntries func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*CustomerPointEntry, err error)

	OrganizationCustomerCouponTemplates func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*CustomerCouponTemplate, err error)

	OrganizationProductCategories func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*ProductCategory, err error)

	OrganizationProductBrands func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*ProductBrand, err error)

	OrganizationSpecificationDefinitions func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*SpecificationDefinition, err error)

	OrganizationProductPackageTemplates func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*ProductPackageTemplate, err error)

	OrganizationProducts func(ctx context.Context, r *GeneratedResolver, obj *Organization) (res []*Product, err error)

	CreateOperatorMembership    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *OperatorMembership, err error)
	UpdateOperatorMembership    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *OperatorMembership, err error)
	DeleteOperatorMemberships   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryOperatorMemberships func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryOperatorMembership     func(ctx context.Context, r *GeneratedResolver, opts QueryOperatorMembershipHandlerOptions) (*OperatorMembership, error)
	QueryOperatorMemberships    func(ctx context.Context, r *GeneratedResolver, opts QueryOperatorMembershipsHandlerOptions) (*OperatorMembershipResultType, error)

	OperatorMembershipAccount func(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (res *Account, err error)

	OperatorMembershipOrganization func(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (res *Organization, err error)

	OperatorMembershipRoles func(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (res []*OperatorRole, err error)

	OperatorMembershipStores func(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (res []*Store, err error)

	OperatorMembershipInvitations func(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (res []*MembershipInvitation, err error)

	CreatePermission    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Permission, err error)
	UpdatePermission    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Permission, err error)
	DeletePermissions   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryPermissions func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryPermission     func(ctx context.Context, r *GeneratedResolver, opts QueryPermissionHandlerOptions) (*Permission, error)
	QueryPermissions    func(ctx context.Context, r *GeneratedResolver, opts QueryPermissionsHandlerOptions) (*PermissionResultType, error)

	PermissionRoles func(ctx context.Context, r *GeneratedResolver, obj *Permission) (res []*OperatorRole, err error)

	CreateOperatorRole    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *OperatorRole, err error)
	UpdateOperatorRole    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *OperatorRole, err error)
	DeleteOperatorRoles   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryOperatorRoles func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryOperatorRole     func(ctx context.Context, r *GeneratedResolver, opts QueryOperatorRoleHandlerOptions) (*OperatorRole, error)
	QueryOperatorRoles    func(ctx context.Context, r *GeneratedResolver, opts QueryOperatorRolesHandlerOptions) (*OperatorRoleResultType, error)

	OperatorRoleOrganization func(ctx context.Context, r *GeneratedResolver, obj *OperatorRole) (res *Organization, err error)

	OperatorRoleMembers func(ctx context.Context, r *GeneratedResolver, obj *OperatorRole) (res []*OperatorMembership, err error)

	OperatorRolePermissions func(ctx context.Context, r *GeneratedResolver, obj *OperatorRole) (res []*Permission, err error)

	CreateStore    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Store, err error)
	UpdateStore    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Store, err error)
	DeleteStores   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStores func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStore     func(ctx context.Context, r *GeneratedResolver, opts QueryStoreHandlerOptions) (*Store, error)
	QueryStores    func(ctx context.Context, r *GeneratedResolver, opts QueryStoresHandlerOptions) (*StoreResultType, error)

	StoreOrganization func(ctx context.Context, r *GeneratedResolver, obj *Store) (res *Organization, err error)

	StoreMembers func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*OperatorMembership, err error)

	StoreReviewedByAccount func(ctx context.Context, r *GeneratedResolver, obj *Store) (res *Account, err error)

	StoreAuditLogs func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*AuditLog, err error)

	StorePaymentConfigs func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*StorePaymentConfig, err error)

	StoreProductListings func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*StoreListing, err error)

	StoreStockMovements func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*StoreStockMovement, err error)

	StoreStocktakes func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*StoreStocktake, err error)

	StorePromotions func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*StorePromotion, err error)

	StoreCouponTemplates func(ctx context.Context, r *GeneratedResolver, obj *Store) (res []*CustomerCouponTemplate, err error)

	CreateSession    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Session, err error)
	UpdateSession    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Session, err error)
	DeleteSessions   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoverySessions func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QuerySession     func(ctx context.Context, r *GeneratedResolver, opts QuerySessionHandlerOptions) (*Session, error)
	QuerySessions    func(ctx context.Context, r *GeneratedResolver, opts QuerySessionsHandlerOptions) (*SessionResultType, error)

	SessionAccount func(ctx context.Context, r *GeneratedResolver, obj *Session) (res *Account, err error)

	SessionOrganization func(ctx context.Context, r *GeneratedResolver, obj *Session) (res *Organization, err error)

	SessionAuditLogs func(ctx context.Context, r *GeneratedResolver, obj *Session) (res []*AuditLog, err error)

	CreateMembershipInvitation    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *MembershipInvitation, err error)
	UpdateMembershipInvitation    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *MembershipInvitation, err error)
	DeleteMembershipInvitations   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryMembershipInvitations func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryMembershipInvitation     func(ctx context.Context, r *GeneratedResolver, opts QueryMembershipInvitationHandlerOptions) (*MembershipInvitation, error)
	QueryMembershipInvitations    func(ctx context.Context, r *GeneratedResolver, opts QueryMembershipInvitationsHandlerOptions) (*MembershipInvitationResultType, error)

	MembershipInvitationMembership func(ctx context.Context, r *GeneratedResolver, obj *MembershipInvitation) (res *OperatorMembership, err error)

	MembershipInvitationInvitedByAccount func(ctx context.Context, r *GeneratedResolver, obj *MembershipInvitation) (res *Account, err error)

	CreateAuditLog    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *AuditLog, err error)
	UpdateAuditLog    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *AuditLog, err error)
	DeleteAuditLogs   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryAuditLogs func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryAuditLog     func(ctx context.Context, r *GeneratedResolver, opts QueryAuditLogHandlerOptions) (*AuditLog, error)
	QueryAuditLogs    func(ctx context.Context, r *GeneratedResolver, opts QueryAuditLogsHandlerOptions) (*AuditLogResultType, error)

	AuditLogActorAccount func(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (res *Account, err error)

	AuditLogSession func(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (res *Session, err error)

	AuditLogOrganization func(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (res *Organization, err error)

	AuditLogStore func(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (res *Store, err error)

	CreateFranchiseOpeningRecord    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *FranchiseOpeningRecord, err error)
	UpdateFranchiseOpeningRecord    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *FranchiseOpeningRecord, err error)
	DeleteFranchiseOpeningRecords   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryFranchiseOpeningRecords func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryFranchiseOpeningRecord     func(ctx context.Context, r *GeneratedResolver, opts QueryFranchiseOpeningRecordHandlerOptions) (*FranchiseOpeningRecord, error)
	QueryFranchiseOpeningRecords    func(ctx context.Context, r *GeneratedResolver, opts QueryFranchiseOpeningRecordsHandlerOptions) (*FranchiseOpeningRecordResultType, error)

	FranchiseOpeningRecordOrganization func(ctx context.Context, r *GeneratedResolver, obj *FranchiseOpeningRecord) (res *Organization, err error)

	FranchiseOpeningRecordInitialAccount func(ctx context.Context, r *GeneratedResolver, obj *FranchiseOpeningRecord) (res *Account, err error)

	FranchiseOpeningRecordRecordedByAccount func(ctx context.Context, r *GeneratedResolver, obj *FranchiseOpeningRecord) (res *Account, err error)

	CreateGlobalPaymentConfig    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *GlobalPaymentConfig, err error)
	UpdateGlobalPaymentConfig    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *GlobalPaymentConfig, err error)
	DeleteGlobalPaymentConfigs   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryGlobalPaymentConfigs func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryGlobalPaymentConfig     func(ctx context.Context, r *GeneratedResolver, opts QueryGlobalPaymentConfigHandlerOptions) (*GlobalPaymentConfig, error)
	QueryGlobalPaymentConfigs    func(ctx context.Context, r *GeneratedResolver, opts QueryGlobalPaymentConfigsHandlerOptions) (*GlobalPaymentConfigResultType, error)

	CreateFranchisePaymentConfig    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *FranchisePaymentConfig, err error)
	UpdateFranchisePaymentConfig    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *FranchisePaymentConfig, err error)
	DeleteFranchisePaymentConfigs   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryFranchisePaymentConfigs func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryFranchisePaymentConfig     func(ctx context.Context, r *GeneratedResolver, opts QueryFranchisePaymentConfigHandlerOptions) (*FranchisePaymentConfig, error)
	QueryFranchisePaymentConfigs    func(ctx context.Context, r *GeneratedResolver, opts QueryFranchisePaymentConfigsHandlerOptions) (*FranchisePaymentConfigResultType, error)

	FranchisePaymentConfigOrganization func(ctx context.Context, r *GeneratedResolver, obj *FranchisePaymentConfig) (res *Organization, err error)

	CreateStorePaymentConfig    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StorePaymentConfig, err error)
	UpdateStorePaymentConfig    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StorePaymentConfig, err error)
	DeleteStorePaymentConfigs   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStorePaymentConfigs func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStorePaymentConfig     func(ctx context.Context, r *GeneratedResolver, opts QueryStorePaymentConfigHandlerOptions) (*StorePaymentConfig, error)
	QueryStorePaymentConfigs    func(ctx context.Context, r *GeneratedResolver, opts QueryStorePaymentConfigsHandlerOptions) (*StorePaymentConfigResultType, error)

	StorePaymentConfigStore func(ctx context.Context, r *GeneratedResolver, obj *StorePaymentConfig) (res *Store, err error)

	CreateCustomerMember    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *CustomerMember, err error)
	UpdateCustomerMember    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *CustomerMember, err error)
	DeleteCustomerMembers   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryCustomerMembers func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryCustomerMember     func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerMemberHandlerOptions) (*CustomerMember, error)
	QueryCustomerMembers    func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerMembersHandlerOptions) (*CustomerMemberResultType, error)

	CustomerMemberOrganization func(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (res *Organization, err error)

	CustomerMemberPointEntries func(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (res []*CustomerPointEntry, err error)

	CustomerMemberCouponGrants func(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (res []*CustomerCouponGrant, err error)

	CustomerMemberCouponDistributionJobs func(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (res []*CustomerCouponDistributionJob, err error)

	CreateCustomerBenefitPolicy     func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *CustomerBenefitPolicy, err error)
	UpdateCustomerBenefitPolicy     func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *CustomerBenefitPolicy, err error)
	DeleteCustomerBenefitPolicies   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryCustomerBenefitPolicies func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryCustomerBenefitPolicy      func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerBenefitPolicyHandlerOptions) (*CustomerBenefitPolicy, error)
	QueryCustomerBenefitPolicies    func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerBenefitPoliciesHandlerOptions) (*CustomerBenefitPolicyResultType, error)

	CustomerBenefitPolicyOrganization func(ctx context.Context, r *GeneratedResolver, obj *CustomerBenefitPolicy) (res *Organization, err error)

	CreateCustomerDailyPointGrantBudget    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *CustomerDailyPointGrantBudget, err error)
	UpdateCustomerDailyPointGrantBudget    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *CustomerDailyPointGrantBudget, err error)
	DeleteCustomerDailyPointGrantBudgets   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryCustomerDailyPointGrantBudgets func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryCustomerDailyPointGrantBudget     func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerDailyPointGrantBudgetHandlerOptions) (*CustomerDailyPointGrantBudget, error)
	QueryCustomerDailyPointGrantBudgets    func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerDailyPointGrantBudgetsHandlerOptions) (*CustomerDailyPointGrantBudgetResultType, error)

	CustomerDailyPointGrantBudgetOrganization func(ctx context.Context, r *GeneratedResolver, obj *CustomerDailyPointGrantBudget) (res *Organization, err error)

	CreateCustomerPointEntry     func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *CustomerPointEntry, err error)
	UpdateCustomerPointEntry     func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *CustomerPointEntry, err error)
	DeleteCustomerPointEntries   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryCustomerPointEntries func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryCustomerPointEntry      func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerPointEntryHandlerOptions) (*CustomerPointEntry, error)
	QueryCustomerPointEntries    func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerPointEntriesHandlerOptions) (*CustomerPointEntryResultType, error)

	CustomerPointEntryMember func(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (res *CustomerMember, err error)

	CustomerPointEntrySourceOrganization func(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (res *Organization, err error)

	CustomerPointEntryReverses func(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (res *CustomerPointEntry, err error)

	CustomerPointEntryReversedBy func(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (res *CustomerPointEntry, err error)

	CreateCustomerCouponTemplate    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *CustomerCouponTemplate, err error)
	UpdateCustomerCouponTemplate    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *CustomerCouponTemplate, err error)
	DeleteCustomerCouponTemplates   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryCustomerCouponTemplates func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryCustomerCouponTemplate     func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponTemplateHandlerOptions) (*CustomerCouponTemplate, error)
	QueryCustomerCouponTemplates    func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponTemplatesHandlerOptions) (*CustomerCouponTemplateResultType, error)

	CustomerCouponTemplateOrganization func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (res *Organization, err error)

	CustomerCouponTemplateApplicableStore func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (res *Store, err error)

	CustomerCouponTemplateGrants func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (res []*CustomerCouponGrant, err error)

	CustomerCouponTemplateDistributionJobs func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (res []*CustomerCouponDistributionJob, err error)

	CreateProductCategory     func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *ProductCategory, err error)
	UpdateProductCategory     func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *ProductCategory, err error)
	DeleteProductCategories   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProductCategories func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProductCategory      func(ctx context.Context, r *GeneratedResolver, opts QueryProductCategoryHandlerOptions) (*ProductCategory, error)
	QueryProductCategories    func(ctx context.Context, r *GeneratedResolver, opts QueryProductCategoriesHandlerOptions) (*ProductCategoryResultType, error)

	ProductCategoryOrganization func(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (res *Organization, err error)

	ProductCategoryParent func(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (res *ProductCategory, err error)

	ProductCategoryChildren func(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (res []*ProductCategory, err error)

	ProductCategoryProducts func(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (res []*Product, err error)

	CreateProductBrand    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *ProductBrand, err error)
	UpdateProductBrand    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *ProductBrand, err error)
	DeleteProductBrands   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProductBrands func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProductBrand     func(ctx context.Context, r *GeneratedResolver, opts QueryProductBrandHandlerOptions) (*ProductBrand, error)
	QueryProductBrands    func(ctx context.Context, r *GeneratedResolver, opts QueryProductBrandsHandlerOptions) (*ProductBrandResultType, error)

	ProductBrandOrganization func(ctx context.Context, r *GeneratedResolver, obj *ProductBrand) (res *Organization, err error)

	ProductBrandProducts func(ctx context.Context, r *GeneratedResolver, obj *ProductBrand) (res []*Product, err error)

	CreateProduct    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Product, err error)
	UpdateProduct    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Product, err error)
	DeleteProducts   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProducts func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProduct     func(ctx context.Context, r *GeneratedResolver, opts QueryProductHandlerOptions) (*Product, error)
	QueryProducts    func(ctx context.Context, r *GeneratedResolver, opts QueryProductsHandlerOptions) (*ProductResultType, error)

	ProductBrand func(ctx context.Context, r *GeneratedResolver, obj *Product) (res *ProductBrand, err error)

	ProductOrganization func(ctx context.Context, r *GeneratedResolver, obj *Product) (res *Organization, err error)

	ProductCategory func(ctx context.Context, r *GeneratedResolver, obj *Product) (res *ProductCategory, err error)

	ProductDefaultPackageTemplate func(ctx context.Context, r *GeneratedResolver, obj *Product) (res *ProductPackageTemplate, err error)

	ProductSpecificationChoices func(ctx context.Context, r *GeneratedResolver, obj *Product) (res []*ProductSpecificationChoice, err error)

	ProductSkus func(ctx context.Context, r *GeneratedResolver, obj *Product) (res []*ProductSku, err error)

	CreateProductSku    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *ProductSku, err error)
	UpdateProductSku    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *ProductSku, err error)
	DeleteProductSkus   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProductSkus func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProductSku     func(ctx context.Context, r *GeneratedResolver, opts QueryProductSkuHandlerOptions) (*ProductSku, error)
	QueryProductSkus    func(ctx context.Context, r *GeneratedResolver, opts QueryProductSkusHandlerOptions) (*ProductSkuResultType, error)

	ProductSkuProduct func(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (res *Product, err error)

	ProductSkuSpecificationValues func(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (res []*ProductSkuSpecificationValue, err error)

	ProductSkuPackages func(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (res []*ProductPackage, err error)

	ProductSkuListings func(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (res []*StoreListing, err error)

	CreateProductPackage    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *ProductPackage, err error)
	UpdateProductPackage    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *ProductPackage, err error)
	DeleteProductPackages   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProductPackages func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProductPackage     func(ctx context.Context, r *GeneratedResolver, opts QueryProductPackageHandlerOptions) (*ProductPackage, error)
	QueryProductPackages    func(ctx context.Context, r *GeneratedResolver, opts QueryProductPackagesHandlerOptions) (*ProductPackageResultType, error)

	ProductPackageSku func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res *ProductSku, err error)

	ProductPackageTemplate func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res *ProductPackageTemplate, err error)

	ProductPackageContainsPackage func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res *ProductPackage, err error)

	ProductPackageContainedByPackages func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res []*ProductPackage, err error)

	ProductPackageOffers func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res []*StorePackageOffer, err error)

	ProductPackageBalances func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res []*StoreStockBalance, err error)

	ProductPackageMovementSources func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res []*StoreStockMovement, err error)

	ProductPackageMovementTargets func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res []*StoreStockMovement, err error)

	ProductPackageStocktakeLines func(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (res []*StoreStocktakeLine, err error)

	CreateSpecificationDefinition    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *SpecificationDefinition, err error)
	UpdateSpecificationDefinition    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *SpecificationDefinition, err error)
	DeleteSpecificationDefinitions   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoverySpecificationDefinitions func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QuerySpecificationDefinition     func(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationDefinitionHandlerOptions) (*SpecificationDefinition, error)
	QuerySpecificationDefinitions    func(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationDefinitionsHandlerOptions) (*SpecificationDefinitionResultType, error)

	SpecificationDefinitionOrganization func(ctx context.Context, r *GeneratedResolver, obj *SpecificationDefinition) (res *Organization, err error)

	SpecificationDefinitionValues func(ctx context.Context, r *GeneratedResolver, obj *SpecificationDefinition) (res []*SpecificationValue, err error)

	CreateSpecificationValue    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *SpecificationValue, err error)
	UpdateSpecificationValue    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *SpecificationValue, err error)
	DeleteSpecificationValues   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoverySpecificationValues func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QuerySpecificationValue     func(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationValueHandlerOptions) (*SpecificationValue, error)
	QuerySpecificationValues    func(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationValuesHandlerOptions) (*SpecificationValueResultType, error)

	SpecificationValueSpecification func(ctx context.Context, r *GeneratedResolver, obj *SpecificationValue) (res *SpecificationDefinition, err error)

	SpecificationValueProductChoices func(ctx context.Context, r *GeneratedResolver, obj *SpecificationValue) (res []*ProductSpecificationChoice, err error)

	SpecificationValueSkuValues func(ctx context.Context, r *GeneratedResolver, obj *SpecificationValue) (res []*ProductSkuSpecificationValue, err error)

	CreateProductSpecificationChoice    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *ProductSpecificationChoice, err error)
	UpdateProductSpecificationChoice    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *ProductSpecificationChoice, err error)
	DeleteProductSpecificationChoices   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProductSpecificationChoices func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProductSpecificationChoice     func(ctx context.Context, r *GeneratedResolver, opts QueryProductSpecificationChoiceHandlerOptions) (*ProductSpecificationChoice, error)
	QueryProductSpecificationChoices    func(ctx context.Context, r *GeneratedResolver, opts QueryProductSpecificationChoicesHandlerOptions) (*ProductSpecificationChoiceResultType, error)

	ProductSpecificationChoiceProduct func(ctx context.Context, r *GeneratedResolver, obj *ProductSpecificationChoice) (res *Product, err error)

	ProductSpecificationChoiceValue func(ctx context.Context, r *GeneratedResolver, obj *ProductSpecificationChoice) (res *SpecificationValue, err error)

	CreateProductSkuSpecificationValue    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *ProductSkuSpecificationValue, err error)
	UpdateProductSkuSpecificationValue    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *ProductSkuSpecificationValue, err error)
	DeleteProductSkuSpecificationValues   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProductSkuSpecificationValues func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProductSkuSpecificationValue     func(ctx context.Context, r *GeneratedResolver, opts QueryProductSkuSpecificationValueHandlerOptions) (*ProductSkuSpecificationValue, error)
	QueryProductSkuSpecificationValues    func(ctx context.Context, r *GeneratedResolver, opts QueryProductSkuSpecificationValuesHandlerOptions) (*ProductSkuSpecificationValueResultType, error)

	ProductSkuSpecificationValueSku func(ctx context.Context, r *GeneratedResolver, obj *ProductSkuSpecificationValue) (res *ProductSku, err error)

	ProductSkuSpecificationValueValue func(ctx context.Context, r *GeneratedResolver, obj *ProductSkuSpecificationValue) (res *SpecificationValue, err error)

	CreateProductPackageTemplate    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *ProductPackageTemplate, err error)
	UpdateProductPackageTemplate    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *ProductPackageTemplate, err error)
	DeleteProductPackageTemplates   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryProductPackageTemplates func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryProductPackageTemplate     func(ctx context.Context, r *GeneratedResolver, opts QueryProductPackageTemplateHandlerOptions) (*ProductPackageTemplate, error)
	QueryProductPackageTemplates    func(ctx context.Context, r *GeneratedResolver, opts QueryProductPackageTemplatesHandlerOptions) (*ProductPackageTemplateResultType, error)

	ProductPackageTemplateOrganization func(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (res *Organization, err error)

	ProductPackageTemplateDefaultProducts func(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (res []*Product, err error)

	ProductPackageTemplateContainsPackage func(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (res *ProductPackageTemplate, err error)

	ProductPackageTemplateContainedByPackages func(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (res []*ProductPackageTemplate, err error)

	ProductPackageTemplateCreatedPackages func(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (res []*ProductPackage, err error)

	CreateStoreListing    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StoreListing, err error)
	UpdateStoreListing    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StoreListing, err error)
	DeleteStoreListings   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStoreListings func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStoreListing     func(ctx context.Context, r *GeneratedResolver, opts QueryStoreListingHandlerOptions) (*StoreListing, error)
	QueryStoreListings    func(ctx context.Context, r *GeneratedResolver, opts QueryStoreListingsHandlerOptions) (*StoreListingResultType, error)

	StoreListingStore func(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (res *Store, err error)

	StoreListingSku func(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (res *ProductSku, err error)

	StoreListingOffers func(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (res []*StorePackageOffer, err error)

	StoreListingBatches func(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (res []*StoreInventoryBatch, err error)

	CreateStorePackageOffer    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StorePackageOffer, err error)
	UpdateStorePackageOffer    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StorePackageOffer, err error)
	DeleteStorePackageOffers   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStorePackageOffers func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStorePackageOffer     func(ctx context.Context, r *GeneratedResolver, opts QueryStorePackageOfferHandlerOptions) (*StorePackageOffer, error)
	QueryStorePackageOffers    func(ctx context.Context, r *GeneratedResolver, opts QueryStorePackageOffersHandlerOptions) (*StorePackageOfferResultType, error)

	StorePackageOfferListing func(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (res *StoreListing, err error)

	StorePackageOfferPackage func(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (res *ProductPackage, err error)

	StorePackageOfferPriceRevisions func(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (res []*StorePriceRevision, err error)

	StorePackageOfferPromotionTargets func(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (res []*StorePromotionTarget, err error)

	CreateStorePriceRevision    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StorePriceRevision, err error)
	UpdateStorePriceRevision    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StorePriceRevision, err error)
	DeleteStorePriceRevisions   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStorePriceRevisions func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStorePriceRevision     func(ctx context.Context, r *GeneratedResolver, opts QueryStorePriceRevisionHandlerOptions) (*StorePriceRevision, error)
	QueryStorePriceRevisions    func(ctx context.Context, r *GeneratedResolver, opts QueryStorePriceRevisionsHandlerOptions) (*StorePriceRevisionResultType, error)

	StorePriceRevisionOffer func(ctx context.Context, r *GeneratedResolver, obj *StorePriceRevision) (res *StorePackageOffer, err error)

	CreateStoreInventoryBatch     func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StoreInventoryBatch, err error)
	UpdateStoreInventoryBatch     func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StoreInventoryBatch, err error)
	DeleteStoreInventoryBatches   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStoreInventoryBatches func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStoreInventoryBatch      func(ctx context.Context, r *GeneratedResolver, opts QueryStoreInventoryBatchHandlerOptions) (*StoreInventoryBatch, error)
	QueryStoreInventoryBatches    func(ctx context.Context, r *GeneratedResolver, opts QueryStoreInventoryBatchesHandlerOptions) (*StoreInventoryBatchResultType, error)

	StoreInventoryBatchListing func(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (res *StoreListing, err error)

	StoreInventoryBatchBalances func(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (res []*StoreStockBalance, err error)

	StoreInventoryBatchMovements func(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (res []*StoreStockMovement, err error)

	StoreInventoryBatchStocktakeLines func(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (res []*StoreStocktakeLine, err error)

	CreateStoreStockBalance    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StoreStockBalance, err error)
	UpdateStoreStockBalance    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StoreStockBalance, err error)
	DeleteStoreStockBalances   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStoreStockBalances func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStoreStockBalance     func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockBalanceHandlerOptions) (*StoreStockBalance, error)
	QueryStoreStockBalances    func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockBalancesHandlerOptions) (*StoreStockBalanceResultType, error)

	StoreStockBalanceBatch func(ctx context.Context, r *GeneratedResolver, obj *StoreStockBalance) (res *StoreInventoryBatch, err error)

	StoreStockBalancePackage func(ctx context.Context, r *GeneratedResolver, obj *StoreStockBalance) (res *ProductPackage, err error)

	CreateStoreStocktake    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StoreStocktake, err error)
	UpdateStoreStocktake    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StoreStocktake, err error)
	DeleteStoreStocktakes   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStoreStocktakes func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStoreStocktake     func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakeHandlerOptions) (*StoreStocktake, error)
	QueryStoreStocktakes    func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakesHandlerOptions) (*StoreStocktakeResultType, error)

	StoreStocktakeStore func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (res *Store, err error)

	StoreStocktakeInitiatedByAccount func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (res *Account, err error)

	StoreStocktakePostedBy func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (res *Account, err error)

	StoreStocktakeLines func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (res []*StoreStocktakeLine, err error)

	CreateStoreStocktakeLine    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StoreStocktakeLine, err error)
	UpdateStoreStocktakeLine    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StoreStocktakeLine, err error)
	DeleteStoreStocktakeLines   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStoreStocktakeLines func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStoreStocktakeLine     func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakeLineHandlerOptions) (*StoreStocktakeLine, error)
	QueryStoreStocktakeLines    func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakeLinesHandlerOptions) (*StoreStocktakeLineResultType, error)

	StoreStocktakeLineStocktake func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (res *StoreStocktake, err error)

	StoreStocktakeLineBatch func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (res *StoreInventoryBatch, err error)

	StoreStocktakeLinePackage func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (res *ProductPackage, err error)

	StoreStocktakeLineMovements func(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (res []*StoreStockMovement, err error)

	CreateStoreStockMovement    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StoreStockMovement, err error)
	UpdateStoreStockMovement    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StoreStockMovement, err error)
	DeleteStoreStockMovements   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStoreStockMovements func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStoreStockMovement     func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockMovementHandlerOptions) (*StoreStockMovement, error)
	QueryStoreStockMovements    func(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockMovementsHandlerOptions) (*StoreStockMovementResultType, error)

	StoreStockMovementStore func(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (res *Store, err error)

	StoreStockMovementBatch func(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (res *StoreInventoryBatch, err error)

	StoreStockMovementSourcePackage func(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (res *ProductPackage, err error)

	StoreStockMovementTargetPackage func(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (res *ProductPackage, err error)

	StoreStockMovementStocktakeLine func(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (res *StoreStocktakeLine, err error)

	CreateStorePromotion    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StorePromotion, err error)
	UpdateStorePromotion    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StorePromotion, err error)
	DeleteStorePromotions   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStorePromotions func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStorePromotion     func(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionHandlerOptions) (*StorePromotion, error)
	QueryStorePromotions    func(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionsHandlerOptions) (*StorePromotionResultType, error)

	StorePromotionStore func(ctx context.Context, r *GeneratedResolver, obj *StorePromotion) (res *Store, err error)

	StorePromotionTargets func(ctx context.Context, r *GeneratedResolver, obj *StorePromotion) (res []*StorePromotionTarget, err error)

	CreateStorePromotionTarget    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StorePromotionTarget, err error)
	UpdateStorePromotionTarget    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StorePromotionTarget, err error)
	DeleteStorePromotionTargets   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryStorePromotionTargets func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryStorePromotionTarget     func(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionTargetHandlerOptions) (*StorePromotionTarget, error)
	QueryStorePromotionTargets    func(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionTargetsHandlerOptions) (*StorePromotionTargetResultType, error)

	StorePromotionTargetPromotion func(ctx context.Context, r *GeneratedResolver, obj *StorePromotionTarget) (res *StorePromotion, err error)

	StorePromotionTargetOffer func(ctx context.Context, r *GeneratedResolver, obj *StorePromotionTarget) (res *StorePackageOffer, err error)

	CreateCustomerCouponGrant    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *CustomerCouponGrant, err error)
	UpdateCustomerCouponGrant    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *CustomerCouponGrant, err error)
	DeleteCustomerCouponGrants   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryCustomerCouponGrants func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryCustomerCouponGrant     func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponGrantHandlerOptions) (*CustomerCouponGrant, error)
	QueryCustomerCouponGrants    func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponGrantsHandlerOptions) (*CustomerCouponGrantResultType, error)

	CustomerCouponGrantMember func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponGrant) (res *CustomerMember, err error)

	CustomerCouponGrantTemplate func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponGrant) (res *CustomerCouponTemplate, err error)

	CreateCustomerCouponDistributionJob    func(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *CustomerCouponDistributionJob, err error)
	UpdateCustomerCouponDistributionJob    func(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *CustomerCouponDistributionJob, err error)
	DeleteCustomerCouponDistributionJobs   func(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error)
	RecoveryCustomerCouponDistributionJobs func(ctx context.Context, r *GeneratedResolver, id []string) (bool, error)
	QueryCustomerCouponDistributionJob     func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponDistributionJobHandlerOptions) (*CustomerCouponDistributionJob, error)
	QueryCustomerCouponDistributionJobs    func(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponDistributionJobsHandlerOptions) (*CustomerCouponDistributionJobResultType, error)

	CustomerCouponDistributionJobTemplate func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponDistributionJob) (res *CustomerCouponTemplate, err error)

	CustomerCouponDistributionJobMember func(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponDistributionJob) (res *CustomerMember, err error)
}

func DefaultResolutionHandlers() ResolutionHandlers {
	handlers := ResolutionHandlers{
		OnEvent: func(ctx context.Context, r *GeneratedResolver, e *Event) error { return nil },

		CreateAccount:    CreateAccountHandler,
		UpdateAccount:    UpdateAccountHandler,
		DeleteAccounts:   DeleteAccountsHandler,
		RecoveryAccounts: RecoveryAccountsHandler,
		QueryAccount:     QueryAccountHandler,
		QueryAccounts:    QueryAccountsHandler,

		AccountMemberships: AccountMembershipsHandler,

		AccountInitializedOrganizations: AccountInitializedOrganizationsHandler,

		AccountOpeningRecords: AccountOpeningRecordsHandler,

		AccountRecordedOpeningRecords: AccountRecordedOpeningRecordsHandler,

		AccountSessions: AccountSessionsHandler,

		AccountReviewedStores: AccountReviewedStoresHandler,

		AccountSentMembershipInvitations: AccountSentMembershipInvitationsHandler,

		AccountAuditLogs: AccountAuditLogsHandler,

		AccountCreatedStocktakes: AccountCreatedStocktakesHandler,

		AccountPostedStocktakes: AccountPostedStocktakesHandler,

		CreateOrganization:    CreateOrganizationHandler,
		UpdateOrganization:    UpdateOrganizationHandler,
		DeleteOrganizations:   DeleteOrganizationsHandler,
		RecoveryOrganizations: RecoveryOrganizationsHandler,
		QueryOrganization:     QueryOrganizationHandler,
		QueryOrganizations:    QueryOrganizationsHandler,

		OrganizationMemberships: OrganizationMembershipsHandler,

		OrganizationInitialAccount: OrganizationInitialAccountHandler,

		OrganizationOpeningRecords: OrganizationOpeningRecordsHandler,

		OrganizationStores: OrganizationStoresHandler,

		OrganizationRoles: OrganizationRolesHandler,

		OrganizationSessions: OrganizationSessionsHandler,

		OrganizationAuditLogs: OrganizationAuditLogsHandler,

		OrganizationPaymentConfigs: OrganizationPaymentConfigsHandler,

		OrganizationCustomerMembers: OrganizationCustomerMembersHandler,

		OrganizationCustomerBenefitPolicies: OrganizationCustomerBenefitPoliciesHandler,

		OrganizationCustomerDailyPointGrantBudgets: OrganizationCustomerDailyPointGrantBudgetsHandler,

		OrganizationCustomerPointEntries: OrganizationCustomerPointEntriesHandler,

		OrganizationCustomerCouponTemplates: OrganizationCustomerCouponTemplatesHandler,

		OrganizationProductCategories: OrganizationProductCategoriesHandler,

		OrganizationProductBrands: OrganizationProductBrandsHandler,

		OrganizationSpecificationDefinitions: OrganizationSpecificationDefinitionsHandler,

		OrganizationProductPackageTemplates: OrganizationProductPackageTemplatesHandler,

		OrganizationProducts: OrganizationProductsHandler,

		CreateOperatorMembership:    CreateOperatorMembershipHandler,
		UpdateOperatorMembership:    UpdateOperatorMembershipHandler,
		DeleteOperatorMemberships:   DeleteOperatorMembershipsHandler,
		RecoveryOperatorMemberships: RecoveryOperatorMembershipsHandler,
		QueryOperatorMembership:     QueryOperatorMembershipHandler,
		QueryOperatorMemberships:    QueryOperatorMembershipsHandler,

		OperatorMembershipAccount: OperatorMembershipAccountHandler,

		OperatorMembershipOrganization: OperatorMembershipOrganizationHandler,

		OperatorMembershipRoles: OperatorMembershipRolesHandler,

		OperatorMembershipStores: OperatorMembershipStoresHandler,

		OperatorMembershipInvitations: OperatorMembershipInvitationsHandler,

		CreatePermission:    CreatePermissionHandler,
		UpdatePermission:    UpdatePermissionHandler,
		DeletePermissions:   DeletePermissionsHandler,
		RecoveryPermissions: RecoveryPermissionsHandler,
		QueryPermission:     QueryPermissionHandler,
		QueryPermissions:    QueryPermissionsHandler,

		PermissionRoles: PermissionRolesHandler,

		CreateOperatorRole:    CreateOperatorRoleHandler,
		UpdateOperatorRole:    UpdateOperatorRoleHandler,
		DeleteOperatorRoles:   DeleteOperatorRolesHandler,
		RecoveryOperatorRoles: RecoveryOperatorRolesHandler,
		QueryOperatorRole:     QueryOperatorRoleHandler,
		QueryOperatorRoles:    QueryOperatorRolesHandler,

		OperatorRoleOrganization: OperatorRoleOrganizationHandler,

		OperatorRoleMembers: OperatorRoleMembersHandler,

		OperatorRolePermissions: OperatorRolePermissionsHandler,

		CreateStore:    CreateStoreHandler,
		UpdateStore:    UpdateStoreHandler,
		DeleteStores:   DeleteStoresHandler,
		RecoveryStores: RecoveryStoresHandler,
		QueryStore:     QueryStoreHandler,
		QueryStores:    QueryStoresHandler,

		StoreOrganization: StoreOrganizationHandler,

		StoreMembers: StoreMembersHandler,

		StoreReviewedByAccount: StoreReviewedByAccountHandler,

		StoreAuditLogs: StoreAuditLogsHandler,

		StorePaymentConfigs: StorePaymentConfigsHandler,

		StoreProductListings: StoreProductListingsHandler,

		StoreStockMovements: StoreStockMovementsHandler,

		StoreStocktakes: StoreStocktakesHandler,

		StorePromotions: StorePromotionsHandler,

		StoreCouponTemplates: StoreCouponTemplatesHandler,

		CreateSession:    CreateSessionHandler,
		UpdateSession:    UpdateSessionHandler,
		DeleteSessions:   DeleteSessionsHandler,
		RecoverySessions: RecoverySessionsHandler,
		QuerySession:     QuerySessionHandler,
		QuerySessions:    QuerySessionsHandler,

		SessionAccount: SessionAccountHandler,

		SessionOrganization: SessionOrganizationHandler,

		SessionAuditLogs: SessionAuditLogsHandler,

		CreateMembershipInvitation:    CreateMembershipInvitationHandler,
		UpdateMembershipInvitation:    UpdateMembershipInvitationHandler,
		DeleteMembershipInvitations:   DeleteMembershipInvitationsHandler,
		RecoveryMembershipInvitations: RecoveryMembershipInvitationsHandler,
		QueryMembershipInvitation:     QueryMembershipInvitationHandler,
		QueryMembershipInvitations:    QueryMembershipInvitationsHandler,

		MembershipInvitationMembership: MembershipInvitationMembershipHandler,

		MembershipInvitationInvitedByAccount: MembershipInvitationInvitedByAccountHandler,

		CreateAuditLog:    CreateAuditLogHandler,
		UpdateAuditLog:    UpdateAuditLogHandler,
		DeleteAuditLogs:   DeleteAuditLogsHandler,
		RecoveryAuditLogs: RecoveryAuditLogsHandler,
		QueryAuditLog:     QueryAuditLogHandler,
		QueryAuditLogs:    QueryAuditLogsHandler,

		AuditLogActorAccount: AuditLogActorAccountHandler,

		AuditLogSession: AuditLogSessionHandler,

		AuditLogOrganization: AuditLogOrganizationHandler,

		AuditLogStore: AuditLogStoreHandler,

		CreateFranchiseOpeningRecord:    CreateFranchiseOpeningRecordHandler,
		UpdateFranchiseOpeningRecord:    UpdateFranchiseOpeningRecordHandler,
		DeleteFranchiseOpeningRecords:   DeleteFranchiseOpeningRecordsHandler,
		RecoveryFranchiseOpeningRecords: RecoveryFranchiseOpeningRecordsHandler,
		QueryFranchiseOpeningRecord:     QueryFranchiseOpeningRecordHandler,
		QueryFranchiseOpeningRecords:    QueryFranchiseOpeningRecordsHandler,

		FranchiseOpeningRecordOrganization: FranchiseOpeningRecordOrganizationHandler,

		FranchiseOpeningRecordInitialAccount: FranchiseOpeningRecordInitialAccountHandler,

		FranchiseOpeningRecordRecordedByAccount: FranchiseOpeningRecordRecordedByAccountHandler,

		CreateGlobalPaymentConfig:    CreateGlobalPaymentConfigHandler,
		UpdateGlobalPaymentConfig:    UpdateGlobalPaymentConfigHandler,
		DeleteGlobalPaymentConfigs:   DeleteGlobalPaymentConfigsHandler,
		RecoveryGlobalPaymentConfigs: RecoveryGlobalPaymentConfigsHandler,
		QueryGlobalPaymentConfig:     QueryGlobalPaymentConfigHandler,
		QueryGlobalPaymentConfigs:    QueryGlobalPaymentConfigsHandler,

		CreateFranchisePaymentConfig:    CreateFranchisePaymentConfigHandler,
		UpdateFranchisePaymentConfig:    UpdateFranchisePaymentConfigHandler,
		DeleteFranchisePaymentConfigs:   DeleteFranchisePaymentConfigsHandler,
		RecoveryFranchisePaymentConfigs: RecoveryFranchisePaymentConfigsHandler,
		QueryFranchisePaymentConfig:     QueryFranchisePaymentConfigHandler,
		QueryFranchisePaymentConfigs:    QueryFranchisePaymentConfigsHandler,

		FranchisePaymentConfigOrganization: FranchisePaymentConfigOrganizationHandler,

		CreateStorePaymentConfig:    CreateStorePaymentConfigHandler,
		UpdateStorePaymentConfig:    UpdateStorePaymentConfigHandler,
		DeleteStorePaymentConfigs:   DeleteStorePaymentConfigsHandler,
		RecoveryStorePaymentConfigs: RecoveryStorePaymentConfigsHandler,
		QueryStorePaymentConfig:     QueryStorePaymentConfigHandler,
		QueryStorePaymentConfigs:    QueryStorePaymentConfigsHandler,

		StorePaymentConfigStore: StorePaymentConfigStoreHandler,

		CreateCustomerMember:    CreateCustomerMemberHandler,
		UpdateCustomerMember:    UpdateCustomerMemberHandler,
		DeleteCustomerMembers:   DeleteCustomerMembersHandler,
		RecoveryCustomerMembers: RecoveryCustomerMembersHandler,
		QueryCustomerMember:     QueryCustomerMemberHandler,
		QueryCustomerMembers:    QueryCustomerMembersHandler,

		CustomerMemberOrganization: CustomerMemberOrganizationHandler,

		CustomerMemberPointEntries: CustomerMemberPointEntriesHandler,

		CustomerMemberCouponGrants: CustomerMemberCouponGrantsHandler,

		CustomerMemberCouponDistributionJobs: CustomerMemberCouponDistributionJobsHandler,

		CreateCustomerBenefitPolicy:     CreateCustomerBenefitPolicyHandler,
		UpdateCustomerBenefitPolicy:     UpdateCustomerBenefitPolicyHandler,
		DeleteCustomerBenefitPolicies:   DeleteCustomerBenefitPoliciesHandler,
		RecoveryCustomerBenefitPolicies: RecoveryCustomerBenefitPoliciesHandler,
		QueryCustomerBenefitPolicy:      QueryCustomerBenefitPolicyHandler,
		QueryCustomerBenefitPolicies:    QueryCustomerBenefitPoliciesHandler,

		CustomerBenefitPolicyOrganization: CustomerBenefitPolicyOrganizationHandler,

		CreateCustomerDailyPointGrantBudget:    CreateCustomerDailyPointGrantBudgetHandler,
		UpdateCustomerDailyPointGrantBudget:    UpdateCustomerDailyPointGrantBudgetHandler,
		DeleteCustomerDailyPointGrantBudgets:   DeleteCustomerDailyPointGrantBudgetsHandler,
		RecoveryCustomerDailyPointGrantBudgets: RecoveryCustomerDailyPointGrantBudgetsHandler,
		QueryCustomerDailyPointGrantBudget:     QueryCustomerDailyPointGrantBudgetHandler,
		QueryCustomerDailyPointGrantBudgets:    QueryCustomerDailyPointGrantBudgetsHandler,

		CustomerDailyPointGrantBudgetOrganization: CustomerDailyPointGrantBudgetOrganizationHandler,

		CreateCustomerPointEntry:     CreateCustomerPointEntryHandler,
		UpdateCustomerPointEntry:     UpdateCustomerPointEntryHandler,
		DeleteCustomerPointEntries:   DeleteCustomerPointEntriesHandler,
		RecoveryCustomerPointEntries: RecoveryCustomerPointEntriesHandler,
		QueryCustomerPointEntry:      QueryCustomerPointEntryHandler,
		QueryCustomerPointEntries:    QueryCustomerPointEntriesHandler,

		CustomerPointEntryMember: CustomerPointEntryMemberHandler,

		CustomerPointEntrySourceOrganization: CustomerPointEntrySourceOrganizationHandler,

		CustomerPointEntryReverses: CustomerPointEntryReversesHandler,

		CustomerPointEntryReversedBy: CustomerPointEntryReversedByHandler,

		CreateCustomerCouponTemplate:    CreateCustomerCouponTemplateHandler,
		UpdateCustomerCouponTemplate:    UpdateCustomerCouponTemplateHandler,
		DeleteCustomerCouponTemplates:   DeleteCustomerCouponTemplatesHandler,
		RecoveryCustomerCouponTemplates: RecoveryCustomerCouponTemplatesHandler,
		QueryCustomerCouponTemplate:     QueryCustomerCouponTemplateHandler,
		QueryCustomerCouponTemplates:    QueryCustomerCouponTemplatesHandler,

		CustomerCouponTemplateOrganization: CustomerCouponTemplateOrganizationHandler,

		CustomerCouponTemplateApplicableStore: CustomerCouponTemplateApplicableStoreHandler,

		CustomerCouponTemplateGrants: CustomerCouponTemplateGrantsHandler,

		CustomerCouponTemplateDistributionJobs: CustomerCouponTemplateDistributionJobsHandler,

		CreateProductCategory:     CreateProductCategoryHandler,
		UpdateProductCategory:     UpdateProductCategoryHandler,
		DeleteProductCategories:   DeleteProductCategoriesHandler,
		RecoveryProductCategories: RecoveryProductCategoriesHandler,
		QueryProductCategory:      QueryProductCategoryHandler,
		QueryProductCategories:    QueryProductCategoriesHandler,

		ProductCategoryOrganization: ProductCategoryOrganizationHandler,

		ProductCategoryParent: ProductCategoryParentHandler,

		ProductCategoryChildren: ProductCategoryChildrenHandler,

		ProductCategoryProducts: ProductCategoryProductsHandler,

		CreateProductBrand:    CreateProductBrandHandler,
		UpdateProductBrand:    UpdateProductBrandHandler,
		DeleteProductBrands:   DeleteProductBrandsHandler,
		RecoveryProductBrands: RecoveryProductBrandsHandler,
		QueryProductBrand:     QueryProductBrandHandler,
		QueryProductBrands:    QueryProductBrandsHandler,

		ProductBrandOrganization: ProductBrandOrganizationHandler,

		ProductBrandProducts: ProductBrandProductsHandler,

		CreateProduct:    CreateProductHandler,
		UpdateProduct:    UpdateProductHandler,
		DeleteProducts:   DeleteProductsHandler,
		RecoveryProducts: RecoveryProductsHandler,
		QueryProduct:     QueryProductHandler,
		QueryProducts:    QueryProductsHandler,

		ProductBrand: ProductBrandHandler,

		ProductOrganization: ProductOrganizationHandler,

		ProductCategory: ProductCategoryHandler,

		ProductDefaultPackageTemplate: ProductDefaultPackageTemplateHandler,

		ProductSpecificationChoices: ProductSpecificationChoicesHandler,

		ProductSkus: ProductSkusHandler,

		CreateProductSku:    CreateProductSkuHandler,
		UpdateProductSku:    UpdateProductSkuHandler,
		DeleteProductSkus:   DeleteProductSkusHandler,
		RecoveryProductSkus: RecoveryProductSkusHandler,
		QueryProductSku:     QueryProductSkuHandler,
		QueryProductSkus:    QueryProductSkusHandler,

		ProductSkuProduct: ProductSkuProductHandler,

		ProductSkuSpecificationValues: ProductSkuSpecificationValuesHandler,

		ProductSkuPackages: ProductSkuPackagesHandler,

		ProductSkuListings: ProductSkuListingsHandler,

		CreateProductPackage:    CreateProductPackageHandler,
		UpdateProductPackage:    UpdateProductPackageHandler,
		DeleteProductPackages:   DeleteProductPackagesHandler,
		RecoveryProductPackages: RecoveryProductPackagesHandler,
		QueryProductPackage:     QueryProductPackageHandler,
		QueryProductPackages:    QueryProductPackagesHandler,

		ProductPackageSku: ProductPackageSkuHandler,

		ProductPackageTemplate: ProductPackageTemplateHandler,

		ProductPackageContainsPackage: ProductPackageContainsPackageHandler,

		ProductPackageContainedByPackages: ProductPackageContainedByPackagesHandler,

		ProductPackageOffers: ProductPackageOffersHandler,

		ProductPackageBalances: ProductPackageBalancesHandler,

		ProductPackageMovementSources: ProductPackageMovementSourcesHandler,

		ProductPackageMovementTargets: ProductPackageMovementTargetsHandler,

		ProductPackageStocktakeLines: ProductPackageStocktakeLinesHandler,

		CreateSpecificationDefinition:    CreateSpecificationDefinitionHandler,
		UpdateSpecificationDefinition:    UpdateSpecificationDefinitionHandler,
		DeleteSpecificationDefinitions:   DeleteSpecificationDefinitionsHandler,
		RecoverySpecificationDefinitions: RecoverySpecificationDefinitionsHandler,
		QuerySpecificationDefinition:     QuerySpecificationDefinitionHandler,
		QuerySpecificationDefinitions:    QuerySpecificationDefinitionsHandler,

		SpecificationDefinitionOrganization: SpecificationDefinitionOrganizationHandler,

		SpecificationDefinitionValues: SpecificationDefinitionValuesHandler,

		CreateSpecificationValue:    CreateSpecificationValueHandler,
		UpdateSpecificationValue:    UpdateSpecificationValueHandler,
		DeleteSpecificationValues:   DeleteSpecificationValuesHandler,
		RecoverySpecificationValues: RecoverySpecificationValuesHandler,
		QuerySpecificationValue:     QuerySpecificationValueHandler,
		QuerySpecificationValues:    QuerySpecificationValuesHandler,

		SpecificationValueSpecification: SpecificationValueSpecificationHandler,

		SpecificationValueProductChoices: SpecificationValueProductChoicesHandler,

		SpecificationValueSkuValues: SpecificationValueSkuValuesHandler,

		CreateProductSpecificationChoice:    CreateProductSpecificationChoiceHandler,
		UpdateProductSpecificationChoice:    UpdateProductSpecificationChoiceHandler,
		DeleteProductSpecificationChoices:   DeleteProductSpecificationChoicesHandler,
		RecoveryProductSpecificationChoices: RecoveryProductSpecificationChoicesHandler,
		QueryProductSpecificationChoice:     QueryProductSpecificationChoiceHandler,
		QueryProductSpecificationChoices:    QueryProductSpecificationChoicesHandler,

		ProductSpecificationChoiceProduct: ProductSpecificationChoiceProductHandler,

		ProductSpecificationChoiceValue: ProductSpecificationChoiceValueHandler,

		CreateProductSkuSpecificationValue:    CreateProductSkuSpecificationValueHandler,
		UpdateProductSkuSpecificationValue:    UpdateProductSkuSpecificationValueHandler,
		DeleteProductSkuSpecificationValues:   DeleteProductSkuSpecificationValuesHandler,
		RecoveryProductSkuSpecificationValues: RecoveryProductSkuSpecificationValuesHandler,
		QueryProductSkuSpecificationValue:     QueryProductSkuSpecificationValueHandler,
		QueryProductSkuSpecificationValues:    QueryProductSkuSpecificationValuesHandler,

		ProductSkuSpecificationValueSku: ProductSkuSpecificationValueSkuHandler,

		ProductSkuSpecificationValueValue: ProductSkuSpecificationValueValueHandler,

		CreateProductPackageTemplate:    CreateProductPackageTemplateHandler,
		UpdateProductPackageTemplate:    UpdateProductPackageTemplateHandler,
		DeleteProductPackageTemplates:   DeleteProductPackageTemplatesHandler,
		RecoveryProductPackageTemplates: RecoveryProductPackageTemplatesHandler,
		QueryProductPackageTemplate:     QueryProductPackageTemplateHandler,
		QueryProductPackageTemplates:    QueryProductPackageTemplatesHandler,

		ProductPackageTemplateOrganization: ProductPackageTemplateOrganizationHandler,

		ProductPackageTemplateDefaultProducts: ProductPackageTemplateDefaultProductsHandler,

		ProductPackageTemplateContainsPackage: ProductPackageTemplateContainsPackageHandler,

		ProductPackageTemplateContainedByPackages: ProductPackageTemplateContainedByPackagesHandler,

		ProductPackageTemplateCreatedPackages: ProductPackageTemplateCreatedPackagesHandler,

		CreateStoreListing:    CreateStoreListingHandler,
		UpdateStoreListing:    UpdateStoreListingHandler,
		DeleteStoreListings:   DeleteStoreListingsHandler,
		RecoveryStoreListings: RecoveryStoreListingsHandler,
		QueryStoreListing:     QueryStoreListingHandler,
		QueryStoreListings:    QueryStoreListingsHandler,

		StoreListingStore: StoreListingStoreHandler,

		StoreListingSku: StoreListingSkuHandler,

		StoreListingOffers: StoreListingOffersHandler,

		StoreListingBatches: StoreListingBatchesHandler,

		CreateStorePackageOffer:    CreateStorePackageOfferHandler,
		UpdateStorePackageOffer:    UpdateStorePackageOfferHandler,
		DeleteStorePackageOffers:   DeleteStorePackageOffersHandler,
		RecoveryStorePackageOffers: RecoveryStorePackageOffersHandler,
		QueryStorePackageOffer:     QueryStorePackageOfferHandler,
		QueryStorePackageOffers:    QueryStorePackageOffersHandler,

		StorePackageOfferListing: StorePackageOfferListingHandler,

		StorePackageOfferPackage: StorePackageOfferPackageHandler,

		StorePackageOfferPriceRevisions: StorePackageOfferPriceRevisionsHandler,

		StorePackageOfferPromotionTargets: StorePackageOfferPromotionTargetsHandler,

		CreateStorePriceRevision:    CreateStorePriceRevisionHandler,
		UpdateStorePriceRevision:    UpdateStorePriceRevisionHandler,
		DeleteStorePriceRevisions:   DeleteStorePriceRevisionsHandler,
		RecoveryStorePriceRevisions: RecoveryStorePriceRevisionsHandler,
		QueryStorePriceRevision:     QueryStorePriceRevisionHandler,
		QueryStorePriceRevisions:    QueryStorePriceRevisionsHandler,

		StorePriceRevisionOffer: StorePriceRevisionOfferHandler,

		CreateStoreInventoryBatch:     CreateStoreInventoryBatchHandler,
		UpdateStoreInventoryBatch:     UpdateStoreInventoryBatchHandler,
		DeleteStoreInventoryBatches:   DeleteStoreInventoryBatchesHandler,
		RecoveryStoreInventoryBatches: RecoveryStoreInventoryBatchesHandler,
		QueryStoreInventoryBatch:      QueryStoreInventoryBatchHandler,
		QueryStoreInventoryBatches:    QueryStoreInventoryBatchesHandler,

		StoreInventoryBatchListing: StoreInventoryBatchListingHandler,

		StoreInventoryBatchBalances: StoreInventoryBatchBalancesHandler,

		StoreInventoryBatchMovements: StoreInventoryBatchMovementsHandler,

		StoreInventoryBatchStocktakeLines: StoreInventoryBatchStocktakeLinesHandler,

		CreateStoreStockBalance:    CreateStoreStockBalanceHandler,
		UpdateStoreStockBalance:    UpdateStoreStockBalanceHandler,
		DeleteStoreStockBalances:   DeleteStoreStockBalancesHandler,
		RecoveryStoreStockBalances: RecoveryStoreStockBalancesHandler,
		QueryStoreStockBalance:     QueryStoreStockBalanceHandler,
		QueryStoreStockBalances:    QueryStoreStockBalancesHandler,

		StoreStockBalanceBatch: StoreStockBalanceBatchHandler,

		StoreStockBalancePackage: StoreStockBalancePackageHandler,

		CreateStoreStocktake:    CreateStoreStocktakeHandler,
		UpdateStoreStocktake:    UpdateStoreStocktakeHandler,
		DeleteStoreStocktakes:   DeleteStoreStocktakesHandler,
		RecoveryStoreStocktakes: RecoveryStoreStocktakesHandler,
		QueryStoreStocktake:     QueryStoreStocktakeHandler,
		QueryStoreStocktakes:    QueryStoreStocktakesHandler,

		StoreStocktakeStore: StoreStocktakeStoreHandler,

		StoreStocktakeInitiatedByAccount: StoreStocktakeInitiatedByAccountHandler,

		StoreStocktakePostedBy: StoreStocktakePostedByHandler,

		StoreStocktakeLines: StoreStocktakeLinesHandler,

		CreateStoreStocktakeLine:    CreateStoreStocktakeLineHandler,
		UpdateStoreStocktakeLine:    UpdateStoreStocktakeLineHandler,
		DeleteStoreStocktakeLines:   DeleteStoreStocktakeLinesHandler,
		RecoveryStoreStocktakeLines: RecoveryStoreStocktakeLinesHandler,
		QueryStoreStocktakeLine:     QueryStoreStocktakeLineHandler,
		QueryStoreStocktakeLines:    QueryStoreStocktakeLinesHandler,

		StoreStocktakeLineStocktake: StoreStocktakeLineStocktakeHandler,

		StoreStocktakeLineBatch: StoreStocktakeLineBatchHandler,

		StoreStocktakeLinePackage: StoreStocktakeLinePackageHandler,

		StoreStocktakeLineMovements: StoreStocktakeLineMovementsHandler,

		CreateStoreStockMovement:    CreateStoreStockMovementHandler,
		UpdateStoreStockMovement:    UpdateStoreStockMovementHandler,
		DeleteStoreStockMovements:   DeleteStoreStockMovementsHandler,
		RecoveryStoreStockMovements: RecoveryStoreStockMovementsHandler,
		QueryStoreStockMovement:     QueryStoreStockMovementHandler,
		QueryStoreStockMovements:    QueryStoreStockMovementsHandler,

		StoreStockMovementStore: StoreStockMovementStoreHandler,

		StoreStockMovementBatch: StoreStockMovementBatchHandler,

		StoreStockMovementSourcePackage: StoreStockMovementSourcePackageHandler,

		StoreStockMovementTargetPackage: StoreStockMovementTargetPackageHandler,

		StoreStockMovementStocktakeLine: StoreStockMovementStocktakeLineHandler,

		CreateStorePromotion:    CreateStorePromotionHandler,
		UpdateStorePromotion:    UpdateStorePromotionHandler,
		DeleteStorePromotions:   DeleteStorePromotionsHandler,
		RecoveryStorePromotions: RecoveryStorePromotionsHandler,
		QueryStorePromotion:     QueryStorePromotionHandler,
		QueryStorePromotions:    QueryStorePromotionsHandler,

		StorePromotionStore: StorePromotionStoreHandler,

		StorePromotionTargets: StorePromotionTargetsHandler,

		CreateStorePromotionTarget:    CreateStorePromotionTargetHandler,
		UpdateStorePromotionTarget:    UpdateStorePromotionTargetHandler,
		DeleteStorePromotionTargets:   DeleteStorePromotionTargetsHandler,
		RecoveryStorePromotionTargets: RecoveryStorePromotionTargetsHandler,
		QueryStorePromotionTarget:     QueryStorePromotionTargetHandler,
		QueryStorePromotionTargets:    QueryStorePromotionTargetsHandler,

		StorePromotionTargetPromotion: StorePromotionTargetPromotionHandler,

		StorePromotionTargetOffer: StorePromotionTargetOfferHandler,

		CreateCustomerCouponGrant:    CreateCustomerCouponGrantHandler,
		UpdateCustomerCouponGrant:    UpdateCustomerCouponGrantHandler,
		DeleteCustomerCouponGrants:   DeleteCustomerCouponGrantsHandler,
		RecoveryCustomerCouponGrants: RecoveryCustomerCouponGrantsHandler,
		QueryCustomerCouponGrant:     QueryCustomerCouponGrantHandler,
		QueryCustomerCouponGrants:    QueryCustomerCouponGrantsHandler,

		CustomerCouponGrantMember: CustomerCouponGrantMemberHandler,

		CustomerCouponGrantTemplate: CustomerCouponGrantTemplateHandler,

		CreateCustomerCouponDistributionJob:    CreateCustomerCouponDistributionJobHandler,
		UpdateCustomerCouponDistributionJob:    UpdateCustomerCouponDistributionJobHandler,
		DeleteCustomerCouponDistributionJobs:   DeleteCustomerCouponDistributionJobsHandler,
		RecoveryCustomerCouponDistributionJobs: RecoveryCustomerCouponDistributionJobsHandler,
		QueryCustomerCouponDistributionJob:     QueryCustomerCouponDistributionJobHandler,
		QueryCustomerCouponDistributionJobs:    QueryCustomerCouponDistributionJobsHandler,

		CustomerCouponDistributionJobTemplate: CustomerCouponDistributionJobTemplateHandler,

		CustomerCouponDistributionJobMember: CustomerCouponDistributionJobMemberHandler,

		WebSocket: WebSocketHandler,
	}
	return handlers
}

type GeneratedResolver struct {
	Handlers        ResolutionHandlers
	DB              *DB
	EventController *EventController
}
