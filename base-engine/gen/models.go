package gen

import (
	"fmt"
	"reflect"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/mitchellh/mapstructure"
)

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}

type AccountResultType struct {
	EntityResultType
}

type Account struct {
	ID                 string        `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Phone              string        `json:"phone" gorm:"type:varchar(32);NOT NULL;uniqueIndex;" validator:"required:true;type:phone;unique:true"`
	DisplayName        string        `json:"displayName" gorm:"type:varchar(64);NOT NULL;" validator:"required:true;minLength:1;maxLength:64"`
	Email              *string       `json:"email" gorm:"type:varchar(128);default:null;" validator:"type:email"`
	Status             AccountStatus `json:"status" gorm:"type:varchar(24);NOT NULL;index;"`
	MustChangePassword bool          `json:"mustChangePassword" gorm:"NOT NULL;default:true;"`
	CredentialVersion  int64         `json:"credentialVersion" gorm:"NOT NULL;default:1;"`
	IsDelete           *int64        `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight             *int64        `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State              *int64        `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy          *string       `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy          *string       `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy          *string       `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt          *int64        `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt          *int64        `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt          int64         `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Memberships []*OperatorMembership `json:"memberships" gorm:"foreignkey:AccountID"`

	InitializedOrganizations []*Organization `json:"initializedOrganizations" gorm:"foreignkey:InitialAccountID"`

	OpeningRecords []*FranchiseOpeningRecord `json:"openingRecords" gorm:"foreignkey:InitialAccountID"`

	RecordedOpeningRecords []*FranchiseOpeningRecord `json:"recordedOpeningRecords" gorm:"foreignkey:RecordedByAccountID"`

	Sessions []*Session `json:"sessions" gorm:"foreignkey:AccountID"`

	ReviewedStores []*Store `json:"reviewedStores" gorm:"foreignkey:ReviewedByAccountID"`

	SentMembershipInvitations []*MembershipInvitation `json:"sentMembershipInvitations" gorm:"foreignkey:InvitedByAccountID"`

	AuditLogs []*AuditLog `json:"auditLogs" gorm:"foreignkey:ActorAccountID"`

	CreatedStocktakes []*StoreStocktake `json:"createdStocktakes" gorm:"foreignkey:InitiatedByAccountID"`

	PostedStocktakes []*StoreStocktake `json:"postedStocktakes" gorm:"foreignkey:PostedByID"`
}

func (m *Account) Is_Entity() {}

type AccountChanges struct {
	ID                 string
	Phone              string
	DisplayName        string
	Email              *string
	Status             AccountStatus
	MustChangePassword bool
	CredentialVersion  int64
	IsDelete           *int64
	Weight             *int64
	State              *int64
	DeletedBy          *string
	UpdatedBy          *string
	CreatedBy          *string
	DeletedAt          *int64
	UpdatedAt          *int64
	CreatedAt          int64

	Memberships               []*OperatorMembership
	InitializedOrganizations  []*Organization
	OpeningRecords            []*FranchiseOpeningRecord
	RecordedOpeningRecords    []*FranchiseOpeningRecord
	Sessions                  []*Session
	ReviewedStores            []*Store
	SentMembershipInvitations []*MembershipInvitation
	AuditLogs                 []*AuditLog
	CreatedStocktakes         []*StoreStocktake
	PostedStocktakes          []*StoreStocktake

	MembershipsIDs               []*string
	InitializedOrganizationsIDs  []*string
	OpeningRecordsIDs            []*string
	RecordedOpeningRecordsIDs    []*string
	SessionsIDs                  []*string
	ReviewedStoresIDs            []*string
	SentMembershipInvitationsIDs []*string
	AuditLogsIDs                 []*string
	CreatedStocktakesIDs         []*string
	PostedStocktakesIDs          []*string
}

type OrganizationResultType struct {
	EntityResultType
}

type Organization struct {
	ID                   string             `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Code                 string             `json:"code" gorm:"type:varchar(32);NOT NULL;uniqueIndex;" validator:"required:true;unique:true;minLength:2;maxLength:32"`
	Name                 string             `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128"`
	Type                 OrganizationType   `json:"type" gorm:"type:varchar(24);NOT NULL;index;"`
	Status               OrganizationStatus `json:"status" gorm:"type:varchar(24);NOT NULL;index;"`
	SuspendedAt          *time.Time         `json:"suspendedAt" gorm:"default:null"`
	SuspensionReasonCode *string            `json:"suspensionReasonCode" gorm:"type:varchar(64);default:null;"`
	InitialAccountID     *string            `json:"initialAccountId" gorm:"type:varchar(36);comment:'initial_account_id';default:null;"`
	IsDelete             *int64             `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight               *int64             `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                *int64             `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy            *string            `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy            *string            `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy            *string            `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt            *int64             `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt            *int64             `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt            int64              `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Memberships []*OperatorMembership `json:"memberships" gorm:"foreignkey:OrganizationID"`

	InitialAccount *Account `json:"initialAccount"`

	OpeningRecords []*FranchiseOpeningRecord `json:"openingRecords" gorm:"foreignkey:OrganizationID"`

	Stores []*Store `json:"stores" gorm:"foreignkey:OrganizationID"`

	Roles []*OperatorRole `json:"roles" gorm:"foreignkey:OrganizationID"`

	Sessions []*Session `json:"sessions" gorm:"foreignkey:OrganizationID"`

	AuditLogs []*AuditLog `json:"auditLogs" gorm:"foreignkey:OrganizationID"`

	PaymentConfigs []*FranchisePaymentConfig `json:"paymentConfigs" gorm:"foreignkey:OrganizationID"`

	CustomerMembers []*CustomerMember `json:"customerMembers" gorm:"foreignkey:OrganizationID"`

	CustomerBenefitPolicies []*CustomerBenefitPolicy `json:"customerBenefitPolicies" gorm:"foreignkey:OrganizationID"`

	CustomerDailyPointGrantBudgets []*CustomerDailyPointGrantBudget `json:"customerDailyPointGrantBudgets" gorm:"foreignkey:OrganizationID"`

	CustomerPointEntries []*CustomerPointEntry `json:"customerPointEntries" gorm:"foreignkey:SourceOrganizationID"`

	CustomerCouponTemplates []*CustomerCouponTemplate `json:"customerCouponTemplates" gorm:"foreignkey:OrganizationID"`

	ProductCategories []*ProductCategory `json:"productCategories" gorm:"foreignkey:OrganizationID"`

	ProductBrands []*ProductBrand `json:"productBrands" gorm:"foreignkey:OrganizationID"`

	SpecificationDefinitions []*SpecificationDefinition `json:"specificationDefinitions" gorm:"foreignkey:OrganizationID"`

	ProductPackageTemplates []*ProductPackageTemplate `json:"productPackageTemplates" gorm:"foreignkey:OrganizationID"`

	Products []*Product `json:"products" gorm:"foreignkey:OrganizationID"`
}

func (m *Organization) Is_Entity() {}

type OrganizationChanges struct {
	ID                   string
	Code                 string
	Name                 string
	Type                 OrganizationType
	Status               OrganizationStatus
	SuspendedAt          *time.Time
	SuspensionReasonCode *string
	InitialAccountID     *string
	IsDelete             *int64
	Weight               *int64
	State                *int64
	DeletedBy            *string
	UpdatedBy            *string
	CreatedBy            *string
	DeletedAt            *int64
	UpdatedAt            *int64
	CreatedAt            int64

	Memberships                    []*OperatorMembership
	InitialAccount                 *Account
	OpeningRecords                 []*FranchiseOpeningRecord
	Stores                         []*Store
	Roles                          []*OperatorRole
	Sessions                       []*Session
	AuditLogs                      []*AuditLog
	PaymentConfigs                 []*FranchisePaymentConfig
	CustomerMembers                []*CustomerMember
	CustomerBenefitPolicies        []*CustomerBenefitPolicy
	CustomerDailyPointGrantBudgets []*CustomerDailyPointGrantBudget
	CustomerPointEntries           []*CustomerPointEntry
	CustomerCouponTemplates        []*CustomerCouponTemplate
	ProductCategories              []*ProductCategory
	ProductBrands                  []*ProductBrand
	SpecificationDefinitions       []*SpecificationDefinition
	ProductPackageTemplates        []*ProductPackageTemplate
	Products                       []*Product

	MembershipsIDs                    []*string
	OpeningRecordsIDs                 []*string
	StoresIDs                         []*string
	RolesIDs                          []*string
	SessionsIDs                       []*string
	AuditLogsIDs                      []*string
	PaymentConfigsIDs                 []*string
	CustomerMembersIDs                []*string
	CustomerBenefitPoliciesIDs        []*string
	CustomerDailyPointGrantBudgetsIDs []*string
	CustomerPointEntriesIDs           []*string
	CustomerCouponTemplatesIDs        []*string
	ProductCategoriesIDs              []*string
	ProductBrandsIDs                  []*string
	SpecificationDefinitionsIDs       []*string
	ProductPackageTemplatesIDs        []*string
	ProductsIDs                       []*string
}

type OperatorMembershipResultType struct {
	EntityResultType
}

type OperatorMembership struct {
	ID              string           `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Status          MembershipStatus `json:"status" gorm:"type:varchar(24);NOT NULL;index;"`
	StoreAccessMode StoreAccessMode  `json:"storeAccessMode" gorm:"type:varchar(24);NOT NULL;"`
	InvitedAt       *time.Time       `json:"invitedAt" gorm:"default:null"`
	AcceptedAt      *time.Time       `json:"acceptedAt" gorm:"default:null"`
	AccountID       string           `json:"accountId" gorm:"type:varchar(36);comment:'account_id';default:null;"`
	OrganizationID  string           `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete        *int64           `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight          *int64           `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State           *int64           `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy       *string          `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy       *string          `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy       *string          `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt       *int64           `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt       *int64           `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt       int64            `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Account *Account `json:"account"`

	Organization *Organization `json:"organization"`

	Roles []*OperatorRole `json:"roles" gorm:"many2many:operatorMembership_roles;jointable_foreignkey:member_id;association_jointable_foreignkey:role_id"`

	Stores []*Store `json:"stores" gorm:"many2many:operatorMembership_stores;jointable_foreignkey:member_id;association_jointable_foreignkey:store_id"`

	Invitations []*MembershipInvitation `json:"invitations" gorm:"foreignkey:MembershipID"`
}

func (m *OperatorMembership) Is_Entity() {}

type OperatorMembershipChanges struct {
	ID              string
	Status          MembershipStatus
	StoreAccessMode StoreAccessMode
	InvitedAt       *time.Time
	AcceptedAt      *time.Time
	AccountID       string
	OrganizationID  string
	IsDelete        *int64
	Weight          *int64
	State           *int64
	DeletedBy       *string
	UpdatedBy       *string
	CreatedBy       *string
	DeletedAt       *int64
	UpdatedAt       *int64
	CreatedAt       int64

	Account      *Account
	Organization *Organization
	Roles        []*OperatorRole
	Stores       []*Store
	Invitations  []*MembershipInvitation

	RolesIDs       []*string
	StoresIDs      []*string
	InvitationsIDs []*string
}

