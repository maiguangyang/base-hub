package gen

import (
	"context"
	"errors"
	"math"

	"base-engine/auth"

	"github.com/99designs/gqlgen/graphql"
	"github.com/graph-gophers/dataloader"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/gorm"
)

type GeneratedQueryResolver struct{ *GeneratedResolver }

type QueryAccountHandlerOptions struct {
	ID     *string
	Filter *AccountFilterType
}

func (r *GeneratedQueryResolver) Account(ctx context.Context, id *string, filter *AccountFilterType) (*Account, error) {
	opts := QueryAccountHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryAccount(ctx, r.GeneratedResolver, opts)
}
func QueryAccountHandler(ctx context.Context, r *GeneratedResolver, opts QueryAccountHandlerOptions) (*Account, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := AccountQueryFilter{}
	rt := &AccountResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("accounts", ctx)+".id = ?", *opts.ID)
	}

	var items []*Account
	giOpts := GetItemsOptions{
		Alias:      TableName("accounts", ctx),
		Preloaders: []string{},
		Item:       &Account{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "Account"}
	}
	return items[0], err
}

type QueryAccountsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*AccountSortType
	Filter      *AccountFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) Accounts(ctx context.Context, current_page *int, per_page *int, q *string, sort []*AccountSortType, filter *AccountFilterType, rand *bool) (*AccountResultType, error) {
	opts := QueryAccountsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryAccounts(ctx, r.GeneratedResolver, opts)
}
func QueryAccountsHandler(ctx context.Context, r *GeneratedResolver, opts QueryAccountsHandlerOptions) (*AccountResultType, error) {
	query := AccountQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &AccountResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedAccountResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedAccountResultTypeResolver) Data(ctx context.Context, obj *AccountResultType) (items []*Account, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("accounts", ctx),
		Preloaders: []string{},
		Item:       &Account{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*Account{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedAccountResultTypeResolver) Total(ctx context.Context, obj *AccountResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("accounts", ctx), &Account{})
}

func (r *GeneratedAccountResultTypeResolver) TotalPage(ctx context.Context, obj *AccountResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedAccountResultTypeResolver) CurrentPage(ctx context.Context, obj *AccountResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedAccountResultTypeResolver) PerPage(ctx context.Context, obj *AccountResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedAccountResolver struct{ *GeneratedResolver }

func (r *GeneratedAccountResolver) Memberships(ctx context.Context, obj *Account) (res []*OperatorMembership, err error) {
	return r.Handlers.AccountMemberships(ctx, r.GeneratedResolver, obj)
}
func AccountMembershipsHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*OperatorMembership, err error) {

	items = []*OperatorMembership{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Memberships"); err != nil {
		return items, errors.New("Memberships " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OperatorMembershipAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*OperatorMembership{}
	if item != nil {
		items = item.([]*OperatorMembership)
	}

	return
}

func (r *GeneratedAccountResolver) MembershipsIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*OperatorMembership{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["AccountAndOperatorMembershipIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*OperatorMembership)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) InitializedOrganizations(ctx context.Context, obj *Account) (res []*Organization, err error) {
	return r.Handlers.AccountInitializedOrganizations(ctx, r.GeneratedResolver, obj)
}
func AccountInitializedOrganizationsHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*Organization, err error) {

	items = []*Organization{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "InitializedOrganizations"); err != nil {
		return items, errors.New("InitializedOrganizations " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationInitialAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Organization{}
	if item != nil {
		items = item.([]*Organization)
	}

	return
}

func (r *GeneratedAccountResolver) InitializedOrganizationsIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*Organization{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["InitialAccountAndOrganizationIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Organization)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) OpeningRecords(ctx context.Context, obj *Account) (res []*FranchiseOpeningRecord, err error) {
	return r.Handlers.AccountOpeningRecords(ctx, r.GeneratedResolver, obj)
}
func AccountOpeningRecordsHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*FranchiseOpeningRecord, err error) {

	items = []*FranchiseOpeningRecord{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "OpeningRecords"); err != nil {
		return items, errors.New("OpeningRecords " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["FranchiseOpeningRecordInitialAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*FranchiseOpeningRecord{}
	if item != nil {
		items = item.([]*FranchiseOpeningRecord)
	}

	return
}

func (r *GeneratedAccountResolver) OpeningRecordsIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*FranchiseOpeningRecord{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["InitialAccountAndFranchiseOpeningRecordIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*FranchiseOpeningRecord)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) RecordedOpeningRecords(ctx context.Context, obj *Account) (res []*FranchiseOpeningRecord, err error) {
	return r.Handlers.AccountRecordedOpeningRecords(ctx, r.GeneratedResolver, obj)
}
func AccountRecordedOpeningRecordsHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*FranchiseOpeningRecord, err error) {

	items = []*FranchiseOpeningRecord{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "RecordedOpeningRecords"); err != nil {
		return items, errors.New("RecordedOpeningRecords " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["FranchiseOpeningRecordRecordedByAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*FranchiseOpeningRecord{}
	if item != nil {
		items = item.([]*FranchiseOpeningRecord)
	}

	return
}

func (r *GeneratedAccountResolver) RecordedOpeningRecordsIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*FranchiseOpeningRecord{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["RecordedByAccountAndFranchiseOpeningRecordIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*FranchiseOpeningRecord)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) Sessions(ctx context.Context, obj *Account) (res []*Session, err error) {
	return r.Handlers.AccountSessions(ctx, r.GeneratedResolver, obj)
}
func AccountSessionsHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*Session, err error) {

	items = []*Session{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Sessions"); err != nil {
		return items, errors.New("Sessions " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SessionAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Session{}
	if item != nil {
		items = item.([]*Session)
	}

	return
}

func (r *GeneratedAccountResolver) SessionsIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*Session{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["AccountAndSessionIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Session)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) ReviewedStores(ctx context.Context, obj *Account) (res []*Store, err error) {
	return r.Handlers.AccountReviewedStores(ctx, r.GeneratedResolver, obj)
}
func AccountReviewedStoresHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*Store, err error) {

	items = []*Store{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ReviewedStores"); err != nil {
		return items, errors.New("ReviewedStores " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreReviewedByAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Store{}
	if item != nil {
		items = item.([]*Store)
	}

	return
}

func (r *GeneratedAccountResolver) ReviewedStoresIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*Store{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ReviewedByAccountAndStoreIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Store)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) SentMembershipInvitations(ctx context.Context, obj *Account) (res []*MembershipInvitation, err error) {
	return r.Handlers.AccountSentMembershipInvitations(ctx, r.GeneratedResolver, obj)
}
func AccountSentMembershipInvitationsHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*MembershipInvitation, err error) {

	items = []*MembershipInvitation{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SentMembershipInvitations"); err != nil {
		return items, errors.New("SentMembershipInvitations " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["MembershipInvitationInvitedByAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*MembershipInvitation{}
	if item != nil {
		items = item.([]*MembershipInvitation)
	}

	return
}

func (r *GeneratedAccountResolver) SentMembershipInvitationsIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*MembershipInvitation{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["InvitedByAccountAndMembershipInvitationIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*MembershipInvitation)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) AuditLogs(ctx context.Context, obj *Account) (res []*AuditLog, err error) {
	return r.Handlers.AccountAuditLogs(ctx, r.GeneratedResolver, obj)
}
func AccountAuditLogsHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*AuditLog, err error) {

	items = []*AuditLog{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "AuditLogs"); err != nil {
		return items, errors.New("AuditLogs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["AuditLogActorAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*AuditLog{}
	if item != nil {
		items = item.([]*AuditLog)
	}

	return
}

func (r *GeneratedAccountResolver) AuditLogsIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*AuditLog{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ActorAccountAndAuditLogIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*AuditLog)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryOrganizationHandlerOptions struct {
	ID     *string
	Filter *OrganizationFilterType
}

func (r *GeneratedQueryResolver) Organization(ctx context.Context, id *string, filter *OrganizationFilterType) (*Organization, error) {
	opts := QueryOrganizationHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryOrganization(ctx, r.GeneratedResolver, opts)
}
func QueryOrganizationHandler(ctx context.Context, r *GeneratedResolver, opts QueryOrganizationHandlerOptions) (*Organization, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := OrganizationQueryFilter{}
	rt := &OrganizationResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("organizations", ctx)+".id = ?", *opts.ID)
	}

	var items []*Organization
	giOpts := GetItemsOptions{
		Alias:      TableName("organizations", ctx),
		Preloaders: []string{},
		Item:       &Organization{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "Organization"}
	}
	return items[0], err
}

type QueryOrganizationsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*OrganizationSortType
	Filter      *OrganizationFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) Organizations(ctx context.Context, current_page *int, per_page *int, q *string, sort []*OrganizationSortType, filter *OrganizationFilterType, rand *bool) (*OrganizationResultType, error) {
	opts := QueryOrganizationsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryOrganizations(ctx, r.GeneratedResolver, opts)
}
func QueryOrganizationsHandler(ctx context.Context, r *GeneratedResolver, opts QueryOrganizationsHandlerOptions) (*OrganizationResultType, error) {
	query := OrganizationQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &OrganizationResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedOrganizationResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedOrganizationResultTypeResolver) Data(ctx context.Context, obj *OrganizationResultType) (items []*Organization, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("organizations", ctx),
		Preloaders: []string{},
		Item:       &Organization{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*Organization{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedOrganizationResultTypeResolver) Total(ctx context.Context, obj *OrganizationResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("organizations", ctx), &Organization{})
}

func (r *GeneratedOrganizationResultTypeResolver) TotalPage(ctx context.Context, obj *OrganizationResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedOrganizationResultTypeResolver) CurrentPage(ctx context.Context, obj *OrganizationResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedOrganizationResultTypeResolver) PerPage(ctx context.Context, obj *OrganizationResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedOrganizationResolver struct{ *GeneratedResolver }

func (r *GeneratedOrganizationResolver) Memberships(ctx context.Context, obj *Organization) (res []*OperatorMembership, err error) {
	return r.Handlers.OrganizationMemberships(ctx, r.GeneratedResolver, obj)
}
func OrganizationMembershipsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*OperatorMembership, err error) {

	items = []*OperatorMembership{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Memberships"); err != nil {
		return items, errors.New("Memberships " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OperatorMembershipOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*OperatorMembership{}
	if item != nil {
		items = item.([]*OperatorMembership)
	}

	return
}

func (r *GeneratedOrganizationResolver) MembershipsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*OperatorMembership{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndOperatorMembershipIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*OperatorMembership)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) InitialAccount(ctx context.Context, obj *Organization) (res *Account, err error) {
	return r.Handlers.OrganizationInitialAccount(ctx, r.GeneratedResolver, obj)
}
func OrganizationInitialAccountHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.InitialAccountID

	if objKey != nil {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Account)

	}

	return
}

func (r *GeneratedOrganizationResolver) OpeningRecords(ctx context.Context, obj *Organization) (res []*FranchiseOpeningRecord, err error) {
	return r.Handlers.OrganizationOpeningRecords(ctx, r.GeneratedResolver, obj)
}
func OrganizationOpeningRecordsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*FranchiseOpeningRecord, err error) {

	items = []*FranchiseOpeningRecord{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "OpeningRecords"); err != nil {
		return items, errors.New("OpeningRecords " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["FranchiseOpeningRecordOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*FranchiseOpeningRecord{}
	if item != nil {
		items = item.([]*FranchiseOpeningRecord)
	}

	return
}

func (r *GeneratedOrganizationResolver) OpeningRecordsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*FranchiseOpeningRecord{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndFranchiseOpeningRecordIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*FranchiseOpeningRecord)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) Stores(ctx context.Context, obj *Organization) (res []*Store, err error) {
	return r.Handlers.OrganizationStores(ctx, r.GeneratedResolver, obj)
}
func OrganizationStoresHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*Store, err error) {

	items = []*Store{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Stores"); err != nil {
		return items, errors.New("Stores " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Store{}
	if item != nil {
		items = item.([]*Store)
	}

	return
}

func (r *GeneratedOrganizationResolver) StoresIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*Store{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndStoreIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Store)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) Roles(ctx context.Context, obj *Organization) (res []*OperatorRole, err error) {
	return r.Handlers.OrganizationRoles(ctx, r.GeneratedResolver, obj)
}
func OrganizationRolesHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*OperatorRole, err error) {

	items = []*OperatorRole{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Roles"); err != nil {
		return items, errors.New("Roles " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OperatorRoleOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*OperatorRole{}
	if item != nil {
		items = item.([]*OperatorRole)
	}

	return
}

func (r *GeneratedOrganizationResolver) RolesIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*OperatorRole{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndOperatorRoleIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*OperatorRole)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) Sessions(ctx context.Context, obj *Organization) (res []*Session, err error) {
	return r.Handlers.OrganizationSessions(ctx, r.GeneratedResolver, obj)
}
func OrganizationSessionsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*Session, err error) {

	items = []*Session{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Sessions"); err != nil {
		return items, errors.New("Sessions " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SessionOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Session{}
	if item != nil {
		items = item.([]*Session)
	}

	return
}

func (r *GeneratedOrganizationResolver) SessionsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*Session{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndSessionIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Session)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) AuditLogs(ctx context.Context, obj *Organization) (res []*AuditLog, err error) {
	return r.Handlers.OrganizationAuditLogs(ctx, r.GeneratedResolver, obj)
}
func OrganizationAuditLogsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*AuditLog, err error) {

	items = []*AuditLog{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "AuditLogs"); err != nil {
		return items, errors.New("AuditLogs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["AuditLogOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*AuditLog{}
	if item != nil {
		items = item.([]*AuditLog)
	}

	return
}

func (r *GeneratedOrganizationResolver) AuditLogsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*AuditLog{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndAuditLogIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*AuditLog)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) PaymentConfigs(ctx context.Context, obj *Organization) (res []*FranchisePaymentConfig, err error) {
	return r.Handlers.OrganizationPaymentConfigs(ctx, r.GeneratedResolver, obj)
}
func OrganizationPaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*FranchisePaymentConfig, err error) {

	items = []*FranchisePaymentConfig{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "PaymentConfigs"); err != nil {
		return items, errors.New("PaymentConfigs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["FranchisePaymentConfigOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*FranchisePaymentConfig{}
	if item != nil {
		items = item.([]*FranchisePaymentConfig)
	}

	return
}

func (r *GeneratedOrganizationResolver) PaymentConfigsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*FranchisePaymentConfig{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndFranchisePaymentConfigIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*FranchisePaymentConfig)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryOperatorMembershipHandlerOptions struct {
	ID     *string
	Filter *OperatorMembershipFilterType
}

func (r *GeneratedQueryResolver) OperatorMembership(ctx context.Context, id *string, filter *OperatorMembershipFilterType) (*OperatorMembership, error) {
	opts := QueryOperatorMembershipHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryOperatorMembership(ctx, r.GeneratedResolver, opts)
}
func QueryOperatorMembershipHandler(ctx context.Context, r *GeneratedResolver, opts QueryOperatorMembershipHandlerOptions) (*OperatorMembership, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := OperatorMembershipQueryFilter{}
	rt := &OperatorMembershipResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("operator_memberships", ctx)+".id = ?", *opts.ID)
	}

	var items []*OperatorMembership
	giOpts := GetItemsOptions{
		Alias:      TableName("operator_memberships", ctx),
		Preloaders: []string{},
		Item:       &OperatorMembership{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "OperatorMembership"}
	}
	return items[0], err
}

type QueryOperatorMembershipsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*OperatorMembershipSortType
	Filter      *OperatorMembershipFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) OperatorMemberships(ctx context.Context, current_page *int, per_page *int, q *string, sort []*OperatorMembershipSortType, filter *OperatorMembershipFilterType, rand *bool) (*OperatorMembershipResultType, error) {
	opts := QueryOperatorMembershipsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryOperatorMemberships(ctx, r.GeneratedResolver, opts)
}
func QueryOperatorMembershipsHandler(ctx context.Context, r *GeneratedResolver, opts QueryOperatorMembershipsHandlerOptions) (*OperatorMembershipResultType, error) {
	query := OperatorMembershipQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &OperatorMembershipResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedOperatorMembershipResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedOperatorMembershipResultTypeResolver) Data(ctx context.Context, obj *OperatorMembershipResultType) (items []*OperatorMembership, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("operator_memberships", ctx),
		Preloaders: []string{},
		Item:       &OperatorMembership{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*OperatorMembership{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedOperatorMembershipResultTypeResolver) Total(ctx context.Context, obj *OperatorMembershipResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("operator_memberships", ctx), &OperatorMembership{})
}

func (r *GeneratedOperatorMembershipResultTypeResolver) TotalPage(ctx context.Context, obj *OperatorMembershipResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedOperatorMembershipResultTypeResolver) CurrentPage(ctx context.Context, obj *OperatorMembershipResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedOperatorMembershipResultTypeResolver) PerPage(ctx context.Context, obj *OperatorMembershipResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedOperatorMembershipResolver struct{ *GeneratedResolver }

func (r *GeneratedOperatorMembershipResolver) Account(ctx context.Context, obj *OperatorMembership) (res *Account, err error) {
	return r.Handlers.OperatorMembershipAccount(ctx, r.GeneratedResolver, obj)
}
func OperatorMembershipAccountHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.AccountID

	if objKey != "" {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Account)

		if items == nil {
			items = &Account{}
		}

	}

	return
}

func (r *GeneratedOperatorMembershipResolver) Organization(ctx context.Context, obj *OperatorMembership) (res *Organization, err error) {
	return r.Handlers.OperatorMembershipOrganization(ctx, r.GeneratedResolver, obj)
}
func OperatorMembershipOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OrganizationID

	if objKey != "" {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Organization)

		if items == nil {
			items = &Organization{}
		}

	}

	return
}

func (r *GeneratedOperatorMembershipResolver) Roles(ctx context.Context, obj *OperatorMembership) (res []*OperatorRole, err error) {
	return r.Handlers.OperatorMembershipRoles(ctx, r.GeneratedResolver, obj)
}
func OperatorMembershipRolesHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (items []*OperatorRole, err error) {

	items = []*OperatorRole{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Roles"); err != nil {
		return items, errors.New("Roles " + err.Error())
	}

	// selects := GetFieldsRequested(ctx, strings.ToLower(TableName("operator_roles", ctx)))
	// wheres  := []string{}
	// values  := []interface{}{}
	// err = tx.Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(obj).Related(&items, "Roles").Error
	// err = r.DB.Query().Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(&OperatorRole{}).Find(&items).Error

	err = r.DB.Query().Model(obj).Order("weight ASC, created_at ASC").Preload("Roles").First(&obj).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return items, nil
		}
		return items, err
	}

	items = obj.Roles

	return
}

func (r *GeneratedOperatorMembershipResolver) RolesIds(ctx context.Context, obj *OperatorMembership) (ids []string, err error) {

	err = r.DB.Query().Order("weight ASC, created_at ASC").Preload("Roles").First(&obj).Error
	if err != nil {
		return
	}

	for _, item := range obj.Roles {
		ids = append(ids, item.ID)
	}

	return
}

func (r *GeneratedOperatorMembershipResolver) Stores(ctx context.Context, obj *OperatorMembership) (res []*Store, err error) {
	return r.Handlers.OperatorMembershipStores(ctx, r.GeneratedResolver, obj)
}
func OperatorMembershipStoresHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (items []*Store, err error) {

	items = []*Store{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Stores"); err != nil {
		return items, errors.New("Stores " + err.Error())
	}

	// selects := GetFieldsRequested(ctx, strings.ToLower(TableName("stores", ctx)))
	// wheres  := []string{}
	// values  := []interface{}{}
	// err = tx.Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(obj).Related(&items, "Stores").Error
	// err = r.DB.Query().Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(&Store{}).Find(&items).Error

	err = r.DB.Query().Model(obj).Order("weight ASC, created_at ASC").Preload("Stores").First(&obj).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return items, nil
		}
		return items, err
	}

	items = obj.Stores

	return
}

func (r *GeneratedOperatorMembershipResolver) StoresIds(ctx context.Context, obj *OperatorMembership) (ids []string, err error) {

	err = r.DB.Query().Order("weight ASC, created_at ASC").Preload("Stores").First(&obj).Error
	if err != nil {
		return
	}

	for _, item := range obj.Stores {
		ids = append(ids, item.ID)
	}

	return
}

func (r *GeneratedOperatorMembershipResolver) Invitations(ctx context.Context, obj *OperatorMembership) (res []*MembershipInvitation, err error) {
	return r.Handlers.OperatorMembershipInvitations(ctx, r.GeneratedResolver, obj)
}
func OperatorMembershipInvitationsHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorMembership) (items []*MembershipInvitation, err error) {

	items = []*MembershipInvitation{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Invitations"); err != nil {
		return items, errors.New("Invitations " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["MembershipInvitationMembership"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*MembershipInvitation{}
	if item != nil {
		items = item.([]*MembershipInvitation)
	}

	return
}

func (r *GeneratedOperatorMembershipResolver) InvitationsIds(ctx context.Context, obj *OperatorMembership) (ids []string, err error) {

	items := []*MembershipInvitation{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["MembershipAndMembershipInvitationIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*MembershipInvitation)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryPermissionHandlerOptions struct {
	ID     *string
	Filter *PermissionFilterType
}

func (r *GeneratedQueryResolver) Permission(ctx context.Context, id *string, filter *PermissionFilterType) (*Permission, error) {
	opts := QueryPermissionHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryPermission(ctx, r.GeneratedResolver, opts)
}
func QueryPermissionHandler(ctx context.Context, r *GeneratedResolver, opts QueryPermissionHandlerOptions) (*Permission, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := PermissionQueryFilter{}
	rt := &PermissionResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("permissions", ctx)+".id = ?", *opts.ID)
	}

	var items []*Permission
	giOpts := GetItemsOptions{
		Alias:      TableName("permissions", ctx),
		Preloaders: []string{},
		Item:       &Permission{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "Permission"}
	}
	return items[0], err
}

type QueryPermissionsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*PermissionSortType
	Filter      *PermissionFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) Permissions(ctx context.Context, current_page *int, per_page *int, q *string, sort []*PermissionSortType, filter *PermissionFilterType, rand *bool) (*PermissionResultType, error) {
	opts := QueryPermissionsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryPermissions(ctx, r.GeneratedResolver, opts)
}
func QueryPermissionsHandler(ctx context.Context, r *GeneratedResolver, opts QueryPermissionsHandlerOptions) (*PermissionResultType, error) {
	query := PermissionQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &PermissionResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedPermissionResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedPermissionResultTypeResolver) Data(ctx context.Context, obj *PermissionResultType) (items []*Permission, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("permissions", ctx),
		Preloaders: []string{},
		Item:       &Permission{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*Permission{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedPermissionResultTypeResolver) Total(ctx context.Context, obj *PermissionResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("permissions", ctx), &Permission{})
}

func (r *GeneratedPermissionResultTypeResolver) TotalPage(ctx context.Context, obj *PermissionResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedPermissionResultTypeResolver) CurrentPage(ctx context.Context, obj *PermissionResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedPermissionResultTypeResolver) PerPage(ctx context.Context, obj *PermissionResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedPermissionResolver struct{ *GeneratedResolver }

func (r *GeneratedPermissionResolver) Roles(ctx context.Context, obj *Permission) (res []*OperatorRole, err error) {
	return r.Handlers.PermissionRoles(ctx, r.GeneratedResolver, obj)
}
func PermissionRolesHandler(ctx context.Context, r *GeneratedResolver, obj *Permission) (items []*OperatorRole, err error) {

	items = []*OperatorRole{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Roles"); err != nil {
		return items, errors.New("Roles " + err.Error())
	}

	// selects := GetFieldsRequested(ctx, strings.ToLower(TableName("operator_roles", ctx)))
	// wheres  := []string{}
	// values  := []interface{}{}
	// err = tx.Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(obj).Related(&items, "Roles").Error
	// err = r.DB.Query().Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(&OperatorRole{}).Find(&items).Error

	err = r.DB.Query().Model(obj).Order("weight ASC, created_at ASC").Preload("Roles").First(&obj).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return items, nil
		}
		return items, err
	}

	items = obj.Roles

	return
}

func (r *GeneratedPermissionResolver) RolesIds(ctx context.Context, obj *Permission) (ids []string, err error) {

	err = r.DB.Query().Order("weight ASC, created_at ASC").Preload("Roles").First(&obj).Error
	if err != nil {
		return
	}

	for _, item := range obj.Roles {
		ids = append(ids, item.ID)
	}

	return
}

type QueryOperatorRoleHandlerOptions struct {
	ID     *string
	Filter *OperatorRoleFilterType
}

func (r *GeneratedQueryResolver) OperatorRole(ctx context.Context, id *string, filter *OperatorRoleFilterType) (*OperatorRole, error) {
	opts := QueryOperatorRoleHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryOperatorRole(ctx, r.GeneratedResolver, opts)
}
func QueryOperatorRoleHandler(ctx context.Context, r *GeneratedResolver, opts QueryOperatorRoleHandlerOptions) (*OperatorRole, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := OperatorRoleQueryFilter{}
	rt := &OperatorRoleResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("operator_roles", ctx)+".id = ?", *opts.ID)
	}

	var items []*OperatorRole
	giOpts := GetItemsOptions{
		Alias:      TableName("operator_roles", ctx),
		Preloaders: []string{},
		Item:       &OperatorRole{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "OperatorRole"}
	}
	return items[0], err
}

type QueryOperatorRolesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*OperatorRoleSortType
	Filter      *OperatorRoleFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) OperatorRoles(ctx context.Context, current_page *int, per_page *int, q *string, sort []*OperatorRoleSortType, filter *OperatorRoleFilterType, rand *bool) (*OperatorRoleResultType, error) {
	opts := QueryOperatorRolesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryOperatorRoles(ctx, r.GeneratedResolver, opts)
}
func QueryOperatorRolesHandler(ctx context.Context, r *GeneratedResolver, opts QueryOperatorRolesHandlerOptions) (*OperatorRoleResultType, error) {
	query := OperatorRoleQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &OperatorRoleResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedOperatorRoleResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedOperatorRoleResultTypeResolver) Data(ctx context.Context, obj *OperatorRoleResultType) (items []*OperatorRole, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("operator_roles", ctx),
		Preloaders: []string{},
		Item:       &OperatorRole{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*OperatorRole{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedOperatorRoleResultTypeResolver) Total(ctx context.Context, obj *OperatorRoleResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("operator_roles", ctx), &OperatorRole{})
}

func (r *GeneratedOperatorRoleResultTypeResolver) TotalPage(ctx context.Context, obj *OperatorRoleResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedOperatorRoleResultTypeResolver) CurrentPage(ctx context.Context, obj *OperatorRoleResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedOperatorRoleResultTypeResolver) PerPage(ctx context.Context, obj *OperatorRoleResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedOperatorRoleResolver struct{ *GeneratedResolver }

func (r *GeneratedOperatorRoleResolver) Organization(ctx context.Context, obj *OperatorRole) (res *Organization, err error) {
	return r.Handlers.OperatorRoleOrganization(ctx, r.GeneratedResolver, obj)
}
func OperatorRoleOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorRole) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OrganizationID

	if objKey != "" {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Organization)

		if items == nil {
			items = &Organization{}
		}

	}

	return
}

func (r *GeneratedOperatorRoleResolver) Members(ctx context.Context, obj *OperatorRole) (res []*OperatorMembership, err error) {
	return r.Handlers.OperatorRoleMembers(ctx, r.GeneratedResolver, obj)
}
func OperatorRoleMembersHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorRole) (items []*OperatorMembership, err error) {

	items = []*OperatorMembership{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Members"); err != nil {
		return items, errors.New("Members " + err.Error())
	}

	// selects := GetFieldsRequested(ctx, strings.ToLower(TableName("operator_memberships", ctx)))
	// wheres  := []string{}
	// values  := []interface{}{}
	// err = tx.Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(obj).Related(&items, "Members").Error
	// err = r.DB.Query().Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(&OperatorMembership{}).Find(&items).Error

	err = r.DB.Query().Model(obj).Order("weight ASC, created_at ASC").Preload("Members").First(&obj).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return items, nil
		}
		return items, err
	}

	items = obj.Members

	return
}

func (r *GeneratedOperatorRoleResolver) MembersIds(ctx context.Context, obj *OperatorRole) (ids []string, err error) {

	err = r.DB.Query().Order("weight ASC, created_at ASC").Preload("Members").First(&obj).Error
	if err != nil {
		return
	}

	for _, item := range obj.Members {
		ids = append(ids, item.ID)
	}

	return
}

func (r *GeneratedOperatorRoleResolver) Permissions(ctx context.Context, obj *OperatorRole) (res []*Permission, err error) {
	return r.Handlers.OperatorRolePermissions(ctx, r.GeneratedResolver, obj)
}
func OperatorRolePermissionsHandler(ctx context.Context, r *GeneratedResolver, obj *OperatorRole) (items []*Permission, err error) {

	items = []*Permission{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Permissions"); err != nil {
		return items, errors.New("Permissions " + err.Error())
	}

	// selects := GetFieldsRequested(ctx, strings.ToLower(TableName("permissions", ctx)))
	// wheres  := []string{}
	// values  := []interface{}{}
	// err = tx.Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(obj).Related(&items, "Permissions").Error
	// err = r.DB.Query().Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(&Permission{}).Find(&items).Error

	err = r.DB.Query().Model(obj).Order("weight ASC, created_at ASC").Preload("Permissions").First(&obj).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return items, nil
		}
		return items, err
	}

	items = obj.Permissions

	return
}

