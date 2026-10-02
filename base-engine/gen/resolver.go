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

		WebSocket: WebSocketHandler,
	}
	return handlers
}

type GeneratedResolver struct {
	Handlers        ResolutionHandlers
	DB              *DB
	EventController *EventController
}
