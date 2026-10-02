package gen

import (
	"context"
	"fmt"
	"strings"
)

// CompactNils removes nil pointers from a slice.
func CompactNils[T any](items []*T) []*T {
	result := make([]*T, 0, len(items))
	for _, item := range items {
		if item != nil {
			result = append(result, item)
		}
	}
	return result
}

func (f *AccountFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *AccountFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("accounts", ctx), wheres, values, joins)
}
func (f *AccountFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Memberships != nil {
		_alias := alias + "_memberships"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"."+"account_id"+" = "+alias+".id")
		err := f.Memberships.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.InitializedOrganizations != nil {
		_alias := alias + "_initializedOrganizations"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+"."+"initial_account_id"+" = "+alias+".id")
		err := f.InitializedOrganizations.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.OpeningRecords != nil {
		_alias := alias + "_openingRecords"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"initial_account_id"+" = "+alias+".id")
		err := f.OpeningRecords.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.RecordedOpeningRecords != nil {
		_alias := alias + "_recordedOpeningRecords"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"recorded_by_account_id"+" = "+alias+".id")
		err := f.RecordedOpeningRecords.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Sessions != nil {
		_alias := alias + "_sessions"
		*joins = append(*joins, "LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+"."+"account_id"+" = "+alias+".id")
		err := f.Sessions.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.ReviewedStores != nil {
		_alias := alias + "_reviewedStores"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"."+"reviewed_by_account_id"+" = "+alias+".id")
		err := f.ReviewedStores.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.SentMembershipInvitations != nil {
		_alias := alias + "_sentMembershipInvitations"
		*joins = append(*joins, "LEFT JOIN "+TableName("membership_invitations", ctx)+" "+_alias+" ON "+_alias+"."+"invited_by_account_id"+" = "+alias+".id")
		err := f.SentMembershipInvitations.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"actor_account_id"+" = "+alias+".id")
		err := f.AuditLogs.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *AccountFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Phone != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, f.Phone)
	}

	if f.PhoneNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, f.PhoneNe)
	}

	if f.PhoneGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, f.PhoneGt)
	}

	if f.PhoneLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, f.PhoneLt)
	}

	if f.PhoneGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, f.PhoneGte)
	}

	if f.PhoneLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, f.PhoneLte)
	}

	if f.PhoneIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, f.PhoneIn)
	}

	if f.PhoneLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.PhoneLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.PhonePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, fmt.Sprintf("%s%%", *f.PhonePrefix))
	}

	if f.PhoneSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("phone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		values = append(values, fmt.Sprintf("%%%s", *f.PhoneSuffix))
	}

	if f.PhoneNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("phone")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("phone"))
		if *f.PhoneNull {
			conditions = append(conditions, aliasPrefix+SnakeString("phone")+" IS NULL"+" OR "+aliasPrefix+SnakeString("phone")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("phone")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("phone")+" <> ''")
		}
	}

	if f.DisplayName != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, f.DisplayName)
	}

	if f.DisplayNameNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, f.DisplayNameNe)
	}

	if f.DisplayNameGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, f.DisplayNameGt)
	}

	if f.DisplayNameLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, f.DisplayNameLt)
	}

	if f.DisplayNameGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, f.DisplayNameGte)
	}

	if f.DisplayNameLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, f.DisplayNameLte)
	}

	if f.DisplayNameIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, f.DisplayNameIn)
	}

	if f.DisplayNameLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.DisplayNameLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.DisplayNamePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, fmt.Sprintf("%s%%", *f.DisplayNamePrefix))
	}

	if f.DisplayNameSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		values = append(values, fmt.Sprintf("%%%s", *f.DisplayNameSuffix))
	}

	if f.DisplayNameNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("displayName")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("displayName"))
		if *f.DisplayNameNull {
			conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" IS NULL"+" OR "+aliasPrefix+SnakeString("displayName")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("displayName")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("displayName")+" <> ''")
		}
	}

	if f.Email != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, f.Email)
	}

	if f.EmailNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, f.EmailNe)
	}

	if f.EmailGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, f.EmailGt)
	}

	if f.EmailLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, f.EmailLt)
	}

	if f.EmailGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, f.EmailGte)
	}

	if f.EmailLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, f.EmailLte)
	}

	if f.EmailIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, f.EmailIn)
	}

	if f.EmailLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.EmailLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.EmailPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, fmt.Sprintf("%s%%", *f.EmailPrefix))
	}

	if f.EmailSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("email")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		values = append(values, fmt.Sprintf("%%%s", *f.EmailSuffix))
	}

	if f.EmailNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("email")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("email"))
		if *f.EmailNull {
			conditions = append(conditions, aliasPrefix+SnakeString("email")+" IS NULL"+" OR "+aliasPrefix+SnakeString("email")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("email")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("email")+" <> ''")
		}
	}

	if f.Status != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.Status)
	}

	if f.StatusNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusNe)
	}

	if f.StatusGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusGt)
	}

	if f.StatusLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusLt)
	}

	if f.StatusGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusGte)
	}

	if f.StatusLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusLte)
	}

	if f.StatusIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusIn)
	}

	if f.StatusNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		if *f.StatusNull {
			conditions = append(conditions, aliasPrefix+SnakeString("status")+" IS NULL"+" OR "+aliasPrefix+SnakeString("status")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("status")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("status")+" <> ''")
		}
	}

	if f.MustChangePassword != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		values = append(values, f.MustChangePassword)
	}

	if f.MustChangePasswordNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		values = append(values, f.MustChangePasswordNe)
	}

	if f.MustChangePasswordGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		values = append(values, f.MustChangePasswordGt)
	}

	if f.MustChangePasswordLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		values = append(values, f.MustChangePasswordLt)
	}

	if f.MustChangePasswordGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		values = append(values, f.MustChangePasswordGte)
	}

	if f.MustChangePasswordLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		values = append(values, f.MustChangePasswordLte)
	}

	if f.MustChangePasswordIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		values = append(values, f.MustChangePasswordIn)
	}

	if f.MustChangePasswordNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("mustChangePassword")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("mustChangePassword"))
		if *f.MustChangePasswordNull {
			conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" IS NULL"+" OR "+aliasPrefix+SnakeString("mustChangePassword")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("mustChangePassword")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("mustChangePassword")+" <> ''")
		}
	}

	if f.CredentialVersion != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersion)
	}

	if f.CredentialVersionNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionNe)
	}

	if f.CredentialVersionGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionGt)
	}

	if f.CredentialVersionLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionLt)
	}

	if f.CredentialVersionGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionGte)
	}

	if f.CredentialVersionLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionLte)
	}

	if f.CredentialVersionIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionIn)
	}

	if f.CredentialVersionNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		if *f.CredentialVersionNull {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" IS NULL"+" OR "+aliasPrefix+SnakeString("credentialVersion")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("credentialVersion")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *AccountFilterType) AndWith(f2 ...*AccountFilterType) *AccountFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &AccountFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *AccountFilterType) OrWith(f2 ...*AccountFilterType) *AccountFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &AccountFilterType{
		Or: append(_f2, f),
	}
}

func (f *OrganizationFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *OrganizationFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("organizations", ctx), wheres, values, joins)
}
func (f *OrganizationFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Memberships != nil {
		_alias := alias + "_memberships"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := f.Memberships.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.InitialAccount != nil {
		_alias := alias + "_initialAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initial_account_id")
		err := f.InitialAccount.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.OpeningRecords != nil {
		_alias := alias + "_openingRecords"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := f.OpeningRecords.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Stores != nil {
		_alias := alias + "_stores"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := f.Stores.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Roles != nil {
		_alias := alias + "_roles"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := f.Roles.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Sessions != nil {
		_alias := alias + "_sessions"
		*joins = append(*joins, "LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := f.Sessions.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := f.AuditLogs.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.PaymentConfigs != nil {
		_alias := alias + "_paymentConfigs"
		*joins = append(*joins, "LEFT JOIN "+TableName("franchise_payment_configs", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")
		err := f.PaymentConfigs.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *OrganizationFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Code != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.Code)
	}

	if f.CodeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeNe)
	}

	if f.CodeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeGt)
	}

	if f.CodeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeLt)
	}

	if f.CodeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeGte)
	}

	if f.CodeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeLte)
	}

	if f.CodeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeIn)
	}

	if f.CodeLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.CodeLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.CodePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, fmt.Sprintf("%s%%", *f.CodePrefix))
	}

	if f.CodeSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, fmt.Sprintf("%%%s", *f.CodeSuffix))
	}

	if f.CodeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		if *f.CodeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("code")+" IS NULL"+" OR "+aliasPrefix+SnakeString("code")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("code")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("code")+" <> ''")
		}
	}

	if f.Name != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.Name)
	}

	if f.NameNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameNe)
	}

	if f.NameGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGt)
	}

	if f.NameLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLt)
	}

	if f.NameGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGte)
	}

	if f.NameLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLte)
	}

	if f.NameIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameIn)
	}

	if f.NameLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.NameLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.NamePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%s%%", *f.NamePrefix))
	}

	if f.NameSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%%%s", *f.NameSuffix))
	}

	if f.NameNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		if *f.NameNull {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NULL"+" OR "+aliasPrefix+SnakeString("name")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("name")+" <> ''")
		}
	}

	if f.Type != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("type")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		values = append(values, f.Type)
	}

	if f.TypeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("type")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		values = append(values, f.TypeNe)
	}

	if f.TypeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("type")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		values = append(values, f.TypeGt)
	}

	if f.TypeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("type")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		values = append(values, f.TypeLt)
	}

	if f.TypeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("type")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		values = append(values, f.TypeGte)
	}

	if f.TypeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("type")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		values = append(values, f.TypeLte)
	}

	if f.TypeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("type")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		values = append(values, f.TypeIn)
	}

	if f.TypeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("type")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("type"))
		if *f.TypeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("type")+" IS NULL"+" OR "+aliasPrefix+SnakeString("type")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("type")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("type")+" <> ''")
		}
	}

	if f.Status != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.Status)
	}

	if f.StatusNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusNe)
	}

	if f.StatusGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusGt)
	}

	if f.StatusLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusLt)
	}

	if f.StatusGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusGte)
	}

	if f.StatusLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusLte)
	}

	if f.StatusIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusIn)
	}

	if f.StatusNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		if *f.StatusNull {
			conditions = append(conditions, aliasPrefix+SnakeString("status")+" IS NULL"+" OR "+aliasPrefix+SnakeString("status")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("status")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("status")+" <> ''")
		}
	}

	if f.SuspendedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		values = append(values, f.SuspendedAt)
	}

	if f.SuspendedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		values = append(values, f.SuspendedAtNe)
	}

	if f.SuspendedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		values = append(values, f.SuspendedAtGt)
	}

	if f.SuspendedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		values = append(values, f.SuspendedAtLt)
	}

	if f.SuspendedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		values = append(values, f.SuspendedAtGte)
	}

	if f.SuspendedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		values = append(values, f.SuspendedAtLte)
	}

	if f.SuspendedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		values = append(values, f.SuspendedAtIn)
	}

	if f.SuspendedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspendedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspendedAt"))
		if *f.SuspendedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("suspendedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("suspendedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("suspendedAt")+" <> ''")
		}
	}

	if f.SuspensionReasonCode != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, f.SuspensionReasonCode)
	}

	if f.SuspensionReasonCodeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, f.SuspensionReasonCodeNe)
	}

	if f.SuspensionReasonCodeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, f.SuspensionReasonCodeGt)
	}

	if f.SuspensionReasonCodeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, f.SuspensionReasonCodeLt)
	}

	if f.SuspensionReasonCodeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, f.SuspensionReasonCodeGte)
	}

	if f.SuspensionReasonCodeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, f.SuspensionReasonCodeLte)
	}

	if f.SuspensionReasonCodeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, f.SuspensionReasonCodeIn)
	}

	if f.SuspensionReasonCodeLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.SuspensionReasonCodeLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.SuspensionReasonCodePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, fmt.Sprintf("%s%%", *f.SuspensionReasonCodePrefix))
	}

	if f.SuspensionReasonCodeSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		values = append(values, fmt.Sprintf("%%%s", *f.SuspensionReasonCodeSuffix))
	}

	if f.SuspensionReasonCodeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("suspensionReasonCode"))
		if *f.SuspensionReasonCodeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" IS NULL"+" OR "+aliasPrefix+SnakeString("suspensionReasonCode")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("suspensionReasonCode")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("suspensionReasonCode")+" <> ''")
		}
	}

	if f.InitialAccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountID)
	}

	if f.InitialAccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDNe)
	}

	if f.InitialAccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDGt)
	}

	if f.InitialAccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDLt)
	}

	if f.InitialAccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDGte)
	}

	if f.InitialAccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDLte)
	}

	if f.InitialAccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDIn)
	}

	if f.InitialAccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		if *f.InitialAccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("initialAccountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("initialAccountId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *OrganizationFilterType) AndWith(f2 ...*OrganizationFilterType) *OrganizationFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &OrganizationFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *OrganizationFilterType) OrWith(f2 ...*OrganizationFilterType) *OrganizationFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &OrganizationFilterType{
		Or: append(_f2, f),
	}
}