func (r *GeneratedOperatorRoleResolver) PermissionsIds(ctx context.Context, obj *OperatorRole) (ids []string, err error) {

	err = r.DB.Query().Order("weight ASC, created_at ASC").Preload("Permissions").First(&obj).Error
	if err != nil {
		return
	}

	for _, item := range obj.Permissions {
		ids = append(ids, item.ID)
	}

	return
}

type QueryStoreHandlerOptions struct {
	ID     *string
	Filter *StoreFilterType
}

func (r *GeneratedQueryResolver) Store(ctx context.Context, id *string, filter *StoreFilterType) (*Store, error) {
	opts := QueryStoreHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStore(ctx, r.GeneratedResolver, opts)
}
func QueryStoreHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreHandlerOptions) (*Store, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StoreQueryFilter{}
	rt := &StoreResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("stores", ctx)+".id = ?", *opts.ID)
	}

	var items []*Store
	giOpts := GetItemsOptions{
		Alias:      TableName("stores", ctx),
		Preloaders: []string{},
		Item:       &Store{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "Store"}
	}
	return items[0], err
}

type QueryStoresHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StoreSortType
	Filter      *StoreFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) Stores(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StoreSortType, filter *StoreFilterType, rand *bool) (*StoreResultType, error) {
	opts := QueryStoresHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStores(ctx, r.GeneratedResolver, opts)
}
func QueryStoresHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoresHandlerOptions) (*StoreResultType, error) {
	query := StoreQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &StoreResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedStoreResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreResultTypeResolver) Data(ctx context.Context, obj *StoreResultType) (items []*Store, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("stores", ctx),
		Preloaders: []string{},
		Item:       &Store{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*Store{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedStoreResultTypeResolver) Total(ctx context.Context, obj *StoreResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("stores", ctx), &Store{})
}

func (r *GeneratedStoreResultTypeResolver) TotalPage(ctx context.Context, obj *StoreResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedStoreResultTypeResolver) CurrentPage(ctx context.Context, obj *StoreResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStoreResultTypeResolver) PerPage(ctx context.Context, obj *StoreResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStoreResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreResolver) Organization(ctx context.Context, obj *Store) (res *Organization, err error) {
	return r.Handlers.StoreOrganization(ctx, r.GeneratedResolver, obj)
}
func StoreOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OrganizationID

	if objKey != "" {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Organization)

		if items == nil {
			items = &Organization{}
		}

	}

	return
}