type PermissionResultType struct {
	EntityResultType
}

type Permission struct {
	ID        string          `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name      string          `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128"`
	Action    string          `json:"action" gorm:"type:varchar(96);NOT NULL;uniqueIndex;" validator:"required:true;unique:true;minLength:3;maxLength:96"`
	Module    string          `json:"module" gorm:"type:varchar(64);NOT NULL;index;"`
	Scope     PermissionScope `json:"scope" gorm:"type:varchar(16);NOT NULL;index;"`
	IsDelete  *int64          `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight    *int64          `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State     *int64          `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy *string         `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy *string         `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy *string         `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt *int64          `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt *int64          `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt int64           `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Roles []*OperatorRole `json:"roles" gorm:"many2many:permission_roles;jointable_foreignkey:permission_id;association_jointable_foreignkey:role_id"`
}

func (m *Permission) Is_Entity() {}

type PermissionChanges struct {
	ID        string
	Name      string
	Action    string
	Module    string
	Scope     PermissionScope
	IsDelete  *int64
	Weight    *int64
	State     *int64
	DeletedBy *string
	UpdatedBy *string
	CreatedBy *string
	DeletedAt *int64
	UpdatedAt *int64
	CreatedAt int64

	Roles []*OperatorRole

	RolesIDs []*string
}

type OperatorRoleResultType struct {
	EntityResultType
}

type OperatorRole struct {
	ID             string   `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name           string   `json:"name" gorm:"type:varchar(64);NOT NULL;" validator:"required:true;minLength:1;maxLength:64;unique:true;uniqueScope:organizationId"`
	Kind           RoleKind `json:"kind" gorm:"type:varchar(32);NOT NULL;index;"`
	OrganizationID string   `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete       *int64   `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight         *int64   `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State          *int64   `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy      *string  `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy      *string  `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy      *string  `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt      *int64   `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt      *int64   `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt      int64    `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	Members []*OperatorMembership `json:"members" gorm:"many2many:operatorMembership_roles;jointable_foreignkey:role_id;association_jointable_foreignkey:member_id"`

	Permissions []*Permission `json:"permissions" gorm:"many2many:permission_roles;jointable_foreignkey:role_id;association_jointable_foreignkey:permission_id"`
}

func (m *OperatorRole) Is_Entity() {}

type OperatorRoleChanges struct {
	ID             string
	Name           string
	Kind           RoleKind
	OrganizationID string
	IsDelete       *int64
	Weight         *int64
	State          *int64
	DeletedBy      *string
	UpdatedBy      *string
	CreatedBy      *string
	DeletedAt      *int64
	UpdatedAt      *int64
	CreatedAt      int64

	Organization *Organization
	Members      []*OperatorMembership
	Permissions  []*Permission

	MembersIDs     []*string
	PermissionsIDs []*string
}

type StoreResultType struct {
	EntityResultType
}

type Store struct {
	ID                      string               `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Code                    string               `json:"code" gorm:"type:varchar(32);NOT NULL;" validator:"required:true;minLength:1;maxLength:32;unique:true;uniqueScope:organizationId"`
	Name                    string               `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128"`
	Lifecycle               StoreLifecycle       `json:"lifecycle" gorm:"type:varchar(32);NOT NULL;index;"`
	SubmittedAt             *time.Time           `json:"submittedAt" gorm:"default:null"`
	ReviewedAt              *time.Time           `json:"reviewedAt" gorm:"default:null"`
	RejectionReason         *string              `json:"rejectionReason" gorm:"type:varchar(512);default:null;" validator:"maxLength:512"`
	ContactPhone            *string              `json:"contactPhone" gorm:"type:varchar(32);default:null;" validator:"maxLength:32"`
	ManagerName             *string              `json:"managerName" gorm:"type:varchar(64);default:null;" validator:"maxLength:64"`
	ManagerPhone            *string              `json:"managerPhone" gorm:"type:varchar(32);default:null;" validator:"maxLength:32"`
	Province                *string              `json:"province" gorm:"type:varchar(64);default:null;" validator:"maxLength:64"`
	City                    *string              `json:"city" gorm:"type:varchar(64);default:null;" validator:"maxLength:64"`
	District                *string              `json:"district" gorm:"type:varchar(64);default:null;" validator:"maxLength:64"`
	Address                 *string              `json:"address" gorm:"type:varchar(256);default:null;" validator:"maxLength:256"`
	BusinessHours           *string              `json:"businessHours" gorm:"type:varchar(64);default:null;" validator:"maxLength:64"`
	BusinessStatus          *StoreBusinessStatus `json:"businessStatus" gorm:"type:varchar(32);default:null;"`
	SupportDineIn           *bool                `json:"supportDineIn" gorm:"default:true;"`
	SupportTakeout          *bool                `json:"supportTakeout" gorm:"default:true;"`
	StoreArea               *float64             `json:"storeArea" gorm:"default:null;"`
	TableCount              *int64               `json:"tableCount" gorm:"default:null;"`
	ReceiptFooter           *string              `json:"receiptFooter" gorm:"type:varchar(256);default:null;" validator:"maxLength:256"`
	BusinessLicenseImageURL *string              `json:"businessLicenseImageUrl" gorm:"type:varchar(512);default:null;" validator:"maxLength:512"`
	OtherDocumentImageURL   *string              `json:"otherDocumentImageUrl" gorm:"type:varchar(512);default:null;" validator:"maxLength:512"`
	OrganizationID          string               `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	ReviewedByAccountID     *string              `json:"reviewedByAccountId" gorm:"type:varchar(36);comment:'reviewed_by_account_id';default:null;"`
	IsDelete                *int64               `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight                  *int64               `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                   *int64               `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy               *string              `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy               *string              `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy               *string              `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt               *int64               `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt               *int64               `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt               int64                `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	Members []*OperatorMembership `json:"members" gorm:"many2many:operatorMembership_stores;jointable_foreignkey:store_id;association_jointable_foreignkey:member_id"`

	ReviewedByAccount *Account `json:"reviewedByAccount"`

	AuditLogs []*AuditLog `json:"auditLogs" gorm:"foreignkey:StoreID"`

	PaymentConfigs []*StorePaymentConfig `json:"paymentConfigs" gorm:"foreignkey:StoreID"`

	ProductListings []*StoreListing `json:"productListings" gorm:"foreignkey:StoreID"`

	StockMovements []*StoreStockMovement `json:"stockMovements" gorm:"foreignkey:StoreID"`

	Stocktakes []*StoreStocktake `json:"stocktakes" gorm:"foreignkey:StoreID"`

	Promotions []*StorePromotion `json:"promotions" gorm:"foreignkey:StoreID"`

	CouponTemplates []*CustomerCouponTemplate `json:"couponTemplates" gorm:"foreignkey:ApplicableStoreID"`
}

func (m *Store) Is_Entity() {}

type StoreChanges struct {
	ID                      string
	Code                    string
	Name                    string
	Lifecycle               StoreLifecycle
	SubmittedAt             *time.Time
	ReviewedAt              *time.Time
	RejectionReason         *string
	ContactPhone            *string
	ManagerName             *string
	ManagerPhone            *string
	Province                *string
	City                    *string
	District                *string
	Address                 *string
	BusinessHours           *string
	BusinessStatus          *StoreBusinessStatus
	SupportDineIn           *bool
	SupportTakeout          *bool
	StoreArea               *float64
	TableCount              *int64
	ReceiptFooter           *string
	BusinessLicenseImageURL *string
	OtherDocumentImageURL   *string
	OrganizationID          string
	ReviewedByAccountID     *string
	IsDelete                *int64
	Weight                  *int64
	State                   *int64
	DeletedBy               *string
	UpdatedBy               *string
	CreatedBy               *string
	DeletedAt               *int64
	UpdatedAt               *int64
	CreatedAt               int64

	Organization      *Organization
	Members           []*OperatorMembership
	ReviewedByAccount *Account
	AuditLogs         []*AuditLog
	PaymentConfigs    []*StorePaymentConfig
	ProductListings   []*StoreListing
	StockMovements    []*StoreStockMovement
	Stocktakes        []*StoreStocktake
	Promotions        []*StorePromotion
	CouponTemplates   []*CustomerCouponTemplate

	MembersIDs         []*string
	AuditLogsIDs       []*string
	PaymentConfigsIDs  []*string
	ProductListingsIDs []*string
	StockMovementsIDs  []*string
	StocktakesIDs      []*string
	PromotionsIDs      []*string
	CouponTemplatesIDs []*string
}

type SessionResultType struct {
	EntityResultType
}

type Session struct {
	ID                string        `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	WorkspaceType     WorkspaceType `json:"workspaceType" gorm:"type:varchar(24);NOT NULL;index;"`
	CredentialVersion int64         `json:"credentialVersion" gorm:"NOT NULL;"`
	ExpiresAt         time.Time     `json:"expiresAt" gorm:"NOT NULL;index;"`
	RevokedAt         *time.Time    `json:"revokedAt" gorm:"default:null"`
	RevocationCode    *string       `json:"revocationCode" gorm:"type:varchar(64);default:null;"`
	LastSeenAt        time.Time     `json:"lastSeenAt" gorm:"NOT NULL;"`
	AccountID         string        `json:"accountId" gorm:"type:varchar(36);comment:'account_id';default:null;"`
	OrganizationID    *string       `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete          *int64        `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight            *int64        `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State             *int64        `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy         *string       `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy         *string       `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy         *string       `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt         *int64        `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt         *int64        `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt         int64         `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Account *Account `json:"account"`

	Organization *Organization `json:"organization"`

	AuditLogs []*AuditLog `json:"auditLogs" gorm:"foreignkey:SessionID"`
}

func (m *Session) Is_Entity() {}

type SessionChanges struct {
	ID                string
	WorkspaceType     WorkspaceType
	CredentialVersion int64
	ExpiresAt         time.Time
	RevokedAt         *time.Time
	RevocationCode    *string
	LastSeenAt        time.Time
	AccountID         string
	OrganizationID    *string
	IsDelete          *int64
	Weight            *int64
	State             *int64
	DeletedBy         *string
	UpdatedBy         *string
	CreatedBy         *string
	DeletedAt         *int64
	UpdatedAt         *int64
	CreatedAt         int64

	Account      *Account
	Organization *Organization
	AuditLogs    []*AuditLog

	AuditLogsIDs []*string
}

type MembershipInvitationResultType struct {
	EntityResultType
}

type MembershipInvitation struct {
	ID                 string     `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	ExpiresAt          time.Time  `json:"expiresAt" gorm:"NOT NULL;index;"`
	AcceptedAt         *time.Time `json:"acceptedAt" gorm:"default:null"`
	RevokedAt          *time.Time `json:"revokedAt" gorm:"default:null"`
	MembershipID       string     `json:"membershipId" gorm:"type:varchar(36);comment:'membership_id';default:null;"`
	InvitedByAccountID string     `json:"invitedByAccountId" gorm:"type:varchar(36);comment:'invited_by_account_id';default:null;"`
	IsDelete           *int64     `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight             *int64     `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State              *int64     `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy          *string    `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy          *string    `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy          *string    `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt          *int64     `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt          *int64     `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt          int64      `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Membership *OperatorMembership `json:"membership"`

	InvitedByAccount *Account `json:"invitedByAccount"`
}

func (m *MembershipInvitation) Is_Entity() {}

type MembershipInvitationChanges struct {
	ID                 string
	ExpiresAt          time.Time
	AcceptedAt         *time.Time
	RevokedAt          *time.Time
	MembershipID       string
	InvitedByAccountID string
	IsDelete           *int64
	Weight             *int64
	State              *int64
	DeletedBy          *string
	UpdatedBy          *string
	CreatedBy          *string
	DeletedAt          *int64
	UpdatedAt          *int64
	CreatedAt          int64

	Membership       *OperatorMembership
	InvitedByAccount *Account
}

type AuditLogResultType struct {
	EntityResultType
}

type AuditLog struct {
	ID             string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Action         string  `json:"action" gorm:"type:varchar(96);NOT NULL;index;"`
	ResourceType   string  `json:"resourceType" gorm:"type:varchar(64);NOT NULL;index;"`
	ResourceID     *string `json:"resourceId" gorm:"type:varchar(36);default:null;index;"`
	ResultCode     string  `json:"resultCode" gorm:"type:varchar(64);NOT NULL;index;"`
	MetadataJSON   *string `json:"metadataJson" gorm:"type:text;"`
	ActorAccountID *string `json:"actorAccountId" gorm:"type:varchar(36);comment:'actor_account_id';default:null;"`
	SessionID      *string `json:"sessionId" gorm:"type:varchar(36);comment:'session_id';default:null;"`
	OrganizationID *string `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	StoreID        *string `json:"storeId" gorm:"type:varchar(36);comment:'store_id';default:null;"`
	IsDelete       *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight         *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State          *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy      *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy      *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy      *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt      *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt      *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt      int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	ActorAccount *Account `json:"actorAccount"`

	Session *Session `json:"session"`

	Organization *Organization `json:"organization"`

	Store *Store `json:"store"`
}