func (f *OperatorMembershipFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *OperatorMembershipFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("operator_memberships", ctx), wheres, values, joins)
}
func (f *OperatorMembershipFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Account != nil {
		_alias := alias + "_account"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"account_id")
		err := f.Account.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := f.Organization.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Roles != nil {
		_alias := alias + "_roles"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"member_id"+" LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"role_id"+" = "+_alias+".id")
		err := f.Roles.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Stores != nil {
		_alias := alias + "_stores"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_stores", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"member_id"+" LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"store_id"+" = "+_alias+".id")
		err := f.Stores.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Invitations != nil {
		_alias := alias + "_invitations"
		*joins = append(*joins, "LEFT JOIN "+TableName("membership_invitations", ctx)+" "+_alias+" ON "+_alias+"."+"membership_id"+" = "+alias+".id")
		err := f.Invitations.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *OperatorMembershipFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Status != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.Status)
	}

	if f.StatusNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusNe)
	}

	if f.StatusGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusGt)
	}

	if f.StatusLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusLt)
	}

	if f.StatusGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusGte)
	}

	if f.StatusLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusLte)
	}

	if f.StatusIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("status")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		values = append(values, f.StatusIn)
	}

	if f.StatusNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("status")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("status"))
		if *f.StatusNull {
			conditions = append(conditions, aliasPrefix+SnakeString("status")+" IS NULL"+" OR "+aliasPrefix+SnakeString("status")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("status")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("status")+" <> ''")
		}
	}

	if f.StoreAccessMode != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		values = append(values, f.StoreAccessMode)
	}

	if f.StoreAccessModeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		values = append(values, f.StoreAccessModeNe)
	}

	if f.StoreAccessModeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		values = append(values, f.StoreAccessModeGt)
	}

	if f.StoreAccessModeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		values = append(values, f.StoreAccessModeLt)
	}

	if f.StoreAccessModeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		values = append(values, f.StoreAccessModeGte)
	}

	if f.StoreAccessModeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		values = append(values, f.StoreAccessModeLte)
	}

	if f.StoreAccessModeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		values = append(values, f.StoreAccessModeIn)
	}

	if f.StoreAccessModeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeAccessMode")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeAccessMode"))
		if *f.StoreAccessModeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" IS NULL"+" OR "+aliasPrefix+SnakeString("storeAccessMode")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("storeAccessMode")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("storeAccessMode")+" <> ''")
		}
	}

	if f.InvitedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		values = append(values, f.InvitedAt)
	}

	if f.InvitedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		values = append(values, f.InvitedAtNe)
	}

	if f.InvitedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		values = append(values, f.InvitedAtGt)
	}

	if f.InvitedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		values = append(values, f.InvitedAtLt)
	}

	if f.InvitedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		values = append(values, f.InvitedAtGte)
	}

	if f.InvitedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		values = append(values, f.InvitedAtLte)
	}

	if f.InvitedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		values = append(values, f.InvitedAtIn)
	}

	if f.InvitedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedAt"))
		if *f.InvitedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("invitedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("invitedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("invitedAt")+" <> ''")
		}
	}

	if f.AcceptedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAt)
	}

	if f.AcceptedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtNe)
	}

	if f.AcceptedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtGt)
	}

	if f.AcceptedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtLt)
	}

	if f.AcceptedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtGte)
	}

	if f.AcceptedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtLte)
	}

	if f.AcceptedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtIn)
	}

	if f.AcceptedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		if *f.AcceptedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("acceptedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("acceptedAt")+" <> ''")
		}
	}

	if f.AccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountID)
	}

	if f.AccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDNe)
	}

	if f.AccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDGt)
	}

	if f.AccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDLt)
	}

	if f.AccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDGte)
	}

	if f.AccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDLte)
	}

	if f.AccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDIn)
	}

	if f.AccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		if *f.AccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("accountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("accountId")+" <> ''")
		}
	}

	if f.OrganizationID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationID)
	}

	if f.OrganizationIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDNe)
	}

	if f.OrganizationIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGt)
	}

	if f.OrganizationIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLt)
	}

	if f.OrganizationIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGte)
	}

	if f.OrganizationIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLte)
	}

	if f.OrganizationIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDIn)
	}

	if f.OrganizationIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		if *f.OrganizationIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *OperatorMembershipFilterType) AndWith(f2 ...*OperatorMembershipFilterType) *OperatorMembershipFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &OperatorMembershipFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *OperatorMembershipFilterType) OrWith(f2 ...*OperatorMembershipFilterType) *OperatorMembershipFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &OperatorMembershipFilterType{
		Or: append(_f2, f),
	}
}

func (f *PermissionFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *PermissionFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("permissions", ctx), wheres, values, joins)
}
func (f *PermissionFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Roles != nil {
		_alias := alias + "_roles"
		*joins = append(*joins, "LEFT JOIN "+TableName("permission_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"permission_id"+" LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"role_id"+" = "+_alias+".id")
		err := f.Roles.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *PermissionFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Name != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.Name)
	}

	if f.NameNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameNe)
	}

	if f.NameGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGt)
	}

	if f.NameLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLt)
	}

	if f.NameGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGte)
	}

	if f.NameLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLte)
	}

	if f.NameIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameIn)
	}

	if f.NameLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.NameLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.NamePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%s%%", *f.NamePrefix))
	}

	if f.NameSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%%%s", *f.NameSuffix))
	}

	if f.NameNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		if *f.NameNull {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NULL"+" OR "+aliasPrefix+SnakeString("name")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("name")+" <> ''")
		}
	}

	if f.Action != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.Action)
	}

	if f.ActionNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionNe)
	}

	if f.ActionGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionGt)
	}

	if f.ActionLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionLt)
	}

	if f.ActionGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionGte)
	}

	if f.ActionLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionLte)
	}

	if f.ActionIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionIn)
	}

	if f.ActionLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ActionLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ActionPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, fmt.Sprintf("%s%%", *f.ActionPrefix))
	}

	if f.ActionSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, fmt.Sprintf("%%%s", *f.ActionSuffix))
	}

	if f.ActionNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		if *f.ActionNull {
			conditions = append(conditions, aliasPrefix+SnakeString("action")+" IS NULL"+" OR "+aliasPrefix+SnakeString("action")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("action")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("action")+" <> ''")
		}
	}

	if f.Module != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, f.Module)
	}

	if f.ModuleNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, f.ModuleNe)
	}

	if f.ModuleGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, f.ModuleGt)
	}

	if f.ModuleLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, f.ModuleLt)
	}

	if f.ModuleGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, f.ModuleGte)
	}

	if f.ModuleLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, f.ModuleLte)
	}

	if f.ModuleIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, f.ModuleIn)
	}

	if f.ModuleLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ModuleLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ModulePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, fmt.Sprintf("%s%%", *f.ModulePrefix))
	}

	if f.ModuleSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("module")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		values = append(values, fmt.Sprintf("%%%s", *f.ModuleSuffix))
	}

	if f.ModuleNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("module")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("module"))
		if *f.ModuleNull {
			conditions = append(conditions, aliasPrefix+SnakeString("module")+" IS NULL"+" OR "+aliasPrefix+SnakeString("module")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("module")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("module")+" <> ''")
		}
	}

	if f.Scope != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("scope")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		values = append(values, f.Scope)
	}

	if f.ScopeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("scope")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		values = append(values, f.ScopeNe)
	}

	if f.ScopeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("scope")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		values = append(values, f.ScopeGt)
	}

	if f.ScopeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("scope")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		values = append(values, f.ScopeLt)
	}

	if f.ScopeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("scope")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		values = append(values, f.ScopeGte)
	}

	if f.ScopeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("scope")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		values = append(values, f.ScopeLte)
	}

	if f.ScopeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("scope")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		values = append(values, f.ScopeIn)
	}

	if f.ScopeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("scope")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("scope"))
		if *f.ScopeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("scope")+" IS NULL"+" OR "+aliasPrefix+SnakeString("scope")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("scope")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("scope")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *PermissionFilterType) AndWith(f2 ...*PermissionFilterType) *PermissionFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &PermissionFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *PermissionFilterType) OrWith(f2 ...*PermissionFilterType) *PermissionFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &PermissionFilterType{
		Or: append(_f2, f),
	}
}

func (f *OperatorRoleFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *OperatorRoleFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("operator_roles", ctx), wheres, values, joins)
}
func (f *OperatorRoleFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := f.Organization.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Members != nil {
		_alias := alias + "_members"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"role_id"+" LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"member_id"+" = "+_alias+".id")
		err := f.Members.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Permissions != nil {
		_alias := alias + "_permissions"
		*joins = append(*joins, "LEFT JOIN "+TableName("permission_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"role_id"+" LEFT JOIN "+TableName("permissions", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"permission_id"+" = "+_alias+".id")
		err := f.Permissions.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *OperatorRoleFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Name != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.Name)
	}

	if f.NameNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameNe)
	}

	if f.NameGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGt)
	}

	if f.NameLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLt)
	}

	if f.NameGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGte)
	}

	if f.NameLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLte)
	}

	if f.NameIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameIn)
	}

	if f.NameLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.NameLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.NamePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%s%%", *f.NamePrefix))
	}

	if f.NameSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%%%s", *f.NameSuffix))
	}

	if f.NameNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		if *f.NameNull {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NULL"+" OR "+aliasPrefix+SnakeString("name")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("name")+" <> ''")
		}
	}

	if f.Kind != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("kind")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		values = append(values, f.Kind)
	}

	if f.KindNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("kind")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		values = append(values, f.KindNe)
	}

	if f.KindGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("kind")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		values = append(values, f.KindGt)
	}

	if f.KindLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("kind")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		values = append(values, f.KindLt)
	}

	if f.KindGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("kind")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		values = append(values, f.KindGte)
	}

	if f.KindLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("kind")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		values = append(values, f.KindLte)
	}

	if f.KindIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("kind")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		values = append(values, f.KindIn)
	}

	if f.KindNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("kind")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("kind"))
		if *f.KindNull {
			conditions = append(conditions, aliasPrefix+SnakeString("kind")+" IS NULL"+" OR "+aliasPrefix+SnakeString("kind")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("kind")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("kind")+" <> ''")
		}
	}

	if f.OrganizationID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationID)
	}

	if f.OrganizationIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDNe)
	}

	if f.OrganizationIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGt)
	}

	if f.OrganizationIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLt)
	}

	if f.OrganizationIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGte)
	}

	if f.OrganizationIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLte)
	}

	if f.OrganizationIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDIn)
	}

	if f.OrganizationIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		if *f.OrganizationIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *OperatorRoleFilterType) AndWith(f2 ...*OperatorRoleFilterType) *OperatorRoleFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &OperatorRoleFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *OperatorRoleFilterType) OrWith(f2 ...*OperatorRoleFilterType) *OperatorRoleFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &OperatorRoleFilterType{
		Or: append(_f2, f),
	}
}