func (r *GeneratedStoreResolver) Members(ctx context.Context, obj *Store) (res []*OperatorMembership, err error) {
	return r.Handlers.StoreMembers(ctx, r.GeneratedResolver, obj)
}
func StoreMembersHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*OperatorMembership, err error) {

	items = []*OperatorMembership{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Members"); err != nil {
		return items, errors.New("Members " + err.Error())
	}

	// selects := GetFieldsRequested(ctx, strings.ToLower(TableName("operator_memberships", ctx)))
	// wheres  := []string{}
	// values  := []interface{}{}
	// err = tx.Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(obj).Related(&items, "Members").Error
	// err = r.DB.Query().Select(selects).Where(strings.Join(wheres, " AND "), values...).Model(&OperatorMembership{}).Find(&items).Error

	err = r.DB.Query().Model(obj).Order("weight ASC, created_at ASC").Preload("Members").First(&obj).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return items, nil
		}
		return items, err
	}

	items = obj.Members

	return
}

func (r *GeneratedStoreResolver) MembersIds(ctx context.Context, obj *Store) (ids []string, err error) {

	err = r.DB.Query().Order("weight ASC, created_at ASC").Preload("Members").First(&obj).Error
	if err != nil {
		return
	}

	for _, item := range obj.Members {
		ids = append(ids, item.ID)
	}

	return
}

