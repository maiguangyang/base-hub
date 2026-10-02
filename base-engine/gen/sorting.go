package gen

import (
	"context"
)

func (s AccountSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("accounts", ctx), sorts, joins)
}
func (s AccountSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Phone != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("phone")+" "+s.Phone.String())
	}

	if s.DisplayName != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("displayName")+" "+s.DisplayName.String())
	}

	if s.Email != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("email")+" "+s.Email.String())
	}

	if s.MustChangePassword != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("mustChangePassword")+" "+s.MustChangePassword.String())
	}

	if s.CredentialVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("credentialVersion")+" "+s.CredentialVersion.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Memberships != nil {
		_alias := alias + "_memberships"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"."+"account_id"+" = "+alias+".id")
		err := s.Memberships.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.InitializedOrganizations != nil {
		_alias := alias + "_initializedOrganizations"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+"."+"initial_account_id"+" = "+alias+".id")
		err := s.InitializedOrganizations.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.OpeningRecords != nil {
		_alias := alias + "_openingRecords"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"initial_account_id"+" = "+alias+".id")
		err := s.OpeningRecords.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.RecordedOpeningRecords != nil {
		_alias := alias + "_recordedOpeningRecords"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"recorded_by_account_id"+" = "+alias+".id")
		err := s.RecordedOpeningRecords.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Sessions != nil {
		_alias := alias + "_sessions"
		*joins = append(*joins, "LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+"."+"account_id"+" = "+alias+".id")
		err := s.Sessions.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ReviewedStores != nil {
		_alias := alias + "_reviewedStores"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"."+"reviewed_by_account_id"+" = "+alias+".id")
		err := s.ReviewedStores.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.SentMembershipInvitations != nil {
		_alias := alias + "_sentMembershipInvitations"
		*joins = append(*joins, "LEFT JOIN "+TableName("membership_invitations", ctx)+" "+_alias+" ON "+_alias+"."+"invited_by_account_id"+" = "+alias+".id")
		err := s.SentMembershipInvitations.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"actor_account_id"+" = "+alias+".id")
		err := s.AuditLogs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CreatedStocktakes != nil {
		_alias := alias + "_createdStocktakes"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+"."+"initiated_by_account_id"+" = "+alias+".id")
		err := s.CreatedStocktakes.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.PostedStocktakes != nil {
		_alias := alias + "_postedStocktakes"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+"."+"posted_by_id"+" = "+alias+".id")
		err := s.PostedStocktakes.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s OrganizationSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("organizations", ctx), sorts, joins)
}
func (s OrganizationSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Code != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("code")+" "+s.Code.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.SuspendedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("suspendedAt")+" "+s.SuspendedAt.String())
	}

	if s.SuspensionReasonCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("suspensionReasonCode")+" "+s.SuspensionReasonCode.String())
	}

	if s.InitialAccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("initialAccountId")+" "+s.InitialAccountID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Memberships != nil {
		_alias := alias + "_memberships"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.Memberships.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.InitialAccount != nil {
		_alias := alias + "_initialAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initial_account_id")
		err := s.InitialAccount.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.OpeningRecords != nil {
		_alias := alias + "_openingRecords"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.OpeningRecords.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Stores != nil {
		_alias := alias + "_stores"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.Stores.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Roles != nil {
		_alias := alias + "_roles"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.Roles.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Sessions != nil {
		_alias := alias + "_sessions"
		*joins = append(*joins, "LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.Sessions.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.AuditLogs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.PaymentConfigs != nil {
		_alias := alias + "_paymentConfigs"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_payment_configs", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.PaymentConfigs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CustomerMembers != nil {
		_alias := alias + "_customerMembers"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.CustomerMembers.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CustomerBenefitPolicies != nil {
		_alias := alias + "_customerBenefitPolicies"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_benefit_policies", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.CustomerBenefitPolicies.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CustomerDailyPointGrantBudgets != nil {
		_alias := alias + "_customerDailyPointGrantBudgets"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_daily_point_grant_budgets", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.CustomerDailyPointGrantBudgets.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CustomerPointEntries != nil {
		_alias := alias + "_customerPointEntries"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+"."+"source_organization_id"+" = "+alias+".id")
		err := s.CustomerPointEntries.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CustomerCouponTemplates != nil {
		_alias := alias + "_customerCouponTemplates"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.CustomerCouponTemplates.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ProductCategories != nil {
		_alias := alias + "_productCategories"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.ProductCategories.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ProductBrands != nil {
		_alias := alias + "_productBrands"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_brands", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.ProductBrands.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.SpecificationDefinitions != nil {
		_alias := alias + "_specificationDefinitions"
		*joins = append(*joins, "LEFT JOIN "+TableName("specification_definitions", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.SpecificationDefinitions.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ProductPackageTemplates != nil {
		_alias := alias + "_productPackageTemplates"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.ProductPackageTemplates.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Products != nil {
		_alias := alias + "_products"
		*joins = append(*joins, "LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := s.Products.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s OperatorMembershipSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("operator_memberships", ctx), sorts, joins)
}
func (s OperatorMembershipSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.InvitedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("invitedAt")+" "+s.InvitedAt.String())
	}

	if s.AcceptedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("acceptedAt")+" "+s.AcceptedAt.String())
	}

	if s.AccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("accountId")+" "+s.AccountID.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Account != nil {
		_alias := alias + "_account"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"account_id")
		err := s.Account.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Roles != nil {
		_alias := alias + "_roles"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"member_id"+" LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"role_id"+" = "+_alias+".id")
		err := s.Roles.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Stores != nil {
		_alias := alias + "_stores"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_stores", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"member_id"+" LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"store_id"+" = "+_alias+".id")
		err := s.Stores.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Invitations != nil {
		_alias := alias + "_invitations"
		*joins = append(*joins, "LEFT JOIN "+TableName("membership_invitations", ctx)+" "+_alias+" ON "+_alias+"."+"membership_id"+" = "+alias+".id")
		err := s.Invitations.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s PermissionSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("permissions", ctx), sorts, joins)
}
func (s PermissionSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.Action != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("action")+" "+s.Action.String())
	}

	if s.Module != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("module")+" "+s.Module.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Roles != nil {
		_alias := alias + "_roles"
		*joins = append(*joins, "LEFT JOIN "+TableName("permission_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"permission_id"+" LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"role_id"+" = "+_alias+".id")
		err := s.Roles.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s OperatorRoleSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("operator_roles", ctx), sorts, joins)
}
func (s OperatorRoleSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Members != nil {
		_alias := alias + "_members"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"role_id"+" LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"member_id"+" = "+_alias+".id")
		err := s.Members.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Permissions != nil {
		_alias := alias + "_permissions"
		*joins = append(*joins, "LEFT JOIN "+TableName("permission_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"role_id"+" LEFT JOIN "+TableName("permissions", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"permission_id"+" = "+_alias+".id")
		err := s.Permissions.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StoreSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("stores", ctx), sorts, joins)
}
func (s StoreSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Code != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("code")+" "+s.Code.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.SubmittedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("submittedAt")+" "+s.SubmittedAt.String())
	}

	if s.ReviewedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reviewedAt")+" "+s.ReviewedAt.String())
	}

	if s.RejectionReason != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("rejectionReason")+" "+s.RejectionReason.String())
	}

	if s.ContactPhone != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("contactPhone")+" "+s.ContactPhone.String())
	}

	if s.ManagerName != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("managerName")+" "+s.ManagerName.String())
	}

	if s.ManagerPhone != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("managerPhone")+" "+s.ManagerPhone.String())
	}

	if s.Province != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("province")+" "+s.Province.String())
	}

	if s.City != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("city")+" "+s.City.String())
	}

	if s.District != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("district")+" "+s.District.String())
	}

	if s.Address != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("address")+" "+s.Address.String())
	}

	if s.BusinessHours != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("businessHours")+" "+s.BusinessHours.String())
	}

	if s.SupportDineIn != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("supportDineIn")+" "+s.SupportDineIn.String())
	}

	if s.SupportTakeout != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("supportTakeout")+" "+s.SupportTakeout.String())
	}

	if s.StoreArea != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storeArea")+" "+s.StoreArea.String())
	}

	if s.TableCount != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("tableCount")+" "+s.TableCount.String())
	}

	if s.ReceiptFooter != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("receiptFooter")+" "+s.ReceiptFooter.String())
	}

	if s.BusinessLicenseImageURL != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("businessLicenseImageUrl")+" "+s.BusinessLicenseImageURL.String())
	}

	if s.OtherDocumentImageURL != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("otherDocumentImageUrl")+" "+s.OtherDocumentImageURL.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.ReviewedByAccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reviewedByAccountId")+" "+s.ReviewedByAccountID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Members != nil {
		_alias := alias + "_members"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_stores", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"store_id"+" LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"member_id"+" = "+_alias+".id")
		err := s.Members.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ReviewedByAccount != nil {
		_alias := alias + "_reviewedByAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"reviewed_by_account_id")
		err := s.ReviewedByAccount.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := s.AuditLogs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.PaymentConfigs != nil {
		_alias := alias + "_paymentConfigs"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_payment_configs", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := s.PaymentConfigs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ProductListings != nil {
		_alias := alias + "_productListings"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := s.ProductListings.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.StockMovements != nil {
		_alias := alias + "_stockMovements"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := s.StockMovements.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Stocktakes != nil {
		_alias := alias + "_stocktakes"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := s.Stocktakes.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Promotions != nil {
		_alias := alias + "_promotions"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_promotions", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := s.Promotions.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CouponTemplates != nil {
		_alias := alias + "_couponTemplates"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+"."+"applicable_store_id"+" = "+alias+".id")
		err := s.CouponTemplates.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SessionSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("sessions", ctx), sorts, joins)
}
func (s SessionSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.CredentialVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("credentialVersion")+" "+s.CredentialVersion.String())
	}

	if s.ExpiresAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("expiresAt")+" "+s.ExpiresAt.String())
	}

	if s.RevokedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("revokedAt")+" "+s.RevokedAt.String())
	}

	if s.RevocationCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("revocationCode")+" "+s.RevocationCode.String())
	}

	if s.LastSeenAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("lastSeenAt")+" "+s.LastSeenAt.String())
	}

	if s.AccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("accountId")+" "+s.AccountID.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Account != nil {
		_alias := alias + "_account"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"account_id")
		err := s.Account.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"session_id"+" = "+alias+".id")
		err := s.AuditLogs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s MembershipInvitationSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("membership_invitations", ctx), sorts, joins)
}
func (s MembershipInvitationSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.ExpiresAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("expiresAt")+" "+s.ExpiresAt.String())
	}

	if s.AcceptedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("acceptedAt")+" "+s.AcceptedAt.String())
	}

	if s.RevokedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("revokedAt")+" "+s.RevokedAt.String())
	}

	if s.MembershipID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("membershipId")+" "+s.MembershipID.String())
	}

	if s.InvitedByAccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("invitedByAccountId")+" "+s.InvitedByAccountID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Membership != nil {
		_alias := alias + "_membership"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"membership_id")
		err := s.Membership.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.InvitedByAccount != nil {
		_alias := alias + "_invitedByAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"invited_by_account_id")
		err := s.InvitedByAccount.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s AuditLogSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("audit_logs", ctx), sorts, joins)
}
func (s AuditLogSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Action != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("action")+" "+s.Action.String())
	}

	if s.ResourceType != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("resourceType")+" "+s.ResourceType.String())
	}

	if s.ResourceID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("resourceId")+" "+s.ResourceID.String())
	}

	if s.ResultCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("resultCode")+" "+s.ResultCode.String())
	}

	if s.MetadataJSON != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("metadataJson")+" "+s.MetadataJSON.String())
	}

	if s.ActorAccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("actorAccountId")+" "+s.ActorAccountID.String())
	}

	if s.SessionID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("sessionId")+" "+s.SessionID.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.StoreID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storeId")+" "+s.StoreID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.ActorAccount != nil {
		_alias := alias + "_actorAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"actor_account_id")
		err := s.ActorAccount.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Session != nil {
		_alias := alias + "_session"
		*joins = append(*joins, "LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"session_id")
		err := s.Session.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := s.Store.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s FranchiseOpeningRecordSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("franchise_opening_records", ctx), sorts, joins)
}
func (s FranchiseOpeningRecordSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.RecordNumber != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("recordNumber")+" "+s.RecordNumber.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.InitialAccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("initialAccountId")+" "+s.InitialAccountID.String())
	}

	if s.RecordedByAccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("recordedByAccountId")+" "+s.RecordedByAccountID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.InitialAccount != nil {
		_alias := alias + "_initialAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initial_account_id")
		err := s.InitialAccount.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.RecordedByAccount != nil {
		_alias := alias + "_recordedByAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"recorded_by_account_id")
		err := s.RecordedByAccount.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s GlobalPaymentConfigSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("global_payment_configs", ctx), sorts, joins)
}
func (s GlobalPaymentConfigSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Channel != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("channel")+" "+s.Channel.String())
	}

	if s.MerchantID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("merchantId")+" "+s.MerchantID.String())
	}

	if s.Environment != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("environment")+" "+s.Environment.String())
	}

	if s.RatePpm != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("ratePpm")+" "+s.RatePpm.String())
	}

	if s.ConfigState != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("configState")+" "+s.ConfigState.String())
	}

	if s.Version != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("version")+" "+s.Version.String())
	}

	if s.KeyID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("keyId")+" "+s.KeyID.String())
	}

	if s.CredentialCiphertext != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("credentialCiphertext")+" "+s.CredentialCiphertext.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	return nil
}