func (m *AuditLog) Is_Entity() {}

type AuditLogChanges struct {
	ID             string
	Action         string
	ResourceType   string
	ResourceID     *string
	ResultCode     string
	MetadataJSON   *string
	ActorAccountID *string
	SessionID      *string
	OrganizationID *string
	StoreID        *string
	IsDelete       *int64
	Weight         *int64
	State          *int64
	DeletedBy      *string
	UpdatedBy      *string
	CreatedBy      *string
	DeletedAt      *int64
	UpdatedAt      *int64
	CreatedAt      int64

	ActorAccount *Account
	Session      *Session
	Organization *Organization
	Store        *Store
}

type FranchiseOpeningRecordResultType struct {
	EntityResultType
}

type FranchiseOpeningRecord struct {
	ID                  string                 `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	RecordNumber        string                 `json:"recordNumber" gorm:"type:varchar(128);NOT NULL;uniqueIndex;" validator:"required:true;unique:true;minLength:1;maxLength:128"`
	Source              FranchiseOpeningSource `json:"source" gorm:"type:varchar(32);NOT NULL;index;"`
	OrganizationID      string                 `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	InitialAccountID    string                 `json:"initialAccountId" gorm:"type:varchar(36);comment:'initial_account_id';default:null;"`
	RecordedByAccountID string                 `json:"recordedByAccountId" gorm:"type:varchar(36);comment:'recorded_by_account_id';default:null;"`
	IsDelete            *int64                 `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight              *int64                 `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State               *int64                 `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy           *string                `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy           *string                `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy           *string                `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt           *int64                 `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt           *int64                 `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt           int64                  `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	InitialAccount *Account `json:"initialAccount"`

	RecordedByAccount *Account `json:"recordedByAccount"`
}

func (m *FranchiseOpeningRecord) Is_Entity() {}

type FranchiseOpeningRecordChanges struct {
	ID                  string
	RecordNumber        string
	Source              FranchiseOpeningSource
	OrganizationID      string
	InitialAccountID    string
	RecordedByAccountID string
	IsDelete            *int64
	Weight              *int64
	State               *int64
	DeletedBy           *string
	UpdatedBy           *string
	CreatedBy           *string
	DeletedAt           *int64
	UpdatedAt           *int64
	CreatedAt           int64

	Organization      *Organization
	InitialAccount    *Account
	RecordedByAccount *Account
}

type GlobalPaymentConfigResultType struct {
	EntityResultType
}

type GlobalPaymentConfig struct {
	ID                   string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Channel              string  `json:"channel" gorm:"type:varchar(16);NOT NULL;" validator:"required:true;unique:true"`
	MerchantID           *string `json:"merchantId" gorm:"type:varchar(64);default:null;"`
	Environment          *string `json:"environment" gorm:"type:varchar(16);default:null;"`
	RatePpm              int64   `json:"ratePpm" gorm:"type:int;NOT NULL;default:0;" validator:"minValue:0;maxValue:1000000"`
	ConfigState          string  `json:"configState" gorm:"type:varchar(16);NOT NULL;"`
	Version              int64   `json:"version" gorm:"type:bigint;NOT NULL;"`
	KeyID                *string `json:"keyId" gorm:"type:varchar(64);default:null;"`
	CredentialCiphertext *string `json:"credentialCiphertext" gorm:"type:longtext;default:null;"`
	IsDelete             *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight               *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy            *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy            *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy            *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt            *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt            *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt            int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`
}

func (m *GlobalPaymentConfig) Is_Entity() {}

type GlobalPaymentConfigChanges struct {
	ID                   string
	Channel              string
	MerchantID           *string
	Environment          *string
	RatePpm              int64
	ConfigState          string
	Version              int64
	KeyID                *string
	CredentialCiphertext *string
	IsDelete             *int64
	Weight               *int64
	State                *int64
	DeletedBy            *string
	UpdatedBy            *string
	CreatedBy            *string
	DeletedAt            *int64
	UpdatedAt            *int64
	CreatedAt            int64
}

type FranchisePaymentConfigResultType struct {
	EntityResultType
}

type FranchisePaymentConfig struct {
	ID                   string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Channel              string  `json:"channel" gorm:"type:varchar(16);NOT NULL;" validator:"required:true;unique:true;uniqueScope:organizationId"`
	MerchantID           *string `json:"merchantId" gorm:"type:varchar(64);default:null;"`
	Environment          *string `json:"environment" gorm:"type:varchar(16);default:null;"`
	RatePpm              int64   `json:"ratePpm" gorm:"type:int;NOT NULL;default:0;" validator:"minValue:0;maxValue:1000000"`
	ConfigState          string  `json:"configState" gorm:"type:varchar(16);NOT NULL;"`
	Version              int64   `json:"version" gorm:"type:bigint;NOT NULL;"`
	KeyID                *string `json:"keyId" gorm:"type:varchar(64);default:null;"`
	CredentialCiphertext *string `json:"credentialCiphertext" gorm:"type:longtext;default:null;"`
	OrganizationID       string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete             *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight               *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy            *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy            *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy            *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt            *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt            *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt            int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`
}

func (m *FranchisePaymentConfig) Is_Entity() {}

type FranchisePaymentConfigChanges struct {
	ID                   string
	Channel              string
	MerchantID           *string
	Environment          *string
	RatePpm              int64
	ConfigState          string
	Version              int64
	KeyID                *string
	CredentialCiphertext *string
	OrganizationID       string
	IsDelete             *int64
	Weight               *int64
	State                *int64
	DeletedBy            *string
	UpdatedBy            *string
	CreatedBy            *string
	DeletedAt            *int64
	UpdatedAt            *int64
	CreatedAt            int64

	Organization *Organization
}

type StorePaymentConfigResultType struct {
	EntityResultType
}

type StorePaymentConfig struct {
	ID                   string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Channel              string  `json:"channel" gorm:"type:varchar(16);NOT NULL;" validator:"required:true;unique:true;uniqueScope:storeId"`
	MerchantID           *string `json:"merchantId" gorm:"type:varchar(64);default:null;"`
	Environment          *string `json:"environment" gorm:"type:varchar(16);default:null;"`
	RatePpm              int64   `json:"ratePpm" gorm:"type:int;NOT NULL;default:0;" validator:"minValue:0;maxValue:1000000"`
	ConfigState          string  `json:"configState" gorm:"type:varchar(16);NOT NULL;"`
	Version              int64   `json:"version" gorm:"type:bigint;NOT NULL;"`
	KeyID                *string `json:"keyId" gorm:"type:varchar(64);default:null;"`
	CredentialCiphertext *string `json:"credentialCiphertext" gorm:"type:longtext;default:null;"`
	StoreID              string  `json:"storeId" gorm:"type:varchar(36);comment:'store_id';default:null;"`
	IsDelete             *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight               *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy            *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy            *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy            *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt            *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt            *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt            int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Store *Store `json:"store"`
}

func (m *StorePaymentConfig) Is_Entity() {}

type StorePaymentConfigChanges struct {
	ID                   string
	Channel              string
	MerchantID           *string
	Environment          *string
	RatePpm              int64
	ConfigState          string
	Version              int64
	KeyID                *string
	CredentialCiphertext *string
	StoreID              string
	IsDelete             *int64
	Weight               *int64
	State                *int64
	DeletedBy            *string
	UpdatedBy            *string
	CreatedBy            *string
	DeletedAt            *int64
	UpdatedAt            *int64
	CreatedAt            int64

	Store *Store
}

type CustomerMemberResultType struct {
	EntityResultType
}

type CustomerMember struct {
	ID                      string               `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	MemberNumber            string               `json:"memberNumber" gorm:"type:varchar(32);NOT NULL;uniqueIndex;"`
	RequestKey              string               `json:"requestKey" gorm:"type:varchar(128);NOT NULL;"`
	Phone                   *string              `json:"phone" gorm:"type:varchar(16);default:null;" validator:"unique:true;uniqueScope:organizationId"`
	Status                  CustomerMemberStatus `json:"status" gorm:"type:varchar(24);NOT NULL;index;"`
	PointsBalance           int64                `json:"pointsBalance" gorm:"NOT NULL;default:0;"`
	PointsFrozen            bool                 `json:"pointsFrozen" gorm:"NOT NULL;default:false;"`
	NoticeVersion           string               `json:"noticeVersion" gorm:"type:varchar(64);NOT NULL;"`
	ProcessingBasisCode     string               `json:"processingBasisCode" gorm:"type:varchar(64);NOT NULL;"`
	EvidenceReference       string               `json:"evidenceReference" gorm:"type:varchar(128);NOT NULL;"`
	CancellationRequestedAt *time.Time           `json:"cancellationRequestedAt" gorm:"default:null"`
	CancelledAt             *time.Time           `json:"cancelledAt" gorm:"default:null"`
	OrganizationID          string               `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete                *int64               `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight                  *int64               `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                   *int64               `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy               *string              `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy               *string              `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy               *string              `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt               *int64               `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt               *int64               `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt               int64                `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	PointEntries []*CustomerPointEntry `json:"pointEntries" gorm:"foreignkey:MemberID"`

	CouponGrants []*CustomerCouponGrant `json:"couponGrants" gorm:"foreignkey:MemberID"`

	CouponDistributionJobs []*CustomerCouponDistributionJob `json:"couponDistributionJobs" gorm:"foreignkey:MemberID"`
}

