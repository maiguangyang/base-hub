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

	return loaders
}