func (s FranchisePaymentConfigSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("franchise_payment_configs", ctx), sorts, joins)
}
func (s FranchisePaymentConfigSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Channel != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("channel")+" "+s.Channel.String())
	}

	if s.MerchantID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("merchantId")+" "+s.MerchantID.String())
	}

	if s.Environment != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("environment")+" "+s.Environment.String())
	}

	if s.RatePpm != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("ratePpm")+" "+s.RatePpm.String())
	}

	if s.ConfigState != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("configState")+" "+s.ConfigState.String())
	}

	if s.Version != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("version")+" "+s.Version.String())
	}

	if s.KeyID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("keyId")+" "+s.KeyID.String())
	}

	if s.CredentialCiphertext != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("credentialCiphertext")+" "+s.CredentialCiphertext.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StorePaymentConfigSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_payment_configs", ctx), sorts, joins)
}
func (s StorePaymentConfigSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Channel != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("channel")+" "+s.Channel.String())
	}

	if s.MerchantID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("merchantId")+" "+s.MerchantID.String())
	}

	if s.Environment != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("environment")+" "+s.Environment.String())
	}

	if s.RatePpm != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("ratePpm")+" "+s.RatePpm.String())
	}

	if s.ConfigState != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("configState")+" "+s.ConfigState.String())
	}

	if s.Version != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("version")+" "+s.Version.String())
	}

	if s.KeyID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("keyId")+" "+s.KeyID.String())
	}

	if s.CredentialCiphertext != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("credentialCiphertext")+" "+s.CredentialCiphertext.String())
	}

	if s.StoreID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storeId")+" "+s.StoreID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := s.Store.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CustomerMemberSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("customer_members", ctx), sorts, joins)
}
func (s CustomerMemberSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.MemberNumber != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("memberNumber")+" "+s.MemberNumber.String())
	}

	if s.RequestKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requestKey")+" "+s.RequestKey.String())
	}

	if s.Phone != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("phone")+" "+s.Phone.String())
	}

	if s.PointsBalance != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("pointsBalance")+" "+s.PointsBalance.String())
	}

	if s.PointsFrozen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("pointsFrozen")+" "+s.PointsFrozen.String())
	}

	if s.NoticeVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("noticeVersion")+" "+s.NoticeVersion.String())
	}

	if s.ProcessingBasisCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("processingBasisCode")+" "+s.ProcessingBasisCode.String())
	}

	if s.EvidenceReference != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("evidenceReference")+" "+s.EvidenceReference.String())
	}

	if s.CancellationRequestedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("cancellationRequestedAt")+" "+s.CancellationRequestedAt.String())
	}

	if s.CancelledAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("cancelledAt")+" "+s.CancelledAt.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.PointEntries != nil {
		_alias := alias + "_pointEntries"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+"."+"member_id"+" = "+alias+".id")
		err := s.PointEntries.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CouponGrants != nil {
		_alias := alias + "_couponGrants"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_grants", ctx)+" "+_alias+" ON "+_alias+"."+"member_id"+" = "+alias+".id")
		err := s.CouponGrants.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CouponDistributionJobs != nil {
		_alias := alias + "_couponDistributionJobs"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_distribution_jobs", ctx)+" "+_alias+" ON "+_alias+"."+"member_id"+" = "+alias+".id")
		err := s.CouponDistributionJobs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CustomerBenefitPolicySortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("customer_benefit_policies", ctx), sorts, joins)
}
func (s CustomerBenefitPolicySortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Version != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("version")+" "+s.Version.String())
	}

	if s.DiscountEnabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("discountEnabled")+" "+s.DiscountEnabled.String())
	}

	if s.DiscountBasisPoints != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("discountBasisPoints")+" "+s.DiscountBasisPoints.String())
	}

	if s.PurchaseEarnEnabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("purchaseEarnEnabled")+" "+s.PurchaseEarnEnabled.String())
	}

	if s.EarnAmountFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("earnAmountFen")+" "+s.EarnAmountFen.String())
	}

	if s.EarnPoints != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("earnPoints")+" "+s.EarnPoints.String())
	}

	if s.RedemptionEnabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("redemptionEnabled")+" "+s.RedemptionEnabled.String())
	}

	if s.RedeemPoints != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("redeemPoints")+" "+s.RedeemPoints.String())
	}

	if s.RedeemAmountFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("redeemAmountFen")+" "+s.RedeemAmountFen.String())
	}

	if s.MaxRedemptionBasisPoints != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("maxRedemptionBasisPoints")+" "+s.MaxRedemptionBasisPoints.String())
	}

	if s.MaxRedemptionPoints != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("maxRedemptionPoints")+" "+s.MaxRedemptionPoints.String())
	}

	if s.ManualGrantMaxSingle != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("manualGrantMaxSingle")+" "+s.ManualGrantMaxSingle.String())
	}

	if s.ManualGrantMaxDaily != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("manualGrantMaxDaily")+" "+s.ManualGrantMaxDaily.String())
	}

	if s.PromotionWithHqCoupon != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("promotionWithHqCoupon")+" "+s.PromotionWithHqCoupon.String())
	}

	if s.PromotionWithStoreCoupon != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("promotionWithStoreCoupon")+" "+s.PromotionWithStoreCoupon.String())
	}

	if s.MemberPriceWithPromotion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("memberPriceWithPromotion")+" "+s.MemberPriceWithPromotion.String())
	}

	if s.PointsWithPromotion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("pointsWithPromotion")+" "+s.PointsWithPromotion.String())
	}

	if s.PointsWithCoupon != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("pointsWithCoupon")+" "+s.PointsWithCoupon.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CustomerDailyPointGrantBudgetSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("customer_daily_point_grant_budgets", ctx), sorts, joins)
}
func (s CustomerDailyPointGrantBudgetSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.BusinessDate != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("businessDate")+" "+s.BusinessDate.String())
	}

	if s.UsedPoints != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("usedPoints")+" "+s.UsedPoints.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CustomerPointEntrySortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("customer_point_entries", ctx), sorts, joins)
}
func (s CustomerPointEntrySortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Delta != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("delta")+" "+s.Delta.String())
	}

	if s.Note != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("note")+" "+s.Note.String())
	}

	if s.RequestKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requestKey")+" "+s.RequestKey.String())
	}

	if s.MemberID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("memberId")+" "+s.MemberID.String())
	}

	if s.SourceOrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("sourceOrganizationId")+" "+s.SourceOrganizationID.String())
	}

	if s.ReversesID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reversesId")+" "+s.ReversesID.String())
	}

	if s.ReversedByID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reversedById")+" "+s.ReversedByID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Member != nil {
		_alias := alias + "_member"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"member_id")
		err := s.Member.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.SourceOrganization != nil {
		_alias := alias + "_sourceOrganization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"source_organization_id")
		err := s.SourceOrganization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Reverses != nil {
		_alias := alias + "_reverses"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"reverses_id")
		err := s.Reverses.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ReversedBy != nil {
		_alias := alias + "_reversedBy"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"reversed_by_id")
		err := s.ReversedBy.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CustomerCouponTemplateSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("customer_coupon_templates", ctx), sorts, joins)
}
func (s CustomerCouponTemplateSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Code != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("code")+" "+s.Code.String())
	}

	if s.RequestKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requestKey")+" "+s.RequestKey.String())
	}

	if s.Title != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("title")+" "+s.Title.String())
	}

	if s.AmountFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("amountFen")+" "+s.AmountFen.String())
	}

	if s.MinSpendFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("minSpendFen")+" "+s.MinSpendFen.String())
	}

	if s.DaysAfterActivation != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("daysAfterActivation")+" "+s.DaysAfterActivation.String())
	}

	if s.EffectiveAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("effectiveAt")+" "+s.EffectiveAt.String())
	}

	if s.DistributionEndsAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("distributionEndsAt")+" "+s.DistributionEndsAt.String())
	}

	if s.PerMemberLimit != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("perMemberLimit")+" "+s.PerMemberLimit.String())
	}

	if s.TotalIssueLimit != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("totalIssueLimit")+" "+s.TotalIssueLimit.String())
	}

	if s.IssuedCount != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("issuedCount")+" "+s.IssuedCount.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.ApplicableStoreID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("applicableStoreId")+" "+s.ApplicableStoreID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ApplicableStore != nil {
		_alias := alias + "_applicableStore"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"applicable_store_id")
		err := s.ApplicableStore.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Grants != nil {
		_alias := alias + "_grants"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_grants", ctx)+" "+_alias+" ON "+_alias+"."+"template_id"+" = "+alias+".id")
		err := s.Grants.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.DistributionJobs != nil {
		_alias := alias + "_distributionJobs"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_distribution_jobs", ctx)+" "+_alias+" ON "+_alias+"."+"template_id"+" = "+alias+".id")
		err := s.DistributionJobs.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductCategorySortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("product_categories", ctx), sorts, joins)
}
func (s ProductCategorySortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.SortOrder != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("sortOrder")+" "+s.SortOrder.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.ParentID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("parentId")+" "+s.ParentID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Parent != nil {
		_alias := alias + "_parent"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"parent_id")
		err := s.Parent.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Children != nil {
		_alias := alias + "_children"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+"."+"parent_id"+" = "+alias+".id")
		err := s.Children.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Products != nil {
		_alias := alias + "_products"
		*joins = append(*joins, "LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"category_id"+" = "+alias+".id")
		err := s.Products.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductBrandSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("product_brands", ctx), sorts, joins)
}
func (s ProductBrandSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Products != nil {
		_alias := alias + "_products"
		*joins = append(*joins, "LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"brand_id"+" = "+alias+".id")
		err := s.Products.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("products", ctx), sorts, joins)
}
func (s ProductSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.Description != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("description")+" "+s.Description.String())
	}

	if s.ImageURL != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("imageUrl")+" "+s.ImageURL.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.BrandID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("brandId")+" "+s.BrandID.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.CategoryID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("categoryId")+" "+s.CategoryID.String())
	}

	if s.DefaultPackageTemplateID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("defaultPackageTemplateId")+" "+s.DefaultPackageTemplateID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Brand != nil {
		_alias := alias + "_brand"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_brands", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"brand_id")
		err := s.Brand.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Category != nil {
		_alias := alias + "_category"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"category_id")
		err := s.Category.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.DefaultPackageTemplate != nil {
		_alias := alias + "_defaultPackageTemplate"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"default_package_template_id")
		err := s.DefaultPackageTemplate.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.SpecificationChoices != nil {
		_alias := alias + "_specificationChoices"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_specification_choices", ctx)+" "+_alias+" ON "+_alias+"."+"product_id"+" = "+alias+".id")
		err := s.SpecificationChoices.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Skus != nil {
		_alias := alias + "_skus"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+"."+"product_id"+" = "+alias+".id")
		err := s.Skus.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductSkuSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("product_skus", ctx), sorts, joins)
}
func (s ProductSkuSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.PublishedPackageSetVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("publishedPackageSetVersion")+" "+s.PublishedPackageSetVersion.String())
	}

	if s.Ingredients != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("ingredients")+" "+s.Ingredients.String())
	}

	if s.Allergens != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("allergens")+" "+s.Allergens.String())
	}

	if s.StorageInstructions != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storageInstructions")+" "+s.StorageInstructions.String())
	}

	if s.ShelfLifeDays != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("shelfLifeDays")+" "+s.ShelfLifeDays.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.SelectionRetired != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("selectionRetired")+" "+s.SelectionRetired.String())
	}

	if s.ProductID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("productId")+" "+s.ProductID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Product != nil {
		_alias := alias + "_product"
		*joins = append(*joins, "LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"product_id")
		err := s.Product.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.SpecificationValues != nil {
		_alias := alias + "_specificationValues"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_sku_specification_values", ctx)+" "+_alias+" ON "+_alias+"."+"sku_id"+" = "+alias+".id")
		err := s.SpecificationValues.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Packages != nil {
		_alias := alias + "_packages"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+"."+"sku_id"+" = "+alias+".id")
		err := s.Packages.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Listings != nil {
		_alias := alias + "_listings"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+"."+"sku_id"+" = "+alias+".id")
		err := s.Listings.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductPackageSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("product_packages", ctx), sorts, joins)
}
func (s ProductPackageSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.Barcode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("barcode")+" "+s.Barcode.String())
	}

	if s.PackageSetVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageSetVersion")+" "+s.PackageSetVersion.String())
	}

	if s.ContainsQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("containsQuantity")+" "+s.ContainsQuantity.String())
	}

	if s.SuggestedPriceFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("suggestedPriceFen")+" "+s.SuggestedPriceFen.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.SkuID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("skuId")+" "+s.SkuID.String())
	}

	if s.TemplateID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("templateId")+" "+s.TemplateID.String())
	}

	if s.ContainsPackageID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("containsPackageId")+" "+s.ContainsPackageID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Sku != nil {
		_alias := alias + "_sku"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"sku_id")
		err := s.Sku.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Template != nil {
		_alias := alias + "_template"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"template_id")
		err := s.Template.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ContainsPackage != nil {
		_alias := alias + "_containsPackage"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"contains_package_id")
		err := s.ContainsPackage.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ContainedByPackages != nil {
		_alias := alias + "_containedByPackages"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+"."+"contains_package_id"+" = "+alias+".id")
		err := s.ContainedByPackages.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Offers != nil {
		_alias := alias + "_offers"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+"."+"package_id"+" = "+alias+".id")
		err := s.Offers.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Balances != nil {
		_alias := alias + "_balances"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stock_balances", ctx)+" "+_alias+" ON "+_alias+"."+"package_id"+" = "+alias+".id")
		err := s.Balances.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.MovementSources != nil {
		_alias := alias + "_movementSources"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"source_package_id"+" = "+alias+".id")
		err := s.MovementSources.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.MovementTargets != nil {
		_alias := alias + "_movementTargets"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"target_package_id"+" = "+alias+".id")
		err := s.MovementTargets.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.StocktakeLines != nil {
		_alias := alias + "_stocktakeLines"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+"."+"package_id"+" = "+alias+".id")
		err := s.StocktakeLines.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SpecificationDefinitionSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("specification_definitions", ctx), sorts, joins)
}
func (s SpecificationDefinitionSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Values != nil {
		_alias := alias + "_values"
		*joins = append(*joins, "LEFT JOIN "+TableName("specification_values", ctx)+" "+_alias+" ON "+_alias+"."+"specification_id"+" = "+alias+".id")
		err := s.Values.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SpecificationValueSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("specification_values", ctx), sorts, joins)
}
func (s SpecificationValueSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.SpecificationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("specificationId")+" "+s.SpecificationID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Specification != nil {
		_alias := alias + "_specification"
		*joins = append(*joins, "LEFT JOIN "+TableName("specification_definitions", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"specification_id")
		err := s.Specification.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ProductChoices != nil {
		_alias := alias + "_productChoices"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_specification_choices", ctx)+" "+_alias+" ON "+_alias+"."+"value_id"+" = "+alias+".id")
		err := s.ProductChoices.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.SkuValues != nil {
		_alias := alias + "_skuValues"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_sku_specification_values", ctx)+" "+_alias+" ON "+_alias+"."+"value_id"+" = "+alias+".id")
		err := s.SkuValues.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductSpecificationChoiceSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("product_specification_choices", ctx), sorts, joins)
}
func (s ProductSpecificationChoiceSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.ProductID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("productId")+" "+s.ProductID.String())
	}

	if s.ValueID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("valueId")+" "+s.ValueID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Product != nil {
		_alias := alias + "_product"
		*joins = append(*joins, "LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"product_id")
		err := s.Product.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Value != nil {
		_alias := alias + "_value"
		*joins = append(*joins, "LEFT JOIN "+TableName("specification_values", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"value_id")
		err := s.Value.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductSkuSpecificationValueSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("product_sku_specification_values", ctx), sorts, joins)
}
func (s ProductSkuSpecificationValueSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.SkuID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("skuId")+" "+s.SkuID.String())
	}

	if s.ValueID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("valueId")+" "+s.ValueID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Sku != nil {
		_alias := alias + "_sku"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"sku_id")
		err := s.Sku.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Value != nil {
		_alias := alias + "_value"
		*joins = append(*joins, "LEFT JOIN "+TableName("specification_values", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"value_id")
		err := s.Value.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ProductPackageTemplateSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("product_package_templates", ctx), sorts, joins)
}
func (s ProductPackageTemplateSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Name != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("name")+" "+s.Name.String())
	}

	if s.ContainsQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("containsQuantity")+" "+s.ContainsQuantity.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.OrganizationID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("organizationId")+" "+s.OrganizationID.String())
	}

	if s.ContainsPackageID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("containsPackageId")+" "+s.ContainsPackageID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := s.Organization.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.DefaultProducts != nil {
		_alias := alias + "_defaultProducts"
		*joins = append(*joins, "LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"default_package_template_id"+" = "+alias+".id")
		err := s.DefaultProducts.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ContainsPackage != nil {
		_alias := alias + "_containsPackage"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"contains_package_id")
		err := s.ContainsPackage.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.ContainedByPackages != nil {
		_alias := alias + "_containedByPackages"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+"."+"contains_package_id"+" = "+alias+".id")
		err := s.ContainedByPackages.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.CreatedPackages != nil {
		_alias := alias + "_createdPackages"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+"."+"template_id"+" = "+alias+".id")
		err := s.CreatedPackages.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StoreListingSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_listings", ctx), sorts, joins)
}
func (s StoreListingSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.SelectedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("selectedAt")+" "+s.SelectedAt.String())
	}

	if s.StoreID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storeId")+" "+s.StoreID.String())
	}

	if s.SkuID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("skuId")+" "+s.SkuID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := s.Store.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Sku != nil {
		_alias := alias + "_sku"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"sku_id")
		err := s.Sku.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Offers != nil {
		_alias := alias + "_offers"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+"."+"listing_id"+" = "+alias+".id")
		err := s.Offers.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Batches != nil {
		_alias := alias + "_batches"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+"."+"listing_id"+" = "+alias+".id")
		err := s.Batches.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StorePackageOfferSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_package_offers", ctx), sorts, joins)
}
func (s StorePackageOfferSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.PriceFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("priceFen")+" "+s.PriceFen.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.ListingID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("listingId")+" "+s.ListingID.String())
	}

	if s.PackageID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageId")+" "+s.PackageID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Listing != nil {
		_alias := alias + "_listing"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"listing_id")
		err := s.Listing.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Package != nil {
		_alias := alias + "_package"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"package_id")
		err := s.Package.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.PriceRevisions != nil {
		_alias := alias + "_priceRevisions"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_price_revisions", ctx)+" "+_alias+" ON "+_alias+"."+"offer_id"+" = "+alias+".id")
		err := s.PriceRevisions.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.PromotionTargets != nil {
		_alias := alias + "_promotionTargets"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_promotion_targets", ctx)+" "+_alias+" ON "+_alias+"."+"offer_id"+" = "+alias+".id")
		err := s.PromotionTargets.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StorePriceRevisionSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_price_revisions", ctx), sorts, joins)
}
func (s StorePriceRevisionSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.PreviousPriceFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("previousPriceFen")+" "+s.PreviousPriceFen.String())
	}

	if s.PriceFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("priceFen")+" "+s.PriceFen.String())
	}

	if s.EffectiveAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("effectiveAt")+" "+s.EffectiveAt.String())
	}

	if s.ReasonCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reasonCode")+" "+s.ReasonCode.String())
	}

	if s.OfferID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("offerId")+" "+s.OfferID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Offer != nil {
		_alias := alias + "_offer"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"offer_id")
		err := s.Offer.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StoreInventoryBatchSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_inventory_batches", ctx), sorts, joins)
}
func (s StoreInventoryBatchSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.BatchNumber != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("batchNumber")+" "+s.BatchNumber.String())
	}

	if s.ProducedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("producedAt")+" "+s.ProducedAt.String())
	}

	if s.ExpiresAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("expiresAt")+" "+s.ExpiresAt.String())
	}

	if s.SourceReference != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("sourceReference")+" "+s.SourceReference.String())
	}

	if s.ListingID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("listingId")+" "+s.ListingID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Listing != nil {
		_alias := alias + "_listing"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"listing_id")
		err := s.Listing.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Balances != nil {
		_alias := alias + "_balances"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stock_balances", ctx)+" "+_alias+" ON "+_alias+"."+"batch_id"+" = "+alias+".id")
		err := s.Balances.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Movements != nil {
		_alias := alias + "_movements"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"batch_id"+" = "+alias+".id")
		err := s.Movements.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.StocktakeLines != nil {
		_alias := alias + "_stocktakeLines"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+"."+"batch_id"+" = "+alias+".id")
		err := s.StocktakeLines.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StoreStockBalanceSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_stock_balances", ctx), sorts, joins)
}
func (s StoreStockBalanceSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.Quantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("quantity")+" "+s.Quantity.String())
	}

	if s.Version != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("version")+" "+s.Version.String())
	}

	if s.BatchID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("batchId")+" "+s.BatchID.String())
	}

	if s.PackageID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageId")+" "+s.PackageID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Batch != nil {
		_alias := alias + "_batch"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"batch_id")
		err := s.Batch.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Package != nil {
		_alias := alias + "_package"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"package_id")
		err := s.Package.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StoreStocktakeSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_stocktakes", ctx), sorts, joins)
}
func (s StoreStocktakeSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.RequestKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requestKey")+" "+s.RequestKey.String())
	}

	if s.ScopeDigest != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("scopeDigest")+" "+s.ScopeDigest.String())
	}

	if s.StartedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("startedAt")+" "+s.StartedAt.String())
	}

	if s.ReviewedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reviewedAt")+" "+s.ReviewedAt.String())
	}

	if s.PostedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("postedAt")+" "+s.PostedAt.String())
	}

	if s.CanceledAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("canceledAt")+" "+s.CanceledAt.String())
	}

	if s.StoreID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storeId")+" "+s.StoreID.String())
	}

	if s.InitiatedByAccountID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("initiatedByAccountId")+" "+s.InitiatedByAccountID.String())
	}

	if s.PostedByID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("postedById")+" "+s.PostedByID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := s.Store.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.InitiatedByAccount != nil {
		_alias := alias + "_initiatedByAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initiated_by_account_id")
		err := s.InitiatedByAccount.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.PostedBy != nil {
		_alias := alias + "_postedBy"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"posted_by_id")
		err := s.PostedBy.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Lines != nil {
		_alias := alias + "_lines"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+"."+"stocktake_id"+" = "+alias+".id")
		err := s.Lines.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StoreStocktakeLineSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_stocktake_lines", ctx), sorts, joins)
}
func (s StoreStocktakeLineSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.SnapshotQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("snapshotQuantity")+" "+s.SnapshotQuantity.String())
	}

	if s.SnapshotVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("snapshotVersion")+" "+s.SnapshotVersion.String())
	}

	if s.PackageSetVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageSetVersion")+" "+s.PackageSetVersion.String())
	}

	if s.BatchNumberSnapshot != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("batchNumberSnapshot")+" "+s.BatchNumberSnapshot.String())
	}

	if s.ExpiresAtSnapshot != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("expiresAtSnapshot")+" "+s.ExpiresAtSnapshot.String())
	}

	if s.PackageNameSnapshot != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageNameSnapshot")+" "+s.PackageNameSnapshot.String())
	}

	if s.PackageEnabledSnapshot != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageEnabledSnapshot")+" "+s.PackageEnabledSnapshot.String())
	}

	if s.CountedQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("countedQuantity")+" "+s.CountedQuantity.String())
	}

	if s.ReasonCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reasonCode")+" "+s.ReasonCode.String())
	}

	if s.ReasonNote != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reasonNote")+" "+s.ReasonNote.String())
	}

	if s.CountedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("countedAt")+" "+s.CountedAt.String())
	}

	if s.NeedsRecount != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("needsRecount")+" "+s.NeedsRecount.String())
	}

	if s.StocktakeID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("stocktakeId")+" "+s.StocktakeID.String())
	}

	if s.BatchID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("batchId")+" "+s.BatchID.String())
	}

	if s.PackageID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageId")+" "+s.PackageID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Stocktake != nil {
		_alias := alias + "_stocktake"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"stocktake_id")
		err := s.Stocktake.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Batch != nil {
		_alias := alias + "_batch"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"batch_id")
		err := s.Batch.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Package != nil {
		_alias := alias + "_package"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"package_id")
		err := s.Package.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Movements != nil {
		_alias := alias + "_movements"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"stocktake_line_id"+" = "+alias+".id")
		err := s.Movements.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StoreStockMovementSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_stock_movements", ctx), sorts, joins)
}
func (s StoreStockMovementSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.RequestKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requestKey")+" "+s.RequestKey.String())
	}

	if s.SourceQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("sourceQuantity")+" "+s.SourceQuantity.String())
	}

	if s.TargetQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("targetQuantity")+" "+s.TargetQuantity.String())
	}

	if s.FactorSnapshot != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("factorSnapshot")+" "+s.FactorSnapshot.String())
	}

	if s.PackageSetVersion != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("packageSetVersion")+" "+s.PackageSetVersion.String())
	}

	if s.ReasonCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("reasonCode")+" "+s.ReasonCode.String())
	}

	if s.OccurredAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("occurredAt")+" "+s.OccurredAt.String())
	}

	if s.StoreID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storeId")+" "+s.StoreID.String())
	}

	if s.BatchID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("batchId")+" "+s.BatchID.String())
	}

	if s.SourcePackageID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("sourcePackageId")+" "+s.SourcePackageID.String())
	}

	if s.TargetPackageID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("targetPackageId")+" "+s.TargetPackageID.String())
	}

	if s.StocktakeLineID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("stocktakeLineId")+" "+s.StocktakeLineID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := s.Store.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Batch != nil {
		_alias := alias + "_batch"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"batch_id")
		err := s.Batch.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.SourcePackage != nil {
		_alias := alias + "_sourcePackage"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"source_package_id")
		err := s.SourcePackage.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.TargetPackage != nil {
		_alias := alias + "_targetPackage"
		*joins = append(*joins, "LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"target_package_id")
		err := s.TargetPackage.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.StocktakeLine != nil {
		_alias := alias + "_stocktakeLine"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"stocktake_line_id")
		err := s.StocktakeLine.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StorePromotionSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_promotions", ctx), sorts, joins)
}
func (s StorePromotionSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.RuleKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("ruleKey")+" "+s.RuleKey.String())
	}

	if s.TimeZone != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("timeZone")+" "+s.TimeZone.String())
	}

	if s.Version != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("version")+" "+s.Version.String())
	}

	if s.Enabled != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("enabled")+" "+s.Enabled.String())
	}

	if s.StartsAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("startsAt")+" "+s.StartsAt.String())
	}

	if s.EndsAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("endsAt")+" "+s.EndsAt.String())
	}

	if s.ThresholdFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("thresholdFen")+" "+s.ThresholdFen.String())
	}

	if s.ThresholdQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("thresholdQuantity")+" "+s.ThresholdQuantity.String())
	}

	if s.DiscountFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("discountFen")+" "+s.DiscountFen.String())
	}

	if s.DiscountBasisPoints != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("discountBasisPoints")+" "+s.DiscountBasisPoints.String())
	}

	if s.FixedPriceFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("fixedPriceFen")+" "+s.FixedPriceFen.String())
	}

	if s.StackWithHqCoupon != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("stackWithHqCoupon")+" "+s.StackWithHqCoupon.String())
	}

	if s.StackWithStoreCoupon != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("stackWithStoreCoupon")+" "+s.StackWithStoreCoupon.String())
	}

	if s.StackWithMemberPrice != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("stackWithMemberPrice")+" "+s.StackWithMemberPrice.String())
	}

	if s.StoreID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("storeId")+" "+s.StoreID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := s.Store.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Targets != nil {
		_alias := alias + "_targets"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_promotion_targets", ctx)+" "+_alias+" ON "+_alias+"."+"promotion_id"+" = "+alias+".id")
		err := s.Targets.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s StorePromotionTargetSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("store_promotion_targets", ctx), sorts, joins)
}
func (s StorePromotionTargetSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.RequiredQuantity != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requiredQuantity")+" "+s.RequiredQuantity.String())
	}

	if s.PromotionID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("promotionId")+" "+s.PromotionID.String())
	}

	if s.OfferID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("offerId")+" "+s.OfferID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Promotion != nil {
		_alias := alias + "_promotion"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_promotions", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"promotion_id")
		err := s.Promotion.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Offer != nil {
		_alias := alias + "_offer"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"offer_id")
		err := s.Offer.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CustomerCouponGrantSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("customer_coupon_grants", ctx), sorts, joins)
}
func (s CustomerCouponGrantSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.AmountFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("amountFen")+" "+s.AmountFen.String())
	}

	if s.MinSpendFen != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("minSpendFen")+" "+s.MinSpendFen.String())
	}

	if s.DaysAfterActivation != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("daysAfterActivation")+" "+s.DaysAfterActivation.String())
	}

	if s.IssuedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("issuedAt")+" "+s.IssuedAt.String())
	}

	if s.ActivatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("activatedAt")+" "+s.ActivatedAt.String())
	}

	if s.ExpiresAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("expiresAt")+" "+s.ExpiresAt.String())
	}

	if s.RevokedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("revokedAt")+" "+s.RevokedAt.String())
	}

	if s.RequestKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requestKey")+" "+s.RequestKey.String())
	}

	if s.IssuerRequestDigest != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("issuerRequestDigest")+" "+s.IssuerRequestDigest.String())
	}

	if s.MemberID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("memberId")+" "+s.MemberID.String())
	}

	if s.TemplateID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("templateId")+" "+s.TemplateID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Member != nil {
		_alias := alias + "_member"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"member_id")
		err := s.Member.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Template != nil {
		_alias := alias + "_template"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"template_id")
		err := s.Template.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CustomerCouponDistributionJobSortType) Apply(ctx context.Context, sorts *[]string, joins *[]string) error {
	return s.ApplyWithAlias(ctx, TableName("customer_coupon_distribution_jobs", ctx), sorts, joins)
}
func (s CustomerCouponDistributionJobSortType) ApplyWithAlias(ctx context.Context, alias string, sorts *[]string, joins *[]string) error {
	aliasPrefix := alias + "."

	if s.ID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("id")+" "+s.ID.String())
	}

	if s.RequestKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("requestKey")+" "+s.RequestKey.String())
	}

	if s.AvailableAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("availableAt")+" "+s.AvailableAt.String())
	}

	if s.LeaseExpiresAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("leaseExpiresAt")+" "+s.LeaseExpiresAt.String())
	}

	if s.LeaseToken != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("leaseToken")+" "+s.LeaseToken.String())
	}

	if s.CursorCreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("cursorCreatedAt")+" "+s.CursorCreatedAt.String())
	}

	if s.CursorKey != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("cursorKey")+" "+s.CursorKey.String())
	}

	if s.Attempts != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("attempts")+" "+s.Attempts.String())
	}

	if s.LastErrorCode != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("lastErrorCode")+" "+s.LastErrorCode.String())
	}

	if s.TemplateID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("templateId")+" "+s.TemplateID.String())
	}

	if s.MemberID != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("memberId")+" "+s.MemberID.String())
	}

	if s.IsDelete != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("isDelete")+" "+s.IsDelete.String())
	}

	if s.Weight != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("weight")+" "+s.Weight.String())
	}

	if s.State != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("state")+" "+s.State.String())
	}

	if s.DeletedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedBy")+" "+s.DeletedBy.String())
	}

	if s.UpdatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedBy")+" "+s.UpdatedBy.String())
	}

	if s.CreatedBy != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdBy")+" "+s.CreatedBy.String())
	}

	if s.DeletedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("deletedAt")+" "+s.DeletedAt.String())
	}

	if s.UpdatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("updatedAt")+" "+s.UpdatedAt.String())
	}

	if s.CreatedAt != nil {
		*sorts = append(*sorts, aliasPrefix+SnakeString("createdAt")+" "+s.CreatedAt.String())
	}

	if s.Template != nil {
		_alias := alias + "_template"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"template_id")
		err := s.Template.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	if s.Member != nil {
		_alias := alias + "_member"
		*joins = append(*joins, "LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"member_id")
		err := s.Member.ApplyWithAlias(ctx, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}

	return nil
}