func (m *CustomerMember) Is_Entity() {}

type CustomerMemberChanges struct {
	ID                      string
	MemberNumber            string
	RequestKey              string
	Phone                   *string
	Status                  CustomerMemberStatus
	PointsBalance           int64
	PointsFrozen            bool
	NoticeVersion           string
	ProcessingBasisCode     string
	EvidenceReference       string
	CancellationRequestedAt *time.Time
	CancelledAt             *time.Time
	OrganizationID          string
	IsDelete                *int64
	Weight                  *int64
	State                   *int64
	DeletedBy               *string
	UpdatedBy               *string
	CreatedBy               *string
	DeletedAt               *int64
	UpdatedAt               *int64
	CreatedAt               int64

	Organization           *Organization
	PointEntries           []*CustomerPointEntry
	CouponGrants           []*CustomerCouponGrant
	CouponDistributionJobs []*CustomerCouponDistributionJob

	PointEntriesIDs           []*string
	CouponGrantsIDs           []*string
	CouponDistributionJobsIDs []*string
}

type CustomerBenefitPolicyResultType struct {
	EntityResultType
}

type CustomerBenefitPolicy struct {
	ID                       string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Version                  int64   `json:"version" gorm:"NOT NULL;default:1;"`
	DiscountEnabled          bool    `json:"discountEnabled" gorm:"NOT NULL;default:false;"`
	DiscountBasisPoints      int64   `json:"discountBasisPoints" gorm:"NOT NULL;default:0;"`
	PurchaseEarnEnabled      bool    `json:"purchaseEarnEnabled" gorm:"NOT NULL;default:false;"`
	EarnAmountFen            int64   `json:"earnAmountFen" gorm:"NOT NULL;default:0;"`
	EarnPoints               int64   `json:"earnPoints" gorm:"NOT NULL;default:0;"`
	RedemptionEnabled        bool    `json:"redemptionEnabled" gorm:"NOT NULL;default:false;"`
	RedeemPoints             int64   `json:"redeemPoints" gorm:"NOT NULL;default:100;"`
	RedeemAmountFen          int64   `json:"redeemAmountFen" gorm:"NOT NULL;default:100;"`
	MaxRedemptionBasisPoints int64   `json:"maxRedemptionBasisPoints" gorm:"NOT NULL;default:0;"`
	MaxRedemptionPoints      int64   `json:"maxRedemptionPoints" gorm:"NOT NULL;default:0;"`
	ManualGrantMaxSingle     int64   `json:"manualGrantMaxSingle" gorm:"NOT NULL;default:0;"`
	ManualGrantMaxDaily      int64   `json:"manualGrantMaxDaily" gorm:"NOT NULL;default:0;"`
	PromotionWithHqCoupon    bool    `json:"promotionWithHqCoupon" gorm:"NOT NULL;default:false;"`
	PromotionWithStoreCoupon bool    `json:"promotionWithStoreCoupon" gorm:"NOT NULL;default:false;"`
	MemberPriceWithPromotion bool    `json:"memberPriceWithPromotion" gorm:"NOT NULL;default:false;"`
	PointsWithPromotion      bool    `json:"pointsWithPromotion" gorm:"NOT NULL;default:false;"`
	PointsWithCoupon         bool    `json:"pointsWithCoupon" gorm:"NOT NULL;default:false;"`
	OrganizationID           string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete                 *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight                   *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                    *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy                *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy                *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy                *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt                *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt                *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt                int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`
}

func (m *CustomerBenefitPolicy) Is_Entity() {}

type CustomerBenefitPolicyChanges struct {
	ID                       string
	Version                  int64
	DiscountEnabled          bool
	DiscountBasisPoints      int64
	PurchaseEarnEnabled      bool
	EarnAmountFen            int64
	EarnPoints               int64
	RedemptionEnabled        bool
	RedeemPoints             int64
	RedeemAmountFen          int64
	MaxRedemptionBasisPoints int64
	MaxRedemptionPoints      int64
	ManualGrantMaxSingle     int64
	ManualGrantMaxDaily      int64
	PromotionWithHqCoupon    bool
	PromotionWithStoreCoupon bool
	MemberPriceWithPromotion bool
	PointsWithPromotion      bool
	PointsWithCoupon         bool
	OrganizationID           string
	IsDelete                 *int64
	Weight                   *int64
	State                    *int64
	DeletedBy                *string
	UpdatedBy                *string
	CreatedBy                *string
	DeletedAt                *int64
	UpdatedAt                *int64
	CreatedAt                int64

	Organization *Organization
}

type CustomerDailyPointGrantBudgetResultType struct {
	EntityResultType
}

type CustomerDailyPointGrantBudget struct {
	ID             string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	BusinessDate   string  `json:"businessDate" gorm:"type:varchar(10);NOT NULL;"`
	UsedPoints     int64   `json:"usedPoints" gorm:"NOT NULL;default:0;"`
	OrganizationID string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete       *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight         *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State          *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy      *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy      *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy      *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt      *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt      *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt      int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`
}

func (m *CustomerDailyPointGrantBudget) Is_Entity() {}

type CustomerDailyPointGrantBudgetChanges struct {
	ID             string
	BusinessDate   string
	UsedPoints     int64
	OrganizationID string
	IsDelete       *int64
	Weight         *int64
	State          *int64
	DeletedBy      *string
	UpdatedBy      *string
	CreatedBy      *string
	DeletedAt      *int64
	UpdatedAt      *int64
	CreatedAt      int64

	Organization *Organization
}

type CustomerPointEntryResultType struct {
	EntityResultType
}

type CustomerPointEntry struct {
	ID                   string                     `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Delta                int64                      `json:"delta" gorm:"NOT NULL;"`
	Source               CustomerPointSource        `json:"source" gorm:"type:varchar(24);NOT NULL;"`
	OperationKind        CustomerPointOperationKind `json:"operationKind" gorm:"type:varchar(24);NOT NULL;"`
	ReasonCode           CustomerPointReasonCode    `json:"reasonCode" gorm:"type:varchar(24);NOT NULL;"`
	Note                 *string                    `json:"note" gorm:"type:varchar(512);default:null;"`
	RequestKey           string                     `json:"requestKey" gorm:"type:varchar(128);NOT NULL;"`
	MemberID             string                     `json:"memberId" gorm:"type:varchar(36);comment:'member_id';default:null;"`
	SourceOrganizationID string                     `json:"sourceOrganizationId" gorm:"type:varchar(36);comment:'source_organization_id';default:null;"`
	ReversesID           *string                    `json:"reversesId" gorm:"type:varchar(36);comment:'reverses_id';default:null;"`
	ReversedByID         *string                    `json:"reversedById" gorm:"type:varchar(36);comment:'reversed_by_id';default:null;"`
	IsDelete             *int64                     `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight               *int64                     `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                *int64                     `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy            *string                    `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy            *string                    `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy            *string                    `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt            *int64                     `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt            *int64                     `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt            int64                      `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Member *CustomerMember `json:"member"`

	SourceOrganization *Organization `json:"sourceOrganization"`

	Reverses *CustomerPointEntry `json:"reverses"`

	ReversedBy *CustomerPointEntry `json:"reversedBy"`
}

func (m *CustomerPointEntry) Is_Entity() {}

type CustomerPointEntryChanges struct {
	ID                   string
	Delta                int64
	Source               CustomerPointSource
	OperationKind        CustomerPointOperationKind
	ReasonCode           CustomerPointReasonCode
	Note                 *string
	RequestKey           string
	MemberID             string
	SourceOrganizationID string
	ReversesID           *string
	ReversedByID         *string
	IsDelete             *int64
	Weight               *int64
	State                *int64
	DeletedBy            *string
	UpdatedBy            *string
	CreatedBy            *string
	DeletedAt            *int64
	UpdatedAt            *int64
	CreatedAt            int64

	Member             *CustomerMember
	SourceOrganization *Organization
	Reverses           *CustomerPointEntry
	ReversedBy         *CustomerPointEntry
}

type CustomerCouponTemplateResultType struct {
	EntityResultType
}

type CustomerCouponTemplate struct {
	ID                  string            `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Code                string            `json:"code" gorm:"type:varchar(32);NOT NULL;"`
	RequestKey          string            `json:"requestKey" gorm:"type:varchar(128);NOT NULL;"`
	Title               string            `json:"title" gorm:"type:varchar(128);NOT NULL;"`
	AmountFen           int64             `json:"amountFen" gorm:"NOT NULL;"`
	MinSpendFen         int64             `json:"minSpendFen" gorm:"NOT NULL;default:0;"`
	DaysAfterActivation int64             `json:"daysAfterActivation" gorm:"NOT NULL;"`
	EffectiveAt         int64             `json:"effectiveAt" gorm:"type:bigint(13);NOT NULL;"`
	DistributionEndsAt  *int64            `json:"distributionEndsAt" gorm:"type:bigint(13);"`
	PerMemberLimit      int64             `json:"perMemberLimit" gorm:"NOT NULL;"`
	TotalIssueLimit     int64             `json:"totalIssueLimit" gorm:"NOT NULL;"`
	IssuedCount         int64             `json:"issuedCount" gorm:"NOT NULL;default:0;"`
	Enabled             bool              `json:"enabled" gorm:"NOT NULL;default:false;"`
	IssuerScope         CouponIssuerScope `json:"issuerScope" gorm:"type:varchar(24);NOT NULL;default:HEADQUARTERS;"`
	OrganizationID      string            `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	ApplicableStoreID   *string           `json:"applicableStoreId" gorm:"type:varchar(36);comment:'applicable_store_id';default:null;"`
	IsDelete            *int64            `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight              *int64            `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State               *int64            `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy           *string           `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy           *string           `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy           *string           `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt           *int64            `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt           *int64            `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt           int64             `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	ApplicableStore *Store `json:"applicableStore"`

	Grants []*CustomerCouponGrant `json:"grants" gorm:"foreignkey:TemplateID"`

	DistributionJobs []*CustomerCouponDistributionJob `json:"distributionJobs" gorm:"foreignkey:TemplateID"`
}

func (m *CustomerCouponTemplate) Is_Entity() {}

type CustomerCouponTemplateChanges struct {
	ID                  string
	Code                string
	RequestKey          string
	Title               string
	AmountFen           int64
	MinSpendFen         int64
	DaysAfterActivation int64
	EffectiveAt         int64
	DistributionEndsAt  *int64
	PerMemberLimit      int64
	TotalIssueLimit     int64
	IssuedCount         int64
	Enabled             bool
	IssuerScope         CouponIssuerScope
	OrganizationID      string
	ApplicableStoreID   *string
	IsDelete            *int64
	Weight              *int64
	State               *int64
	DeletedBy           *string
	UpdatedBy           *string
	CreatedBy           *string
	DeletedAt           *int64
	UpdatedAt           *int64
	CreatedAt           int64

	Organization     *Organization
	ApplicableStore  *Store
	Grants           []*CustomerCouponGrant
	DistributionJobs []*CustomerCouponDistributionJob

	GrantsIDs           []*string
	DistributionJobsIDs []*string
}

type ProductCategoryResultType struct {
	EntityResultType
}

type ProductCategory struct {
	ID             string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name           string  `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128"`
	SortOrder      int64   `json:"sortOrder" gorm:"NOT NULL;default:0;"`
	Enabled        bool    `json:"enabled" gorm:"NOT NULL;default:true;"`
	OrganizationID string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	ParentID       *string `json:"parentId" gorm:"type:varchar(36);comment:'parent_id';default:null;"`
	IsDelete       *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight         *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State          *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy      *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy      *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy      *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt      *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt      *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt      int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	Parent *ProductCategory `json:"parent"`

	Children []*ProductCategory `json:"children" gorm:"foreignkey:ParentID"`

	Products []*Product `json:"products" gorm:"foreignkey:CategoryID"`
}