func (f *StoreFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *StoreFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("stores", ctx), wheres, values, joins)
}
func (f *StoreFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := f.Organization.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Members != nil {
		_alias := alias + "_members"
		*joins = append(*joins, "LEFT JOIN "+TableName("operatorMembership_stores", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"store_id"+" LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"member_id"+" = "+_alias+".id")
		err := f.Members.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.ReviewedByAccount != nil {
		_alias := alias + "_reviewedByAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"reviewed_by_account_id")
		err := f.ReviewedByAccount.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := f.AuditLogs.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.PaymentConfigs != nil {
		_alias := alias + "_paymentConfigs"
		*joins = append(*joins, "LEFT JOIN "+TableName("store_payment_configs", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")
		err := f.PaymentConfigs.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *StoreFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Code != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.Code)
	}

	if f.CodeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeNe)
	}

	if f.CodeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeGt)
	}

	if f.CodeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeLt)
	}

	if f.CodeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeGte)
	}

	if f.CodeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeLte)
	}

	if f.CodeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, f.CodeIn)
	}

	if f.CodeLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.CodeLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.CodePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, fmt.Sprintf("%s%%", *f.CodePrefix))
	}

	if f.CodeSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("code")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		values = append(values, fmt.Sprintf("%%%s", *f.CodeSuffix))
	}

	if f.CodeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("code")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("code"))
		if *f.CodeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("code")+" IS NULL"+" OR "+aliasPrefix+SnakeString("code")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("code")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("code")+" <> ''")
		}
	}

	if f.Name != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.Name)
	}

	if f.NameNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameNe)
	}

	if f.NameGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGt)
	}

	if f.NameLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLt)
	}

	if f.NameGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameGte)
	}

	if f.NameLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameLte)
	}

	if f.NameIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, f.NameIn)
	}

	if f.NameLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.NameLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.NamePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%s%%", *f.NamePrefix))
	}

	if f.NameSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("name")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		values = append(values, fmt.Sprintf("%%%s", *f.NameSuffix))
	}

	if f.NameNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("name")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("name"))
		if *f.NameNull {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NULL"+" OR "+aliasPrefix+SnakeString("name")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("name")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("name")+" <> ''")
		}
	}

	if f.Lifecycle != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		values = append(values, f.Lifecycle)
	}

	if f.LifecycleNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		values = append(values, f.LifecycleNe)
	}

	if f.LifecycleGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		values = append(values, f.LifecycleGt)
	}

	if f.LifecycleLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		values = append(values, f.LifecycleLt)
	}

	if f.LifecycleGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		values = append(values, f.LifecycleGte)
	}

	if f.LifecycleLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		values = append(values, f.LifecycleLte)
	}

	if f.LifecycleIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		values = append(values, f.LifecycleIn)
	}

	if f.LifecycleNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lifecycle")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lifecycle"))
		if *f.LifecycleNull {
			conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" IS NULL"+" OR "+aliasPrefix+SnakeString("lifecycle")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("lifecycle")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("lifecycle")+" <> ''")
		}
	}

	if f.SubmittedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		values = append(values, f.SubmittedAt)
	}

	if f.SubmittedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		values = append(values, f.SubmittedAtNe)
	}

	if f.SubmittedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		values = append(values, f.SubmittedAtGt)
	}

	if f.SubmittedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		values = append(values, f.SubmittedAtLt)
	}

	if f.SubmittedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		values = append(values, f.SubmittedAtGte)
	}

	if f.SubmittedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		values = append(values, f.SubmittedAtLte)
	}

	if f.SubmittedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		values = append(values, f.SubmittedAtIn)
	}

	if f.SubmittedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("submittedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("submittedAt"))
		if *f.SubmittedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("submittedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("submittedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("submittedAt")+" <> ''")
		}
	}

	if f.ReviewedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		values = append(values, f.ReviewedAt)
	}

	if f.ReviewedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		values = append(values, f.ReviewedAtNe)
	}

	if f.ReviewedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		values = append(values, f.ReviewedAtGt)
	}

	if f.ReviewedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		values = append(values, f.ReviewedAtLt)
	}

	if f.ReviewedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		values = append(values, f.ReviewedAtGte)
	}

	if f.ReviewedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		values = append(values, f.ReviewedAtLte)
	}

	if f.ReviewedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		values = append(values, f.ReviewedAtIn)
	}

	if f.ReviewedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedAt"))
		if *f.ReviewedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("reviewedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("reviewedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("reviewedAt")+" <> ''")
		}
	}

	if f.RejectionReason != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, f.RejectionReason)
	}

	if f.RejectionReasonNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, f.RejectionReasonNe)
	}

	if f.RejectionReasonGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, f.RejectionReasonGt)
	}

	if f.RejectionReasonLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, f.RejectionReasonLt)
	}

	if f.RejectionReasonGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, f.RejectionReasonGte)
	}

	if f.RejectionReasonLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, f.RejectionReasonLte)
	}

	if f.RejectionReasonIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, f.RejectionReasonIn)
	}

	if f.RejectionReasonLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.RejectionReasonLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.RejectionReasonPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, fmt.Sprintf("%s%%", *f.RejectionReasonPrefix))
	}

	if f.RejectionReasonSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		values = append(values, fmt.Sprintf("%%%s", *f.RejectionReasonSuffix))
	}

	if f.RejectionReasonNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("rejectionReason")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("rejectionReason"))
		if *f.RejectionReasonNull {
			conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" IS NULL"+" OR "+aliasPrefix+SnakeString("rejectionReason")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("rejectionReason")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("rejectionReason")+" <> ''")
		}
	}

	if f.ContactPhone != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, f.ContactPhone)
	}

	if f.ContactPhoneNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, f.ContactPhoneNe)
	}

	if f.ContactPhoneGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, f.ContactPhoneGt)
	}

	if f.ContactPhoneLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, f.ContactPhoneLt)
	}

	if f.ContactPhoneGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, f.ContactPhoneGte)
	}

	if f.ContactPhoneLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, f.ContactPhoneLte)
	}

	if f.ContactPhoneIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, f.ContactPhoneIn)
	}

	if f.ContactPhoneLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ContactPhoneLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ContactPhonePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, fmt.Sprintf("%s%%", *f.ContactPhonePrefix))
	}

	if f.ContactPhoneSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		values = append(values, fmt.Sprintf("%%%s", *f.ContactPhoneSuffix))
	}

	if f.ContactPhoneNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("contactPhone")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("contactPhone"))
		if *f.ContactPhoneNull {
			conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" IS NULL"+" OR "+aliasPrefix+SnakeString("contactPhone")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("contactPhone")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("contactPhone")+" <> ''")
		}
	}

	if f.ManagerName != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, f.ManagerName)
	}

	if f.ManagerNameNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, f.ManagerNameNe)
	}

	if f.ManagerNameGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, f.ManagerNameGt)
	}

	if f.ManagerNameLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, f.ManagerNameLt)
	}

	if f.ManagerNameGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, f.ManagerNameGte)
	}

	if f.ManagerNameLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, f.ManagerNameLte)
	}

	if f.ManagerNameIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, f.ManagerNameIn)
	}

	if f.ManagerNameLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ManagerNameLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ManagerNamePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, fmt.Sprintf("%s%%", *f.ManagerNamePrefix))
	}

	if f.ManagerNameSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		values = append(values, fmt.Sprintf("%%%s", *f.ManagerNameSuffix))
	}

	if f.ManagerNameNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerName")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerName"))
		if *f.ManagerNameNull {
			conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" IS NULL"+" OR "+aliasPrefix+SnakeString("managerName")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("managerName")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("managerName")+" <> ''")
		}
	}

	if f.ManagerPhone != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, f.ManagerPhone)
	}

	if f.ManagerPhoneNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, f.ManagerPhoneNe)
	}

	if f.ManagerPhoneGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, f.ManagerPhoneGt)
	}

	if f.ManagerPhoneLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, f.ManagerPhoneLt)
	}

	if f.ManagerPhoneGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, f.ManagerPhoneGte)
	}

	if f.ManagerPhoneLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, f.ManagerPhoneLte)
	}

	if f.ManagerPhoneIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, f.ManagerPhoneIn)
	}

	if f.ManagerPhoneLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ManagerPhoneLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ManagerPhonePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, fmt.Sprintf("%s%%", *f.ManagerPhonePrefix))
	}

	if f.ManagerPhoneSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		values = append(values, fmt.Sprintf("%%%s", *f.ManagerPhoneSuffix))
	}

	if f.ManagerPhoneNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("managerPhone")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("managerPhone"))
		if *f.ManagerPhoneNull {
			conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" IS NULL"+" OR "+aliasPrefix+SnakeString("managerPhone")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("managerPhone")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("managerPhone")+" <> ''")
		}
	}

	if f.Province != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, f.Province)
	}

	if f.ProvinceNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, f.ProvinceNe)
	}

	if f.ProvinceGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, f.ProvinceGt)
	}

	if f.ProvinceLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, f.ProvinceLt)
	}

	if f.ProvinceGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, f.ProvinceGte)
	}

	if f.ProvinceLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, f.ProvinceLte)
	}

	if f.ProvinceIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, f.ProvinceIn)
	}

	if f.ProvinceLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ProvinceLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ProvincePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, fmt.Sprintf("%s%%", *f.ProvincePrefix))
	}

	if f.ProvinceSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("province")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		values = append(values, fmt.Sprintf("%%%s", *f.ProvinceSuffix))
	}

	if f.ProvinceNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("province")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("province"))
		if *f.ProvinceNull {
			conditions = append(conditions, aliasPrefix+SnakeString("province")+" IS NULL"+" OR "+aliasPrefix+SnakeString("province")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("province")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("province")+" <> ''")
		}
	}

	if f.City != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, f.City)
	}

	if f.CityNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, f.CityNe)
	}

	if f.CityGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, f.CityGt)
	}

	if f.CityLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, f.CityLt)
	}

	if f.CityGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, f.CityGte)
	}

	if f.CityLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, f.CityLte)
	}

	if f.CityIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, f.CityIn)
	}

	if f.CityLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.CityLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.CityPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, fmt.Sprintf("%s%%", *f.CityPrefix))
	}

	if f.CitySuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("city")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		values = append(values, fmt.Sprintf("%%%s", *f.CitySuffix))
	}

	if f.CityNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("city")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("city"))
		if *f.CityNull {
			conditions = append(conditions, aliasPrefix+SnakeString("city")+" IS NULL"+" OR "+aliasPrefix+SnakeString("city")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("city")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("city")+" <> ''")
		}
	}

	if f.District != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, f.District)
	}

	if f.DistrictNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, f.DistrictNe)
	}

	if f.DistrictGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, f.DistrictGt)
	}

	if f.DistrictLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, f.DistrictLt)
	}

	if f.DistrictGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, f.DistrictGte)
	}

	if f.DistrictLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, f.DistrictLte)
	}

	if f.DistrictIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, f.DistrictIn)
	}

	if f.DistrictLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.DistrictLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.DistrictPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, fmt.Sprintf("%s%%", *f.DistrictPrefix))
	}

	if f.DistrictSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("district")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		values = append(values, fmt.Sprintf("%%%s", *f.DistrictSuffix))
	}

	if f.DistrictNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("district")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("district"))
		if *f.DistrictNull {
			conditions = append(conditions, aliasPrefix+SnakeString("district")+" IS NULL"+" OR "+aliasPrefix+SnakeString("district")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("district")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("district")+" <> ''")
		}
	}

	if f.Address != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, f.Address)
	}

	if f.AddressNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, f.AddressNe)
	}

	if f.AddressGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, f.AddressGt)
	}

	if f.AddressLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, f.AddressLt)
	}

	if f.AddressGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, f.AddressGte)
	}

	if f.AddressLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, f.AddressLte)
	}

	if f.AddressIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, f.AddressIn)
	}

	if f.AddressLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.AddressLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.AddressPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, fmt.Sprintf("%s%%", *f.AddressPrefix))
	}

	if f.AddressSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("address")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		values = append(values, fmt.Sprintf("%%%s", *f.AddressSuffix))
	}

	if f.AddressNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("address")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("address"))
		if *f.AddressNull {
			conditions = append(conditions, aliasPrefix+SnakeString("address")+" IS NULL"+" OR "+aliasPrefix+SnakeString("address")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("address")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("address")+" <> ''")
		}
	}

	if f.BusinessHours != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, f.BusinessHours)
	}

	if f.BusinessHoursNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, f.BusinessHoursNe)
	}

	if f.BusinessHoursGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, f.BusinessHoursGt)
	}

	if f.BusinessHoursLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, f.BusinessHoursLt)
	}

	if f.BusinessHoursGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, f.BusinessHoursGte)
	}

	if f.BusinessHoursLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, f.BusinessHoursLte)
	}

	if f.BusinessHoursIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, f.BusinessHoursIn)
	}

	if f.BusinessHoursLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.BusinessHoursLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.BusinessHoursPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, fmt.Sprintf("%s%%", *f.BusinessHoursPrefix))
	}

	if f.BusinessHoursSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		values = append(values, fmt.Sprintf("%%%s", *f.BusinessHoursSuffix))
	}

	if f.BusinessHoursNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessHours")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessHours"))
		if *f.BusinessHoursNull {
			conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" IS NULL"+" OR "+aliasPrefix+SnakeString("businessHours")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("businessHours")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("businessHours")+" <> ''")
		}
	}

	if f.BusinessStatus != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		values = append(values, f.BusinessStatus)
	}

	if f.BusinessStatusNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		values = append(values, f.BusinessStatusNe)
	}

	if f.BusinessStatusGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		values = append(values, f.BusinessStatusGt)
	}

	if f.BusinessStatusLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		values = append(values, f.BusinessStatusLt)
	}

	if f.BusinessStatusGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		values = append(values, f.BusinessStatusGte)
	}

	if f.BusinessStatusLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		values = append(values, f.BusinessStatusLte)
	}

	if f.BusinessStatusIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		values = append(values, f.BusinessStatusIn)
	}

	if f.BusinessStatusNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessStatus")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessStatus"))
		if *f.BusinessStatusNull {
			conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" IS NULL"+" OR "+aliasPrefix+SnakeString("businessStatus")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("businessStatus")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("businessStatus")+" <> ''")
		}
	}

	if f.SupportDineIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		values = append(values, f.SupportDineIn)
	}

	if f.SupportDineInNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		values = append(values, f.SupportDineInNe)
	}

	if f.SupportDineInGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		values = append(values, f.SupportDineInGt)
	}

	if f.SupportDineInLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		values = append(values, f.SupportDineInLt)
	}

	if f.SupportDineInGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		values = append(values, f.SupportDineInGte)
	}

	if f.SupportDineInLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		values = append(values, f.SupportDineInLte)
	}

	if f.SupportDineInIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		values = append(values, f.SupportDineInIn)
	}

	if f.SupportDineInNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportDineIn")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportDineIn"))
		if *f.SupportDineInNull {
			conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" IS NULL"+" OR "+aliasPrefix+SnakeString("supportDineIn")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("supportDineIn")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("supportDineIn")+" <> ''")
		}
	}

	if f.SupportTakeout != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		values = append(values, f.SupportTakeout)
	}

	if f.SupportTakeoutNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		values = append(values, f.SupportTakeoutNe)
	}

	if f.SupportTakeoutGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		values = append(values, f.SupportTakeoutGt)
	}

	if f.SupportTakeoutLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		values = append(values, f.SupportTakeoutLt)
	}

	if f.SupportTakeoutGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		values = append(values, f.SupportTakeoutGte)
	}

	if f.SupportTakeoutLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		values = append(values, f.SupportTakeoutLte)
	}

	if f.SupportTakeoutIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		values = append(values, f.SupportTakeoutIn)
	}

	if f.SupportTakeoutNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("supportTakeout")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("supportTakeout"))
		if *f.SupportTakeoutNull {
			conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" IS NULL"+" OR "+aliasPrefix+SnakeString("supportTakeout")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("supportTakeout")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("supportTakeout")+" <> ''")
		}
	}

	if f.StoreArea != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		values = append(values, f.StoreArea)
	}

	if f.StoreAreaNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		values = append(values, f.StoreAreaNe)
	}

	if f.StoreAreaGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		values = append(values, f.StoreAreaGt)
	}

	if f.StoreAreaLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		values = append(values, f.StoreAreaLt)
	}

	if f.StoreAreaGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		values = append(values, f.StoreAreaGte)
	}

	if f.StoreAreaLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		values = append(values, f.StoreAreaLte)
	}

	if f.StoreAreaIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		values = append(values, f.StoreAreaIn)
	}

	if f.StoreAreaNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeArea")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeArea"))
		if *f.StoreAreaNull {
			conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" IS NULL"+" OR "+aliasPrefix+SnakeString("storeArea")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("storeArea")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("storeArea")+" <> ''")
		}
	}

	if f.TableCount != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		values = append(values, f.TableCount)
	}

	if f.TableCountNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		values = append(values, f.TableCountNe)
	}

	if f.TableCountGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		values = append(values, f.TableCountGt)
	}

	if f.TableCountLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		values = append(values, f.TableCountLt)
	}

	if f.TableCountGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		values = append(values, f.TableCountGte)
	}

	if f.TableCountLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		values = append(values, f.TableCountLte)
	}

	if f.TableCountIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		values = append(values, f.TableCountIn)
	}

	if f.TableCountNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("tableCount")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("tableCount"))
		if *f.TableCountNull {
			conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" IS NULL"+" OR "+aliasPrefix+SnakeString("tableCount")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("tableCount")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("tableCount")+" <> ''")
		}
	}

	if f.ReceiptFooter != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, f.ReceiptFooter)
	}

	if f.ReceiptFooterNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, f.ReceiptFooterNe)
	}

	if f.ReceiptFooterGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, f.ReceiptFooterGt)
	}

	if f.ReceiptFooterLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, f.ReceiptFooterLt)
	}

	if f.ReceiptFooterGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, f.ReceiptFooterGte)
	}

	if f.ReceiptFooterLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, f.ReceiptFooterLte)
	}

	if f.ReceiptFooterIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, f.ReceiptFooterIn)
	}

	if f.ReceiptFooterLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ReceiptFooterLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ReceiptFooterPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, fmt.Sprintf("%s%%", *f.ReceiptFooterPrefix))
	}

	if f.ReceiptFooterSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		values = append(values, fmt.Sprintf("%%%s", *f.ReceiptFooterSuffix))
	}

	if f.ReceiptFooterNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("receiptFooter")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("receiptFooter"))
		if *f.ReceiptFooterNull {
			conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" IS NULL"+" OR "+aliasPrefix+SnakeString("receiptFooter")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("receiptFooter")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("receiptFooter")+" <> ''")
		}
	}

	if f.BusinessLicenseImageURL != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, f.BusinessLicenseImageURL)
	}

	if f.BusinessLicenseImageURLNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, f.BusinessLicenseImageURLNe)
	}

	if f.BusinessLicenseImageURLGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, f.BusinessLicenseImageURLGt)
	}

	if f.BusinessLicenseImageURLLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, f.BusinessLicenseImageURLLt)
	}

	if f.BusinessLicenseImageURLGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, f.BusinessLicenseImageURLGte)
	}

	if f.BusinessLicenseImageURLLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, f.BusinessLicenseImageURLLte)
	}

	if f.BusinessLicenseImageURLIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, f.BusinessLicenseImageURLIn)
	}

	if f.BusinessLicenseImageURLLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.BusinessLicenseImageURLLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.BusinessLicenseImageURLPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, fmt.Sprintf("%s%%", *f.BusinessLicenseImageURLPrefix))
	}

	if f.BusinessLicenseImageURLSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		values = append(values, fmt.Sprintf("%%%s", *f.BusinessLicenseImageURLSuffix))
	}

	if f.BusinessLicenseImageURLNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("businessLicenseImageUrl"))
		if *f.BusinessLicenseImageURLNull {
			conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" IS NULL"+" OR "+aliasPrefix+SnakeString("businessLicenseImageUrl")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("businessLicenseImageUrl")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("businessLicenseImageUrl")+" <> ''")
		}
	}

	if f.OtherDocumentImageURL != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, f.OtherDocumentImageURL)
	}

	if f.OtherDocumentImageURLNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, f.OtherDocumentImageURLNe)
	}

	if f.OtherDocumentImageURLGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, f.OtherDocumentImageURLGt)
	}

	if f.OtherDocumentImageURLLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, f.OtherDocumentImageURLLt)
	}

	if f.OtherDocumentImageURLGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, f.OtherDocumentImageURLGte)
	}

	if f.OtherDocumentImageURLLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, f.OtherDocumentImageURLLte)
	}

	if f.OtherDocumentImageURLIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, f.OtherDocumentImageURLIn)
	}

	if f.OtherDocumentImageURLLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.OtherDocumentImageURLLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.OtherDocumentImageURLPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, fmt.Sprintf("%s%%", *f.OtherDocumentImageURLPrefix))
	}

	if f.OtherDocumentImageURLSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		values = append(values, fmt.Sprintf("%%%s", *f.OtherDocumentImageURLSuffix))
	}

	if f.OtherDocumentImageURLNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("otherDocumentImageUrl"))
		if *f.OtherDocumentImageURLNull {
			conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" IS NULL"+" OR "+aliasPrefix+SnakeString("otherDocumentImageUrl")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("otherDocumentImageUrl")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("otherDocumentImageUrl")+" <> ''")
		}
	}

	if f.OrganizationID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationID)
	}

	if f.OrganizationIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDNe)
	}

	if f.OrganizationIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGt)
	}

	if f.OrganizationIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLt)
	}

	if f.OrganizationIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGte)
	}

	if f.OrganizationIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLte)
	}

	if f.OrganizationIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDIn)
	}

	if f.OrganizationIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		if *f.OrganizationIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" <> ''")
		}
	}

	if f.ReviewedByAccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		values = append(values, f.ReviewedByAccountID)
	}

	if f.ReviewedByAccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		values = append(values, f.ReviewedByAccountIDNe)
	}

	if f.ReviewedByAccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		values = append(values, f.ReviewedByAccountIDGt)
	}

	if f.ReviewedByAccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		values = append(values, f.ReviewedByAccountIDLt)
	}

	if f.ReviewedByAccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		values = append(values, f.ReviewedByAccountIDGte)
	}

	if f.ReviewedByAccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		values = append(values, f.ReviewedByAccountIDLte)
	}

	if f.ReviewedByAccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		values = append(values, f.ReviewedByAccountIDIn)
	}

	if f.ReviewedByAccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("reviewedByAccountId"))
		if *f.ReviewedByAccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("reviewedByAccountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("reviewedByAccountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("reviewedByAccountId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *StoreFilterType) AndWith(f2 ...*StoreFilterType) *StoreFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &StoreFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *StoreFilterType) OrWith(f2 ...*StoreFilterType) *StoreFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &StoreFilterType{
		Or: append(_f2, f),
	}
}