func (r *GeneratedStoreResolver) ReviewedByAccount(ctx context.Context, obj *Store) (res *Account, err error) {
	return r.Handlers.StoreReviewedByAccount(ctx, r.GeneratedResolver, obj)
}
func StoreReviewedByAccountHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ReviewedByAccountID

	if objKey != nil {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Account)

	}

	return
}

func (r *GeneratedStoreResolver) AuditLogs(ctx context.Context, obj *Store) (res []*AuditLog, err error) {
	return r.Handlers.StoreAuditLogs(ctx, r.GeneratedResolver, obj)
}
func StoreAuditLogsHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*AuditLog, err error) {

	items = []*AuditLog{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "AuditLogs"); err != nil {
		return items, errors.New("AuditLogs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["AuditLogStore"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*AuditLog{}
	if item != nil {
		items = item.([]*AuditLog)
	}

	return
}

func (r *GeneratedStoreResolver) AuditLogsIds(ctx context.Context, obj *Store) (ids []string, err error) {

	items := []*AuditLog{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreAndAuditLogIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*AuditLog)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreResolver) PaymentConfigs(ctx context.Context, obj *Store) (res []*StorePaymentConfig, err error) {
	return r.Handlers.StorePaymentConfigs(ctx, r.GeneratedResolver, obj)
}
func StorePaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*StorePaymentConfig, err error) {

	items = []*StorePaymentConfig{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "PaymentConfigs"); err != nil {
		return items, errors.New("PaymentConfigs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StorePaymentConfigStore"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StorePaymentConfig{}
	if item != nil {
		items = item.([]*StorePaymentConfig)
	}

	return
}

func (r *GeneratedStoreResolver) PaymentConfigsIds(ctx context.Context, obj *Store) (ids []string, err error) {

	items := []*StorePaymentConfig{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreAndStorePaymentConfigIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StorePaymentConfig)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QuerySessionHandlerOptions struct {
	ID     *string
	Filter *SessionFilterType
}

func (r *GeneratedQueryResolver) Session(ctx context.Context, id *string, filter *SessionFilterType) (*Session, error) {
	opts := QuerySessionHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QuerySession(ctx, r.GeneratedResolver, opts)
}
func QuerySessionHandler(ctx context.Context, r *GeneratedResolver, opts QuerySessionHandlerOptions) (*Session, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := SessionQueryFilter{}
	rt := &SessionResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("sessions", ctx)+".id = ?", *opts.ID)
	}

	var items []*Session
	giOpts := GetItemsOptions{
		Alias:      TableName("sessions", ctx),
		Preloaders: []string{},
		Item:       &Session{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "Session"}
	}
	return items[0], err
}

type QuerySessionsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*SessionSortType
	Filter      *SessionFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) Sessions(ctx context.Context, current_page *int, per_page *int, q *string, sort []*SessionSortType, filter *SessionFilterType, rand *bool) (*SessionResultType, error) {
	opts := QuerySessionsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QuerySessions(ctx, r.GeneratedResolver, opts)
}
func QuerySessionsHandler(ctx context.Context, r *GeneratedResolver, opts QuerySessionsHandlerOptions) (*SessionResultType, error) {
	query := SessionQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &SessionResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedSessionResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedSessionResultTypeResolver) Data(ctx context.Context, obj *SessionResultType) (items []*Session, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("sessions", ctx),
		Preloaders: []string{},
		Item:       &Session{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*Session{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedSessionResultTypeResolver) Total(ctx context.Context, obj *SessionResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("sessions", ctx), &Session{})
}

func (r *GeneratedSessionResultTypeResolver) TotalPage(ctx context.Context, obj *SessionResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedSessionResultTypeResolver) CurrentPage(ctx context.Context, obj *SessionResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedSessionResultTypeResolver) PerPage(ctx context.Context, obj *SessionResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedSessionResolver struct{ *GeneratedResolver }

func (r *GeneratedSessionResolver) Account(ctx context.Context, obj *Session) (res *Account, err error) {
	return r.Handlers.SessionAccount(ctx, r.GeneratedResolver, obj)
}
func SessionAccountHandler(ctx context.Context, r *GeneratedResolver, obj *Session) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.AccountID

	if objKey != "" {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Account)

		if items == nil {
			items = &Account{}
		}

	}

	return
}

func (r *GeneratedSessionResolver) Organization(ctx context.Context, obj *Session) (res *Organization, err error) {
	return r.Handlers.SessionOrganization(ctx, r.GeneratedResolver, obj)
}
func SessionOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *Session) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OrganizationID

	if objKey != nil {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Organization)

	}

	return
}

