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

	MembershipsIDs               []*string
	InitializedOrganizationsIDs  []*string
	OpeningRecordsIDs            []*string
	RecordedOpeningRecordsIDs    []*string
	SessionsIDs                  []*string
	ReviewedStoresIDs            []*string
	SentMembershipInvitationsIDs []*string
	AuditLogsIDs                 []*string
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

	Memberships    []*OperatorMembership
	InitialAccount *Account
	OpeningRecords []*FranchiseOpeningRecord
	Stores         []*Store
	Roles          []*OperatorRole
	Sessions       []*Session
	AuditLogs      []*AuditLog
	PaymentConfigs []*FranchisePaymentConfig

	MembershipsIDs    []*string
	OpeningRecordsIDs []*string
	StoresIDs         []*string
	RolesIDs          []*string
	SessionsIDs       []*string
	AuditLogsIDs      []*string
	PaymentConfigsIDs []*string
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

	MembersIDs        []*string
	AuditLogsIDs      []*string
	PaymentConfigsIDs []*string
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
