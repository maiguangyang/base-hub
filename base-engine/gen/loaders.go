package gen

import (
	"context"
	"errors"

	"github.com/graph-gophers/dataloader"
	"gorm.io/gorm"
)

func GetLoaders(db *DB) map[string]*dataloader.Loader {
	loaders := map[string]*dataloader.Loader{}

	accountsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Account{}
		selects := GetFieldsRequested(ctx, TableName("accounts", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("accounts", ctx)+".*") == -1 && IndexOf(selects, TableName("accounts", ctx)+".id") == -1 {
			selects = append(selects, TableName("accounts", ctx)+".id")
		}

		res := db.Query().Table(TableName("accounts", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]Account, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Account with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["Account"] = dataloader.NewBatchedLoader(accountsBatchFn, dataloader.WithClearCacheOnBatch())

	organizationsInitialAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Organization{}
		selects := GetFieldsRequested(ctx, TableName("organizations", ctx))

		if IndexOf(selects, TableName("organizations", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("organizations", ctx)+".id") == -1 {
				selects = append(selects, "organizations"+".id")
			}

			if IndexOf(selects, TableName("organizations", ctx)+".initial_account_id") == -1 {
				selects = append(selects, TableName("organizations", ctx)+".initial_account_id")
			}
		}

		res := db.Query().Table(TableName("organizations", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "initial_account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Organization, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.InitialAccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Organization{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Organization with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["OrganizationInitialAccount"] = dataloader.NewBatchedLoader(organizationsInitialAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["InitialAccountAndOrganizationIds"] = dataloader.NewBatchedLoader(organizationsInitialAccountBatchFn, dataloader.WithClearCacheOnBatch())

	organizationsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Organization{}
		selects := GetFieldsRequested(ctx, TableName("organizations", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("organizations", ctx)+".*") == -1 && IndexOf(selects, TableName("organizations", ctx)+".id") == -1 {
			selects = append(selects, TableName("organizations", ctx)+".id")
		}

		res := db.Query().Table(TableName("organizations", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]Organization, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Organization with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["Organization"] = dataloader.NewBatchedLoader(organizationsBatchFn, dataloader.WithClearCacheOnBatch())

	operator_membershipsAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]OperatorMembership{}
		selects := GetFieldsRequested(ctx, TableName("operator_memberships", ctx))

		if IndexOf(selects, TableName("operator_memberships", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("operator_memberships", ctx)+".id") == -1 {
				selects = append(selects, "operator_memberships"+".id")
			}

			if IndexOf(selects, TableName("operator_memberships", ctx)+".account_id") == -1 {
				selects = append(selects, TableName("operator_memberships", ctx)+".account_id")
			}
		}

		res := db.Query().Table(TableName("operator_memberships", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*OperatorMembership, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.AccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*OperatorMembership{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("OperatorMembership with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["OperatorMembershipAccount"] = dataloader.NewBatchedLoader(operator_membershipsAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["AccountAndOperatorMembershipIds"] = dataloader.NewBatchedLoader(operator_membershipsAccountBatchFn, dataloader.WithClearCacheOnBatch())

	operator_membershipsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]OperatorMembership{}
		selects := GetFieldsRequested(ctx, TableName("operator_memberships", ctx))

		if IndexOf(selects, TableName("operator_memberships", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("operator_memberships", ctx)+".id") == -1 {
				selects = append(selects, "operator_memberships"+".id")
			}

			if IndexOf(selects, TableName("operator_memberships", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("operator_memberships", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("operator_memberships", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*OperatorMembership, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*OperatorMembership{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("OperatorMembership with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["OperatorMembershipOrganization"] = dataloader.NewBatchedLoader(operator_membershipsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndOperatorMembershipIds"] = dataloader.NewBatchedLoader(operator_membershipsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	operator_membershipsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]OperatorMembership{}
		selects := GetFieldsRequested(ctx, TableName("operator_memberships", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("operator_memberships", ctx)+".*") == -1 && IndexOf(selects, TableName("operator_memberships", ctx)+".id") == -1 {
			selects = append(selects, TableName("operator_memberships", ctx)+".id")
		}

		res := db.Query().Table(TableName("operator_memberships", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]OperatorMembership, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("OperatorMembership with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["OperatorMembership"] = dataloader.NewBatchedLoader(operator_membershipsBatchFn, dataloader.WithClearCacheOnBatch())

	permissionsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Permission{}
		selects := GetFieldsRequested(ctx, TableName("permissions", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("permissions", ctx)+".*") == -1 && IndexOf(selects, TableName("permissions", ctx)+".id") == -1 {
			selects = append(selects, TableName("permissions", ctx)+".id")
		}

		res := db.Query().Table(TableName("permissions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]Permission, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Permission with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["Permission"] = dataloader.NewBatchedLoader(permissionsBatchFn, dataloader.WithClearCacheOnBatch())

	operator_rolesOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]OperatorRole{}
		selects := GetFieldsRequested(ctx, TableName("operator_roles", ctx))

		if IndexOf(selects, TableName("operator_roles", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("operator_roles", ctx)+".id") == -1 {
				selects = append(selects, "operator_roles"+".id")
			}

			if IndexOf(selects, TableName("operator_roles", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("operator_roles", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("operator_roles", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*OperatorRole, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*OperatorRole{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("OperatorRole with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["OperatorRoleOrganization"] = dataloader.NewBatchedLoader(operator_rolesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndOperatorRoleIds"] = dataloader.NewBatchedLoader(operator_rolesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	operator_rolesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]OperatorRole{}
		selects := GetFieldsRequested(ctx, TableName("operator_roles", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("operator_roles", ctx)+".*") == -1 && IndexOf(selects, TableName("operator_roles", ctx)+".id") == -1 {
			selects = append(selects, TableName("operator_roles", ctx)+".id")
		}

		res := db.Query().Table(TableName("operator_roles", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]OperatorRole, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("OperatorRole with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["OperatorRole"] = dataloader.NewBatchedLoader(operator_rolesBatchFn, dataloader.WithClearCacheOnBatch())

	storesOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Store{}
		selects := GetFieldsRequested(ctx, TableName("stores", ctx))

		if IndexOf(selects, TableName("stores", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("stores", ctx)+".id") == -1 {
				selects = append(selects, "stores"+".id")
			}

			if IndexOf(selects, TableName("stores", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("stores", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("stores", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Store, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Store{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Store with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreOrganization"] = dataloader.NewBatchedLoader(storesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndStoreIds"] = dataloader.NewBatchedLoader(storesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	storesReviewedByAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Store{}
		selects := GetFieldsRequested(ctx, TableName("stores", ctx))

		if IndexOf(selects, TableName("stores", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("stores", ctx)+".id") == -1 {
				selects = append(selects, "stores"+".id")
			}

			if IndexOf(selects, TableName("stores", ctx)+".reviewed_by_account_id") == -1 {
				selects = append(selects, TableName("stores", ctx)+".reviewed_by_account_id")
			}
		}

		res := db.Query().Table(TableName("stores", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "reviewed_by_account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Store, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ReviewedByAccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Store{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Store with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreReviewedByAccount"] = dataloader.NewBatchedLoader(storesReviewedByAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ReviewedByAccountAndStoreIds"] = dataloader.NewBatchedLoader(storesReviewedByAccountBatchFn, dataloader.WithClearCacheOnBatch())

	storesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Store{}
		selects := GetFieldsRequested(ctx, TableName("stores", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("stores", ctx)+".*") == -1 && IndexOf(selects, TableName("stores", ctx)+".id") == -1 {
			selects = append(selects, TableName("stores", ctx)+".id")
		}

		res := db.Query().Table(TableName("stores", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]Store, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Store with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["Store"] = dataloader.NewBatchedLoader(storesBatchFn, dataloader.WithClearCacheOnBatch())

	sessionsAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Session{}
		selects := GetFieldsRequested(ctx, TableName("sessions", ctx))

		if IndexOf(selects, TableName("sessions", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("sessions", ctx)+".id") == -1 {
				selects = append(selects, "sessions"+".id")
			}

			if IndexOf(selects, TableName("sessions", ctx)+".account_id") == -1 {
				selects = append(selects, TableName("sessions", ctx)+".account_id")
			}
		}

		res := db.Query().Table(TableName("sessions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Session, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.AccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Session{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Session with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["SessionAccount"] = dataloader.NewBatchedLoader(sessionsAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["AccountAndSessionIds"] = dataloader.NewBatchedLoader(sessionsAccountBatchFn, dataloader.WithClearCacheOnBatch())

	sessionsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Session{}
		selects := GetFieldsRequested(ctx, TableName("sessions", ctx))

		if IndexOf(selects, TableName("sessions", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("sessions", ctx)+".id") == -1 {
				selects = append(selects, "sessions"+".id")
			}

			if IndexOf(selects, TableName("sessions", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("sessions", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("sessions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Session, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Session{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Session with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["SessionOrganization"] = dataloader.NewBatchedLoader(sessionsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndSessionIds"] = dataloader.NewBatchedLoader(sessionsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	sessionsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Session{}
		selects := GetFieldsRequested(ctx, TableName("sessions", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("sessions", ctx)+".*") == -1 && IndexOf(selects, TableName("sessions", ctx)+".id") == -1 {
			selects = append(selects, TableName("sessions", ctx)+".id")
		}

		res := db.Query().Table(TableName("sessions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]Session, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Session with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["Session"] = dataloader.NewBatchedLoader(sessionsBatchFn, dataloader.WithClearCacheOnBatch())

	membership_invitationsMembershipBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]MembershipInvitation{}
		selects := GetFieldsRequested(ctx, TableName("membership_invitations", ctx))

		if IndexOf(selects, TableName("membership_invitations", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("membership_invitations", ctx)+".id") == -1 {
				selects = append(selects, "membership_invitations"+".id")
			}

			if IndexOf(selects, TableName("membership_invitations", ctx)+".membership_id") == -1 {
				selects = append(selects, TableName("membership_invitations", ctx)+".membership_id")
			}
		}

		res := db.Query().Table(TableName("membership_invitations", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "membership_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*MembershipInvitation, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.MembershipID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*MembershipInvitation{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("MembershipInvitation with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["MembershipInvitationMembership"] = dataloader.NewBatchedLoader(membership_invitationsMembershipBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["MembershipAndMembershipInvitationIds"] = dataloader.NewBatchedLoader(membership_invitationsMembershipBatchFn, dataloader.WithClearCacheOnBatch())

	membership_invitationsInvitedByAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]MembershipInvitation{}
		selects := GetFieldsRequested(ctx, TableName("membership_invitations", ctx))

		if IndexOf(selects, TableName("membership_invitations", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("membership_invitations", ctx)+".id") == -1 {
				selects = append(selects, "membership_invitations"+".id")
			}

			if IndexOf(selects, TableName("membership_invitations", ctx)+".invited_by_account_id") == -1 {
				selects = append(selects, TableName("membership_invitations", ctx)+".invited_by_account_id")
			}
		}

		res := db.Query().Table(TableName("membership_invitations", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "invited_by_account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*MembershipInvitation, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.InvitedByAccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*MembershipInvitation{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("MembershipInvitation with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["MembershipInvitationInvitedByAccount"] = dataloader.NewBatchedLoader(membership_invitationsInvitedByAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["InvitedByAccountAndMembershipInvitationIds"] = dataloader.NewBatchedLoader(membership_invitationsInvitedByAccountBatchFn, dataloader.WithClearCacheOnBatch())

	membership_invitationsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]MembershipInvitation{}
		selects := GetFieldsRequested(ctx, TableName("membership_invitations", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("membership_invitations", ctx)+".*") == -1 && IndexOf(selects, TableName("membership_invitations", ctx)+".id") == -1 {
			selects = append(selects, TableName("membership_invitations", ctx)+".id")
		}

		res := db.Query().Table(TableName("membership_invitations", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]MembershipInvitation, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("MembershipInvitation with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["MembershipInvitation"] = dataloader.NewBatchedLoader(membership_invitationsBatchFn, dataloader.WithClearCacheOnBatch())

	audit_logsActorAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]AuditLog{}
		selects := GetFieldsRequested(ctx, TableName("audit_logs", ctx))

		if IndexOf(selects, TableName("audit_logs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("audit_logs", ctx)+".id") == -1 {
				selects = append(selects, "audit_logs"+".id")
			}

			if IndexOf(selects, TableName("audit_logs", ctx)+".actor_account_id") == -1 {
				selects = append(selects, TableName("audit_logs", ctx)+".actor_account_id")
			}
		}

		res := db.Query().Table(TableName("audit_logs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "actor_account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*AuditLog, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ActorAccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*AuditLog{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("AuditLog with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["AuditLogActorAccount"] = dataloader.NewBatchedLoader(audit_logsActorAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ActorAccountAndAuditLogIds"] = dataloader.NewBatchedLoader(audit_logsActorAccountBatchFn, dataloader.WithClearCacheOnBatch())

	audit_logsSessionBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]AuditLog{}
		selects := GetFieldsRequested(ctx, TableName("audit_logs", ctx))

		if IndexOf(selects, TableName("audit_logs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("audit_logs", ctx)+".id") == -1 {
				selects = append(selects, "audit_logs"+".id")
			}

			if IndexOf(selects, TableName("audit_logs", ctx)+".session_id") == -1 {
				selects = append(selects, TableName("audit_logs", ctx)+".session_id")
			}
		}

		res := db.Query().Table(TableName("audit_logs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "session_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*AuditLog, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.SessionID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*AuditLog{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("AuditLog with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["AuditLogSession"] = dataloader.NewBatchedLoader(audit_logsSessionBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["SessionAndAuditLogIds"] = dataloader.NewBatchedLoader(audit_logsSessionBatchFn, dataloader.WithClearCacheOnBatch())

	audit_logsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]AuditLog{}
		selects := GetFieldsRequested(ctx, TableName("audit_logs", ctx))

		if IndexOf(selects, TableName("audit_logs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("audit_logs", ctx)+".id") == -1 {
				selects = append(selects, "audit_logs"+".id")
			}

			if IndexOf(selects, TableName("audit_logs", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("audit_logs", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("audit_logs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*AuditLog, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*AuditLog{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("AuditLog with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["AuditLogOrganization"] = dataloader.NewBatchedLoader(audit_logsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndAuditLogIds"] = dataloader.NewBatchedLoader(audit_logsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	audit_logsStoreBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]AuditLog{}
		selects := GetFieldsRequested(ctx, TableName("audit_logs", ctx))

		if IndexOf(selects, TableName("audit_logs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("audit_logs", ctx)+".id") == -1 {
				selects = append(selects, "audit_logs"+".id")
			}

			if IndexOf(selects, TableName("audit_logs", ctx)+".store_id") == -1 {
				selects = append(selects, TableName("audit_logs", ctx)+".store_id")
			}
		}

		res := db.Query().Table(TableName("audit_logs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "store_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*AuditLog, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.StoreID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*AuditLog{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("AuditLog with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["AuditLogStore"] = dataloader.NewBatchedLoader(audit_logsStoreBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StoreAndAuditLogIds"] = dataloader.NewBatchedLoader(audit_logsStoreBatchFn, dataloader.WithClearCacheOnBatch())

	audit_logsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]AuditLog{}
		selects := GetFieldsRequested(ctx, TableName("audit_logs", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("audit_logs", ctx)+".*") == -1 && IndexOf(selects, TableName("audit_logs", ctx)+".id") == -1 {
			selects = append(selects, TableName("audit_logs", ctx)+".id")
		}

		res := db.Query().Table(TableName("audit_logs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]AuditLog, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("AuditLog with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["AuditLog"] = dataloader.NewBatchedLoader(audit_logsBatchFn, dataloader.WithClearCacheOnBatch())

	franchise_opening_recordsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]FranchiseOpeningRecord{}
		selects := GetFieldsRequested(ctx, TableName("franchise_opening_records", ctx))

		if IndexOf(selects, TableName("franchise_opening_records", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("franchise_opening_records", ctx)+".id") == -1 {
				selects = append(selects, "franchise_opening_records"+".id")
			}

			if IndexOf(selects, TableName("franchise_opening_records", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("franchise_opening_records", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("franchise_opening_records", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*FranchiseOpeningRecord, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*FranchiseOpeningRecord{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("FranchiseOpeningRecord with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["FranchiseOpeningRecordOrganization"] = dataloader.NewBatchedLoader(franchise_opening_recordsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndFranchiseOpeningRecordIds"] = dataloader.NewBatchedLoader(franchise_opening_recordsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	franchise_opening_recordsInitialAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]FranchiseOpeningRecord{}
		selects := GetFieldsRequested(ctx, TableName("franchise_opening_records", ctx))

		if IndexOf(selects, TableName("franchise_opening_records", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("franchise_opening_records", ctx)+".id") == -1 {
				selects = append(selects, "franchise_opening_records"+".id")
			}

			if IndexOf(selects, TableName("franchise_opening_records", ctx)+".initial_account_id") == -1 {
				selects = append(selects, TableName("franchise_opening_records", ctx)+".initial_account_id")
			}
		}

		res := db.Query().Table(TableName("franchise_opening_records", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "initial_account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*FranchiseOpeningRecord, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.InitialAccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*FranchiseOpeningRecord{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("FranchiseOpeningRecord with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["FranchiseOpeningRecordInitialAccount"] = dataloader.NewBatchedLoader(franchise_opening_recordsInitialAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["InitialAccountAndFranchiseOpeningRecordIds"] = dataloader.NewBatchedLoader(franchise_opening_recordsInitialAccountBatchFn, dataloader.WithClearCacheOnBatch())

	franchise_opening_recordsRecordedByAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]FranchiseOpeningRecord{}
		selects := GetFieldsRequested(ctx, TableName("franchise_opening_records", ctx))

		if IndexOf(selects, TableName("franchise_opening_records", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("franchise_opening_records", ctx)+".id") == -1 {
				selects = append(selects, "franchise_opening_records"+".id")
			}

			if IndexOf(selects, TableName("franchise_opening_records", ctx)+".recorded_by_account_id") == -1 {
				selects = append(selects, TableName("franchise_opening_records", ctx)+".recorded_by_account_id")
			}
		}

		res := db.Query().Table(TableName("franchise_opening_records", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "recorded_by_account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*FranchiseOpeningRecord, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.RecordedByAccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*FranchiseOpeningRecord{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("FranchiseOpeningRecord with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["FranchiseOpeningRecordRecordedByAccount"] = dataloader.NewBatchedLoader(franchise_opening_recordsRecordedByAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["RecordedByAccountAndFranchiseOpeningRecordIds"] = dataloader.NewBatchedLoader(franchise_opening_recordsRecordedByAccountBatchFn, dataloader.WithClearCacheOnBatch())

	franchise_opening_recordsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]FranchiseOpeningRecord{}
		selects := GetFieldsRequested(ctx, TableName("franchise_opening_records", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("franchise_opening_records", ctx)+".*") == -1 && IndexOf(selects, TableName("franchise_opening_records", ctx)+".id") == -1 {
			selects = append(selects, TableName("franchise_opening_records", ctx)+".id")
		}

		res := db.Query().Table(TableName("franchise_opening_records", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]FranchiseOpeningRecord, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("FranchiseOpeningRecord with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["FranchiseOpeningRecord"] = dataloader.NewBatchedLoader(franchise_opening_recordsBatchFn, dataloader.WithClearCacheOnBatch())

	global_payment_configsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]GlobalPaymentConfig{}
		selects := GetFieldsRequested(ctx, TableName("global_payment_configs", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("global_payment_configs", ctx)+".*") == -1 && IndexOf(selects, TableName("global_payment_configs", ctx)+".id") == -1 {
			selects = append(selects, TableName("global_payment_configs", ctx)+".id")
		}

		res := db.Query().Table(TableName("global_payment_configs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]GlobalPaymentConfig, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("GlobalPaymentConfig with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["GlobalPaymentConfig"] = dataloader.NewBatchedLoader(global_payment_configsBatchFn, dataloader.WithClearCacheOnBatch())

	franchise_payment_configsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]FranchisePaymentConfig{}
		selects := GetFieldsRequested(ctx, TableName("franchise_payment_configs", ctx))

		if IndexOf(selects, TableName("franchise_payment_configs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("franchise_payment_configs", ctx)+".id") == -1 {
				selects = append(selects, "franchise_payment_configs"+".id")
			}

			if IndexOf(selects, TableName("franchise_payment_configs", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("franchise_payment_configs", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("franchise_payment_configs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*FranchisePaymentConfig, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*FranchisePaymentConfig{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("FranchisePaymentConfig with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["FranchisePaymentConfigOrganization"] = dataloader.NewBatchedLoader(franchise_payment_configsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndFranchisePaymentConfigIds"] = dataloader.NewBatchedLoader(franchise_payment_configsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	franchise_payment_configsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]FranchisePaymentConfig{}
		selects := GetFieldsRequested(ctx, TableName("franchise_payment_configs", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("franchise_payment_configs", ctx)+".*") == -1 && IndexOf(selects, TableName("franchise_payment_configs", ctx)+".id") == -1 {
			selects = append(selects, TableName("franchise_payment_configs", ctx)+".id")
		}

		res := db.Query().Table(TableName("franchise_payment_configs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]FranchisePaymentConfig, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("FranchisePaymentConfig with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["FranchisePaymentConfig"] = dataloader.NewBatchedLoader(franchise_payment_configsBatchFn, dataloader.WithClearCacheOnBatch())

	store_payment_configsStoreBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePaymentConfig{}
		selects := GetFieldsRequested(ctx, TableName("store_payment_configs", ctx))

		if IndexOf(selects, TableName("store_payment_configs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_payment_configs", ctx)+".id") == -1 {
				selects = append(selects, "store_payment_configs"+".id")
			}

			if IndexOf(selects, TableName("store_payment_configs", ctx)+".store_id") == -1 {
				selects = append(selects, TableName("store_payment_configs", ctx)+".store_id")
			}
		}

		res := db.Query().Table(TableName("store_payment_configs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "store_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StorePaymentConfig, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.StoreID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StorePaymentConfig{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePaymentConfig with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StorePaymentConfigStore"] = dataloader.NewBatchedLoader(store_payment_configsStoreBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StoreAndStorePaymentConfigIds"] = dataloader.NewBatchedLoader(store_payment_configsStoreBatchFn, dataloader.WithClearCacheOnBatch())

	store_payment_configsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePaymentConfig{}
		selects := GetFieldsRequested(ctx, TableName("store_payment_configs", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_payment_configs", ctx)+".*") == -1 && IndexOf(selects, TableName("store_payment_configs", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_payment_configs", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_payment_configs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StorePaymentConfig, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePaymentConfig with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StorePaymentConfig"] = dataloader.NewBatchedLoader(store_payment_configsBatchFn, dataloader.WithClearCacheOnBatch())

	customer_membersOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerMember{}
		selects := GetFieldsRequested(ctx, TableName("customer_members", ctx))

		if IndexOf(selects, TableName("customer_members", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_members", ctx)+".id") == -1 {
				selects = append(selects, "customer_members"+".id")
			}

			if IndexOf(selects, TableName("customer_members", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("customer_members", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("customer_members", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerMember, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerMember{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerMember with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerMemberOrganization"] = dataloader.NewBatchedLoader(customer_membersOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndCustomerMemberIds"] = dataloader.NewBatchedLoader(customer_membersOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	customer_membersBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerMember{}
		selects := GetFieldsRequested(ctx, TableName("customer_members", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("customer_members", ctx)+".*") == -1 && IndexOf(selects, TableName("customer_members", ctx)+".id") == -1 {
			selects = append(selects, TableName("customer_members", ctx)+".id")
		}

		res := db.Query().Table(TableName("customer_members", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]CustomerMember, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerMember with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["CustomerMember"] = dataloader.NewBatchedLoader(customer_membersBatchFn, dataloader.WithClearCacheOnBatch())

	customer_benefit_policiesOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerBenefitPolicy{}
		selects := GetFieldsRequested(ctx, TableName("customer_benefit_policies", ctx))

		if IndexOf(selects, TableName("customer_benefit_policies", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_benefit_policies", ctx)+".id") == -1 {
				selects = append(selects, "customer_benefit_policies"+".id")
			}

			if IndexOf(selects, TableName("customer_benefit_policies", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("customer_benefit_policies", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("customer_benefit_policies", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerBenefitPolicy, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerBenefitPolicy{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerBenefitPolicy with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerBenefitPolicyOrganization"] = dataloader.NewBatchedLoader(customer_benefit_policiesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndCustomerBenefitPolicyIds"] = dataloader.NewBatchedLoader(customer_benefit_policiesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	customer_benefit_policiesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerBenefitPolicy{}
		selects := GetFieldsRequested(ctx, TableName("customer_benefit_policies", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("customer_benefit_policies", ctx)+".*") == -1 && IndexOf(selects, TableName("customer_benefit_policies", ctx)+".id") == -1 {
			selects = append(selects, TableName("customer_benefit_policies", ctx)+".id")
		}

		res := db.Query().Table(TableName("customer_benefit_policies", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]CustomerBenefitPolicy, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerBenefitPolicy with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["CustomerBenefitPolicy"] = dataloader.NewBatchedLoader(customer_benefit_policiesBatchFn, dataloader.WithClearCacheOnBatch())

	customer_daily_point_grant_budgetsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerDailyPointGrantBudget{}
		selects := GetFieldsRequested(ctx, TableName("customer_daily_point_grant_budgets", ctx))

		if IndexOf(selects, TableName("customer_daily_point_grant_budgets", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_daily_point_grant_budgets", ctx)+".id") == -1 {
				selects = append(selects, "customer_daily_point_grant_budgets"+".id")
			}

			if IndexOf(selects, TableName("customer_daily_point_grant_budgets", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("customer_daily_point_grant_budgets", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("customer_daily_point_grant_budgets", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerDailyPointGrantBudget, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerDailyPointGrantBudget{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerDailyPointGrantBudget with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerDailyPointGrantBudgetOrganization"] = dataloader.NewBatchedLoader(customer_daily_point_grant_budgetsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndCustomerDailyPointGrantBudgetIds"] = dataloader.NewBatchedLoader(customer_daily_point_grant_budgetsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	customer_daily_point_grant_budgetsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerDailyPointGrantBudget{}
		selects := GetFieldsRequested(ctx, TableName("customer_daily_point_grant_budgets", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("customer_daily_point_grant_budgets", ctx)+".*") == -1 && IndexOf(selects, TableName("customer_daily_point_grant_budgets", ctx)+".id") == -1 {
			selects = append(selects, TableName("customer_daily_point_grant_budgets", ctx)+".id")
		}

		res := db.Query().Table(TableName("customer_daily_point_grant_budgets", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]CustomerDailyPointGrantBudget, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerDailyPointGrantBudget with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["CustomerDailyPointGrantBudget"] = dataloader.NewBatchedLoader(customer_daily_point_grant_budgetsBatchFn, dataloader.WithClearCacheOnBatch())

	customer_point_entriesMemberBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerPointEntry{}
		selects := GetFieldsRequested(ctx, TableName("customer_point_entries", ctx))

		if IndexOf(selects, TableName("customer_point_entries", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_point_entries", ctx)+".id") == -1 {
				selects = append(selects, "customer_point_entries"+".id")
			}

			if IndexOf(selects, TableName("customer_point_entries", ctx)+".member_id") == -1 {
				selects = append(selects, TableName("customer_point_entries", ctx)+".member_id")
			}
		}

		res := db.Query().Table(TableName("customer_point_entries", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "member_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerPointEntry, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.MemberID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerPointEntry{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerPointEntry with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerPointEntryMember"] = dataloader.NewBatchedLoader(customer_point_entriesMemberBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["MemberAndCustomerPointEntryIds"] = dataloader.NewBatchedLoader(customer_point_entriesMemberBatchFn, dataloader.WithClearCacheOnBatch())

	customer_point_entriesSourceOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerPointEntry{}
		selects := GetFieldsRequested(ctx, TableName("customer_point_entries", ctx))

		if IndexOf(selects, TableName("customer_point_entries", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_point_entries", ctx)+".id") == -1 {
				selects = append(selects, "customer_point_entries"+".id")
			}

			if IndexOf(selects, TableName("customer_point_entries", ctx)+".source_organization_id") == -1 {
				selects = append(selects, TableName("customer_point_entries", ctx)+".source_organization_id")
			}
		}

		res := db.Query().Table(TableName("customer_point_entries", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "source_organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerPointEntry, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.SourceOrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerPointEntry{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerPointEntry with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerPointEntrySourceOrganization"] = dataloader.NewBatchedLoader(customer_point_entriesSourceOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["SourceOrganizationAndCustomerPointEntryIds"] = dataloader.NewBatchedLoader(customer_point_entriesSourceOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	customer_point_entriesReversesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerPointEntry{}
		selects := GetFieldsRequested(ctx, TableName("customer_point_entries", ctx))

		if IndexOf(selects, TableName("customer_point_entries", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_point_entries", ctx)+".id") == -1 {
				selects = append(selects, "customer_point_entries"+".id")
			}

			if IndexOf(selects, TableName("customer_point_entries", ctx)+".reverses_id") == -1 {
				selects = append(selects, TableName("customer_point_entries", ctx)+".reverses_id")
			}
		}

		res := db.Query().Table(TableName("customer_point_entries", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "reverses_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerPointEntry, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ReversesID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerPointEntry{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerPointEntry with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerPointEntryReverses"] = dataloader.NewBatchedLoader(customer_point_entriesReversesBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ReversesAndCustomerPointEntryIds"] = dataloader.NewBatchedLoader(customer_point_entriesReversesBatchFn, dataloader.WithClearCacheOnBatch())

	customer_point_entriesReversedByBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerPointEntry{}
		selects := GetFieldsRequested(ctx, TableName("customer_point_entries", ctx))

		if IndexOf(selects, TableName("customer_point_entries", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_point_entries", ctx)+".id") == -1 {
				selects = append(selects, "customer_point_entries"+".id")
			}

			if IndexOf(selects, TableName("customer_point_entries", ctx)+".reversed_by_id") == -1 {
				selects = append(selects, TableName("customer_point_entries", ctx)+".reversed_by_id")
			}
		}

		res := db.Query().Table(TableName("customer_point_entries", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "reversed_by_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerPointEntry, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ReversedByID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerPointEntry{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerPointEntry with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerPointEntryReversedBy"] = dataloader.NewBatchedLoader(customer_point_entriesReversedByBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ReversedByAndCustomerPointEntryIds"] = dataloader.NewBatchedLoader(customer_point_entriesReversedByBatchFn, dataloader.WithClearCacheOnBatch())

	customer_point_entriesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerPointEntry{}
		selects := GetFieldsRequested(ctx, TableName("customer_point_entries", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("customer_point_entries", ctx)+".*") == -1 && IndexOf(selects, TableName("customer_point_entries", ctx)+".id") == -1 {
			selects = append(selects, TableName("customer_point_entries", ctx)+".id")
		}

		res := db.Query().Table(TableName("customer_point_entries", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]CustomerPointEntry, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerPointEntry with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["CustomerPointEntry"] = dataloader.NewBatchedLoader(customer_point_entriesBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_templatesOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponTemplate{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_templates", ctx))

		if IndexOf(selects, TableName("customer_coupon_templates", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_coupon_templates", ctx)+".id") == -1 {
				selects = append(selects, "customer_coupon_templates"+".id")
			}

			if IndexOf(selects, TableName("customer_coupon_templates", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("customer_coupon_templates", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("customer_coupon_templates", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerCouponTemplate, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerCouponTemplate{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponTemplate with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerCouponTemplateOrganization"] = dataloader.NewBatchedLoader(customer_coupon_templatesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndCustomerCouponTemplateIds"] = dataloader.NewBatchedLoader(customer_coupon_templatesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_templatesApplicableStoreBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponTemplate{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_templates", ctx))

		if IndexOf(selects, TableName("customer_coupon_templates", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_coupon_templates", ctx)+".id") == -1 {
				selects = append(selects, "customer_coupon_templates"+".id")
			}

			if IndexOf(selects, TableName("customer_coupon_templates", ctx)+".applicable_store_id") == -1 {
				selects = append(selects, TableName("customer_coupon_templates", ctx)+".applicable_store_id")
			}
		}

		res := db.Query().Table(TableName("customer_coupon_templates", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "applicable_store_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerCouponTemplate, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ApplicableStoreID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerCouponTemplate{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponTemplate with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerCouponTemplateApplicableStore"] = dataloader.NewBatchedLoader(customer_coupon_templatesApplicableStoreBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ApplicableStoreAndCustomerCouponTemplateIds"] = dataloader.NewBatchedLoader(customer_coupon_templatesApplicableStoreBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_templatesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponTemplate{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_templates", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("customer_coupon_templates", ctx)+".*") == -1 && IndexOf(selects, TableName("customer_coupon_templates", ctx)+".id") == -1 {
			selects = append(selects, TableName("customer_coupon_templates", ctx)+".id")
		}

		res := db.Query().Table(TableName("customer_coupon_templates", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]CustomerCouponTemplate, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponTemplate with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["CustomerCouponTemplate"] = dataloader.NewBatchedLoader(customer_coupon_templatesBatchFn, dataloader.WithClearCacheOnBatch())

	product_categoriesOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductCategory{}
		selects := GetFieldsRequested(ctx, TableName("product_categories", ctx))

		if IndexOf(selects, TableName("product_categories", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_categories", ctx)+".id") == -1 {
				selects = append(selects, "product_categories"+".id")
			}

			if IndexOf(selects, TableName("product_categories", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("product_categories", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("product_categories", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductCategory, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductCategory{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductCategory with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductCategoryOrganization"] = dataloader.NewBatchedLoader(product_categoriesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndProductCategoryIds"] = dataloader.NewBatchedLoader(product_categoriesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	product_categoriesParentBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductCategory{}
		selects := GetFieldsRequested(ctx, TableName("product_categories", ctx))

		if IndexOf(selects, TableName("product_categories", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_categories", ctx)+".id") == -1 {
				selects = append(selects, "product_categories"+".id")
			}

			if IndexOf(selects, TableName("product_categories", ctx)+".parent_id") == -1 {
				selects = append(selects, TableName("product_categories", ctx)+".parent_id")
			}
		}

		res := db.Query().Table(TableName("product_categories", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "parent_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductCategory, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ParentID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductCategory{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductCategory with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductCategoryParent"] = dataloader.NewBatchedLoader(product_categoriesParentBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ParentAndProductCategoryIds"] = dataloader.NewBatchedLoader(product_categoriesParentBatchFn, dataloader.WithClearCacheOnBatch())

	product_categoriesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductCategory{}
		selects := GetFieldsRequested(ctx, TableName("product_categories", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("product_categories", ctx)+".*") == -1 && IndexOf(selects, TableName("product_categories", ctx)+".id") == -1 {
			selects = append(selects, TableName("product_categories", ctx)+".id")
		}

		res := db.Query().Table(TableName("product_categories", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]ProductCategory, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductCategory with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["ProductCategory"] = dataloader.NewBatchedLoader(product_categoriesBatchFn, dataloader.WithClearCacheOnBatch())

	product_brandsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductBrand{}
		selects := GetFieldsRequested(ctx, TableName("product_brands", ctx))

		if IndexOf(selects, TableName("product_brands", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_brands", ctx)+".id") == -1 {
				selects = append(selects, "product_brands"+".id")
			}

			if IndexOf(selects, TableName("product_brands", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("product_brands", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("product_brands", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductBrand, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductBrand{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductBrand with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductBrandOrganization"] = dataloader.NewBatchedLoader(product_brandsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndProductBrandIds"] = dataloader.NewBatchedLoader(product_brandsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	product_brandsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductBrand{}
		selects := GetFieldsRequested(ctx, TableName("product_brands", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("product_brands", ctx)+".*") == -1 && IndexOf(selects, TableName("product_brands", ctx)+".id") == -1 {
			selects = append(selects, TableName("product_brands", ctx)+".id")
		}

		res := db.Query().Table(TableName("product_brands", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]ProductBrand, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductBrand with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["ProductBrand"] = dataloader.NewBatchedLoader(product_brandsBatchFn, dataloader.WithClearCacheOnBatch())

	productsBrandBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Product{}
		selects := GetFieldsRequested(ctx, TableName("products", ctx))

		if IndexOf(selects, TableName("products", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("products", ctx)+".id") == -1 {
				selects = append(selects, "products"+".id")
			}

			if IndexOf(selects, TableName("products", ctx)+".brand_id") == -1 {
				selects = append(selects, TableName("products", ctx)+".brand_id")
			}
		}

		res := db.Query().Table(TableName("products", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "brand_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Product, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.BrandID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Product{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Product with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductBrand"] = dataloader.NewBatchedLoader(productsBrandBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["BrandAndProductIds"] = dataloader.NewBatchedLoader(productsBrandBatchFn, dataloader.WithClearCacheOnBatch())

	productsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Product{}
		selects := GetFieldsRequested(ctx, TableName("products", ctx))

		if IndexOf(selects, TableName("products", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("products", ctx)+".id") == -1 {
				selects = append(selects, "products"+".id")
			}

			if IndexOf(selects, TableName("products", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("products", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("products", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Product, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Product{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Product with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductOrganization"] = dataloader.NewBatchedLoader(productsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndProductIds"] = dataloader.NewBatchedLoader(productsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	productsCategoryBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Product{}
		selects := GetFieldsRequested(ctx, TableName("products", ctx))

		if IndexOf(selects, TableName("products", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("products", ctx)+".id") == -1 {
				selects = append(selects, "products"+".id")
			}

			if IndexOf(selects, TableName("products", ctx)+".category_id") == -1 {
				selects = append(selects, TableName("products", ctx)+".category_id")
			}
		}

		res := db.Query().Table(TableName("products", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "category_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Product, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.CategoryID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Product{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Product with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductCategory"] = dataloader.NewBatchedLoader(productsCategoryBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["CategoryAndProductIds"] = dataloader.NewBatchedLoader(productsCategoryBatchFn, dataloader.WithClearCacheOnBatch())

	productsDefaultPackageTemplateBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Product{}
		selects := GetFieldsRequested(ctx, TableName("products", ctx))

		if IndexOf(selects, TableName("products", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("products", ctx)+".id") == -1 {
				selects = append(selects, "products"+".id")
			}

			if IndexOf(selects, TableName("products", ctx)+".default_package_template_id") == -1 {
				selects = append(selects, TableName("products", ctx)+".default_package_template_id")
			}
		}

		res := db.Query().Table(TableName("products", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "default_package_template_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*Product, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.DefaultPackageTemplateID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*Product{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Product with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductDefaultPackageTemplate"] = dataloader.NewBatchedLoader(productsDefaultPackageTemplateBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["DefaultPackageTemplateAndProductIds"] = dataloader.NewBatchedLoader(productsDefaultPackageTemplateBatchFn, dataloader.WithClearCacheOnBatch())

	productsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]Product{}
		selects := GetFieldsRequested(ctx, TableName("products", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("products", ctx)+".*") == -1 && IndexOf(selects, TableName("products", ctx)+".id") == -1 {
			selects = append(selects, TableName("products", ctx)+".id")
		}

		res := db.Query().Table(TableName("products", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]Product, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("Product with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["Product"] = dataloader.NewBatchedLoader(productsBatchFn, dataloader.WithClearCacheOnBatch())

	product_skusProductBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSku{}
		selects := GetFieldsRequested(ctx, TableName("product_skus", ctx))

		if IndexOf(selects, TableName("product_skus", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_skus", ctx)+".id") == -1 {
				selects = append(selects, "product_skus"+".id")
			}

			if IndexOf(selects, TableName("product_skus", ctx)+".product_id") == -1 {
				selects = append(selects, TableName("product_skus", ctx)+".product_id")
			}
		}

		res := db.Query().Table(TableName("product_skus", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "product_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductSku, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.ProductID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductSku{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSku with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductSkuProduct"] = dataloader.NewBatchedLoader(product_skusProductBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ProductAndProductSkuIds"] = dataloader.NewBatchedLoader(product_skusProductBatchFn, dataloader.WithClearCacheOnBatch())

	product_skusBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSku{}
		selects := GetFieldsRequested(ctx, TableName("product_skus", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("product_skus", ctx)+".*") == -1 && IndexOf(selects, TableName("product_skus", ctx)+".id") == -1 {
			selects = append(selects, TableName("product_skus", ctx)+".id")
		}

		res := db.Query().Table(TableName("product_skus", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]ProductSku, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSku with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["ProductSku"] = dataloader.NewBatchedLoader(product_skusBatchFn, dataloader.WithClearCacheOnBatch())

	product_packagesSkuBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductPackage{}
		selects := GetFieldsRequested(ctx, TableName("product_packages", ctx))

		if IndexOf(selects, TableName("product_packages", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_packages", ctx)+".id") == -1 {
				selects = append(selects, "product_packages"+".id")
			}

			if IndexOf(selects, TableName("product_packages", ctx)+".sku_id") == -1 {
				selects = append(selects, TableName("product_packages", ctx)+".sku_id")
			}
		}

		res := db.Query().Table(TableName("product_packages", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "sku_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductPackage, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.SkuID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductPackage{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductPackage with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductPackageSku"] = dataloader.NewBatchedLoader(product_packagesSkuBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["SkuAndProductPackageIds"] = dataloader.NewBatchedLoader(product_packagesSkuBatchFn, dataloader.WithClearCacheOnBatch())

	product_packagesTemplateBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductPackage{}
		selects := GetFieldsRequested(ctx, TableName("product_packages", ctx))

		if IndexOf(selects, TableName("product_packages", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_packages", ctx)+".id") == -1 {
				selects = append(selects, "product_packages"+".id")
			}

			if IndexOf(selects, TableName("product_packages", ctx)+".template_id") == -1 {
				selects = append(selects, TableName("product_packages", ctx)+".template_id")
			}
		}

		res := db.Query().Table(TableName("product_packages", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "template_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductPackage, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.TemplateID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductPackage{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductPackage with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductPackageTemplate"] = dataloader.NewBatchedLoader(product_packagesTemplateBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["TemplateAndProductPackageIds"] = dataloader.NewBatchedLoader(product_packagesTemplateBatchFn, dataloader.WithClearCacheOnBatch())

	product_packagesContainsPackageBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductPackage{}
		selects := GetFieldsRequested(ctx, TableName("product_packages", ctx))

		if IndexOf(selects, TableName("product_packages", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_packages", ctx)+".id") == -1 {
				selects = append(selects, "product_packages"+".id")
			}

			if IndexOf(selects, TableName("product_packages", ctx)+".contains_package_id") == -1 {
				selects = append(selects, TableName("product_packages", ctx)+".contains_package_id")
			}
		}

		res := db.Query().Table(TableName("product_packages", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "contains_package_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductPackage, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ContainsPackageID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductPackage{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductPackage with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductPackageContainsPackage"] = dataloader.NewBatchedLoader(product_packagesContainsPackageBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ContainsPackageAndProductPackageIds"] = dataloader.NewBatchedLoader(product_packagesContainsPackageBatchFn, dataloader.WithClearCacheOnBatch())

	product_packagesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductPackage{}
		selects := GetFieldsRequested(ctx, TableName("product_packages", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("product_packages", ctx)+".*") == -1 && IndexOf(selects, TableName("product_packages", ctx)+".id") == -1 {
			selects = append(selects, TableName("product_packages", ctx)+".id")
		}

		res := db.Query().Table(TableName("product_packages", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]ProductPackage, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductPackage with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["ProductPackage"] = dataloader.NewBatchedLoader(product_packagesBatchFn, dataloader.WithClearCacheOnBatch())

	specification_definitionsOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]SpecificationDefinition{}
		selects := GetFieldsRequested(ctx, TableName("specification_definitions", ctx))

		if IndexOf(selects, TableName("specification_definitions", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("specification_definitions", ctx)+".id") == -1 {
				selects = append(selects, "specification_definitions"+".id")
			}

			if IndexOf(selects, TableName("specification_definitions", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("specification_definitions", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("specification_definitions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*SpecificationDefinition, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*SpecificationDefinition{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("SpecificationDefinition with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["SpecificationDefinitionOrganization"] = dataloader.NewBatchedLoader(specification_definitionsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndSpecificationDefinitionIds"] = dataloader.NewBatchedLoader(specification_definitionsOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	specification_definitionsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]SpecificationDefinition{}
		selects := GetFieldsRequested(ctx, TableName("specification_definitions", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("specification_definitions", ctx)+".*") == -1 && IndexOf(selects, TableName("specification_definitions", ctx)+".id") == -1 {
			selects = append(selects, TableName("specification_definitions", ctx)+".id")
		}

		res := db.Query().Table(TableName("specification_definitions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]SpecificationDefinition, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("SpecificationDefinition with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["SpecificationDefinition"] = dataloader.NewBatchedLoader(specification_definitionsBatchFn, dataloader.WithClearCacheOnBatch())

	specification_valuesSpecificationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]SpecificationValue{}
		selects := GetFieldsRequested(ctx, TableName("specification_values", ctx))

		if IndexOf(selects, TableName("specification_values", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("specification_values", ctx)+".id") == -1 {
				selects = append(selects, "specification_values"+".id")
			}

			if IndexOf(selects, TableName("specification_values", ctx)+".specification_id") == -1 {
				selects = append(selects, TableName("specification_values", ctx)+".specification_id")
			}
		}

		res := db.Query().Table(TableName("specification_values", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "specification_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*SpecificationValue, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.SpecificationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*SpecificationValue{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("SpecificationValue with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["SpecificationValueSpecification"] = dataloader.NewBatchedLoader(specification_valuesSpecificationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["SpecificationAndSpecificationValueIds"] = dataloader.NewBatchedLoader(specification_valuesSpecificationBatchFn, dataloader.WithClearCacheOnBatch())

	specification_valuesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]SpecificationValue{}
		selects := GetFieldsRequested(ctx, TableName("specification_values", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("specification_values", ctx)+".*") == -1 && IndexOf(selects, TableName("specification_values", ctx)+".id") == -1 {
			selects = append(selects, TableName("specification_values", ctx)+".id")
		}

		res := db.Query().Table(TableName("specification_values", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]SpecificationValue, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("SpecificationValue with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["SpecificationValue"] = dataloader.NewBatchedLoader(specification_valuesBatchFn, dataloader.WithClearCacheOnBatch())

	product_specification_choicesProductBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSpecificationChoice{}
		selects := GetFieldsRequested(ctx, TableName("product_specification_choices", ctx))

		if IndexOf(selects, TableName("product_specification_choices", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_specification_choices", ctx)+".id") == -1 {
				selects = append(selects, "product_specification_choices"+".id")
			}

			if IndexOf(selects, TableName("product_specification_choices", ctx)+".product_id") == -1 {
				selects = append(selects, TableName("product_specification_choices", ctx)+".product_id")
			}
		}

		res := db.Query().Table(TableName("product_specification_choices", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "product_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductSpecificationChoice, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.ProductID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductSpecificationChoice{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSpecificationChoice with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductSpecificationChoiceProduct"] = dataloader.NewBatchedLoader(product_specification_choicesProductBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ProductAndProductSpecificationChoiceIds"] = dataloader.NewBatchedLoader(product_specification_choicesProductBatchFn, dataloader.WithClearCacheOnBatch())

	product_specification_choicesValueBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSpecificationChoice{}
		selects := GetFieldsRequested(ctx, TableName("product_specification_choices", ctx))

		if IndexOf(selects, TableName("product_specification_choices", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_specification_choices", ctx)+".id") == -1 {
				selects = append(selects, "product_specification_choices"+".id")
			}

			if IndexOf(selects, TableName("product_specification_choices", ctx)+".value_id") == -1 {
				selects = append(selects, TableName("product_specification_choices", ctx)+".value_id")
			}
		}

		res := db.Query().Table(TableName("product_specification_choices", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "value_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductSpecificationChoice, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.ValueID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductSpecificationChoice{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSpecificationChoice with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductSpecificationChoiceValue"] = dataloader.NewBatchedLoader(product_specification_choicesValueBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ValueAndProductSpecificationChoiceIds"] = dataloader.NewBatchedLoader(product_specification_choicesValueBatchFn, dataloader.WithClearCacheOnBatch())

	product_specification_choicesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSpecificationChoice{}
		selects := GetFieldsRequested(ctx, TableName("product_specification_choices", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("product_specification_choices", ctx)+".*") == -1 && IndexOf(selects, TableName("product_specification_choices", ctx)+".id") == -1 {
			selects = append(selects, TableName("product_specification_choices", ctx)+".id")
		}

		res := db.Query().Table(TableName("product_specification_choices", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]ProductSpecificationChoice, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSpecificationChoice with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["ProductSpecificationChoice"] = dataloader.NewBatchedLoader(product_specification_choicesBatchFn, dataloader.WithClearCacheOnBatch())

	product_sku_specification_valuesSkuBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSkuSpecificationValue{}
		selects := GetFieldsRequested(ctx, TableName("product_sku_specification_values", ctx))

		if IndexOf(selects, TableName("product_sku_specification_values", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_sku_specification_values", ctx)+".id") == -1 {
				selects = append(selects, "product_sku_specification_values"+".id")
			}

			if IndexOf(selects, TableName("product_sku_specification_values", ctx)+".sku_id") == -1 {
				selects = append(selects, TableName("product_sku_specification_values", ctx)+".sku_id")
			}
		}

		res := db.Query().Table(TableName("product_sku_specification_values", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "sku_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductSkuSpecificationValue, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.SkuID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductSkuSpecificationValue{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSkuSpecificationValue with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductSkuSpecificationValueSku"] = dataloader.NewBatchedLoader(product_sku_specification_valuesSkuBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["SkuAndProductSkuSpecificationValueIds"] = dataloader.NewBatchedLoader(product_sku_specification_valuesSkuBatchFn, dataloader.WithClearCacheOnBatch())

	product_sku_specification_valuesValueBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSkuSpecificationValue{}
		selects := GetFieldsRequested(ctx, TableName("product_sku_specification_values", ctx))

		if IndexOf(selects, TableName("product_sku_specification_values", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_sku_specification_values", ctx)+".id") == -1 {
				selects = append(selects, "product_sku_specification_values"+".id")
			}

			if IndexOf(selects, TableName("product_sku_specification_values", ctx)+".value_id") == -1 {
				selects = append(selects, TableName("product_sku_specification_values", ctx)+".value_id")
			}
		}

		res := db.Query().Table(TableName("product_sku_specification_values", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "value_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductSkuSpecificationValue, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.ValueID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductSkuSpecificationValue{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSkuSpecificationValue with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductSkuSpecificationValueValue"] = dataloader.NewBatchedLoader(product_sku_specification_valuesValueBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ValueAndProductSkuSpecificationValueIds"] = dataloader.NewBatchedLoader(product_sku_specification_valuesValueBatchFn, dataloader.WithClearCacheOnBatch())

	product_sku_specification_valuesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductSkuSpecificationValue{}
		selects := GetFieldsRequested(ctx, TableName("product_sku_specification_values", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("product_sku_specification_values", ctx)+".*") == -1 && IndexOf(selects, TableName("product_sku_specification_values", ctx)+".id") == -1 {
			selects = append(selects, TableName("product_sku_specification_values", ctx)+".id")
		}

		res := db.Query().Table(TableName("product_sku_specification_values", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]ProductSkuSpecificationValue, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductSkuSpecificationValue with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["ProductSkuSpecificationValue"] = dataloader.NewBatchedLoader(product_sku_specification_valuesBatchFn, dataloader.WithClearCacheOnBatch())

	product_package_templatesOrganizationBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductPackageTemplate{}
		selects := GetFieldsRequested(ctx, TableName("product_package_templates", ctx))

		if IndexOf(selects, TableName("product_package_templates", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_package_templates", ctx)+".id") == -1 {
				selects = append(selects, "product_package_templates"+".id")
			}

			if IndexOf(selects, TableName("product_package_templates", ctx)+".organization_id") == -1 {
				selects = append(selects, TableName("product_package_templates", ctx)+".organization_id")
			}
		}

		res := db.Query().Table(TableName("product_package_templates", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "organization_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductPackageTemplate, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OrganizationID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductPackageTemplate{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductPackageTemplate with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductPackageTemplateOrganization"] = dataloader.NewBatchedLoader(product_package_templatesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OrganizationAndProductPackageTemplateIds"] = dataloader.NewBatchedLoader(product_package_templatesOrganizationBatchFn, dataloader.WithClearCacheOnBatch())

	product_package_templatesContainsPackageBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductPackageTemplate{}
		selects := GetFieldsRequested(ctx, TableName("product_package_templates", ctx))

		if IndexOf(selects, TableName("product_package_templates", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("product_package_templates", ctx)+".id") == -1 {
				selects = append(selects, "product_package_templates"+".id")
			}

			if IndexOf(selects, TableName("product_package_templates", ctx)+".contains_package_id") == -1 {
				selects = append(selects, TableName("product_package_templates", ctx)+".contains_package_id")
			}
		}

		res := db.Query().Table(TableName("product_package_templates", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "contains_package_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*ProductPackageTemplate, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.ContainsPackageID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*ProductPackageTemplate{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductPackageTemplate with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["ProductPackageTemplateContainsPackage"] = dataloader.NewBatchedLoader(product_package_templatesContainsPackageBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ContainsPackageAndProductPackageTemplateIds"] = dataloader.NewBatchedLoader(product_package_templatesContainsPackageBatchFn, dataloader.WithClearCacheOnBatch())

	product_package_templatesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]ProductPackageTemplate{}
		selects := GetFieldsRequested(ctx, TableName("product_package_templates", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("product_package_templates", ctx)+".*") == -1 && IndexOf(selects, TableName("product_package_templates", ctx)+".id") == -1 {
			selects = append(selects, TableName("product_package_templates", ctx)+".id")
		}

		res := db.Query().Table(TableName("product_package_templates", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]ProductPackageTemplate, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("ProductPackageTemplate with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["ProductPackageTemplate"] = dataloader.NewBatchedLoader(product_package_templatesBatchFn, dataloader.WithClearCacheOnBatch())

	store_listingsStoreBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreListing{}
		selects := GetFieldsRequested(ctx, TableName("store_listings", ctx))

		if IndexOf(selects, TableName("store_listings", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_listings", ctx)+".id") == -1 {
				selects = append(selects, "store_listings"+".id")
			}

			if IndexOf(selects, TableName("store_listings", ctx)+".store_id") == -1 {
				selects = append(selects, TableName("store_listings", ctx)+".store_id")
			}
		}

		res := db.Query().Table(TableName("store_listings", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "store_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreListing, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.StoreID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreListing{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreListing with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreListingStore"] = dataloader.NewBatchedLoader(store_listingsStoreBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StoreAndStoreListingIds"] = dataloader.NewBatchedLoader(store_listingsStoreBatchFn, dataloader.WithClearCacheOnBatch())

	store_listingsSkuBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreListing{}
		selects := GetFieldsRequested(ctx, TableName("store_listings", ctx))

		if IndexOf(selects, TableName("store_listings", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_listings", ctx)+".id") == -1 {
				selects = append(selects, "store_listings"+".id")
			}

			if IndexOf(selects, TableName("store_listings", ctx)+".sku_id") == -1 {
				selects = append(selects, TableName("store_listings", ctx)+".sku_id")
			}
		}

		res := db.Query().Table(TableName("store_listings", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "sku_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreListing, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.SkuID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreListing{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreListing with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreListingSku"] = dataloader.NewBatchedLoader(store_listingsSkuBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["SkuAndStoreListingIds"] = dataloader.NewBatchedLoader(store_listingsSkuBatchFn, dataloader.WithClearCacheOnBatch())

	store_listingsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreListing{}
		selects := GetFieldsRequested(ctx, TableName("store_listings", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_listings", ctx)+".*") == -1 && IndexOf(selects, TableName("store_listings", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_listings", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_listings", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StoreListing, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreListing with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StoreListing"] = dataloader.NewBatchedLoader(store_listingsBatchFn, dataloader.WithClearCacheOnBatch())

	store_package_offersListingBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePackageOffer{}
		selects := GetFieldsRequested(ctx, TableName("store_package_offers", ctx))

		if IndexOf(selects, TableName("store_package_offers", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_package_offers", ctx)+".id") == -1 {
				selects = append(selects, "store_package_offers"+".id")
			}

			if IndexOf(selects, TableName("store_package_offers", ctx)+".listing_id") == -1 {
				selects = append(selects, TableName("store_package_offers", ctx)+".listing_id")
			}
		}

		res := db.Query().Table(TableName("store_package_offers", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "listing_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StorePackageOffer, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.ListingID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StorePackageOffer{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePackageOffer with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StorePackageOfferListing"] = dataloader.NewBatchedLoader(store_package_offersListingBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ListingAndStorePackageOfferIds"] = dataloader.NewBatchedLoader(store_package_offersListingBatchFn, dataloader.WithClearCacheOnBatch())

	store_package_offersPackageBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePackageOffer{}
		selects := GetFieldsRequested(ctx, TableName("store_package_offers", ctx))

		if IndexOf(selects, TableName("store_package_offers", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_package_offers", ctx)+".id") == -1 {
				selects = append(selects, "store_package_offers"+".id")
			}

			if IndexOf(selects, TableName("store_package_offers", ctx)+".package_id") == -1 {
				selects = append(selects, TableName("store_package_offers", ctx)+".package_id")
			}
		}

		res := db.Query().Table(TableName("store_package_offers", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "package_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StorePackageOffer, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.PackageID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StorePackageOffer{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePackageOffer with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StorePackageOfferPackage"] = dataloader.NewBatchedLoader(store_package_offersPackageBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["PackageAndStorePackageOfferIds"] = dataloader.NewBatchedLoader(store_package_offersPackageBatchFn, dataloader.WithClearCacheOnBatch())

	store_package_offersBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePackageOffer{}
		selects := GetFieldsRequested(ctx, TableName("store_package_offers", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_package_offers", ctx)+".*") == -1 && IndexOf(selects, TableName("store_package_offers", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_package_offers", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_package_offers", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StorePackageOffer, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePackageOffer with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StorePackageOffer"] = dataloader.NewBatchedLoader(store_package_offersBatchFn, dataloader.WithClearCacheOnBatch())

	store_price_revisionsOfferBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePriceRevision{}
		selects := GetFieldsRequested(ctx, TableName("store_price_revisions", ctx))

		if IndexOf(selects, TableName("store_price_revisions", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_price_revisions", ctx)+".id") == -1 {
				selects = append(selects, "store_price_revisions"+".id")
			}

			if IndexOf(selects, TableName("store_price_revisions", ctx)+".offer_id") == -1 {
				selects = append(selects, TableName("store_price_revisions", ctx)+".offer_id")
			}
		}

		res := db.Query().Table(TableName("store_price_revisions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "offer_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StorePriceRevision, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OfferID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StorePriceRevision{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePriceRevision with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StorePriceRevisionOffer"] = dataloader.NewBatchedLoader(store_price_revisionsOfferBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OfferAndStorePriceRevisionIds"] = dataloader.NewBatchedLoader(store_price_revisionsOfferBatchFn, dataloader.WithClearCacheOnBatch())

	store_price_revisionsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePriceRevision{}
		selects := GetFieldsRequested(ctx, TableName("store_price_revisions", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_price_revisions", ctx)+".*") == -1 && IndexOf(selects, TableName("store_price_revisions", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_price_revisions", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_price_revisions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StorePriceRevision, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePriceRevision with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StorePriceRevision"] = dataloader.NewBatchedLoader(store_price_revisionsBatchFn, dataloader.WithClearCacheOnBatch())

	store_inventory_batchesListingBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreInventoryBatch{}
		selects := GetFieldsRequested(ctx, TableName("store_inventory_batches", ctx))

		if IndexOf(selects, TableName("store_inventory_batches", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_inventory_batches", ctx)+".id") == -1 {
				selects = append(selects, "store_inventory_batches"+".id")
			}

			if IndexOf(selects, TableName("store_inventory_batches", ctx)+".listing_id") == -1 {
				selects = append(selects, TableName("store_inventory_batches", ctx)+".listing_id")
			}
		}

		res := db.Query().Table(TableName("store_inventory_batches", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "listing_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreInventoryBatch, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.ListingID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreInventoryBatch{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreInventoryBatch with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreInventoryBatchListing"] = dataloader.NewBatchedLoader(store_inventory_batchesListingBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["ListingAndStoreInventoryBatchIds"] = dataloader.NewBatchedLoader(store_inventory_batchesListingBatchFn, dataloader.WithClearCacheOnBatch())

	store_inventory_batchesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreInventoryBatch{}
		selects := GetFieldsRequested(ctx, TableName("store_inventory_batches", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_inventory_batches", ctx)+".*") == -1 && IndexOf(selects, TableName("store_inventory_batches", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_inventory_batches", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_inventory_batches", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StoreInventoryBatch, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreInventoryBatch with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StoreInventoryBatch"] = dataloader.NewBatchedLoader(store_inventory_batchesBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_balancesBatchBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockBalance{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_balances", ctx))

		if IndexOf(selects, TableName("store_stock_balances", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stock_balances", ctx)+".id") == -1 {
				selects = append(selects, "store_stock_balances"+".id")
			}

			if IndexOf(selects, TableName("store_stock_balances", ctx)+".batch_id") == -1 {
				selects = append(selects, TableName("store_stock_balances", ctx)+".batch_id")
			}
		}

		res := db.Query().Table(TableName("store_stock_balances", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "batch_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStockBalance, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.BatchID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStockBalance{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockBalance with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStockBalanceBatch"] = dataloader.NewBatchedLoader(store_stock_balancesBatchBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["BatchAndStoreStockBalanceIds"] = dataloader.NewBatchedLoader(store_stock_balancesBatchBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_balancesPackageBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockBalance{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_balances", ctx))

		if IndexOf(selects, TableName("store_stock_balances", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stock_balances", ctx)+".id") == -1 {
				selects = append(selects, "store_stock_balances"+".id")
			}

			if IndexOf(selects, TableName("store_stock_balances", ctx)+".package_id") == -1 {
				selects = append(selects, TableName("store_stock_balances", ctx)+".package_id")
			}
		}

		res := db.Query().Table(TableName("store_stock_balances", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "package_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStockBalance, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.PackageID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStockBalance{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockBalance with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStockBalancePackage"] = dataloader.NewBatchedLoader(store_stock_balancesPackageBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["PackageAndStoreStockBalanceIds"] = dataloader.NewBatchedLoader(store_stock_balancesPackageBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_balancesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockBalance{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_balances", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_stock_balances", ctx)+".*") == -1 && IndexOf(selects, TableName("store_stock_balances", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_stock_balances", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_stock_balances", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StoreStockBalance, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockBalance with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StoreStockBalance"] = dataloader.NewBatchedLoader(store_stock_balancesBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktakesStoreBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktake{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktakes", ctx))

		if IndexOf(selects, TableName("store_stocktakes", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stocktakes", ctx)+".id") == -1 {
				selects = append(selects, "store_stocktakes"+".id")
			}

			if IndexOf(selects, TableName("store_stocktakes", ctx)+".store_id") == -1 {
				selects = append(selects, TableName("store_stocktakes", ctx)+".store_id")
			}
		}

		res := db.Query().Table(TableName("store_stocktakes", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "store_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStocktake, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.StoreID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStocktake{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktake with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStocktakeStore"] = dataloader.NewBatchedLoader(store_stocktakesStoreBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StoreAndStoreStocktakeIds"] = dataloader.NewBatchedLoader(store_stocktakesStoreBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktakesInitiatedByAccountBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktake{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktakes", ctx))

		if IndexOf(selects, TableName("store_stocktakes", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stocktakes", ctx)+".id") == -1 {
				selects = append(selects, "store_stocktakes"+".id")
			}

			if IndexOf(selects, TableName("store_stocktakes", ctx)+".initiated_by_account_id") == -1 {
				selects = append(selects, TableName("store_stocktakes", ctx)+".initiated_by_account_id")
			}
		}

		res := db.Query().Table(TableName("store_stocktakes", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "initiated_by_account_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStocktake, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.InitiatedByAccountID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStocktake{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktake with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStocktakeInitiatedByAccount"] = dataloader.NewBatchedLoader(store_stocktakesInitiatedByAccountBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["InitiatedByAccountAndStoreStocktakeIds"] = dataloader.NewBatchedLoader(store_stocktakesInitiatedByAccountBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktakesPostedByBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktake{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktakes", ctx))

		if IndexOf(selects, TableName("store_stocktakes", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stocktakes", ctx)+".id") == -1 {
				selects = append(selects, "store_stocktakes"+".id")
			}

			if IndexOf(selects, TableName("store_stocktakes", ctx)+".posted_by_id") == -1 {
				selects = append(selects, TableName("store_stocktakes", ctx)+".posted_by_id")
			}
		}

		res := db.Query().Table(TableName("store_stocktakes", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "posted_by_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStocktake, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.PostedByID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStocktake{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktake with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStocktakePostedBy"] = dataloader.NewBatchedLoader(store_stocktakesPostedByBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["PostedByAndStoreStocktakeIds"] = dataloader.NewBatchedLoader(store_stocktakesPostedByBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktakesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktake{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktakes", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_stocktakes", ctx)+".*") == -1 && IndexOf(selects, TableName("store_stocktakes", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_stocktakes", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_stocktakes", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StoreStocktake, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktake with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StoreStocktake"] = dataloader.NewBatchedLoader(store_stocktakesBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktake_linesStocktakeBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktakeLine{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktake_lines", ctx))

		if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".id") == -1 {
				selects = append(selects, "store_stocktake_lines"+".id")
			}

			if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".stocktake_id") == -1 {
				selects = append(selects, TableName("store_stocktake_lines", ctx)+".stocktake_id")
			}
		}

		res := db.Query().Table(TableName("store_stocktake_lines", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "stocktake_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStocktakeLine, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.StocktakeID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStocktakeLine{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktakeLine with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStocktakeLineStocktake"] = dataloader.NewBatchedLoader(store_stocktake_linesStocktakeBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StocktakeAndStoreStocktakeLineIds"] = dataloader.NewBatchedLoader(store_stocktake_linesStocktakeBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktake_linesBatchBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktakeLine{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktake_lines", ctx))

		if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".id") == -1 {
				selects = append(selects, "store_stocktake_lines"+".id")
			}

			if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".batch_id") == -1 {
				selects = append(selects, TableName("store_stocktake_lines", ctx)+".batch_id")
			}
		}

		res := db.Query().Table(TableName("store_stocktake_lines", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "batch_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStocktakeLine, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.BatchID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStocktakeLine{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktakeLine with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStocktakeLineBatch"] = dataloader.NewBatchedLoader(store_stocktake_linesBatchBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["BatchAndStoreStocktakeLineIds"] = dataloader.NewBatchedLoader(store_stocktake_linesBatchBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktake_linesPackageBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktakeLine{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktake_lines", ctx))

		if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".id") == -1 {
				selects = append(selects, "store_stocktake_lines"+".id")
			}

			if IndexOf(selects, TableName("store_stocktake_lines", ctx)+".package_id") == -1 {
				selects = append(selects, TableName("store_stocktake_lines", ctx)+".package_id")
			}
		}

		res := db.Query().Table(TableName("store_stocktake_lines", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "package_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStocktakeLine, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.PackageID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStocktakeLine{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktakeLine with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStocktakeLinePackage"] = dataloader.NewBatchedLoader(store_stocktake_linesPackageBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["PackageAndStoreStocktakeLineIds"] = dataloader.NewBatchedLoader(store_stocktake_linesPackageBatchFn, dataloader.WithClearCacheOnBatch())

	store_stocktake_linesBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStocktakeLine{}
		selects := GetFieldsRequested(ctx, TableName("store_stocktake_lines", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_stocktake_lines", ctx)+".*") == -1 && IndexOf(selects, TableName("store_stocktake_lines", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_stocktake_lines", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_stocktake_lines", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StoreStocktakeLine, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStocktakeLine with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StoreStocktakeLine"] = dataloader.NewBatchedLoader(store_stocktake_linesBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_movementsStoreBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockMovement{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_movements", ctx))

		if IndexOf(selects, TableName("store_stock_movements", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stock_movements", ctx)+".id") == -1 {
				selects = append(selects, "store_stock_movements"+".id")
			}

			if IndexOf(selects, TableName("store_stock_movements", ctx)+".store_id") == -1 {
				selects = append(selects, TableName("store_stock_movements", ctx)+".store_id")
			}
		}

		res := db.Query().Table(TableName("store_stock_movements", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "store_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStockMovement, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.StoreID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStockMovement{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockMovement with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStockMovementStore"] = dataloader.NewBatchedLoader(store_stock_movementsStoreBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StoreAndStoreStockMovementIds"] = dataloader.NewBatchedLoader(store_stock_movementsStoreBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_movementsBatchBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockMovement{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_movements", ctx))

		if IndexOf(selects, TableName("store_stock_movements", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stock_movements", ctx)+".id") == -1 {
				selects = append(selects, "store_stock_movements"+".id")
			}

			if IndexOf(selects, TableName("store_stock_movements", ctx)+".batch_id") == -1 {
				selects = append(selects, TableName("store_stock_movements", ctx)+".batch_id")
			}
		}

		res := db.Query().Table(TableName("store_stock_movements", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "batch_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStockMovement, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.BatchID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStockMovement{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockMovement with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStockMovementBatch"] = dataloader.NewBatchedLoader(store_stock_movementsBatchBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["BatchAndStoreStockMovementIds"] = dataloader.NewBatchedLoader(store_stock_movementsBatchBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_movementsSourcePackageBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockMovement{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_movements", ctx))

		if IndexOf(selects, TableName("store_stock_movements", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stock_movements", ctx)+".id") == -1 {
				selects = append(selects, "store_stock_movements"+".id")
			}

			if IndexOf(selects, TableName("store_stock_movements", ctx)+".source_package_id") == -1 {
				selects = append(selects, TableName("store_stock_movements", ctx)+".source_package_id")
			}
		}

		res := db.Query().Table(TableName("store_stock_movements", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "source_package_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStockMovement, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.SourcePackageID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStockMovement{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockMovement with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStockMovementSourcePackage"] = dataloader.NewBatchedLoader(store_stock_movementsSourcePackageBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["SourcePackageAndStoreStockMovementIds"] = dataloader.NewBatchedLoader(store_stock_movementsSourcePackageBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_movementsTargetPackageBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockMovement{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_movements", ctx))

		if IndexOf(selects, TableName("store_stock_movements", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stock_movements", ctx)+".id") == -1 {
				selects = append(selects, "store_stock_movements"+".id")
			}

			if IndexOf(selects, TableName("store_stock_movements", ctx)+".target_package_id") == -1 {
				selects = append(selects, TableName("store_stock_movements", ctx)+".target_package_id")
			}
		}

		res := db.Query().Table(TableName("store_stock_movements", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "target_package_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStockMovement, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.TargetPackageID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStockMovement{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockMovement with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStockMovementTargetPackage"] = dataloader.NewBatchedLoader(store_stock_movementsTargetPackageBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["TargetPackageAndStoreStockMovementIds"] = dataloader.NewBatchedLoader(store_stock_movementsTargetPackageBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_movementsStocktakeLineBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockMovement{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_movements", ctx))

		if IndexOf(selects, TableName("store_stock_movements", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_stock_movements", ctx)+".id") == -1 {
				selects = append(selects, "store_stock_movements"+".id")
			}

			if IndexOf(selects, TableName("store_stock_movements", ctx)+".stocktake_line_id") == -1 {
				selects = append(selects, TableName("store_stock_movements", ctx)+".stocktake_line_id")
			}
		}

		res := db.Query().Table(TableName("store_stock_movements", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "stocktake_line_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StoreStockMovement, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.StocktakeLineID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StoreStockMovement{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockMovement with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StoreStockMovementStocktakeLine"] = dataloader.NewBatchedLoader(store_stock_movementsStocktakeLineBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StocktakeLineAndStoreStockMovementIds"] = dataloader.NewBatchedLoader(store_stock_movementsStocktakeLineBatchFn, dataloader.WithClearCacheOnBatch())

	store_stock_movementsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StoreStockMovement{}
		selects := GetFieldsRequested(ctx, TableName("store_stock_movements", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_stock_movements", ctx)+".*") == -1 && IndexOf(selects, TableName("store_stock_movements", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_stock_movements", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_stock_movements", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StoreStockMovement, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StoreStockMovement with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StoreStockMovement"] = dataloader.NewBatchedLoader(store_stock_movementsBatchFn, dataloader.WithClearCacheOnBatch())

	store_promotionsStoreBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePromotion{}
		selects := GetFieldsRequested(ctx, TableName("store_promotions", ctx))

		if IndexOf(selects, TableName("store_promotions", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_promotions", ctx)+".id") == -1 {
				selects = append(selects, "store_promotions"+".id")
			}

			if IndexOf(selects, TableName("store_promotions", ctx)+".store_id") == -1 {
				selects = append(selects, TableName("store_promotions", ctx)+".store_id")
			}
		}

		res := db.Query().Table(TableName("store_promotions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "store_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StorePromotion, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.StoreID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StorePromotion{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePromotion with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StorePromotionStore"] = dataloader.NewBatchedLoader(store_promotionsStoreBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["StoreAndStorePromotionIds"] = dataloader.NewBatchedLoader(store_promotionsStoreBatchFn, dataloader.WithClearCacheOnBatch())

	store_promotionsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePromotion{}
		selects := GetFieldsRequested(ctx, TableName("store_promotions", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_promotions", ctx)+".*") == -1 && IndexOf(selects, TableName("store_promotions", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_promotions", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_promotions", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StorePromotion, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePromotion with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StorePromotion"] = dataloader.NewBatchedLoader(store_promotionsBatchFn, dataloader.WithClearCacheOnBatch())

	store_promotion_targetsPromotionBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePromotionTarget{}
		selects := GetFieldsRequested(ctx, TableName("store_promotion_targets", ctx))

		if IndexOf(selects, TableName("store_promotion_targets", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_promotion_targets", ctx)+".id") == -1 {
				selects = append(selects, "store_promotion_targets"+".id")
			}

			if IndexOf(selects, TableName("store_promotion_targets", ctx)+".promotion_id") == -1 {
				selects = append(selects, TableName("store_promotion_targets", ctx)+".promotion_id")
			}
		}

		res := db.Query().Table(TableName("store_promotion_targets", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "promotion_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StorePromotionTarget, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.PromotionID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StorePromotionTarget{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePromotionTarget with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StorePromotionTargetPromotion"] = dataloader.NewBatchedLoader(store_promotion_targetsPromotionBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["PromotionAndStorePromotionTargetIds"] = dataloader.NewBatchedLoader(store_promotion_targetsPromotionBatchFn, dataloader.WithClearCacheOnBatch())

	store_promotion_targetsOfferBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePromotionTarget{}
		selects := GetFieldsRequested(ctx, TableName("store_promotion_targets", ctx))

		if IndexOf(selects, TableName("store_promotion_targets", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("store_promotion_targets", ctx)+".id") == -1 {
				selects = append(selects, "store_promotion_targets"+".id")
			}

			if IndexOf(selects, TableName("store_promotion_targets", ctx)+".offer_id") == -1 {
				selects = append(selects, TableName("store_promotion_targets", ctx)+".offer_id")
			}
		}

		res := db.Query().Table(TableName("store_promotion_targets", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "offer_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*StorePromotionTarget, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.OfferID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*StorePromotionTarget{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePromotionTarget with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["StorePromotionTargetOffer"] = dataloader.NewBatchedLoader(store_promotion_targetsOfferBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["OfferAndStorePromotionTargetIds"] = dataloader.NewBatchedLoader(store_promotion_targetsOfferBatchFn, dataloader.WithClearCacheOnBatch())

	store_promotion_targetsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]StorePromotionTarget{}
		selects := GetFieldsRequested(ctx, TableName("store_promotion_targets", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("store_promotion_targets", ctx)+".*") == -1 && IndexOf(selects, TableName("store_promotion_targets", ctx)+".id") == -1 {
			selects = append(selects, TableName("store_promotion_targets", ctx)+".id")
		}

		res := db.Query().Table(TableName("store_promotion_targets", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]StorePromotionTarget, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("StorePromotionTarget with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["StorePromotionTarget"] = dataloader.NewBatchedLoader(store_promotion_targetsBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_grantsMemberBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponGrant{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_grants", ctx))

		if IndexOf(selects, TableName("customer_coupon_grants", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_coupon_grants", ctx)+".id") == -1 {
				selects = append(selects, "customer_coupon_grants"+".id")
			}

			if IndexOf(selects, TableName("customer_coupon_grants", ctx)+".member_id") == -1 {
				selects = append(selects, TableName("customer_coupon_grants", ctx)+".member_id")
			}
		}

		res := db.Query().Table(TableName("customer_coupon_grants", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "member_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerCouponGrant, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.MemberID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerCouponGrant{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponGrant with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerCouponGrantMember"] = dataloader.NewBatchedLoader(customer_coupon_grantsMemberBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["MemberAndCustomerCouponGrantIds"] = dataloader.NewBatchedLoader(customer_coupon_grantsMemberBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_grantsTemplateBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponGrant{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_grants", ctx))

		if IndexOf(selects, TableName("customer_coupon_grants", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_coupon_grants", ctx)+".id") == -1 {
				selects = append(selects, "customer_coupon_grants"+".id")
			}

			if IndexOf(selects, TableName("customer_coupon_grants", ctx)+".template_id") == -1 {
				selects = append(selects, TableName("customer_coupon_grants", ctx)+".template_id")
			}
		}

		res := db.Query().Table(TableName("customer_coupon_grants", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "template_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerCouponGrant, len(keys))
		for _, v := range *items {
			item := v

			mapKey := item.TemplateID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerCouponGrant{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponGrant with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerCouponGrantTemplate"] = dataloader.NewBatchedLoader(customer_coupon_grantsTemplateBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["TemplateAndCustomerCouponGrantIds"] = dataloader.NewBatchedLoader(customer_coupon_grantsTemplateBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_grantsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponGrant{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_grants", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("customer_coupon_grants", ctx)+".*") == -1 && IndexOf(selects, TableName("customer_coupon_grants", ctx)+".id") == -1 {
			selects = append(selects, TableName("customer_coupon_grants", ctx)+".id")
		}

		res := db.Query().Table(TableName("customer_coupon_grants", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]CustomerCouponGrant, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponGrant with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["CustomerCouponGrant"] = dataloader.NewBatchedLoader(customer_coupon_grantsBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_distribution_jobsTemplateBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponDistributionJob{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_distribution_jobs", ctx))

		if IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".id") == -1 {
				selects = append(selects, "customer_coupon_distribution_jobs"+".id")
			}

			if IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".template_id") == -1 {
				selects = append(selects, TableName("customer_coupon_distribution_jobs", ctx)+".template_id")
			}
		}

		res := db.Query().Table(TableName("customer_coupon_distribution_jobs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "template_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerCouponDistributionJob, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.TemplateID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerCouponDistributionJob{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponDistributionJob with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerCouponDistributionJobTemplate"] = dataloader.NewBatchedLoader(customer_coupon_distribution_jobsTemplateBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["TemplateAndCustomerCouponDistributionJobIds"] = dataloader.NewBatchedLoader(customer_coupon_distribution_jobsTemplateBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_distribution_jobsMemberBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponDistributionJob{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_distribution_jobs", ctx))

		if IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".*") == -1 {
			if IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".id") == -1 {
				selects = append(selects, "customer_coupon_distribution_jobs"+".id")
			}

			if IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".member_id") == -1 {
				selects = append(selects, TableName("customer_coupon_distribution_jobs", ctx)+".member_id")
			}
		}

		res := db.Query().Table(TableName("customer_coupon_distribution_jobs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "member_id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string][]*CustomerCouponDistributionJob, len(keys))
		for _, v := range *items {
			item := v

			mapKey := *item.MemberID

			if itemMap[mapKey] == nil {
				itemMap[mapKey] = []*CustomerCouponDistributionJob{}
			}
			itemMap[mapKey] = append(itemMap[mapKey], &item)
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponDistributionJob with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  item,
					Error: nil,
				})
			}
		}
		return results
	}
	loaders["CustomerCouponDistributionJobMember"] = dataloader.NewBatchedLoader(customer_coupon_distribution_jobsMemberBatchFn, dataloader.WithClearCacheOnBatch())
	loaders["MemberAndCustomerCouponDistributionJobIds"] = dataloader.NewBatchedLoader(customer_coupon_distribution_jobsMemberBatchFn, dataloader.WithClearCacheOnBatch())

	customer_coupon_distribution_jobsBatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]CustomerCouponDistributionJob{}
		selects := GetFieldsRequested(ctx, TableName("customer_coupon_distribution_jobs", ctx))
		if len(selects) > 0 && IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".*") == -1 && IndexOf(selects, TableName("customer_coupon_distribution_jobs", ctx)+".id") == -1 {
			selects = append(selects, TableName("customer_coupon_distribution_jobs", ctx)+".id")
		}

		res := db.Query().Table(TableName("customer_coupon_distribution_jobs", ctx)).Select(selects).Order("weight ASC, created_at ASC").Find(items, "id IN (?)", ids)
		if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*dataloader.Result{
				{Error: res.Error},
			}
		}

		itemMap := make(map[string]CustomerCouponDistributionJob, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: nil,
					// Error: fmt.Errorf("CustomerCouponDistributionJob with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["CustomerCouponDistributionJob"] = dataloader.NewBatchedLoader(customer_coupon_distribution_jobsBatchFn, dataloader.WithClearCacheOnBatch())

	return loaders
}