func (r *GeneratedSessionResolver) AuditLogs(ctx context.Context, obj *Session) (res []*AuditLog, err error) {
	return r.Handlers.SessionAuditLogs(ctx, r.GeneratedResolver, obj)
}
func SessionAuditLogsHandler(ctx context.Context, r *GeneratedResolver, obj *Session) (items []*AuditLog, err error) {

	items = []*AuditLog{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "AuditLogs"); err != nil {
		return items, errors.New("AuditLogs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["AuditLogSession"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*AuditLog{}
	if item != nil {
		items = item.([]*AuditLog)
	}

	return
}

func (r *GeneratedSessionResolver) AuditLogsIds(ctx context.Context, obj *Session) (ids []string, err error) {

	items := []*AuditLog{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SessionAndAuditLogIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*AuditLog)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryMembershipInvitationHandlerOptions struct {
	ID     *string
	Filter *MembershipInvitationFilterType
}

func (r *GeneratedQueryResolver) MembershipInvitation(ctx context.Context, id *string, filter *MembershipInvitationFilterType) (*MembershipInvitation, error) {
	opts := QueryMembershipInvitationHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryMembershipInvitation(ctx, r.GeneratedResolver, opts)
}
func QueryMembershipInvitationHandler(ctx context.Context, r *GeneratedResolver, opts QueryMembershipInvitationHandlerOptions) (*MembershipInvitation, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := MembershipInvitationQueryFilter{}
	rt := &MembershipInvitationResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("membership_invitations", ctx)+".id = ?", *opts.ID)
	}

	var items []*MembershipInvitation
	giOpts := GetItemsOptions{
		Alias:      TableName("membership_invitations", ctx),
		Preloaders: []string{},
		Item:       &MembershipInvitation{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "MembershipInvitation"}
	}
	return items[0], err
}

type QueryMembershipInvitationsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*MembershipInvitationSortType
	Filter      *MembershipInvitationFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) MembershipInvitations(ctx context.Context, current_page *int, per_page *int, q *string, sort []*MembershipInvitationSortType, filter *MembershipInvitationFilterType, rand *bool) (*MembershipInvitationResultType, error) {
	opts := QueryMembershipInvitationsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryMembershipInvitations(ctx, r.GeneratedResolver, opts)
}
func QueryMembershipInvitationsHandler(ctx context.Context, r *GeneratedResolver, opts QueryMembershipInvitationsHandlerOptions) (*MembershipInvitationResultType, error) {
	query := MembershipInvitationQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &MembershipInvitationResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedMembershipInvitationResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedMembershipInvitationResultTypeResolver) Data(ctx context.Context, obj *MembershipInvitationResultType) (items []*MembershipInvitation, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("membership_invitations", ctx),
		Preloaders: []string{},
		Item:       &MembershipInvitation{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*MembershipInvitation{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedMembershipInvitationResultTypeResolver) Total(ctx context.Context, obj *MembershipInvitationResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("membership_invitations", ctx), &MembershipInvitation{})
}

func (r *GeneratedMembershipInvitationResultTypeResolver) TotalPage(ctx context.Context, obj *MembershipInvitationResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedMembershipInvitationResultTypeResolver) CurrentPage(ctx context.Context, obj *MembershipInvitationResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedMembershipInvitationResultTypeResolver) PerPage(ctx context.Context, obj *MembershipInvitationResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedMembershipInvitationResolver struct{ *GeneratedResolver }

func (r *GeneratedMembershipInvitationResolver) Membership(ctx context.Context, obj *MembershipInvitation) (res *OperatorMembership, err error) {
	return r.Handlers.MembershipInvitationMembership(ctx, r.GeneratedResolver, obj)
}
func MembershipInvitationMembershipHandler(ctx context.Context, r *GeneratedResolver, obj *MembershipInvitation) (items *OperatorMembership, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
		return items, errors.New("OperatorMembership " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.MembershipID

	if objKey != "" {
		item, _ := loaders["OperatorMembership"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*OperatorMembership)

		if items == nil {
			items = &OperatorMembership{}
		}

	}

	return
}

func (r *GeneratedMembershipInvitationResolver) InvitedByAccount(ctx context.Context, obj *MembershipInvitation) (res *Account, err error) {
	return r.Handlers.MembershipInvitationInvitedByAccount(ctx, r.GeneratedResolver, obj)
}
func MembershipInvitationInvitedByAccountHandler(ctx context.Context, r *GeneratedResolver, obj *MembershipInvitation) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.InvitedByAccountID

	if objKey != "" {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Account)

		if items == nil {
			items = &Account{}
		}

	}

	return
}

type QueryAuditLogHandlerOptions struct {
	ID     *string
	Filter *AuditLogFilterType
}

func (r *GeneratedQueryResolver) AuditLog(ctx context.Context, id *string, filter *AuditLogFilterType) (*AuditLog, error) {
	opts := QueryAuditLogHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryAuditLog(ctx, r.GeneratedResolver, opts)
}
func QueryAuditLogHandler(ctx context.Context, r *GeneratedResolver, opts QueryAuditLogHandlerOptions) (*AuditLog, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := AuditLogQueryFilter{}
	rt := &AuditLogResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("audit_logs", ctx)+".id = ?", *opts.ID)
	}

	var items []*AuditLog
	giOpts := GetItemsOptions{
		Alias:      TableName("audit_logs", ctx),
		Preloaders: []string{},
		Item:       &AuditLog{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "AuditLog"}
	}
	return items[0], err
}

type QueryAuditLogsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*AuditLogSortType
	Filter      *AuditLogFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) AuditLogs(ctx context.Context, current_page *int, per_page *int, q *string, sort []*AuditLogSortType, filter *AuditLogFilterType, rand *bool) (*AuditLogResultType, error) {
	opts := QueryAuditLogsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryAuditLogs(ctx, r.GeneratedResolver, opts)
}
func QueryAuditLogsHandler(ctx context.Context, r *GeneratedResolver, opts QueryAuditLogsHandlerOptions) (*AuditLogResultType, error) {
	query := AuditLogQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &AuditLogResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedAuditLogResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedAuditLogResultTypeResolver) Data(ctx context.Context, obj *AuditLogResultType) (items []*AuditLog, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("audit_logs", ctx),
		Preloaders: []string{},
		Item:       &AuditLog{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*AuditLog{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedAuditLogResultTypeResolver) Total(ctx context.Context, obj *AuditLogResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("audit_logs", ctx), &AuditLog{})
}

func (r *GeneratedAuditLogResultTypeResolver) TotalPage(ctx context.Context, obj *AuditLogResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedAuditLogResultTypeResolver) CurrentPage(ctx context.Context, obj *AuditLogResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedAuditLogResultTypeResolver) PerPage(ctx context.Context, obj *AuditLogResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedAuditLogResolver struct{ *GeneratedResolver }

func (r *GeneratedAuditLogResolver) ActorAccount(ctx context.Context, obj *AuditLog) (res *Account, err error) {
	return r.Handlers.AuditLogActorAccount(ctx, r.GeneratedResolver, obj)
}
func AuditLogActorAccountHandler(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ActorAccountID

	if objKey != nil {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Account)

	}

	return
}

func (r *GeneratedAuditLogResolver) Session(ctx context.Context, obj *AuditLog) (res *Session, err error) {
	return r.Handlers.AuditLogSession(ctx, r.GeneratedResolver, obj)
}
func AuditLogSessionHandler(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (items *Session, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
		return items, errors.New("Session " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.SessionID

	if objKey != nil {
		item, _ := loaders["Session"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Session)

	}

	return
}

func (r *GeneratedAuditLogResolver) Organization(ctx context.Context, obj *AuditLog) (res *Organization, err error) {
	return r.Handlers.AuditLogOrganization(ctx, r.GeneratedResolver, obj)
}
func AuditLogOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OrganizationID

	if objKey != nil {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Organization)

	}

	return
}

func (r *GeneratedAuditLogResolver) Store(ctx context.Context, obj *AuditLog) (res *Store, err error) {
	return r.Handlers.AuditLogStore(ctx, r.GeneratedResolver, obj)
}
func AuditLogStoreHandler(ctx context.Context, r *GeneratedResolver, obj *AuditLog) (items *Store, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
		return items, errors.New("Store " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.StoreID

	if objKey != nil {
		item, _ := loaders["Store"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Store)

	}

	return
}

type QueryFranchiseOpeningRecordHandlerOptions struct {
	ID     *string
	Filter *FranchiseOpeningRecordFilterType
}

func (r *GeneratedQueryResolver) FranchiseOpeningRecord(ctx context.Context, id *string, filter *FranchiseOpeningRecordFilterType) (*FranchiseOpeningRecord, error) {
	opts := QueryFranchiseOpeningRecordHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryFranchiseOpeningRecord(ctx, r.GeneratedResolver, opts)
}
func QueryFranchiseOpeningRecordHandler(ctx context.Context, r *GeneratedResolver, opts QueryFranchiseOpeningRecordHandlerOptions) (*FranchiseOpeningRecord, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := FranchiseOpeningRecordQueryFilter{}
	rt := &FranchiseOpeningRecordResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("franchise_opening_records", ctx)+".id = ?", *opts.ID)
	}

	var items []*FranchiseOpeningRecord
	giOpts := GetItemsOptions{
		Alias:      TableName("franchise_opening_records", ctx),
		Preloaders: []string{},
		Item:       &FranchiseOpeningRecord{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "FranchiseOpeningRecord"}
	}
	return items[0], err
}

type QueryFranchiseOpeningRecordsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*FranchiseOpeningRecordSortType
	Filter      *FranchiseOpeningRecordFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) FranchiseOpeningRecords(ctx context.Context, current_page *int, per_page *int, q *string, sort []*FranchiseOpeningRecordSortType, filter *FranchiseOpeningRecordFilterType, rand *bool) (*FranchiseOpeningRecordResultType, error) {
	opts := QueryFranchiseOpeningRecordsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryFranchiseOpeningRecords(ctx, r.GeneratedResolver, opts)
}
func QueryFranchiseOpeningRecordsHandler(ctx context.Context, r *GeneratedResolver, opts QueryFranchiseOpeningRecordsHandlerOptions) (*FranchiseOpeningRecordResultType, error) {
	query := FranchiseOpeningRecordQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &FranchiseOpeningRecordResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedFranchiseOpeningRecordResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedFranchiseOpeningRecordResultTypeResolver) Data(ctx context.Context, obj *FranchiseOpeningRecordResultType) (items []*FranchiseOpeningRecord, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("franchise_opening_records", ctx),
		Preloaders: []string{},
		Item:       &FranchiseOpeningRecord{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*FranchiseOpeningRecord{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedFranchiseOpeningRecordResultTypeResolver) Total(ctx context.Context, obj *FranchiseOpeningRecordResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("franchise_opening_records", ctx), &FranchiseOpeningRecord{})
}

func (r *GeneratedFranchiseOpeningRecordResultTypeResolver) TotalPage(ctx context.Context, obj *FranchiseOpeningRecordResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedFranchiseOpeningRecordResultTypeResolver) CurrentPage(ctx context.Context, obj *FranchiseOpeningRecordResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedFranchiseOpeningRecordResultTypeResolver) PerPage(ctx context.Context, obj *FranchiseOpeningRecordResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedFranchiseOpeningRecordResolver struct{ *GeneratedResolver }

func (r *GeneratedFranchiseOpeningRecordResolver) Organization(ctx context.Context, obj *FranchiseOpeningRecord) (res *Organization, err error) {
	return r.Handlers.FranchiseOpeningRecordOrganization(ctx, r.GeneratedResolver, obj)
}
func FranchiseOpeningRecordOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *FranchiseOpeningRecord) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OrganizationID

	if objKey != "" {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Organization)

		if items == nil {
			items = &Organization{}
		}

	}

	return
}

func (r *GeneratedFranchiseOpeningRecordResolver) InitialAccount(ctx context.Context, obj *FranchiseOpeningRecord) (res *Account, err error) {
	return r.Handlers.FranchiseOpeningRecordInitialAccount(ctx, r.GeneratedResolver, obj)
}
func FranchiseOpeningRecordInitialAccountHandler(ctx context.Context, r *GeneratedResolver, obj *FranchiseOpeningRecord) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.InitialAccountID

	if objKey != "" {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Account)

		if items == nil {
			items = &Account{}
		}

	}

	return
}

func (r *GeneratedFranchiseOpeningRecordResolver) RecordedByAccount(ctx context.Context, obj *FranchiseOpeningRecord) (res *Account, err error) {
	return r.Handlers.FranchiseOpeningRecordRecordedByAccount(ctx, r.GeneratedResolver, obj)
}
func FranchiseOpeningRecordRecordedByAccountHandler(ctx context.Context, r *GeneratedResolver, obj *FranchiseOpeningRecord) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.RecordedByAccountID

	if objKey != "" {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Account)

		if items == nil {
			items = &Account{}
		}

	}

	return
}

type QueryGlobalPaymentConfigHandlerOptions struct {
	ID     *string
	Filter *GlobalPaymentConfigFilterType
}

func (r *GeneratedQueryResolver) GlobalPaymentConfig(ctx context.Context, id *string, filter *GlobalPaymentConfigFilterType) (*GlobalPaymentConfig, error) {
	opts := QueryGlobalPaymentConfigHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryGlobalPaymentConfig(ctx, r.GeneratedResolver, opts)
}
func QueryGlobalPaymentConfigHandler(ctx context.Context, r *GeneratedResolver, opts QueryGlobalPaymentConfigHandlerOptions) (*GlobalPaymentConfig, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := GlobalPaymentConfigQueryFilter{}
	rt := &GlobalPaymentConfigResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("global_payment_configs", ctx)+".id = ?", *opts.ID)
	}

	var items []*GlobalPaymentConfig
	giOpts := GetItemsOptions{
		Alias:      TableName("global_payment_configs", ctx),
		Preloaders: []string{},
		Item:       &GlobalPaymentConfig{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "GlobalPaymentConfig"}
	}
	return items[0], err
}

type QueryGlobalPaymentConfigsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*GlobalPaymentConfigSortType
	Filter      *GlobalPaymentConfigFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) GlobalPaymentConfigs(ctx context.Context, current_page *int, per_page *int, q *string, sort []*GlobalPaymentConfigSortType, filter *GlobalPaymentConfigFilterType, rand *bool) (*GlobalPaymentConfigResultType, error) {
	opts := QueryGlobalPaymentConfigsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryGlobalPaymentConfigs(ctx, r.GeneratedResolver, opts)
}
func QueryGlobalPaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, opts QueryGlobalPaymentConfigsHandlerOptions) (*GlobalPaymentConfigResultType, error) {
	query := GlobalPaymentConfigQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &GlobalPaymentConfigResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedGlobalPaymentConfigResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedGlobalPaymentConfigResultTypeResolver) Data(ctx context.Context, obj *GlobalPaymentConfigResultType) (items []*GlobalPaymentConfig, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("global_payment_configs", ctx),
		Preloaders: []string{},
		Item:       &GlobalPaymentConfig{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*GlobalPaymentConfig{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedGlobalPaymentConfigResultTypeResolver) Total(ctx context.Context, obj *GlobalPaymentConfigResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("global_payment_configs", ctx), &GlobalPaymentConfig{})
}

func (r *GeneratedGlobalPaymentConfigResultTypeResolver) TotalPage(ctx context.Context, obj *GlobalPaymentConfigResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedGlobalPaymentConfigResultTypeResolver) CurrentPage(ctx context.Context, obj *GlobalPaymentConfigResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedGlobalPaymentConfigResultTypeResolver) PerPage(ctx context.Context, obj *GlobalPaymentConfigResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type QueryFranchisePaymentConfigHandlerOptions struct {
	ID     *string
	Filter *FranchisePaymentConfigFilterType
}

func (r *GeneratedQueryResolver) FranchisePaymentConfig(ctx context.Context, id *string, filter *FranchisePaymentConfigFilterType) (*FranchisePaymentConfig, error) {
	opts := QueryFranchisePaymentConfigHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryFranchisePaymentConfig(ctx, r.GeneratedResolver, opts)
}
func QueryFranchisePaymentConfigHandler(ctx context.Context, r *GeneratedResolver, opts QueryFranchisePaymentConfigHandlerOptions) (*FranchisePaymentConfig, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := FranchisePaymentConfigQueryFilter{}
	rt := &FranchisePaymentConfigResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("franchise_payment_configs", ctx)+".id = ?", *opts.ID)
	}

	var items []*FranchisePaymentConfig
	giOpts := GetItemsOptions{
		Alias:      TableName("franchise_payment_configs", ctx),
		Preloaders: []string{},
		Item:       &FranchisePaymentConfig{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "FranchisePaymentConfig"}
	}
	return items[0], err
}

type QueryFranchisePaymentConfigsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*FranchisePaymentConfigSortType
	Filter      *FranchisePaymentConfigFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) FranchisePaymentConfigs(ctx context.Context, current_page *int, per_page *int, q *string, sort []*FranchisePaymentConfigSortType, filter *FranchisePaymentConfigFilterType, rand *bool) (*FranchisePaymentConfigResultType, error) {
	opts := QueryFranchisePaymentConfigsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryFranchisePaymentConfigs(ctx, r.GeneratedResolver, opts)
}
func QueryFranchisePaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, opts QueryFranchisePaymentConfigsHandlerOptions) (*FranchisePaymentConfigResultType, error) {
	query := FranchisePaymentConfigQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &FranchisePaymentConfigResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedFranchisePaymentConfigResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedFranchisePaymentConfigResultTypeResolver) Data(ctx context.Context, obj *FranchisePaymentConfigResultType) (items []*FranchisePaymentConfig, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("franchise_payment_configs", ctx),
		Preloaders: []string{},
		Item:       &FranchisePaymentConfig{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*FranchisePaymentConfig{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedFranchisePaymentConfigResultTypeResolver) Total(ctx context.Context, obj *FranchisePaymentConfigResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("franchise_payment_configs", ctx), &FranchisePaymentConfig{})
}

func (r *GeneratedFranchisePaymentConfigResultTypeResolver) TotalPage(ctx context.Context, obj *FranchisePaymentConfigResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedFranchisePaymentConfigResultTypeResolver) CurrentPage(ctx context.Context, obj *FranchisePaymentConfigResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedFranchisePaymentConfigResultTypeResolver) PerPage(ctx context.Context, obj *FranchisePaymentConfigResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedFranchisePaymentConfigResolver struct{ *GeneratedResolver }

func (r *GeneratedFranchisePaymentConfigResolver) Organization(ctx context.Context, obj *FranchisePaymentConfig) (res *Organization, err error) {
	return r.Handlers.FranchisePaymentConfigOrganization(ctx, r.GeneratedResolver, obj)
}
func FranchisePaymentConfigOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *FranchisePaymentConfig) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OrganizationID

	if objKey != "" {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Organization)

		if items == nil {
			items = &Organization{}
		}

	}

	return
}

type QueryStorePaymentConfigHandlerOptions struct {
	ID     *string
	Filter *StorePaymentConfigFilterType
}

func (r *GeneratedQueryResolver) StorePaymentConfig(ctx context.Context, id *string, filter *StorePaymentConfigFilterType) (*StorePaymentConfig, error) {
	opts := QueryStorePaymentConfigHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStorePaymentConfig(ctx, r.GeneratedResolver, opts)
}
func QueryStorePaymentConfigHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePaymentConfigHandlerOptions) (*StorePaymentConfig, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StorePaymentConfigQueryFilter{}
	rt := &StorePaymentConfigResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_payment_configs", ctx)+".id = ?", *opts.ID)
	}

	var items []*StorePaymentConfig
	giOpts := GetItemsOptions{
		Alias:      TableName("store_payment_configs", ctx),
		Preloaders: []string{},
		Item:       &StorePaymentConfig{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StorePaymentConfig"}
	}
	return items[0], err
}

type QueryStorePaymentConfigsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StorePaymentConfigSortType
	Filter      *StorePaymentConfigFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StorePaymentConfigs(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StorePaymentConfigSortType, filter *StorePaymentConfigFilterType, rand *bool) (*StorePaymentConfigResultType, error) {
	opts := QueryStorePaymentConfigsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStorePaymentConfigs(ctx, r.GeneratedResolver, opts)
}
func QueryStorePaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePaymentConfigsHandlerOptions) (*StorePaymentConfigResultType, error) {
	query := StorePaymentConfigQueryFilter{opts.Q}

	var selectionSet *ast.SelectionSet
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			if f.Field.Name == "data" {
				selectionSet = &f.Field.SelectionSet
			}
		}
	}()

	_sort := []EntitySort{}
	for _, sort := range opts.Sort {
		_sort = append(_sort, sort)
	}

	return &StorePaymentConfigResultType{
		EntityResultType: EntityResultType{
			CurrentPage:  opts.CurrentPage,
			PerPage:      opts.PerPage,
			Rand:         opts.Rand,
			Query:        &query,
			Sort:         _sort,
			Filter:       opts.Filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type GeneratedStorePaymentConfigResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePaymentConfigResultTypeResolver) Data(ctx context.Context, obj *StorePaymentConfigResultType) (items []*StorePaymentConfig, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_payment_configs", ctx),
		Preloaders: []string{},
		Item:       &StorePaymentConfig{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StorePaymentConfig{}
	idMap := map[string]bool{}
	for _, item := range items {
		if _, ok := idMap[item.ID]; !ok {
			idMap[item.ID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	items = uniqueItems

	return
}

func (r *GeneratedStorePaymentConfigResultTypeResolver) Total(ctx context.Context, obj *StorePaymentConfigResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_payment_configs", ctx), &StorePaymentConfig{})
}

func (r *GeneratedStorePaymentConfigResultTypeResolver) TotalPage(ctx context.Context, obj *StorePaymentConfigResultType) (count int, err error) {
	total, _ := r.Total(ctx, obj)
	perPage, _ := r.PerPage(ctx, obj)
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPage < 0 {
		totalPage = 0
	} else if perPage <= 0 {
		totalPage = total
	}

	return totalPage, nil
}

func (r *GeneratedStorePaymentConfigResultTypeResolver) CurrentPage(ctx context.Context, obj *StorePaymentConfigResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStorePaymentConfigResultTypeResolver) PerPage(ctx context.Context, obj *StorePaymentConfigResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStorePaymentConfigResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePaymentConfigResolver) Store(ctx context.Context, obj *StorePaymentConfig) (res *Store, err error) {
	return r.Handlers.StorePaymentConfigStore(ctx, r.GeneratedResolver, obj)
}
func StorePaymentConfigStoreHandler(ctx context.Context, r *GeneratedResolver, obj *StorePaymentConfig) (items *Store, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
		return items, errors.New("Store " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.StoreID

	if objKey != "" {
		item, _ := loaders["Store"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Store)

		if items == nil {
			items = &Store{}
		}

	}

	return
}