func (f *SessionFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *SessionFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("sessions", ctx), wheres, values, joins)
}
func (f *SessionFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Account != nil {
		_alias := alias + "_account"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"account_id")
		err := f.Account.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := f.Organization.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.AuditLogs != nil {
		_alias := alias + "_auditLogs"
		*joins = append(*joins, "LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"session_id"+" = "+alias+".id")
		err := f.AuditLogs.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *SessionFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.WorkspaceType != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		values = append(values, f.WorkspaceType)
	}

	if f.WorkspaceTypeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		values = append(values, f.WorkspaceTypeNe)
	}

	if f.WorkspaceTypeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		values = append(values, f.WorkspaceTypeGt)
	}

	if f.WorkspaceTypeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		values = append(values, f.WorkspaceTypeLt)
	}

	if f.WorkspaceTypeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		values = append(values, f.WorkspaceTypeGte)
	}

	if f.WorkspaceTypeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		values = append(values, f.WorkspaceTypeLte)
	}

	if f.WorkspaceTypeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		values = append(values, f.WorkspaceTypeIn)
	}

	if f.WorkspaceTypeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("workspaceType")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("workspaceType"))
		if *f.WorkspaceTypeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" IS NULL"+" OR "+aliasPrefix+SnakeString("workspaceType")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("workspaceType")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("workspaceType")+" <> ''")
		}
	}

	if f.CredentialVersion != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersion)
	}

	if f.CredentialVersionNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionNe)
	}

	if f.CredentialVersionGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionGt)
	}

	if f.CredentialVersionLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionLt)
	}

	if f.CredentialVersionGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionGte)
	}

	if f.CredentialVersionLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionLte)
	}

	if f.CredentialVersionIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		values = append(values, f.CredentialVersionIn)
	}

	if f.CredentialVersionNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialVersion")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialVersion"))
		if *f.CredentialVersionNull {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" IS NULL"+" OR "+aliasPrefix+SnakeString("credentialVersion")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialVersion")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("credentialVersion")+" <> ''")
		}
	}

	if f.ExpiresAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAt)
	}

	if f.ExpiresAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtNe)
	}

	if f.ExpiresAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtGt)
	}

	if f.ExpiresAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtLt)
	}

	if f.ExpiresAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtGte)
	}

	if f.ExpiresAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtLte)
	}

	if f.ExpiresAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtIn)
	}

	if f.ExpiresAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		if *f.ExpiresAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("expiresAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("expiresAt")+" <> ''")
		}
	}

	if f.RevokedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAt)
	}

	if f.RevokedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtNe)
	}

	if f.RevokedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtGt)
	}

	if f.RevokedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtLt)
	}

	if f.RevokedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtGte)
	}

	if f.RevokedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtLte)
	}

	if f.RevokedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtIn)
	}

	if f.RevokedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		if *f.RevokedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("revokedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("revokedAt")+" <> ''")
		}
	}

	if f.RevocationCode != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, f.RevocationCode)
	}

	if f.RevocationCodeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, f.RevocationCodeNe)
	}

	if f.RevocationCodeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, f.RevocationCodeGt)
	}

	if f.RevocationCodeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, f.RevocationCodeLt)
	}

	if f.RevocationCodeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, f.RevocationCodeGte)
	}

	if f.RevocationCodeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, f.RevocationCodeLte)
	}

	if f.RevocationCodeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, f.RevocationCodeIn)
	}

	if f.RevocationCodeLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.RevocationCodeLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.RevocationCodePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, fmt.Sprintf("%s%%", *f.RevocationCodePrefix))
	}

	if f.RevocationCodeSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		values = append(values, fmt.Sprintf("%%%s", *f.RevocationCodeSuffix))
	}

	if f.RevocationCodeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revocationCode")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revocationCode"))
		if *f.RevocationCodeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" IS NULL"+" OR "+aliasPrefix+SnakeString("revocationCode")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("revocationCode")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("revocationCode")+" <> ''")
		}
	}

	if f.LastSeenAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		values = append(values, f.LastSeenAt)
	}

	if f.LastSeenAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		values = append(values, f.LastSeenAtNe)
	}

	if f.LastSeenAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		values = append(values, f.LastSeenAtGt)
	}

	if f.LastSeenAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		values = append(values, f.LastSeenAtLt)
	}

	if f.LastSeenAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		values = append(values, f.LastSeenAtGte)
	}

	if f.LastSeenAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		values = append(values, f.LastSeenAtLte)
	}

	if f.LastSeenAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		values = append(values, f.LastSeenAtIn)
	}

	if f.LastSeenAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("lastSeenAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("lastSeenAt"))
		if *f.LastSeenAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("lastSeenAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("lastSeenAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("lastSeenAt")+" <> ''")
		}
	}

	if f.AccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountID)
	}

	if f.AccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDNe)
	}

	if f.AccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDGt)
	}

	if f.AccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDLt)
	}

	if f.AccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDGte)
	}

	if f.AccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDLte)
	}

	if f.AccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		values = append(values, f.AccountIDIn)
	}

	if f.AccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("accountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("accountId"))
		if *f.AccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("accountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("accountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("accountId")+" <> ''")
		}
	}

	if f.OrganizationID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationID)
	}

	if f.OrganizationIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDNe)
	}

	if f.OrganizationIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGt)
	}

	if f.OrganizationIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLt)
	}

	if f.OrganizationIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGte)
	}

	if f.OrganizationIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLte)
	}

	if f.OrganizationIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDIn)
	}

	if f.OrganizationIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		if *f.OrganizationIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *SessionFilterType) AndWith(f2 ...*SessionFilterType) *SessionFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &SessionFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *SessionFilterType) OrWith(f2 ...*SessionFilterType) *SessionFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &SessionFilterType{
		Or: append(_f2, f),
	}
}