func (m *ProductCategory) Is_Entity() {}

type ProductCategoryChanges struct {
	ID             string
	Name           string
	SortOrder      int64
	Enabled        bool
	OrganizationID string
	ParentID       *string
	IsDelete       *int64
	Weight         *int64
	State          *int64
	DeletedBy      *string
	UpdatedBy      *string
	CreatedBy      *string
	DeletedAt      *int64
	UpdatedAt      *int64
	CreatedAt      int64

	Organization *Organization
	Parent       *ProductCategory
	Children     []*ProductCategory
	Products     []*Product

	ChildrenIDs []*string
	ProductsIDs []*string
}

type ProductBrandResultType struct {
	EntityResultType
}

type ProductBrand struct {
	ID             string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name           string  `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128;unique:true;uniqueScope:organizationId"`
	Enabled        bool    `json:"enabled" gorm:"NOT NULL;default:true;"`
	OrganizationID string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete       *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight         *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State          *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy      *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy      *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy      *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt      *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt      *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt      int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	Products []*Product `json:"products" gorm:"foreignkey:BrandID"`
}

func (m *ProductBrand) Is_Entity() {}

type ProductBrandChanges struct {
	ID             string
	Name           string
	Enabled        bool
	OrganizationID string
	IsDelete       *int64
	Weight         *int64
	State          *int64
	DeletedBy      *string
	UpdatedBy      *string
	CreatedBy      *string
	DeletedAt      *int64
	UpdatedAt      *int64
	CreatedAt      int64

	Organization *Organization
	Products     []*Product

	ProductsIDs []*string
}

type ProductResultType struct {
	EntityResultType
}

type Product struct {
	ID                       string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name                     string  `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128"`
	Description              *string `json:"description" gorm:"type:text;"`
	ImageURL                 *string `json:"imageUrl" gorm:"type:varchar(512);default:null;" validator:"maxLength:512"`
	Enabled                  bool    `json:"enabled" gorm:"NOT NULL;default:true;"`
	BrandID                  *string `json:"brandId" gorm:"type:varchar(36);comment:'brand_id';default:null;"`
	OrganizationID           string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	CategoryID               string  `json:"categoryId" gorm:"type:varchar(36);comment:'category_id';default:null;"`
	DefaultPackageTemplateID *string `json:"defaultPackageTemplateId" gorm:"type:varchar(36);comment:'default_package_template_id';default:null;"`
	IsDelete                 *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight                   *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                    *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy                *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy                *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy                *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt                *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt                *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt                int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Brand *ProductBrand `json:"brand"`

	Organization *Organization `json:"organization"`

	Category *ProductCategory `json:"category"`

	DefaultPackageTemplate *ProductPackageTemplate `json:"defaultPackageTemplate"`

	SpecificationChoices []*ProductSpecificationChoice `json:"specificationChoices" gorm:"foreignkey:ProductID"`

	Skus []*ProductSku `json:"skus" gorm:"foreignkey:ProductID"`
}

func (m *Product) Is_Entity() {}

type ProductChanges struct {
	ID                       string
	Name                     string
	Description              *string
	ImageURL                 *string
	Enabled                  bool
	BrandID                  *string
	OrganizationID           string
	CategoryID               string
	DefaultPackageTemplateID *string
	IsDelete                 *int64
	Weight                   *int64
	State                    *int64
	DeletedBy                *string
	UpdatedBy                *string
	CreatedBy                *string
	DeletedAt                *int64
	UpdatedAt                *int64
	CreatedAt                int64

	Brand                  *ProductBrand
	Organization           *Organization
	Category               *ProductCategory
	DefaultPackageTemplate *ProductPackageTemplate
	SpecificationChoices   []*ProductSpecificationChoice
	Skus                   []*ProductSku

	SpecificationChoicesIDs []*string
	SkusIDs                 []*string
}

type ProductSkuResultType struct {
	EntityResultType
}

type ProductSku struct {
	ID                         string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name                       string  `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128"`
	PublishedPackageSetVersion int64   `json:"publishedPackageSetVersion" gorm:"NOT NULL;default:0;"`
	Ingredients                *string `json:"ingredients" gorm:"type:text;"`
	Allergens                  *string `json:"allergens" gorm:"type:text;"`
	StorageInstructions        *string `json:"storageInstructions" gorm:"type:varchar(512);default:null;" validator:"maxLength:512"`
	ShelfLifeDays              *int64  `json:"shelfLifeDays" gorm:"default:null;"`
	Enabled                    bool    `json:"enabled" gorm:"NOT NULL;default:true;"`
	SelectionRetired           bool    `json:"selectionRetired" gorm:"NOT NULL;default:false;"`
	ProductID                  string  `json:"productId" gorm:"type:varchar(36);comment:'product_id';default:null;"`
	IsDelete                   *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight                     *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                      *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy                  *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy                  *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy                  *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt                  *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt                  *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt                  int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Product *Product `json:"product"`

	SpecificationValues []*ProductSkuSpecificationValue `json:"specificationValues" gorm:"foreignkey:SkuID"`

	Packages []*ProductPackage `json:"packages" gorm:"foreignkey:SkuID"`

	Listings []*StoreListing `json:"listings" gorm:"foreignkey:SkuID"`
}

func (m *ProductSku) Is_Entity() {}

type ProductSkuChanges struct {
	ID                         string
	Name                       string
	PublishedPackageSetVersion int64
	Ingredients                *string
	Allergens                  *string
	StorageInstructions        *string
	ShelfLifeDays              *int64
	Enabled                    bool
	SelectionRetired           bool
	ProductID                  string
	IsDelete                   *int64
	Weight                     *int64
	State                      *int64
	DeletedBy                  *string
	UpdatedBy                  *string
	CreatedBy                  *string
	DeletedAt                  *int64
	UpdatedAt                  *int64
	CreatedAt                  int64

	Product             *Product
	SpecificationValues []*ProductSkuSpecificationValue
	Packages            []*ProductPackage
	Listings            []*StoreListing

	SpecificationValuesIDs []*string
	PackagesIDs            []*string
	ListingsIDs            []*string
}

type ProductPackageResultType struct {
	EntityResultType
}

type ProductPackage struct {
	ID                string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name              string  `json:"name" gorm:"type:varchar(64);NOT NULL;" validator:"required:true;minLength:1;maxLength:64"`
	Barcode           *string `json:"barcode" gorm:"type:varchar(64);default:null;" validator:"maxLength:64"`
	PackageSetVersion int64   `json:"packageSetVersion" gorm:"NOT NULL;default:1;"`
	ContainsQuantity  *int64  `json:"containsQuantity" gorm:"default:null;"`
	SuggestedPriceFen *int64  `json:"suggestedPriceFen" gorm:"default:null;"`
	Enabled           bool    `json:"enabled" gorm:"NOT NULL;default:false;"`
	SkuID             string  `json:"skuId" gorm:"type:varchar(36);comment:'sku_id';default:null;"`
	TemplateID        *string `json:"templateId" gorm:"type:varchar(36);comment:'template_id';default:null;"`
	ContainsPackageID *string `json:"containsPackageId" gorm:"type:varchar(36);comment:'contains_package_id';default:null;"`
	IsDelete          *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight            *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State             *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy         *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy         *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy         *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt         *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt         *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt         int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Sku *ProductSku `json:"sku"`

	Template *ProductPackageTemplate `json:"template"`

	ContainsPackage *ProductPackage `json:"containsPackage"`

	ContainedByPackages []*ProductPackage `json:"containedByPackages" gorm:"foreignkey:ContainsPackageID"`

	Offers []*StorePackageOffer `json:"offers" gorm:"foreignkey:PackageID"`

	Balances []*StoreStockBalance `json:"balances" gorm:"foreignkey:PackageID"`

	MovementSources []*StoreStockMovement `json:"movementSources" gorm:"foreignkey:SourcePackageID"`

	MovementTargets []*StoreStockMovement `json:"movementTargets" gorm:"foreignkey:TargetPackageID"`

	StocktakeLines []*StoreStocktakeLine `json:"stocktakeLines" gorm:"foreignkey:PackageID"`
}

func (m *ProductPackage) Is_Entity() {}

type ProductPackageChanges struct {
	ID                string
	Name              string
	Barcode           *string
	PackageSetVersion int64
	ContainsQuantity  *int64
	SuggestedPriceFen *int64
	Enabled           bool
	SkuID             string
	TemplateID        *string
	ContainsPackageID *string
	IsDelete          *int64
	Weight            *int64
	State             *int64
	DeletedBy         *string
	UpdatedBy         *string
	CreatedBy         *string
	DeletedAt         *int64
	UpdatedAt         *int64
	CreatedAt         int64

	Sku                 *ProductSku
	Template            *ProductPackageTemplate
	ContainsPackage     *ProductPackage
	ContainedByPackages []*ProductPackage
	Offers              []*StorePackageOffer
	Balances            []*StoreStockBalance
	MovementSources     []*StoreStockMovement
	MovementTargets     []*StoreStockMovement
	StocktakeLines      []*StoreStocktakeLine

	ContainedByPackagesIDs []*string
	OffersIDs              []*string
	BalancesIDs            []*string
	MovementSourcesIDs     []*string
	MovementTargetsIDs     []*string
	StocktakeLinesIDs      []*string
}

type SpecificationDefinitionResultType struct {
	EntityResultType
}

type SpecificationDefinition struct {
	ID             string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name           string  `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128;unique:true;uniqueScope:organizationId"`
	Enabled        bool    `json:"enabled" gorm:"NOT NULL;default:true;"`
	OrganizationID string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	IsDelete       *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight         *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State          *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy      *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy      *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy      *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt      *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt      *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt      int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	Values []*SpecificationValue `json:"values" gorm:"foreignkey:SpecificationID"`
}

func (m *SpecificationDefinition) Is_Entity() {}

type SpecificationDefinitionChanges struct {
	ID             string
	Name           string
	Enabled        bool
	OrganizationID string
	IsDelete       *int64
	Weight         *int64
	State          *int64
	DeletedBy      *string
	UpdatedBy      *string
	CreatedBy      *string
	DeletedAt      *int64
	UpdatedAt      *int64
	CreatedAt      int64

	Organization *Organization
	Values       []*SpecificationValue

	ValuesIDs []*string
}

type SpecificationValueResultType struct {
	EntityResultType
}

