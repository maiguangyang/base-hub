package gen

import (
	"context"
	"errors"
	"math"

	"github.com/99designs/gqlgen/graphql"
	"github.com/graph-gophers/dataloader"
	"base-engine/auth"
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

func (r *GeneratedAccountResolver) CreatedStocktakes(ctx context.Context, obj *Account) (res []*StoreStocktake, err error) {
	return r.Handlers.AccountCreatedStocktakes(ctx, r.GeneratedResolver, obj)
}
func AccountCreatedStocktakesHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*StoreStocktake, err error) {

	items = []*StoreStocktake{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CreatedStocktakes"); err != nil {
		return items, errors.New("CreatedStocktakes " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStocktakeInitiatedByAccount"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStocktake{}
	if item != nil {
		items = item.([]*StoreStocktake)
	}

	return
}

func (r *GeneratedAccountResolver) CreatedStocktakesIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*StoreStocktake{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["InitiatedByAccountAndStoreStocktakeIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStocktake)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedAccountResolver) PostedStocktakes(ctx context.Context, obj *Account) (res []*StoreStocktake, err error) {
	return r.Handlers.AccountPostedStocktakes(ctx, r.GeneratedResolver, obj)
}
func AccountPostedStocktakesHandler(ctx context.Context, r *GeneratedResolver, obj *Account) (items []*StoreStocktake, err error) {

	items = []*StoreStocktake{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "PostedStocktakes"); err != nil {
		return items, errors.New("PostedStocktakes " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStocktakePostedBy"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStocktake{}
	if item != nil {
		items = item.([]*StoreStocktake)
	}

	return
}

func (r *GeneratedAccountResolver) PostedStocktakesIds(ctx context.Context, obj *Account) (ids []string, err error) {

	items := []*StoreStocktake{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["PostedByAndStoreStocktakeIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStocktake)
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

func (r *GeneratedOrganizationResolver) CustomerMembers(ctx context.Context, obj *Organization) (res []*CustomerMember, err error) {
	return r.Handlers.OrganizationCustomerMembers(ctx, r.GeneratedResolver, obj)
}
func OrganizationCustomerMembersHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*CustomerMember, err error) {

	items = []*CustomerMember{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerMembers"); err != nil {
		return items, errors.New("CustomerMembers " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerMemberOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerMember{}
	if item != nil {
		items = item.([]*CustomerMember)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerMembersIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*CustomerMember{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndCustomerMemberIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerMember)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerBenefitPolicies(ctx context.Context, obj *Organization) (res []*CustomerBenefitPolicy, err error) {
	return r.Handlers.OrganizationCustomerBenefitPolicies(ctx, r.GeneratedResolver, obj)
}
func OrganizationCustomerBenefitPoliciesHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*CustomerBenefitPolicy, err error) {

	items = []*CustomerBenefitPolicy{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerBenefitPolicies"); err != nil {
		return items, errors.New("CustomerBenefitPolicies " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerBenefitPolicyOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerBenefitPolicy{}
	if item != nil {
		items = item.([]*CustomerBenefitPolicy)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerBenefitPoliciesIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*CustomerBenefitPolicy{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndCustomerBenefitPolicyIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerBenefitPolicy)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerDailyPointGrantBudgets(ctx context.Context, obj *Organization) (res []*CustomerDailyPointGrantBudget, err error) {
	return r.Handlers.OrganizationCustomerDailyPointGrantBudgets(ctx, r.GeneratedResolver, obj)
}
func OrganizationCustomerDailyPointGrantBudgetsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*CustomerDailyPointGrantBudget, err error) {

	items = []*CustomerDailyPointGrantBudget{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerDailyPointGrantBudgets"); err != nil {
		return items, errors.New("CustomerDailyPointGrantBudgets " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerDailyPointGrantBudgetOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerDailyPointGrantBudget{}
	if item != nil {
		items = item.([]*CustomerDailyPointGrantBudget)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerDailyPointGrantBudgetsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*CustomerDailyPointGrantBudget{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndCustomerDailyPointGrantBudgetIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerDailyPointGrantBudget)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerPointEntries(ctx context.Context, obj *Organization) (res []*CustomerPointEntry, err error) {
	return r.Handlers.OrganizationCustomerPointEntries(ctx, r.GeneratedResolver, obj)
}
func OrganizationCustomerPointEntriesHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*CustomerPointEntry, err error) {

	items = []*CustomerPointEntry{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerPointEntries"); err != nil {
		return items, errors.New("CustomerPointEntries " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerPointEntrySourceOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerPointEntry{}
	if item != nil {
		items = item.([]*CustomerPointEntry)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerPointEntriesIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*CustomerPointEntry{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SourceOrganizationAndCustomerPointEntryIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerPointEntry)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerCouponTemplates(ctx context.Context, obj *Organization) (res []*CustomerCouponTemplate, err error) {
	return r.Handlers.OrganizationCustomerCouponTemplates(ctx, r.GeneratedResolver, obj)
}
func OrganizationCustomerCouponTemplatesHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*CustomerCouponTemplate, err error) {

	items = []*CustomerCouponTemplate{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerCouponTemplates"); err != nil {
		return items, errors.New("CustomerCouponTemplates " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerCouponTemplateOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerCouponTemplate{}
	if item != nil {
		items = item.([]*CustomerCouponTemplate)
	}

	return
}

func (r *GeneratedOrganizationResolver) CustomerCouponTemplatesIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*CustomerCouponTemplate{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndCustomerCouponTemplateIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerCouponTemplate)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) ProductCategories(ctx context.Context, obj *Organization) (res []*ProductCategory, err error) {
	return r.Handlers.OrganizationProductCategories(ctx, r.GeneratedResolver, obj)
}
func OrganizationProductCategoriesHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*ProductCategory, err error) {

	items = []*ProductCategory{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductCategories"); err != nil {
		return items, errors.New("ProductCategories " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductCategoryOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductCategory{}
	if item != nil {
		items = item.([]*ProductCategory)
	}

	return
}

func (r *GeneratedOrganizationResolver) ProductCategoriesIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*ProductCategory{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndProductCategoryIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductCategory)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) ProductBrands(ctx context.Context, obj *Organization) (res []*ProductBrand, err error) {
	return r.Handlers.OrganizationProductBrands(ctx, r.GeneratedResolver, obj)
}
func OrganizationProductBrandsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*ProductBrand, err error) {

	items = []*ProductBrand{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductBrands"); err != nil {
		return items, errors.New("ProductBrands " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductBrandOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductBrand{}
	if item != nil {
		items = item.([]*ProductBrand)
	}

	return
}

func (r *GeneratedOrganizationResolver) ProductBrandsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*ProductBrand{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndProductBrandIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductBrand)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) SpecificationDefinitions(ctx context.Context, obj *Organization) (res []*SpecificationDefinition, err error) {
	return r.Handlers.OrganizationSpecificationDefinitions(ctx, r.GeneratedResolver, obj)
}
func OrganizationSpecificationDefinitionsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*SpecificationDefinition, err error) {

	items = []*SpecificationDefinition{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SpecificationDefinitions"); err != nil {
		return items, errors.New("SpecificationDefinitions " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SpecificationDefinitionOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*SpecificationDefinition{}
	if item != nil {
		items = item.([]*SpecificationDefinition)
	}

	return
}

func (r *GeneratedOrganizationResolver) SpecificationDefinitionsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*SpecificationDefinition{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndSpecificationDefinitionIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*SpecificationDefinition)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) ProductPackageTemplates(ctx context.Context, obj *Organization) (res []*ProductPackageTemplate, err error) {
	return r.Handlers.OrganizationProductPackageTemplates(ctx, r.GeneratedResolver, obj)
}
func OrganizationProductPackageTemplatesHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*ProductPackageTemplate, err error) {

	items = []*ProductPackageTemplate{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackageTemplates"); err != nil {
		return items, errors.New("ProductPackageTemplates " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductPackageTemplateOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductPackageTemplate{}
	if item != nil {
		items = item.([]*ProductPackageTemplate)
	}

	return
}

func (r *GeneratedOrganizationResolver) ProductPackageTemplatesIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*ProductPackageTemplate{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndProductPackageTemplateIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductPackageTemplate)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedOrganizationResolver) Products(ctx context.Context, obj *Organization) (res []*Product, err error) {
	return r.Handlers.OrganizationProducts(ctx, r.GeneratedResolver, obj)
}
func OrganizationProductsHandler(ctx context.Context, r *GeneratedResolver, obj *Organization) (items []*Product, err error) {

	items = []*Product{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Products"); err != nil {
		return items, errors.New("Products " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductOrganization"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Product{}
	if item != nil {
		items = item.([]*Product)
	}

	return
}

func (r *GeneratedOrganizationResolver) ProductsIds(ctx context.Context, obj *Organization) (ids []string, err error) {

	items := []*Product{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OrganizationAndProductIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Product)
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

func (r *GeneratedStoreResolver) ProductListings(ctx context.Context, obj *Store) (res []*StoreListing, err error) {
	return r.Handlers.StoreProductListings(ctx, r.GeneratedResolver, obj)
}
func StoreProductListingsHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*StoreListing, err error) {

	items = []*StoreListing{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductListings"); err != nil {
		return items, errors.New("ProductListings " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreListingStore"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreListing{}
	if item != nil {
		items = item.([]*StoreListing)
	}

	return
}

func (r *GeneratedStoreResolver) ProductListingsIds(ctx context.Context, obj *Store) (ids []string, err error) {

	items := []*StoreListing{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreAndStoreListingIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreListing)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreResolver) StockMovements(ctx context.Context, obj *Store) (res []*StoreStockMovement, err error) {
	return r.Handlers.StoreStockMovements(ctx, r.GeneratedResolver, obj)
}
func StoreStockMovementsHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*StoreStockMovement, err error) {

	items = []*StoreStockMovement{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StockMovements"); err != nil {
		return items, errors.New("StockMovements " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStockMovementStore"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStockMovement{}
	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	return
}

func (r *GeneratedStoreResolver) StockMovementsIds(ctx context.Context, obj *Store) (ids []string, err error) {

	items := []*StoreStockMovement{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreAndStoreStockMovementIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreResolver) Stocktakes(ctx context.Context, obj *Store) (res []*StoreStocktake, err error) {
	return r.Handlers.StoreStocktakes(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakesHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*StoreStocktake, err error) {

	items = []*StoreStocktake{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Stocktakes"); err != nil {
		return items, errors.New("Stocktakes " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStocktakeStore"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStocktake{}
	if item != nil {
		items = item.([]*StoreStocktake)
	}

	return
}

func (r *GeneratedStoreResolver) StocktakesIds(ctx context.Context, obj *Store) (ids []string, err error) {

	items := []*StoreStocktake{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreAndStoreStocktakeIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStocktake)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreResolver) Promotions(ctx context.Context, obj *Store) (res []*StorePromotion, err error) {
	return r.Handlers.StorePromotions(ctx, r.GeneratedResolver, obj)
}
func StorePromotionsHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*StorePromotion, err error) {

	items = []*StorePromotion{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Promotions"); err != nil {
		return items, errors.New("Promotions " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StorePromotionStore"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StorePromotion{}
	if item != nil {
		items = item.([]*StorePromotion)
	}

	return
}

func (r *GeneratedStoreResolver) PromotionsIds(ctx context.Context, obj *Store) (ids []string, err error) {

	items := []*StorePromotion{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreAndStorePromotionIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StorePromotion)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreResolver) CouponTemplates(ctx context.Context, obj *Store) (res []*CustomerCouponTemplate, err error) {
	return r.Handlers.StoreCouponTemplates(ctx, r.GeneratedResolver, obj)
}
func StoreCouponTemplatesHandler(ctx context.Context, r *GeneratedResolver, obj *Store) (items []*CustomerCouponTemplate, err error) {

	items = []*CustomerCouponTemplate{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CouponTemplates"); err != nil {
		return items, errors.New("CouponTemplates " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerCouponTemplateApplicableStore"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerCouponTemplate{}
	if item != nil {
		items = item.([]*CustomerCouponTemplate)
	}

	return
}

func (r *GeneratedStoreResolver) CouponTemplatesIds(ctx context.Context, obj *Store) (ids []string, err error) {

	items := []*CustomerCouponTemplate{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ApplicableStoreAndCustomerCouponTemplateIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerCouponTemplate)
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

type QueryCustomerMemberHandlerOptions struct {
	ID     *string
	Filter *CustomerMemberFilterType
}

func (r *GeneratedQueryResolver) CustomerMember(ctx context.Context, id *string, filter *CustomerMemberFilterType) (*CustomerMember, error) {
	opts := QueryCustomerMemberHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryCustomerMember(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerMemberHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerMemberHandlerOptions) (*CustomerMember, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := CustomerMemberQueryFilter{}
	rt := &CustomerMemberResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("customer_members", ctx)+".id = ?", *opts.ID)
	}

	var items []*CustomerMember
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_members", ctx),
		Preloaders: []string{},
		Item:       &CustomerMember{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "CustomerMember"}
	}
	return items[0], err
}

type QueryCustomerMembersHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*CustomerMemberSortType
	Filter      *CustomerMemberFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) CustomerMembers(ctx context.Context, current_page *int, per_page *int, q *string, sort []*CustomerMemberSortType, filter *CustomerMemberFilterType, rand *bool) (*CustomerMemberResultType, error) {
	opts := QueryCustomerMembersHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryCustomerMembers(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerMembersHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerMembersHandlerOptions) (*CustomerMemberResultType, error) {
	query := CustomerMemberQueryFilter{opts.Q}

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

	return &CustomerMemberResultType{
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

type GeneratedCustomerMemberResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerMemberResultTypeResolver) Data(ctx context.Context, obj *CustomerMemberResultType) (items []*CustomerMember, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_members", ctx),
		Preloaders: []string{},
		Item:       &CustomerMember{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*CustomerMember{}
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

func (r *GeneratedCustomerMemberResultTypeResolver) Total(ctx context.Context, obj *CustomerMemberResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("customer_members", ctx), &CustomerMember{})
}

func (r *GeneratedCustomerMemberResultTypeResolver) TotalPage(ctx context.Context, obj *CustomerMemberResultType) (count int, err error) {
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

func (r *GeneratedCustomerMemberResultTypeResolver) CurrentPage(ctx context.Context, obj *CustomerMemberResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedCustomerMemberResultTypeResolver) PerPage(ctx context.Context, obj *CustomerMemberResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedCustomerMemberResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerMemberResolver) Organization(ctx context.Context, obj *CustomerMember) (res *Organization, err error) {
	return r.Handlers.CustomerMemberOrganization(ctx, r.GeneratedResolver, obj)
}
func CustomerMemberOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (items *Organization, err error) {

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

func (r *GeneratedCustomerMemberResolver) PointEntries(ctx context.Context, obj *CustomerMember) (res []*CustomerPointEntry, err error) {
	return r.Handlers.CustomerMemberPointEntries(ctx, r.GeneratedResolver, obj)
}
func CustomerMemberPointEntriesHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (items []*CustomerPointEntry, err error) {

	items = []*CustomerPointEntry{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "PointEntries"); err != nil {
		return items, errors.New("PointEntries " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerPointEntryMember"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerPointEntry{}
	if item != nil {
		items = item.([]*CustomerPointEntry)
	}

	return
}

func (r *GeneratedCustomerMemberResolver) PointEntriesIds(ctx context.Context, obj *CustomerMember) (ids []string, err error) {

	items := []*CustomerPointEntry{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["MemberAndCustomerPointEntryIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerPointEntry)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedCustomerMemberResolver) CouponGrants(ctx context.Context, obj *CustomerMember) (res []*CustomerCouponGrant, err error) {
	return r.Handlers.CustomerMemberCouponGrants(ctx, r.GeneratedResolver, obj)
}
func CustomerMemberCouponGrantsHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (items []*CustomerCouponGrant, err error) {

	items = []*CustomerCouponGrant{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CouponGrants"); err != nil {
		return items, errors.New("CouponGrants " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerCouponGrantMember"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerCouponGrant{}
	if item != nil {
		items = item.([]*CustomerCouponGrant)
	}

	return
}

func (r *GeneratedCustomerMemberResolver) CouponGrantsIds(ctx context.Context, obj *CustomerMember) (ids []string, err error) {

	items := []*CustomerCouponGrant{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["MemberAndCustomerCouponGrantIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerCouponGrant)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedCustomerMemberResolver) CouponDistributionJobs(ctx context.Context, obj *CustomerMember) (res []*CustomerCouponDistributionJob, err error) {
	return r.Handlers.CustomerMemberCouponDistributionJobs(ctx, r.GeneratedResolver, obj)
}
func CustomerMemberCouponDistributionJobsHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerMember) (items []*CustomerCouponDistributionJob, err error) {

	items = []*CustomerCouponDistributionJob{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CouponDistributionJobs"); err != nil {
		return items, errors.New("CouponDistributionJobs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerCouponDistributionJobMember"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerCouponDistributionJob{}
	if item != nil {
		items = item.([]*CustomerCouponDistributionJob)
	}

	return
}

func (r *GeneratedCustomerMemberResolver) CouponDistributionJobsIds(ctx context.Context, obj *CustomerMember) (ids []string, err error) {
	if _, err := auth.RequirePrincipal(ctx); err != nil {
		return nil, err
	}
	if denied := auth.NewError(auth.CodePermissionDenied); denied != nil {
		return nil, denied
	}

	items := []*CustomerCouponDistributionJob{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["MemberAndCustomerCouponDistributionJobIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerCouponDistributionJob)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryCustomerBenefitPolicyHandlerOptions struct {
	ID     *string
	Filter *CustomerBenefitPolicyFilterType
}

func (r *GeneratedQueryResolver) CustomerBenefitPolicy(ctx context.Context, id *string, filter *CustomerBenefitPolicyFilterType) (*CustomerBenefitPolicy, error) {
	opts := QueryCustomerBenefitPolicyHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryCustomerBenefitPolicy(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerBenefitPolicyHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerBenefitPolicyHandlerOptions) (*CustomerBenefitPolicy, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := CustomerBenefitPolicyQueryFilter{}
	rt := &CustomerBenefitPolicyResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("customer_benefit_policies", ctx)+".id = ?", *opts.ID)
	}

	var items []*CustomerBenefitPolicy
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_benefit_policies", ctx),
		Preloaders: []string{},
		Item:       &CustomerBenefitPolicy{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "CustomerBenefitPolicy"}
	}
	return items[0], err
}

type QueryCustomerBenefitPoliciesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*CustomerBenefitPolicySortType
	Filter      *CustomerBenefitPolicyFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) CustomerBenefitPolicies(ctx context.Context, current_page *int, per_page *int, q *string, sort []*CustomerBenefitPolicySortType, filter *CustomerBenefitPolicyFilterType, rand *bool) (*CustomerBenefitPolicyResultType, error) {
	opts := QueryCustomerBenefitPoliciesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryCustomerBenefitPolicies(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerBenefitPoliciesHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerBenefitPoliciesHandlerOptions) (*CustomerBenefitPolicyResultType, error) {
	query := CustomerBenefitPolicyQueryFilter{opts.Q}

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

	return &CustomerBenefitPolicyResultType{
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

type GeneratedCustomerBenefitPolicyResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerBenefitPolicyResultTypeResolver) Data(ctx context.Context, obj *CustomerBenefitPolicyResultType) (items []*CustomerBenefitPolicy, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_benefit_policies", ctx),
		Preloaders: []string{},
		Item:       &CustomerBenefitPolicy{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*CustomerBenefitPolicy{}
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

func (r *GeneratedCustomerBenefitPolicyResultTypeResolver) Total(ctx context.Context, obj *CustomerBenefitPolicyResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("customer_benefit_policies", ctx), &CustomerBenefitPolicy{})
}

func (r *GeneratedCustomerBenefitPolicyResultTypeResolver) TotalPage(ctx context.Context, obj *CustomerBenefitPolicyResultType) (count int, err error) {
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

func (r *GeneratedCustomerBenefitPolicyResultTypeResolver) CurrentPage(ctx context.Context, obj *CustomerBenefitPolicyResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedCustomerBenefitPolicyResultTypeResolver) PerPage(ctx context.Context, obj *CustomerBenefitPolicyResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedCustomerBenefitPolicyResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerBenefitPolicyResolver) Organization(ctx context.Context, obj *CustomerBenefitPolicy) (res *Organization, err error) {
	return r.Handlers.CustomerBenefitPolicyOrganization(ctx, r.GeneratedResolver, obj)
}
func CustomerBenefitPolicyOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerBenefitPolicy) (items *Organization, err error) {

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

type QueryCustomerDailyPointGrantBudgetHandlerOptions struct {
	ID     *string
	Filter *CustomerDailyPointGrantBudgetFilterType
}

func (r *GeneratedQueryResolver) CustomerDailyPointGrantBudget(ctx context.Context, id *string, filter *CustomerDailyPointGrantBudgetFilterType) (*CustomerDailyPointGrantBudget, error) {
	opts := QueryCustomerDailyPointGrantBudgetHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryCustomerDailyPointGrantBudget(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerDailyPointGrantBudgetHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerDailyPointGrantBudgetHandlerOptions) (*CustomerDailyPointGrantBudget, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := CustomerDailyPointGrantBudgetQueryFilter{}
	rt := &CustomerDailyPointGrantBudgetResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("customer_daily_point_grant_budgets", ctx)+".id = ?", *opts.ID)
	}

	var items []*CustomerDailyPointGrantBudget
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_daily_point_grant_budgets", ctx),
		Preloaders: []string{},
		Item:       &CustomerDailyPointGrantBudget{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "CustomerDailyPointGrantBudget"}
	}
	return items[0], err
}

type QueryCustomerDailyPointGrantBudgetsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*CustomerDailyPointGrantBudgetSortType
	Filter      *CustomerDailyPointGrantBudgetFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) CustomerDailyPointGrantBudgets(ctx context.Context, current_page *int, per_page *int, q *string, sort []*CustomerDailyPointGrantBudgetSortType, filter *CustomerDailyPointGrantBudgetFilterType, rand *bool) (*CustomerDailyPointGrantBudgetResultType, error) {
	opts := QueryCustomerDailyPointGrantBudgetsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryCustomerDailyPointGrantBudgets(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerDailyPointGrantBudgetsHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerDailyPointGrantBudgetsHandlerOptions) (*CustomerDailyPointGrantBudgetResultType, error) {
	query := CustomerDailyPointGrantBudgetQueryFilter{opts.Q}

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

	return &CustomerDailyPointGrantBudgetResultType{
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

type GeneratedCustomerDailyPointGrantBudgetResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerDailyPointGrantBudgetResultTypeResolver) Data(ctx context.Context, obj *CustomerDailyPointGrantBudgetResultType) (items []*CustomerDailyPointGrantBudget, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_daily_point_grant_budgets", ctx),
		Preloaders: []string{},
		Item:       &CustomerDailyPointGrantBudget{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*CustomerDailyPointGrantBudget{}
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

func (r *GeneratedCustomerDailyPointGrantBudgetResultTypeResolver) Total(ctx context.Context, obj *CustomerDailyPointGrantBudgetResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("customer_daily_point_grant_budgets", ctx), &CustomerDailyPointGrantBudget{})
}

func (r *GeneratedCustomerDailyPointGrantBudgetResultTypeResolver) TotalPage(ctx context.Context, obj *CustomerDailyPointGrantBudgetResultType) (count int, err error) {
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

func (r *GeneratedCustomerDailyPointGrantBudgetResultTypeResolver) CurrentPage(ctx context.Context, obj *CustomerDailyPointGrantBudgetResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedCustomerDailyPointGrantBudgetResultTypeResolver) PerPage(ctx context.Context, obj *CustomerDailyPointGrantBudgetResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedCustomerDailyPointGrantBudgetResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerDailyPointGrantBudgetResolver) Organization(ctx context.Context, obj *CustomerDailyPointGrantBudget) (res *Organization, err error) {
	return r.Handlers.CustomerDailyPointGrantBudgetOrganization(ctx, r.GeneratedResolver, obj)
}
func CustomerDailyPointGrantBudgetOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerDailyPointGrantBudget) (items *Organization, err error) {

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

type QueryCustomerPointEntryHandlerOptions struct {
	ID     *string
	Filter *CustomerPointEntryFilterType
}

func (r *GeneratedQueryResolver) CustomerPointEntry(ctx context.Context, id *string, filter *CustomerPointEntryFilterType) (*CustomerPointEntry, error) {
	opts := QueryCustomerPointEntryHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryCustomerPointEntry(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerPointEntryHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerPointEntryHandlerOptions) (*CustomerPointEntry, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := CustomerPointEntryQueryFilter{}
	rt := &CustomerPointEntryResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("customer_point_entries", ctx)+".id = ?", *opts.ID)
	}

	var items []*CustomerPointEntry
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_point_entries", ctx),
		Preloaders: []string{},
		Item:       &CustomerPointEntry{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "CustomerPointEntry"}
	}
	return items[0], err
}

type QueryCustomerPointEntriesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*CustomerPointEntrySortType
	Filter      *CustomerPointEntryFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) CustomerPointEntries(ctx context.Context, current_page *int, per_page *int, q *string, sort []*CustomerPointEntrySortType, filter *CustomerPointEntryFilterType, rand *bool) (*CustomerPointEntryResultType, error) {
	opts := QueryCustomerPointEntriesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryCustomerPointEntries(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerPointEntriesHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerPointEntriesHandlerOptions) (*CustomerPointEntryResultType, error) {
	query := CustomerPointEntryQueryFilter{opts.Q}

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

	return &CustomerPointEntryResultType{
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

type GeneratedCustomerPointEntryResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerPointEntryResultTypeResolver) Data(ctx context.Context, obj *CustomerPointEntryResultType) (items []*CustomerPointEntry, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_point_entries", ctx),
		Preloaders: []string{},
		Item:       &CustomerPointEntry{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*CustomerPointEntry{}
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

func (r *GeneratedCustomerPointEntryResultTypeResolver) Total(ctx context.Context, obj *CustomerPointEntryResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("customer_point_entries", ctx), &CustomerPointEntry{})
}

func (r *GeneratedCustomerPointEntryResultTypeResolver) TotalPage(ctx context.Context, obj *CustomerPointEntryResultType) (count int, err error) {
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

func (r *GeneratedCustomerPointEntryResultTypeResolver) CurrentPage(ctx context.Context, obj *CustomerPointEntryResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedCustomerPointEntryResultTypeResolver) PerPage(ctx context.Context, obj *CustomerPointEntryResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedCustomerPointEntryResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerPointEntryResolver) Member(ctx context.Context, obj *CustomerPointEntry) (res *CustomerMember, err error) {
	return r.Handlers.CustomerPointEntryMember(ctx, r.GeneratedResolver, obj)
}
func CustomerPointEntryMemberHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (items *CustomerMember, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerMember"); err != nil {
		return items, errors.New("CustomerMember " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.MemberID

	if objKey != "" {
		item, _ := loaders["CustomerMember"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*CustomerMember)

		if items == nil {
			items = &CustomerMember{}
		}

	}

	return
}

func (r *GeneratedCustomerPointEntryResolver) SourceOrganization(ctx context.Context, obj *CustomerPointEntry) (res *Organization, err error) {
	return r.Handlers.CustomerPointEntrySourceOrganization(ctx, r.GeneratedResolver, obj)
}
func CustomerPointEntrySourceOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (items *Organization, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
		return items, errors.New("Organization " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.SourceOrganizationID

	if objKey != "" {
		item, _ := loaders["Organization"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Organization)

		if items == nil {
			items = &Organization{}
		}

	}

	return
}

func (r *GeneratedCustomerPointEntryResolver) Reverses(ctx context.Context, obj *CustomerPointEntry) (res *CustomerPointEntry, err error) {
	return r.Handlers.CustomerPointEntryReverses(ctx, r.GeneratedResolver, obj)
}
func CustomerPointEntryReversesHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (items *CustomerPointEntry, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerPointEntry"); err != nil {
		return items, errors.New("CustomerPointEntry " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ReversesID

	if objKey != nil {
		item, _ := loaders["CustomerPointEntry"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*CustomerPointEntry)

	}

	return
}

func (r *GeneratedCustomerPointEntryResolver) ReversedBy(ctx context.Context, obj *CustomerPointEntry) (res *CustomerPointEntry, err error) {
	return r.Handlers.CustomerPointEntryReversedBy(ctx, r.GeneratedResolver, obj)
}
func CustomerPointEntryReversedByHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerPointEntry) (items *CustomerPointEntry, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerPointEntry"); err != nil {
		return items, errors.New("CustomerPointEntry " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ReversedByID

	if objKey != nil {
		item, _ := loaders["CustomerPointEntry"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*CustomerPointEntry)

	}

	return
}

type QueryCustomerCouponTemplateHandlerOptions struct {
	ID     *string
	Filter *CustomerCouponTemplateFilterType
}

func (r *GeneratedQueryResolver) CustomerCouponTemplate(ctx context.Context, id *string, filter *CustomerCouponTemplateFilterType) (*CustomerCouponTemplate, error) {
	opts := QueryCustomerCouponTemplateHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryCustomerCouponTemplate(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerCouponTemplateHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponTemplateHandlerOptions) (*CustomerCouponTemplate, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := CustomerCouponTemplateQueryFilter{}
	rt := &CustomerCouponTemplateResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("customer_coupon_templates", ctx)+".id = ?", *opts.ID)
	}

	var items []*CustomerCouponTemplate
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_coupon_templates", ctx),
		Preloaders: []string{},
		Item:       &CustomerCouponTemplate{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "CustomerCouponTemplate"}
	}
	return items[0], err
}

type QueryCustomerCouponTemplatesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*CustomerCouponTemplateSortType
	Filter      *CustomerCouponTemplateFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) CustomerCouponTemplates(ctx context.Context, current_page *int, per_page *int, q *string, sort []*CustomerCouponTemplateSortType, filter *CustomerCouponTemplateFilterType, rand *bool) (*CustomerCouponTemplateResultType, error) {
	opts := QueryCustomerCouponTemplatesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryCustomerCouponTemplates(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerCouponTemplatesHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponTemplatesHandlerOptions) (*CustomerCouponTemplateResultType, error) {
	query := CustomerCouponTemplateQueryFilter{opts.Q}

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

	return &CustomerCouponTemplateResultType{
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

type GeneratedCustomerCouponTemplateResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerCouponTemplateResultTypeResolver) Data(ctx context.Context, obj *CustomerCouponTemplateResultType) (items []*CustomerCouponTemplate, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_coupon_templates", ctx),
		Preloaders: []string{},
		Item:       &CustomerCouponTemplate{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*CustomerCouponTemplate{}
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

func (r *GeneratedCustomerCouponTemplateResultTypeResolver) Total(ctx context.Context, obj *CustomerCouponTemplateResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("customer_coupon_templates", ctx), &CustomerCouponTemplate{})
}

func (r *GeneratedCustomerCouponTemplateResultTypeResolver) TotalPage(ctx context.Context, obj *CustomerCouponTemplateResultType) (count int, err error) {
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

func (r *GeneratedCustomerCouponTemplateResultTypeResolver) CurrentPage(ctx context.Context, obj *CustomerCouponTemplateResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedCustomerCouponTemplateResultTypeResolver) PerPage(ctx context.Context, obj *CustomerCouponTemplateResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedCustomerCouponTemplateResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerCouponTemplateResolver) Organization(ctx context.Context, obj *CustomerCouponTemplate) (res *Organization, err error) {
	return r.Handlers.CustomerCouponTemplateOrganization(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponTemplateOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (items *Organization, err error) {

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

func (r *GeneratedCustomerCouponTemplateResolver) ApplicableStore(ctx context.Context, obj *CustomerCouponTemplate) (res *Store, err error) {
	return r.Handlers.CustomerCouponTemplateApplicableStore(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponTemplateApplicableStoreHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (items *Store, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
		return items, errors.New("Store " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ApplicableStoreID

	if objKey != nil {
		item, _ := loaders["Store"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Store)

	}

	return
}

func (r *GeneratedCustomerCouponTemplateResolver) Grants(ctx context.Context, obj *CustomerCouponTemplate) (res []*CustomerCouponGrant, err error) {
	return r.Handlers.CustomerCouponTemplateGrants(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponTemplateGrantsHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (items []*CustomerCouponGrant, err error) {

	items = []*CustomerCouponGrant{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Grants"); err != nil {
		return items, errors.New("Grants " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerCouponGrantTemplate"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerCouponGrant{}
	if item != nil {
		items = item.([]*CustomerCouponGrant)
	}

	return
}

func (r *GeneratedCustomerCouponTemplateResolver) GrantsIds(ctx context.Context, obj *CustomerCouponTemplate) (ids []string, err error) {

	items := []*CustomerCouponGrant{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["TemplateAndCustomerCouponGrantIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerCouponGrant)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedCustomerCouponTemplateResolver) DistributionJobs(ctx context.Context, obj *CustomerCouponTemplate) (res []*CustomerCouponDistributionJob, err error) {
	return r.Handlers.CustomerCouponTemplateDistributionJobs(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponTemplateDistributionJobsHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponTemplate) (items []*CustomerCouponDistributionJob, err error) {

	items = []*CustomerCouponDistributionJob{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "DistributionJobs"); err != nil {
		return items, errors.New("DistributionJobs " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CustomerCouponDistributionJobTemplate"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*CustomerCouponDistributionJob{}
	if item != nil {
		items = item.([]*CustomerCouponDistributionJob)
	}

	return
}

func (r *GeneratedCustomerCouponTemplateResolver) DistributionJobsIds(ctx context.Context, obj *CustomerCouponTemplate) (ids []string, err error) {
	if _, err := auth.RequirePrincipal(ctx); err != nil {
		return nil, err
	}
	if denied := auth.NewError(auth.CodePermissionDenied); denied != nil {
		return nil, denied
	}

	items := []*CustomerCouponDistributionJob{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["TemplateAndCustomerCouponDistributionJobIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*CustomerCouponDistributionJob)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryProductCategoryHandlerOptions struct {
	ID     *string
	Filter *ProductCategoryFilterType
}

func (r *GeneratedQueryResolver) ProductCategory(ctx context.Context, id *string, filter *ProductCategoryFilterType) (*ProductCategory, error) {
	opts := QueryProductCategoryHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProductCategory(ctx, r.GeneratedResolver, opts)
}
func QueryProductCategoryHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductCategoryHandlerOptions) (*ProductCategory, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductCategoryQueryFilter{}
	rt := &ProductCategoryResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("product_categories", ctx)+".id = ?", *opts.ID)
	}

	var items []*ProductCategory
	giOpts := GetItemsOptions{
		Alias:      TableName("product_categories", ctx),
		Preloaders: []string{},
		Item:       &ProductCategory{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "ProductCategory"}
	}
	return items[0], err
}

type QueryProductCategoriesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductCategorySortType
	Filter      *ProductCategoryFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) ProductCategories(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductCategorySortType, filter *ProductCategoryFilterType, rand *bool) (*ProductCategoryResultType, error) {
	opts := QueryProductCategoriesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProductCategories(ctx, r.GeneratedResolver, opts)
}
func QueryProductCategoriesHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductCategoriesHandlerOptions) (*ProductCategoryResultType, error) {
	query := ProductCategoryQueryFilter{opts.Q}

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

	return &ProductCategoryResultType{
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

type GeneratedProductCategoryResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductCategoryResultTypeResolver) Data(ctx context.Context, obj *ProductCategoryResultType) (items []*ProductCategory, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("product_categories", ctx),
		Preloaders: []string{},
		Item:       &ProductCategory{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*ProductCategory{}
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

func (r *GeneratedProductCategoryResultTypeResolver) Total(ctx context.Context, obj *ProductCategoryResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("product_categories", ctx), &ProductCategory{})
}

func (r *GeneratedProductCategoryResultTypeResolver) TotalPage(ctx context.Context, obj *ProductCategoryResultType) (count int, err error) {
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

func (r *GeneratedProductCategoryResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductCategoryResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductCategoryResultTypeResolver) PerPage(ctx context.Context, obj *ProductCategoryResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductCategoryResolver struct{ *GeneratedResolver }

func (r *GeneratedProductCategoryResolver) Organization(ctx context.Context, obj *ProductCategory) (res *Organization, err error) {
	return r.Handlers.ProductCategoryOrganization(ctx, r.GeneratedResolver, obj)
}
func ProductCategoryOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (items *Organization, err error) {

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

func (r *GeneratedProductCategoryResolver) Parent(ctx context.Context, obj *ProductCategory) (res *ProductCategory, err error) {
	return r.Handlers.ProductCategoryParent(ctx, r.GeneratedResolver, obj)
}
func ProductCategoryParentHandler(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (items *ProductCategory, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductCategory"); err != nil {
		return items, errors.New("ProductCategory " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ParentID

	if objKey != nil {
		item, _ := loaders["ProductCategory"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductCategory)

	}

	return
}

func (r *GeneratedProductCategoryResolver) Children(ctx context.Context, obj *ProductCategory) (res []*ProductCategory, err error) {
	return r.Handlers.ProductCategoryChildren(ctx, r.GeneratedResolver, obj)
}
func ProductCategoryChildrenHandler(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (items []*ProductCategory, err error) {

	items = []*ProductCategory{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Children"); err != nil {
		return items, errors.New("Children " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductCategoryParent"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductCategory{}
	if item != nil {
		items = item.([]*ProductCategory)
	}

	return
}

func (r *GeneratedProductCategoryResolver) ChildrenIds(ctx context.Context, obj *ProductCategory) (ids []string, err error) {

	items := []*ProductCategory{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ParentAndProductCategoryIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductCategory)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductCategoryResolver) Products(ctx context.Context, obj *ProductCategory) (res []*Product, err error) {
	return r.Handlers.ProductCategoryProducts(ctx, r.GeneratedResolver, obj)
}
func ProductCategoryProductsHandler(ctx context.Context, r *GeneratedResolver, obj *ProductCategory) (items []*Product, err error) {

	items = []*Product{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Products"); err != nil {
		return items, errors.New("Products " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductCategory"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Product{}
	if item != nil {
		items = item.([]*Product)
	}

	return
}

func (r *GeneratedProductCategoryResolver) ProductsIds(ctx context.Context, obj *ProductCategory) (ids []string, err error) {

	items := []*Product{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["CategoryAndProductIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Product)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryProductBrandHandlerOptions struct {
	ID     *string
	Filter *ProductBrandFilterType
}

func (r *GeneratedQueryResolver) ProductBrand(ctx context.Context, id *string, filter *ProductBrandFilterType) (*ProductBrand, error) {
	opts := QueryProductBrandHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProductBrand(ctx, r.GeneratedResolver, opts)
}
func QueryProductBrandHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductBrandHandlerOptions) (*ProductBrand, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductBrandQueryFilter{}
	rt := &ProductBrandResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("product_brands", ctx)+".id = ?", *opts.ID)
	}

	var items []*ProductBrand
	giOpts := GetItemsOptions{
		Alias:      TableName("product_brands", ctx),
		Preloaders: []string{},
		Item:       &ProductBrand{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "ProductBrand"}
	}
	return items[0], err
}

type QueryProductBrandsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductBrandSortType
	Filter      *ProductBrandFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) ProductBrands(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductBrandSortType, filter *ProductBrandFilterType, rand *bool) (*ProductBrandResultType, error) {
	opts := QueryProductBrandsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProductBrands(ctx, r.GeneratedResolver, opts)
}
func QueryProductBrandsHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductBrandsHandlerOptions) (*ProductBrandResultType, error) {
	query := ProductBrandQueryFilter{opts.Q}

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

	return &ProductBrandResultType{
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

type GeneratedProductBrandResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductBrandResultTypeResolver) Data(ctx context.Context, obj *ProductBrandResultType) (items []*ProductBrand, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("product_brands", ctx),
		Preloaders: []string{},
		Item:       &ProductBrand{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*ProductBrand{}
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

func (r *GeneratedProductBrandResultTypeResolver) Total(ctx context.Context, obj *ProductBrandResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("product_brands", ctx), &ProductBrand{})
}

func (r *GeneratedProductBrandResultTypeResolver) TotalPage(ctx context.Context, obj *ProductBrandResultType) (count int, err error) {
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

func (r *GeneratedProductBrandResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductBrandResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductBrandResultTypeResolver) PerPage(ctx context.Context, obj *ProductBrandResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductBrandResolver struct{ *GeneratedResolver }

func (r *GeneratedProductBrandResolver) Organization(ctx context.Context, obj *ProductBrand) (res *Organization, err error) {
	return r.Handlers.ProductBrandOrganization(ctx, r.GeneratedResolver, obj)
}
func ProductBrandOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *ProductBrand) (items *Organization, err error) {

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

func (r *GeneratedProductBrandResolver) Products(ctx context.Context, obj *ProductBrand) (res []*Product, err error) {
	return r.Handlers.ProductBrandProducts(ctx, r.GeneratedResolver, obj)
}
func ProductBrandProductsHandler(ctx context.Context, r *GeneratedResolver, obj *ProductBrand) (items []*Product, err error) {

	items = []*Product{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Products"); err != nil {
		return items, errors.New("Products " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductBrand"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Product{}
	if item != nil {
		items = item.([]*Product)
	}

	return
}

func (r *GeneratedProductBrandResolver) ProductsIds(ctx context.Context, obj *ProductBrand) (ids []string, err error) {

	items := []*Product{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["BrandAndProductIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Product)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryProductHandlerOptions struct {
	ID     *string
	Filter *ProductFilterType
}

func (r *GeneratedQueryResolver) Product(ctx context.Context, id *string, filter *ProductFilterType) (*Product, error) {
	opts := QueryProductHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProduct(ctx, r.GeneratedResolver, opts)
}
func QueryProductHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductHandlerOptions) (*Product, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductQueryFilter{}
	rt := &ProductResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("products", ctx)+".id = ?", *opts.ID)
	}

	var items []*Product
	giOpts := GetItemsOptions{
		Alias:      TableName("products", ctx),
		Preloaders: []string{},
		Item:       &Product{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "Product"}
	}
	return items[0], err
}

type QueryProductsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductSortType
	Filter      *ProductFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) Products(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductSortType, filter *ProductFilterType, rand *bool) (*ProductResultType, error) {
	opts := QueryProductsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProducts(ctx, r.GeneratedResolver, opts)
}
func QueryProductsHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductsHandlerOptions) (*ProductResultType, error) {
	query := ProductQueryFilter{opts.Q}

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

	return &ProductResultType{
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

type GeneratedProductResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductResultTypeResolver) Data(ctx context.Context, obj *ProductResultType) (items []*Product, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("products", ctx),
		Preloaders: []string{},
		Item:       &Product{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*Product{}
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

func (r *GeneratedProductResultTypeResolver) Total(ctx context.Context, obj *ProductResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("products", ctx), &Product{})
}

func (r *GeneratedProductResultTypeResolver) TotalPage(ctx context.Context, obj *ProductResultType) (count int, err error) {
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

func (r *GeneratedProductResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductResultTypeResolver) PerPage(ctx context.Context, obj *ProductResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductResolver struct{ *GeneratedResolver }

func (r *GeneratedProductResolver) Brand(ctx context.Context, obj *Product) (res *ProductBrand, err error) {
	return r.Handlers.ProductBrand(ctx, r.GeneratedResolver, obj)
}
func ProductBrandHandler(ctx context.Context, r *GeneratedResolver, obj *Product) (items *ProductBrand, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductBrand"); err != nil {
		return items, errors.New("ProductBrand " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.BrandID

	if objKey != nil {
		item, _ := loaders["ProductBrand"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductBrand)

	}

	return
}

func (r *GeneratedProductResolver) Organization(ctx context.Context, obj *Product) (res *Organization, err error) {
	return r.Handlers.ProductOrganization(ctx, r.GeneratedResolver, obj)
}
func ProductOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *Product) (items *Organization, err error) {

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

func (r *GeneratedProductResolver) Category(ctx context.Context, obj *Product) (res *ProductCategory, err error) {
	return r.Handlers.ProductCategory(ctx, r.GeneratedResolver, obj)
}
func ProductCategoryHandler(ctx context.Context, r *GeneratedResolver, obj *Product) (items *ProductCategory, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductCategory"); err != nil {
		return items, errors.New("ProductCategory " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.CategoryID

	if objKey != "" {
		item, _ := loaders["ProductCategory"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*ProductCategory)

		if items == nil {
			items = &ProductCategory{}
		}

	}

	return
}

func (r *GeneratedProductResolver) DefaultPackageTemplate(ctx context.Context, obj *Product) (res *ProductPackageTemplate, err error) {
	return r.Handlers.ProductDefaultPackageTemplate(ctx, r.GeneratedResolver, obj)
}
func ProductDefaultPackageTemplateHandler(ctx context.Context, r *GeneratedResolver, obj *Product) (items *ProductPackageTemplate, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackageTemplate"); err != nil {
		return items, errors.New("ProductPackageTemplate " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.DefaultPackageTemplateID

	if objKey != nil {
		item, _ := loaders["ProductPackageTemplate"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductPackageTemplate)

	}

	return
}

func (r *GeneratedProductResolver) SpecificationChoices(ctx context.Context, obj *Product) (res []*ProductSpecificationChoice, err error) {
	return r.Handlers.ProductSpecificationChoices(ctx, r.GeneratedResolver, obj)
}
func ProductSpecificationChoicesHandler(ctx context.Context, r *GeneratedResolver, obj *Product) (items []*ProductSpecificationChoice, err error) {

	items = []*ProductSpecificationChoice{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SpecificationChoices"); err != nil {
		return items, errors.New("SpecificationChoices " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductSpecificationChoiceProduct"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductSpecificationChoice{}
	if item != nil {
		items = item.([]*ProductSpecificationChoice)
	}

	return
}

func (r *GeneratedProductResolver) SpecificationChoicesIds(ctx context.Context, obj *Product) (ids []string, err error) {

	items := []*ProductSpecificationChoice{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductAndProductSpecificationChoiceIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductSpecificationChoice)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductResolver) Skus(ctx context.Context, obj *Product) (res []*ProductSku, err error) {
	return r.Handlers.ProductSkus(ctx, r.GeneratedResolver, obj)
}
func ProductSkusHandler(ctx context.Context, r *GeneratedResolver, obj *Product) (items []*ProductSku, err error) {

	items = []*ProductSku{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Skus"); err != nil {
		return items, errors.New("Skus " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductSkuProduct"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductSku{}
	if item != nil {
		items = item.([]*ProductSku)
	}

	return
}

func (r *GeneratedProductResolver) SkusIds(ctx context.Context, obj *Product) (ids []string, err error) {

	items := []*ProductSku{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductAndProductSkuIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductSku)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryProductSkuHandlerOptions struct {
	ID     *string
	Filter *ProductSkuFilterType
}

func (r *GeneratedQueryResolver) ProductSku(ctx context.Context, id *string, filter *ProductSkuFilterType) (*ProductSku, error) {
	opts := QueryProductSkuHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProductSku(ctx, r.GeneratedResolver, opts)
}
func QueryProductSkuHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductSkuHandlerOptions) (*ProductSku, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductSkuQueryFilter{}
	rt := &ProductSkuResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("product_skus", ctx)+".id = ?", *opts.ID)
	}

	var items []*ProductSku
	giOpts := GetItemsOptions{
		Alias:      TableName("product_skus", ctx),
		Preloaders: []string{},
		Item:       &ProductSku{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "ProductSku"}
	}
	return items[0], err
}

type QueryProductSkusHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductSkuSortType
	Filter      *ProductSkuFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) ProductSkus(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductSkuSortType, filter *ProductSkuFilterType, rand *bool) (*ProductSkuResultType, error) {
	opts := QueryProductSkusHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProductSkus(ctx, r.GeneratedResolver, opts)
}
func QueryProductSkusHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductSkusHandlerOptions) (*ProductSkuResultType, error) {
	query := ProductSkuQueryFilter{opts.Q}

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

	return &ProductSkuResultType{
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

type GeneratedProductSkuResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductSkuResultTypeResolver) Data(ctx context.Context, obj *ProductSkuResultType) (items []*ProductSku, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("product_skus", ctx),
		Preloaders: []string{},
		Item:       &ProductSku{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*ProductSku{}
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

func (r *GeneratedProductSkuResultTypeResolver) Total(ctx context.Context, obj *ProductSkuResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("product_skus", ctx), &ProductSku{})
}

func (r *GeneratedProductSkuResultTypeResolver) TotalPage(ctx context.Context, obj *ProductSkuResultType) (count int, err error) {
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

func (r *GeneratedProductSkuResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductSkuResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductSkuResultTypeResolver) PerPage(ctx context.Context, obj *ProductSkuResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductSkuResolver struct{ *GeneratedResolver }

func (r *GeneratedProductSkuResolver) Product(ctx context.Context, obj *ProductSku) (res *Product, err error) {
	return r.Handlers.ProductSkuProduct(ctx, r.GeneratedResolver, obj)
}
func ProductSkuProductHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (items *Product, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Product"); err != nil {
		return items, errors.New("Product " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ProductID

	if objKey != "" {
		item, _ := loaders["Product"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Product)

		if items == nil {
			items = &Product{}
		}

	}

	return
}

func (r *GeneratedProductSkuResolver) SpecificationValues(ctx context.Context, obj *ProductSku) (res []*ProductSkuSpecificationValue, err error) {
	return r.Handlers.ProductSkuSpecificationValues(ctx, r.GeneratedResolver, obj)
}
func ProductSkuSpecificationValuesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (items []*ProductSkuSpecificationValue, err error) {

	items = []*ProductSkuSpecificationValue{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SpecificationValues"); err != nil {
		return items, errors.New("SpecificationValues " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductSkuSpecificationValueSku"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductSkuSpecificationValue{}
	if item != nil {
		items = item.([]*ProductSkuSpecificationValue)
	}

	return
}

func (r *GeneratedProductSkuResolver) SpecificationValuesIds(ctx context.Context, obj *ProductSku) (ids []string, err error) {

	items := []*ProductSkuSpecificationValue{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SkuAndProductSkuSpecificationValueIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductSkuSpecificationValue)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductSkuResolver) Packages(ctx context.Context, obj *ProductSku) (res []*ProductPackage, err error) {
	return r.Handlers.ProductSkuPackages(ctx, r.GeneratedResolver, obj)
}
func ProductSkuPackagesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (items []*ProductPackage, err error) {

	items = []*ProductPackage{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Packages"); err != nil {
		return items, errors.New("Packages " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductPackageSku"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductPackage{}
	if item != nil {
		items = item.([]*ProductPackage)
	}

	return
}

func (r *GeneratedProductSkuResolver) PackagesIds(ctx context.Context, obj *ProductSku) (ids []string, err error) {

	items := []*ProductPackage{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SkuAndProductPackageIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductPackage)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductSkuResolver) Listings(ctx context.Context, obj *ProductSku) (res []*StoreListing, err error) {
	return r.Handlers.ProductSkuListings(ctx, r.GeneratedResolver, obj)
}
func ProductSkuListingsHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSku) (items []*StoreListing, err error) {

	items = []*StoreListing{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Listings"); err != nil {
		return items, errors.New("Listings " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreListingSku"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreListing{}
	if item != nil {
		items = item.([]*StoreListing)
	}

	return
}

func (r *GeneratedProductSkuResolver) ListingsIds(ctx context.Context, obj *ProductSku) (ids []string, err error) {

	items := []*StoreListing{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SkuAndStoreListingIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreListing)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryProductPackageHandlerOptions struct {
	ID     *string
	Filter *ProductPackageFilterType
}

func (r *GeneratedQueryResolver) ProductPackage(ctx context.Context, id *string, filter *ProductPackageFilterType) (*ProductPackage, error) {
	opts := QueryProductPackageHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProductPackage(ctx, r.GeneratedResolver, opts)
}
func QueryProductPackageHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductPackageHandlerOptions) (*ProductPackage, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductPackageQueryFilter{}
	rt := &ProductPackageResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("product_packages", ctx)+".id = ?", *opts.ID)
	}

	var items []*ProductPackage
	giOpts := GetItemsOptions{
		Alias:      TableName("product_packages", ctx),
		Preloaders: []string{},
		Item:       &ProductPackage{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "ProductPackage"}
	}
	return items[0], err
}

type QueryProductPackagesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductPackageSortType
	Filter      *ProductPackageFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) ProductPackages(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductPackageSortType, filter *ProductPackageFilterType, rand *bool) (*ProductPackageResultType, error) {
	opts := QueryProductPackagesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProductPackages(ctx, r.GeneratedResolver, opts)
}
func QueryProductPackagesHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductPackagesHandlerOptions) (*ProductPackageResultType, error) {
	query := ProductPackageQueryFilter{opts.Q}

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

	return &ProductPackageResultType{
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

type GeneratedProductPackageResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductPackageResultTypeResolver) Data(ctx context.Context, obj *ProductPackageResultType) (items []*ProductPackage, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("product_packages", ctx),
		Preloaders: []string{},
		Item:       &ProductPackage{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*ProductPackage{}
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

func (r *GeneratedProductPackageResultTypeResolver) Total(ctx context.Context, obj *ProductPackageResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("product_packages", ctx), &ProductPackage{})
}

func (r *GeneratedProductPackageResultTypeResolver) TotalPage(ctx context.Context, obj *ProductPackageResultType) (count int, err error) {
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

func (r *GeneratedProductPackageResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductPackageResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductPackageResultTypeResolver) PerPage(ctx context.Context, obj *ProductPackageResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductPackageResolver struct{ *GeneratedResolver }

func (r *GeneratedProductPackageResolver) Sku(ctx context.Context, obj *ProductPackage) (res *ProductSku, err error) {
	return r.Handlers.ProductPackageSku(ctx, r.GeneratedResolver, obj)
}
func ProductPackageSkuHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items *ProductSku, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductSku"); err != nil {
		return items, errors.New("ProductSku " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.SkuID

	if objKey != "" {
		item, _ := loaders["ProductSku"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*ProductSku)

		if items == nil {
			items = &ProductSku{}
		}

	}

	return
}

func (r *GeneratedProductPackageResolver) Template(ctx context.Context, obj *ProductPackage) (res *ProductPackageTemplate, err error) {
	return r.Handlers.ProductPackageTemplate(ctx, r.GeneratedResolver, obj)
}
func ProductPackageTemplateHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items *ProductPackageTemplate, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackageTemplate"); err != nil {
		return items, errors.New("ProductPackageTemplate " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.TemplateID

	if objKey != nil {
		item, _ := loaders["ProductPackageTemplate"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductPackageTemplate)

	}

	return
}

func (r *GeneratedProductPackageResolver) ContainsPackage(ctx context.Context, obj *ProductPackage) (res *ProductPackage, err error) {
	return r.Handlers.ProductPackageContainsPackage(ctx, r.GeneratedResolver, obj)
}
func ProductPackageContainsPackageHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items *ProductPackage, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackage"); err != nil {
		return items, errors.New("ProductPackage " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ContainsPackageID

	if objKey != nil {
		item, _ := loaders["ProductPackage"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductPackage)

	}

	return
}

func (r *GeneratedProductPackageResolver) ContainedByPackages(ctx context.Context, obj *ProductPackage) (res []*ProductPackage, err error) {
	return r.Handlers.ProductPackageContainedByPackages(ctx, r.GeneratedResolver, obj)
}
func ProductPackageContainedByPackagesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items []*ProductPackage, err error) {

	items = []*ProductPackage{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ContainedByPackages"); err != nil {
		return items, errors.New("ContainedByPackages " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductPackageContainsPackage"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductPackage{}
	if item != nil {
		items = item.([]*ProductPackage)
	}

	return
}

func (r *GeneratedProductPackageResolver) ContainedByPackagesIds(ctx context.Context, obj *ProductPackage) (ids []string, err error) {

	items := []*ProductPackage{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ContainsPackageAndProductPackageIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductPackage)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductPackageResolver) Offers(ctx context.Context, obj *ProductPackage) (res []*StorePackageOffer, err error) {
	return r.Handlers.ProductPackageOffers(ctx, r.GeneratedResolver, obj)
}
func ProductPackageOffersHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items []*StorePackageOffer, err error) {

	items = []*StorePackageOffer{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Offers"); err != nil {
		return items, errors.New("Offers " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StorePackageOfferPackage"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StorePackageOffer{}
	if item != nil {
		items = item.([]*StorePackageOffer)
	}

	return
}

func (r *GeneratedProductPackageResolver) OffersIds(ctx context.Context, obj *ProductPackage) (ids []string, err error) {

	items := []*StorePackageOffer{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["PackageAndStorePackageOfferIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StorePackageOffer)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductPackageResolver) Balances(ctx context.Context, obj *ProductPackage) (res []*StoreStockBalance, err error) {
	return r.Handlers.ProductPackageBalances(ctx, r.GeneratedResolver, obj)
}
func ProductPackageBalancesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items []*StoreStockBalance, err error) {

	items = []*StoreStockBalance{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Balances"); err != nil {
		return items, errors.New("Balances " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStockBalancePackage"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStockBalance{}
	if item != nil {
		items = item.([]*StoreStockBalance)
	}

	return
}

func (r *GeneratedProductPackageResolver) BalancesIds(ctx context.Context, obj *ProductPackage) (ids []string, err error) {

	items := []*StoreStockBalance{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["PackageAndStoreStockBalanceIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStockBalance)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductPackageResolver) MovementSources(ctx context.Context, obj *ProductPackage) (res []*StoreStockMovement, err error) {
	return r.Handlers.ProductPackageMovementSources(ctx, r.GeneratedResolver, obj)
}
func ProductPackageMovementSourcesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items []*StoreStockMovement, err error) {

	items = []*StoreStockMovement{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "MovementSources"); err != nil {
		return items, errors.New("MovementSources " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStockMovementSourcePackage"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStockMovement{}
	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	return
}

func (r *GeneratedProductPackageResolver) MovementSourcesIds(ctx context.Context, obj *ProductPackage) (ids []string, err error) {

	items := []*StoreStockMovement{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SourcePackageAndStoreStockMovementIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductPackageResolver) MovementTargets(ctx context.Context, obj *ProductPackage) (res []*StoreStockMovement, err error) {
	return r.Handlers.ProductPackageMovementTargets(ctx, r.GeneratedResolver, obj)
}
func ProductPackageMovementTargetsHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items []*StoreStockMovement, err error) {

	items = []*StoreStockMovement{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "MovementTargets"); err != nil {
		return items, errors.New("MovementTargets " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStockMovementTargetPackage"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStockMovement{}
	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	return
}

func (r *GeneratedProductPackageResolver) MovementTargetsIds(ctx context.Context, obj *ProductPackage) (ids []string, err error) {

	items := []*StoreStockMovement{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["TargetPackageAndStoreStockMovementIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductPackageResolver) StocktakeLines(ctx context.Context, obj *ProductPackage) (res []*StoreStocktakeLine, err error) {
	return r.Handlers.ProductPackageStocktakeLines(ctx, r.GeneratedResolver, obj)
}
func ProductPackageStocktakeLinesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackage) (items []*StoreStocktakeLine, err error) {

	items = []*StoreStocktakeLine{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StocktakeLines"); err != nil {
		return items, errors.New("StocktakeLines " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStocktakeLinePackage"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStocktakeLine{}
	if item != nil {
		items = item.([]*StoreStocktakeLine)
	}

	return
}

func (r *GeneratedProductPackageResolver) StocktakeLinesIds(ctx context.Context, obj *ProductPackage) (ids []string, err error) {

	items := []*StoreStocktakeLine{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["PackageAndStoreStocktakeLineIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStocktakeLine)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QuerySpecificationDefinitionHandlerOptions struct {
	ID     *string
	Filter *SpecificationDefinitionFilterType
}

func (r *GeneratedQueryResolver) SpecificationDefinition(ctx context.Context, id *string, filter *SpecificationDefinitionFilterType) (*SpecificationDefinition, error) {
	opts := QuerySpecificationDefinitionHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QuerySpecificationDefinition(ctx, r.GeneratedResolver, opts)
}
func QuerySpecificationDefinitionHandler(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationDefinitionHandlerOptions) (*SpecificationDefinition, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := SpecificationDefinitionQueryFilter{}
	rt := &SpecificationDefinitionResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("specification_definitions", ctx)+".id = ?", *opts.ID)
	}

	var items []*SpecificationDefinition
	giOpts := GetItemsOptions{
		Alias:      TableName("specification_definitions", ctx),
		Preloaders: []string{},
		Item:       &SpecificationDefinition{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "SpecificationDefinition"}
	}
	return items[0], err
}

type QuerySpecificationDefinitionsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*SpecificationDefinitionSortType
	Filter      *SpecificationDefinitionFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) SpecificationDefinitions(ctx context.Context, current_page *int, per_page *int, q *string, sort []*SpecificationDefinitionSortType, filter *SpecificationDefinitionFilterType, rand *bool) (*SpecificationDefinitionResultType, error) {
	opts := QuerySpecificationDefinitionsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QuerySpecificationDefinitions(ctx, r.GeneratedResolver, opts)
}
func QuerySpecificationDefinitionsHandler(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationDefinitionsHandlerOptions) (*SpecificationDefinitionResultType, error) {
	query := SpecificationDefinitionQueryFilter{opts.Q}

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

	return &SpecificationDefinitionResultType{
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

type GeneratedSpecificationDefinitionResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedSpecificationDefinitionResultTypeResolver) Data(ctx context.Context, obj *SpecificationDefinitionResultType) (items []*SpecificationDefinition, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("specification_definitions", ctx),
		Preloaders: []string{},
		Item:       &SpecificationDefinition{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*SpecificationDefinition{}
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

func (r *GeneratedSpecificationDefinitionResultTypeResolver) Total(ctx context.Context, obj *SpecificationDefinitionResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("specification_definitions", ctx), &SpecificationDefinition{})
}

func (r *GeneratedSpecificationDefinitionResultTypeResolver) TotalPage(ctx context.Context, obj *SpecificationDefinitionResultType) (count int, err error) {
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

func (r *GeneratedSpecificationDefinitionResultTypeResolver) CurrentPage(ctx context.Context, obj *SpecificationDefinitionResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedSpecificationDefinitionResultTypeResolver) PerPage(ctx context.Context, obj *SpecificationDefinitionResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedSpecificationDefinitionResolver struct{ *GeneratedResolver }

func (r *GeneratedSpecificationDefinitionResolver) Organization(ctx context.Context, obj *SpecificationDefinition) (res *Organization, err error) {
	return r.Handlers.SpecificationDefinitionOrganization(ctx, r.GeneratedResolver, obj)
}
func SpecificationDefinitionOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *SpecificationDefinition) (items *Organization, err error) {

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

func (r *GeneratedSpecificationDefinitionResolver) Values(ctx context.Context, obj *SpecificationDefinition) (res []*SpecificationValue, err error) {
	return r.Handlers.SpecificationDefinitionValues(ctx, r.GeneratedResolver, obj)
}
func SpecificationDefinitionValuesHandler(ctx context.Context, r *GeneratedResolver, obj *SpecificationDefinition) (items []*SpecificationValue, err error) {

	items = []*SpecificationValue{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Values"); err != nil {
		return items, errors.New("Values " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SpecificationValueSpecification"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*SpecificationValue{}
	if item != nil {
		items = item.([]*SpecificationValue)
	}

	return
}

func (r *GeneratedSpecificationDefinitionResolver) ValuesIds(ctx context.Context, obj *SpecificationDefinition) (ids []string, err error) {

	items := []*SpecificationValue{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["SpecificationAndSpecificationValueIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*SpecificationValue)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QuerySpecificationValueHandlerOptions struct {
	ID     *string
	Filter *SpecificationValueFilterType
}

func (r *GeneratedQueryResolver) SpecificationValue(ctx context.Context, id *string, filter *SpecificationValueFilterType) (*SpecificationValue, error) {
	opts := QuerySpecificationValueHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QuerySpecificationValue(ctx, r.GeneratedResolver, opts)
}
func QuerySpecificationValueHandler(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationValueHandlerOptions) (*SpecificationValue, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := SpecificationValueQueryFilter{}
	rt := &SpecificationValueResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("specification_values", ctx)+".id = ?", *opts.ID)
	}

	var items []*SpecificationValue
	giOpts := GetItemsOptions{
		Alias:      TableName("specification_values", ctx),
		Preloaders: []string{},
		Item:       &SpecificationValue{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "SpecificationValue"}
	}
	return items[0], err
}

type QuerySpecificationValuesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*SpecificationValueSortType
	Filter      *SpecificationValueFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) SpecificationValues(ctx context.Context, current_page *int, per_page *int, q *string, sort []*SpecificationValueSortType, filter *SpecificationValueFilterType, rand *bool) (*SpecificationValueResultType, error) {
	opts := QuerySpecificationValuesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QuerySpecificationValues(ctx, r.GeneratedResolver, opts)
}
func QuerySpecificationValuesHandler(ctx context.Context, r *GeneratedResolver, opts QuerySpecificationValuesHandlerOptions) (*SpecificationValueResultType, error) {
	query := SpecificationValueQueryFilter{opts.Q}

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

	return &SpecificationValueResultType{
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

type GeneratedSpecificationValueResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedSpecificationValueResultTypeResolver) Data(ctx context.Context, obj *SpecificationValueResultType) (items []*SpecificationValue, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("specification_values", ctx),
		Preloaders: []string{},
		Item:       &SpecificationValue{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*SpecificationValue{}
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

func (r *GeneratedSpecificationValueResultTypeResolver) Total(ctx context.Context, obj *SpecificationValueResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("specification_values", ctx), &SpecificationValue{})
}

func (r *GeneratedSpecificationValueResultTypeResolver) TotalPage(ctx context.Context, obj *SpecificationValueResultType) (count int, err error) {
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

func (r *GeneratedSpecificationValueResultTypeResolver) CurrentPage(ctx context.Context, obj *SpecificationValueResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedSpecificationValueResultTypeResolver) PerPage(ctx context.Context, obj *SpecificationValueResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedSpecificationValueResolver struct{ *GeneratedResolver }

func (r *GeneratedSpecificationValueResolver) Specification(ctx context.Context, obj *SpecificationValue) (res *SpecificationDefinition, err error) {
	return r.Handlers.SpecificationValueSpecification(ctx, r.GeneratedResolver, obj)
}
func SpecificationValueSpecificationHandler(ctx context.Context, r *GeneratedResolver, obj *SpecificationValue) (items *SpecificationDefinition, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SpecificationDefinition"); err != nil {
		return items, errors.New("SpecificationDefinition " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.SpecificationID

	if objKey != "" {
		item, _ := loaders["SpecificationDefinition"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*SpecificationDefinition)

		if items == nil {
			items = &SpecificationDefinition{}
		}

	}

	return
}

func (r *GeneratedSpecificationValueResolver) ProductChoices(ctx context.Context, obj *SpecificationValue) (res []*ProductSpecificationChoice, err error) {
	return r.Handlers.SpecificationValueProductChoices(ctx, r.GeneratedResolver, obj)
}
func SpecificationValueProductChoicesHandler(ctx context.Context, r *GeneratedResolver, obj *SpecificationValue) (items []*ProductSpecificationChoice, err error) {

	items = []*ProductSpecificationChoice{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductChoices"); err != nil {
		return items, errors.New("ProductChoices " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductSpecificationChoiceValue"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductSpecificationChoice{}
	if item != nil {
		items = item.([]*ProductSpecificationChoice)
	}

	return
}

func (r *GeneratedSpecificationValueResolver) ProductChoicesIds(ctx context.Context, obj *SpecificationValue) (ids []string, err error) {

	items := []*ProductSpecificationChoice{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ValueAndProductSpecificationChoiceIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductSpecificationChoice)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedSpecificationValueResolver) SkuValues(ctx context.Context, obj *SpecificationValue) (res []*ProductSkuSpecificationValue, err error) {
	return r.Handlers.SpecificationValueSkuValues(ctx, r.GeneratedResolver, obj)
}
func SpecificationValueSkuValuesHandler(ctx context.Context, r *GeneratedResolver, obj *SpecificationValue) (items []*ProductSkuSpecificationValue, err error) {

	items = []*ProductSkuSpecificationValue{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SkuValues"); err != nil {
		return items, errors.New("SkuValues " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductSkuSpecificationValueValue"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductSkuSpecificationValue{}
	if item != nil {
		items = item.([]*ProductSkuSpecificationValue)
	}

	return
}

func (r *GeneratedSpecificationValueResolver) SkuValuesIds(ctx context.Context, obj *SpecificationValue) (ids []string, err error) {

	items := []*ProductSkuSpecificationValue{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ValueAndProductSkuSpecificationValueIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductSkuSpecificationValue)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryProductSpecificationChoiceHandlerOptions struct {
	ID     *string
	Filter *ProductSpecificationChoiceFilterType
}

func (r *GeneratedQueryResolver) ProductSpecificationChoice(ctx context.Context, id *string, filter *ProductSpecificationChoiceFilterType) (*ProductSpecificationChoice, error) {
	opts := QueryProductSpecificationChoiceHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProductSpecificationChoice(ctx, r.GeneratedResolver, opts)
}
func QueryProductSpecificationChoiceHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductSpecificationChoiceHandlerOptions) (*ProductSpecificationChoice, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductSpecificationChoiceQueryFilter{}
	rt := &ProductSpecificationChoiceResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("product_specification_choices", ctx)+".id = ?", *opts.ID)
	}

	var items []*ProductSpecificationChoice
	giOpts := GetItemsOptions{
		Alias:      TableName("product_specification_choices", ctx),
		Preloaders: []string{},
		Item:       &ProductSpecificationChoice{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "ProductSpecificationChoice"}
	}
	return items[0], err
}

type QueryProductSpecificationChoicesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductSpecificationChoiceSortType
	Filter      *ProductSpecificationChoiceFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) ProductSpecificationChoices(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductSpecificationChoiceSortType, filter *ProductSpecificationChoiceFilterType, rand *bool) (*ProductSpecificationChoiceResultType, error) {
	opts := QueryProductSpecificationChoicesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProductSpecificationChoices(ctx, r.GeneratedResolver, opts)
}
func QueryProductSpecificationChoicesHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductSpecificationChoicesHandlerOptions) (*ProductSpecificationChoiceResultType, error) {
	query := ProductSpecificationChoiceQueryFilter{opts.Q}

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

	return &ProductSpecificationChoiceResultType{
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

type GeneratedProductSpecificationChoiceResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductSpecificationChoiceResultTypeResolver) Data(ctx context.Context, obj *ProductSpecificationChoiceResultType) (items []*ProductSpecificationChoice, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("product_specification_choices", ctx),
		Preloaders: []string{},
		Item:       &ProductSpecificationChoice{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*ProductSpecificationChoice{}
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

func (r *GeneratedProductSpecificationChoiceResultTypeResolver) Total(ctx context.Context, obj *ProductSpecificationChoiceResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("product_specification_choices", ctx), &ProductSpecificationChoice{})
}

func (r *GeneratedProductSpecificationChoiceResultTypeResolver) TotalPage(ctx context.Context, obj *ProductSpecificationChoiceResultType) (count int, err error) {
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

func (r *GeneratedProductSpecificationChoiceResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductSpecificationChoiceResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductSpecificationChoiceResultTypeResolver) PerPage(ctx context.Context, obj *ProductSpecificationChoiceResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductSpecificationChoiceResolver struct{ *GeneratedResolver }

func (r *GeneratedProductSpecificationChoiceResolver) Product(ctx context.Context, obj *ProductSpecificationChoice) (res *Product, err error) {
	return r.Handlers.ProductSpecificationChoiceProduct(ctx, r.GeneratedResolver, obj)
}
func ProductSpecificationChoiceProductHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSpecificationChoice) (items *Product, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Product"); err != nil {
		return items, errors.New("Product " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ProductID

	if objKey != "" {
		item, _ := loaders["Product"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Product)

		if items == nil {
			items = &Product{}
		}

	}

	return
}

func (r *GeneratedProductSpecificationChoiceResolver) Value(ctx context.Context, obj *ProductSpecificationChoice) (res *SpecificationValue, err error) {
	return r.Handlers.ProductSpecificationChoiceValue(ctx, r.GeneratedResolver, obj)
}
func ProductSpecificationChoiceValueHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSpecificationChoice) (items *SpecificationValue, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SpecificationValue"); err != nil {
		return items, errors.New("SpecificationValue " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ValueID

	if objKey != "" {
		item, _ := loaders["SpecificationValue"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*SpecificationValue)

		if items == nil {
			items = &SpecificationValue{}
		}

	}

	return
}

type QueryProductSkuSpecificationValueHandlerOptions struct {
	ID     *string
	Filter *ProductSkuSpecificationValueFilterType
}

func (r *GeneratedQueryResolver) ProductSkuSpecificationValue(ctx context.Context, id *string, filter *ProductSkuSpecificationValueFilterType) (*ProductSkuSpecificationValue, error) {
	opts := QueryProductSkuSpecificationValueHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProductSkuSpecificationValue(ctx, r.GeneratedResolver, opts)
}
func QueryProductSkuSpecificationValueHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductSkuSpecificationValueHandlerOptions) (*ProductSkuSpecificationValue, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductSkuSpecificationValueQueryFilter{}
	rt := &ProductSkuSpecificationValueResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("product_sku_specification_values", ctx)+".id = ?", *opts.ID)
	}

	var items []*ProductSkuSpecificationValue
	giOpts := GetItemsOptions{
		Alias:      TableName("product_sku_specification_values", ctx),
		Preloaders: []string{},
		Item:       &ProductSkuSpecificationValue{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "ProductSkuSpecificationValue"}
	}
	return items[0], err
}

type QueryProductSkuSpecificationValuesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductSkuSpecificationValueSortType
	Filter      *ProductSkuSpecificationValueFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) ProductSkuSpecificationValues(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductSkuSpecificationValueSortType, filter *ProductSkuSpecificationValueFilterType, rand *bool) (*ProductSkuSpecificationValueResultType, error) {
	opts := QueryProductSkuSpecificationValuesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProductSkuSpecificationValues(ctx, r.GeneratedResolver, opts)
}
func QueryProductSkuSpecificationValuesHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductSkuSpecificationValuesHandlerOptions) (*ProductSkuSpecificationValueResultType, error) {
	query := ProductSkuSpecificationValueQueryFilter{opts.Q}

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

	return &ProductSkuSpecificationValueResultType{
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

type GeneratedProductSkuSpecificationValueResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductSkuSpecificationValueResultTypeResolver) Data(ctx context.Context, obj *ProductSkuSpecificationValueResultType) (items []*ProductSkuSpecificationValue, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("product_sku_specification_values", ctx),
		Preloaders: []string{},
		Item:       &ProductSkuSpecificationValue{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*ProductSkuSpecificationValue{}
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

func (r *GeneratedProductSkuSpecificationValueResultTypeResolver) Total(ctx context.Context, obj *ProductSkuSpecificationValueResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("product_sku_specification_values", ctx), &ProductSkuSpecificationValue{})
}

func (r *GeneratedProductSkuSpecificationValueResultTypeResolver) TotalPage(ctx context.Context, obj *ProductSkuSpecificationValueResultType) (count int, err error) {
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

func (r *GeneratedProductSkuSpecificationValueResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductSkuSpecificationValueResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductSkuSpecificationValueResultTypeResolver) PerPage(ctx context.Context, obj *ProductSkuSpecificationValueResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductSkuSpecificationValueResolver struct{ *GeneratedResolver }

func (r *GeneratedProductSkuSpecificationValueResolver) Sku(ctx context.Context, obj *ProductSkuSpecificationValue) (res *ProductSku, err error) {
	return r.Handlers.ProductSkuSpecificationValueSku(ctx, r.GeneratedResolver, obj)
}
func ProductSkuSpecificationValueSkuHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSkuSpecificationValue) (items *ProductSku, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductSku"); err != nil {
		return items, errors.New("ProductSku " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.SkuID

	if objKey != "" {
		item, _ := loaders["ProductSku"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*ProductSku)

		if items == nil {
			items = &ProductSku{}
		}

	}

	return
}

func (r *GeneratedProductSkuSpecificationValueResolver) Value(ctx context.Context, obj *ProductSkuSpecificationValue) (res *SpecificationValue, err error) {
	return r.Handlers.ProductSkuSpecificationValueValue(ctx, r.GeneratedResolver, obj)
}
func ProductSkuSpecificationValueValueHandler(ctx context.Context, r *GeneratedResolver, obj *ProductSkuSpecificationValue) (items *SpecificationValue, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "SpecificationValue"); err != nil {
		return items, errors.New("SpecificationValue " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ValueID

	if objKey != "" {
		item, _ := loaders["SpecificationValue"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*SpecificationValue)

		if items == nil {
			items = &SpecificationValue{}
		}

	}

	return
}

type QueryProductPackageTemplateHandlerOptions struct {
	ID     *string
	Filter *ProductPackageTemplateFilterType
}

func (r *GeneratedQueryResolver) ProductPackageTemplate(ctx context.Context, id *string, filter *ProductPackageTemplateFilterType) (*ProductPackageTemplate, error) {
	opts := QueryProductPackageTemplateHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryProductPackageTemplate(ctx, r.GeneratedResolver, opts)
}
func QueryProductPackageTemplateHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductPackageTemplateHandlerOptions) (*ProductPackageTemplate, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := ProductPackageTemplateQueryFilter{}
	rt := &ProductPackageTemplateResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("product_package_templates", ctx)+".id = ?", *opts.ID)
	}

	var items []*ProductPackageTemplate
	giOpts := GetItemsOptions{
		Alias:      TableName("product_package_templates", ctx),
		Preloaders: []string{},
		Item:       &ProductPackageTemplate{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "ProductPackageTemplate"}
	}
	return items[0], err
}

type QueryProductPackageTemplatesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*ProductPackageTemplateSortType
	Filter      *ProductPackageTemplateFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) ProductPackageTemplates(ctx context.Context, current_page *int, per_page *int, q *string, sort []*ProductPackageTemplateSortType, filter *ProductPackageTemplateFilterType, rand *bool) (*ProductPackageTemplateResultType, error) {
	opts := QueryProductPackageTemplatesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryProductPackageTemplates(ctx, r.GeneratedResolver, opts)
}
func QueryProductPackageTemplatesHandler(ctx context.Context, r *GeneratedResolver, opts QueryProductPackageTemplatesHandlerOptions) (*ProductPackageTemplateResultType, error) {
	query := ProductPackageTemplateQueryFilter{opts.Q}

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

	return &ProductPackageTemplateResultType{
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

type GeneratedProductPackageTemplateResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedProductPackageTemplateResultTypeResolver) Data(ctx context.Context, obj *ProductPackageTemplateResultType) (items []*ProductPackageTemplate, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("product_package_templates", ctx),
		Preloaders: []string{},
		Item:       &ProductPackageTemplate{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*ProductPackageTemplate{}
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

func (r *GeneratedProductPackageTemplateResultTypeResolver) Total(ctx context.Context, obj *ProductPackageTemplateResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("product_package_templates", ctx), &ProductPackageTemplate{})
}

func (r *GeneratedProductPackageTemplateResultTypeResolver) TotalPage(ctx context.Context, obj *ProductPackageTemplateResultType) (count int, err error) {
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

func (r *GeneratedProductPackageTemplateResultTypeResolver) CurrentPage(ctx context.Context, obj *ProductPackageTemplateResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedProductPackageTemplateResultTypeResolver) PerPage(ctx context.Context, obj *ProductPackageTemplateResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedProductPackageTemplateResolver struct{ *GeneratedResolver }

func (r *GeneratedProductPackageTemplateResolver) Organization(ctx context.Context, obj *ProductPackageTemplate) (res *Organization, err error) {
	return r.Handlers.ProductPackageTemplateOrganization(ctx, r.GeneratedResolver, obj)
}
func ProductPackageTemplateOrganizationHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (items *Organization, err error) {

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

func (r *GeneratedProductPackageTemplateResolver) DefaultProducts(ctx context.Context, obj *ProductPackageTemplate) (res []*Product, err error) {
	return r.Handlers.ProductPackageTemplateDefaultProducts(ctx, r.GeneratedResolver, obj)
}
func ProductPackageTemplateDefaultProductsHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (items []*Product, err error) {

	items = []*Product{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "DefaultProducts"); err != nil {
		return items, errors.New("DefaultProducts " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductDefaultPackageTemplate"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*Product{}
	if item != nil {
		items = item.([]*Product)
	}

	return
}

func (r *GeneratedProductPackageTemplateResolver) DefaultProductsIds(ctx context.Context, obj *ProductPackageTemplate) (ids []string, err error) {

	items := []*Product{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["DefaultPackageTemplateAndProductIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*Product)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductPackageTemplateResolver) ContainsPackage(ctx context.Context, obj *ProductPackageTemplate) (res *ProductPackageTemplate, err error) {
	return r.Handlers.ProductPackageTemplateContainsPackage(ctx, r.GeneratedResolver, obj)
}
func ProductPackageTemplateContainsPackageHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (items *ProductPackageTemplate, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackageTemplate"); err != nil {
		return items, errors.New("ProductPackageTemplate " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ContainsPackageID

	if objKey != nil {
		item, _ := loaders["ProductPackageTemplate"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductPackageTemplate)

	}

	return
}

func (r *GeneratedProductPackageTemplateResolver) ContainedByPackages(ctx context.Context, obj *ProductPackageTemplate) (res []*ProductPackageTemplate, err error) {
	return r.Handlers.ProductPackageTemplateContainedByPackages(ctx, r.GeneratedResolver, obj)
}
func ProductPackageTemplateContainedByPackagesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (items []*ProductPackageTemplate, err error) {

	items = []*ProductPackageTemplate{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ContainedByPackages"); err != nil {
		return items, errors.New("ContainedByPackages " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductPackageTemplateContainsPackage"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductPackageTemplate{}
	if item != nil {
		items = item.([]*ProductPackageTemplate)
	}

	return
}

func (r *GeneratedProductPackageTemplateResolver) ContainedByPackagesIds(ctx context.Context, obj *ProductPackageTemplate) (ids []string, err error) {

	items := []*ProductPackageTemplate{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ContainsPackageAndProductPackageTemplateIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductPackageTemplate)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedProductPackageTemplateResolver) CreatedPackages(ctx context.Context, obj *ProductPackageTemplate) (res []*ProductPackage, err error) {
	return r.Handlers.ProductPackageTemplateCreatedPackages(ctx, r.GeneratedResolver, obj)
}
func ProductPackageTemplateCreatedPackagesHandler(ctx context.Context, r *GeneratedResolver, obj *ProductPackageTemplate) (items []*ProductPackage, err error) {

	items = []*ProductPackage{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CreatedPackages"); err != nil {
		return items, errors.New("CreatedPackages " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ProductPackageTemplate"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*ProductPackage{}
	if item != nil {
		items = item.([]*ProductPackage)
	}

	return
}

func (r *GeneratedProductPackageTemplateResolver) CreatedPackagesIds(ctx context.Context, obj *ProductPackageTemplate) (ids []string, err error) {

	items := []*ProductPackage{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["TemplateAndProductPackageIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*ProductPackage)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryStoreListingHandlerOptions struct {
	ID     *string
	Filter *StoreListingFilterType
}

func (r *GeneratedQueryResolver) StoreListing(ctx context.Context, id *string, filter *StoreListingFilterType) (*StoreListing, error) {
	opts := QueryStoreListingHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStoreListing(ctx, r.GeneratedResolver, opts)
}
func QueryStoreListingHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreListingHandlerOptions) (*StoreListing, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StoreListingQueryFilter{}
	rt := &StoreListingResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_listings", ctx)+".id = ?", *opts.ID)
	}

	var items []*StoreListing
	giOpts := GetItemsOptions{
		Alias:      TableName("store_listings", ctx),
		Preloaders: []string{},
		Item:       &StoreListing{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StoreListing"}
	}
	return items[0], err
}

type QueryStoreListingsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StoreListingSortType
	Filter      *StoreListingFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StoreListings(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StoreListingSortType, filter *StoreListingFilterType, rand *bool) (*StoreListingResultType, error) {
	opts := QueryStoreListingsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStoreListings(ctx, r.GeneratedResolver, opts)
}
func QueryStoreListingsHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreListingsHandlerOptions) (*StoreListingResultType, error) {
	query := StoreListingQueryFilter{opts.Q}

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

	return &StoreListingResultType{
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

type GeneratedStoreListingResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreListingResultTypeResolver) Data(ctx context.Context, obj *StoreListingResultType) (items []*StoreListing, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_listings", ctx),
		Preloaders: []string{},
		Item:       &StoreListing{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StoreListing{}
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

func (r *GeneratedStoreListingResultTypeResolver) Total(ctx context.Context, obj *StoreListingResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_listings", ctx), &StoreListing{})
}

func (r *GeneratedStoreListingResultTypeResolver) TotalPage(ctx context.Context, obj *StoreListingResultType) (count int, err error) {
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

func (r *GeneratedStoreListingResultTypeResolver) CurrentPage(ctx context.Context, obj *StoreListingResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStoreListingResultTypeResolver) PerPage(ctx context.Context, obj *StoreListingResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStoreListingResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreListingResolver) Store(ctx context.Context, obj *StoreListing) (res *Store, err error) {
	return r.Handlers.StoreListingStore(ctx, r.GeneratedResolver, obj)
}
func StoreListingStoreHandler(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (items *Store, err error) {

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

func (r *GeneratedStoreListingResolver) Sku(ctx context.Context, obj *StoreListing) (res *ProductSku, err error) {
	return r.Handlers.StoreListingSku(ctx, r.GeneratedResolver, obj)
}
func StoreListingSkuHandler(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (items *ProductSku, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductSku"); err != nil {
		return items, errors.New("ProductSku " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.SkuID

	if objKey != "" {
		item, _ := loaders["ProductSku"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*ProductSku)

		if items == nil {
			items = &ProductSku{}
		}

	}

	return
}

func (r *GeneratedStoreListingResolver) Offers(ctx context.Context, obj *StoreListing) (res []*StorePackageOffer, err error) {
	return r.Handlers.StoreListingOffers(ctx, r.GeneratedResolver, obj)
}
func StoreListingOffersHandler(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (items []*StorePackageOffer, err error) {

	items = []*StorePackageOffer{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Offers"); err != nil {
		return items, errors.New("Offers " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StorePackageOfferListing"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StorePackageOffer{}
	if item != nil {
		items = item.([]*StorePackageOffer)
	}

	return
}

func (r *GeneratedStoreListingResolver) OffersIds(ctx context.Context, obj *StoreListing) (ids []string, err error) {

	items := []*StorePackageOffer{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ListingAndStorePackageOfferIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StorePackageOffer)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreListingResolver) Batches(ctx context.Context, obj *StoreListing) (res []*StoreInventoryBatch, err error) {
	return r.Handlers.StoreListingBatches(ctx, r.GeneratedResolver, obj)
}
func StoreListingBatchesHandler(ctx context.Context, r *GeneratedResolver, obj *StoreListing) (items []*StoreInventoryBatch, err error) {

	items = []*StoreInventoryBatch{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Batches"); err != nil {
		return items, errors.New("Batches " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreInventoryBatchListing"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreInventoryBatch{}
	if item != nil {
		items = item.([]*StoreInventoryBatch)
	}

	return
}

func (r *GeneratedStoreListingResolver) BatchesIds(ctx context.Context, obj *StoreListing) (ids []string, err error) {

	items := []*StoreInventoryBatch{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["ListingAndStoreInventoryBatchIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreInventoryBatch)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryStorePackageOfferHandlerOptions struct {
	ID     *string
	Filter *StorePackageOfferFilterType
}

func (r *GeneratedQueryResolver) StorePackageOffer(ctx context.Context, id *string, filter *StorePackageOfferFilterType) (*StorePackageOffer, error) {
	opts := QueryStorePackageOfferHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStorePackageOffer(ctx, r.GeneratedResolver, opts)
}
func QueryStorePackageOfferHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePackageOfferHandlerOptions) (*StorePackageOffer, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StorePackageOfferQueryFilter{}
	rt := &StorePackageOfferResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_package_offers", ctx)+".id = ?", *opts.ID)
	}

	var items []*StorePackageOffer
	giOpts := GetItemsOptions{
		Alias:      TableName("store_package_offers", ctx),
		Preloaders: []string{},
		Item:       &StorePackageOffer{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StorePackageOffer"}
	}
	return items[0], err
}

type QueryStorePackageOffersHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StorePackageOfferSortType
	Filter      *StorePackageOfferFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StorePackageOffers(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StorePackageOfferSortType, filter *StorePackageOfferFilterType, rand *bool) (*StorePackageOfferResultType, error) {
	opts := QueryStorePackageOffersHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStorePackageOffers(ctx, r.GeneratedResolver, opts)
}
func QueryStorePackageOffersHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePackageOffersHandlerOptions) (*StorePackageOfferResultType, error) {
	query := StorePackageOfferQueryFilter{opts.Q}

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

	return &StorePackageOfferResultType{
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

type GeneratedStorePackageOfferResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePackageOfferResultTypeResolver) Data(ctx context.Context, obj *StorePackageOfferResultType) (items []*StorePackageOffer, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_package_offers", ctx),
		Preloaders: []string{},
		Item:       &StorePackageOffer{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StorePackageOffer{}
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

func (r *GeneratedStorePackageOfferResultTypeResolver) Total(ctx context.Context, obj *StorePackageOfferResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_package_offers", ctx), &StorePackageOffer{})
}

func (r *GeneratedStorePackageOfferResultTypeResolver) TotalPage(ctx context.Context, obj *StorePackageOfferResultType) (count int, err error) {
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

func (r *GeneratedStorePackageOfferResultTypeResolver) CurrentPage(ctx context.Context, obj *StorePackageOfferResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStorePackageOfferResultTypeResolver) PerPage(ctx context.Context, obj *StorePackageOfferResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStorePackageOfferResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePackageOfferResolver) Listing(ctx context.Context, obj *StorePackageOffer) (res *StoreListing, err error) {
	return r.Handlers.StorePackageOfferListing(ctx, r.GeneratedResolver, obj)
}
func StorePackageOfferListingHandler(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (items *StoreListing, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StoreListing"); err != nil {
		return items, errors.New("StoreListing " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ListingID

	if objKey != "" {
		item, _ := loaders["StoreListing"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StoreListing)

		if items == nil {
			items = &StoreListing{}
		}

	}

	return
}

func (r *GeneratedStorePackageOfferResolver) Package(ctx context.Context, obj *StorePackageOffer) (res *ProductPackage, err error) {
	return r.Handlers.StorePackageOfferPackage(ctx, r.GeneratedResolver, obj)
}
func StorePackageOfferPackageHandler(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (items *ProductPackage, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackage"); err != nil {
		return items, errors.New("ProductPackage " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.PackageID

	if objKey != "" {
		item, _ := loaders["ProductPackage"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*ProductPackage)

		if items == nil {
			items = &ProductPackage{}
		}

	}

	return
}

func (r *GeneratedStorePackageOfferResolver) PriceRevisions(ctx context.Context, obj *StorePackageOffer) (res []*StorePriceRevision, err error) {
	return r.Handlers.StorePackageOfferPriceRevisions(ctx, r.GeneratedResolver, obj)
}
func StorePackageOfferPriceRevisionsHandler(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (items []*StorePriceRevision, err error) {

	items = []*StorePriceRevision{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "PriceRevisions"); err != nil {
		return items, errors.New("PriceRevisions " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StorePriceRevisionOffer"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StorePriceRevision{}
	if item != nil {
		items = item.([]*StorePriceRevision)
	}

	return
}

func (r *GeneratedStorePackageOfferResolver) PriceRevisionsIds(ctx context.Context, obj *StorePackageOffer) (ids []string, err error) {

	items := []*StorePriceRevision{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OfferAndStorePriceRevisionIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StorePriceRevision)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStorePackageOfferResolver) PromotionTargets(ctx context.Context, obj *StorePackageOffer) (res []*StorePromotionTarget, err error) {
	return r.Handlers.StorePackageOfferPromotionTargets(ctx, r.GeneratedResolver, obj)
}
func StorePackageOfferPromotionTargetsHandler(ctx context.Context, r *GeneratedResolver, obj *StorePackageOffer) (items []*StorePromotionTarget, err error) {

	items = []*StorePromotionTarget{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "PromotionTargets"); err != nil {
		return items, errors.New("PromotionTargets " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StorePromotionTargetOffer"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StorePromotionTarget{}
	if item != nil {
		items = item.([]*StorePromotionTarget)
	}

	return
}

func (r *GeneratedStorePackageOfferResolver) PromotionTargetsIds(ctx context.Context, obj *StorePackageOffer) (ids []string, err error) {

	items := []*StorePromotionTarget{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["OfferAndStorePromotionTargetIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StorePromotionTarget)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryStorePriceRevisionHandlerOptions struct {
	ID     *string
	Filter *StorePriceRevisionFilterType
}

func (r *GeneratedQueryResolver) StorePriceRevision(ctx context.Context, id *string, filter *StorePriceRevisionFilterType) (*StorePriceRevision, error) {
	opts := QueryStorePriceRevisionHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStorePriceRevision(ctx, r.GeneratedResolver, opts)
}
func QueryStorePriceRevisionHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePriceRevisionHandlerOptions) (*StorePriceRevision, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StorePriceRevisionQueryFilter{}
	rt := &StorePriceRevisionResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_price_revisions", ctx)+".id = ?", *opts.ID)
	}

	var items []*StorePriceRevision
	giOpts := GetItemsOptions{
		Alias:      TableName("store_price_revisions", ctx),
		Preloaders: []string{},
		Item:       &StorePriceRevision{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StorePriceRevision"}
	}
	return items[0], err
}

type QueryStorePriceRevisionsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StorePriceRevisionSortType
	Filter      *StorePriceRevisionFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StorePriceRevisions(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StorePriceRevisionSortType, filter *StorePriceRevisionFilterType, rand *bool) (*StorePriceRevisionResultType, error) {
	opts := QueryStorePriceRevisionsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStorePriceRevisions(ctx, r.GeneratedResolver, opts)
}
func QueryStorePriceRevisionsHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePriceRevisionsHandlerOptions) (*StorePriceRevisionResultType, error) {
	query := StorePriceRevisionQueryFilter{opts.Q}

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

	return &StorePriceRevisionResultType{
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

type GeneratedStorePriceRevisionResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePriceRevisionResultTypeResolver) Data(ctx context.Context, obj *StorePriceRevisionResultType) (items []*StorePriceRevision, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_price_revisions", ctx),
		Preloaders: []string{},
		Item:       &StorePriceRevision{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StorePriceRevision{}
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

func (r *GeneratedStorePriceRevisionResultTypeResolver) Total(ctx context.Context, obj *StorePriceRevisionResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_price_revisions", ctx), &StorePriceRevision{})
}

func (r *GeneratedStorePriceRevisionResultTypeResolver) TotalPage(ctx context.Context, obj *StorePriceRevisionResultType) (count int, err error) {
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

func (r *GeneratedStorePriceRevisionResultTypeResolver) CurrentPage(ctx context.Context, obj *StorePriceRevisionResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStorePriceRevisionResultTypeResolver) PerPage(ctx context.Context, obj *StorePriceRevisionResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStorePriceRevisionResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePriceRevisionResolver) Offer(ctx context.Context, obj *StorePriceRevision) (res *StorePackageOffer, err error) {
	return r.Handlers.StorePriceRevisionOffer(ctx, r.GeneratedResolver, obj)
}
func StorePriceRevisionOfferHandler(ctx context.Context, r *GeneratedResolver, obj *StorePriceRevision) (items *StorePackageOffer, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StorePackageOffer"); err != nil {
		return items, errors.New("StorePackageOffer " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OfferID

	if objKey != "" {
		item, _ := loaders["StorePackageOffer"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StorePackageOffer)

		if items == nil {
			items = &StorePackageOffer{}
		}

	}

	return
}

type QueryStoreInventoryBatchHandlerOptions struct {
	ID     *string
	Filter *StoreInventoryBatchFilterType
}

func (r *GeneratedQueryResolver) StoreInventoryBatch(ctx context.Context, id *string, filter *StoreInventoryBatchFilterType) (*StoreInventoryBatch, error) {
	opts := QueryStoreInventoryBatchHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStoreInventoryBatch(ctx, r.GeneratedResolver, opts)
}
func QueryStoreInventoryBatchHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreInventoryBatchHandlerOptions) (*StoreInventoryBatch, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StoreInventoryBatchQueryFilter{}
	rt := &StoreInventoryBatchResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_inventory_batches", ctx)+".id = ?", *opts.ID)
	}

	var items []*StoreInventoryBatch
	giOpts := GetItemsOptions{
		Alias:      TableName("store_inventory_batches", ctx),
		Preloaders: []string{},
		Item:       &StoreInventoryBatch{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StoreInventoryBatch"}
	}
	return items[0], err
}

type QueryStoreInventoryBatchesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StoreInventoryBatchSortType
	Filter      *StoreInventoryBatchFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StoreInventoryBatches(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StoreInventoryBatchSortType, filter *StoreInventoryBatchFilterType, rand *bool) (*StoreInventoryBatchResultType, error) {
	opts := QueryStoreInventoryBatchesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStoreInventoryBatches(ctx, r.GeneratedResolver, opts)
}
func QueryStoreInventoryBatchesHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreInventoryBatchesHandlerOptions) (*StoreInventoryBatchResultType, error) {
	query := StoreInventoryBatchQueryFilter{opts.Q}

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

	return &StoreInventoryBatchResultType{
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

type GeneratedStoreInventoryBatchResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreInventoryBatchResultTypeResolver) Data(ctx context.Context, obj *StoreInventoryBatchResultType) (items []*StoreInventoryBatch, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_inventory_batches", ctx),
		Preloaders: []string{},
		Item:       &StoreInventoryBatch{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StoreInventoryBatch{}
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

func (r *GeneratedStoreInventoryBatchResultTypeResolver) Total(ctx context.Context, obj *StoreInventoryBatchResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_inventory_batches", ctx), &StoreInventoryBatch{})
}

func (r *GeneratedStoreInventoryBatchResultTypeResolver) TotalPage(ctx context.Context, obj *StoreInventoryBatchResultType) (count int, err error) {
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

func (r *GeneratedStoreInventoryBatchResultTypeResolver) CurrentPage(ctx context.Context, obj *StoreInventoryBatchResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStoreInventoryBatchResultTypeResolver) PerPage(ctx context.Context, obj *StoreInventoryBatchResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStoreInventoryBatchResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreInventoryBatchResolver) Listing(ctx context.Context, obj *StoreInventoryBatch) (res *StoreListing, err error) {
	return r.Handlers.StoreInventoryBatchListing(ctx, r.GeneratedResolver, obj)
}
func StoreInventoryBatchListingHandler(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (items *StoreListing, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StoreListing"); err != nil {
		return items, errors.New("StoreListing " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.ListingID

	if objKey != "" {
		item, _ := loaders["StoreListing"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StoreListing)

		if items == nil {
			items = &StoreListing{}
		}

	}

	return
}

func (r *GeneratedStoreInventoryBatchResolver) Balances(ctx context.Context, obj *StoreInventoryBatch) (res []*StoreStockBalance, err error) {
	return r.Handlers.StoreInventoryBatchBalances(ctx, r.GeneratedResolver, obj)
}
func StoreInventoryBatchBalancesHandler(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (items []*StoreStockBalance, err error) {

	items = []*StoreStockBalance{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Balances"); err != nil {
		return items, errors.New("Balances " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStockBalanceBatch"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStockBalance{}
	if item != nil {
		items = item.([]*StoreStockBalance)
	}

	return
}

func (r *GeneratedStoreInventoryBatchResolver) BalancesIds(ctx context.Context, obj *StoreInventoryBatch) (ids []string, err error) {

	items := []*StoreStockBalance{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["BatchAndStoreStockBalanceIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStockBalance)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreInventoryBatchResolver) Movements(ctx context.Context, obj *StoreInventoryBatch) (res []*StoreStockMovement, err error) {
	return r.Handlers.StoreInventoryBatchMovements(ctx, r.GeneratedResolver, obj)
}
func StoreInventoryBatchMovementsHandler(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (items []*StoreStockMovement, err error) {

	items = []*StoreStockMovement{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Movements"); err != nil {
		return items, errors.New("Movements " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStockMovementBatch"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStockMovement{}
	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	return
}

func (r *GeneratedStoreInventoryBatchResolver) MovementsIds(ctx context.Context, obj *StoreInventoryBatch) (ids []string, err error) {

	items := []*StoreStockMovement{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["BatchAndStoreStockMovementIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

func (r *GeneratedStoreInventoryBatchResolver) StocktakeLines(ctx context.Context, obj *StoreInventoryBatch) (res []*StoreStocktakeLine, err error) {
	return r.Handlers.StoreInventoryBatchStocktakeLines(ctx, r.GeneratedResolver, obj)
}
func StoreInventoryBatchStocktakeLinesHandler(ctx context.Context, r *GeneratedResolver, obj *StoreInventoryBatch) (items []*StoreStocktakeLine, err error) {

	items = []*StoreStocktakeLine{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StocktakeLines"); err != nil {
		return items, errors.New("StocktakeLines " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStocktakeLineBatch"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStocktakeLine{}
	if item != nil {
		items = item.([]*StoreStocktakeLine)
	}

	return
}

func (r *GeneratedStoreInventoryBatchResolver) StocktakeLinesIds(ctx context.Context, obj *StoreInventoryBatch) (ids []string, err error) {

	items := []*StoreStocktakeLine{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["BatchAndStoreStocktakeLineIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStocktakeLine)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryStoreStockBalanceHandlerOptions struct {
	ID     *string
	Filter *StoreStockBalanceFilterType
}

func (r *GeneratedQueryResolver) StoreStockBalance(ctx context.Context, id *string, filter *StoreStockBalanceFilterType) (*StoreStockBalance, error) {
	opts := QueryStoreStockBalanceHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStoreStockBalance(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStockBalanceHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockBalanceHandlerOptions) (*StoreStockBalance, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StoreStockBalanceQueryFilter{}
	rt := &StoreStockBalanceResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_stock_balances", ctx)+".id = ?", *opts.ID)
	}

	var items []*StoreStockBalance
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stock_balances", ctx),
		Preloaders: []string{},
		Item:       &StoreStockBalance{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StoreStockBalance"}
	}
	return items[0], err
}

type QueryStoreStockBalancesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StoreStockBalanceSortType
	Filter      *StoreStockBalanceFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StoreStockBalances(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StoreStockBalanceSortType, filter *StoreStockBalanceFilterType, rand *bool) (*StoreStockBalanceResultType, error) {
	opts := QueryStoreStockBalancesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStoreStockBalances(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStockBalancesHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockBalancesHandlerOptions) (*StoreStockBalanceResultType, error) {
	query := StoreStockBalanceQueryFilter{opts.Q}

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

	return &StoreStockBalanceResultType{
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

type GeneratedStoreStockBalanceResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStockBalanceResultTypeResolver) Data(ctx context.Context, obj *StoreStockBalanceResultType) (items []*StoreStockBalance, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stock_balances", ctx),
		Preloaders: []string{},
		Item:       &StoreStockBalance{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StoreStockBalance{}
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

func (r *GeneratedStoreStockBalanceResultTypeResolver) Total(ctx context.Context, obj *StoreStockBalanceResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_stock_balances", ctx), &StoreStockBalance{})
}

func (r *GeneratedStoreStockBalanceResultTypeResolver) TotalPage(ctx context.Context, obj *StoreStockBalanceResultType) (count int, err error) {
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

func (r *GeneratedStoreStockBalanceResultTypeResolver) CurrentPage(ctx context.Context, obj *StoreStockBalanceResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStoreStockBalanceResultTypeResolver) PerPage(ctx context.Context, obj *StoreStockBalanceResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStoreStockBalanceResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStockBalanceResolver) Batch(ctx context.Context, obj *StoreStockBalance) (res *StoreInventoryBatch, err error) {
	return r.Handlers.StoreStockBalanceBatch(ctx, r.GeneratedResolver, obj)
}
func StoreStockBalanceBatchHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStockBalance) (items *StoreInventoryBatch, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StoreInventoryBatch"); err != nil {
		return items, errors.New("StoreInventoryBatch " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.BatchID

	if objKey != "" {
		item, _ := loaders["StoreInventoryBatch"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StoreInventoryBatch)

		if items == nil {
			items = &StoreInventoryBatch{}
		}

	}

	return
}

func (r *GeneratedStoreStockBalanceResolver) Package(ctx context.Context, obj *StoreStockBalance) (res *ProductPackage, err error) {
	return r.Handlers.StoreStockBalancePackage(ctx, r.GeneratedResolver, obj)
}
func StoreStockBalancePackageHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStockBalance) (items *ProductPackage, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackage"); err != nil {
		return items, errors.New("ProductPackage " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.PackageID

	if objKey != "" {
		item, _ := loaders["ProductPackage"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*ProductPackage)

		if items == nil {
			items = &ProductPackage{}
		}

	}

	return
}

type QueryStoreStocktakeHandlerOptions struct {
	ID     *string
	Filter *StoreStocktakeFilterType
}

func (r *GeneratedQueryResolver) StoreStocktake(ctx context.Context, id *string, filter *StoreStocktakeFilterType) (*StoreStocktake, error) {
	opts := QueryStoreStocktakeHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStoreStocktake(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStocktakeHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakeHandlerOptions) (*StoreStocktake, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StoreStocktakeQueryFilter{}
	rt := &StoreStocktakeResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_stocktakes", ctx)+".id = ?", *opts.ID)
	}

	var items []*StoreStocktake
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stocktakes", ctx),
		Preloaders: []string{},
		Item:       &StoreStocktake{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StoreStocktake"}
	}
	return items[0], err
}

type QueryStoreStocktakesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StoreStocktakeSortType
	Filter      *StoreStocktakeFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StoreStocktakes(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StoreStocktakeSortType, filter *StoreStocktakeFilterType, rand *bool) (*StoreStocktakeResultType, error) {
	opts := QueryStoreStocktakesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStoreStocktakes(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStocktakesHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakesHandlerOptions) (*StoreStocktakeResultType, error) {
	query := StoreStocktakeQueryFilter{opts.Q}

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

	return &StoreStocktakeResultType{
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

type GeneratedStoreStocktakeResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStocktakeResultTypeResolver) Data(ctx context.Context, obj *StoreStocktakeResultType) (items []*StoreStocktake, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stocktakes", ctx),
		Preloaders: []string{},
		Item:       &StoreStocktake{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StoreStocktake{}
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

func (r *GeneratedStoreStocktakeResultTypeResolver) Total(ctx context.Context, obj *StoreStocktakeResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_stocktakes", ctx), &StoreStocktake{})
}

func (r *GeneratedStoreStocktakeResultTypeResolver) TotalPage(ctx context.Context, obj *StoreStocktakeResultType) (count int, err error) {
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

func (r *GeneratedStoreStocktakeResultTypeResolver) CurrentPage(ctx context.Context, obj *StoreStocktakeResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStoreStocktakeResultTypeResolver) PerPage(ctx context.Context, obj *StoreStocktakeResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStoreStocktakeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStocktakeResolver) Store(ctx context.Context, obj *StoreStocktake) (res *Store, err error) {
	return r.Handlers.StoreStocktakeStore(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakeStoreHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (items *Store, err error) {

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

func (r *GeneratedStoreStocktakeResolver) InitiatedByAccount(ctx context.Context, obj *StoreStocktake) (res *Account, err error) {
	return r.Handlers.StoreStocktakeInitiatedByAccount(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakeInitiatedByAccountHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.InitiatedByAccountID

	if objKey != "" {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*Account)

		if items == nil {
			items = &Account{}
		}

	}

	return
}

func (r *GeneratedStoreStocktakeResolver) PostedBy(ctx context.Context, obj *StoreStocktake) (res *Account, err error) {
	return r.Handlers.StoreStocktakePostedBy(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakePostedByHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (items *Account, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Account"); err != nil {
		return items, errors.New("Account " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.PostedByID

	if objKey != nil {
		item, _ := loaders["Account"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*Account)

	}

	return
}

func (r *GeneratedStoreStocktakeResolver) Lines(ctx context.Context, obj *StoreStocktake) (res []*StoreStocktakeLine, err error) {
	return r.Handlers.StoreStocktakeLines(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakeLinesHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktake) (items []*StoreStocktakeLine, err error) {

	items = []*StoreStocktakeLine{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Lines"); err != nil {
		return items, errors.New("Lines " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStocktakeLineStocktake"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStocktakeLine{}
	if item != nil {
		items = item.([]*StoreStocktakeLine)
	}

	return
}

func (r *GeneratedStoreStocktakeResolver) LinesIds(ctx context.Context, obj *StoreStocktake) (ids []string, err error) {

	items := []*StoreStocktakeLine{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StocktakeAndStoreStocktakeLineIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStocktakeLine)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryStoreStocktakeLineHandlerOptions struct {
	ID     *string
	Filter *StoreStocktakeLineFilterType
}

func (r *GeneratedQueryResolver) StoreStocktakeLine(ctx context.Context, id *string, filter *StoreStocktakeLineFilterType) (*StoreStocktakeLine, error) {
	opts := QueryStoreStocktakeLineHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStoreStocktakeLine(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStocktakeLineHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakeLineHandlerOptions) (*StoreStocktakeLine, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StoreStocktakeLineQueryFilter{}
	rt := &StoreStocktakeLineResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_stocktake_lines", ctx)+".id = ?", *opts.ID)
	}

	var items []*StoreStocktakeLine
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stocktake_lines", ctx),
		Preloaders: []string{},
		Item:       &StoreStocktakeLine{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StoreStocktakeLine"}
	}
	return items[0], err
}

type QueryStoreStocktakeLinesHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StoreStocktakeLineSortType
	Filter      *StoreStocktakeLineFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StoreStocktakeLines(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StoreStocktakeLineSortType, filter *StoreStocktakeLineFilterType, rand *bool) (*StoreStocktakeLineResultType, error) {
	opts := QueryStoreStocktakeLinesHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStoreStocktakeLines(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStocktakeLinesHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStocktakeLinesHandlerOptions) (*StoreStocktakeLineResultType, error) {
	query := StoreStocktakeLineQueryFilter{opts.Q}

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

	return &StoreStocktakeLineResultType{
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

type GeneratedStoreStocktakeLineResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStocktakeLineResultTypeResolver) Data(ctx context.Context, obj *StoreStocktakeLineResultType) (items []*StoreStocktakeLine, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stocktake_lines", ctx),
		Preloaders: []string{},
		Item:       &StoreStocktakeLine{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StoreStocktakeLine{}
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

func (r *GeneratedStoreStocktakeLineResultTypeResolver) Total(ctx context.Context, obj *StoreStocktakeLineResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_stocktake_lines", ctx), &StoreStocktakeLine{})
}

func (r *GeneratedStoreStocktakeLineResultTypeResolver) TotalPage(ctx context.Context, obj *StoreStocktakeLineResultType) (count int, err error) {
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

func (r *GeneratedStoreStocktakeLineResultTypeResolver) CurrentPage(ctx context.Context, obj *StoreStocktakeLineResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStoreStocktakeLineResultTypeResolver) PerPage(ctx context.Context, obj *StoreStocktakeLineResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStoreStocktakeLineResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStocktakeLineResolver) Stocktake(ctx context.Context, obj *StoreStocktakeLine) (res *StoreStocktake, err error) {
	return r.Handlers.StoreStocktakeLineStocktake(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakeLineStocktakeHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (items *StoreStocktake, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StoreStocktake"); err != nil {
		return items, errors.New("StoreStocktake " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.StocktakeID

	if objKey != "" {
		item, _ := loaders["StoreStocktake"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StoreStocktake)

		if items == nil {
			items = &StoreStocktake{}
		}

	}

	return
}

func (r *GeneratedStoreStocktakeLineResolver) Batch(ctx context.Context, obj *StoreStocktakeLine) (res *StoreInventoryBatch, err error) {
	return r.Handlers.StoreStocktakeLineBatch(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakeLineBatchHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (items *StoreInventoryBatch, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StoreInventoryBatch"); err != nil {
		return items, errors.New("StoreInventoryBatch " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.BatchID

	if objKey != "" {
		item, _ := loaders["StoreInventoryBatch"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StoreInventoryBatch)

		if items == nil {
			items = &StoreInventoryBatch{}
		}

	}

	return
}

func (r *GeneratedStoreStocktakeLineResolver) Package(ctx context.Context, obj *StoreStocktakeLine) (res *ProductPackage, err error) {
	return r.Handlers.StoreStocktakeLinePackage(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakeLinePackageHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (items *ProductPackage, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackage"); err != nil {
		return items, errors.New("ProductPackage " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.PackageID

	if objKey != "" {
		item, _ := loaders["ProductPackage"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*ProductPackage)

		if items == nil {
			items = &ProductPackage{}
		}

	}

	return
}

func (r *GeneratedStoreStocktakeLineResolver) Movements(ctx context.Context, obj *StoreStocktakeLine) (res []*StoreStockMovement, err error) {
	return r.Handlers.StoreStocktakeLineMovements(ctx, r.GeneratedResolver, obj)
}
func StoreStocktakeLineMovementsHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStocktakeLine) (items []*StoreStockMovement, err error) {

	items = []*StoreStockMovement{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Movements"); err != nil {
		return items, errors.New("Movements " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StoreStockMovementStocktakeLine"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StoreStockMovement{}
	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	return
}

func (r *GeneratedStoreStocktakeLineResolver) MovementsIds(ctx context.Context, obj *StoreStocktakeLine) (ids []string, err error) {

	items := []*StoreStockMovement{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StocktakeLineAndStoreStockMovementIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StoreStockMovement)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryStoreStockMovementHandlerOptions struct {
	ID     *string
	Filter *StoreStockMovementFilterType
}

func (r *GeneratedQueryResolver) StoreStockMovement(ctx context.Context, id *string, filter *StoreStockMovementFilterType) (*StoreStockMovement, error) {
	opts := QueryStoreStockMovementHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStoreStockMovement(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStockMovementHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockMovementHandlerOptions) (*StoreStockMovement, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StoreStockMovementQueryFilter{}
	rt := &StoreStockMovementResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_stock_movements", ctx)+".id = ?", *opts.ID)
	}

	var items []*StoreStockMovement
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stock_movements", ctx),
		Preloaders: []string{},
		Item:       &StoreStockMovement{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StoreStockMovement"}
	}
	return items[0], err
}

type QueryStoreStockMovementsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StoreStockMovementSortType
	Filter      *StoreStockMovementFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StoreStockMovements(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StoreStockMovementSortType, filter *StoreStockMovementFilterType, rand *bool) (*StoreStockMovementResultType, error) {
	opts := QueryStoreStockMovementsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStoreStockMovements(ctx, r.GeneratedResolver, opts)
}
func QueryStoreStockMovementsHandler(ctx context.Context, r *GeneratedResolver, opts QueryStoreStockMovementsHandlerOptions) (*StoreStockMovementResultType, error) {
	query := StoreStockMovementQueryFilter{opts.Q}

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

	return &StoreStockMovementResultType{
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

type GeneratedStoreStockMovementResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStockMovementResultTypeResolver) Data(ctx context.Context, obj *StoreStockMovementResultType) (items []*StoreStockMovement, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_stock_movements", ctx),
		Preloaders: []string{},
		Item:       &StoreStockMovement{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StoreStockMovement{}
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

func (r *GeneratedStoreStockMovementResultTypeResolver) Total(ctx context.Context, obj *StoreStockMovementResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_stock_movements", ctx), &StoreStockMovement{})
}

func (r *GeneratedStoreStockMovementResultTypeResolver) TotalPage(ctx context.Context, obj *StoreStockMovementResultType) (count int, err error) {
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

func (r *GeneratedStoreStockMovementResultTypeResolver) CurrentPage(ctx context.Context, obj *StoreStockMovementResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStoreStockMovementResultTypeResolver) PerPage(ctx context.Context, obj *StoreStockMovementResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStoreStockMovementResolver struct{ *GeneratedResolver }

func (r *GeneratedStoreStockMovementResolver) Store(ctx context.Context, obj *StoreStockMovement) (res *Store, err error) {
	return r.Handlers.StoreStockMovementStore(ctx, r.GeneratedResolver, obj)
}
func StoreStockMovementStoreHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (items *Store, err error) {

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

func (r *GeneratedStoreStockMovementResolver) Batch(ctx context.Context, obj *StoreStockMovement) (res *StoreInventoryBatch, err error) {
	return r.Handlers.StoreStockMovementBatch(ctx, r.GeneratedResolver, obj)
}
func StoreStockMovementBatchHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (items *StoreInventoryBatch, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StoreInventoryBatch"); err != nil {
		return items, errors.New("StoreInventoryBatch " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.BatchID

	if objKey != "" {
		item, _ := loaders["StoreInventoryBatch"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StoreInventoryBatch)

		if items == nil {
			items = &StoreInventoryBatch{}
		}

	}

	return
}

func (r *GeneratedStoreStockMovementResolver) SourcePackage(ctx context.Context, obj *StoreStockMovement) (res *ProductPackage, err error) {
	return r.Handlers.StoreStockMovementSourcePackage(ctx, r.GeneratedResolver, obj)
}
func StoreStockMovementSourcePackageHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (items *ProductPackage, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackage"); err != nil {
		return items, errors.New("ProductPackage " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.SourcePackageID

	if objKey != nil {
		item, _ := loaders["ProductPackage"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductPackage)

	}

	return
}

func (r *GeneratedStoreStockMovementResolver) TargetPackage(ctx context.Context, obj *StoreStockMovement) (res *ProductPackage, err error) {
	return r.Handlers.StoreStockMovementTargetPackage(ctx, r.GeneratedResolver, obj)
}
func StoreStockMovementTargetPackageHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (items *ProductPackage, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "ProductPackage"); err != nil {
		return items, errors.New("ProductPackage " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.TargetPackageID

	if objKey != nil {
		item, _ := loaders["ProductPackage"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*ProductPackage)

	}

	return
}

func (r *GeneratedStoreStockMovementResolver) StocktakeLine(ctx context.Context, obj *StoreStockMovement) (res *StoreStocktakeLine, err error) {
	return r.Handlers.StoreStockMovementStocktakeLine(ctx, r.GeneratedResolver, obj)
}
func StoreStockMovementStocktakeLineHandler(ctx context.Context, r *GeneratedResolver, obj *StoreStockMovement) (items *StoreStocktakeLine, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StoreStocktakeLine"); err != nil {
		return items, errors.New("StoreStocktakeLine " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.StocktakeLineID

	if objKey != nil {
		item, _ := loaders["StoreStocktakeLine"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*StoreStocktakeLine)

	}

	return
}

type QueryStorePromotionHandlerOptions struct {
	ID     *string
	Filter *StorePromotionFilterType
}

func (r *GeneratedQueryResolver) StorePromotion(ctx context.Context, id *string, filter *StorePromotionFilterType) (*StorePromotion, error) {
	opts := QueryStorePromotionHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStorePromotion(ctx, r.GeneratedResolver, opts)
}
func QueryStorePromotionHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionHandlerOptions) (*StorePromotion, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StorePromotionQueryFilter{}
	rt := &StorePromotionResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_promotions", ctx)+".id = ?", *opts.ID)
	}

	var items []*StorePromotion
	giOpts := GetItemsOptions{
		Alias:      TableName("store_promotions", ctx),
		Preloaders: []string{},
		Item:       &StorePromotion{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StorePromotion"}
	}
	return items[0], err
}

type QueryStorePromotionsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StorePromotionSortType
	Filter      *StorePromotionFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StorePromotions(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StorePromotionSortType, filter *StorePromotionFilterType, rand *bool) (*StorePromotionResultType, error) {
	opts := QueryStorePromotionsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStorePromotions(ctx, r.GeneratedResolver, opts)
}
func QueryStorePromotionsHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionsHandlerOptions) (*StorePromotionResultType, error) {
	query := StorePromotionQueryFilter{opts.Q}

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

	return &StorePromotionResultType{
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

type GeneratedStorePromotionResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePromotionResultTypeResolver) Data(ctx context.Context, obj *StorePromotionResultType) (items []*StorePromotion, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_promotions", ctx),
		Preloaders: []string{},
		Item:       &StorePromotion{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StorePromotion{}
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

func (r *GeneratedStorePromotionResultTypeResolver) Total(ctx context.Context, obj *StorePromotionResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_promotions", ctx), &StorePromotion{})
}

func (r *GeneratedStorePromotionResultTypeResolver) TotalPage(ctx context.Context, obj *StorePromotionResultType) (count int, err error) {
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

func (r *GeneratedStorePromotionResultTypeResolver) CurrentPage(ctx context.Context, obj *StorePromotionResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStorePromotionResultTypeResolver) PerPage(ctx context.Context, obj *StorePromotionResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStorePromotionResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePromotionResolver) Store(ctx context.Context, obj *StorePromotion) (res *Store, err error) {
	return r.Handlers.StorePromotionStore(ctx, r.GeneratedResolver, obj)
}
func StorePromotionStoreHandler(ctx context.Context, r *GeneratedResolver, obj *StorePromotion) (items *Store, err error) {

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

func (r *GeneratedStorePromotionResolver) Targets(ctx context.Context, obj *StorePromotion) (res []*StorePromotionTarget, err error) {
	return r.Handlers.StorePromotionTargets(ctx, r.GeneratedResolver, obj)
}
func StorePromotionTargetsHandler(ctx context.Context, r *GeneratedResolver, obj *StorePromotion) (items []*StorePromotionTarget, err error) {

	items = []*StorePromotionTarget{}

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "Targets"); err != nil {
		return items, errors.New("Targets " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["StorePromotionTargetPromotion"].Load(ctx, dataloader.StringKey(obj.ID))()
	items = []*StorePromotionTarget{}
	if item != nil {
		items = item.([]*StorePromotionTarget)
	}

	return
}

func (r *GeneratedStorePromotionResolver) TargetsIds(ctx context.Context, obj *StorePromotion) (ids []string, err error) {

	items := []*StorePromotionTarget{}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	item, _ := loaders["PromotionAndStorePromotionTargetIds"].Load(ctx, dataloader.StringKey(obj.ID))()

	if item != nil {
		items = item.([]*StorePromotionTarget)
	}

	for _, v := range items {
		ids = append(ids, v.ID)
	}

	return
}

type QueryStorePromotionTargetHandlerOptions struct {
	ID     *string
	Filter *StorePromotionTargetFilterType
}

func (r *GeneratedQueryResolver) StorePromotionTarget(ctx context.Context, id *string, filter *StorePromotionTargetFilterType) (*StorePromotionTarget, error) {
	opts := QueryStorePromotionTargetHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryStorePromotionTarget(ctx, r.GeneratedResolver, opts)
}
func QueryStorePromotionTargetHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionTargetHandlerOptions) (*StorePromotionTarget, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := StorePromotionTargetQueryFilter{}
	rt := &StorePromotionTargetResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("store_promotion_targets", ctx)+".id = ?", *opts.ID)
	}

	var items []*StorePromotionTarget
	giOpts := GetItemsOptions{
		Alias:      TableName("store_promotion_targets", ctx),
		Preloaders: []string{},
		Item:       &StorePromotionTarget{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "StorePromotionTarget"}
	}
	return items[0], err
}

type QueryStorePromotionTargetsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*StorePromotionTargetSortType
	Filter      *StorePromotionTargetFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) StorePromotionTargets(ctx context.Context, current_page *int, per_page *int, q *string, sort []*StorePromotionTargetSortType, filter *StorePromotionTargetFilterType, rand *bool) (*StorePromotionTargetResultType, error) {
	opts := QueryStorePromotionTargetsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryStorePromotionTargets(ctx, r.GeneratedResolver, opts)
}
func QueryStorePromotionTargetsHandler(ctx context.Context, r *GeneratedResolver, opts QueryStorePromotionTargetsHandlerOptions) (*StorePromotionTargetResultType, error) {
	query := StorePromotionTargetQueryFilter{opts.Q}

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

	return &StorePromotionTargetResultType{
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

type GeneratedStorePromotionTargetResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePromotionTargetResultTypeResolver) Data(ctx context.Context, obj *StorePromotionTargetResultType) (items []*StorePromotionTarget, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("store_promotion_targets", ctx),
		Preloaders: []string{},
		Item:       &StorePromotionTarget{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*StorePromotionTarget{}
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

func (r *GeneratedStorePromotionTargetResultTypeResolver) Total(ctx context.Context, obj *StorePromotionTargetResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("store_promotion_targets", ctx), &StorePromotionTarget{})
}

func (r *GeneratedStorePromotionTargetResultTypeResolver) TotalPage(ctx context.Context, obj *StorePromotionTargetResultType) (count int, err error) {
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

func (r *GeneratedStorePromotionTargetResultTypeResolver) CurrentPage(ctx context.Context, obj *StorePromotionTargetResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedStorePromotionTargetResultTypeResolver) PerPage(ctx context.Context, obj *StorePromotionTargetResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedStorePromotionTargetResolver struct{ *GeneratedResolver }

func (r *GeneratedStorePromotionTargetResolver) Promotion(ctx context.Context, obj *StorePromotionTarget) (res *StorePromotion, err error) {
	return r.Handlers.StorePromotionTargetPromotion(ctx, r.GeneratedResolver, obj)
}
func StorePromotionTargetPromotionHandler(ctx context.Context, r *GeneratedResolver, obj *StorePromotionTarget) (items *StorePromotion, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StorePromotion"); err != nil {
		return items, errors.New("StorePromotion " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.PromotionID

	if objKey != "" {
		item, _ := loaders["StorePromotion"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StorePromotion)

		if items == nil {
			items = &StorePromotion{}
		}

	}

	return
}

func (r *GeneratedStorePromotionTargetResolver) Offer(ctx context.Context, obj *StorePromotionTarget) (res *StorePackageOffer, err error) {
	return r.Handlers.StorePromotionTargetOffer(ctx, r.GeneratedResolver, obj)
}
func StorePromotionTargetOfferHandler(ctx context.Context, r *GeneratedResolver, obj *StorePromotionTarget) (items *StorePackageOffer, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "StorePackageOffer"); err != nil {
		return items, errors.New("StorePackageOffer " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.OfferID

	if objKey != "" {
		item, _ := loaders["StorePackageOffer"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*StorePackageOffer)

		if items == nil {
			items = &StorePackageOffer{}
		}

	}

	return
}

type QueryCustomerCouponGrantHandlerOptions struct {
	ID     *string
	Filter *CustomerCouponGrantFilterType
}

func (r *GeneratedQueryResolver) CustomerCouponGrant(ctx context.Context, id *string, filter *CustomerCouponGrantFilterType) (*CustomerCouponGrant, error) {
	opts := QueryCustomerCouponGrantHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryCustomerCouponGrant(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerCouponGrantHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponGrantHandlerOptions) (*CustomerCouponGrant, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := CustomerCouponGrantQueryFilter{}
	rt := &CustomerCouponGrantResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("customer_coupon_grants", ctx)+".id = ?", *opts.ID)
	}

	var items []*CustomerCouponGrant
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_coupon_grants", ctx),
		Preloaders: []string{},
		Item:       &CustomerCouponGrant{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "CustomerCouponGrant"}
	}
	return items[0], err
}

type QueryCustomerCouponGrantsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*CustomerCouponGrantSortType
	Filter      *CustomerCouponGrantFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) CustomerCouponGrants(ctx context.Context, current_page *int, per_page *int, q *string, sort []*CustomerCouponGrantSortType, filter *CustomerCouponGrantFilterType, rand *bool) (*CustomerCouponGrantResultType, error) {
	opts := QueryCustomerCouponGrantsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryCustomerCouponGrants(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerCouponGrantsHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponGrantsHandlerOptions) (*CustomerCouponGrantResultType, error) {
	query := CustomerCouponGrantQueryFilter{opts.Q}

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

	return &CustomerCouponGrantResultType{
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

type GeneratedCustomerCouponGrantResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerCouponGrantResultTypeResolver) Data(ctx context.Context, obj *CustomerCouponGrantResultType) (items []*CustomerCouponGrant, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_coupon_grants", ctx),
		Preloaders: []string{},
		Item:       &CustomerCouponGrant{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*CustomerCouponGrant{}
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

func (r *GeneratedCustomerCouponGrantResultTypeResolver) Total(ctx context.Context, obj *CustomerCouponGrantResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("customer_coupon_grants", ctx), &CustomerCouponGrant{})
}

func (r *GeneratedCustomerCouponGrantResultTypeResolver) TotalPage(ctx context.Context, obj *CustomerCouponGrantResultType) (count int, err error) {
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

func (r *GeneratedCustomerCouponGrantResultTypeResolver) CurrentPage(ctx context.Context, obj *CustomerCouponGrantResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedCustomerCouponGrantResultTypeResolver) PerPage(ctx context.Context, obj *CustomerCouponGrantResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedCustomerCouponGrantResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerCouponGrantResolver) Member(ctx context.Context, obj *CustomerCouponGrant) (res *CustomerMember, err error) {
	return r.Handlers.CustomerCouponGrantMember(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponGrantMemberHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponGrant) (items *CustomerMember, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerMember"); err != nil {
		return items, errors.New("CustomerMember " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.MemberID

	if objKey != "" {
		item, _ := loaders["CustomerMember"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*CustomerMember)

		if items == nil {
			items = &CustomerMember{}
		}

	}

	return
}

func (r *GeneratedCustomerCouponGrantResolver) Template(ctx context.Context, obj *CustomerCouponGrant) (res *CustomerCouponTemplate, err error) {
	return r.Handlers.CustomerCouponGrantTemplate(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponGrantTemplateHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponGrant) (items *CustomerCouponTemplate, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerCouponTemplate"); err != nil {
		return items, errors.New("CustomerCouponTemplate " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.TemplateID

	if objKey != "" {
		item, _ := loaders["CustomerCouponTemplate"].Load(ctx, dataloader.StringKey(objKey))()

		items, _ = item.(*CustomerCouponTemplate)

		if items == nil {
			items = &CustomerCouponTemplate{}
		}

	}

	return
}

type QueryCustomerCouponDistributionJobHandlerOptions struct {
	ID     *string
	Filter *CustomerCouponDistributionJobFilterType
}

func (r *GeneratedQueryResolver) CustomerCouponDistributionJob(ctx context.Context, id *string, filter *CustomerCouponDistributionJobFilterType) (*CustomerCouponDistributionJob, error) {
	opts := QueryCustomerCouponDistributionJobHandlerOptions{
		ID:     id,
		Filter: filter,
	}
	return r.Handlers.QueryCustomerCouponDistributionJob(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerCouponDistributionJobHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponDistributionJobHandlerOptions) (*CustomerCouponDistributionJob, error) {
	selection := []ast.Selection{}
	func() {
		defer func() { recover() }()
		for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
			selection = append(selection, f.Field)
		}
	}()
	selectionSet := ast.SelectionSet(selection)

	query := CustomerCouponDistributionJobQueryFilter{}
	rt := &CustomerCouponDistributionJobResultType{
		EntityResultType: EntityResultType{
			Query:        &query,
			Filter:       opts.Filter,
			SelectionSet: &selectionSet,
		},
	}
	qb := r.DB.Query()
	if opts.ID != nil {
		qb = qb.Where(TableName("customer_coupon_distribution_jobs", ctx)+".id = ?", *opts.ID)
	}

	var items []*CustomerCouponDistributionJob
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_coupon_distribution_jobs", ctx),
		Preloaders: []string{},
		Item:       &CustomerCouponDistributionJob{},
	}
	err := rt.GetData(ctx, qb, giOpts, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, &NotFoundError{Entity: "CustomerCouponDistributionJob"}
	}
	return items[0], err
}

type QueryCustomerCouponDistributionJobsHandlerOptions struct {
	CurrentPage *int
	PerPage     *int
	Q           *string
	Sort        []*CustomerCouponDistributionJobSortType
	Filter      *CustomerCouponDistributionJobFilterType
	Rand        *bool
}

func (r *GeneratedQueryResolver) CustomerCouponDistributionJobs(ctx context.Context, current_page *int, per_page *int, q *string, sort []*CustomerCouponDistributionJobSortType, filter *CustomerCouponDistributionJobFilterType, rand *bool) (*CustomerCouponDistributionJobResultType, error) {
	opts := QueryCustomerCouponDistributionJobsHandlerOptions{
		CurrentPage: current_page,
		PerPage:     per_page,
		Q:           q,
		Sort:        sort,
		Filter:      filter,
		Rand:        rand,
	}
	return r.Handlers.QueryCustomerCouponDistributionJobs(ctx, r.GeneratedResolver, opts)
}
func QueryCustomerCouponDistributionJobsHandler(ctx context.Context, r *GeneratedResolver, opts QueryCustomerCouponDistributionJobsHandlerOptions) (*CustomerCouponDistributionJobResultType, error) {
	query := CustomerCouponDistributionJobQueryFilter{opts.Q}

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

	return &CustomerCouponDistributionJobResultType{
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

type GeneratedCustomerCouponDistributionJobResultTypeResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerCouponDistributionJobResultTypeResolver) Data(ctx context.Context, obj *CustomerCouponDistributionJobResultType) (items []*CustomerCouponDistributionJob, err error) {
	giOpts := GetItemsOptions{
		Alias:      TableName("customer_coupon_distribution_jobs", ctx),
		Preloaders: []string{},
		Item:       &CustomerCouponDistributionJob{},
	}
	err = obj.GetData(ctx, r.DB.db, giOpts, &items)

	uniqueItems := []*CustomerCouponDistributionJob{}
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

func (r *GeneratedCustomerCouponDistributionJobResultTypeResolver) Total(ctx context.Context, obj *CustomerCouponDistributionJobResultType) (count int, err error) {
	return obj.GetTotal(ctx, r.DB.db, TableName("customer_coupon_distribution_jobs", ctx), &CustomerCouponDistributionJob{})
}

func (r *GeneratedCustomerCouponDistributionJobResultTypeResolver) TotalPage(ctx context.Context, obj *CustomerCouponDistributionJobResultType) (count int, err error) {
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

func (r *GeneratedCustomerCouponDistributionJobResultTypeResolver) CurrentPage(ctx context.Context, obj *CustomerCouponDistributionJobResultType) (count int, err error) {
	return int(*obj.EntityResultType.CurrentPage), nil
}

func (r *GeneratedCustomerCouponDistributionJobResultTypeResolver) PerPage(ctx context.Context, obj *CustomerCouponDistributionJobResultType) (count int, err error) {
	return int(*obj.EntityResultType.PerPage), nil
}

type GeneratedCustomerCouponDistributionJobResolver struct{ *GeneratedResolver }

func (r *GeneratedCustomerCouponDistributionJobResolver) Template(ctx context.Context, obj *CustomerCouponDistributionJob) (res *CustomerCouponTemplate, err error) {
	return r.Handlers.CustomerCouponDistributionJobTemplate(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponDistributionJobTemplateHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponDistributionJob) (items *CustomerCouponTemplate, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerCouponTemplate"); err != nil {
		return items, errors.New("CustomerCouponTemplate " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.TemplateID

	if objKey != nil {
		item, _ := loaders["CustomerCouponTemplate"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*CustomerCouponTemplate)

	}

	return
}

func (r *GeneratedCustomerCouponDistributionJobResolver) Member(ctx context.Context, obj *CustomerCouponDistributionJob) (res *CustomerMember, err error) {
	return r.Handlers.CustomerCouponDistributionJobMember(ctx, r.GeneratedResolver, obj)
}
func CustomerCouponDistributionJobMemberHandler(ctx context.Context, r *GeneratedResolver, obj *CustomerCouponDistributionJob) (items *CustomerMember, err error) {

	// 判断是否有详情权限
	if err := auth.CheckAuthorization(ctx, "CustomerMember"); err != nil {
		return items, errors.New("CustomerMember " + err.Error())
	}

	loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
	objKey := obj.MemberID

	if objKey != nil {
		item, _ := loaders["CustomerMember"].Load(ctx, dataloader.StringKey(*objKey))()

		items, _ = item.(*CustomerMember)

	}

	return
}