func (f *MembershipInvitationFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *MembershipInvitationFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("membership_invitations", ctx), wheres, values, joins)
}
func (f *MembershipInvitationFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Membership != nil {
		_alias := alias + "_membership"
		*joins = append(*joins, "LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"membership_id")
		err := f.Membership.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.InvitedByAccount != nil {
		_alias := alias + "_invitedByAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"invited_by_account_id")
		err := f.InvitedByAccount.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *MembershipInvitationFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.ExpiresAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAt)
	}

	if f.ExpiresAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtNe)
	}

	if f.ExpiresAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtGt)
	}

	if f.ExpiresAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtLt)
	}

	if f.ExpiresAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtGte)
	}

	if f.ExpiresAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtLte)
	}

	if f.ExpiresAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		values = append(values, f.ExpiresAtIn)
	}

	if f.ExpiresAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("expiresAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("expiresAt"))
		if *f.ExpiresAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("expiresAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("expiresAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("expiresAt")+" <> ''")
		}
	}

	if f.AcceptedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAt)
	}

	if f.AcceptedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtNe)
	}

	if f.AcceptedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtGt)
	}

	if f.AcceptedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtLt)
	}

	if f.AcceptedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtGte)
	}

	if f.AcceptedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtLte)
	}

	if f.AcceptedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		values = append(values, f.AcceptedAtIn)
	}

	if f.AcceptedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("acceptedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("acceptedAt"))
		if *f.AcceptedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("acceptedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("acceptedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("acceptedAt")+" <> ''")
		}
	}

	if f.RevokedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAt)
	}

	if f.RevokedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtNe)
	}

	if f.RevokedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtGt)
	}

	if f.RevokedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtLt)
	}

	if f.RevokedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtGte)
	}

	if f.RevokedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtLte)
	}

	if f.RevokedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		values = append(values, f.RevokedAtIn)
	}

	if f.RevokedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("revokedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("revokedAt"))
		if *f.RevokedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("revokedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("revokedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("revokedAt")+" <> ''")
		}
	}

	if f.MembershipID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		values = append(values, f.MembershipID)
	}

	if f.MembershipIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		values = append(values, f.MembershipIDNe)
	}

	if f.MembershipIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		values = append(values, f.MembershipIDGt)
	}

	if f.MembershipIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		values = append(values, f.MembershipIDLt)
	}

	if f.MembershipIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		values = append(values, f.MembershipIDGte)
	}

	if f.MembershipIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		values = append(values, f.MembershipIDLte)
	}

	if f.MembershipIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		values = append(values, f.MembershipIDIn)
	}

	if f.MembershipIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("membershipId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("membershipId"))
		if *f.MembershipIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("membershipId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("membershipId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("membershipId")+" <> ''")
		}
	}

	if f.InvitedByAccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		values = append(values, f.InvitedByAccountID)
	}

	if f.InvitedByAccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		values = append(values, f.InvitedByAccountIDNe)
	}

	if f.InvitedByAccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		values = append(values, f.InvitedByAccountIDGt)
	}

	if f.InvitedByAccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		values = append(values, f.InvitedByAccountIDLt)
	}

	if f.InvitedByAccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		values = append(values, f.InvitedByAccountIDGte)
	}

	if f.InvitedByAccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		values = append(values, f.InvitedByAccountIDLte)
	}

	if f.InvitedByAccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		values = append(values, f.InvitedByAccountIDIn)
	}

	if f.InvitedByAccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("invitedByAccountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("invitedByAccountId"))
		if *f.InvitedByAccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("invitedByAccountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("invitedByAccountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("invitedByAccountId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *MembershipInvitationFilterType) AndWith(f2 ...*MembershipInvitationFilterType) *MembershipInvitationFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &MembershipInvitationFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *MembershipInvitationFilterType) OrWith(f2 ...*MembershipInvitationFilterType) *MembershipInvitationFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &MembershipInvitationFilterType{
		Or: append(_f2, f),
	}
}

func (f *AuditLogFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *AuditLogFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("audit_logs", ctx), wheres, values, joins)
}
func (f *AuditLogFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.ActorAccount != nil {
		_alias := alias + "_actorAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"actor_account_id")
		err := f.ActorAccount.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Session != nil {
		_alias := alias + "_session"
		*joins = append(*joins, "LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"session_id")
		err := f.Session.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := f.Organization.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := f.Store.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *AuditLogFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Action != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.Action)
	}

	if f.ActionNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionNe)
	}

	if f.ActionGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionGt)
	}

	if f.ActionLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionLt)
	}

	if f.ActionGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionGte)
	}

	if f.ActionLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionLte)
	}

	if f.ActionIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, f.ActionIn)
	}

	if f.ActionLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ActionLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ActionPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, fmt.Sprintf("%s%%", *f.ActionPrefix))
	}

	if f.ActionSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("action")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		values = append(values, fmt.Sprintf("%%%s", *f.ActionSuffix))
	}

	if f.ActionNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("action")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("action"))
		if *f.ActionNull {
			conditions = append(conditions, aliasPrefix+SnakeString("action")+" IS NULL"+" OR "+aliasPrefix+SnakeString("action")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("action")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("action")+" <> ''")
		}
	}

	if f.ResourceType != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, f.ResourceType)
	}

	if f.ResourceTypeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, f.ResourceTypeNe)
	}

	if f.ResourceTypeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, f.ResourceTypeGt)
	}

	if f.ResourceTypeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, f.ResourceTypeLt)
	}

	if f.ResourceTypeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, f.ResourceTypeGte)
	}

	if f.ResourceTypeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, f.ResourceTypeLte)
	}

	if f.ResourceTypeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, f.ResourceTypeIn)
	}

	if f.ResourceTypeLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ResourceTypeLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ResourceTypePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, fmt.Sprintf("%s%%", *f.ResourceTypePrefix))
	}

	if f.ResourceTypeSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		values = append(values, fmt.Sprintf("%%%s", *f.ResourceTypeSuffix))
	}

	if f.ResourceTypeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceType")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceType"))
		if *f.ResourceTypeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" IS NULL"+" OR "+aliasPrefix+SnakeString("resourceType")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("resourceType")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("resourceType")+" <> ''")
		}
	}

	if f.ResourceID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, f.ResourceID)
	}

	if f.ResourceIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, f.ResourceIDNe)
	}

	if f.ResourceIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, f.ResourceIDGt)
	}

	if f.ResourceIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, f.ResourceIDLt)
	}

	if f.ResourceIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, f.ResourceIDGte)
	}

	if f.ResourceIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, f.ResourceIDLte)
	}

	if f.ResourceIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, f.ResourceIDIn)
	}

	if f.ResourceIDLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ResourceIDLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ResourceIDPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, fmt.Sprintf("%s%%", *f.ResourceIDPrefix))
	}

	if f.ResourceIDSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		values = append(values, fmt.Sprintf("%%%s", *f.ResourceIDSuffix))
	}

	if f.ResourceIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resourceId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resourceId"))
		if *f.ResourceIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("resourceId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("resourceId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("resourceId")+" <> ''")
		}
	}

	if f.ResultCode != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, f.ResultCode)
	}

	if f.ResultCodeNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, f.ResultCodeNe)
	}

	if f.ResultCodeGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, f.ResultCodeGt)
	}

	if f.ResultCodeLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, f.ResultCodeLt)
	}

	if f.ResultCodeGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, f.ResultCodeGte)
	}

	if f.ResultCodeLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, f.ResultCodeLte)
	}

	if f.ResultCodeIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, f.ResultCodeIn)
	}

	if f.ResultCodeLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ResultCodeLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ResultCodePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, fmt.Sprintf("%s%%", *f.ResultCodePrefix))
	}

	if f.ResultCodeSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		values = append(values, fmt.Sprintf("%%%s", *f.ResultCodeSuffix))
	}

	if f.ResultCodeNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("resultCode")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("resultCode"))
		if *f.ResultCodeNull {
			conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" IS NULL"+" OR "+aliasPrefix+SnakeString("resultCode")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("resultCode")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("resultCode")+" <> ''")
		}
	}

	if f.MetadataJSON != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, f.MetadataJSON)
	}

	if f.MetadataJSONNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, f.MetadataJSONNe)
	}

	if f.MetadataJSONGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, f.MetadataJSONGt)
	}

	if f.MetadataJSONLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, f.MetadataJSONLt)
	}

	if f.MetadataJSONGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, f.MetadataJSONGte)
	}

	if f.MetadataJSONLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, f.MetadataJSONLte)
	}

	if f.MetadataJSONIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, f.MetadataJSONIn)
	}

	if f.MetadataJSONLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.MetadataJSONLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.MetadataJSONPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, fmt.Sprintf("%s%%", *f.MetadataJSONPrefix))
	}

	if f.MetadataJSONSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		values = append(values, fmt.Sprintf("%%%s", *f.MetadataJSONSuffix))
	}

	if f.MetadataJSONNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("metadataJson")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("metadataJson"))
		if *f.MetadataJSONNull {
			conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" IS NULL"+" OR "+aliasPrefix+SnakeString("metadataJson")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("metadataJson")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("metadataJson")+" <> ''")
		}
	}

	if f.ActorAccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		values = append(values, f.ActorAccountID)
	}

	if f.ActorAccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		values = append(values, f.ActorAccountIDNe)
	}

	if f.ActorAccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		values = append(values, f.ActorAccountIDGt)
	}

	if f.ActorAccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		values = append(values, f.ActorAccountIDLt)
	}

	if f.ActorAccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		values = append(values, f.ActorAccountIDGte)
	}

	if f.ActorAccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		values = append(values, f.ActorAccountIDLte)
	}

	if f.ActorAccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		values = append(values, f.ActorAccountIDIn)
	}

	if f.ActorAccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("actorAccountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("actorAccountId"))
		if *f.ActorAccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("actorAccountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("actorAccountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("actorAccountId")+" <> ''")
		}
	}

	if f.SessionID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		values = append(values, f.SessionID)
	}

	if f.SessionIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		values = append(values, f.SessionIDNe)
	}

	if f.SessionIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		values = append(values, f.SessionIDGt)
	}

	if f.SessionIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		values = append(values, f.SessionIDLt)
	}

	if f.SessionIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		values = append(values, f.SessionIDGte)
	}

	if f.SessionIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		values = append(values, f.SessionIDLte)
	}

	if f.SessionIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		values = append(values, f.SessionIDIn)
	}

	if f.SessionIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("sessionId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("sessionId"))
		if *f.SessionIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("sessionId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("sessionId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("sessionId")+" <> ''")
		}
	}

	if f.OrganizationID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationID)
	}

	if f.OrganizationIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDNe)
	}

	if f.OrganizationIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGt)
	}

	if f.OrganizationIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLt)
	}

	if f.OrganizationIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGte)
	}

	if f.OrganizationIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLte)
	}

	if f.OrganizationIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDIn)
	}

	if f.OrganizationIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		if *f.OrganizationIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" <> ''")
		}
	}

	if f.StoreID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreID)
	}

	if f.StoreIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDNe)
	}

	if f.StoreIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDGt)
	}

	if f.StoreIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDLt)
	}

	if f.StoreIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDGte)
	}

	if f.StoreIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDLte)
	}

	if f.StoreIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDIn)
	}

	if f.StoreIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		if *f.StoreIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("storeId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("storeId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *AuditLogFilterType) AndWith(f2 ...*AuditLogFilterType) *AuditLogFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &AuditLogFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *AuditLogFilterType) OrWith(f2 ...*AuditLogFilterType) *AuditLogFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &AuditLogFilterType{
		Or: append(_f2, f),
	}
}

func (f *FranchiseOpeningRecordFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *FranchiseOpeningRecordFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("franchise_opening_records", ctx), wheres, values, joins)
}
func (f *FranchiseOpeningRecordFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := f.Organization.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.InitialAccount != nil {
		_alias := alias + "_initialAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initial_account_id")
		err := f.InitialAccount.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	if f.RecordedByAccount != nil {
		_alias := alias + "_recordedByAccount"
		*joins = append(*joins, "LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"recorded_by_account_id")
		err := f.RecordedByAccount.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FranchiseOpeningRecordFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.RecordNumber != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, f.RecordNumber)
	}

	if f.RecordNumberNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, f.RecordNumberNe)
	}

	if f.RecordNumberGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, f.RecordNumberGt)
	}

	if f.RecordNumberLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, f.RecordNumberLt)
	}

	if f.RecordNumberGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, f.RecordNumberGte)
	}

	if f.RecordNumberLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, f.RecordNumberLte)
	}

	if f.RecordNumberIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, f.RecordNumberIn)
	}

	if f.RecordNumberLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.RecordNumberLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.RecordNumberPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, fmt.Sprintf("%s%%", *f.RecordNumberPrefix))
	}

	if f.RecordNumberSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		values = append(values, fmt.Sprintf("%%%s", *f.RecordNumberSuffix))
	}

	if f.RecordNumberNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordNumber")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordNumber"))
		if *f.RecordNumberNull {
			conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" IS NULL"+" OR "+aliasPrefix+SnakeString("recordNumber")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("recordNumber")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("recordNumber")+" <> ''")
		}
	}

	if f.Source != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("source")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		values = append(values, f.Source)
	}

	if f.SourceNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("source")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		values = append(values, f.SourceNe)
	}

	if f.SourceGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("source")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		values = append(values, f.SourceGt)
	}

	if f.SourceLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("source")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		values = append(values, f.SourceLt)
	}

	if f.SourceGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("source")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		values = append(values, f.SourceGte)
	}

	if f.SourceLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("source")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		values = append(values, f.SourceLte)
	}

	if f.SourceIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("source")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		values = append(values, f.SourceIn)
	}

	if f.SourceNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("source")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("source"))
		if *f.SourceNull {
			conditions = append(conditions, aliasPrefix+SnakeString("source")+" IS NULL"+" OR "+aliasPrefix+SnakeString("source")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("source")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("source")+" <> ''")
		}
	}

	if f.OrganizationID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationID)
	}

	if f.OrganizationIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDNe)
	}

	if f.OrganizationIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGt)
	}

	if f.OrganizationIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLt)
	}

	if f.OrganizationIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGte)
	}

	if f.OrganizationIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLte)
	}

	if f.OrganizationIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDIn)
	}

	if f.OrganizationIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		if *f.OrganizationIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" <> ''")
		}
	}

	if f.InitialAccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountID)
	}

	if f.InitialAccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDNe)
	}

	if f.InitialAccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDGt)
	}

	if f.InitialAccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDLt)
	}

	if f.InitialAccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDGte)
	}

	if f.InitialAccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDLte)
	}

	if f.InitialAccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		values = append(values, f.InitialAccountIDIn)
	}

	if f.InitialAccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("initialAccountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("initialAccountId"))
		if *f.InitialAccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("initialAccountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("initialAccountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("initialAccountId")+" <> ''")
		}
	}

	if f.RecordedByAccountID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		values = append(values, f.RecordedByAccountID)
	}

	if f.RecordedByAccountIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		values = append(values, f.RecordedByAccountIDNe)
	}

	if f.RecordedByAccountIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		values = append(values, f.RecordedByAccountIDGt)
	}

	if f.RecordedByAccountIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		values = append(values, f.RecordedByAccountIDLt)
	}

	if f.RecordedByAccountIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		values = append(values, f.RecordedByAccountIDGte)
	}

	if f.RecordedByAccountIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		values = append(values, f.RecordedByAccountIDLte)
	}

	if f.RecordedByAccountIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		values = append(values, f.RecordedByAccountIDIn)
	}

	if f.RecordedByAccountIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("recordedByAccountId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("recordedByAccountId"))
		if *f.RecordedByAccountIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("recordedByAccountId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("recordedByAccountId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("recordedByAccountId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *FranchiseOpeningRecordFilterType) AndWith(f2 ...*FranchiseOpeningRecordFilterType) *FranchiseOpeningRecordFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &FranchiseOpeningRecordFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *FranchiseOpeningRecordFilterType) OrWith(f2 ...*FranchiseOpeningRecordFilterType) *FranchiseOpeningRecordFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &FranchiseOpeningRecordFilterType{
		Or: append(_f2, f),
	}
}

func (f *GlobalPaymentConfigFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *GlobalPaymentConfigFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("global_payment_configs", ctx), wheres, values, joins)
}
func (f *GlobalPaymentConfigFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	return nil
}