type SpecificationValue struct {
	ID              string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name            string  `json:"name" gorm:"type:varchar(128);NOT NULL;" validator:"required:true;minLength:1;maxLength:128;unique:true;uniqueScope:specificationId"`
	Enabled         bool    `json:"enabled" gorm:"NOT NULL;default:true;"`
	SpecificationID string  `json:"specificationId" gorm:"type:varchar(36);comment:'specification_id';default:null;"`
	IsDelete        *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight          *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State           *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy       *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy       *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy       *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt       *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt       *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt       int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Specification *SpecificationDefinition `json:"specification"`

	ProductChoices []*ProductSpecificationChoice `json:"productChoices" gorm:"foreignkey:ValueID"`

	SkuValues []*ProductSkuSpecificationValue `json:"skuValues" gorm:"foreignkey:ValueID"`
}

func (m *SpecificationValue) Is_Entity() {}

type SpecificationValueChanges struct {
	ID              string
	Name            string
	Enabled         bool
	SpecificationID string
	IsDelete        *int64
	Weight          *int64
	State           *int64
	DeletedBy       *string
	UpdatedBy       *string
	CreatedBy       *string
	DeletedAt       *int64
	UpdatedAt       *int64
	CreatedAt       int64

	Specification  *SpecificationDefinition
	ProductChoices []*ProductSpecificationChoice
	SkuValues      []*ProductSkuSpecificationValue

	ProductChoicesIDs []*string
	SkuValuesIDs      []*string
}

type ProductSpecificationChoiceResultType struct {
	EntityResultType
}

type ProductSpecificationChoice struct {
	ID        string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	ProductID string  `json:"productId" gorm:"type:varchar(36);comment:'product_id';default:null;"`
	ValueID   string  `json:"valueId" gorm:"type:varchar(36);comment:'value_id';default:null;"`
	IsDelete  *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight    *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State     *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Product *Product `json:"product"`

	Value *SpecificationValue `json:"value"`
}

func (m *ProductSpecificationChoice) Is_Entity() {}

type ProductSpecificationChoiceChanges struct {
	ID        string
	ProductID string
	ValueID   string
	IsDelete  *int64
	Weight    *int64
	State     *int64
	DeletedBy *string
	UpdatedBy *string
	CreatedBy *string
	DeletedAt *int64
	UpdatedAt *int64
	CreatedAt int64

	Product *Product
	Value   *SpecificationValue
}

type ProductSkuSpecificationValueResultType struct {
	EntityResultType
}

type ProductSkuSpecificationValue struct {
	ID        string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	SkuID     string  `json:"skuId" gorm:"type:varchar(36);comment:'sku_id';default:null;"`
	ValueID   string  `json:"valueId" gorm:"type:varchar(36);comment:'value_id';default:null;"`
	IsDelete  *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight    *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State     *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Sku *ProductSku `json:"sku"`

	Value *SpecificationValue `json:"value"`
}

func (m *ProductSkuSpecificationValue) Is_Entity() {}

type ProductSkuSpecificationValueChanges struct {
	ID        string
	SkuID     string
	ValueID   string
	IsDelete  *int64
	Weight    *int64
	State     *int64
	DeletedBy *string
	UpdatedBy *string
	CreatedBy *string
	DeletedAt *int64
	UpdatedAt *int64
	CreatedAt int64

	Sku   *ProductSku
	Value *SpecificationValue
}

type ProductPackageTemplateResultType struct {
	EntityResultType
}

type ProductPackageTemplate struct {
	ID                string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Name              string  `json:"name" gorm:"type:varchar(64);NOT NULL;" validator:"required:true;minLength:1;maxLength:64"`
	ContainsQuantity  *int64  `json:"containsQuantity" gorm:"default:null;"`
	Enabled           bool    `json:"enabled" gorm:"NOT NULL;default:true;"`
	OrganizationID    string  `json:"organizationId" gorm:"type:varchar(36);comment:'organization_id';default:null;"`
	ContainsPackageID *string `json:"containsPackageId" gorm:"type:varchar(36);comment:'contains_package_id';default:null;"`
	IsDelete          *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight            *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State             *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy         *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy         *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy         *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt         *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt         *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt         int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Organization *Organization `json:"organization"`

	DefaultProducts []*Product `json:"defaultProducts" gorm:"foreignkey:DefaultPackageTemplateID"`

	ContainsPackage *ProductPackageTemplate `json:"containsPackage"`

	ContainedByPackages []*ProductPackageTemplate `json:"containedByPackages" gorm:"foreignkey:ContainsPackageID"`

	CreatedPackages []*ProductPackage `json:"createdPackages" gorm:"foreignkey:TemplateID"`
}

func (m *ProductPackageTemplate) Is_Entity() {}

type ProductPackageTemplateChanges struct {
	ID                string
	Name              string
	ContainsQuantity  *int64
	Enabled           bool
	OrganizationID    string
	ContainsPackageID *string
	IsDelete          *int64
	Weight            *int64
	State             *int64
	DeletedBy         *string
	UpdatedBy         *string
	CreatedBy         *string
	DeletedAt         *int64
	UpdatedAt         *int64
	CreatedAt         int64

	Organization        *Organization
	DefaultProducts     []*Product
	ContainsPackage     *ProductPackageTemplate
	ContainedByPackages []*ProductPackageTemplate
	CreatedPackages     []*ProductPackage

	DefaultProductsIDs     []*string
	ContainedByPackagesIDs []*string
	CreatedPackagesIDs     []*string
}

type StoreListingResultType struct {
	EntityResultType
}

type StoreListing struct {
	ID         string    `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Enabled    bool      `json:"enabled" gorm:"NOT NULL;default:false;"`
	SelectedAt time.Time `json:"selectedAt" gorm:"default:null"`
	StoreID    string    `json:"storeId" gorm:"type:varchar(36);comment:'store_id';default:null;"`
	SkuID      string    `json:"skuId" gorm:"type:varchar(36);comment:'sku_id';default:null;"`
	IsDelete   *int64    `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight     *int64    `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State      *int64    `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy  *string   `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy  *string   `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy  *string   `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt  *int64    `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt  *int64    `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt  int64     `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Store *Store `json:"store"`

	Sku *ProductSku `json:"sku"`

	Offers []*StorePackageOffer `json:"offers" gorm:"foreignkey:ListingID"`

	Batches []*StoreInventoryBatch `json:"batches" gorm:"foreignkey:ListingID"`
}

func (m *StoreListing) Is_Entity() {}

type StoreListingChanges struct {
	ID         string
	Enabled    bool
	SelectedAt time.Time
	StoreID    string
	SkuID      string
	IsDelete   *int64
	Weight     *int64
	State      *int64
	DeletedBy  *string
	UpdatedBy  *string
	CreatedBy  *string
	DeletedAt  *int64
	UpdatedAt  *int64
	CreatedAt  int64

	Store   *Store
	Sku     *ProductSku
	Offers  []*StorePackageOffer
	Batches []*StoreInventoryBatch

	OffersIDs  []*string
	BatchesIDs []*string
}

type StorePackageOfferResultType struct {
	EntityResultType
}

type StorePackageOffer struct {
	ID        string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	PriceFen  int64   `json:"priceFen" gorm:"NOT NULL;default:0;"`
	Enabled   bool    `json:"enabled" gorm:"NOT NULL;default:false;"`
	ListingID string  `json:"listingId" gorm:"type:varchar(36);comment:'listing_id';default:null;"`
	PackageID string  `json:"packageId" gorm:"type:varchar(36);comment:'package_id';default:null;"`
	IsDelete  *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight    *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State     *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Listing *StoreListing `json:"listing"`

	Package *ProductPackage `json:"package"`

	PriceRevisions []*StorePriceRevision `json:"priceRevisions" gorm:"foreignkey:OfferID"`

	PromotionTargets []*StorePromotionTarget `json:"promotionTargets" gorm:"foreignkey:OfferID"`
}

func (m *StorePackageOffer) Is_Entity() {}

type StorePackageOfferChanges struct {
	ID        string
	PriceFen  int64
	Enabled   bool
	ListingID string
	PackageID string
	IsDelete  *int64
	Weight    *int64
	State     *int64
	DeletedBy *string
	UpdatedBy *string
	CreatedBy *string
	DeletedAt *int64
	UpdatedAt *int64
	CreatedAt int64

	Listing          *StoreListing
	Package          *ProductPackage
	PriceRevisions   []*StorePriceRevision
	PromotionTargets []*StorePromotionTarget

	PriceRevisionsIDs   []*string
	PromotionTargetsIDs []*string
}

type StorePriceRevisionResultType struct {
	EntityResultType
}

type StorePriceRevision struct {
	ID               string    `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	PreviousPriceFen *int64    `json:"previousPriceFen" gorm:"default:null;"`
	PriceFen         int64     `json:"priceFen" gorm:"NOT NULL;"`
	EffectiveAt      time.Time `json:"effectiveAt" gorm:"default:null"`
	ReasonCode       string    `json:"reasonCode" gorm:"type:varchar(64);NOT NULL;"`
	OfferID          string    `json:"offerId" gorm:"type:varchar(36);comment:'offer_id';default:null;"`
	IsDelete         *int64    `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight           *int64    `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State            *int64    `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy        *string   `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy        *string   `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy        *string   `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt        *int64    `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt        *int64    `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt        int64     `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Offer *StorePackageOffer `json:"offer"`
}

func (m *StorePriceRevision) Is_Entity() {}

type StorePriceRevisionChanges struct {
	ID               string
	PreviousPriceFen *int64
	PriceFen         int64
	EffectiveAt      time.Time
	ReasonCode       string
	OfferID          string
	IsDelete         *int64
	Weight           *int64
	State            *int64
	DeletedBy        *string
	UpdatedBy        *string
	CreatedBy        *string
	DeletedAt        *int64
	UpdatedAt        *int64
	CreatedAt        int64

	Offer *StorePackageOffer
}

type StoreInventoryBatchResultType struct {
	EntityResultType
}

type StoreInventoryBatch struct {
	ID              string     `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	BatchNumber     string     `json:"batchNumber" gorm:"type:varchar(128);NOT NULL;"`
	ProducedAt      *time.Time `json:"producedAt" gorm:"default:null"`
	ExpiresAt       *time.Time `json:"expiresAt" gorm:"default:null"`
	SourceReference *string    `json:"sourceReference" gorm:"type:varchar(128);default:null;" validator:"maxLength:128"`
	ListingID       string     `json:"listingId" gorm:"type:varchar(36);comment:'listing_id';default:null;"`
	IsDelete        *int64     `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight          *int64     `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State           *int64     `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy       *string    `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy       *string    `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy       *string    `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt       *int64     `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt       *int64     `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt       int64      `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Listing *StoreListing `json:"listing"`

	Balances []*StoreStockBalance `json:"balances" gorm:"foreignkey:BatchID"`

	Movements []*StoreStockMovement `json:"movements" gorm:"foreignkey:BatchID"`

	StocktakeLines []*StoreStocktakeLine `json:"stocktakeLines" gorm:"foreignkey:BatchID"`
}

func (m *StoreInventoryBatch) Is_Entity() {}

