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