func (f *GlobalPaymentConfigFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Channel != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.Channel)
	}

	if f.ChannelNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelNe)
	}

	if f.ChannelGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelGt)
	}

	if f.ChannelLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelLt)
	}

	if f.ChannelGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelGte)
	}

	if f.ChannelLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelLte)
	}

	if f.ChannelIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelIn)
	}

	if f.ChannelLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ChannelLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ChannelPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, fmt.Sprintf("%s%%", *f.ChannelPrefix))
	}

	if f.ChannelSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, fmt.Sprintf("%%%s", *f.ChannelSuffix))
	}

	if f.ChannelNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		if *f.ChannelNull {
			conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IS NULL"+" OR "+aliasPrefix+SnakeString("channel")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("channel")+" <> ''")
		}
	}

	if f.MerchantID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantID)
	}

	if f.MerchantIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDNe)
	}

	if f.MerchantIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDGt)
	}

	if f.MerchantIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDLt)
	}

	if f.MerchantIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDGte)
	}

	if f.MerchantIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDLte)
	}

	if f.MerchantIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDIn)
	}

	if f.MerchantIDLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.MerchantIDLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.MerchantIDPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, fmt.Sprintf("%s%%", *f.MerchantIDPrefix))
	}

	if f.MerchantIDSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, fmt.Sprintf("%%%s", *f.MerchantIDSuffix))
	}

	if f.MerchantIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		if *f.MerchantIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("merchantId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("merchantId")+" <> ''")
		}
	}

	if f.Environment != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.Environment)
	}

	if f.EnvironmentNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentNe)
	}

	if f.EnvironmentGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentGt)
	}

	if f.EnvironmentLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentLt)
	}

	if f.EnvironmentGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentGte)
	}

	if f.EnvironmentLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentLte)
	}

	if f.EnvironmentIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentIn)
	}

	if f.EnvironmentLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.EnvironmentLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.EnvironmentPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, fmt.Sprintf("%s%%", *f.EnvironmentPrefix))
	}

	if f.EnvironmentSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, fmt.Sprintf("%%%s", *f.EnvironmentSuffix))
	}

	if f.EnvironmentNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		if *f.EnvironmentNull {
			conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IS NULL"+" OR "+aliasPrefix+SnakeString("environment")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("environment")+" <> ''")
		}
	}

	if f.RatePpm != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpm)
	}

	if f.RatePpmNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmNe)
	}

	if f.RatePpmGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmGt)
	}

	if f.RatePpmLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmLt)
	}

	if f.RatePpmGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmGte)
	}

	if f.RatePpmLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmLte)
	}

	if f.RatePpmIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmIn)
	}

	if f.RatePpmNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		if *f.RatePpmNull {
			conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IS NULL"+" OR "+aliasPrefix+SnakeString("ratePpm")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("ratePpm")+" <> ''")
		}
	}

	if f.ConfigState != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigState)
	}

	if f.ConfigStateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateNe)
	}

	if f.ConfigStateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateGt)
	}

	if f.ConfigStateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateLt)
	}

	if f.ConfigStateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateGte)
	}

	if f.ConfigStateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateLte)
	}

	if f.ConfigStateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateIn)
	}

	if f.ConfigStateLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ConfigStateLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ConfigStatePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, fmt.Sprintf("%s%%", *f.ConfigStatePrefix))
	}

	if f.ConfigStateSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, fmt.Sprintf("%%%s", *f.ConfigStateSuffix))
	}

	if f.ConfigStateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		if *f.ConfigStateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IS NULL"+" OR "+aliasPrefix+SnakeString("configState")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("configState")+" <> ''")
		}
	}

	if f.Version != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.Version)
	}

	if f.VersionNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionNe)
	}

	if f.VersionGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionGt)
	}

	if f.VersionLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionLt)
	}

	if f.VersionGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionGte)
	}

	if f.VersionLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionLte)
	}

	if f.VersionIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionIn)
	}

	if f.VersionNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		if *f.VersionNull {
			conditions = append(conditions, aliasPrefix+SnakeString("version")+" IS NULL"+" OR "+aliasPrefix+SnakeString("version")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("version")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("version")+" <> ''")
		}
	}

	if f.KeyID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyID)
	}

	if f.KeyIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDNe)
	}

	if f.KeyIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDGt)
	}

	if f.KeyIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDLt)
	}

	if f.KeyIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDGte)
	}

	if f.KeyIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDLte)
	}

	if f.KeyIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDIn)
	}

	if f.KeyIDLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.KeyIDLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.KeyIDPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, fmt.Sprintf("%s%%", *f.KeyIDPrefix))
	}

	if f.KeyIDSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, fmt.Sprintf("%%%s", *f.KeyIDSuffix))
	}

	if f.KeyIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		if *f.KeyIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("keyId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("keyId")+" <> ''")
		}
	}

	if f.CredentialCiphertext != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertext)
	}

	if f.CredentialCiphertextNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextNe)
	}

	if f.CredentialCiphertextGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextGt)
	}

	if f.CredentialCiphertextLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextLt)
	}

	if f.CredentialCiphertextGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextGte)
	}

	if f.CredentialCiphertextLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextLte)
	}

	if f.CredentialCiphertextIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextIn)
	}

	if f.CredentialCiphertextLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.CredentialCiphertextLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.CredentialCiphertextPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, fmt.Sprintf("%s%%", *f.CredentialCiphertextPrefix))
	}

	if f.CredentialCiphertextSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, fmt.Sprintf("%%%s", *f.CredentialCiphertextSuffix))
	}

	if f.CredentialCiphertextNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		if *f.CredentialCiphertextNull {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IS NULL"+" OR "+aliasPrefix+SnakeString("credentialCiphertext")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("credentialCiphertext")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *GlobalPaymentConfigFilterType) AndWith(f2 ...*GlobalPaymentConfigFilterType) *GlobalPaymentConfigFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &GlobalPaymentConfigFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *GlobalPaymentConfigFilterType) OrWith(f2 ...*GlobalPaymentConfigFilterType) *GlobalPaymentConfigFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &GlobalPaymentConfigFilterType{
		Or: append(_f2, f),
	}
}

func (f *FranchisePaymentConfigFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *FranchisePaymentConfigFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("franchise_payment_configs", ctx), wheres, values, joins)
}
func (f *FranchisePaymentConfigFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Organization != nil {
		_alias := alias + "_organization"
		*joins = append(*joins, "LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")
		err := f.Organization.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FranchisePaymentConfigFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Channel != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.Channel)
	}

	if f.ChannelNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelNe)
	}

	if f.ChannelGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelGt)
	}

	if f.ChannelLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelLt)
	}

	if f.ChannelGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelGte)
	}

	if f.ChannelLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelLte)
	}

	if f.ChannelIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelIn)
	}

	if f.ChannelLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ChannelLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ChannelPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, fmt.Sprintf("%s%%", *f.ChannelPrefix))
	}

	if f.ChannelSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, fmt.Sprintf("%%%s", *f.ChannelSuffix))
	}

	if f.ChannelNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		if *f.ChannelNull {
			conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IS NULL"+" OR "+aliasPrefix+SnakeString("channel")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("channel")+" <> ''")
		}
	}

	if f.MerchantID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantID)
	}

	if f.MerchantIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDNe)
	}

	if f.MerchantIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDGt)
	}

	if f.MerchantIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDLt)
	}

	if f.MerchantIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDGte)
	}

	if f.MerchantIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDLte)
	}

	if f.MerchantIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDIn)
	}

	if f.MerchantIDLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.MerchantIDLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.MerchantIDPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, fmt.Sprintf("%s%%", *f.MerchantIDPrefix))
	}

	if f.MerchantIDSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, fmt.Sprintf("%%%s", *f.MerchantIDSuffix))
	}

	if f.MerchantIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		if *f.MerchantIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("merchantId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("merchantId")+" <> ''")
		}
	}

	if f.Environment != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.Environment)
	}

	if f.EnvironmentNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentNe)
	}

	if f.EnvironmentGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentGt)
	}

	if f.EnvironmentLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentLt)
	}

	if f.EnvironmentGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentGte)
	}

	if f.EnvironmentLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentLte)
	}

	if f.EnvironmentIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentIn)
	}

	if f.EnvironmentLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.EnvironmentLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.EnvironmentPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, fmt.Sprintf("%s%%", *f.EnvironmentPrefix))
	}

	if f.EnvironmentSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, fmt.Sprintf("%%%s", *f.EnvironmentSuffix))
	}

	if f.EnvironmentNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		if *f.EnvironmentNull {
			conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IS NULL"+" OR "+aliasPrefix+SnakeString("environment")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("environment")+" <> ''")
		}
	}

	if f.RatePpm != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpm)
	}

	if f.RatePpmNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmNe)
	}

	if f.RatePpmGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmGt)
	}

	if f.RatePpmLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmLt)
	}

	if f.RatePpmGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmGte)
	}

	if f.RatePpmLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmLte)
	}

	if f.RatePpmIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmIn)
	}

	if f.RatePpmNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		if *f.RatePpmNull {
			conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IS NULL"+" OR "+aliasPrefix+SnakeString("ratePpm")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("ratePpm")+" <> ''")
		}
	}

	if f.ConfigState != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigState)
	}

	if f.ConfigStateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateNe)
	}

	if f.ConfigStateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateGt)
	}

	if f.ConfigStateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateLt)
	}

	if f.ConfigStateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateGte)
	}

	if f.ConfigStateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateLte)
	}

	if f.ConfigStateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateIn)
	}

	if f.ConfigStateLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ConfigStateLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ConfigStatePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, fmt.Sprintf("%s%%", *f.ConfigStatePrefix))
	}

	if f.ConfigStateSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, fmt.Sprintf("%%%s", *f.ConfigStateSuffix))
	}

	if f.ConfigStateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		if *f.ConfigStateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IS NULL"+" OR "+aliasPrefix+SnakeString("configState")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("configState")+" <> ''")
		}
	}

	if f.Version != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.Version)
	}

	if f.VersionNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionNe)
	}

	if f.VersionGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionGt)
	}

	if f.VersionLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionLt)
	}

	if f.VersionGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionGte)
	}

	if f.VersionLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionLte)
	}

	if f.VersionIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionIn)
	}

	if f.VersionNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		if *f.VersionNull {
			conditions = append(conditions, aliasPrefix+SnakeString("version")+" IS NULL"+" OR "+aliasPrefix+SnakeString("version")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("version")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("version")+" <> ''")
		}
	}

	if f.KeyID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyID)
	}

	if f.KeyIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDNe)
	}

	if f.KeyIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDGt)
	}

	if f.KeyIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDLt)
	}

	if f.KeyIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDGte)
	}

	if f.KeyIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDLte)
	}

	if f.KeyIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDIn)
	}

	if f.KeyIDLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.KeyIDLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.KeyIDPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, fmt.Sprintf("%s%%", *f.KeyIDPrefix))
	}

	if f.KeyIDSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, fmt.Sprintf("%%%s", *f.KeyIDSuffix))
	}

	if f.KeyIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		if *f.KeyIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("keyId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("keyId")+" <> ''")
		}
	}

	if f.CredentialCiphertext != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertext)
	}

	if f.CredentialCiphertextNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextNe)
	}

	if f.CredentialCiphertextGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextGt)
	}

	if f.CredentialCiphertextLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextLt)
	}

	if f.CredentialCiphertextGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextGte)
	}

	if f.CredentialCiphertextLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextLte)
	}

	if f.CredentialCiphertextIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextIn)
	}

	if f.CredentialCiphertextLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.CredentialCiphertextLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.CredentialCiphertextPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, fmt.Sprintf("%s%%", *f.CredentialCiphertextPrefix))
	}

	if f.CredentialCiphertextSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, fmt.Sprintf("%%%s", *f.CredentialCiphertextSuffix))
	}

	if f.CredentialCiphertextNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		if *f.CredentialCiphertextNull {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IS NULL"+" OR "+aliasPrefix+SnakeString("credentialCiphertext")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("credentialCiphertext")+" <> ''")
		}
	}

	if f.OrganizationID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationID)
	}

	if f.OrganizationIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDNe)
	}

	if f.OrganizationIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGt)
	}

	if f.OrganizationIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLt)
	}

	if f.OrganizationIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDGte)
	}

	if f.OrganizationIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDLte)
	}

	if f.OrganizationIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		values = append(values, f.OrganizationIDIn)
	}

	if f.OrganizationIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("organizationId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("organizationId"))
		if *f.OrganizationIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("organizationId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("organizationId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *FranchisePaymentConfigFilterType) AndWith(f2 ...*FranchisePaymentConfigFilterType) *FranchisePaymentConfigFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &FranchisePaymentConfigFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *FranchisePaymentConfigFilterType) OrWith(f2 ...*FranchisePaymentConfigFilterType) *FranchisePaymentConfigFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &FranchisePaymentConfigFilterType{
		Or: append(_f2, f),
	}
}