type StoreInventoryBatchChanges struct {
	ID              string
	BatchNumber     string
	ProducedAt      *time.Time
	ExpiresAt       *time.Time
	SourceReference *string
	ListingID       string
	IsDelete        *int64
	Weight          *int64
	State           *int64
	DeletedBy       *string
	UpdatedBy       *string
	CreatedBy       *string
	DeletedAt       *int64
	UpdatedAt       *int64
	CreatedAt       int64

	Listing        *StoreListing
	Balances       []*StoreStockBalance
	Movements      []*StoreStockMovement
	StocktakeLines []*StoreStocktakeLine

	BalancesIDs       []*string
	MovementsIDs      []*string
	StocktakeLinesIDs []*string
}

type StoreStockBalanceResultType struct {
	EntityResultType
}

type StoreStockBalance struct {
	ID        string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Quantity  int64   `json:"quantity" gorm:"NOT NULL;default:0;check:chk_store_stock_quantity_nonnegative,quantity >= 0;"`
	Version   int64   `json:"version" gorm:"NOT NULL;default:0;"`
	BatchID   string  `json:"batchId" gorm:"type:varchar(36);comment:'batch_id';default:null;"`
	PackageID string  `json:"packageId" gorm:"type:varchar(36);comment:'package_id';default:null;"`
	IsDelete  *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight    *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State     *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Batch *StoreInventoryBatch `json:"batch"`

	Package *ProductPackage `json:"package"`
}

func (m *StoreStockBalance) Is_Entity() {}

type StoreStockBalanceChanges struct {
	ID        string
	Quantity  int64
	Version   int64
	BatchID   string
	PackageID string
	IsDelete  *int64
	Weight    *int64
	State     *int64
	DeletedBy *string
	UpdatedBy *string
	CreatedBy *string
	DeletedAt *int64
	UpdatedAt *int64
	CreatedAt int64

	Batch   *StoreInventoryBatch
	Package *ProductPackage
}

type StoreStocktakeResultType struct {
	EntityResultType
}

type StoreStocktake struct {
	ID                   string          `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Status               StocktakeStatus `json:"status" gorm:"type:varchar(16);NOT NULL;"`
	RequestKey           string          `json:"requestKey" gorm:"type:varchar(128);NOT NULL;"`
	ScopeDigest          string          `json:"scopeDigest" gorm:"type:varchar(64);NOT NULL;"`
	StartedAt            time.Time       `json:"startedAt" gorm:"default:null"`
	ReviewedAt           *time.Time      `json:"reviewedAt" gorm:"default:null"`
	PostedAt             *time.Time      `json:"postedAt" gorm:"default:null"`
	CanceledAt           *time.Time      `json:"canceledAt" gorm:"default:null"`
	StoreID              string          `json:"storeId" gorm:"type:varchar(36);comment:'store_id';default:null;"`
	InitiatedByAccountID string          `json:"initiatedByAccountId" gorm:"type:varchar(36);comment:'initiated_by_account_id';default:null;"`
	PostedByID           *string         `json:"postedById" gorm:"type:varchar(36);comment:'posted_by_id';default:null;"`
	IsDelete             *int64          `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight               *int64          `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                *int64          `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy            *string         `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy            *string         `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy            *string         `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt            *int64          `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt            *int64          `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt            int64           `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Store *Store `json:"store"`

	InitiatedByAccount *Account `json:"initiatedByAccount"`

	PostedBy *Account `json:"postedBy"`

	Lines []*StoreStocktakeLine `json:"lines" gorm:"foreignkey:StocktakeID"`
}

func (m *StoreStocktake) Is_Entity() {}

type StoreStocktakeChanges struct {
	ID                   string
	Status               StocktakeStatus
	RequestKey           string
	ScopeDigest          string
	StartedAt            time.Time
	ReviewedAt           *time.Time
	PostedAt             *time.Time
	CanceledAt           *time.Time
	StoreID              string
	InitiatedByAccountID string
	PostedByID           *string
	IsDelete             *int64
	Weight               *int64
	State                *int64
	DeletedBy            *string
	UpdatedBy            *string
	CreatedBy            *string
	DeletedAt            *int64
	UpdatedAt            *int64
	CreatedAt            int64

	Store              *Store
	InitiatedByAccount *Account
	PostedBy           *Account
	Lines              []*StoreStocktakeLine

	LinesIDs []*string
}

type StoreStocktakeLineResultType struct {
	EntityResultType
}

type StoreStocktakeLine struct {
	ID                     string     `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	SnapshotQuantity       int64      `json:"snapshotQuantity" gorm:"NOT NULL;"`
	SnapshotVersion        int64      `json:"snapshotVersion" gorm:"NOT NULL;"`
	PackageSetVersion      int64      `json:"packageSetVersion" gorm:"NOT NULL;"`
	BatchNumberSnapshot    string     `json:"batchNumberSnapshot" gorm:"type:varchar(128);NOT NULL;"`
	ExpiresAtSnapshot      *time.Time `json:"expiresAtSnapshot" gorm:"default:null"`
	PackageNameSnapshot    string     `json:"packageNameSnapshot" gorm:"type:varchar(64);NOT NULL;"`
	PackageEnabledSnapshot bool       `json:"packageEnabledSnapshot" gorm:"NOT NULL;"`
	CountedQuantity        *int64     `json:"countedQuantity" gorm:"default:null;"`
	ReasonCode             *string    `json:"reasonCode" gorm:"type:varchar(64);default:null;"`
	ReasonNote             *string    `json:"reasonNote" gorm:"type:varchar(256);default:null;"`
	CountedAt              *time.Time `json:"countedAt" gorm:"default:null"`
	NeedsRecount           bool       `json:"needsRecount" gorm:"NOT NULL;default:false;"`
	StocktakeID            string     `json:"stocktakeId" gorm:"type:varchar(36);comment:'stocktake_id';default:null;"`
	BatchID                string     `json:"batchId" gorm:"type:varchar(36);comment:'batch_id';default:null;"`
	PackageID              string     `json:"packageId" gorm:"type:varchar(36);comment:'package_id';default:null;"`
	IsDelete               *int64     `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight                 *int64     `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                  *int64     `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy              *string    `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy              *string    `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy              *string    `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt              *int64     `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt              *int64     `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt              int64      `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Stocktake *StoreStocktake `json:"stocktake"`

	Batch *StoreInventoryBatch `json:"batch"`

	Package *ProductPackage `json:"package"`

	Movements []*StoreStockMovement `json:"movements" gorm:"foreignkey:StocktakeLineID"`
}

func (m *StoreStocktakeLine) Is_Entity() {}

type StoreStocktakeLineChanges struct {
	ID                     string
	SnapshotQuantity       int64
	SnapshotVersion        int64
	PackageSetVersion      int64
	BatchNumberSnapshot    string
	ExpiresAtSnapshot      *time.Time
	PackageNameSnapshot    string
	PackageEnabledSnapshot bool
	CountedQuantity        *int64
	ReasonCode             *string
	ReasonNote             *string
	CountedAt              *time.Time
	NeedsRecount           bool
	StocktakeID            string
	BatchID                string
	PackageID              string
	IsDelete               *int64
	Weight                 *int64
	State                  *int64
	DeletedBy              *string
	UpdatedBy              *string
	CreatedBy              *string
	DeletedAt              *int64
	UpdatedAt              *int64
	CreatedAt              int64

	Stocktake *StoreStocktake
	Batch     *StoreInventoryBatch
	Package   *ProductPackage
	Movements []*StoreStockMovement

	MovementsIDs []*string
}

type StoreStockMovementResultType struct {
	EntityResultType
}

type StoreStockMovement struct {
	ID                string            `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Kind              StockMovementKind `json:"kind" gorm:"type:varchar(32);NOT NULL;"`
	RequestKey        string            `json:"requestKey" gorm:"type:varchar(128);NOT NULL;"`
	SourceQuantity    *int64            `json:"sourceQuantity" gorm:"default:null;"`
	TargetQuantity    *int64            `json:"targetQuantity" gorm:"default:null;"`
	FactorSnapshot    *int64            `json:"factorSnapshot" gorm:"default:null;"`
	PackageSetVersion int64             `json:"packageSetVersion" gorm:"NOT NULL;"`
	ReasonCode        *string           `json:"reasonCode" gorm:"type:varchar(64);default:null;"`
	OccurredAt        time.Time         `json:"occurredAt" gorm:"default:null"`
	StoreID           string            `json:"storeId" gorm:"type:varchar(36);comment:'store_id';default:null;"`
	BatchID           string            `json:"batchId" gorm:"type:varchar(36);comment:'batch_id';default:null;"`
	SourcePackageID   *string           `json:"sourcePackageId" gorm:"type:varchar(36);comment:'source_package_id';default:null;"`
	TargetPackageID   *string           `json:"targetPackageId" gorm:"type:varchar(36);comment:'target_package_id';default:null;"`
	StocktakeLineID   *string           `json:"stocktakeLineId" gorm:"type:varchar(36);comment:'stocktake_line_id';default:null;"`
	IsDelete          *int64            `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight            *int64            `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State             *int64            `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy         *string           `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy         *string           `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy         *string           `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt         *int64            `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt         *int64            `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt         int64             `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Store *Store `json:"store"`

	Batch *StoreInventoryBatch `json:"batch"`

	SourcePackage *ProductPackage `json:"sourcePackage"`

	TargetPackage *ProductPackage `json:"targetPackage"`

	StocktakeLine *StoreStocktakeLine `json:"stocktakeLine"`
}

func (m *StoreStockMovement) Is_Entity() {}

type StoreStockMovementChanges struct {
	ID                string
	Kind              StockMovementKind
	RequestKey        string
	SourceQuantity    *int64
	TargetQuantity    *int64
	FactorSnapshot    *int64
	PackageSetVersion int64
	ReasonCode        *string
	OccurredAt        time.Time
	StoreID           string
	BatchID           string
	SourcePackageID   *string
	TargetPackageID   *string
	StocktakeLineID   *string
	IsDelete          *int64
	Weight            *int64
	State             *int64
	DeletedBy         *string
	UpdatedBy         *string
	CreatedBy         *string
	DeletedAt         *int64
	UpdatedAt         *int64
	CreatedAt         int64

	Store         *Store
	Batch         *StoreInventoryBatch
	SourcePackage *ProductPackage
	TargetPackage *ProductPackage
	StocktakeLine *StoreStocktakeLine
}

type StorePromotionResultType struct {
	EntityResultType
}

