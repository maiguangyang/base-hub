package gen

import (
	"context"
	"fmt"
	"strings"
	"time"

	"base-engine/auth"
	"base-engine/utils"

	"github.com/gofrs/uuid"
	"gorm.io/gorm/clause"
)

// ============================================================
// 类型定义
// ============================================================

// GeneratedMutationResolver 生成的 Mutation 解析器
type GeneratedMutationResolver struct{ *GeneratedResolver }

// MutationEvents 变更事件集合
type MutationEvents struct {
	Events []Event
}

// ============================================================
// 实体 Mutation 解析器
// ============================================================

// ============================================================
// Account - Create
// ============================================================

// CreateAccount 创建 Account 实体的解析器入口
func (r *GeneratedMutationResolver) CreateAccount(ctx context.Context, input map[string]interface{}) (item *Account, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateAccount(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateAccountHandler 处理 Account 创建逻辑
func CreateAccountHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Account, err error) {
	item = &Account{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "Account",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes AccountChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	// ToMany: memberships - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["memberships"]) && !utils.IsNil(input["membershipsIds"]) {
		return nil, fmt.Errorf("membershipsIds and memberships cannot coexist")
	}

	// ToMany: initializedOrganizations - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["initializedOrganizations"]) && !utils.IsNil(input["initializedOrganizationsIds"]) {
		return nil, fmt.Errorf("initializedOrganizationsIds and initializedOrganizations cannot coexist")
	}

	// ToMany: openingRecords - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["openingRecords"]) && !utils.IsNil(input["openingRecordsIds"]) {
		return nil, fmt.Errorf("openingRecordsIds and openingRecords cannot coexist")
	}

	// ToMany: recordedOpeningRecords - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["recordedOpeningRecords"]) && !utils.IsNil(input["recordedOpeningRecordsIds"]) {
		return nil, fmt.Errorf("recordedOpeningRecordsIds and recordedOpeningRecords cannot coexist")
	}

	// ToMany: sessions - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["sessions"]) && !utils.IsNil(input["sessionsIds"]) {
		return nil, fmt.Errorf("sessionsIds and sessions cannot coexist")
	}

	// ToMany: reviewedStores - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["reviewedStores"]) && !utils.IsNil(input["reviewedStoresIds"]) {
		return nil, fmt.Errorf("reviewedStoresIds and reviewedStores cannot coexist")
	}

	// ToMany: sentMembershipInvitations - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["sentMembershipInvitations"]) && !utils.IsNil(input["sentMembershipInvitationsIds"]) {
		return nil, fmt.Errorf("sentMembershipInvitationsIds and sentMembershipInvitations cannot coexist")
	}

	// ToMany: auditLogs - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["phone"]; ok {

		item.Phone = changes.Phone

		event.AddNewValue("phone", changes.Phone)
	}

	if _, ok := input["displayName"]; ok {

		item.DisplayName = changes.DisplayName

		event.AddNewValue("displayName", changes.DisplayName)
	}

	if _, ok := input["email"]; ok && changes.Email != nil {

		item.Email = changes.Email

		event.AddNewValue("email", changes.Email)
	}

	if _, ok := input["status"]; ok {

		item.Status = changes.Status

		event.AddNewValue("status", changes.Status)
	}

	if _, ok := input["mustChangePassword"]; ok {

		item.MustChangePassword = changes.MustChangePassword

		event.AddNewValue("mustChangePassword", changes.MustChangePassword)
	}

	if _, ok := input["credentialVersion"]; ok {

		item.CredentialVersion = changes.CredentialVersion

		event.AddNewValue("credentialVersion", changes.CredentialVersion)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("accounts", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// ---------- ToMany: memberships (OneToMany 外键在 OperatorMembership.account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["membershipsIds"]; ok && !utils.IsNil(input["membershipsIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membershipsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("memberships", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["memberships"]; ok && !utils.IsNil(input["memberships"]) {
		newMemberships := []*OperatorMembership{}
		updateMemberships := []*OperatorMembership{}

		hasCreateMemberships := false
		hasUpdateMemberships := false

		for index, v := range changes.Memberships {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMemberships {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMemberships = true
				}

				membershipsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membershipsInput["id"].(string), membershipsInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateMemberships = append(updateMemberships, v)
			} else {
				// 创建新记录
				if !hasCreateMemberships {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMemberships = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.AccountID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMemberships = append(newMemberships, v)
			}
		}

		allItems := append(updateMemberships, newMemberships...)

		event.AddNewValue("memberships", allItems)
	}

	// ---------- ToMany: initializedOrganizations (OneToMany 外键在 Organization.initial_account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["initializedOrganizationsIds"]; ok && !utils.IsNil(input["initializedOrganizationsIds"]) {
		items := []*Organization{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
				return item, fmt.Errorf("Organization Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("initializedOrganizationsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("initializedOrganizations", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["initializedOrganizations"]; ok && !utils.IsNil(input["initializedOrganizations"]) {
		newInitializedOrganizations := []*Organization{}
		updateInitializedOrganizations := []*Organization{}

		hasCreateInitializedOrganizations := false
		hasUpdateInitializedOrganizations := false

		for index, v := range changes.InitializedOrganizations {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateInitializedOrganizations {
					if err := auth.CheckAuthorization(ctx, "UpdateOrganization"); err != nil {
						return item, fmt.Errorf("UpdateOrganization: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
						return item, fmt.Errorf("Organization Detail: %w", err)
					}
					hasUpdateInitializedOrganizations = true
				}

				initializedOrganizationsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOrganization(ctx, r, initializedOrganizationsInput["id"].(string), initializedOrganizationsInput); err != nil {
					return item, fmt.Errorf("Organization ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateInitializedOrganizations = append(updateInitializedOrganizations, v)
			} else {
				// 创建新记录
				if !hasCreateInitializedOrganizations {
					if err := auth.CheckAuthorization(ctx, "CreateOrganization"); err != nil {
						return item, fmt.Errorf("CreateOrganization: %w", err)
					}
					hasCreateInitializedOrganizations = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.InitialAccountID = &item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("organizations", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newInitializedOrganizations = append(newInitializedOrganizations, v)
			}
		}

		allItems := append(updateInitializedOrganizations, newInitializedOrganizations...)

		event.AddNewValue("initializedOrganizations", allItems)
	}

	// ---------- ToMany: openingRecords (OneToMany 外键在 FranchiseOpeningRecord.initial_account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["openingRecordsIds"]; ok && !utils.IsNil(input["openingRecordsIds"]) {
		items := []*FranchiseOpeningRecord{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
				return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("openingRecordsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("openingRecords", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["openingRecords"]; ok && !utils.IsNil(input["openingRecords"]) {
		newOpeningRecords := []*FranchiseOpeningRecord{}
		updateOpeningRecords := []*FranchiseOpeningRecord{}

		hasCreateOpeningRecords := false
		hasUpdateOpeningRecords := false

		for index, v := range changes.OpeningRecords {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("UpdateFranchiseOpeningRecord: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
					}
					hasUpdateOpeningRecords = true
				}

				openingRecordsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchiseOpeningRecord(ctx, r, openingRecordsInput["id"].(string), openingRecordsInput); err != nil {
					return item, fmt.Errorf("FranchiseOpeningRecord ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateOpeningRecords = append(updateOpeningRecords, v)
			} else {
				// 创建新记录
				if !hasCreateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "CreateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("CreateFranchiseOpeningRecord: %w", err)
					}
					hasCreateOpeningRecords = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.InitialAccountID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("franchise_opening_records", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newOpeningRecords = append(newOpeningRecords, v)
			}
		}

		allItems := append(updateOpeningRecords, newOpeningRecords...)

		event.AddNewValue("openingRecords", allItems)
	}

	// ---------- ToMany: recordedOpeningRecords (OneToMany 外键在 FranchiseOpeningRecord.recorded_by_account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["recordedOpeningRecordsIds"]; ok && !utils.IsNil(input["recordedOpeningRecordsIds"]) {
		items := []*FranchiseOpeningRecord{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
				return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("recordedOpeningRecordsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("recorded_by_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("recordedOpeningRecords", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["recordedOpeningRecords"]; ok && !utils.IsNil(input["recordedOpeningRecords"]) {
		newRecordedOpeningRecords := []*FranchiseOpeningRecord{}
		updateRecordedOpeningRecords := []*FranchiseOpeningRecord{}

		hasCreateRecordedOpeningRecords := false
		hasUpdateRecordedOpeningRecords := false

		for index, v := range changes.RecordedOpeningRecords {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRecordedOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("UpdateFranchiseOpeningRecord: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
					}
					hasUpdateRecordedOpeningRecords = true
				}

				recordedOpeningRecordsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchiseOpeningRecord(ctx, r, recordedOpeningRecordsInput["id"].(string), recordedOpeningRecordsInput); err != nil {
					return item, fmt.Errorf("FranchiseOpeningRecord ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("recorded_by_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateRecordedOpeningRecords = append(updateRecordedOpeningRecords, v)
			} else {
				// 创建新记录
				if !hasCreateRecordedOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "CreateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("CreateFranchiseOpeningRecord: %w", err)
					}
					hasCreateRecordedOpeningRecords = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.RecordedByAccountID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("franchise_opening_records", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRecordedOpeningRecords = append(newRecordedOpeningRecords, v)
			}
		}

		allItems := append(updateRecordedOpeningRecords, newRecordedOpeningRecords...)

		event.AddNewValue("recordedOpeningRecords", allItems)
	}

	// ---------- ToMany: sessions (OneToMany 外键在 Session.account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["sessionsIds"]; ok && !utils.IsNil(input["sessionsIds"]) {
		items := []*Session{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
				return item, fmt.Errorf("Session Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("sessionsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("sessions", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["sessions"]; ok && !utils.IsNil(input["sessions"]) {
		newSessions := []*Session{}
		updateSessions := []*Session{}

		hasCreateSessions := false
		hasUpdateSessions := false

		for index, v := range changes.Sessions {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateSessions {
					if err := auth.CheckAuthorization(ctx, "UpdateSession"); err != nil {
						return item, fmt.Errorf("UpdateSession: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
						return item, fmt.Errorf("Session Detail: %w", err)
					}
					hasUpdateSessions = true
				}

				sessionsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateSession(ctx, r, sessionsInput["id"].(string), sessionsInput); err != nil {
					return item, fmt.Errorf("Session ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateSessions = append(updateSessions, v)
			} else {
				// 创建新记录
				if !hasCreateSessions {
					if err := auth.CheckAuthorization(ctx, "CreateSession"); err != nil {
						return item, fmt.Errorf("CreateSession: %w", err)
					}
					hasCreateSessions = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.AccountID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("sessions", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newSessions = append(newSessions, v)
			}
		}

		allItems := append(updateSessions, newSessions...)

		event.AddNewValue("sessions", allItems)
	}

	// ---------- ToMany: reviewedStores (OneToMany 外键在 Store.reviewed_by_account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["reviewedStoresIds"]; ok && !utils.IsNil(input["reviewedStoresIds"]) {
		items := []*Store{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
				return item, fmt.Errorf("Store Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("reviewedStoresIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("reviewed_by_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("reviewedStores", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["reviewedStores"]; ok && !utils.IsNil(input["reviewedStores"]) {
		newReviewedStores := []*Store{}
		updateReviewedStores := []*Store{}

		hasCreateReviewedStores := false
		hasUpdateReviewedStores := false

		for index, v := range changes.ReviewedStores {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateReviewedStores {
					if err := auth.CheckAuthorization(ctx, "UpdateStore"); err != nil {
						return item, fmt.Errorf("UpdateStore: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
						return item, fmt.Errorf("Store Detail: %w", err)
					}
					hasUpdateReviewedStores = true
				}

				reviewedStoresInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStore(ctx, r, reviewedStoresInput["id"].(string), reviewedStoresInput); err != nil {
					return item, fmt.Errorf("Store ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("reviewed_by_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateReviewedStores = append(updateReviewedStores, v)
			} else {
				// 创建新记录
				if !hasCreateReviewedStores {
					if err := auth.CheckAuthorization(ctx, "CreateStore"); err != nil {
						return item, fmt.Errorf("CreateStore: %w", err)
					}
					hasCreateReviewedStores = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.ReviewedByAccountID = &item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("stores", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newReviewedStores = append(newReviewedStores, v)
			}
		}

		allItems := append(updateReviewedStores, newReviewedStores...)

		event.AddNewValue("reviewedStores", allItems)
	}

	// ---------- ToMany: sentMembershipInvitations (OneToMany 外键在 MembershipInvitation.invited_by_account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["sentMembershipInvitationsIds"]; ok && !utils.IsNil(input["sentMembershipInvitationsIds"]) {
		items := []*MembershipInvitation{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
				return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("sentMembershipInvitationsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("invited_by_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("sentMembershipInvitations", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["sentMembershipInvitations"]; ok && !utils.IsNil(input["sentMembershipInvitations"]) {
		newSentMembershipInvitations := []*MembershipInvitation{}
		updateSentMembershipInvitations := []*MembershipInvitation{}

		hasCreateSentMembershipInvitations := false
		hasUpdateSentMembershipInvitations := false

		for index, v := range changes.SentMembershipInvitations {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateSentMembershipInvitations {
					if err := auth.CheckAuthorization(ctx, "UpdateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("UpdateMembershipInvitation: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
						return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
					}
					hasUpdateSentMembershipInvitations = true
				}

				sentMembershipInvitationsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateMembershipInvitation(ctx, r, sentMembershipInvitationsInput["id"].(string), sentMembershipInvitationsInput); err != nil {
					return item, fmt.Errorf("MembershipInvitation ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("invited_by_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateSentMembershipInvitations = append(updateSentMembershipInvitations, v)
			} else {
				// 创建新记录
				if !hasCreateSentMembershipInvitations {
					if err := auth.CheckAuthorization(ctx, "CreateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("CreateMembershipInvitation: %w", err)
					}
					hasCreateSentMembershipInvitations = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.InvitedByAccountID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("membership_invitations", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newSentMembershipInvitations = append(newSentMembershipInvitations, v)
			}
		}

		allItems := append(updateSentMembershipInvitations, newSentMembershipInvitations...)

		event.AddNewValue("sentMembershipInvitations", allItems)
	}

	// ---------- ToMany: auditLogs (OneToMany 外键在 AuditLog.actor_account_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("actor_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("actor_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.ActorAccountID = &item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Account - Update
// ============================================================

// UpdateAccount 更新 Account 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateAccount(ctx context.Context, id string, input map[string]interface{}) (item *Account, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateAccount(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateAccountHandler 处理 Account 更新逻辑
func UpdateAccountHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Account, err error) {
	item = &Account{}
	newItem := &Account{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "Account",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes AccountChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["memberships"]) && !utils.IsNil(input["membershipsIds"]) {
		return nil, fmt.Errorf("membershipsIds and memberships cannot coexist")
	}

	if !utils.IsNil(input["initializedOrganizations"]) && !utils.IsNil(input["initializedOrganizationsIds"]) {
		return nil, fmt.Errorf("initializedOrganizationsIds and initializedOrganizations cannot coexist")
	}

	if !utils.IsNil(input["openingRecords"]) && !utils.IsNil(input["openingRecordsIds"]) {
		return nil, fmt.Errorf("openingRecordsIds and openingRecords cannot coexist")
	}

	if !utils.IsNil(input["recordedOpeningRecords"]) && !utils.IsNil(input["recordedOpeningRecordsIds"]) {
		return nil, fmt.Errorf("recordedOpeningRecordsIds and recordedOpeningRecords cannot coexist")
	}

	if !utils.IsNil(input["sessions"]) && !utils.IsNil(input["sessionsIds"]) {
		return nil, fmt.Errorf("sessionsIds and sessions cannot coexist")
	}

	if !utils.IsNil(input["reviewedStores"]) && !utils.IsNil(input["reviewedStoresIds"]) {
		return nil, fmt.Errorf("reviewedStoresIds and reviewedStores cannot coexist")
	}

	if !utils.IsNil(input["sentMembershipInvitations"]) && !utils.IsNil(input["sentMembershipInvitationsIds"]) {
		return nil, fmt.Errorf("sentMembershipInvitationsIds and sentMembershipInvitations cannot coexist")
	}

	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("accounts", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["phone"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Phone != changes.Phone {

			event.AddOldValue("phone", item.Phone)
			event.AddNewValue("phone", changes.Phone)

			item.Phone = changes.Phone
			newItem.Phone = changes.Phone
			changedFields = append(changedFields, "phone")
			isChange = true
		}
	}

	if _, ok := input["displayName"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.DisplayName != changes.DisplayName {

			event.AddOldValue("displayName", item.DisplayName)
			event.AddNewValue("displayName", changes.DisplayName)

			item.DisplayName = changes.DisplayName
			newItem.DisplayName = changes.DisplayName
			changedFields = append(changedFields, "display_name")
			isChange = true
		}
	}

	if _, ok := input["email"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Email != changes.Email) && (item.Email == nil || changes.Email == nil || *item.Email != *changes.Email) {

			event.AddOldValue("email", item.Email)
			event.AddNewValue("email", changes.Email)

			item.Email = changes.Email
			newItem.Email = changes.Email
			changedFields = append(changedFields, "email")
			isChange = true
		}
	}

	if _, ok := input["status"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Status != changes.Status {

			event.AddOldValue("status", item.Status)
			event.AddNewValue("status", changes.Status)

			item.Status = changes.Status
			newItem.Status = changes.Status
			changedFields = append(changedFields, "status")
			isChange = true
		}
	}

	if _, ok := input["mustChangePassword"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.MustChangePassword != changes.MustChangePassword {

			event.AddOldValue("mustChangePassword", item.MustChangePassword)
			event.AddNewValue("mustChangePassword", changes.MustChangePassword)

			item.MustChangePassword = changes.MustChangePassword
			newItem.MustChangePassword = changes.MustChangePassword
			changedFields = append(changedFields, "must_change_password")
			isChange = true
		}
	}

	if _, ok := input["credentialVersion"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.CredentialVersion != changes.CredentialVersion {

			event.AddOldValue("credentialVersion", item.CredentialVersion)
			event.AddNewValue("credentialVersion", changes.CredentialVersion)

			item.CredentialVersion = changes.CredentialVersion
			newItem.CredentialVersion = changes.CredentialVersion
			changedFields = append(changedFields, "credential_version")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("accounts", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// ---------- ToMany: memberships ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["membershipsIds"]; ok && !utils.IsNil(input["membershipsIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membershipsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&OperatorMembership{}).Where("account_id = ?", item.ID).Update("account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&OperatorMembership{}).Where("account_id = ?", item.ID).Update("account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("memberships", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["memberships"]; ok && !utils.IsNil(input["memberships"]) {
		newMemberships := []*OperatorMembership{}
		updateMemberships := []*OperatorMembership{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&OperatorMembership{}).Where("account_id = ?", item.ID).Update("account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateMemberships := false
		hasUpdateMemberships := false

		for index, v := range changes.Memberships {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMemberships {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMemberships = true
				}

				membershipsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membershipsInput["id"].(string), membershipsInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateMemberships = append(updateMemberships, v)
			} else {
				// 创建新记录
				if !hasCreateMemberships {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMemberships = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.AccountID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMemberships = append(newMemberships, v)
			}
		}

		allItems := append(updateMemberships, newMemberships...)

		event.AddNewValue("memberships", allItems)
	}

	// ---------- ToMany: initializedOrganizations ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["initializedOrganizationsIds"]; ok && !utils.IsNil(input["initializedOrganizationsIds"]) {
		items := []*Organization{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
				return item, fmt.Errorf("Organization Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("initializedOrganizationsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&Organization{}).Where("initial_account_id = ?", item.ID).Update("initial_account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&Organization{}).Where("initial_account_id = ?", item.ID).Update("initial_account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("initializedOrganizations", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["initializedOrganizations"]; ok && !utils.IsNil(input["initializedOrganizations"]) {
		newInitializedOrganizations := []*Organization{}
		updateInitializedOrganizations := []*Organization{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&Organization{}).Where("initial_account_id = ?", item.ID).Update("initial_account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateInitializedOrganizations := false
		hasUpdateInitializedOrganizations := false

		for index, v := range changes.InitializedOrganizations {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateInitializedOrganizations {
					if err := auth.CheckAuthorization(ctx, "UpdateOrganization"); err != nil {
						return item, fmt.Errorf("UpdateOrganization: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Organization"); err != nil {
						return item, fmt.Errorf("Organization Detail: %w", err)
					}
					hasUpdateInitializedOrganizations = true
				}

				initializedOrganizationsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOrganization(ctx, r, initializedOrganizationsInput["id"].(string), initializedOrganizationsInput); err != nil {
					return item, fmt.Errorf("Organization ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateInitializedOrganizations = append(updateInitializedOrganizations, v)
			} else {
				// 创建新记录
				if !hasCreateInitializedOrganizations {
					if err := auth.CheckAuthorization(ctx, "CreateOrganization"); err != nil {
						return item, fmt.Errorf("CreateOrganization: %w", err)
					}
					hasCreateInitializedOrganizations = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.InitialAccountID = &item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("organizations", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newInitializedOrganizations = append(newInitializedOrganizations, v)
			}
		}

		allItems := append(updateInitializedOrganizations, newInitializedOrganizations...)

		event.AddNewValue("initializedOrganizations", allItems)
	}

	// ---------- ToMany: openingRecords ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["openingRecordsIds"]; ok && !utils.IsNil(input["openingRecordsIds"]) {
		items := []*FranchiseOpeningRecord{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
				return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("openingRecordsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&FranchiseOpeningRecord{}).Where("initial_account_id = ?", item.ID).Update("initial_account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&FranchiseOpeningRecord{}).Where("initial_account_id = ?", item.ID).Update("initial_account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("openingRecords", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["openingRecords"]; ok && !utils.IsNil(input["openingRecords"]) {
		newOpeningRecords := []*FranchiseOpeningRecord{}
		updateOpeningRecords := []*FranchiseOpeningRecord{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&FranchiseOpeningRecord{}).Where("initial_account_id = ?", item.ID).Update("initial_account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateOpeningRecords := false
		hasUpdateOpeningRecords := false

		for index, v := range changes.OpeningRecords {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("UpdateFranchiseOpeningRecord: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
					}
					hasUpdateOpeningRecords = true
				}

				openingRecordsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchiseOpeningRecord(ctx, r, openingRecordsInput["id"].(string), openingRecordsInput); err != nil {
					return item, fmt.Errorf("FranchiseOpeningRecord ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("initial_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateOpeningRecords = append(updateOpeningRecords, v)
			} else {
				// 创建新记录
				if !hasCreateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "CreateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("CreateFranchiseOpeningRecord: %w", err)
					}
					hasCreateOpeningRecords = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.InitialAccountID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("franchise_opening_records", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newOpeningRecords = append(newOpeningRecords, v)
			}
		}

		allItems := append(updateOpeningRecords, newOpeningRecords...)

		event.AddNewValue("openingRecords", allItems)
	}

	// ---------- ToMany: recordedOpeningRecords ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["recordedOpeningRecordsIds"]; ok && !utils.IsNil(input["recordedOpeningRecordsIds"]) {
		items := []*FranchiseOpeningRecord{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
				return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("recordedOpeningRecordsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&FranchiseOpeningRecord{}).Where("recorded_by_account_id = ?", item.ID).Update("recorded_by_account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("recorded_by_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&FranchiseOpeningRecord{}).Where("recorded_by_account_id = ?", item.ID).Update("recorded_by_account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("recordedOpeningRecords", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["recordedOpeningRecords"]; ok && !utils.IsNil(input["recordedOpeningRecords"]) {
		newRecordedOpeningRecords := []*FranchiseOpeningRecord{}
		updateRecordedOpeningRecords := []*FranchiseOpeningRecord{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&FranchiseOpeningRecord{}).Where("recorded_by_account_id = ?", item.ID).Update("recorded_by_account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateRecordedOpeningRecords := false
		hasUpdateRecordedOpeningRecords := false

		for index, v := range changes.RecordedOpeningRecords {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRecordedOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("UpdateFranchiseOpeningRecord: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
					}
					hasUpdateRecordedOpeningRecords = true
				}

				recordedOpeningRecordsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchiseOpeningRecord(ctx, r, recordedOpeningRecordsInput["id"].(string), recordedOpeningRecordsInput); err != nil {
					return item, fmt.Errorf("FranchiseOpeningRecord ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("recorded_by_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateRecordedOpeningRecords = append(updateRecordedOpeningRecords, v)
			} else {
				// 创建新记录
				if !hasCreateRecordedOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "CreateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("CreateFranchiseOpeningRecord: %w", err)
					}
					hasCreateRecordedOpeningRecords = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.RecordedByAccountID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("franchise_opening_records", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRecordedOpeningRecords = append(newRecordedOpeningRecords, v)
			}
		}

		allItems := append(updateRecordedOpeningRecords, newRecordedOpeningRecords...)

		event.AddNewValue("recordedOpeningRecords", allItems)
	}

	// ---------- ToMany: sessions ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["sessionsIds"]; ok && !utils.IsNil(input["sessionsIds"]) {
		items := []*Session{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
				return item, fmt.Errorf("Session Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("sessionsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&Session{}).Where("account_id = ?", item.ID).Update("account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&Session{}).Where("account_id = ?", item.ID).Update("account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("sessions", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["sessions"]; ok && !utils.IsNil(input["sessions"]) {
		newSessions := []*Session{}
		updateSessions := []*Session{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&Session{}).Where("account_id = ?", item.ID).Update("account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateSessions := false
		hasUpdateSessions := false

		for index, v := range changes.Sessions {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateSessions {
					if err := auth.CheckAuthorization(ctx, "UpdateSession"); err != nil {
						return item, fmt.Errorf("UpdateSession: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
						return item, fmt.Errorf("Session Detail: %w", err)
					}
					hasUpdateSessions = true
				}

				sessionsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateSession(ctx, r, sessionsInput["id"].(string), sessionsInput); err != nil {
					return item, fmt.Errorf("Session ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateSessions = append(updateSessions, v)
			} else {
				// 创建新记录
				if !hasCreateSessions {
					if err := auth.CheckAuthorization(ctx, "CreateSession"); err != nil {
						return item, fmt.Errorf("CreateSession: %w", err)
					}
					hasCreateSessions = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.AccountID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("sessions", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newSessions = append(newSessions, v)
			}
		}

		allItems := append(updateSessions, newSessions...)

		event.AddNewValue("sessions", allItems)
	}

	// ---------- ToMany: reviewedStores ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["reviewedStoresIds"]; ok && !utils.IsNil(input["reviewedStoresIds"]) {
		items := []*Store{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
				return item, fmt.Errorf("Store Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("reviewedStoresIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&Store{}).Where("reviewed_by_account_id = ?", item.ID).Update("reviewed_by_account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("reviewed_by_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&Store{}).Where("reviewed_by_account_id = ?", item.ID).Update("reviewed_by_account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("reviewedStores", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["reviewedStores"]; ok && !utils.IsNil(input["reviewedStores"]) {
		newReviewedStores := []*Store{}
		updateReviewedStores := []*Store{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&Store{}).Where("reviewed_by_account_id = ?", item.ID).Update("reviewed_by_account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateReviewedStores := false
		hasUpdateReviewedStores := false

		for index, v := range changes.ReviewedStores {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateReviewedStores {
					if err := auth.CheckAuthorization(ctx, "UpdateStore"); err != nil {
						return item, fmt.Errorf("UpdateStore: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
						return item, fmt.Errorf("Store Detail: %w", err)
					}
					hasUpdateReviewedStores = true
				}

				reviewedStoresInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStore(ctx, r, reviewedStoresInput["id"].(string), reviewedStoresInput); err != nil {
					return item, fmt.Errorf("Store ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("reviewed_by_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateReviewedStores = append(updateReviewedStores, v)
			} else {
				// 创建新记录
				if !hasCreateReviewedStores {
					if err := auth.CheckAuthorization(ctx, "CreateStore"); err != nil {
						return item, fmt.Errorf("CreateStore: %w", err)
					}
					hasCreateReviewedStores = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.ReviewedByAccountID = &item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("stores", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newReviewedStores = append(newReviewedStores, v)
			}
		}

		allItems := append(updateReviewedStores, newReviewedStores...)

		event.AddNewValue("reviewedStores", allItems)
	}

	// ---------- ToMany: sentMembershipInvitations ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["sentMembershipInvitationsIds"]; ok && !utils.IsNil(input["sentMembershipInvitationsIds"]) {
		items := []*MembershipInvitation{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
				return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("sentMembershipInvitationsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&MembershipInvitation{}).Where("invited_by_account_id = ?", item.ID).Update("invited_by_account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("invited_by_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&MembershipInvitation{}).Where("invited_by_account_id = ?", item.ID).Update("invited_by_account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("sentMembershipInvitations", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["sentMembershipInvitations"]; ok && !utils.IsNil(input["sentMembershipInvitations"]) {
		newSentMembershipInvitations := []*MembershipInvitation{}
		updateSentMembershipInvitations := []*MembershipInvitation{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&MembershipInvitation{}).Where("invited_by_account_id = ?", item.ID).Update("invited_by_account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateSentMembershipInvitations := false
		hasUpdateSentMembershipInvitations := false

		for index, v := range changes.SentMembershipInvitations {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateSentMembershipInvitations {
					if err := auth.CheckAuthorization(ctx, "UpdateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("UpdateMembershipInvitation: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
						return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
					}
					hasUpdateSentMembershipInvitations = true
				}

				sentMembershipInvitationsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateMembershipInvitation(ctx, r, sentMembershipInvitationsInput["id"].(string), sentMembershipInvitationsInput); err != nil {
					return item, fmt.Errorf("MembershipInvitation ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("invited_by_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateSentMembershipInvitations = append(updateSentMembershipInvitations, v)
			} else {
				// 创建新记录
				if !hasCreateSentMembershipInvitations {
					if err := auth.CheckAuthorization(ctx, "CreateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("CreateMembershipInvitation: %w", err)
					}
					hasCreateSentMembershipInvitations = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.InvitedByAccountID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("membership_invitations", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newSentMembershipInvitations = append(newSentMembershipInvitations, v)
			}
		}

		allItems := append(updateSentMembershipInvitations, newSentMembershipInvitations...)

		event.AddNewValue("sentMembershipInvitations", allItems)
	}

	// ---------- ToMany: auditLogs ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&AuditLog{}).Where("actor_account_id = ?", item.ID).Update("actor_account_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("actor_account_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&AuditLog{}).Where("actor_account_id = ?", item.ID).Update("actor_account_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&AuditLog{}).Where("actor_account_id = ?", item.ID).Update("actor_account_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("actor_account_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.ActorAccountID = &item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Account - Delete
// ============================================================

// DeleteAccountFunc 执行删除或恢复操作
func DeleteAccountFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &Account{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("accounts", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "Account",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("accounts", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("accounts", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("accounts", ctx)).Updates(Account{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteAccounts 批量删除 Account 实体
func (r *GeneratedMutationResolver) DeleteAccounts(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteAccounts(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteAccountsHandler 处理批量删除逻辑
func DeleteAccountsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteAccountFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Account - Recovery
// ============================================================

// RecoveryAccounts 批量恢复 Account 实体
func (r *GeneratedMutationResolver) RecoveryAccounts(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryAccounts(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryAccountsHandler 处理批量恢复逻辑
func RecoveryAccountsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteAccountFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Organization - Create
// ============================================================

// CreateOrganization 创建 Organization 实体的解析器入口
func (r *GeneratedMutationResolver) CreateOrganization(ctx context.Context, input map[string]interface{}) (item *Organization, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateOrganization(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateOrganizationHandler 处理 Organization 创建逻辑
func CreateOrganizationHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Organization, err error) {
	item = &Organization{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "Organization",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes OrganizationChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	// ToMany: memberships - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["memberships"]) && !utils.IsNil(input["membershipsIds"]) {
		return nil, fmt.Errorf("membershipsIds and memberships cannot coexist")
	}

	if !utils.IsNil(input["initialAccount"]) && !utils.IsNil(input["initialAccountId"]) {
		return nil, fmt.Errorf("initialAccountId and initialAccount cannot coexist")
	}

	// ToMany: openingRecords - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["openingRecords"]) && !utils.IsNil(input["openingRecordsIds"]) {
		return nil, fmt.Errorf("openingRecordsIds and openingRecords cannot coexist")
	}

	// ToMany: stores - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["stores"]) && !utils.IsNil(input["storesIds"]) {
		return nil, fmt.Errorf("storesIds and stores cannot coexist")
	}

	// ToMany: roles - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["roles"]) && !utils.IsNil(input["rolesIds"]) {
		return nil, fmt.Errorf("rolesIds and roles cannot coexist")
	}

	// ToMany: sessions - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["sessions"]) && !utils.IsNil(input["sessionsIds"]) {
		return nil, fmt.Errorf("sessionsIds and sessions cannot coexist")
	}

	// ToMany: auditLogs - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	// ToMany: paymentConfigs - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["paymentConfigs"]) && !utils.IsNil(input["paymentConfigsIds"]) {
		return nil, fmt.Errorf("paymentConfigsIds and paymentConfigs cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["code"]; ok {

		item.Code = changes.Code

		event.AddNewValue("code", changes.Code)
	}

	if _, ok := input["name"]; ok {

		item.Name = changes.Name

		event.AddNewValue("name", changes.Name)
	}

	if _, ok := input["type"]; ok {

		item.Type = changes.Type

		event.AddNewValue("type", changes.Type)
	}

	if _, ok := input["status"]; ok {

		item.Status = changes.Status

		event.AddNewValue("status", changes.Status)
	}

	if _, ok := input["suspendedAt"]; ok && changes.SuspendedAt != nil {

		item.SuspendedAt = changes.SuspendedAt

		event.AddNewValue("suspendedAt", changes.SuspendedAt)
	}

	if _, ok := input["suspensionReasonCode"]; ok && changes.SuspensionReasonCode != nil {

		item.SuspensionReasonCode = changes.SuspensionReasonCode

		event.AddNewValue("suspensionReasonCode", changes.SuspensionReasonCode)
	}

	if _, ok := input["initialAccountId"]; ok && changes.InitialAccountID != nil {

		if !utils.IsNil(input["initialAccountId"]) {
			if err := tx.Select("id").Where("id = ?", input["initialAccountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("initialAccountId: %w", err)
			}
		}

		item.InitialAccountID = changes.InitialAccountID

		event.AddNewValue("initialAccountId", changes.InitialAccountID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("organizations", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// ---------- ToMany: memberships (OneToMany 外键在 OperatorMembership.organization_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["membershipsIds"]; ok && !utils.IsNil(input["membershipsIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membershipsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("memberships", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["memberships"]; ok && !utils.IsNil(input["memberships"]) {
		newMemberships := []*OperatorMembership{}
		updateMemberships := []*OperatorMembership{}

		hasCreateMemberships := false
		hasUpdateMemberships := false

		for index, v := range changes.Memberships {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMemberships {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMemberships = true
				}

				membershipsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membershipsInput["id"].(string), membershipsInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateMemberships = append(updateMemberships, v)
			} else {
				// 创建新记录
				if !hasCreateMemberships {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMemberships = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.OrganizationID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMemberships = append(newMemberships, v)
			}
		}

		allItems := append(updateMemberships, newMemberships...)

		event.AddNewValue("memberships", allItems)
	}

	// ---------- ToMany: openingRecords (OneToMany 外键在 FranchiseOpeningRecord.organization_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["openingRecordsIds"]; ok && !utils.IsNil(input["openingRecordsIds"]) {
		items := []*FranchiseOpeningRecord{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
				return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("openingRecordsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("openingRecords", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["openingRecords"]; ok && !utils.IsNil(input["openingRecords"]) {
		newOpeningRecords := []*FranchiseOpeningRecord{}
		updateOpeningRecords := []*FranchiseOpeningRecord{}

		hasCreateOpeningRecords := false
		hasUpdateOpeningRecords := false

		for index, v := range changes.OpeningRecords {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("UpdateFranchiseOpeningRecord: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
					}
					hasUpdateOpeningRecords = true
				}

				openingRecordsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchiseOpeningRecord(ctx, r, openingRecordsInput["id"].(string), openingRecordsInput); err != nil {
					return item, fmt.Errorf("FranchiseOpeningRecord ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateOpeningRecords = append(updateOpeningRecords, v)
			} else {
				// 创建新记录
				if !hasCreateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "CreateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("CreateFranchiseOpeningRecord: %w", err)
					}
					hasCreateOpeningRecords = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.OrganizationID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("franchise_opening_records", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newOpeningRecords = append(newOpeningRecords, v)
			}
		}

		allItems := append(updateOpeningRecords, newOpeningRecords...)

		event.AddNewValue("openingRecords", allItems)
	}

	// ---------- ToMany: stores (OneToMany 外键在 Store.organization_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["storesIds"]; ok && !utils.IsNil(input["storesIds"]) {
		items := []*Store{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
				return item, fmt.Errorf("Store Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("storesIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("stores", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["stores"]; ok && !utils.IsNil(input["stores"]) {
		newStores := []*Store{}
		updateStores := []*Store{}

		hasCreateStores := false
		hasUpdateStores := false

		for index, v := range changes.Stores {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateStores {
					if err := auth.CheckAuthorization(ctx, "UpdateStore"); err != nil {
						return item, fmt.Errorf("UpdateStore: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
						return item, fmt.Errorf("Store Detail: %w", err)
					}
					hasUpdateStores = true
				}

				storesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStore(ctx, r, storesInput["id"].(string), storesInput); err != nil {
					return item, fmt.Errorf("Store ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateStores = append(updateStores, v)
			} else {
				// 创建新记录
				if !hasCreateStores {
					if err := auth.CheckAuthorization(ctx, "CreateStore"); err != nil {
						return item, fmt.Errorf("CreateStore: %w", err)
					}
					hasCreateStores = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.OrganizationID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("stores", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newStores = append(newStores, v)
			}
		}

		allItems := append(updateStores, newStores...)

		event.AddNewValue("stores", allItems)
	}

	// ---------- ToMany: roles (OneToMany 外键在 OperatorRole.organization_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["rolesIds"]; ok && !utils.IsNil(input["rolesIds"]) {
		items := []*OperatorRole{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
				return item, fmt.Errorf("OperatorRole Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("rolesIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("roles", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["roles"]; ok && !utils.IsNil(input["roles"]) {
		newRoles := []*OperatorRole{}
		updateRoles := []*OperatorRole{}

		hasCreateRoles := false
		hasUpdateRoles := false

		for index, v := range changes.Roles {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRoles {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorRole"); err != nil {
						return item, fmt.Errorf("UpdateOperatorRole: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
						return item, fmt.Errorf("OperatorRole Detail: %w", err)
					}
					hasUpdateRoles = true
				}

				rolesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorRole(ctx, r, rolesInput["id"].(string), rolesInput); err != nil {
					return item, fmt.Errorf("OperatorRole ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateRoles = append(updateRoles, v)
			} else {
				// 创建新记录
				if !hasCreateRoles {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorRole"); err != nil {
						return item, fmt.Errorf("CreateOperatorRole: %w", err)
					}
					hasCreateRoles = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.OrganizationID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("operator_roles", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRoles = append(newRoles, v)
			}
		}

		allItems := append(updateRoles, newRoles...)

		event.AddNewValue("roles", allItems)
	}

	// ---------- ToMany: sessions (OneToMany 外键在 Session.organization_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["sessionsIds"]; ok && !utils.IsNil(input["sessionsIds"]) {
		items := []*Session{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
				return item, fmt.Errorf("Session Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("sessionsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("sessions", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["sessions"]; ok && !utils.IsNil(input["sessions"]) {
		newSessions := []*Session{}
		updateSessions := []*Session{}

		hasCreateSessions := false
		hasUpdateSessions := false

		for index, v := range changes.Sessions {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateSessions {
					if err := auth.CheckAuthorization(ctx, "UpdateSession"); err != nil {
						return item, fmt.Errorf("UpdateSession: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
						return item, fmt.Errorf("Session Detail: %w", err)
					}
					hasUpdateSessions = true
				}

				sessionsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateSession(ctx, r, sessionsInput["id"].(string), sessionsInput); err != nil {
					return item, fmt.Errorf("Session ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateSessions = append(updateSessions, v)
			} else {
				// 创建新记录
				if !hasCreateSessions {
					if err := auth.CheckAuthorization(ctx, "CreateSession"); err != nil {
						return item, fmt.Errorf("CreateSession: %w", err)
					}
					hasCreateSessions = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.OrganizationID = &item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("sessions", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newSessions = append(newSessions, v)
			}
		}

		allItems := append(updateSessions, newSessions...)

		event.AddNewValue("sessions", allItems)
	}

	// ---------- ToMany: auditLogs (OneToMany 外键在 AuditLog.organization_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.OrganizationID = &item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// ---------- ToMany: paymentConfigs (OneToMany 外键在 FranchisePaymentConfig.organization_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["paymentConfigsIds"]; ok && !utils.IsNil(input["paymentConfigsIds"]) {
		items := []*FranchisePaymentConfig{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "FranchisePaymentConfig"); err != nil {
				return item, fmt.Errorf("FranchisePaymentConfig Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("paymentConfigsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("paymentConfigs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["paymentConfigs"]; ok && !utils.IsNil(input["paymentConfigs"]) {
		newPaymentConfigs := []*FranchisePaymentConfig{}
		updatePaymentConfigs := []*FranchisePaymentConfig{}

		hasCreatePaymentConfigs := false
		hasUpdatePaymentConfigs := false

		for index, v := range changes.PaymentConfigs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchisePaymentConfig"); err != nil {
						return item, fmt.Errorf("UpdateFranchisePaymentConfig: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchisePaymentConfig"); err != nil {
						return item, fmt.Errorf("FranchisePaymentConfig Detail: %w", err)
					}
					hasUpdatePaymentConfigs = true
				}

				paymentConfigsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchisePaymentConfig(ctx, r, paymentConfigsInput["id"].(string), paymentConfigsInput); err != nil {
					return item, fmt.Errorf("FranchisePaymentConfig ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updatePaymentConfigs = append(updatePaymentConfigs, v)
			} else {
				// 创建新记录
				if !hasCreatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "CreateFranchisePaymentConfig"); err != nil {
						return item, fmt.Errorf("CreateFranchisePaymentConfig: %w", err)
					}
					hasCreatePaymentConfigs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.OrganizationID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("franchise_payment_configs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newPaymentConfigs = append(newPaymentConfigs, v)
			}
		}

		allItems := append(updatePaymentConfigs, newPaymentConfigs...)

		event.AddNewValue("paymentConfigs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Organization - Update
// ============================================================

// UpdateOrganization 更新 Organization 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateOrganization(ctx context.Context, id string, input map[string]interface{}) (item *Organization, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateOrganization(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateOrganizationHandler 处理 Organization 更新逻辑
func UpdateOrganizationHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Organization, err error) {
	item = &Organization{}
	newItem := &Organization{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "Organization",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes OrganizationChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["memberships"]) && !utils.IsNil(input["membershipsIds"]) {
		return nil, fmt.Errorf("membershipsIds and memberships cannot coexist")
	}

	if !utils.IsNil(input["initialAccount"]) && !utils.IsNil(input["initialAccountId"]) {
		return nil, fmt.Errorf("initialAccountId and initialAccount cannot coexist")
	}

	if !utils.IsNil(input["openingRecords"]) && !utils.IsNil(input["openingRecordsIds"]) {
		return nil, fmt.Errorf("openingRecordsIds and openingRecords cannot coexist")
	}

	if !utils.IsNil(input["stores"]) && !utils.IsNil(input["storesIds"]) {
		return nil, fmt.Errorf("storesIds and stores cannot coexist")
	}

	if !utils.IsNil(input["roles"]) && !utils.IsNil(input["rolesIds"]) {
		return nil, fmt.Errorf("rolesIds and roles cannot coexist")
	}

	if !utils.IsNil(input["sessions"]) && !utils.IsNil(input["sessionsIds"]) {
		return nil, fmt.Errorf("sessionsIds and sessions cannot coexist")
	}

	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	if !utils.IsNil(input["paymentConfigs"]) && !utils.IsNil(input["paymentConfigsIds"]) {
		return nil, fmt.Errorf("paymentConfigsIds and paymentConfigs cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("organizations", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["code"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Code != changes.Code {

			event.AddOldValue("code", item.Code)
			event.AddNewValue("code", changes.Code)

			item.Code = changes.Code
			newItem.Code = changes.Code
			changedFields = append(changedFields, "code")
			isChange = true
		}
	}

	if _, ok := input["name"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Name != changes.Name {

			event.AddOldValue("name", item.Name)
			event.AddNewValue("name", changes.Name)

			item.Name = changes.Name
			newItem.Name = changes.Name
			changedFields = append(changedFields, "name")
			isChange = true
		}
	}

	if _, ok := input["type"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Type != changes.Type {

			event.AddOldValue("type", item.Type)
			event.AddNewValue("type", changes.Type)

			item.Type = changes.Type
			newItem.Type = changes.Type
			changedFields = append(changedFields, "type")
			isChange = true
		}
	}

	if _, ok := input["status"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Status != changes.Status {

			event.AddOldValue("status", item.Status)
			event.AddNewValue("status", changes.Status)

			item.Status = changes.Status
			newItem.Status = changes.Status
			changedFields = append(changedFields, "status")
			isChange = true
		}
	}

	if _, ok := input["suspendedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.SuspendedAt != changes.SuspendedAt) && (item.SuspendedAt == nil || changes.SuspendedAt == nil || *item.SuspendedAt != *changes.SuspendedAt) {

			event.AddOldValue("suspendedAt", item.SuspendedAt)
			event.AddNewValue("suspendedAt", changes.SuspendedAt)

			item.SuspendedAt = changes.SuspendedAt
			newItem.SuspendedAt = changes.SuspendedAt
			changedFields = append(changedFields, "suspended_at")
			isChange = true
		}
	}

	if _, ok := input["suspensionReasonCode"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.SuspensionReasonCode != changes.SuspensionReasonCode) && (item.SuspensionReasonCode == nil || changes.SuspensionReasonCode == nil || *item.SuspensionReasonCode != *changes.SuspensionReasonCode) {

			event.AddOldValue("suspensionReasonCode", item.SuspensionReasonCode)
			event.AddNewValue("suspensionReasonCode", changes.SuspensionReasonCode)

			item.SuspensionReasonCode = changes.SuspensionReasonCode
			newItem.SuspensionReasonCode = changes.SuspensionReasonCode
			changedFields = append(changedFields, "suspension_reason_code")
			isChange = true
		}
	}

	if _, ok := input["initialAccountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.InitialAccountID != changes.InitialAccountID) && (item.InitialAccountID == nil || changes.InitialAccountID == nil || *item.InitialAccountID != *changes.InitialAccountID) {

			if !utils.IsNil(input["initialAccountId"]) {
				if err := tx.Select("id").Where("id = ?", input["initialAccountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("initialAccountId: %w", err)
				}
			}

			event.AddOldValue("initialAccountId", item.InitialAccountID)
			event.AddNewValue("initialAccountId", changes.InitialAccountID)

			item.InitialAccountID = changes.InitialAccountID
			newItem.InitialAccountID = changes.InitialAccountID
			changedFields = append(changedFields, "initial_account_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("organizations", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// ---------- ToMany: memberships ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["membershipsIds"]; ok && !utils.IsNil(input["membershipsIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membershipsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&OperatorMembership{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&OperatorMembership{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("memberships", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["memberships"]; ok && !utils.IsNil(input["memberships"]) {
		newMemberships := []*OperatorMembership{}
		updateMemberships := []*OperatorMembership{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&OperatorMembership{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateMemberships := false
		hasUpdateMemberships := false

		for index, v := range changes.Memberships {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMemberships {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMemberships = true
				}

				membershipsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membershipsInput["id"].(string), membershipsInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateMemberships = append(updateMemberships, v)
			} else {
				// 创建新记录
				if !hasCreateMemberships {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMemberships = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.OrganizationID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMemberships = append(newMemberships, v)
			}
		}

		allItems := append(updateMemberships, newMemberships...)

		event.AddNewValue("memberships", allItems)
	}

	// ---------- ToMany: openingRecords ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["openingRecordsIds"]; ok && !utils.IsNil(input["openingRecordsIds"]) {
		items := []*FranchiseOpeningRecord{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
				return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("openingRecordsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&FranchiseOpeningRecord{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&FranchiseOpeningRecord{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("openingRecords", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["openingRecords"]; ok && !utils.IsNil(input["openingRecords"]) {
		newOpeningRecords := []*FranchiseOpeningRecord{}
		updateOpeningRecords := []*FranchiseOpeningRecord{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&FranchiseOpeningRecord{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateOpeningRecords := false
		hasUpdateOpeningRecords := false

		for index, v := range changes.OpeningRecords {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("UpdateFranchiseOpeningRecord: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("FranchiseOpeningRecord Detail: %w", err)
					}
					hasUpdateOpeningRecords = true
				}

				openingRecordsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchiseOpeningRecord(ctx, r, openingRecordsInput["id"].(string), openingRecordsInput); err != nil {
					return item, fmt.Errorf("FranchiseOpeningRecord ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateOpeningRecords = append(updateOpeningRecords, v)
			} else {
				// 创建新记录
				if !hasCreateOpeningRecords {
					if err := auth.CheckAuthorization(ctx, "CreateFranchiseOpeningRecord"); err != nil {
						return item, fmt.Errorf("CreateFranchiseOpeningRecord: %w", err)
					}
					hasCreateOpeningRecords = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.OrganizationID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("franchise_opening_records", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newOpeningRecords = append(newOpeningRecords, v)
			}
		}

		allItems := append(updateOpeningRecords, newOpeningRecords...)

		event.AddNewValue("openingRecords", allItems)
	}

	// ---------- ToMany: stores ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["storesIds"]; ok && !utils.IsNil(input["storesIds"]) {
		items := []*Store{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
				return item, fmt.Errorf("Store Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("storesIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&Store{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&Store{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("stores", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["stores"]; ok && !utils.IsNil(input["stores"]) {
		newStores := []*Store{}
		updateStores := []*Store{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&Store{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateStores := false
		hasUpdateStores := false

		for index, v := range changes.Stores {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateStores {
					if err := auth.CheckAuthorization(ctx, "UpdateStore"); err != nil {
						return item, fmt.Errorf("UpdateStore: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
						return item, fmt.Errorf("Store Detail: %w", err)
					}
					hasUpdateStores = true
				}

				storesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStore(ctx, r, storesInput["id"].(string), storesInput); err != nil {
					return item, fmt.Errorf("Store ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateStores = append(updateStores, v)
			} else {
				// 创建新记录
				if !hasCreateStores {
					if err := auth.CheckAuthorization(ctx, "CreateStore"); err != nil {
						return item, fmt.Errorf("CreateStore: %w", err)
					}
					hasCreateStores = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.OrganizationID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("stores", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newStores = append(newStores, v)
			}
		}

		allItems := append(updateStores, newStores...)

		event.AddNewValue("stores", allItems)
	}

	// ---------- ToMany: roles ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["rolesIds"]; ok && !utils.IsNil(input["rolesIds"]) {
		items := []*OperatorRole{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
				return item, fmt.Errorf("OperatorRole Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("rolesIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&OperatorRole{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&OperatorRole{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("roles", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["roles"]; ok && !utils.IsNil(input["roles"]) {
		newRoles := []*OperatorRole{}
		updateRoles := []*OperatorRole{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&OperatorRole{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateRoles := false
		hasUpdateRoles := false

		for index, v := range changes.Roles {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRoles {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorRole"); err != nil {
						return item, fmt.Errorf("UpdateOperatorRole: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
						return item, fmt.Errorf("OperatorRole Detail: %w", err)
					}
					hasUpdateRoles = true
				}

				rolesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorRole(ctx, r, rolesInput["id"].(string), rolesInput); err != nil {
					return item, fmt.Errorf("OperatorRole ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateRoles = append(updateRoles, v)
			} else {
				// 创建新记录
				if !hasCreateRoles {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorRole"); err != nil {
						return item, fmt.Errorf("CreateOperatorRole: %w", err)
					}
					hasCreateRoles = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.OrganizationID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("operator_roles", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRoles = append(newRoles, v)
			}
		}

		allItems := append(updateRoles, newRoles...)

		event.AddNewValue("roles", allItems)
	}

	// ---------- ToMany: sessions ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["sessionsIds"]; ok && !utils.IsNil(input["sessionsIds"]) {
		items := []*Session{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
				return item, fmt.Errorf("Session Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("sessionsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&Session{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&Session{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("sessions", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["sessions"]; ok && !utils.IsNil(input["sessions"]) {
		newSessions := []*Session{}
		updateSessions := []*Session{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&Session{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateSessions := false
		hasUpdateSessions := false

		for index, v := range changes.Sessions {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateSessions {
					if err := auth.CheckAuthorization(ctx, "UpdateSession"); err != nil {
						return item, fmt.Errorf("UpdateSession: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Session"); err != nil {
						return item, fmt.Errorf("Session Detail: %w", err)
					}
					hasUpdateSessions = true
				}

				sessionsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateSession(ctx, r, sessionsInput["id"].(string), sessionsInput); err != nil {
					return item, fmt.Errorf("Session ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateSessions = append(updateSessions, v)
			} else {
				// 创建新记录
				if !hasCreateSessions {
					if err := auth.CheckAuthorization(ctx, "CreateSession"); err != nil {
						return item, fmt.Errorf("CreateSession: %w", err)
					}
					hasCreateSessions = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.OrganizationID = &item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("sessions", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newSessions = append(newSessions, v)
			}
		}

		allItems := append(updateSessions, newSessions...)

		event.AddNewValue("sessions", allItems)
	}

	// ---------- ToMany: auditLogs ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&AuditLog{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&AuditLog{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&AuditLog{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.OrganizationID = &item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// ---------- ToMany: paymentConfigs ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["paymentConfigsIds"]; ok && !utils.IsNil(input["paymentConfigsIds"]) {
		items := []*FranchisePaymentConfig{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "FranchisePaymentConfig"); err != nil {
				return item, fmt.Errorf("FranchisePaymentConfig Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("paymentConfigsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&FranchisePaymentConfig{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&FranchisePaymentConfig{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("paymentConfigs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["paymentConfigs"]; ok && !utils.IsNil(input["paymentConfigs"]) {
		newPaymentConfigs := []*FranchisePaymentConfig{}
		updatePaymentConfigs := []*FranchisePaymentConfig{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&FranchisePaymentConfig{}).Where("organization_id = ?", item.ID).Update("organization_id", nil).Error; err != nil {
			return item, err
		}

		hasCreatePaymentConfigs := false
		hasUpdatePaymentConfigs := false

		for index, v := range changes.PaymentConfigs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "UpdateFranchisePaymentConfig"); err != nil {
						return item, fmt.Errorf("UpdateFranchisePaymentConfig: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "FranchisePaymentConfig"); err != nil {
						return item, fmt.Errorf("FranchisePaymentConfig Detail: %w", err)
					}
					hasUpdatePaymentConfigs = true
				}

				paymentConfigsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateFranchisePaymentConfig(ctx, r, paymentConfigsInput["id"].(string), paymentConfigsInput); err != nil {
					return item, fmt.Errorf("FranchisePaymentConfig ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("organization_id", item.ID).Error; err != nil {
					return item, err
				}

				updatePaymentConfigs = append(updatePaymentConfigs, v)
			} else {
				// 创建新记录
				if !hasCreatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "CreateFranchisePaymentConfig"); err != nil {
						return item, fmt.Errorf("CreateFranchisePaymentConfig: %w", err)
					}
					hasCreatePaymentConfigs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.OrganizationID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("franchise_payment_configs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newPaymentConfigs = append(newPaymentConfigs, v)
			}
		}

		allItems := append(updatePaymentConfigs, newPaymentConfigs...)

		event.AddNewValue("paymentConfigs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Organization - Delete
// ============================================================

// DeleteOrganizationFunc 执行删除或恢复操作
func DeleteOrganizationFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &Organization{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("organizations", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "Organization",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("organizations", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("organizations", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("organizations", ctx)).Updates(Organization{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteOrganizations 批量删除 Organization 实体
func (r *GeneratedMutationResolver) DeleteOrganizations(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteOrganizations(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteOrganizationsHandler 处理批量删除逻辑
func DeleteOrganizationsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteOrganizationFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Organization - Recovery
// ============================================================

// RecoveryOrganizations 批量恢复 Organization 实体
func (r *GeneratedMutationResolver) RecoveryOrganizations(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryOrganizations(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryOrganizationsHandler 处理批量恢复逻辑
func RecoveryOrganizationsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteOrganizationFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// OperatorMembership - Create
// ============================================================

// CreateOperatorMembership 创建 OperatorMembership 实体的解析器入口
func (r *GeneratedMutationResolver) CreateOperatorMembership(ctx context.Context, input map[string]interface{}) (item *OperatorMembership, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateOperatorMembership(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateOperatorMembershipHandler 处理 OperatorMembership 创建逻辑
func CreateOperatorMembershipHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *OperatorMembership, err error) {
	item = &OperatorMembership{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "OperatorMembership",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes OperatorMembershipChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["account"]) && !utils.IsNil(input["accountId"]) {
		return nil, fmt.Errorf("accountId and account cannot coexist")
	}

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	// ToMany: roles - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["roles"]) && !utils.IsNil(input["rolesIds"]) {
		return nil, fmt.Errorf("rolesIds and roles cannot coexist")
	}

	// ToMany: stores - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["stores"]) && !utils.IsNil(input["storesIds"]) {
		return nil, fmt.Errorf("storesIds and stores cannot coexist")
	}

	// ToMany: invitations - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["invitations"]) && !utils.IsNil(input["invitationsIds"]) {
		return nil, fmt.Errorf("invitationsIds and invitations cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["status"]; ok {

		item.Status = changes.Status

		event.AddNewValue("status", changes.Status)
	}

	if _, ok := input["storeAccessMode"]; ok {

		item.StoreAccessMode = changes.StoreAccessMode

		event.AddNewValue("storeAccessMode", changes.StoreAccessMode)
	}

	if _, ok := input["invitedAt"]; ok && changes.InvitedAt != nil {

		item.InvitedAt = changes.InvitedAt

		event.AddNewValue("invitedAt", changes.InvitedAt)
	}

	if _, ok := input["acceptedAt"]; ok && changes.AcceptedAt != nil {

		item.AcceptedAt = changes.AcceptedAt

		event.AddNewValue("acceptedAt", changes.AcceptedAt)
	}

	if _, ok := input["accountId"]; ok {

		if !utils.IsNil(input["accountId"]) {
			if err := tx.Select("id").Where("id = ?", input["accountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("accountId: %w", err)
			}
		}

		item.AccountID = changes.AccountID

		event.AddNewValue("accountId", changes.AccountID)
	}

	if _, ok := input["organizationId"]; ok {

		if !utils.IsNil(input["organizationId"]) {
			if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
				return nil, fmt.Errorf("organizationId: %w", err)
			}
		}

		item.OrganizationID = changes.OrganizationID

		event.AddNewValue("organizationId", changes.OrganizationID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// ---------- ToMany: roles (ManyToMany 中间表) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["rolesIds"]; ok && !utils.IsNil(input["rolesIds"]) {
		items := []*OperatorRole{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
				return item, fmt.Errorf("OperatorRole Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("rolesIds %s not found", strings.Join(differenceIds, ","))
			}

			// ManyToMany: 使用 Replace 更新中间表
			if err := tx.Model(item).Association("Roles").Replace(items); err != nil {
				return item, err
			}

		}
		event.AddNewValue("roles", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["roles"]; ok && !utils.IsNil(input["roles"]) {
		newRoles := []*OperatorRole{}
		updateRoles := []*OperatorRole{}

		hasCreateRoles := false
		hasUpdateRoles := false

		for index, v := range changes.Roles {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRoles {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorRole"); err != nil {
						return item, fmt.Errorf("UpdateOperatorRole: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
						return item, fmt.Errorf("OperatorRole Detail: %w", err)
					}
					hasUpdateRoles = true
				}

				rolesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorRole(ctx, r, rolesInput["id"].(string), rolesInput); err != nil {
					return item, fmt.Errorf("OperatorRole ID %s: %w", v.ID, err)
				}

				updateRoles = append(updateRoles, v)
			} else {
				// 创建新记录
				if !hasCreateRoles {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorRole"); err != nil {
						return item, fmt.Errorf("CreateOperatorRole: %w", err)
					}
					hasCreateRoles = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("operator_roles", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRoles = append(newRoles, v)
			}
		}

		allItems := append(updateRoles, newRoles...)

		// ManyToMany: 使用 Replace 更新中间表
		if err := tx.Model(item).Association("Roles").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("roles", allItems)
	}

	// ---------- ToMany: stores (ManyToMany 中间表) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["storesIds"]; ok && !utils.IsNil(input["storesIds"]) {
		items := []*Store{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
				return item, fmt.Errorf("Store Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("storesIds %s not found", strings.Join(differenceIds, ","))
			}

			// ManyToMany: 使用 Replace 更新中间表
			if err := tx.Model(item).Association("Stores").Replace(items); err != nil {
				return item, err
			}

		}
		event.AddNewValue("stores", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["stores"]; ok && !utils.IsNil(input["stores"]) {
		newStores := []*Store{}
		updateStores := []*Store{}

		hasCreateStores := false
		hasUpdateStores := false

		for index, v := range changes.Stores {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateStores {
					if err := auth.CheckAuthorization(ctx, "UpdateStore"); err != nil {
						return item, fmt.Errorf("UpdateStore: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
						return item, fmt.Errorf("Store Detail: %w", err)
					}
					hasUpdateStores = true
				}

				storesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStore(ctx, r, storesInput["id"].(string), storesInput); err != nil {
					return item, fmt.Errorf("Store ID %s: %w", v.ID, err)
				}

				updateStores = append(updateStores, v)
			} else {
				// 创建新记录
				if !hasCreateStores {
					if err := auth.CheckAuthorization(ctx, "CreateStore"); err != nil {
						return item, fmt.Errorf("CreateStore: %w", err)
					}
					hasCreateStores = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("stores", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newStores = append(newStores, v)
			}
		}

		allItems := append(updateStores, newStores...)

		// ManyToMany: 使用 Replace 更新中间表
		if err := tx.Model(item).Association("Stores").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("stores", allItems)
	}

	// ---------- ToMany: invitations (OneToMany 外键在 MembershipInvitation.membership_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["invitationsIds"]; ok && !utils.IsNil(input["invitationsIds"]) {
		items := []*MembershipInvitation{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
				return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("invitationsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("membership_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("invitations", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["invitations"]; ok && !utils.IsNil(input["invitations"]) {
		newInvitations := []*MembershipInvitation{}
		updateInvitations := []*MembershipInvitation{}

		hasCreateInvitations := false
		hasUpdateInvitations := false

		for index, v := range changes.Invitations {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateInvitations {
					if err := auth.CheckAuthorization(ctx, "UpdateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("UpdateMembershipInvitation: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
						return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
					}
					hasUpdateInvitations = true
				}

				invitationsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateMembershipInvitation(ctx, r, invitationsInput["id"].(string), invitationsInput); err != nil {
					return item, fmt.Errorf("MembershipInvitation ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("membership_id", item.ID).Error; err != nil {
					return item, err
				}

				updateInvitations = append(updateInvitations, v)
			} else {
				// 创建新记录
				if !hasCreateInvitations {
					if err := auth.CheckAuthorization(ctx, "CreateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("CreateMembershipInvitation: %w", err)
					}
					hasCreateInvitations = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.MembershipID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("membership_invitations", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newInvitations = append(newInvitations, v)
			}
		}

		allItems := append(updateInvitations, newInvitations...)

		event.AddNewValue("invitations", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// OperatorMembership - Update
// ============================================================

// UpdateOperatorMembership 更新 OperatorMembership 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateOperatorMembership(ctx context.Context, id string, input map[string]interface{}) (item *OperatorMembership, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateOperatorMembership(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateOperatorMembershipHandler 处理 OperatorMembership 更新逻辑
func UpdateOperatorMembershipHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *OperatorMembership, err error) {
	item = &OperatorMembership{}
	newItem := &OperatorMembership{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "OperatorMembership",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes OperatorMembershipChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["account"]) && !utils.IsNil(input["accountId"]) {
		return nil, fmt.Errorf("accountId and account cannot coexist")
	}

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["roles"]) && !utils.IsNil(input["rolesIds"]) {
		return nil, fmt.Errorf("rolesIds and roles cannot coexist")
	}

	if !utils.IsNil(input["stores"]) && !utils.IsNil(input["storesIds"]) {
		return nil, fmt.Errorf("storesIds and stores cannot coexist")
	}

	if !utils.IsNil(input["invitations"]) && !utils.IsNil(input["invitationsIds"]) {
		return nil, fmt.Errorf("invitationsIds and invitations cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("operator_memberships", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["status"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Status != changes.Status {

			event.AddOldValue("status", item.Status)
			event.AddNewValue("status", changes.Status)

			item.Status = changes.Status
			newItem.Status = changes.Status
			changedFields = append(changedFields, "status")
			isChange = true
		}
	}

	if _, ok := input["storeAccessMode"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.StoreAccessMode != changes.StoreAccessMode {

			event.AddOldValue("storeAccessMode", item.StoreAccessMode)
			event.AddNewValue("storeAccessMode", changes.StoreAccessMode)

			item.StoreAccessMode = changes.StoreAccessMode
			newItem.StoreAccessMode = changes.StoreAccessMode
			changedFields = append(changedFields, "store_access_mode")
			isChange = true
		}
	}

	if _, ok := input["invitedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.InvitedAt != changes.InvitedAt) && (item.InvitedAt == nil || changes.InvitedAt == nil || *item.InvitedAt != *changes.InvitedAt) {

			event.AddOldValue("invitedAt", item.InvitedAt)
			event.AddNewValue("invitedAt", changes.InvitedAt)

			item.InvitedAt = changes.InvitedAt
			newItem.InvitedAt = changes.InvitedAt
			changedFields = append(changedFields, "invited_at")
			isChange = true
		}
	}

	if _, ok := input["acceptedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.AcceptedAt != changes.AcceptedAt) && (item.AcceptedAt == nil || changes.AcceptedAt == nil || *item.AcceptedAt != *changes.AcceptedAt) {

			event.AddOldValue("acceptedAt", item.AcceptedAt)
			event.AddNewValue("acceptedAt", changes.AcceptedAt)

			item.AcceptedAt = changes.AcceptedAt
			newItem.AcceptedAt = changes.AcceptedAt
			changedFields = append(changedFields, "accepted_at")
			isChange = true
		}
	}

	if _, ok := input["accountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.AccountID != changes.AccountID {

			if !utils.IsNil(input["accountId"]) {
				if err := tx.Select("id").Where("id = ?", input["accountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("accountId: %w", err)
				}
			}

			event.AddOldValue("accountId", item.AccountID)
			event.AddNewValue("accountId", changes.AccountID)

			item.AccountID = changes.AccountID
			newItem.AccountID = changes.AccountID
			changedFields = append(changedFields, "account_id")
			isChange = true
		}
	}

	if _, ok := input["organizationId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.OrganizationID != changes.OrganizationID {

			if !utils.IsNil(input["organizationId"]) {
				if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
					return nil, fmt.Errorf("organizationId: %w", err)
				}
			}

			event.AddOldValue("organizationId", item.OrganizationID)
			event.AddNewValue("organizationId", changes.OrganizationID)

			item.OrganizationID = changes.OrganizationID
			newItem.OrganizationID = changes.OrganizationID
			changedFields = append(changedFields, "organization_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("operator_memberships", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// ---------- ToMany: roles ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["rolesIds"]; ok && !utils.IsNil(input["rolesIds"]) {
		items := []*OperatorRole{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
				return item, fmt.Errorf("OperatorRole Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("rolesIds %s not found", strings.Join(differenceIds, ","))
			}

			if err := tx.Model(item).Association("Roles").Replace(items); err != nil {
				return item, err
			}

		} else {
			// 清空关联

			if err := tx.Model(item).Association("Roles").Clear(); err != nil {
				return item, err
			}

		}
		event.AddNewValue("roles", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["roles"]; ok && !utils.IsNil(input["roles"]) {
		newRoles := []*OperatorRole{}
		updateRoles := []*OperatorRole{}

		hasCreateRoles := false
		hasUpdateRoles := false

		for index, v := range changes.Roles {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRoles {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorRole"); err != nil {
						return item, fmt.Errorf("UpdateOperatorRole: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
						return item, fmt.Errorf("OperatorRole Detail: %w", err)
					}
					hasUpdateRoles = true
				}

				rolesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorRole(ctx, r, rolesInput["id"].(string), rolesInput); err != nil {
					return item, fmt.Errorf("OperatorRole ID %s: %w", v.ID, err)
				}

				updateRoles = append(updateRoles, v)
			} else {
				// 创建新记录
				if !hasCreateRoles {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorRole"); err != nil {
						return item, fmt.Errorf("CreateOperatorRole: %w", err)
					}
					hasCreateRoles = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				if err := tx.Omit(clause.Associations).Table(TableName("operator_roles", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRoles = append(newRoles, v)
			}
		}

		allItems := append(updateRoles, newRoles...)

		if err := tx.Model(item).Association("Roles").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("roles", allItems)
	}

	// ---------- ToMany: stores ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["storesIds"]; ok && !utils.IsNil(input["storesIds"]) {
		items := []*Store{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
				return item, fmt.Errorf("Store Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("storesIds %s not found", strings.Join(differenceIds, ","))
			}

			if err := tx.Model(item).Association("Stores").Replace(items); err != nil {
				return item, err
			}

		} else {
			// 清空关联

			if err := tx.Model(item).Association("Stores").Clear(); err != nil {
				return item, err
			}

		}
		event.AddNewValue("stores", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["stores"]; ok && !utils.IsNil(input["stores"]) {
		newStores := []*Store{}
		updateStores := []*Store{}

		hasCreateStores := false
		hasUpdateStores := false

		for index, v := range changes.Stores {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateStores {
					if err := auth.CheckAuthorization(ctx, "UpdateStore"); err != nil {
						return item, fmt.Errorf("UpdateStore: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Store"); err != nil {
						return item, fmt.Errorf("Store Detail: %w", err)
					}
					hasUpdateStores = true
				}

				storesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStore(ctx, r, storesInput["id"].(string), storesInput); err != nil {
					return item, fmt.Errorf("Store ID %s: %w", v.ID, err)
				}

				updateStores = append(updateStores, v)
			} else {
				// 创建新记录
				if !hasCreateStores {
					if err := auth.CheckAuthorization(ctx, "CreateStore"); err != nil {
						return item, fmt.Errorf("CreateStore: %w", err)
					}
					hasCreateStores = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				if err := tx.Omit(clause.Associations).Table(TableName("stores", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newStores = append(newStores, v)
			}
		}

		allItems := append(updateStores, newStores...)

		if err := tx.Model(item).Association("Stores").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("stores", allItems)
	}

	// ---------- ToMany: invitations ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["invitationsIds"]; ok && !utils.IsNil(input["invitationsIds"]) {
		items := []*MembershipInvitation{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
				return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("invitationsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&MembershipInvitation{}).Where("membership_id = ?", item.ID).Update("membership_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("membership_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&MembershipInvitation{}).Where("membership_id = ?", item.ID).Update("membership_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("invitations", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["invitations"]; ok && !utils.IsNil(input["invitations"]) {
		newInvitations := []*MembershipInvitation{}
		updateInvitations := []*MembershipInvitation{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&MembershipInvitation{}).Where("membership_id = ?", item.ID).Update("membership_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateInvitations := false
		hasUpdateInvitations := false

		for index, v := range changes.Invitations {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateInvitations {
					if err := auth.CheckAuthorization(ctx, "UpdateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("UpdateMembershipInvitation: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "MembershipInvitation"); err != nil {
						return item, fmt.Errorf("MembershipInvitation Detail: %w", err)
					}
					hasUpdateInvitations = true
				}

				invitationsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateMembershipInvitation(ctx, r, invitationsInput["id"].(string), invitationsInput); err != nil {
					return item, fmt.Errorf("MembershipInvitation ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("membership_id", item.ID).Error; err != nil {
					return item, err
				}

				updateInvitations = append(updateInvitations, v)
			} else {
				// 创建新记录
				if !hasCreateInvitations {
					if err := auth.CheckAuthorization(ctx, "CreateMembershipInvitation"); err != nil {
						return item, fmt.Errorf("CreateMembershipInvitation: %w", err)
					}
					hasCreateInvitations = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.MembershipID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("membership_invitations", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newInvitations = append(newInvitations, v)
			}
		}

		allItems := append(updateInvitations, newInvitations...)

		event.AddNewValue("invitations", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// OperatorMembership - Delete
// ============================================================

// DeleteOperatorMembershipFunc 执行删除或恢复操作
func DeleteOperatorMembershipFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &OperatorMembership{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("operator_memberships", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "OperatorMembership",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("operator_memberships", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("operator_memberships", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("operator_memberships", ctx)).Updates(OperatorMembership{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteOperatorMemberships 批量删除 OperatorMembership 实体
func (r *GeneratedMutationResolver) DeleteOperatorMemberships(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteOperatorMemberships(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteOperatorMembershipsHandler 处理批量删除逻辑
func DeleteOperatorMembershipsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteOperatorMembershipFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// OperatorMembership - Recovery
// ============================================================

// RecoveryOperatorMemberships 批量恢复 OperatorMembership 实体
func (r *GeneratedMutationResolver) RecoveryOperatorMemberships(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryOperatorMemberships(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryOperatorMembershipsHandler 处理批量恢复逻辑
func RecoveryOperatorMembershipsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteOperatorMembershipFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Permission - Create
// ============================================================

// CreatePermission 创建 Permission 实体的解析器入口
func (r *GeneratedMutationResolver) CreatePermission(ctx context.Context, input map[string]interface{}) (item *Permission, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreatePermission(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreatePermissionHandler 处理 Permission 创建逻辑
func CreatePermissionHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Permission, err error) {
	item = &Permission{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "Permission",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes PermissionChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	// ToMany: roles - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["roles"]) && !utils.IsNil(input["rolesIds"]) {
		return nil, fmt.Errorf("rolesIds and roles cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["name"]; ok {

		item.Name = changes.Name

		event.AddNewValue("name", changes.Name)
	}

	if _, ok := input["action"]; ok {

		item.Action = changes.Action

		event.AddNewValue("action", changes.Action)
	}

	if _, ok := input["module"]; ok {

		item.Module = changes.Module

		event.AddNewValue("module", changes.Module)
	}

	if _, ok := input["scope"]; ok {

		item.Scope = changes.Scope

		event.AddNewValue("scope", changes.Scope)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("permissions", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// ---------- ToMany: roles (ManyToMany 中间表) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["rolesIds"]; ok && !utils.IsNil(input["rolesIds"]) {
		items := []*OperatorRole{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
				return item, fmt.Errorf("OperatorRole Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("rolesIds %s not found", strings.Join(differenceIds, ","))
			}

			// ManyToMany: 使用 Replace 更新中间表
			if err := tx.Model(item).Association("Roles").Replace(items); err != nil {
				return item, err
			}

		}
		event.AddNewValue("roles", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["roles"]; ok && !utils.IsNil(input["roles"]) {
		newRoles := []*OperatorRole{}
		updateRoles := []*OperatorRole{}

		hasCreateRoles := false
		hasUpdateRoles := false

		for index, v := range changes.Roles {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRoles {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorRole"); err != nil {
						return item, fmt.Errorf("UpdateOperatorRole: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
						return item, fmt.Errorf("OperatorRole Detail: %w", err)
					}
					hasUpdateRoles = true
				}

				rolesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorRole(ctx, r, rolesInput["id"].(string), rolesInput); err != nil {
					return item, fmt.Errorf("OperatorRole ID %s: %w", v.ID, err)
				}

				updateRoles = append(updateRoles, v)
			} else {
				// 创建新记录
				if !hasCreateRoles {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorRole"); err != nil {
						return item, fmt.Errorf("CreateOperatorRole: %w", err)
					}
					hasCreateRoles = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("operator_roles", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRoles = append(newRoles, v)
			}
		}

		allItems := append(updateRoles, newRoles...)

		// ManyToMany: 使用 Replace 更新中间表
		if err := tx.Model(item).Association("Roles").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("roles", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Permission - Update
// ============================================================

// UpdatePermission 更新 Permission 实体的解析器入口
func (r *GeneratedMutationResolver) UpdatePermission(ctx context.Context, id string, input map[string]interface{}) (item *Permission, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdatePermission(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdatePermissionHandler 处理 Permission 更新逻辑
func UpdatePermissionHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Permission, err error) {
	item = &Permission{}
	newItem := &Permission{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "Permission",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes PermissionChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["roles"]) && !utils.IsNil(input["rolesIds"]) {
		return nil, fmt.Errorf("rolesIds and roles cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("permissions", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["name"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Name != changes.Name {

			event.AddOldValue("name", item.Name)
			event.AddNewValue("name", changes.Name)

			item.Name = changes.Name
			newItem.Name = changes.Name
			changedFields = append(changedFields, "name")
			isChange = true
		}
	}

	if _, ok := input["action"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Action != changes.Action {

			event.AddOldValue("action", item.Action)
			event.AddNewValue("action", changes.Action)

			item.Action = changes.Action
			newItem.Action = changes.Action
			changedFields = append(changedFields, "action")
			isChange = true
		}
	}

	if _, ok := input["module"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Module != changes.Module {

			event.AddOldValue("module", item.Module)
			event.AddNewValue("module", changes.Module)

			item.Module = changes.Module
			newItem.Module = changes.Module
			changedFields = append(changedFields, "module")
			isChange = true
		}
	}

	if _, ok := input["scope"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Scope != changes.Scope {

			event.AddOldValue("scope", item.Scope)
			event.AddNewValue("scope", changes.Scope)

			item.Scope = changes.Scope
			newItem.Scope = changes.Scope
			changedFields = append(changedFields, "scope")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("permissions", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// ---------- ToMany: roles ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["rolesIds"]; ok && !utils.IsNil(input["rolesIds"]) {
		items := []*OperatorRole{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
				return item, fmt.Errorf("OperatorRole Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("rolesIds %s not found", strings.Join(differenceIds, ","))
			}

			if err := tx.Model(item).Association("Roles").Replace(items); err != nil {
				return item, err
			}

		} else {
			// 清空关联

			if err := tx.Model(item).Association("Roles").Clear(); err != nil {
				return item, err
			}

		}
		event.AddNewValue("roles", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["roles"]; ok && !utils.IsNil(input["roles"]) {
		newRoles := []*OperatorRole{}
		updateRoles := []*OperatorRole{}

		hasCreateRoles := false
		hasUpdateRoles := false

		for index, v := range changes.Roles {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateRoles {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorRole"); err != nil {
						return item, fmt.Errorf("UpdateOperatorRole: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorRole"); err != nil {
						return item, fmt.Errorf("OperatorRole Detail: %w", err)
					}
					hasUpdateRoles = true
				}

				rolesInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorRole(ctx, r, rolesInput["id"].(string), rolesInput); err != nil {
					return item, fmt.Errorf("OperatorRole ID %s: %w", v.ID, err)
				}

				updateRoles = append(updateRoles, v)
			} else {
				// 创建新记录
				if !hasCreateRoles {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorRole"); err != nil {
						return item, fmt.Errorf("CreateOperatorRole: %w", err)
					}
					hasCreateRoles = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				if err := tx.Omit(clause.Associations).Table(TableName("operator_roles", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newRoles = append(newRoles, v)
			}
		}

		allItems := append(updateRoles, newRoles...)

		if err := tx.Model(item).Association("Roles").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("roles", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Permission - Delete
// ============================================================

// DeletePermissionFunc 执行删除或恢复操作
func DeletePermissionFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &Permission{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("permissions", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "Permission",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("permissions", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("permissions", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("permissions", ctx)).Updates(Permission{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeletePermissions 批量删除 Permission 实体
func (r *GeneratedMutationResolver) DeletePermissions(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeletePermissions(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeletePermissionsHandler 处理批量删除逻辑
func DeletePermissionsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeletePermissionFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Permission - Recovery
// ============================================================

// RecoveryPermissions 批量恢复 Permission 实体
func (r *GeneratedMutationResolver) RecoveryPermissions(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryPermissions(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryPermissionsHandler 处理批量恢复逻辑
func RecoveryPermissionsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeletePermissionFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// OperatorRole - Create
// ============================================================

// CreateOperatorRole 创建 OperatorRole 实体的解析器入口
func (r *GeneratedMutationResolver) CreateOperatorRole(ctx context.Context, input map[string]interface{}) (item *OperatorRole, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateOperatorRole(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateOperatorRoleHandler 处理 OperatorRole 创建逻辑
func CreateOperatorRoleHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *OperatorRole, err error) {
	item = &OperatorRole{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "OperatorRole",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes OperatorRoleChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	// ToMany: members - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["members"]) && !utils.IsNil(input["membersIds"]) {
		return nil, fmt.Errorf("membersIds and members cannot coexist")
	}

	// ToMany: permissions - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["permissions"]) && !utils.IsNil(input["permissionsIds"]) {
		return nil, fmt.Errorf("permissionsIds and permissions cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["name"]; ok {

		item.Name = changes.Name

		event.AddNewValue("name", changes.Name)
	}

	if _, ok := input["kind"]; ok {

		item.Kind = changes.Kind

		event.AddNewValue("kind", changes.Kind)
	}

	if _, ok := input["organizationId"]; ok {

		if !utils.IsNil(input["organizationId"]) {
			if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
				return nil, fmt.Errorf("organizationId: %w", err)
			}
		}

		item.OrganizationID = changes.OrganizationID

		event.AddNewValue("organizationId", changes.OrganizationID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("operator_roles", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// ---------- ToMany: members (ManyToMany 中间表) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["membersIds"]; ok && !utils.IsNil(input["membersIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membersIds %s not found", strings.Join(differenceIds, ","))
			}

			// ManyToMany: 使用 Replace 更新中间表
			if err := tx.Model(item).Association("Members").Replace(items); err != nil {
				return item, err
			}

		}
		event.AddNewValue("members", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["members"]; ok && !utils.IsNil(input["members"]) {
		newMembers := []*OperatorMembership{}
		updateMembers := []*OperatorMembership{}

		hasCreateMembers := false
		hasUpdateMembers := false

		for index, v := range changes.Members {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMembers {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMembers = true
				}

				membersInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membersInput["id"].(string), membersInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				updateMembers = append(updateMembers, v)
			} else {
				// 创建新记录
				if !hasCreateMembers {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMembers = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMembers = append(newMembers, v)
			}
		}

		allItems := append(updateMembers, newMembers...)

		// ManyToMany: 使用 Replace 更新中间表
		if err := tx.Model(item).Association("Members").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("members", allItems)
	}

	// ---------- ToMany: permissions (ManyToMany 中间表) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["permissionsIds"]; ok && !utils.IsNil(input["permissionsIds"]) {
		items := []*Permission{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "Permission"); err != nil {
				return item, fmt.Errorf("Permission Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("permissionsIds %s not found", strings.Join(differenceIds, ","))
			}

			// ManyToMany: 使用 Replace 更新中间表
			if err := tx.Model(item).Association("Permissions").Replace(items); err != nil {
				return item, err
			}

		}
		event.AddNewValue("permissions", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["permissions"]; ok && !utils.IsNil(input["permissions"]) {
		newPermissions := []*Permission{}
		updatePermissions := []*Permission{}

		hasCreatePermissions := false
		hasUpdatePermissions := false

		for index, v := range changes.Permissions {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdatePermissions {
					if err := auth.CheckAuthorization(ctx, "UpdatePermission"); err != nil {
						return item, fmt.Errorf("UpdatePermission: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Permission"); err != nil {
						return item, fmt.Errorf("Permission Detail: %w", err)
					}
					hasUpdatePermissions = true
				}

				permissionsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdatePermission(ctx, r, permissionsInput["id"].(string), permissionsInput); err != nil {
					return item, fmt.Errorf("Permission ID %s: %w", v.ID, err)
				}

				updatePermissions = append(updatePermissions, v)
			} else {
				// 创建新记录
				if !hasCreatePermissions {
					if err := auth.CheckAuthorization(ctx, "CreatePermission"); err != nil {
						return item, fmt.Errorf("CreatePermission: %w", err)
					}
					hasCreatePermissions = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("permissions", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newPermissions = append(newPermissions, v)
			}
		}

		allItems := append(updatePermissions, newPermissions...)

		// ManyToMany: 使用 Replace 更新中间表
		if err := tx.Model(item).Association("Permissions").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("permissions", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// OperatorRole - Update
// ============================================================

// UpdateOperatorRole 更新 OperatorRole 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateOperatorRole(ctx context.Context, id string, input map[string]interface{}) (item *OperatorRole, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateOperatorRole(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateOperatorRoleHandler 处理 OperatorRole 更新逻辑
func UpdateOperatorRoleHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *OperatorRole, err error) {
	item = &OperatorRole{}
	newItem := &OperatorRole{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "OperatorRole",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes OperatorRoleChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["members"]) && !utils.IsNil(input["membersIds"]) {
		return nil, fmt.Errorf("membersIds and members cannot coexist")
	}

	if !utils.IsNil(input["permissions"]) && !utils.IsNil(input["permissionsIds"]) {
		return nil, fmt.Errorf("permissionsIds and permissions cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("operator_roles", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["name"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Name != changes.Name {

			event.AddOldValue("name", item.Name)
			event.AddNewValue("name", changes.Name)

			item.Name = changes.Name
			newItem.Name = changes.Name
			changedFields = append(changedFields, "name")
			isChange = true
		}
	}

	if _, ok := input["kind"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Kind != changes.Kind {

			event.AddOldValue("kind", item.Kind)
			event.AddNewValue("kind", changes.Kind)

			item.Kind = changes.Kind
			newItem.Kind = changes.Kind
			changedFields = append(changedFields, "kind")
			isChange = true
		}
	}

	if _, ok := input["organizationId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.OrganizationID != changes.OrganizationID {

			if !utils.IsNil(input["organizationId"]) {
				if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
					return nil, fmt.Errorf("organizationId: %w", err)
				}
			}

			event.AddOldValue("organizationId", item.OrganizationID)
			event.AddNewValue("organizationId", changes.OrganizationID)

			item.OrganizationID = changes.OrganizationID
			newItem.OrganizationID = changes.OrganizationID
			changedFields = append(changedFields, "organization_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("operator_roles", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// ---------- ToMany: members ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["membersIds"]; ok && !utils.IsNil(input["membersIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membersIds %s not found", strings.Join(differenceIds, ","))
			}

			if err := tx.Model(item).Association("Members").Replace(items); err != nil {
				return item, err
			}

		} else {
			// 清空关联

			if err := tx.Model(item).Association("Members").Clear(); err != nil {
				return item, err
			}

		}
		event.AddNewValue("members", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["members"]; ok && !utils.IsNil(input["members"]) {
		newMembers := []*OperatorMembership{}
		updateMembers := []*OperatorMembership{}

		hasCreateMembers := false
		hasUpdateMembers := false

		for index, v := range changes.Members {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMembers {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMembers = true
				}

				membersInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membersInput["id"].(string), membersInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				updateMembers = append(updateMembers, v)
			} else {
				// 创建新记录
				if !hasCreateMembers {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMembers = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMembers = append(newMembers, v)
			}
		}

		allItems := append(updateMembers, newMembers...)

		if err := tx.Model(item).Association("Members").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("members", allItems)
	}

	// ---------- ToMany: permissions ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["permissionsIds"]; ok && !utils.IsNil(input["permissionsIds"]) {
		items := []*Permission{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "Permission"); err != nil {
				return item, fmt.Errorf("Permission Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("permissionsIds %s not found", strings.Join(differenceIds, ","))
			}

			if err := tx.Model(item).Association("Permissions").Replace(items); err != nil {
				return item, err
			}

		} else {
			// 清空关联

			if err := tx.Model(item).Association("Permissions").Clear(); err != nil {
				return item, err
			}

		}
		event.AddNewValue("permissions", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["permissions"]; ok && !utils.IsNil(input["permissions"]) {
		newPermissions := []*Permission{}
		updatePermissions := []*Permission{}

		hasCreatePermissions := false
		hasUpdatePermissions := false

		for index, v := range changes.Permissions {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdatePermissions {
					if err := auth.CheckAuthorization(ctx, "UpdatePermission"); err != nil {
						return item, fmt.Errorf("UpdatePermission: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "Permission"); err != nil {
						return item, fmt.Errorf("Permission Detail: %w", err)
					}
					hasUpdatePermissions = true
				}

				permissionsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdatePermission(ctx, r, permissionsInput["id"].(string), permissionsInput); err != nil {
					return item, fmt.Errorf("Permission ID %s: %w", v.ID, err)
				}

				updatePermissions = append(updatePermissions, v)
			} else {
				// 创建新记录
				if !hasCreatePermissions {
					if err := auth.CheckAuthorization(ctx, "CreatePermission"); err != nil {
						return item, fmt.Errorf("CreatePermission: %w", err)
					}
					hasCreatePermissions = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				if err := tx.Omit(clause.Associations).Table(TableName("permissions", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newPermissions = append(newPermissions, v)
			}
		}

		allItems := append(updatePermissions, newPermissions...)

		if err := tx.Model(item).Association("Permissions").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("permissions", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// OperatorRole - Delete
// ============================================================

// DeleteOperatorRoleFunc 执行删除或恢复操作
func DeleteOperatorRoleFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &OperatorRole{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("operator_roles", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "OperatorRole",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("operator_roles", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("operator_roles", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("operator_roles", ctx)).Updates(OperatorRole{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteOperatorRoles 批量删除 OperatorRole 实体
func (r *GeneratedMutationResolver) DeleteOperatorRoles(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteOperatorRoles(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteOperatorRolesHandler 处理批量删除逻辑
func DeleteOperatorRolesHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteOperatorRoleFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// OperatorRole - Recovery
// ============================================================

// RecoveryOperatorRoles 批量恢复 OperatorRole 实体
func (r *GeneratedMutationResolver) RecoveryOperatorRoles(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryOperatorRoles(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryOperatorRolesHandler 处理批量恢复逻辑
func RecoveryOperatorRolesHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteOperatorRoleFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Store - Create
// ============================================================

// CreateStore 创建 Store 实体的解析器入口
func (r *GeneratedMutationResolver) CreateStore(ctx context.Context, input map[string]interface{}) (item *Store, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateStore(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateStoreHandler 处理 Store 创建逻辑
func CreateStoreHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Store, err error) {
	item = &Store{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "Store",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes StoreChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	// ToMany: members - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["members"]) && !utils.IsNil(input["membersIds"]) {
		return nil, fmt.Errorf("membersIds and members cannot coexist")
	}

	if !utils.IsNil(input["reviewedByAccount"]) && !utils.IsNil(input["reviewedByAccountId"]) {
		return nil, fmt.Errorf("reviewedByAccountId and reviewedByAccount cannot coexist")
	}

	// ToMany: auditLogs - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	// ToMany: paymentConfigs - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["paymentConfigs"]) && !utils.IsNil(input["paymentConfigsIds"]) {
		return nil, fmt.Errorf("paymentConfigsIds and paymentConfigs cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["code"]; ok {

		item.Code = changes.Code

		event.AddNewValue("code", changes.Code)
	}

	if _, ok := input["name"]; ok {

		item.Name = changes.Name

		event.AddNewValue("name", changes.Name)
	}

	if _, ok := input["lifecycle"]; ok {

		item.Lifecycle = changes.Lifecycle

		event.AddNewValue("lifecycle", changes.Lifecycle)
	}

	if _, ok := input["submittedAt"]; ok && changes.SubmittedAt != nil {

		item.SubmittedAt = changes.SubmittedAt

		event.AddNewValue("submittedAt", changes.SubmittedAt)
	}

	if _, ok := input["reviewedAt"]; ok && changes.ReviewedAt != nil {

		item.ReviewedAt = changes.ReviewedAt

		event.AddNewValue("reviewedAt", changes.ReviewedAt)
	}

	if _, ok := input["rejectionReason"]; ok && changes.RejectionReason != nil {

		item.RejectionReason = changes.RejectionReason

		event.AddNewValue("rejectionReason", changes.RejectionReason)
	}

	if _, ok := input["contactPhone"]; ok && changes.ContactPhone != nil {

		item.ContactPhone = changes.ContactPhone

		event.AddNewValue("contactPhone", changes.ContactPhone)
	}

	if _, ok := input["managerName"]; ok && changes.ManagerName != nil {

		item.ManagerName = changes.ManagerName

		event.AddNewValue("managerName", changes.ManagerName)
	}

	if _, ok := input["managerPhone"]; ok && changes.ManagerPhone != nil {

		item.ManagerPhone = changes.ManagerPhone

		event.AddNewValue("managerPhone", changes.ManagerPhone)
	}

	if _, ok := input["province"]; ok && changes.Province != nil {

		item.Province = changes.Province

		event.AddNewValue("province", changes.Province)
	}

	if _, ok := input["city"]; ok && changes.City != nil {

		item.City = changes.City

		event.AddNewValue("city", changes.City)
	}

	if _, ok := input["district"]; ok && changes.District != nil {

		item.District = changes.District

		event.AddNewValue("district", changes.District)
	}

	if _, ok := input["address"]; ok && changes.Address != nil {

		item.Address = changes.Address

		event.AddNewValue("address", changes.Address)
	}

	if _, ok := input["businessHours"]; ok && changes.BusinessHours != nil {

		item.BusinessHours = changes.BusinessHours

		event.AddNewValue("businessHours", changes.BusinessHours)
	}

	if _, ok := input["businessStatus"]; ok && changes.BusinessStatus != nil {

		item.BusinessStatus = changes.BusinessStatus

		event.AddNewValue("businessStatus", changes.BusinessStatus)
	}

	if _, ok := input["supportDineIn"]; ok && changes.SupportDineIn != nil {

		item.SupportDineIn = changes.SupportDineIn

		event.AddNewValue("supportDineIn", changes.SupportDineIn)
	}

	if _, ok := input["supportTakeout"]; ok && changes.SupportTakeout != nil {

		item.SupportTakeout = changes.SupportTakeout

		event.AddNewValue("supportTakeout", changes.SupportTakeout)
	}

	if _, ok := input["storeArea"]; ok && changes.StoreArea != nil {

		item.StoreArea = changes.StoreArea

		event.AddNewValue("storeArea", changes.StoreArea)
	}

	if _, ok := input["tableCount"]; ok && changes.TableCount != nil {

		item.TableCount = changes.TableCount

		event.AddNewValue("tableCount", changes.TableCount)
	}

	if _, ok := input["receiptFooter"]; ok && changes.ReceiptFooter != nil {

		item.ReceiptFooter = changes.ReceiptFooter

		event.AddNewValue("receiptFooter", changes.ReceiptFooter)
	}

	if _, ok := input["businessLicenseImageUrl"]; ok && changes.BusinessLicenseImageURL != nil {

		item.BusinessLicenseImageURL = changes.BusinessLicenseImageURL

		event.AddNewValue("businessLicenseImageUrl", changes.BusinessLicenseImageURL)
	}

	if _, ok := input["otherDocumentImageUrl"]; ok && changes.OtherDocumentImageURL != nil {

		item.OtherDocumentImageURL = changes.OtherDocumentImageURL

		event.AddNewValue("otherDocumentImageUrl", changes.OtherDocumentImageURL)
	}

	if _, ok := input["organizationId"]; ok {

		if !utils.IsNil(input["organizationId"]) {
			if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
				return nil, fmt.Errorf("organizationId: %w", err)
			}
		}

		item.OrganizationID = changes.OrganizationID

		event.AddNewValue("organizationId", changes.OrganizationID)
	}

	if _, ok := input["reviewedByAccountId"]; ok && changes.ReviewedByAccountID != nil {

		if !utils.IsNil(input["reviewedByAccountId"]) {
			if err := tx.Select("id").Where("id = ?", input["reviewedByAccountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("reviewedByAccountId: %w", err)
			}
		}

		item.ReviewedByAccountID = changes.ReviewedByAccountID

		event.AddNewValue("reviewedByAccountId", changes.ReviewedByAccountID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("stores", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// ---------- ToMany: members (ManyToMany 中间表) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["membersIds"]; ok && !utils.IsNil(input["membersIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membersIds %s not found", strings.Join(differenceIds, ","))
			}

			// ManyToMany: 使用 Replace 更新中间表
			if err := tx.Model(item).Association("Members").Replace(items); err != nil {
				return item, err
			}

		}
		event.AddNewValue("members", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["members"]; ok && !utils.IsNil(input["members"]) {
		newMembers := []*OperatorMembership{}
		updateMembers := []*OperatorMembership{}

		hasCreateMembers := false
		hasUpdateMembers := false

		for index, v := range changes.Members {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMembers {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMembers = true
				}

				membersInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membersInput["id"].(string), membersInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				updateMembers = append(updateMembers, v)
			} else {
				// 创建新记录
				if !hasCreateMembers {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMembers = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMembers = append(newMembers, v)
			}
		}

		allItems := append(updateMembers, newMembers...)

		// ManyToMany: 使用 Replace 更新中间表
		if err := tx.Model(item).Association("Members").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("members", allItems)
	}

	// ---------- ToMany: auditLogs (OneToMany 外键在 AuditLog.store_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.StoreID = &item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// ---------- ToMany: paymentConfigs (OneToMany 外键在 StorePaymentConfig.store_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["paymentConfigsIds"]; ok && !utils.IsNil(input["paymentConfigsIds"]) {
		items := []*StorePaymentConfig{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "StorePaymentConfig"); err != nil {
				return item, fmt.Errorf("StorePaymentConfig Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("paymentConfigsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("paymentConfigs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["paymentConfigs"]; ok && !utils.IsNil(input["paymentConfigs"]) {
		newPaymentConfigs := []*StorePaymentConfig{}
		updatePaymentConfigs := []*StorePaymentConfig{}

		hasCreatePaymentConfigs := false
		hasUpdatePaymentConfigs := false

		for index, v := range changes.PaymentConfigs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "UpdateStorePaymentConfig"); err != nil {
						return item, fmt.Errorf("UpdateStorePaymentConfig: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "StorePaymentConfig"); err != nil {
						return item, fmt.Errorf("StorePaymentConfig Detail: %w", err)
					}
					hasUpdatePaymentConfigs = true
				}

				paymentConfigsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStorePaymentConfig(ctx, r, paymentConfigsInput["id"].(string), paymentConfigsInput); err != nil {
					return item, fmt.Errorf("StorePaymentConfig ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}

				updatePaymentConfigs = append(updatePaymentConfigs, v)
			} else {
				// 创建新记录
				if !hasCreatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "CreateStorePaymentConfig"); err != nil {
						return item, fmt.Errorf("CreateStorePaymentConfig: %w", err)
					}
					hasCreatePaymentConfigs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.StoreID = item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("store_payment_configs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newPaymentConfigs = append(newPaymentConfigs, v)
			}
		}

		allItems := append(updatePaymentConfigs, newPaymentConfigs...)

		event.AddNewValue("paymentConfigs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Store - Update
// ============================================================

// UpdateStore 更新 Store 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateStore(ctx context.Context, id string, input map[string]interface{}) (item *Store, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateStore(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateStoreHandler 处理 Store 更新逻辑
func UpdateStoreHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Store, err error) {
	item = &Store{}
	newItem := &Store{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "Store",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes StoreChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["members"]) && !utils.IsNil(input["membersIds"]) {
		return nil, fmt.Errorf("membersIds and members cannot coexist")
	}

	if !utils.IsNil(input["reviewedByAccount"]) && !utils.IsNil(input["reviewedByAccountId"]) {
		return nil, fmt.Errorf("reviewedByAccountId and reviewedByAccount cannot coexist")
	}

	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	if !utils.IsNil(input["paymentConfigs"]) && !utils.IsNil(input["paymentConfigsIds"]) {
		return nil, fmt.Errorf("paymentConfigsIds and paymentConfigs cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("stores", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["code"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Code != changes.Code {

			event.AddOldValue("code", item.Code)
			event.AddNewValue("code", changes.Code)

			item.Code = changes.Code
			newItem.Code = changes.Code
			changedFields = append(changedFields, "code")
			isChange = true
		}
	}

	if _, ok := input["name"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Name != changes.Name {

			event.AddOldValue("name", item.Name)
			event.AddNewValue("name", changes.Name)

			item.Name = changes.Name
			newItem.Name = changes.Name
			changedFields = append(changedFields, "name")
			isChange = true
		}
	}

	if _, ok := input["lifecycle"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Lifecycle != changes.Lifecycle {

			event.AddOldValue("lifecycle", item.Lifecycle)
			event.AddNewValue("lifecycle", changes.Lifecycle)

			item.Lifecycle = changes.Lifecycle
			newItem.Lifecycle = changes.Lifecycle
			changedFields = append(changedFields, "lifecycle")
			isChange = true
		}
	}

	if _, ok := input["submittedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.SubmittedAt != changes.SubmittedAt) && (item.SubmittedAt == nil || changes.SubmittedAt == nil || *item.SubmittedAt != *changes.SubmittedAt) {

			event.AddOldValue("submittedAt", item.SubmittedAt)
			event.AddNewValue("submittedAt", changes.SubmittedAt)

			item.SubmittedAt = changes.SubmittedAt
			newItem.SubmittedAt = changes.SubmittedAt
			changedFields = append(changedFields, "submitted_at")
			isChange = true
		}
	}

	if _, ok := input["reviewedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ReviewedAt != changes.ReviewedAt) && (item.ReviewedAt == nil || changes.ReviewedAt == nil || *item.ReviewedAt != *changes.ReviewedAt) {

			event.AddOldValue("reviewedAt", item.ReviewedAt)
			event.AddNewValue("reviewedAt", changes.ReviewedAt)

			item.ReviewedAt = changes.ReviewedAt
			newItem.ReviewedAt = changes.ReviewedAt
			changedFields = append(changedFields, "reviewed_at")
			isChange = true
		}
	}

	if _, ok := input["rejectionReason"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.RejectionReason != changes.RejectionReason) && (item.RejectionReason == nil || changes.RejectionReason == nil || *item.RejectionReason != *changes.RejectionReason) {

			event.AddOldValue("rejectionReason", item.RejectionReason)
			event.AddNewValue("rejectionReason", changes.RejectionReason)

			item.RejectionReason = changes.RejectionReason
			newItem.RejectionReason = changes.RejectionReason
			changedFields = append(changedFields, "rejection_reason")
			isChange = true
		}
	}

	if _, ok := input["contactPhone"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ContactPhone != changes.ContactPhone) && (item.ContactPhone == nil || changes.ContactPhone == nil || *item.ContactPhone != *changes.ContactPhone) {

			event.AddOldValue("contactPhone", item.ContactPhone)
			event.AddNewValue("contactPhone", changes.ContactPhone)

			item.ContactPhone = changes.ContactPhone
			newItem.ContactPhone = changes.ContactPhone
			changedFields = append(changedFields, "contact_phone")
			isChange = true
		}
	}

	if _, ok := input["managerName"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ManagerName != changes.ManagerName) && (item.ManagerName == nil || changes.ManagerName == nil || *item.ManagerName != *changes.ManagerName) {

			event.AddOldValue("managerName", item.ManagerName)
			event.AddNewValue("managerName", changes.ManagerName)

			item.ManagerName = changes.ManagerName
			newItem.ManagerName = changes.ManagerName
			changedFields = append(changedFields, "manager_name")
			isChange = true
		}
	}

	if _, ok := input["managerPhone"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ManagerPhone != changes.ManagerPhone) && (item.ManagerPhone == nil || changes.ManagerPhone == nil || *item.ManagerPhone != *changes.ManagerPhone) {

			event.AddOldValue("managerPhone", item.ManagerPhone)
			event.AddNewValue("managerPhone", changes.ManagerPhone)

			item.ManagerPhone = changes.ManagerPhone
			newItem.ManagerPhone = changes.ManagerPhone
			changedFields = append(changedFields, "manager_phone")
			isChange = true
		}
	}

	if _, ok := input["province"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Province != changes.Province) && (item.Province == nil || changes.Province == nil || *item.Province != *changes.Province) {

			event.AddOldValue("province", item.Province)
			event.AddNewValue("province", changes.Province)

			item.Province = changes.Province
			newItem.Province = changes.Province
			changedFields = append(changedFields, "province")
			isChange = true
		}
	}

	if _, ok := input["city"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.City != changes.City) && (item.City == nil || changes.City == nil || *item.City != *changes.City) {

			event.AddOldValue("city", item.City)
			event.AddNewValue("city", changes.City)

			item.City = changes.City
			newItem.City = changes.City
			changedFields = append(changedFields, "city")
			isChange = true
		}
	}

	if _, ok := input["district"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.District != changes.District) && (item.District == nil || changes.District == nil || *item.District != *changes.District) {

			event.AddOldValue("district", item.District)
			event.AddNewValue("district", changes.District)

			item.District = changes.District
			newItem.District = changes.District
			changedFields = append(changedFields, "district")
			isChange = true
		}
	}

	if _, ok := input["address"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Address != changes.Address) && (item.Address == nil || changes.Address == nil || *item.Address != *changes.Address) {

			event.AddOldValue("address", item.Address)
			event.AddNewValue("address", changes.Address)

			item.Address = changes.Address
			newItem.Address = changes.Address
			changedFields = append(changedFields, "address")
			isChange = true
		}
	}

	if _, ok := input["businessHours"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.BusinessHours != changes.BusinessHours) && (item.BusinessHours == nil || changes.BusinessHours == nil || *item.BusinessHours != *changes.BusinessHours) {

			event.AddOldValue("businessHours", item.BusinessHours)
			event.AddNewValue("businessHours", changes.BusinessHours)

			item.BusinessHours = changes.BusinessHours
			newItem.BusinessHours = changes.BusinessHours
			changedFields = append(changedFields, "business_hours")
			isChange = true
		}
	}

	if _, ok := input["businessStatus"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.BusinessStatus != changes.BusinessStatus) && (item.BusinessStatus == nil || changes.BusinessStatus == nil || *item.BusinessStatus != *changes.BusinessStatus) {

			event.AddOldValue("businessStatus", item.BusinessStatus)
			event.AddNewValue("businessStatus", changes.BusinessStatus)

			item.BusinessStatus = changes.BusinessStatus
			newItem.BusinessStatus = changes.BusinessStatus
			changedFields = append(changedFields, "business_status")
			isChange = true
		}
	}

	if _, ok := input["supportDineIn"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.SupportDineIn != changes.SupportDineIn) && (item.SupportDineIn == nil || changes.SupportDineIn == nil || *item.SupportDineIn != *changes.SupportDineIn) {

			event.AddOldValue("supportDineIn", item.SupportDineIn)
			event.AddNewValue("supportDineIn", changes.SupportDineIn)

			item.SupportDineIn = changes.SupportDineIn
			newItem.SupportDineIn = changes.SupportDineIn
			changedFields = append(changedFields, "support_dine_in")
			isChange = true
		}
	}

	if _, ok := input["supportTakeout"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.SupportTakeout != changes.SupportTakeout) && (item.SupportTakeout == nil || changes.SupportTakeout == nil || *item.SupportTakeout != *changes.SupportTakeout) {

			event.AddOldValue("supportTakeout", item.SupportTakeout)
			event.AddNewValue("supportTakeout", changes.SupportTakeout)

			item.SupportTakeout = changes.SupportTakeout
			newItem.SupportTakeout = changes.SupportTakeout
			changedFields = append(changedFields, "support_takeout")
			isChange = true
		}
	}

	if _, ok := input["storeArea"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.StoreArea != changes.StoreArea) && (item.StoreArea == nil || changes.StoreArea == nil || *item.StoreArea != *changes.StoreArea) {

			event.AddOldValue("storeArea", item.StoreArea)
			event.AddNewValue("storeArea", changes.StoreArea)

			item.StoreArea = changes.StoreArea
			newItem.StoreArea = changes.StoreArea
			changedFields = append(changedFields, "store_area")
			isChange = true
		}
	}

	if _, ok := input["tableCount"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.TableCount != changes.TableCount) && (item.TableCount == nil || changes.TableCount == nil || *item.TableCount != *changes.TableCount) {

			event.AddOldValue("tableCount", item.TableCount)
			event.AddNewValue("tableCount", changes.TableCount)

			item.TableCount = changes.TableCount
			newItem.TableCount = changes.TableCount
			changedFields = append(changedFields, "table_count")
			isChange = true
		}
	}

	if _, ok := input["receiptFooter"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ReceiptFooter != changes.ReceiptFooter) && (item.ReceiptFooter == nil || changes.ReceiptFooter == nil || *item.ReceiptFooter != *changes.ReceiptFooter) {

			event.AddOldValue("receiptFooter", item.ReceiptFooter)
			event.AddNewValue("receiptFooter", changes.ReceiptFooter)

			item.ReceiptFooter = changes.ReceiptFooter
			newItem.ReceiptFooter = changes.ReceiptFooter
			changedFields = append(changedFields, "receipt_footer")
			isChange = true
		}
	}

	if _, ok := input["businessLicenseImageUrl"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.BusinessLicenseImageURL != changes.BusinessLicenseImageURL) && (item.BusinessLicenseImageURL == nil || changes.BusinessLicenseImageURL == nil || *item.BusinessLicenseImageURL != *changes.BusinessLicenseImageURL) {

			event.AddOldValue("businessLicenseImageUrl", item.BusinessLicenseImageURL)
			event.AddNewValue("businessLicenseImageUrl", changes.BusinessLicenseImageURL)

			item.BusinessLicenseImageURL = changes.BusinessLicenseImageURL
			newItem.BusinessLicenseImageURL = changes.BusinessLicenseImageURL
			changedFields = append(changedFields, "business_license_image_url")
			isChange = true
		}
	}

	if _, ok := input["otherDocumentImageUrl"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.OtherDocumentImageURL != changes.OtherDocumentImageURL) && (item.OtherDocumentImageURL == nil || changes.OtherDocumentImageURL == nil || *item.OtherDocumentImageURL != *changes.OtherDocumentImageURL) {

			event.AddOldValue("otherDocumentImageUrl", item.OtherDocumentImageURL)
			event.AddNewValue("otherDocumentImageUrl", changes.OtherDocumentImageURL)

			item.OtherDocumentImageURL = changes.OtherDocumentImageURL
			newItem.OtherDocumentImageURL = changes.OtherDocumentImageURL
			changedFields = append(changedFields, "other_document_image_url")
			isChange = true
		}
	}

	if _, ok := input["organizationId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.OrganizationID != changes.OrganizationID {

			if !utils.IsNil(input["organizationId"]) {
				if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
					return nil, fmt.Errorf("organizationId: %w", err)
				}
			}

			event.AddOldValue("organizationId", item.OrganizationID)
			event.AddNewValue("organizationId", changes.OrganizationID)

			item.OrganizationID = changes.OrganizationID
			newItem.OrganizationID = changes.OrganizationID
			changedFields = append(changedFields, "organization_id")
			isChange = true
		}
	}

	if _, ok := input["reviewedByAccountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ReviewedByAccountID != changes.ReviewedByAccountID) && (item.ReviewedByAccountID == nil || changes.ReviewedByAccountID == nil || *item.ReviewedByAccountID != *changes.ReviewedByAccountID) {

			if !utils.IsNil(input["reviewedByAccountId"]) {
				if err := tx.Select("id").Where("id = ?", input["reviewedByAccountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("reviewedByAccountId: %w", err)
				}
			}

			event.AddOldValue("reviewedByAccountId", item.ReviewedByAccountID)
			event.AddNewValue("reviewedByAccountId", changes.ReviewedByAccountID)

			item.ReviewedByAccountID = changes.ReviewedByAccountID
			newItem.ReviewedByAccountID = changes.ReviewedByAccountID
			changedFields = append(changedFields, "reviewed_by_account_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("stores", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// ---------- ToMany: members ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["membersIds"]; ok && !utils.IsNil(input["membersIds"]) {
		items := []*OperatorMembership{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
				return item, fmt.Errorf("OperatorMembership Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("membersIds %s not found", strings.Join(differenceIds, ","))
			}

			if err := tx.Model(item).Association("Members").Replace(items); err != nil {
				return item, err
			}

		} else {
			// 清空关联

			if err := tx.Model(item).Association("Members").Clear(); err != nil {
				return item, err
			}

		}
		event.AddNewValue("members", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["members"]; ok && !utils.IsNil(input["members"]) {
		newMembers := []*OperatorMembership{}
		updateMembers := []*OperatorMembership{}

		hasCreateMembers := false
		hasUpdateMembers := false

		for index, v := range changes.Members {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateMembers {
					if err := auth.CheckAuthorization(ctx, "UpdateOperatorMembership"); err != nil {
						return item, fmt.Errorf("UpdateOperatorMembership: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "OperatorMembership"); err != nil {
						return item, fmt.Errorf("OperatorMembership Detail: %w", err)
					}
					hasUpdateMembers = true
				}

				membersInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateOperatorMembership(ctx, r, membersInput["id"].(string), membersInput); err != nil {
					return item, fmt.Errorf("OperatorMembership ID %s: %w", v.ID, err)
				}

				updateMembers = append(updateMembers, v)
			} else {
				// 创建新记录
				if !hasCreateMembers {
					if err := auth.CheckAuthorization(ctx, "CreateOperatorMembership"); err != nil {
						return item, fmt.Errorf("CreateOperatorMembership: %w", err)
					}
					hasCreateMembers = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				if err := tx.Omit(clause.Associations).Table(TableName("operator_memberships", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newMembers = append(newMembers, v)
			}
		}

		allItems := append(updateMembers, newMembers...)

		if err := tx.Model(item).Association("Members").Replace(allItems); err != nil {
			return item, err
		}

		event.AddNewValue("members", allItems)
	}

	// ---------- ToMany: auditLogs ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&AuditLog{}).Where("store_id = ?", item.ID).Update("store_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&AuditLog{}).Where("store_id = ?", item.ID).Update("store_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&AuditLog{}).Where("store_id = ?", item.ID).Update("store_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.StoreID = &item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// ---------- ToMany: paymentConfigs ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["paymentConfigsIds"]; ok && !utils.IsNil(input["paymentConfigsIds"]) {
		items := []*StorePaymentConfig{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "StorePaymentConfig"); err != nil {
				return item, fmt.Errorf("StorePaymentConfig Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("paymentConfigsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&StorePaymentConfig{}).Where("store_id = ?", item.ID).Update("store_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&StorePaymentConfig{}).Where("store_id = ?", item.ID).Update("store_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("paymentConfigs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["paymentConfigs"]; ok && !utils.IsNil(input["paymentConfigs"]) {
		newPaymentConfigs := []*StorePaymentConfig{}
		updatePaymentConfigs := []*StorePaymentConfig{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&StorePaymentConfig{}).Where("store_id = ?", item.ID).Update("store_id", nil).Error; err != nil {
			return item, err
		}

		hasCreatePaymentConfigs := false
		hasUpdatePaymentConfigs := false

		for index, v := range changes.PaymentConfigs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "UpdateStorePaymentConfig"); err != nil {
						return item, fmt.Errorf("UpdateStorePaymentConfig: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "StorePaymentConfig"); err != nil {
						return item, fmt.Errorf("StorePaymentConfig Detail: %w", err)
					}
					hasUpdatePaymentConfigs = true
				}

				paymentConfigsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateStorePaymentConfig(ctx, r, paymentConfigsInput["id"].(string), paymentConfigsInput); err != nil {
					return item, fmt.Errorf("StorePaymentConfig ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("store_id", item.ID).Error; err != nil {
					return item, err
				}

				updatePaymentConfigs = append(updatePaymentConfigs, v)
			} else {
				// 创建新记录
				if !hasCreatePaymentConfigs {
					if err := auth.CheckAuthorization(ctx, "CreateStorePaymentConfig"); err != nil {
						return item, fmt.Errorf("CreateStorePaymentConfig: %w", err)
					}
					hasCreatePaymentConfigs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.StoreID = item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("store_payment_configs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newPaymentConfigs = append(newPaymentConfigs, v)
			}
		}

		allItems := append(updatePaymentConfigs, newPaymentConfigs...)

		event.AddNewValue("paymentConfigs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Store - Delete
// ============================================================

// DeleteStoreFunc 执行删除或恢复操作
func DeleteStoreFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &Store{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("stores", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "Store",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("stores", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("stores", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("stores", ctx)).Updates(Store{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteStores 批量删除 Store 实体
func (r *GeneratedMutationResolver) DeleteStores(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteStores(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteStoresHandler 处理批量删除逻辑
func DeleteStoresHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteStoreFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Store - Recovery
// ============================================================

// RecoveryStores 批量恢复 Store 实体
func (r *GeneratedMutationResolver) RecoveryStores(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryStores(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryStoresHandler 处理批量恢复逻辑
func RecoveryStoresHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteStoreFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Session - Create
// ============================================================

// CreateSession 创建 Session 实体的解析器入口
func (r *GeneratedMutationResolver) CreateSession(ctx context.Context, input map[string]interface{}) (item *Session, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateSession(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateSessionHandler 处理 Session 创建逻辑
func CreateSessionHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *Session, err error) {
	item = &Session{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "Session",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes SessionChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["account"]) && !utils.IsNil(input["accountId"]) {
		return nil, fmt.Errorf("accountId and account cannot coexist")
	}

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	// ToMany: auditLogs - 不能同时传入 IDs 和嵌套对象
	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["workspaceType"]; ok {

		item.WorkspaceType = changes.WorkspaceType

		event.AddNewValue("workspaceType", changes.WorkspaceType)
	}

	if _, ok := input["credentialVersion"]; ok {

		item.CredentialVersion = changes.CredentialVersion

		event.AddNewValue("credentialVersion", changes.CredentialVersion)
	}

	if _, ok := input["expiresAt"]; ok {

		item.ExpiresAt = changes.ExpiresAt

		event.AddNewValue("expiresAt", changes.ExpiresAt)
	}

	if _, ok := input["revokedAt"]; ok && changes.RevokedAt != nil {

		item.RevokedAt = changes.RevokedAt

		event.AddNewValue("revokedAt", changes.RevokedAt)
	}

	if _, ok := input["revocationCode"]; ok && changes.RevocationCode != nil {

		item.RevocationCode = changes.RevocationCode

		event.AddNewValue("revocationCode", changes.RevocationCode)
	}

	if _, ok := input["lastSeenAt"]; ok {

		item.LastSeenAt = changes.LastSeenAt

		event.AddNewValue("lastSeenAt", changes.LastSeenAt)
	}

	if _, ok := input["accountId"]; ok {

		if !utils.IsNil(input["accountId"]) {
			if err := tx.Select("id").Where("id = ?", input["accountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("accountId: %w", err)
			}
		}

		item.AccountID = changes.AccountID

		event.AddNewValue("accountId", changes.AccountID)
	}

	if _, ok := input["organizationId"]; ok && changes.OrganizationID != nil {

		if !utils.IsNil(input["organizationId"]) {
			if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
				return nil, fmt.Errorf("organizationId: %w", err)
			}
		}

		item.OrganizationID = changes.OrganizationID

		event.AddNewValue("organizationId", changes.OrganizationID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("sessions", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// ---------- ToMany: auditLogs (OneToMany 外键在 AuditLog.session_id) ----------

	// 方式1：通过 IDs 关联现有记录
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			// 权限检查
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}

			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}

			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			// 验证所有 ID 都存在
			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 更新关联记录的外键
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("session_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				// OneToMany: 设置外键指向当前实体
				if err := tx.Model(v).Update("session_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				// OneToMany: 设置外键指向当前实体

				v.SessionID = &item.ID

				// 保存新记录
				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Session - Update
// ============================================================

// UpdateSession 更新 Session 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateSession(ctx context.Context, id string, input map[string]interface{}) (item *Session, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateSession(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateSessionHandler 处理 Session 更新逻辑
func UpdateSessionHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *Session, err error) {
	item = &Session{}
	newItem := &Session{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "Session",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes SessionChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["account"]) && !utils.IsNil(input["accountId"]) {
		return nil, fmt.Errorf("accountId and account cannot coexist")
	}

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["auditLogs"]) && !utils.IsNil(input["auditLogsIds"]) {
		return nil, fmt.Errorf("auditLogsIds and auditLogs cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("sessions", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["workspaceType"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.WorkspaceType != changes.WorkspaceType {

			event.AddOldValue("workspaceType", item.WorkspaceType)
			event.AddNewValue("workspaceType", changes.WorkspaceType)

			item.WorkspaceType = changes.WorkspaceType
			newItem.WorkspaceType = changes.WorkspaceType
			changedFields = append(changedFields, "workspace_type")
			isChange = true
		}
	}

	if _, ok := input["credentialVersion"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.CredentialVersion != changes.CredentialVersion {

			event.AddOldValue("credentialVersion", item.CredentialVersion)
			event.AddNewValue("credentialVersion", changes.CredentialVersion)

			item.CredentialVersion = changes.CredentialVersion
			newItem.CredentialVersion = changes.CredentialVersion
			changedFields = append(changedFields, "credential_version")
			isChange = true
		}
	}

	if _, ok := input["expiresAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ExpiresAt != changes.ExpiresAt {

			event.AddOldValue("expiresAt", item.ExpiresAt)
			event.AddNewValue("expiresAt", changes.ExpiresAt)

			item.ExpiresAt = changes.ExpiresAt
			newItem.ExpiresAt = changes.ExpiresAt
			changedFields = append(changedFields, "expires_at")
			isChange = true
		}
	}

	if _, ok := input["revokedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.RevokedAt != changes.RevokedAt) && (item.RevokedAt == nil || changes.RevokedAt == nil || *item.RevokedAt != *changes.RevokedAt) {

			event.AddOldValue("revokedAt", item.RevokedAt)
			event.AddNewValue("revokedAt", changes.RevokedAt)

			item.RevokedAt = changes.RevokedAt
			newItem.RevokedAt = changes.RevokedAt
			changedFields = append(changedFields, "revoked_at")
			isChange = true
		}
	}

	if _, ok := input["revocationCode"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.RevocationCode != changes.RevocationCode) && (item.RevocationCode == nil || changes.RevocationCode == nil || *item.RevocationCode != *changes.RevocationCode) {

			event.AddOldValue("revocationCode", item.RevocationCode)
			event.AddNewValue("revocationCode", changes.RevocationCode)

			item.RevocationCode = changes.RevocationCode
			newItem.RevocationCode = changes.RevocationCode
			changedFields = append(changedFields, "revocation_code")
			isChange = true
		}
	}

	if _, ok := input["lastSeenAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.LastSeenAt != changes.LastSeenAt {

			event.AddOldValue("lastSeenAt", item.LastSeenAt)
			event.AddNewValue("lastSeenAt", changes.LastSeenAt)

			item.LastSeenAt = changes.LastSeenAt
			newItem.LastSeenAt = changes.LastSeenAt
			changedFields = append(changedFields, "last_seen_at")
			isChange = true
		}
	}

	if _, ok := input["accountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.AccountID != changes.AccountID {

			if !utils.IsNil(input["accountId"]) {
				if err := tx.Select("id").Where("id = ?", input["accountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("accountId: %w", err)
				}
			}

			event.AddOldValue("accountId", item.AccountID)
			event.AddNewValue("accountId", changes.AccountID)

			item.AccountID = changes.AccountID
			newItem.AccountID = changes.AccountID
			changedFields = append(changedFields, "account_id")
			isChange = true
		}
	}

	if _, ok := input["organizationId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.OrganizationID != changes.OrganizationID) && (item.OrganizationID == nil || changes.OrganizationID == nil || *item.OrganizationID != *changes.OrganizationID) {

			if !utils.IsNil(input["organizationId"]) {
				if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
					return nil, fmt.Errorf("organizationId: %w", err)
				}
			}

			event.AddOldValue("organizationId", item.OrganizationID)
			event.AddNewValue("organizationId", changes.OrganizationID)

			item.OrganizationID = changes.OrganizationID
			newItem.OrganizationID = changes.OrganizationID
			changedFields = append(changedFields, "organization_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("sessions", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// ---------- ToMany: auditLogs ----------

	// 方式1：通过 IDs 关联
	if ids, ok := input["auditLogsIds"]; ok && !utils.IsNil(input["auditLogsIds"]) {
		items := []*AuditLog{}
		itemIds := []string{}
		findIds := []string{}

		for _, v := range ids.([]string) {
			itemIds = append(itemIds, v)
		}

		if len(itemIds) > 0 {
			if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
				return item, fmt.Errorf("AuditLog Detail: %w", err)
			}
			if err := tx.Find(&items, "id IN (?)", itemIds).Error; err != nil {
				return item, err
			}
			for _, v := range items {
				findIds = append(findIds, v.ID)
			}

			differenceIds := utils.Difference(itemIds, findIds)
			if len(differenceIds) > 0 {
				return item, fmt.Errorf("auditLogsIds %s not found", strings.Join(differenceIds, ","))
			}

			// OneToMany: 先清除旧关联，再设置新关联
			if err := tx.Model(&AuditLog{}).Where("session_id = ?", item.ID).Update("session_id", nil).Error; err != nil {
				return item, err
			}
			for _, relItem := range items {
				if err := tx.Model(relItem).Update("session_id", item.ID).Error; err != nil {
					return item, err
				}
			}

		} else {
			// 清空关联

			if err := tx.Model(&AuditLog{}).Where("session_id = ?", item.ID).Update("session_id", nil).Error; err != nil {
				return item, err
			}

		}
		event.AddNewValue("auditLogs", items)
	}

	// 方式2：通过嵌套对象创建/更新
	if _, ok := input["auditLogs"]; ok && !utils.IsNil(input["auditLogs"]) {
		newAuditLogs := []*AuditLog{}
		updateAuditLogs := []*AuditLog{}

		// OneToMany: 先清除旧关联（与 IDs 方式行为一致）
		if err := tx.Model(&AuditLog{}).Where("session_id = ?", item.ID).Update("session_id", nil).Error; err != nil {
			return item, err
		}

		hasCreateAuditLogs := false
		hasUpdateAuditLogs := false

		for index, v := range changes.AuditLogs {
			weight := int64(index + 1)
			v.Weight = &weight

			if !utils.IsEmpty(v.ID) {
				// 更新现有记录
				v.UpdatedAt = &timestampMillis
				v.UpdatedBy = principalID

				if !hasUpdateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "UpdateAuditLog"); err != nil {
						return item, fmt.Errorf("UpdateAuditLog: %w", err)
					}
					if err := auth.CheckAuthorization(ctx, "AuditLog"); err != nil {
						return item, fmt.Errorf("AuditLog Detail: %w", err)
					}
					hasUpdateAuditLogs = true
				}

				auditLogsInput := utils.StructToMap(*v)
				if _, err := r.Handlers.UpdateAuditLog(ctx, r, auditLogsInput["id"].(string), auditLogsInput); err != nil {
					return item, fmt.Errorf("AuditLog ID %s: %w", v.ID, err)
				}

				if err := tx.Model(v).Update("session_id", item.ID).Error; err != nil {
					return item, err
				}

				updateAuditLogs = append(updateAuditLogs, v)
			} else {
				// 创建新记录
				if !hasCreateAuditLogs {
					if err := auth.CheckAuthorization(ctx, "CreateAuditLog"); err != nil {
						return item, fmt.Errorf("CreateAuditLog: %w", err)
					}
					hasCreateAuditLogs = true
				}

				v.ID = uuid.Must(uuid.NewV4()).String()
				v.CreatedAt = timestampMillis
				v.CreatedBy = principalID

				v.SessionID = &item.ID

				if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(v).Error; err != nil {
					return item, err
				}

				newAuditLogs = append(newAuditLogs, v)
			}
		}

		allItems := append(updateAuditLogs, newAuditLogs...)

		event.AddNewValue("auditLogs", allItems)
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// Session - Delete
// ============================================================

// DeleteSessionFunc 执行删除或恢复操作
func DeleteSessionFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &Session{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("sessions", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "Session",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("sessions", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("sessions", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("sessions", ctx)).Updates(Session{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteSessions 批量删除 Session 实体
func (r *GeneratedMutationResolver) DeleteSessions(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteSessions(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteSessionsHandler 处理批量删除逻辑
func DeleteSessionsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteSessionFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// Session - Recovery
// ============================================================

// RecoverySessions 批量恢复 Session 实体
func (r *GeneratedMutationResolver) RecoverySessions(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoverySessions(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoverySessionsHandler 处理批量恢复逻辑
func RecoverySessionsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteSessionFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// MembershipInvitation - Create
// ============================================================

// CreateMembershipInvitation 创建 MembershipInvitation 实体的解析器入口
func (r *GeneratedMutationResolver) CreateMembershipInvitation(ctx context.Context, input map[string]interface{}) (item *MembershipInvitation, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateMembershipInvitation(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateMembershipInvitationHandler 处理 MembershipInvitation 创建逻辑
func CreateMembershipInvitationHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *MembershipInvitation, err error) {
	item = &MembershipInvitation{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "MembershipInvitation",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes MembershipInvitationChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["membership"]) && !utils.IsNil(input["membershipId"]) {
		return nil, fmt.Errorf("membershipId and membership cannot coexist")
	}

	if !utils.IsNil(input["invitedByAccount"]) && !utils.IsNil(input["invitedByAccountId"]) {
		return nil, fmt.Errorf("invitedByAccountId and invitedByAccount cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["expiresAt"]; ok {

		item.ExpiresAt = changes.ExpiresAt

		event.AddNewValue("expiresAt", changes.ExpiresAt)
	}

	if _, ok := input["acceptedAt"]; ok && changes.AcceptedAt != nil {

		item.AcceptedAt = changes.AcceptedAt

		event.AddNewValue("acceptedAt", changes.AcceptedAt)
	}

	if _, ok := input["revokedAt"]; ok && changes.RevokedAt != nil {

		item.RevokedAt = changes.RevokedAt

		event.AddNewValue("revokedAt", changes.RevokedAt)
	}

	if _, ok := input["membershipId"]; ok {

		if !utils.IsNil(input["membershipId"]) {
			if err := tx.Select("id").Where("id = ?", input["membershipId"]).First(&OperatorMembership{}).Error; err != nil {
				return nil, fmt.Errorf("membershipId: %w", err)
			}
		}

		item.MembershipID = changes.MembershipID

		event.AddNewValue("membershipId", changes.MembershipID)
	}

	if _, ok := input["invitedByAccountId"]; ok {

		if !utils.IsNil(input["invitedByAccountId"]) {
			if err := tx.Select("id").Where("id = ?", input["invitedByAccountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("invitedByAccountId: %w", err)
			}
		}

		item.InvitedByAccountID = changes.InvitedByAccountID

		event.AddNewValue("invitedByAccountId", changes.InvitedByAccountID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("membership_invitations", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// MembershipInvitation - Update
// ============================================================

// UpdateMembershipInvitation 更新 MembershipInvitation 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateMembershipInvitation(ctx context.Context, id string, input map[string]interface{}) (item *MembershipInvitation, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateMembershipInvitation(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateMembershipInvitationHandler 处理 MembershipInvitation 更新逻辑
func UpdateMembershipInvitationHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *MembershipInvitation, err error) {
	item = &MembershipInvitation{}
	newItem := &MembershipInvitation{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "MembershipInvitation",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes MembershipInvitationChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["membership"]) && !utils.IsNil(input["membershipId"]) {
		return nil, fmt.Errorf("membershipId and membership cannot coexist")
	}

	if !utils.IsNil(input["invitedByAccount"]) && !utils.IsNil(input["invitedByAccountId"]) {
		return nil, fmt.Errorf("invitedByAccountId and invitedByAccount cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("membership_invitations", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["expiresAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ExpiresAt != changes.ExpiresAt {

			event.AddOldValue("expiresAt", item.ExpiresAt)
			event.AddNewValue("expiresAt", changes.ExpiresAt)

			item.ExpiresAt = changes.ExpiresAt
			newItem.ExpiresAt = changes.ExpiresAt
			changedFields = append(changedFields, "expires_at")
			isChange = true
		}
	}

	if _, ok := input["acceptedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.AcceptedAt != changes.AcceptedAt) && (item.AcceptedAt == nil || changes.AcceptedAt == nil || *item.AcceptedAt != *changes.AcceptedAt) {

			event.AddOldValue("acceptedAt", item.AcceptedAt)
			event.AddNewValue("acceptedAt", changes.AcceptedAt)

			item.AcceptedAt = changes.AcceptedAt
			newItem.AcceptedAt = changes.AcceptedAt
			changedFields = append(changedFields, "accepted_at")
			isChange = true
		}
	}

	if _, ok := input["revokedAt"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.RevokedAt != changes.RevokedAt) && (item.RevokedAt == nil || changes.RevokedAt == nil || *item.RevokedAt != *changes.RevokedAt) {

			event.AddOldValue("revokedAt", item.RevokedAt)
			event.AddNewValue("revokedAt", changes.RevokedAt)

			item.RevokedAt = changes.RevokedAt
			newItem.RevokedAt = changes.RevokedAt
			changedFields = append(changedFields, "revoked_at")
			isChange = true
		}
	}

	if _, ok := input["membershipId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.MembershipID != changes.MembershipID {

			if !utils.IsNil(input["membershipId"]) {
				if err := tx.Select("id").Where("id = ?", input["membershipId"]).First(&OperatorMembership{}).Error; err != nil {
					return nil, fmt.Errorf("membershipId: %w", err)
				}
			}

			event.AddOldValue("membershipId", item.MembershipID)
			event.AddNewValue("membershipId", changes.MembershipID)

			item.MembershipID = changes.MembershipID
			newItem.MembershipID = changes.MembershipID
			changedFields = append(changedFields, "membership_id")
			isChange = true
		}
	}

	if _, ok := input["invitedByAccountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.InvitedByAccountID != changes.InvitedByAccountID {

			if !utils.IsNil(input["invitedByAccountId"]) {
				if err := tx.Select("id").Where("id = ?", input["invitedByAccountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("invitedByAccountId: %w", err)
				}
			}

			event.AddOldValue("invitedByAccountId", item.InvitedByAccountID)
			event.AddNewValue("invitedByAccountId", changes.InvitedByAccountID)

			item.InvitedByAccountID = changes.InvitedByAccountID
			newItem.InvitedByAccountID = changes.InvitedByAccountID
			changedFields = append(changedFields, "invited_by_account_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("membership_invitations", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// MembershipInvitation - Delete
// ============================================================

// DeleteMembershipInvitationFunc 执行删除或恢复操作
func DeleteMembershipInvitationFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &MembershipInvitation{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("membership_invitations", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "MembershipInvitation",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("membership_invitations", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("membership_invitations", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("membership_invitations", ctx)).Updates(MembershipInvitation{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteMembershipInvitations 批量删除 MembershipInvitation 实体
func (r *GeneratedMutationResolver) DeleteMembershipInvitations(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteMembershipInvitations(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteMembershipInvitationsHandler 处理批量删除逻辑
func DeleteMembershipInvitationsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteMembershipInvitationFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// MembershipInvitation - Recovery
// ============================================================

// RecoveryMembershipInvitations 批量恢复 MembershipInvitation 实体
func (r *GeneratedMutationResolver) RecoveryMembershipInvitations(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryMembershipInvitations(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryMembershipInvitationsHandler 处理批量恢复逻辑
func RecoveryMembershipInvitationsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteMembershipInvitationFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// AuditLog - Create
// ============================================================

// CreateAuditLog 创建 AuditLog 实体的解析器入口
func (r *GeneratedMutationResolver) CreateAuditLog(ctx context.Context, input map[string]interface{}) (item *AuditLog, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateAuditLog(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateAuditLogHandler 处理 AuditLog 创建逻辑
func CreateAuditLogHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *AuditLog, err error) {
	item = &AuditLog{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "AuditLog",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes AuditLogChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["actorAccount"]) && !utils.IsNil(input["actorAccountId"]) {
		return nil, fmt.Errorf("actorAccountId and actorAccount cannot coexist")
	}

	if !utils.IsNil(input["session"]) && !utils.IsNil(input["sessionId"]) {
		return nil, fmt.Errorf("sessionId and session cannot coexist")
	}

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["store"]) && !utils.IsNil(input["storeId"]) {
		return nil, fmt.Errorf("storeId and store cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["action"]; ok {

		item.Action = changes.Action

		event.AddNewValue("action", changes.Action)
	}

	if _, ok := input["resourceType"]; ok {

		item.ResourceType = changes.ResourceType

		event.AddNewValue("resourceType", changes.ResourceType)
	}

	if _, ok := input["resourceId"]; ok && changes.ResourceID != nil {

		item.ResourceID = changes.ResourceID

		event.AddNewValue("resourceId", changes.ResourceID)
	}

	if _, ok := input["resultCode"]; ok {

		item.ResultCode = changes.ResultCode

		event.AddNewValue("resultCode", changes.ResultCode)
	}

	if _, ok := input["metadataJson"]; ok && changes.MetadataJSON != nil {

		item.MetadataJSON = changes.MetadataJSON

		event.AddNewValue("metadataJson", changes.MetadataJSON)
	}

	if _, ok := input["actorAccountId"]; ok && changes.ActorAccountID != nil {

		if !utils.IsNil(input["actorAccountId"]) {
			if err := tx.Select("id").Where("id = ?", input["actorAccountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("actorAccountId: %w", err)
			}
		}

		item.ActorAccountID = changes.ActorAccountID

		event.AddNewValue("actorAccountId", changes.ActorAccountID)
	}

	if _, ok := input["sessionId"]; ok && changes.SessionID != nil {

		if !utils.IsNil(input["sessionId"]) {
			if err := tx.Select("id").Where("id = ?", input["sessionId"]).First(&Session{}).Error; err != nil {
				return nil, fmt.Errorf("sessionId: %w", err)
			}
		}

		item.SessionID = changes.SessionID

		event.AddNewValue("sessionId", changes.SessionID)
	}

	if _, ok := input["organizationId"]; ok && changes.OrganizationID != nil {

		if !utils.IsNil(input["organizationId"]) {
			if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
				return nil, fmt.Errorf("organizationId: %w", err)
			}
		}

		item.OrganizationID = changes.OrganizationID

		event.AddNewValue("organizationId", changes.OrganizationID)
	}

	if _, ok := input["storeId"]; ok && changes.StoreID != nil {

		if !utils.IsNil(input["storeId"]) {
			if err := tx.Select("id").Where("id = ?", input["storeId"]).First(&Store{}).Error; err != nil {
				return nil, fmt.Errorf("storeId: %w", err)
			}
		}

		item.StoreID = changes.StoreID

		event.AddNewValue("storeId", changes.StoreID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("audit_logs", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// AuditLog - Update
// ============================================================

// UpdateAuditLog 更新 AuditLog 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateAuditLog(ctx context.Context, id string, input map[string]interface{}) (item *AuditLog, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateAuditLog(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateAuditLogHandler 处理 AuditLog 更新逻辑
func UpdateAuditLogHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *AuditLog, err error) {
	item = &AuditLog{}
	newItem := &AuditLog{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "AuditLog",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes AuditLogChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["actorAccount"]) && !utils.IsNil(input["actorAccountId"]) {
		return nil, fmt.Errorf("actorAccountId and actorAccount cannot coexist")
	}

	if !utils.IsNil(input["session"]) && !utils.IsNil(input["sessionId"]) {
		return nil, fmt.Errorf("sessionId and session cannot coexist")
	}

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["store"]) && !utils.IsNil(input["storeId"]) {
		return nil, fmt.Errorf("storeId and store cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("audit_logs", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["action"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Action != changes.Action {

			event.AddOldValue("action", item.Action)
			event.AddNewValue("action", changes.Action)

			item.Action = changes.Action
			newItem.Action = changes.Action
			changedFields = append(changedFields, "action")
			isChange = true
		}
	}

	if _, ok := input["resourceType"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ResourceType != changes.ResourceType {

			event.AddOldValue("resourceType", item.ResourceType)
			event.AddNewValue("resourceType", changes.ResourceType)

			item.ResourceType = changes.ResourceType
			newItem.ResourceType = changes.ResourceType
			changedFields = append(changedFields, "resource_type")
			isChange = true
		}
	}

	if _, ok := input["resourceId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ResourceID != changes.ResourceID) && (item.ResourceID == nil || changes.ResourceID == nil || *item.ResourceID != *changes.ResourceID) {

			event.AddOldValue("resourceId", item.ResourceID)
			event.AddNewValue("resourceId", changes.ResourceID)

			item.ResourceID = changes.ResourceID
			newItem.ResourceID = changes.ResourceID
			changedFields = append(changedFields, "resource_id")
			isChange = true
		}
	}

	if _, ok := input["resultCode"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ResultCode != changes.ResultCode {

			event.AddOldValue("resultCode", item.ResultCode)
			event.AddNewValue("resultCode", changes.ResultCode)

			item.ResultCode = changes.ResultCode
			newItem.ResultCode = changes.ResultCode
			changedFields = append(changedFields, "result_code")
			isChange = true
		}
	}

	if _, ok := input["metadataJson"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.MetadataJSON != changes.MetadataJSON) && (item.MetadataJSON == nil || changes.MetadataJSON == nil || *item.MetadataJSON != *changes.MetadataJSON) {

			event.AddOldValue("metadataJson", item.MetadataJSON)
			event.AddNewValue("metadataJson", changes.MetadataJSON)

			item.MetadataJSON = changes.MetadataJSON
			newItem.MetadataJSON = changes.MetadataJSON
			changedFields = append(changedFields, "metadata_json")
			isChange = true
		}
	}

	if _, ok := input["actorAccountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.ActorAccountID != changes.ActorAccountID) && (item.ActorAccountID == nil || changes.ActorAccountID == nil || *item.ActorAccountID != *changes.ActorAccountID) {

			if !utils.IsNil(input["actorAccountId"]) {
				if err := tx.Select("id").Where("id = ?", input["actorAccountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("actorAccountId: %w", err)
				}
			}

			event.AddOldValue("actorAccountId", item.ActorAccountID)
			event.AddNewValue("actorAccountId", changes.ActorAccountID)

			item.ActorAccountID = changes.ActorAccountID
			newItem.ActorAccountID = changes.ActorAccountID
			changedFields = append(changedFields, "actor_account_id")
			isChange = true
		}
	}

	if _, ok := input["sessionId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.SessionID != changes.SessionID) && (item.SessionID == nil || changes.SessionID == nil || *item.SessionID != *changes.SessionID) {

			if !utils.IsNil(input["sessionId"]) {
				if err := tx.Select("id").Where("id = ?", input["sessionId"]).First(&Session{}).Error; err != nil {
					return nil, fmt.Errorf("sessionId: %w", err)
				}
			}

			event.AddOldValue("sessionId", item.SessionID)
			event.AddNewValue("sessionId", changes.SessionID)

			item.SessionID = changes.SessionID
			newItem.SessionID = changes.SessionID
			changedFields = append(changedFields, "session_id")
			isChange = true
		}
	}

	if _, ok := input["organizationId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.OrganizationID != changes.OrganizationID) && (item.OrganizationID == nil || changes.OrganizationID == nil || *item.OrganizationID != *changes.OrganizationID) {

			if !utils.IsNil(input["organizationId"]) {
				if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
					return nil, fmt.Errorf("organizationId: %w", err)
				}
			}

			event.AddOldValue("organizationId", item.OrganizationID)
			event.AddNewValue("organizationId", changes.OrganizationID)

			item.OrganizationID = changes.OrganizationID
			newItem.OrganizationID = changes.OrganizationID
			changedFields = append(changedFields, "organization_id")
			isChange = true
		}
	}

	if _, ok := input["storeId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.StoreID != changes.StoreID) && (item.StoreID == nil || changes.StoreID == nil || *item.StoreID != *changes.StoreID) {

			if !utils.IsNil(input["storeId"]) {
				if err := tx.Select("id").Where("id = ?", input["storeId"]).First(&Store{}).Error; err != nil {
					return nil, fmt.Errorf("storeId: %w", err)
				}
			}

			event.AddOldValue("storeId", item.StoreID)
			event.AddNewValue("storeId", changes.StoreID)

			item.StoreID = changes.StoreID
			newItem.StoreID = changes.StoreID
			changedFields = append(changedFields, "store_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("audit_logs", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// AuditLog - Delete
// ============================================================

// DeleteAuditLogFunc 执行删除或恢复操作
func DeleteAuditLogFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &AuditLog{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("audit_logs", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "AuditLog",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("audit_logs", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("audit_logs", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("audit_logs", ctx)).Updates(AuditLog{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteAuditLogs 批量删除 AuditLog 实体
func (r *GeneratedMutationResolver) DeleteAuditLogs(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteAuditLogs(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteAuditLogsHandler 处理批量删除逻辑
func DeleteAuditLogsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteAuditLogFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// AuditLog - Recovery
// ============================================================

// RecoveryAuditLogs 批量恢复 AuditLog 实体
func (r *GeneratedMutationResolver) RecoveryAuditLogs(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryAuditLogs(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryAuditLogsHandler 处理批量恢复逻辑
func RecoveryAuditLogsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteAuditLogFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// FranchiseOpeningRecord - Create
// ============================================================

// CreateFranchiseOpeningRecord 创建 FranchiseOpeningRecord 实体的解析器入口
func (r *GeneratedMutationResolver) CreateFranchiseOpeningRecord(ctx context.Context, input map[string]interface{}) (item *FranchiseOpeningRecord, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateFranchiseOpeningRecord(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateFranchiseOpeningRecordHandler 处理 FranchiseOpeningRecord 创建逻辑
func CreateFranchiseOpeningRecordHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *FranchiseOpeningRecord, err error) {
	item = &FranchiseOpeningRecord{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "FranchiseOpeningRecord",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes FranchiseOpeningRecordChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["initialAccount"]) && !utils.IsNil(input["initialAccountId"]) {
		return nil, fmt.Errorf("initialAccountId and initialAccount cannot coexist")
	}

	if !utils.IsNil(input["recordedByAccount"]) && !utils.IsNil(input["recordedByAccountId"]) {
		return nil, fmt.Errorf("recordedByAccountId and recordedByAccount cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["recordNumber"]; ok {

		item.RecordNumber = changes.RecordNumber

		event.AddNewValue("recordNumber", changes.RecordNumber)
	}

	if _, ok := input["source"]; ok {

		item.Source = changes.Source

		event.AddNewValue("source", changes.Source)
	}

	if _, ok := input["organizationId"]; ok {

		if !utils.IsNil(input["organizationId"]) {
			if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
				return nil, fmt.Errorf("organizationId: %w", err)
			}
		}

		item.OrganizationID = changes.OrganizationID

		event.AddNewValue("organizationId", changes.OrganizationID)
	}

	if _, ok := input["initialAccountId"]; ok {

		if !utils.IsNil(input["initialAccountId"]) {
			if err := tx.Select("id").Where("id = ?", input["initialAccountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("initialAccountId: %w", err)
			}
		}

		item.InitialAccountID = changes.InitialAccountID

		event.AddNewValue("initialAccountId", changes.InitialAccountID)
	}

	if _, ok := input["recordedByAccountId"]; ok {

		if !utils.IsNil(input["recordedByAccountId"]) {
			if err := tx.Select("id").Where("id = ?", input["recordedByAccountId"]).First(&Account{}).Error; err != nil {
				return nil, fmt.Errorf("recordedByAccountId: %w", err)
			}
		}

		item.RecordedByAccountID = changes.RecordedByAccountID

		event.AddNewValue("recordedByAccountId", changes.RecordedByAccountID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("franchise_opening_records", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// FranchiseOpeningRecord - Update
// ============================================================

// UpdateFranchiseOpeningRecord 更新 FranchiseOpeningRecord 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateFranchiseOpeningRecord(ctx context.Context, id string, input map[string]interface{}) (item *FranchiseOpeningRecord, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateFranchiseOpeningRecord(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateFranchiseOpeningRecordHandler 处理 FranchiseOpeningRecord 更新逻辑
func UpdateFranchiseOpeningRecordHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *FranchiseOpeningRecord, err error) {
	item = &FranchiseOpeningRecord{}
	newItem := &FranchiseOpeningRecord{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "FranchiseOpeningRecord",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes FranchiseOpeningRecordChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	if !utils.IsNil(input["initialAccount"]) && !utils.IsNil(input["initialAccountId"]) {
		return nil, fmt.Errorf("initialAccountId and initialAccount cannot coexist")
	}

	if !utils.IsNil(input["recordedByAccount"]) && !utils.IsNil(input["recordedByAccountId"]) {
		return nil, fmt.Errorf("recordedByAccountId and recordedByAccount cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("franchise_opening_records", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["recordNumber"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.RecordNumber != changes.RecordNumber {

			event.AddOldValue("recordNumber", item.RecordNumber)
			event.AddNewValue("recordNumber", changes.RecordNumber)

			item.RecordNumber = changes.RecordNumber
			newItem.RecordNumber = changes.RecordNumber
			changedFields = append(changedFields, "record_number")
			isChange = true
		}
	}

	if _, ok := input["source"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Source != changes.Source {

			event.AddOldValue("source", item.Source)
			event.AddNewValue("source", changes.Source)

			item.Source = changes.Source
			newItem.Source = changes.Source
			changedFields = append(changedFields, "source")
			isChange = true
		}
	}

	if _, ok := input["organizationId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.OrganizationID != changes.OrganizationID {

			if !utils.IsNil(input["organizationId"]) {
				if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
					return nil, fmt.Errorf("organizationId: %w", err)
				}
			}

			event.AddOldValue("organizationId", item.OrganizationID)
			event.AddNewValue("organizationId", changes.OrganizationID)

			item.OrganizationID = changes.OrganizationID
			newItem.OrganizationID = changes.OrganizationID
			changedFields = append(changedFields, "organization_id")
			isChange = true
		}
	}

	if _, ok := input["initialAccountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.InitialAccountID != changes.InitialAccountID {

			if !utils.IsNil(input["initialAccountId"]) {
				if err := tx.Select("id").Where("id = ?", input["initialAccountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("initialAccountId: %w", err)
				}
			}

			event.AddOldValue("initialAccountId", item.InitialAccountID)
			event.AddNewValue("initialAccountId", changes.InitialAccountID)

			item.InitialAccountID = changes.InitialAccountID
			newItem.InitialAccountID = changes.InitialAccountID
			changedFields = append(changedFields, "initial_account_id")
			isChange = true
		}
	}

	if _, ok := input["recordedByAccountId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.RecordedByAccountID != changes.RecordedByAccountID {

			if !utils.IsNil(input["recordedByAccountId"]) {
				if err := tx.Select("id").Where("id = ?", input["recordedByAccountId"]).First(&Account{}).Error; err != nil {
					return nil, fmt.Errorf("recordedByAccountId: %w", err)
				}
			}

			event.AddOldValue("recordedByAccountId", item.RecordedByAccountID)
			event.AddNewValue("recordedByAccountId", changes.RecordedByAccountID)

			item.RecordedByAccountID = changes.RecordedByAccountID
			newItem.RecordedByAccountID = changes.RecordedByAccountID
			changedFields = append(changedFields, "recorded_by_account_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("franchise_opening_records", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// FranchiseOpeningRecord - Delete
// ============================================================

// DeleteFranchiseOpeningRecordFunc 执行删除或恢复操作
func DeleteFranchiseOpeningRecordFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &FranchiseOpeningRecord{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("franchise_opening_records", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "FranchiseOpeningRecord",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("franchise_opening_records", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("franchise_opening_records", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("franchise_opening_records", ctx)).Updates(FranchiseOpeningRecord{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteFranchiseOpeningRecords 批量删除 FranchiseOpeningRecord 实体
func (r *GeneratedMutationResolver) DeleteFranchiseOpeningRecords(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteFranchiseOpeningRecords(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteFranchiseOpeningRecordsHandler 处理批量删除逻辑
func DeleteFranchiseOpeningRecordsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteFranchiseOpeningRecordFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// FranchiseOpeningRecord - Recovery
// ============================================================

// RecoveryFranchiseOpeningRecords 批量恢复 FranchiseOpeningRecord 实体
func (r *GeneratedMutationResolver) RecoveryFranchiseOpeningRecords(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryFranchiseOpeningRecords(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryFranchiseOpeningRecordsHandler 处理批量恢复逻辑
func RecoveryFranchiseOpeningRecordsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteFranchiseOpeningRecordFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// GlobalPaymentConfig - Create
// ============================================================

// CreateGlobalPaymentConfig 创建 GlobalPaymentConfig 实体的解析器入口
func (r *GeneratedMutationResolver) CreateGlobalPaymentConfig(ctx context.Context, input map[string]interface{}) (item *GlobalPaymentConfig, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateGlobalPaymentConfig(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateGlobalPaymentConfigHandler 处理 GlobalPaymentConfig 创建逻辑
func CreateGlobalPaymentConfigHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *GlobalPaymentConfig, err error) {
	item = &GlobalPaymentConfig{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "GlobalPaymentConfig",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes GlobalPaymentConfigChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["channel"]; ok {

		item.Channel = changes.Channel

		event.AddNewValue("channel", changes.Channel)
	}

	if _, ok := input["merchantId"]; ok && changes.MerchantID != nil {

		item.MerchantID = changes.MerchantID

		event.AddNewValue("merchantId", changes.MerchantID)
	}

	if _, ok := input["environment"]; ok && changes.Environment != nil {

		item.Environment = changes.Environment

		event.AddNewValue("environment", changes.Environment)
	}

	if _, ok := input["ratePpm"]; ok {

		item.RatePpm = changes.RatePpm

		event.AddNewValue("ratePpm", changes.RatePpm)
	}

	if _, ok := input["configState"]; ok {

		item.ConfigState = changes.ConfigState

		event.AddNewValue("configState", changes.ConfigState)
	}

	if _, ok := input["version"]; ok {

		item.Version = changes.Version

		event.AddNewValue("version", changes.Version)
	}

	if _, ok := input["keyId"]; ok && changes.KeyID != nil {

		item.KeyID = changes.KeyID

		event.AddNewValue("keyId", changes.KeyID)
	}

	if _, ok := input["credentialCiphertext"]; ok && changes.CredentialCiphertext != nil {

		item.CredentialCiphertext = changes.CredentialCiphertext

		event.AddNewValue("credentialCiphertext", changes.CredentialCiphertext)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("global_payment_configs", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// GlobalPaymentConfig - Update
// ============================================================

// UpdateGlobalPaymentConfig 更新 GlobalPaymentConfig 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateGlobalPaymentConfig(ctx context.Context, id string, input map[string]interface{}) (item *GlobalPaymentConfig, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateGlobalPaymentConfig(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateGlobalPaymentConfigHandler 处理 GlobalPaymentConfig 更新逻辑
func UpdateGlobalPaymentConfigHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *GlobalPaymentConfig, err error) {
	item = &GlobalPaymentConfig{}
	newItem := &GlobalPaymentConfig{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "GlobalPaymentConfig",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes GlobalPaymentConfigChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("global_payment_configs", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["channel"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Channel != changes.Channel {

			event.AddOldValue("channel", item.Channel)
			event.AddNewValue("channel", changes.Channel)

			item.Channel = changes.Channel
			newItem.Channel = changes.Channel
			changedFields = append(changedFields, "channel")
			isChange = true
		}
	}

	if _, ok := input["merchantId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.MerchantID != changes.MerchantID) && (item.MerchantID == nil || changes.MerchantID == nil || *item.MerchantID != *changes.MerchantID) {

			event.AddOldValue("merchantId", item.MerchantID)
			event.AddNewValue("merchantId", changes.MerchantID)

			item.MerchantID = changes.MerchantID
			newItem.MerchantID = changes.MerchantID
			changedFields = append(changedFields, "merchant_id")
			isChange = true
		}
	}

	if _, ok := input["environment"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Environment != changes.Environment) && (item.Environment == nil || changes.Environment == nil || *item.Environment != *changes.Environment) {

			event.AddOldValue("environment", item.Environment)
			event.AddNewValue("environment", changes.Environment)

			item.Environment = changes.Environment
			newItem.Environment = changes.Environment
			changedFields = append(changedFields, "environment")
			isChange = true
		}
	}

	if _, ok := input["ratePpm"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.RatePpm != changes.RatePpm {

			event.AddOldValue("ratePpm", item.RatePpm)
			event.AddNewValue("ratePpm", changes.RatePpm)

			item.RatePpm = changes.RatePpm
			newItem.RatePpm = changes.RatePpm
			changedFields = append(changedFields, "rate_ppm")
			isChange = true
		}
	}

	if _, ok := input["configState"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ConfigState != changes.ConfigState {

			event.AddOldValue("configState", item.ConfigState)
			event.AddNewValue("configState", changes.ConfigState)

			item.ConfigState = changes.ConfigState
			newItem.ConfigState = changes.ConfigState
			changedFields = append(changedFields, "config_state")
			isChange = true
		}
	}

	if _, ok := input["version"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Version != changes.Version {

			event.AddOldValue("version", item.Version)
			event.AddNewValue("version", changes.Version)

			item.Version = changes.Version
			newItem.Version = changes.Version
			changedFields = append(changedFields, "version")
			isChange = true
		}
	}

	if _, ok := input["keyId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.KeyID != changes.KeyID) && (item.KeyID == nil || changes.KeyID == nil || *item.KeyID != *changes.KeyID) {

			event.AddOldValue("keyId", item.KeyID)
			event.AddNewValue("keyId", changes.KeyID)

			item.KeyID = changes.KeyID
			newItem.KeyID = changes.KeyID
			changedFields = append(changedFields, "key_id")
			isChange = true
		}
	}

	if _, ok := input["credentialCiphertext"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.CredentialCiphertext != changes.CredentialCiphertext) && (item.CredentialCiphertext == nil || changes.CredentialCiphertext == nil || *item.CredentialCiphertext != *changes.CredentialCiphertext) {

			event.AddOldValue("credentialCiphertext", item.CredentialCiphertext)
			event.AddNewValue("credentialCiphertext", changes.CredentialCiphertext)

			item.CredentialCiphertext = changes.CredentialCiphertext
			newItem.CredentialCiphertext = changes.CredentialCiphertext
			changedFields = append(changedFields, "credential_ciphertext")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("global_payment_configs", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// GlobalPaymentConfig - Delete
// ============================================================

// DeleteGlobalPaymentConfigFunc 执行删除或恢复操作
func DeleteGlobalPaymentConfigFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &GlobalPaymentConfig{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("global_payment_configs", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "GlobalPaymentConfig",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("global_payment_configs", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("global_payment_configs", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("global_payment_configs", ctx)).Updates(GlobalPaymentConfig{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteGlobalPaymentConfigs 批量删除 GlobalPaymentConfig 实体
func (r *GeneratedMutationResolver) DeleteGlobalPaymentConfigs(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteGlobalPaymentConfigs(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteGlobalPaymentConfigsHandler 处理批量删除逻辑
func DeleteGlobalPaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteGlobalPaymentConfigFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// GlobalPaymentConfig - Recovery
// ============================================================

// RecoveryGlobalPaymentConfigs 批量恢复 GlobalPaymentConfig 实体
func (r *GeneratedMutationResolver) RecoveryGlobalPaymentConfigs(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryGlobalPaymentConfigs(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryGlobalPaymentConfigsHandler 处理批量恢复逻辑
func RecoveryGlobalPaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteGlobalPaymentConfigFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// FranchisePaymentConfig - Create
// ============================================================

// CreateFranchisePaymentConfig 创建 FranchisePaymentConfig 实体的解析器入口
func (r *GeneratedMutationResolver) CreateFranchisePaymentConfig(ctx context.Context, input map[string]interface{}) (item *FranchisePaymentConfig, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateFranchisePaymentConfig(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateFranchisePaymentConfigHandler 处理 FranchisePaymentConfig 创建逻辑
func CreateFranchisePaymentConfigHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *FranchisePaymentConfig, err error) {
	item = &FranchisePaymentConfig{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "FranchisePaymentConfig",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes FranchisePaymentConfigChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["channel"]; ok {

		item.Channel = changes.Channel

		event.AddNewValue("channel", changes.Channel)
	}

	if _, ok := input["merchantId"]; ok && changes.MerchantID != nil {

		item.MerchantID = changes.MerchantID

		event.AddNewValue("merchantId", changes.MerchantID)
	}

	if _, ok := input["environment"]; ok && changes.Environment != nil {

		item.Environment = changes.Environment

		event.AddNewValue("environment", changes.Environment)
	}

	if _, ok := input["ratePpm"]; ok {

		item.RatePpm = changes.RatePpm

		event.AddNewValue("ratePpm", changes.RatePpm)
	}

	if _, ok := input["configState"]; ok {

		item.ConfigState = changes.ConfigState

		event.AddNewValue("configState", changes.ConfigState)
	}

	if _, ok := input["version"]; ok {

		item.Version = changes.Version

		event.AddNewValue("version", changes.Version)
	}

	if _, ok := input["keyId"]; ok && changes.KeyID != nil {

		item.KeyID = changes.KeyID

		event.AddNewValue("keyId", changes.KeyID)
	}

	if _, ok := input["credentialCiphertext"]; ok && changes.CredentialCiphertext != nil {

		item.CredentialCiphertext = changes.CredentialCiphertext

		event.AddNewValue("credentialCiphertext", changes.CredentialCiphertext)
	}

	if _, ok := input["organizationId"]; ok {

		if !utils.IsNil(input["organizationId"]) {
			if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
				return nil, fmt.Errorf("organizationId: %w", err)
			}
		}

		item.OrganizationID = changes.OrganizationID

		event.AddNewValue("organizationId", changes.OrganizationID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("franchise_payment_configs", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// FranchisePaymentConfig - Update
// ============================================================

// UpdateFranchisePaymentConfig 更新 FranchisePaymentConfig 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateFranchisePaymentConfig(ctx context.Context, id string, input map[string]interface{}) (item *FranchisePaymentConfig, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateFranchisePaymentConfig(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateFranchisePaymentConfigHandler 处理 FranchisePaymentConfig 更新逻辑
func UpdateFranchisePaymentConfigHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *FranchisePaymentConfig, err error) {
	item = &FranchisePaymentConfig{}
	newItem := &FranchisePaymentConfig{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "FranchisePaymentConfig",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes FranchisePaymentConfigChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["organization"]) && !utils.IsNil(input["organizationId"]) {
		return nil, fmt.Errorf("organizationId and organization cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("franchise_payment_configs", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["channel"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Channel != changes.Channel {

			event.AddOldValue("channel", item.Channel)
			event.AddNewValue("channel", changes.Channel)

			item.Channel = changes.Channel
			newItem.Channel = changes.Channel
			changedFields = append(changedFields, "channel")
			isChange = true
		}
	}

	if _, ok := input["merchantId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.MerchantID != changes.MerchantID) && (item.MerchantID == nil || changes.MerchantID == nil || *item.MerchantID != *changes.MerchantID) {

			event.AddOldValue("merchantId", item.MerchantID)
			event.AddNewValue("merchantId", changes.MerchantID)

			item.MerchantID = changes.MerchantID
			newItem.MerchantID = changes.MerchantID
			changedFields = append(changedFields, "merchant_id")
			isChange = true
		}
	}

	if _, ok := input["environment"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Environment != changes.Environment) && (item.Environment == nil || changes.Environment == nil || *item.Environment != *changes.Environment) {

			event.AddOldValue("environment", item.Environment)
			event.AddNewValue("environment", changes.Environment)

			item.Environment = changes.Environment
			newItem.Environment = changes.Environment
			changedFields = append(changedFields, "environment")
			isChange = true
		}
	}

	if _, ok := input["ratePpm"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.RatePpm != changes.RatePpm {

			event.AddOldValue("ratePpm", item.RatePpm)
			event.AddNewValue("ratePpm", changes.RatePpm)

			item.RatePpm = changes.RatePpm
			newItem.RatePpm = changes.RatePpm
			changedFields = append(changedFields, "rate_ppm")
			isChange = true
		}
	}

	if _, ok := input["configState"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ConfigState != changes.ConfigState {

			event.AddOldValue("configState", item.ConfigState)
			event.AddNewValue("configState", changes.ConfigState)

			item.ConfigState = changes.ConfigState
			newItem.ConfigState = changes.ConfigState
			changedFields = append(changedFields, "config_state")
			isChange = true
		}
	}

	if _, ok := input["version"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Version != changes.Version {

			event.AddOldValue("version", item.Version)
			event.AddNewValue("version", changes.Version)

			item.Version = changes.Version
			newItem.Version = changes.Version
			changedFields = append(changedFields, "version")
			isChange = true
		}
	}

	if _, ok := input["keyId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.KeyID != changes.KeyID) && (item.KeyID == nil || changes.KeyID == nil || *item.KeyID != *changes.KeyID) {

			event.AddOldValue("keyId", item.KeyID)
			event.AddNewValue("keyId", changes.KeyID)

			item.KeyID = changes.KeyID
			newItem.KeyID = changes.KeyID
			changedFields = append(changedFields, "key_id")
			isChange = true
		}
	}

	if _, ok := input["credentialCiphertext"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.CredentialCiphertext != changes.CredentialCiphertext) && (item.CredentialCiphertext == nil || changes.CredentialCiphertext == nil || *item.CredentialCiphertext != *changes.CredentialCiphertext) {

			event.AddOldValue("credentialCiphertext", item.CredentialCiphertext)
			event.AddNewValue("credentialCiphertext", changes.CredentialCiphertext)

			item.CredentialCiphertext = changes.CredentialCiphertext
			newItem.CredentialCiphertext = changes.CredentialCiphertext
			changedFields = append(changedFields, "credential_ciphertext")
			isChange = true
		}
	}

	if _, ok := input["organizationId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.OrganizationID != changes.OrganizationID {

			if !utils.IsNil(input["organizationId"]) {
				if err := tx.Select("id").Where("id = ?", input["organizationId"]).First(&Organization{}).Error; err != nil {
					return nil, fmt.Errorf("organizationId: %w", err)
				}
			}

			event.AddOldValue("organizationId", item.OrganizationID)
			event.AddNewValue("organizationId", changes.OrganizationID)

			item.OrganizationID = changes.OrganizationID
			newItem.OrganizationID = changes.OrganizationID
			changedFields = append(changedFields, "organization_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("franchise_payment_configs", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// FranchisePaymentConfig - Delete
// ============================================================

// DeleteFranchisePaymentConfigFunc 执行删除或恢复操作
func DeleteFranchisePaymentConfigFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &FranchisePaymentConfig{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("franchise_payment_configs", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "FranchisePaymentConfig",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("franchise_payment_configs", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("franchise_payment_configs", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("franchise_payment_configs", ctx)).Updates(FranchisePaymentConfig{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteFranchisePaymentConfigs 批量删除 FranchisePaymentConfig 实体
func (r *GeneratedMutationResolver) DeleteFranchisePaymentConfigs(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteFranchisePaymentConfigs(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteFranchisePaymentConfigsHandler 处理批量删除逻辑
func DeleteFranchisePaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteFranchisePaymentConfigFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// FranchisePaymentConfig - Recovery
// ============================================================

// RecoveryFranchisePaymentConfigs 批量恢复 FranchisePaymentConfig 实体
func (r *GeneratedMutationResolver) RecoveryFranchisePaymentConfigs(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryFranchisePaymentConfigs(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryFranchisePaymentConfigsHandler 处理批量恢复逻辑
func RecoveryFranchisePaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteFranchisePaymentConfigFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// StorePaymentConfig - Create
// ============================================================

// CreateStorePaymentConfig 创建 StorePaymentConfig 实体的解析器入口
func (r *GeneratedMutationResolver) CreateStorePaymentConfig(ctx context.Context, input map[string]interface{}) (item *StorePaymentConfig, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.CreateStorePaymentConfig(ctx, r.GeneratedResolver, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// CreateStorePaymentConfigHandler 处理 StorePaymentConfig 创建逻辑
func CreateStorePaymentConfigHandler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *StorePaymentConfig, err error) {
	item = &StorePaymentConfig{}
	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeCreated,
		Entity:      "StorePaymentConfig",
		EntityID:    item.ID,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes StorePaymentConfigChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// 设置基础字段
	item.ID = uuid.Must(uuid.NewV4()).String()
	item.CreatedAt = timestampMillis
	item.CreatedBy = principalID

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["store"]) && !utils.IsNil(input["storeId"]) {
		return nil, fmt.Errorf("storeId and store cannot coexist")
	}

	// ========== 处理 ManyToOne/OneToOne 关系（当前实体持有外键） ==========

	// ========== 处理普通字段 ==========

	if _, ok := input["channel"]; ok {

		item.Channel = changes.Channel

		event.AddNewValue("channel", changes.Channel)
	}

	if _, ok := input["merchantId"]; ok && changes.MerchantID != nil {

		item.MerchantID = changes.MerchantID

		event.AddNewValue("merchantId", changes.MerchantID)
	}

	if _, ok := input["environment"]; ok && changes.Environment != nil {

		item.Environment = changes.Environment

		event.AddNewValue("environment", changes.Environment)
	}

	if _, ok := input["ratePpm"]; ok {

		item.RatePpm = changes.RatePpm

		event.AddNewValue("ratePpm", changes.RatePpm)
	}

	if _, ok := input["configState"]; ok {

		item.ConfigState = changes.ConfigState

		event.AddNewValue("configState", changes.ConfigState)
	}

	if _, ok := input["version"]; ok {

		item.Version = changes.Version

		event.AddNewValue("version", changes.Version)
	}

	if _, ok := input["keyId"]; ok && changes.KeyID != nil {

		item.KeyID = changes.KeyID

		event.AddNewValue("keyId", changes.KeyID)
	}

	if _, ok := input["credentialCiphertext"]; ok && changes.CredentialCiphertext != nil {

		item.CredentialCiphertext = changes.CredentialCiphertext

		event.AddNewValue("credentialCiphertext", changes.CredentialCiphertext)
	}

	if _, ok := input["storeId"]; ok {

		if !utils.IsNil(input["storeId"]) {
			if err := tx.Select("id").Where("id = ?", input["storeId"]).First(&Store{}).Error; err != nil {
				return nil, fmt.Errorf("storeId: %w", err)
			}
		}

		item.StoreID = changes.StoreID

		event.AddNewValue("storeId", changes.StoreID)
	}

	if _, ok := input["isDelete"]; ok && changes.IsDelete != nil {

		item.IsDelete = changes.IsDelete

		event.AddNewValue("isDelete", changes.IsDelete)
	}

	if _, ok := input["weight"]; ok && changes.Weight != nil {

		item.Weight = changes.Weight

		event.AddNewValue("weight", changes.Weight)
	}

	if _, ok := input["state"]; ok && changes.State != nil {

		item.State = changes.State

		event.AddNewValue("state", changes.State)
	}

	// ========== 保存主实体 ==========
	if err := tx.Omit(clause.Associations).Table(TableName("store_payment_configs", ctx)).Create(item).Error; err != nil {
		return item, err
	}

	// ========== 处理 OneToMany/ManyToMany 关系（关联表持有外键或中间表） ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// StorePaymentConfig - Update
// ============================================================

// UpdateStorePaymentConfig 更新 StorePaymentConfig 实体的解析器入口
func (r *GeneratedMutationResolver) UpdateStorePaymentConfig(ctx context.Context, id string, input map[string]interface{}) (item *StorePaymentConfig, err error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err = r.Handlers.UpdateStorePaymentConfig(ctx, r.GeneratedResolver, id, input)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return
}

// UpdateStorePaymentConfigHandler 处理 StorePaymentConfig 更新逻辑
func UpdateStorePaymentConfigHandler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *StorePaymentConfig, err error) {
	item = &StorePaymentConfig{}
	newItem := &StorePaymentConfig{}
	isChange := false

	now := time.Now()
	timestampMillis := now.UnixNano() / 1e6
	principalID := GetPrincipalIDFromContext(ctx)
	tx := GetTransaction(ctx)

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeUpdated,
		Entity:      "StorePaymentConfig",
		EntityID:    id,
		Date:        timestampMillis,
		PrincipalID: principalID,
	})

	// 解析输入变更
	var changes StorePaymentConfigChanges
	if err = ApplyChanges(input, &changes); err != nil {
		return
	}

	// 验证必填字段
	if err = CheckStructFieldIsEmpty(item, input); err != nil {
		return nil, err
	}

	// ========== 验证关系字段冲突 ==========

	if !utils.IsNil(input["store"]) && !utils.IsNil(input["storeId"]) {
		return nil, fmt.Errorf("storeId and store cannot coexist")
	}

	// 获取现有实体
	if err = GetItem(ctx, tx, TableName("store_payment_configs", ctx), item, &id); err != nil {
		return nil, err
	}

	// 设置审计字段
	newItem.UpdatedAt = &timestampMillis
	newItem.UpdatedBy = principalID

	// 字段变更追踪
	changedFields := []string{}

	// ========== 处理 ManyToOne/OneToOne 关系 ==========

	// ========== 处理普通字段 ==========
	// changedFields := []string{} (Moved to top)

	if _, ok := input["id"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ID != changes.ID {

			event.AddOldValue("id", item.ID)
			event.AddNewValue("id", changes.ID)

			item.ID = changes.ID
			newItem.ID = changes.ID
			changedFields = append(changedFields, "id")
			isChange = true
		}
	}

	if _, ok := input["channel"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Channel != changes.Channel {

			event.AddOldValue("channel", item.Channel)
			event.AddNewValue("channel", changes.Channel)

			item.Channel = changes.Channel
			newItem.Channel = changes.Channel
			changedFields = append(changedFields, "channel")
			isChange = true
		}
	}

	if _, ok := input["merchantId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.MerchantID != changes.MerchantID) && (item.MerchantID == nil || changes.MerchantID == nil || *item.MerchantID != *changes.MerchantID) {

			event.AddOldValue("merchantId", item.MerchantID)
			event.AddNewValue("merchantId", changes.MerchantID)

			item.MerchantID = changes.MerchantID
			newItem.MerchantID = changes.MerchantID
			changedFields = append(changedFields, "merchant_id")
			isChange = true
		}
	}

	if _, ok := input["environment"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Environment != changes.Environment) && (item.Environment == nil || changes.Environment == nil || *item.Environment != *changes.Environment) {

			event.AddOldValue("environment", item.Environment)
			event.AddNewValue("environment", changes.Environment)

			item.Environment = changes.Environment
			newItem.Environment = changes.Environment
			changedFields = append(changedFields, "environment")
			isChange = true
		}
	}

	if _, ok := input["ratePpm"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.RatePpm != changes.RatePpm {

			event.AddOldValue("ratePpm", item.RatePpm)
			event.AddNewValue("ratePpm", changes.RatePpm)

			item.RatePpm = changes.RatePpm
			newItem.RatePpm = changes.RatePpm
			changedFields = append(changedFields, "rate_ppm")
			isChange = true
		}
	}

	if _, ok := input["configState"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.ConfigState != changes.ConfigState {

			event.AddOldValue("configState", item.ConfigState)
			event.AddNewValue("configState", changes.ConfigState)

			item.ConfigState = changes.ConfigState
			newItem.ConfigState = changes.ConfigState
			changedFields = append(changedFields, "config_state")
			isChange = true
		}
	}

	if _, ok := input["version"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.Version != changes.Version {

			event.AddOldValue("version", item.Version)
			event.AddNewValue("version", changes.Version)

			item.Version = changes.Version
			newItem.Version = changes.Version
			changedFields = append(changedFields, "version")
			isChange = true
		}
	}

	if _, ok := input["keyId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.KeyID != changes.KeyID) && (item.KeyID == nil || changes.KeyID == nil || *item.KeyID != *changes.KeyID) {

			event.AddOldValue("keyId", item.KeyID)
			event.AddNewValue("keyId", changes.KeyID)

			item.KeyID = changes.KeyID
			newItem.KeyID = changes.KeyID
			changedFields = append(changedFields, "key_id")
			isChange = true
		}
	}

	if _, ok := input["credentialCiphertext"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.CredentialCiphertext != changes.CredentialCiphertext) && (item.CredentialCiphertext == nil || changes.CredentialCiphertext == nil || *item.CredentialCiphertext != *changes.CredentialCiphertext) {

			event.AddOldValue("credentialCiphertext", item.CredentialCiphertext)
			event.AddNewValue("credentialCiphertext", changes.CredentialCiphertext)

			item.CredentialCiphertext = changes.CredentialCiphertext
			newItem.CredentialCiphertext = changes.CredentialCiphertext
			changedFields = append(changedFields, "credential_ciphertext")
			isChange = true
		}
	}

	if _, ok := input["storeId"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if item.StoreID != changes.StoreID {

			if !utils.IsNil(input["storeId"]) {
				if err := tx.Select("id").Where("id = ?", input["storeId"]).First(&Store{}).Error; err != nil {
					return nil, fmt.Errorf("storeId: %w", err)
				}
			}

			event.AddOldValue("storeId", item.StoreID)
			event.AddNewValue("storeId", changes.StoreID)

			item.StoreID = changes.StoreID
			newItem.StoreID = changes.StoreID
			changedFields = append(changedFields, "store_id")
			isChange = true
		}
	}

	if _, ok := input["isDelete"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.IsDelete != changes.IsDelete) && (item.IsDelete == nil || changes.IsDelete == nil || *item.IsDelete != *changes.IsDelete) {

			event.AddOldValue("isDelete", item.IsDelete)
			event.AddNewValue("isDelete", changes.IsDelete)

			item.IsDelete = changes.IsDelete
			newItem.IsDelete = changes.IsDelete
			changedFields = append(changedFields, "is_delete")
			isChange = true
		}
	}

	if _, ok := input["weight"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.Weight != changes.Weight) && (item.Weight == nil || changes.Weight == nil || *item.Weight != *changes.Weight) {

			event.AddOldValue("weight", item.Weight)
			event.AddNewValue("weight", changes.Weight)

			item.Weight = changes.Weight
			newItem.Weight = changes.Weight
			changedFields = append(changedFields, "weight")
			isChange = true
		}
	}

	if _, ok := input["state"]; ok {
		// 只要 input 中包含该字段，且值发生了变化（包括变为 null），就进行更新
		if (item.State != changes.State) && (item.State == nil || changes.State == nil || *item.State != *changes.State) {

			event.AddOldValue("state", item.State)
			event.AddNewValue("state", changes.State)

			item.State = changes.State
			newItem.State = changes.State
			changedFields = append(changedFields, "state")
			isChange = true
		}
	}

	// ========== 保存主实体变更 ==========
	if isChange {
		changedFields = append(changedFields, "updated_at", "updated_by")

		if err := tx.Table(TableName("store_payment_configs", ctx)).Where("id = ?", id).Select(changedFields).Updates(newItem).Error; err != nil {
			return item, err
		}
	}

	// ========== 处理 OneToMany/ManyToMany 关系 ==========

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// ============================================================
// StorePaymentConfig - Delete
// ============================================================

// DeleteStorePaymentConfigFunc 执行删除或恢复操作
func DeleteStorePaymentConfigFunc(ctx context.Context, r *GeneratedResolver, id string, operationType string, unscoped *bool) (err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item := &StorePaymentConfig{}
	now := time.Now()
	tx := GetTransaction(ctx)

	// 检查主从关系约束

	// 确定操作类型
	var status int64 = 1
	var isDelete int64 = 2
	if operationType == "recovery" {
		isDelete = 1
		status = 2
	}

	// 获取现有实体
	if err = tx.Unscoped().Table(TableName("store_payment_configs", ctx)).Where("is_delete = ? and id = ?", status, id).First(item).Error; err != nil {
		return err
	}

	deletedAt := now.UnixNano() / 1e6

	// 创建事件记录
	event := NewEvent(EventMetadata{
		Type:        EventTypeDeleted,
		Entity:      "StorePaymentConfig",
		EntityID:    id,
		Date:        deletedAt,
		PrincipalID: principalID,
	})

	// 执行删除或恢复
	if operationType == "recovery" {
		if err := tx.Unscoped().Table(TableName("store_payment_configs", ctx)).Model(&item).Updates(map[string]interface{}{
			"IsDelete":  1,
			"DeletedAt": nil,
			"DeletedBy": nil,
		}).Error; err != nil {
			return err
		}
	} else {
		if unscoped != nil && *unscoped {
			// 物理删除
			if err := tx.Unscoped().Table(TableName("store_payment_configs", ctx)).Model(&item).Delete(item).Error; err != nil {
				return err
			}
		} else {
			// 软删除
			if err := tx.Model(&item).Table(TableName("store_payment_configs", ctx)).Updates(StorePaymentConfig{
				IsDelete:  &isDelete,
				DeletedAt: &deletedAt,
				DeletedBy: principalID,
				UpdatedBy: principalID,
			}).Error; err != nil {
				return err
			}
		}
	}

	// 记录事件
	if len(event.Changes) > 0 {
		AddMutationEvent(ctx, event)
	}

	return
}

// DeleteStorePaymentConfigs 批量删除 StorePaymentConfig 实体
func (r *GeneratedMutationResolver) DeleteStorePaymentConfigs(ctx context.Context, id []string, unscoped *bool) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteStorePaymentConfigs(ctx, r.GeneratedResolver, id, unscoped)
	if err != nil {
		RollbackMutationContext(ctx, r.GeneratedResolver)
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// DeleteStorePaymentConfigsHandler 处理批量删除逻辑
func DeleteStorePaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, id []string, unscoped *bool) (bool, error) {
	for _, itemID := range id {
		if err := DeleteStorePaymentConfigFunc(ctx, r, itemID, "delete", unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}

// ============================================================
// StorePaymentConfig - Recovery
// ============================================================

// RecoveryStorePaymentConfigs 批量恢复 StorePaymentConfig 实体
func (r *GeneratedMutationResolver) RecoveryStorePaymentConfigs(ctx context.Context, id []string) (bool, error) {
	ctx = EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.RecoveryStorePaymentConfigs(ctx, r.GeneratedResolver, id)
	if err != nil {
		return done, err
	}
	err = FinishMutationContext(ctx, r.GeneratedResolver)
	return done, err
}

// RecoveryStorePaymentConfigsHandler 处理批量恢复逻辑
func RecoveryStorePaymentConfigsHandler(ctx context.Context, r *GeneratedResolver, id []string) (bool, error) {
	unscoped := false
	for _, itemID := range id {
		if err := DeleteStorePaymentConfigFunc(ctx, r, itemID, "recovery", &unscoped); err != nil {
			return false, err
		}
	}
	return true, nil
}