func (f *StorePaymentConfigFilterType) IsEmpty(ctx context.Context) bool {
	wheres := []string{}
	values := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, "companies", &wheres, &values, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0
}
func (f *StorePaymentConfigFilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, TableName("store_payment_configs", ctx), wheres, values, joins)
}
func (f *StorePaymentConfigFilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."

	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_cs := []string{}
			err := or.ApplyWithAlias(ctx, alias, &_cs, &vs, &js)
			if err != nil {
				return err
			}
			cs = append(cs, strings.Join(_cs, " AND "))
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		if len(cs) > 0 {
			*wheres = append(*wheres, strings.Join(cs, " AND "))
		}
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}

	if f.Store != nil {
		_alias := alias + "_store"
		*joins = append(*joins, "LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")
		err := f.Store.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *StorePaymentConfigFilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}
	whereConditions := []string{}

	if f.ID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.ID)
	}

	if f.IDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDNe)
	}

	if f.IDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGt)
	}

	if f.IDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLt)
	}

	if f.IDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDGte)
	}

	if f.IDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDLte)
	}

	if f.IDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("id")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		values = append(values, f.IDIn)
	}

	if f.IDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("id")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("id"))
		if *f.IDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NULL"+" OR "+aliasPrefix+SnakeString("id")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("id")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("id")+" <> ''")
		}
	}

	if f.Channel != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.Channel)
	}

	if f.ChannelNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelNe)
	}

	if f.ChannelGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelGt)
	}

	if f.ChannelLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelLt)
	}

	if f.ChannelGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelGte)
	}

	if f.ChannelLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelLte)
	}

	if f.ChannelIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, f.ChannelIn)
	}

	if f.ChannelLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ChannelLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ChannelPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, fmt.Sprintf("%s%%", *f.ChannelPrefix))
	}

	if f.ChannelSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("channel")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		values = append(values, fmt.Sprintf("%%%s", *f.ChannelSuffix))
	}

	if f.ChannelNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("channel")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("channel"))
		if *f.ChannelNull {
			conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IS NULL"+" OR "+aliasPrefix+SnakeString("channel")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("channel")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("channel")+" <> ''")
		}
	}

	if f.MerchantID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantID)
	}

	if f.MerchantIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDNe)
	}

	if f.MerchantIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDGt)
	}

	if f.MerchantIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDLt)
	}

	if f.MerchantIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDGte)
	}

	if f.MerchantIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDLte)
	}

	if f.MerchantIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, f.MerchantIDIn)
	}

	if f.MerchantIDLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.MerchantIDLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.MerchantIDPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, fmt.Sprintf("%s%%", *f.MerchantIDPrefix))
	}

	if f.MerchantIDSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		values = append(values, fmt.Sprintf("%%%s", *f.MerchantIDSuffix))
	}

	if f.MerchantIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("merchantId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("merchantId"))
		if *f.MerchantIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("merchantId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("merchantId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("merchantId")+" <> ''")
		}
	}

	if f.Environment != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.Environment)
	}

	if f.EnvironmentNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentNe)
	}

	if f.EnvironmentGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentGt)
	}

	if f.EnvironmentLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentLt)
	}

	if f.EnvironmentGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentGte)
	}

	if f.EnvironmentLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentLte)
	}

	if f.EnvironmentIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, f.EnvironmentIn)
	}

	if f.EnvironmentLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.EnvironmentLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.EnvironmentPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, fmt.Sprintf("%s%%", *f.EnvironmentPrefix))
	}

	if f.EnvironmentSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("environment")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		values = append(values, fmt.Sprintf("%%%s", *f.EnvironmentSuffix))
	}

	if f.EnvironmentNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("environment")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("environment"))
		if *f.EnvironmentNull {
			conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IS NULL"+" OR "+aliasPrefix+SnakeString("environment")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("environment")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("environment")+" <> ''")
		}
	}

	if f.RatePpm != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpm)
	}

	if f.RatePpmNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmNe)
	}

	if f.RatePpmGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmGt)
	}

	if f.RatePpmLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmLt)
	}

	if f.RatePpmGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmGte)
	}

	if f.RatePpmLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmLte)
	}

	if f.RatePpmIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		values = append(values, f.RatePpmIn)
	}

	if f.RatePpmNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("ratePpm")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("ratePpm"))
		if *f.RatePpmNull {
			conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IS NULL"+" OR "+aliasPrefix+SnakeString("ratePpm")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("ratePpm")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("ratePpm")+" <> ''")
		}
	}

	if f.ConfigState != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigState)
	}

	if f.ConfigStateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateNe)
	}

	if f.ConfigStateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateGt)
	}

	if f.ConfigStateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateLt)
	}

	if f.ConfigStateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateGte)
	}

	if f.ConfigStateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateLte)
	}

	if f.ConfigStateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, f.ConfigStateIn)
	}

	if f.ConfigStateLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.ConfigStateLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.ConfigStatePrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, fmt.Sprintf("%s%%", *f.ConfigStatePrefix))
	}

	if f.ConfigStateSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("configState")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		values = append(values, fmt.Sprintf("%%%s", *f.ConfigStateSuffix))
	}

	if f.ConfigStateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("configState")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("configState"))
		if *f.ConfigStateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IS NULL"+" OR "+aliasPrefix+SnakeString("configState")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("configState")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("configState")+" <> ''")
		}
	}

	if f.Version != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.Version)
	}

	if f.VersionNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionNe)
	}

	if f.VersionGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionGt)
	}

	if f.VersionLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionLt)
	}

	if f.VersionGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionGte)
	}

	if f.VersionLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionLte)
	}

	if f.VersionIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("version")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		values = append(values, f.VersionIn)
	}

	if f.VersionNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("version")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("version"))
		if *f.VersionNull {
			conditions = append(conditions, aliasPrefix+SnakeString("version")+" IS NULL"+" OR "+aliasPrefix+SnakeString("version")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("version")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("version")+" <> ''")
		}
	}

	if f.KeyID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyID)
	}

	if f.KeyIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDNe)
	}

	if f.KeyIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDGt)
	}

	if f.KeyIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDLt)
	}

	if f.KeyIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDGte)
	}

	if f.KeyIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDLte)
	}

	if f.KeyIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, f.KeyIDIn)
	}

	if f.KeyIDLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.KeyIDLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.KeyIDPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, fmt.Sprintf("%s%%", *f.KeyIDPrefix))
	}

	if f.KeyIDSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		values = append(values, fmt.Sprintf("%%%s", *f.KeyIDSuffix))
	}

	if f.KeyIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("keyId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("keyId"))
		if *f.KeyIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("keyId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("keyId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("keyId")+" <> ''")
		}
	}

	if f.CredentialCiphertext != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertext)
	}

	if f.CredentialCiphertextNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextNe)
	}

	if f.CredentialCiphertextGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextGt)
	}

	if f.CredentialCiphertextLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextLt)
	}

	if f.CredentialCiphertextGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextGte)
	}

	if f.CredentialCiphertextLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextLte)
	}

	if f.CredentialCiphertextIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, f.CredentialCiphertextIn)
	}

	if f.CredentialCiphertextLike != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, "%"+strings.Replace(strings.Replace(*f.CredentialCiphertextLike, "?", "_", -1), "*", "%", -1)+"%")
	}

	if f.CredentialCiphertextPrefix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, fmt.Sprintf("%s%%", *f.CredentialCiphertextPrefix))
	}

	if f.CredentialCiphertextSuffix != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" LIKE ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		values = append(values, fmt.Sprintf("%%%s", *f.CredentialCiphertextSuffix))
	}

	if f.CredentialCiphertextNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("credentialCiphertext")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("credentialCiphertext"))
		if *f.CredentialCiphertextNull {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IS NULL"+" OR "+aliasPrefix+SnakeString("credentialCiphertext")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("credentialCiphertext")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("credentialCiphertext")+" <> ''")
		}
	}

	if f.StoreID != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreID)
	}

	if f.StoreIDNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDNe)
	}

	if f.StoreIDGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDGt)
	}

	if f.StoreIDLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDLt)
	}

	if f.StoreIDGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDGte)
	}

	if f.StoreIDLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDLte)
	}

	if f.StoreIDIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		values = append(values, f.StoreIDIn)
	}

	if f.StoreIDNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("storeId")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("storeId"))
		if *f.StoreIDNull {
			conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" IS NULL"+" OR "+aliasPrefix+SnakeString("storeId")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("storeId")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("storeId")+" <> ''")
		}
	}

	if f.IsDelete != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDelete)
	}

	if f.IsDeleteNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteNe)
	}

	if f.IsDeleteGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGt)
	}

	if f.IsDeleteLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLt)
	}

	if f.IsDeleteGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteGte)
	}

	if f.IsDeleteLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteLte)
	}

	if f.IsDeleteIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		values = append(values, f.IsDeleteIn)
	}

	if f.IsDeleteNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("isDelete")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("isDelete"))
		if *f.IsDeleteNull {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("isDelete")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("isDelete")+" <> ''")
		}
	}

	if f.Weight != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.Weight)
	}

	if f.WeightNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightNe)
	}

	if f.WeightGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGt)
	}

	if f.WeightLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLt)
	}

	if f.WeightGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightGte)
	}

	if f.WeightLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightLte)
	}

	if f.WeightIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		values = append(values, f.WeightIn)
	}

	if f.WeightNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("weight")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("weight"))
		if *f.WeightNull {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NULL"+" OR "+aliasPrefix+SnakeString("weight")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("weight")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("weight")+" <> ''")
		}
	}

	if f.State != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.State)
	}

	if f.StateNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateNe)
	}

	if f.StateGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGt)
	}

	if f.StateLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLt)
	}

	if f.StateGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateGte)
	}

	if f.StateLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateLte)
	}

	if f.StateIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("state")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		values = append(values, f.StateIn)
	}

	if f.StateNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("state")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("state"))
		if *f.StateNull {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NULL"+" OR "+aliasPrefix+SnakeString("state")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("state")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("state")+" <> ''")
		}
	}

	if f.DeletedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedBy)
	}

	if f.DeletedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByNe)
	}

	if f.DeletedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGt)
	}

	if f.DeletedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLt)
	}

	if f.DeletedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByGte)
	}

	if f.DeletedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByLte)
	}

	if f.DeletedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		values = append(values, f.DeletedByIn)
	}

	if f.DeletedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedBy"))
		if *f.DeletedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedBy")+" <> ''")
		}
	}

	if f.UpdatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedBy)
	}

	if f.UpdatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByNe)
	}

	if f.UpdatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGt)
	}

	if f.UpdatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLt)
	}

	if f.UpdatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByGte)
	}

	if f.UpdatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByLte)
	}

	if f.UpdatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		values = append(values, f.UpdatedByIn)
	}

	if f.UpdatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedBy"))
		if *f.UpdatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedBy")+" <> ''")
		}
	}

	if f.CreatedBy != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedBy)
	}

	if f.CreatedByNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByNe)
	}

	if f.CreatedByGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGt)
	}

	if f.CreatedByLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLt)
	}

	if f.CreatedByGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByGte)
	}

	if f.CreatedByLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByLte)
	}

	if f.CreatedByIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		values = append(values, f.CreatedByIn)
	}

	if f.CreatedByNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdBy")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdBy"))
		if *f.CreatedByNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdBy")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdBy")+" <> ''")
		}
	}

	if f.DeletedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAt)
	}

	if f.DeletedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtNe)
	}

	if f.DeletedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGt)
	}

	if f.DeletedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLt)
	}

	if f.DeletedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtGte)
	}

	if f.DeletedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtLte)
	}

	if f.DeletedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		values = append(values, f.DeletedAtIn)
	}

	if f.DeletedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("deletedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("deletedAt"))
		if *f.DeletedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("deletedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("deletedAt")+" <> ''")
		}
	}

	if f.UpdatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAt)
	}

	if f.UpdatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtNe)
	}

	if f.UpdatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLt)
	}

	if f.UpdatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtGte)
	}

	if f.UpdatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtLte)
	}

	if f.UpdatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		values = append(values, f.UpdatedAtIn)
	}

	if f.UpdatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("updatedAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("updatedAt"))
		if *f.UpdatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("updatedAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("updatedAt")+" <> ''")
		}
	}

	if f.CreatedAt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" = ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAt)
	}

	if f.CreatedAtNe != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" != ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtNe)
	}

	if f.CreatedAtGt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" > ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" < ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLt)
	}

	if f.CreatedAtGte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" >= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtGte)
	}

	if f.CreatedAtLte != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" <= ?")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtLte)
	}

	if f.CreatedAtIn != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IN (?)")
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		values = append(values, f.CreatedAtIn)
	}

	if f.CreatedAtNull != nil && IndexOf(whereConditions, aliasPrefix+SnakeString("createdAt")) == -1 {
		whereConditions = append(whereConditions, aliasPrefix+SnakeString("createdAt"))
		if *f.CreatedAtNull {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" =''")
		} else {
			conditions = append(conditions, aliasPrefix+SnakeString("createdAt")+" IS NOT NULL"+" OR "+aliasPrefix+SnakeString("createdAt")+" <> ''")
		}
	}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *StorePaymentConfigFilterType) AndWith(f2 ...*StorePaymentConfigFilterType) *StorePaymentConfigFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &StorePaymentConfigFilterType{
		And: append(_f2, f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *StorePaymentConfigFilterType) OrWith(f2 ...*StorePaymentConfigFilterType) *StorePaymentConfigFilterType {
	_f2 := CompactNils(f2)
	if len(_f2) == 0 {
		return f
	}
	return &StorePaymentConfigFilterType{
		Or: append(_f2, f),
	}
}