type StorePromotion struct {
	ID                   string             `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	RuleKey              string             `json:"ruleKey" gorm:"type:varchar(36);NOT NULL;"`
	TimeZone             string             `json:"timeZone" gorm:"type:varchar(64);NOT NULL;"`
	Kind                 StorePromotionKind `json:"kind" gorm:"type:varchar(32);NOT NULL;"`
	Version              int64              `json:"version" gorm:"NOT NULL;default:1;"`
	Enabled              bool               `json:"enabled" gorm:"NOT NULL;default:false;"`
	StartsAt             time.Time          `json:"startsAt" gorm:"default:null"`
	EndsAt               time.Time          `json:"endsAt" gorm:"default:null"`
	ThresholdFen         *int64             `json:"thresholdFen" gorm:"default:null;"`
	ThresholdQuantity    *int64             `json:"thresholdQuantity" gorm:"default:null;"`
	DiscountFen          *int64             `json:"discountFen" gorm:"default:null;"`
	DiscountBasisPoints  *int64             `json:"discountBasisPoints" gorm:"default:null;"`
	FixedPriceFen        *int64             `json:"fixedPriceFen" gorm:"default:null;"`
	StackWithHqCoupon    bool               `json:"stackWithHqCoupon" gorm:"NOT NULL;default:false;"`
	StackWithStoreCoupon bool               `json:"stackWithStoreCoupon" gorm:"NOT NULL;default:false;"`
	StackWithMemberPrice bool               `json:"stackWithMemberPrice" gorm:"NOT NULL;default:false;"`
	StoreID              string             `json:"storeId" gorm:"type:varchar(36);comment:'store_id';default:null;"`
	IsDelete             *int64             `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight               *int64             `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State                *int64             `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy            *string            `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy            *string            `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy            *string            `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt            *int64             `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt            *int64             `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt            int64              `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Store *Store `json:"store"`

	Targets []*StorePromotionTarget `json:"targets" gorm:"foreignkey:PromotionID"`
}

func (m *StorePromotion) Is_Entity() {}

type StorePromotionChanges struct {
	ID                   string
	RuleKey              string
	TimeZone             string
	Kind                 StorePromotionKind
	Version              int64
	Enabled              bool
	StartsAt             time.Time
	EndsAt               time.Time
	ThresholdFen         *int64
	ThresholdQuantity    *int64
	DiscountFen          *int64
	DiscountBasisPoints  *int64
	FixedPriceFen        *int64
	StackWithHqCoupon    bool
	StackWithStoreCoupon bool
	StackWithMemberPrice bool
	StoreID              string
	IsDelete             *int64
	Weight               *int64
	State                *int64
	DeletedBy            *string
	UpdatedBy            *string
	CreatedBy            *string
	DeletedAt            *int64
	UpdatedAt            *int64
	CreatedAt            int64

	Store   *Store
	Targets []*StorePromotionTarget

	TargetsIDs []*string
}

type StorePromotionTargetResultType struct {
	EntityResultType
}

type StorePromotionTarget struct {
	ID               string  `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	RequiredQuantity int64   `json:"requiredQuantity" gorm:"NOT NULL;default:1;"`
	PromotionID      string  `json:"promotionId" gorm:"type:varchar(36);comment:'promotion_id';default:null;"`
	OfferID          string  `json:"offerId" gorm:"type:varchar(36);comment:'offer_id';default:null;"`
	IsDelete         *int64  `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight           *int64  `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State            *int64  `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy        *string `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy        *string `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy        *string `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt        *int64  `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt        *int64  `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt        int64   `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Promotion *StorePromotion `json:"promotion"`

	Offer *StorePackageOffer `json:"offer"`
}

func (m *StorePromotionTarget) Is_Entity() {}

type StorePromotionTargetChanges struct {
	ID               string
	RequiredQuantity int64
	PromotionID      string
	OfferID          string
	IsDelete         *int64
	Weight           *int64
	State            *int64
	DeletedBy        *string
	UpdatedBy        *string
	CreatedBy        *string
	DeletedAt        *int64
	UpdatedAt        *int64
	CreatedAt        int64

	Promotion *StorePromotion
	Offer     *StorePackageOffer
}

type CustomerCouponGrantResultType struct {
	EntityResultType
}

type CustomerCouponGrant struct {
	ID                  string                    `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Status              CustomerCouponGrantStatus `json:"status" gorm:"type:varchar(24);NOT NULL;index;"`
	AmountFen           int64                     `json:"amountFen" gorm:"NOT NULL;"`
	MinSpendFen         int64                     `json:"minSpendFen" gorm:"NOT NULL;"`
	DaysAfterActivation int64                     `json:"daysAfterActivation" gorm:"NOT NULL;"`
	IssuedAt            int64                     `json:"issuedAt" gorm:"type:bigint(13);NOT NULL;"`
	ActivatedAt         *int64                    `json:"activatedAt" gorm:"type:bigint(13);"`
	ExpiresAt           *int64                    `json:"expiresAt" gorm:"type:bigint(13);"`
	RevokedAt           *int64                    `json:"revokedAt" gorm:"type:bigint(13);"`
	RequestKey          string                    `json:"requestKey" gorm:"type:varchar(128);NOT NULL;"`
	IssuerRequestDigest *string                   `json:"issuerRequestDigest" gorm:"type:varchar(64);default:null;"`
	MemberID            string                    `json:"memberId" gorm:"type:varchar(36);comment:'member_id';default:null;"`
	TemplateID          string                    `json:"templateId" gorm:"type:varchar(36);comment:'template_id';default:null;"`
	IsDelete            *int64                    `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight              *int64                    `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State               *int64                    `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy           *string                   `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy           *string                   `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy           *string                   `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt           *int64                    `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt           *int64                    `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt           int64                     `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Member *CustomerMember `json:"member"`

	Template *CustomerCouponTemplate `json:"template"`
}

func (m *CustomerCouponGrant) Is_Entity() {}

type CustomerCouponGrantChanges struct {
	ID                  string
	Status              CustomerCouponGrantStatus
	AmountFen           int64
	MinSpendFen         int64
	DaysAfterActivation int64
	IssuedAt            int64
	ActivatedAt         *int64
	ExpiresAt           *int64
	RevokedAt           *int64
	RequestKey          string
	IssuerRequestDigest *string
	MemberID            string
	TemplateID          string
	IsDelete            *int64
	Weight              *int64
	State               *int64
	DeletedBy           *string
	UpdatedBy           *string
	CreatedBy           *string
	DeletedAt           *int64
	UpdatedAt           *int64
	CreatedAt           int64

	Member   *CustomerMember
	Template *CustomerCouponTemplate
}

type CustomerCouponDistributionJobResultType struct {
	EntityResultType
}

type CustomerCouponDistributionJob struct {
	ID              string                              `json:"id" gorm:"type:varchar(36);comment:'uuid';primaryKey;uniqueIndex;NOT NULL;"`
	Kind            CustomerCouponDistributionJobKind   `json:"kind" gorm:"type:varchar(24);NOT NULL;"`
	Status          CustomerCouponDistributionJobStatus `json:"status" gorm:"type:varchar(24);NOT NULL;"`
	RequestKey      string                              `json:"requestKey" gorm:"type:varchar(128);NOT NULL;"`
	AvailableAt     int64                               `json:"availableAt" gorm:"type:bigint(13);NOT NULL;"`
	LeaseExpiresAt  *int64                              `json:"leaseExpiresAt" gorm:"type:bigint(13);"`
	LeaseToken      *string                             `json:"leaseToken" gorm:"type:varchar(64);default:null;"`
	CursorCreatedAt *int64                              `json:"cursorCreatedAt" gorm:"type:bigint(13);"`
	CursorKey       *string                             `json:"cursorKey" gorm:"type:varchar(64);default:null;"`
	Attempts        int64                               `json:"attempts" gorm:"NOT NULL;default:0;"`
	LastErrorCode   *string                             `json:"lastErrorCode" gorm:"type:varchar(64);default:null;"`
	TemplateID      *string                             `json:"templateId" gorm:"type:varchar(36);comment:'template_id';default:null;"`
	MemberID        *string                             `json:"memberId" gorm:"type:varchar(36);comment:'member_id';default:null;"`
	IsDelete        *int64                              `json:"isDelete" gorm:"type:int(2);comment:'是否删除：1/正常、2/删除';default:1;index:is_delete;"`
	Weight          *int64                              `json:"weight" gorm:"type:int(11);comment:'权重：用来排序';default:1;index:weight;"`
	State           *int64                              `json:"state" gorm:"type:int(2);comment:'状态：1/正常、2/禁用';default:1;index:state;"`
	DeletedBy       *string                             `json:"deletedBy" gorm:"type:varchar(36);comment:'deleted_by';default:null;index:deleted_by;"`
	UpdatedBy       *string                             `json:"updatedBy" gorm:"type:varchar(36);comment:'updated_by';default:null;index:updated_by;"`
	CreatedBy       *string                             `json:"createdBy" gorm:"type:varchar(36);comment:'created_by';default:null;index:created_by;"`
	DeletedAt       *int64                              `json:"deletedAt" gorm:"type:bigint(13);comment:'deleted_at';default:null;"`
	UpdatedAt       *int64                              `json:"updatedAt" gorm:"type:bigint(13);comment:'updated_at';default:null; autoUpdateTime:milli;"`
	CreatedAt       int64                               `json:"createdAt" gorm:"type:bigint(13);comment:'created_at';default:null; autoCreateTime:milli;"`

	Template *CustomerCouponTemplate `json:"template"`

	Member *CustomerMember `json:"member"`
}

func (m *CustomerCouponDistributionJob) Is_Entity() {}

type CustomerCouponDistributionJobChanges struct {
	ID              string
	Kind            CustomerCouponDistributionJobKind
	Status          CustomerCouponDistributionJobStatus
	RequestKey      string
	AvailableAt     int64
	LeaseExpiresAt  *int64
	LeaseToken      *string
	CursorCreatedAt *int64
	CursorKey       *string
	Attempts        int64
	LastErrorCode   *string
	TemplateID      *string
	MemberID        *string
	IsDelete        *int64
	Weight          *int64
	State           *int64
	DeletedBy       *string
	UpdatedBy       *string
	CreatedBy       *string
	DeletedAt       *int64
	UpdatedAt       *int64
	CreatedAt       int64

	Template *CustomerCouponTemplate
	Member   *CustomerMember
}

// used to convert map[string]interface{} to EntityChanges struct
func ApplyChanges(changes map[string]interface{}, to interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		TagName:     "json",
		Result:      to,
		ZeroFields:  true,
		// This is needed to get mapstructure to call the gqlgen unmarshaler func for custom scalars (eg Date)
		DecodeHook: func(a reflect.Type, b reflect.Type, v interface{}) (interface{}, error) {
			if a == b {
				return v, nil
			}

			if b == reflect.TypeOf(time.Time{}) {
				switch a.Kind() {
				case reflect.String:
					return time.Parse(time.RFC3339, v.(string))
				case reflect.Float64:
					return time.Unix(0, int64(v.(float64))*int64(time.Millisecond)), nil
				case reflect.Int64:
					return time.Unix(0, v.(int64)*int64(time.Millisecond)), nil
				case reflect.Struct:
					if t, ok := v.(time.Time); ok {
						return t, nil
					}
					return v, nil
				default:
					return v, fmt.Errorf("Unable to parse date from %v", v)
				}
			}

			if reflect.PtrTo(b).Implements(reflect.TypeOf((*graphql.Unmarshaler)(nil)).Elem()) {
				resultType := reflect.New(b)
				result := resultType.MethodByName("UnmarshalGQL").Call([]reflect.Value{reflect.ValueOf(v)})
				err, _ := result[0].Interface().(error)
				return resultType.Elem().Interface(), err
			}

			return v, nil
		},
	})

	if err != nil {
		return err
	}

	return dec.Decode(changes)
}
