package gen

import (
	"context"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/gorm"
)

type AccountQueryFilter struct {
	Query *string
}

func (qf *AccountQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("accounts", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *AccountQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["phone"]; ok {

		column := alias + "." + SnakeString("phone")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["displayName"]; ok {

		column := alias + "." + SnakeString("displayName")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["email"]; ok {

		column := alias + "." + SnakeString("email")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["credentialVersion"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("credentialVersion")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["memberships"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_memberships"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"."+"account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorMembershipQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["initializedOrganizations"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_initializedOrganizations"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+"."+"initial_account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["openingRecords"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_openingRecords"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"initial_account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := FranchiseOpeningRecordQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["recordedOpeningRecords"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_recordedOpeningRecords"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"recorded_by_account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := FranchiseOpeningRecordQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["sessions"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sessions"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+"."+"account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SessionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["reviewedStores"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_reviewedStores"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"."+"reviewed_by_account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["sentMembershipInvitations"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sentMembershipInvitations"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("membership_invitations", ctx)+" "+_alias+" ON "+_alias+"."+"invited_by_account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := MembershipInvitationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["auditLogs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_auditLogs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"actor_account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AuditLogQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["createdStocktakes"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_createdStocktakes"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+"."+"initiated_by_account_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["postedStocktakes"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_postedStocktakes"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+"."+"posted_by_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type OrganizationQueryFilter struct {
	Query *string
}

func (qf *OrganizationQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("organizations", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *OrganizationQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["code"]; ok {

		column := alias + "." + SnakeString("code")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["suspensionReasonCode"]; ok {

		column := alias + "." + SnakeString("suspensionReasonCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["memberships"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_memberships"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorMembershipQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["initialAccount"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_initialAccount"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initial_account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["openingRecords"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_openingRecords"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("franchise_opening_records", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := FranchiseOpeningRecordQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["stores"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stores"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["roles"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_roles"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorRoleQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["sessions"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sessions"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SessionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["auditLogs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_auditLogs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AuditLogQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["paymentConfigs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_paymentConfigs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("franchise_payment_configs", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := FranchisePaymentConfigQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["customerMembers"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_customerMembers"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerMemberQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["customerBenefitPolicies"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_customerBenefitPolicies"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_benefit_policies", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerBenefitPolicyQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["customerDailyPointGrantBudgets"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_customerDailyPointGrantBudgets"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_daily_point_grant_budgets", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerDailyPointGrantBudgetQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["customerPointEntries"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_customerPointEntries"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+"."+"source_organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerPointEntryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["customerCouponTemplates"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_customerCouponTemplates"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["productCategories"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_productCategories"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductCategoryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["productBrands"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_productBrands"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_brands", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductBrandQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["specificationDefinitions"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_specificationDefinitions"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("specification_definitions", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SpecificationDefinitionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["productPackageTemplates"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_productPackageTemplates"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["products"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_products"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"organization_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type OperatorMembershipQueryFilter struct {
	Query *string
}

func (qf *OperatorMembershipQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("operator_memberships", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *OperatorMembershipQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["account"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_account"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["roles"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_roles"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operatorMembership_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"member_id"+" LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"role_id"+" = "+_alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorRoleQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["stores"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stores"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operatorMembership_stores", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"member_id"+" LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"store_id"+" = "+_alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["invitations"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_invitations"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("membership_invitations", ctx)+" "+_alias+" ON "+_alias+"."+"membership_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := MembershipInvitationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type PermissionQueryFilter struct {
	Query *string
}

func (qf *PermissionQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("permissions", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *PermissionQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["action"]; ok {

		column := alias + "." + SnakeString("action")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["module"]; ok {

		column := alias + "." + SnakeString("module")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["roles"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_roles"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("permission_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"permission_id"+" LEFT JOIN "+TableName("operator_roles", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"role_id"+" = "+_alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorRoleQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type OperatorRoleQueryFilter struct {
	Query *string
}

func (qf *OperatorRoleQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("operator_roles", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *OperatorRoleQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["members"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_members"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operatorMembership_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"role_id"+" LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"member_id"+" = "+_alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorMembershipQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["permissions"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_permissions"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("permission_roles", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"role_id"+" LEFT JOIN "+TableName("permissions", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"permission_id"+" = "+_alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := PermissionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StoreQueryFilter struct {
	Query *string
}

func (qf *StoreQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("stores", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StoreQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["code"]; ok {

		column := alias + "." + SnakeString("code")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["rejectionReason"]; ok {

		column := alias + "." + SnakeString("rejectionReason")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["contactPhone"]; ok {

		column := alias + "." + SnakeString("contactPhone")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["managerName"]; ok {

		column := alias + "." + SnakeString("managerName")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["managerPhone"]; ok {

		column := alias + "." + SnakeString("managerPhone")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["province"]; ok {

		column := alias + "." + SnakeString("province")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["city"]; ok {

		column := alias + "." + SnakeString("city")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["district"]; ok {

		column := alias + "." + SnakeString("district")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["address"]; ok {

		column := alias + "." + SnakeString("address")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["businessHours"]; ok {

		column := alias + "." + SnakeString("businessHours")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["storeArea"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("storeArea")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["tableCount"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("tableCount")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["receiptFooter"]; ok {

		column := alias + "." + SnakeString("receiptFooter")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["businessLicenseImageUrl"]; ok {

		column := alias + "." + SnakeString("businessLicenseImageUrl")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["otherDocumentImageUrl"]; ok {

		column := alias + "." + SnakeString("otherDocumentImageUrl")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["members"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_members"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operatorMembership_stores", ctx)+" "+_alias+"_jointable"+" ON "+alias+".id = "+_alias+"_jointable"+"."+"store_id"+" LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+"_jointable"+"."+"member_id"+" = "+_alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorMembershipQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["reviewedByAccount"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_reviewedByAccount"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"reviewed_by_account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["auditLogs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_auditLogs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AuditLogQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["paymentConfigs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_paymentConfigs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_payment_configs", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePaymentConfigQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["productListings"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_productListings"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreListingQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["stockMovements"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stockMovements"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStockMovementQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["stocktakes"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stocktakes"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["promotions"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_promotions"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_promotions", ctx)+" "+_alias+" ON "+_alias+"."+"store_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePromotionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["couponTemplates"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_couponTemplates"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+"."+"applicable_store_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type SessionQueryFilter struct {
	Query *string
}

func (qf *SessionQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("sessions", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *SessionQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["credentialVersion"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("credentialVersion")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["revocationCode"]; ok {

		column := alias + "." + SnakeString("revocationCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["account"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_account"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["auditLogs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_auditLogs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("audit_logs", ctx)+" "+_alias+" ON "+_alias+"."+"session_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AuditLogQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type MembershipInvitationQueryFilter struct {
	Query *string
}

func (qf *MembershipInvitationQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("membership_invitations", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *MembershipInvitationQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["membership"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_membership"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("operator_memberships", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"membership_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OperatorMembershipQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["invitedByAccount"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_invitedByAccount"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"invited_by_account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type AuditLogQueryFilter struct {
	Query *string
}

func (qf *AuditLogQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("audit_logs", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *AuditLogQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["action"]; ok {

		column := alias + "." + SnakeString("action")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["resourceType"]; ok {

		column := alias + "." + SnakeString("resourceType")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["resourceId"]; ok {

		column := alias + "." + SnakeString("resourceId")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["resultCode"]; ok {

		column := alias + "." + SnakeString("resultCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["metadataJson"]; ok {

		column := alias + "." + SnakeString("metadataJson")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["actorAccount"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_actorAccount"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"actor_account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["session"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_session"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("sessions", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"session_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SessionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["store"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_store"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type FranchiseOpeningRecordQueryFilter struct {
	Query *string
}

func (qf *FranchiseOpeningRecordQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("franchise_opening_records", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *FranchiseOpeningRecordQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["recordNumber"]; ok {

		column := alias + "." + SnakeString("recordNumber")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["initialAccount"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_initialAccount"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initial_account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["recordedByAccount"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_recordedByAccount"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"recorded_by_account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type GlobalPaymentConfigQueryFilter struct {
	Query *string
}

func (qf *GlobalPaymentConfigQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("global_payment_configs", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *GlobalPaymentConfigQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["channel"]; ok {

		column := alias + "." + SnakeString("channel")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["merchantId"]; ok {

		column := alias + "." + SnakeString("merchantId")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["environment"]; ok {

		column := alias + "." + SnakeString("environment")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["ratePpm"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("ratePpm")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["configState"]; ok {

		column := alias + "." + SnakeString("configState")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["version"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("version")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["keyId"]; ok {

		column := alias + "." + SnakeString("keyId")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["credentialCiphertext"]; ok {

		column := alias + "." + SnakeString("credentialCiphertext")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	return nil
}

type FranchisePaymentConfigQueryFilter struct {
	Query *string
}

func (qf *FranchisePaymentConfigQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("franchise_payment_configs", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *FranchisePaymentConfigQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["channel"]; ok {

		column := alias + "." + SnakeString("channel")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["merchantId"]; ok {

		column := alias + "." + SnakeString("merchantId")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["environment"]; ok {

		column := alias + "." + SnakeString("environment")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["ratePpm"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("ratePpm")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["configState"]; ok {

		column := alias + "." + SnakeString("configState")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["version"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("version")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["keyId"]; ok {

		column := alias + "." + SnakeString("keyId")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["credentialCiphertext"]; ok {

		column := alias + "." + SnakeString("credentialCiphertext")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StorePaymentConfigQueryFilter struct {
	Query *string
}

func (qf *StorePaymentConfigQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_payment_configs", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StorePaymentConfigQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["channel"]; ok {

		column := alias + "." + SnakeString("channel")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["merchantId"]; ok {

		column := alias + "." + SnakeString("merchantId")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["environment"]; ok {

		column := alias + "." + SnakeString("environment")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["ratePpm"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("ratePpm")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["configState"]; ok {

		column := alias + "." + SnakeString("configState")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["version"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("version")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["keyId"]; ok {

		column := alias + "." + SnakeString("keyId")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["credentialCiphertext"]; ok {

		column := alias + "." + SnakeString("credentialCiphertext")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["store"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_store"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type CustomerMemberQueryFilter struct {
	Query *string
}

func (qf *CustomerMemberQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("customer_members", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *CustomerMemberQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["memberNumber"]; ok {

		column := alias + "." + SnakeString("memberNumber")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["requestKey"]; ok {

		column := alias + "." + SnakeString("requestKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["phone"]; ok {

		column := alias + "." + SnakeString("phone")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["pointsBalance"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("pointsBalance")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["noticeVersion"]; ok {

		column := alias + "." + SnakeString("noticeVersion")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["processingBasisCode"]; ok {

		column := alias + "." + SnakeString("processingBasisCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["evidenceReference"]; ok {

		column := alias + "." + SnakeString("evidenceReference")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["pointEntries"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_pointEntries"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+"."+"member_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerPointEntryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["couponGrants"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_couponGrants"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_grants", ctx)+" "+_alias+" ON "+_alias+"."+"member_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponGrantQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["couponDistributionJobs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_couponDistributionJobs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_distribution_jobs", ctx)+" "+_alias+" ON "+_alias+"."+"member_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponDistributionJobQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type CustomerBenefitPolicyQueryFilter struct {
	Query *string
}

func (qf *CustomerBenefitPolicyQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("customer_benefit_policies", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *CustomerBenefitPolicyQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["version"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("version")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["discountBasisPoints"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("discountBasisPoints")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["earnAmountFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("earnAmountFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["earnPoints"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("earnPoints")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["redeemPoints"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("redeemPoints")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["redeemAmountFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("redeemAmountFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["maxRedemptionBasisPoints"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("maxRedemptionBasisPoints")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["maxRedemptionPoints"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("maxRedemptionPoints")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["manualGrantMaxSingle"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("manualGrantMaxSingle")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["manualGrantMaxDaily"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("manualGrantMaxDaily")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type CustomerDailyPointGrantBudgetQueryFilter struct {
	Query *string
}

func (qf *CustomerDailyPointGrantBudgetQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("customer_daily_point_grant_budgets", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *CustomerDailyPointGrantBudgetQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["businessDate"]; ok {

		column := alias + "." + SnakeString("businessDate")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["usedPoints"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("usedPoints")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type CustomerPointEntryQueryFilter struct {
	Query *string
}

func (qf *CustomerPointEntryQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("customer_point_entries", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *CustomerPointEntryQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["delta"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("delta")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["note"]; ok {

		column := alias + "." + SnakeString("note")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["requestKey"]; ok {

		column := alias + "." + SnakeString("requestKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["member"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_member"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"member_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerMemberQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["sourceOrganization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sourceOrganization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"source_organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["reverses"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_reverses"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"reverses_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerPointEntryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["reversedBy"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_reversedBy"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_point_entries", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"reversed_by_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerPointEntryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type CustomerCouponTemplateQueryFilter struct {
	Query *string
}

func (qf *CustomerCouponTemplateQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("customer_coupon_templates", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *CustomerCouponTemplateQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["code"]; ok {

		column := alias + "." + SnakeString("code")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["requestKey"]; ok {

		column := alias + "." + SnakeString("requestKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["title"]; ok {

		column := alias + "." + SnakeString("title")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["amountFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("amountFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["minSpendFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("minSpendFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["daysAfterActivation"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("daysAfterActivation")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["effectiveAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("effectiveAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["distributionEndsAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("distributionEndsAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["perMemberLimit"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("perMemberLimit")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["totalIssueLimit"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("totalIssueLimit")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["issuedCount"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("issuedCount")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["applicableStore"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_applicableStore"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"applicable_store_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["grants"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_grants"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_grants", ctx)+" "+_alias+" ON "+_alias+"."+"template_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponGrantQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["distributionJobs"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_distributionJobs"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_distribution_jobs", ctx)+" "+_alias+" ON "+_alias+"."+"template_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponDistributionJobQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductCategoryQueryFilter struct {
	Query *string
}

func (qf *ProductCategoryQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("product_categories", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductCategoryQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["sortOrder"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("sortOrder")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["parent"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_parent"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"parent_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductCategoryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["children"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_children"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+"."+"parent_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductCategoryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["products"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_products"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"category_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductBrandQueryFilter struct {
	Query *string
}

func (qf *ProductBrandQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("product_brands", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductBrandQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["products"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_products"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"brand_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductQueryFilter struct {
	Query *string
}

func (qf *ProductQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("products", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["description"]; ok {

		column := alias + "." + SnakeString("description")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["imageUrl"]; ok {

		column := alias + "." + SnakeString("imageUrl")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["brand"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_brand"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_brands", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"brand_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductBrandQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["category"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_category"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_categories", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"category_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductCategoryQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["defaultPackageTemplate"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_defaultPackageTemplate"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"default_package_template_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["specificationChoices"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_specificationChoices"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_specification_choices", ctx)+" "+_alias+" ON "+_alias+"."+"product_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSpecificationChoiceQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["skus"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_skus"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+"."+"product_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSkuQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductSkuQueryFilter struct {
	Query *string
}

func (qf *ProductSkuQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("product_skus", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductSkuQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["publishedPackageSetVersion"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("publishedPackageSetVersion")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["ingredients"]; ok {

		column := alias + "." + SnakeString("ingredients")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["allergens"]; ok {

		column := alias + "." + SnakeString("allergens")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["storageInstructions"]; ok {

		column := alias + "." + SnakeString("storageInstructions")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["shelfLifeDays"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("shelfLifeDays")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["product"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_product"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"product_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["specificationValues"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_specificationValues"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_sku_specification_values", ctx)+" "+_alias+" ON "+_alias+"."+"sku_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSkuSpecificationValueQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["packages"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_packages"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+"."+"sku_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["listings"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_listings"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+"."+"sku_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreListingQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductPackageQueryFilter struct {
	Query *string
}

func (qf *ProductPackageQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("product_packages", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductPackageQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["barcode"]; ok {

		column := alias + "." + SnakeString("barcode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["packageSetVersion"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("packageSetVersion")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["containsQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("containsQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["suggestedPriceFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("suggestedPriceFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["sku"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sku"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"sku_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSkuQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["template"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_template"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"template_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["containsPackage"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_containsPackage"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"contains_package_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["containedByPackages"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_containedByPackages"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+"."+"contains_package_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["offers"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_offers"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+"."+"package_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePackageOfferQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["balances"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_balances"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stock_balances", ctx)+" "+_alias+" ON "+_alias+"."+"package_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStockBalanceQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["movementSources"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_movementSources"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"source_package_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStockMovementQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["movementTargets"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_movementTargets"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"target_package_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStockMovementQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["stocktakeLines"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stocktakeLines"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+"."+"package_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeLineQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type SpecificationDefinitionQueryFilter struct {
	Query *string
}

func (qf *SpecificationDefinitionQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("specification_definitions", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *SpecificationDefinitionQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["values"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_values"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("specification_values", ctx)+" "+_alias+" ON "+_alias+"."+"specification_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SpecificationValueQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type SpecificationValueQueryFilter struct {
	Query *string
}

func (qf *SpecificationValueQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("specification_values", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *SpecificationValueQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["specification"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_specification"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("specification_definitions", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"specification_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SpecificationDefinitionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["productChoices"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_productChoices"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_specification_choices", ctx)+" "+_alias+" ON "+_alias+"."+"value_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSpecificationChoiceQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["skuValues"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_skuValues"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_sku_specification_values", ctx)+" "+_alias+" ON "+_alias+"."+"value_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSkuSpecificationValueQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductSpecificationChoiceQueryFilter struct {
	Query *string
}

func (qf *ProductSpecificationChoiceQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("product_specification_choices", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductSpecificationChoiceQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["product"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_product"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"product_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["value"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_value"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("specification_values", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"value_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SpecificationValueQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductSkuSpecificationValueQueryFilter struct {
	Query *string
}

func (qf *ProductSkuSpecificationValueQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("product_sku_specification_values", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductSkuSpecificationValueQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["sku"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sku"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"sku_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSkuQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["value"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_value"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("specification_values", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"value_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := SpecificationValueQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProductPackageTemplateQueryFilter struct {
	Query *string
}

func (qf *ProductPackageTemplateQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("product_package_templates", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *ProductPackageTemplateQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["name"]; ok {

		column := alias + "." + SnakeString("name")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["containsQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("containsQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["organization"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_organization"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("organizations", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"organization_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := OrganizationQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["defaultProducts"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_defaultProducts"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("products", ctx)+" "+_alias+" ON "+_alias+"."+"default_package_template_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["containsPackage"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_containsPackage"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"contains_package_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["containedByPackages"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_containedByPackages"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_package_templates", ctx)+" "+_alias+" ON "+_alias+"."+"contains_package_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["createdPackages"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_createdPackages"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+"."+"template_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StoreListingQueryFilter struct {
	Query *string
}

func (qf *StoreListingQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_listings", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StoreListingQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["store"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_store"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["sku"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sku"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_skus", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"sku_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductSkuQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["offers"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_offers"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+"."+"listing_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePackageOfferQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["batches"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_batches"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+"."+"listing_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreInventoryBatchQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StorePackageOfferQueryFilter struct {
	Query *string
}

func (qf *StorePackageOfferQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_package_offers", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StorePackageOfferQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["priceFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("priceFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["listing"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_listing"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"listing_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreListingQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["package"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_package"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"package_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["priceRevisions"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_priceRevisions"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_price_revisions", ctx)+" "+_alias+" ON "+_alias+"."+"offer_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePriceRevisionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["promotionTargets"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_promotionTargets"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_promotion_targets", ctx)+" "+_alias+" ON "+_alias+"."+"offer_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePromotionTargetQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StorePriceRevisionQueryFilter struct {
	Query *string
}

func (qf *StorePriceRevisionQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_price_revisions", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StorePriceRevisionQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["previousPriceFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("previousPriceFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["priceFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("priceFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["reasonCode"]; ok {

		column := alias + "." + SnakeString("reasonCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["offer"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_offer"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"offer_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePackageOfferQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StoreInventoryBatchQueryFilter struct {
	Query *string
}

func (qf *StoreInventoryBatchQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_inventory_batches", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StoreInventoryBatchQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["batchNumber"]; ok {

		column := alias + "." + SnakeString("batchNumber")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["sourceReference"]; ok {

		column := alias + "." + SnakeString("sourceReference")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["listing"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_listing"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_listings", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"listing_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreListingQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["balances"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_balances"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stock_balances", ctx)+" "+_alias+" ON "+_alias+"."+"batch_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStockBalanceQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["movements"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_movements"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"batch_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStockMovementQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["stocktakeLines"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stocktakeLines"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+"."+"batch_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeLineQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StoreStockBalanceQueryFilter struct {
	Query *string
}

func (qf *StoreStockBalanceQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_stock_balances", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StoreStockBalanceQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["quantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("quantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["version"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("version")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["batch"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_batch"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"batch_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreInventoryBatchQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["package"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_package"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"package_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StoreStocktakeQueryFilter struct {
	Query *string
}

func (qf *StoreStocktakeQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_stocktakes", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StoreStocktakeQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["requestKey"]; ok {

		column := alias + "." + SnakeString("requestKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["scopeDigest"]; ok {

		column := alias + "." + SnakeString("scopeDigest")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["store"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_store"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["initiatedByAccount"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_initiatedByAccount"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"initiated_by_account_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["postedBy"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_postedBy"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("accounts", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"posted_by_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := AccountQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["lines"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_lines"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+"."+"stocktake_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeLineQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StoreStocktakeLineQueryFilter struct {
	Query *string
}

func (qf *StoreStocktakeLineQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_stocktake_lines", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StoreStocktakeLineQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["snapshotQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("snapshotQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["snapshotVersion"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("snapshotVersion")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["packageSetVersion"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("packageSetVersion")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["batchNumberSnapshot"]; ok {

		column := alias + "." + SnakeString("batchNumberSnapshot")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["packageNameSnapshot"]; ok {

		column := alias + "." + SnakeString("packageNameSnapshot")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["countedQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("countedQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["reasonCode"]; ok {

		column := alias + "." + SnakeString("reasonCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["reasonNote"]; ok {

		column := alias + "." + SnakeString("reasonNote")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["stocktake"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stocktake"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktakes", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"stocktake_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["batch"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_batch"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"batch_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreInventoryBatchQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["package"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_package"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"package_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["movements"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_movements"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stock_movements", ctx)+" "+_alias+" ON "+_alias+"."+"stocktake_line_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStockMovementQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StoreStockMovementQueryFilter struct {
	Query *string
}

func (qf *StoreStockMovementQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_stock_movements", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StoreStockMovementQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["requestKey"]; ok {

		column := alias + "." + SnakeString("requestKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["sourceQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("sourceQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["targetQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("targetQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["factorSnapshot"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("factorSnapshot")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["packageSetVersion"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("packageSetVersion")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["reasonCode"]; ok {

		column := alias + "." + SnakeString("reasonCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["store"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_store"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["batch"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_batch"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_inventory_batches", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"batch_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreInventoryBatchQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["sourcePackage"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_sourcePackage"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"source_package_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["targetPackage"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_targetPackage"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("product_packages", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"target_package_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := ProductPackageQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["stocktakeLine"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_stocktakeLine"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_stocktake_lines", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"stocktake_line_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreStocktakeLineQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StorePromotionQueryFilter struct {
	Query *string
}

func (qf *StorePromotionQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_promotions", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StorePromotionQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["ruleKey"]; ok {

		column := alias + "." + SnakeString("ruleKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["timeZone"]; ok {

		column := alias + "." + SnakeString("timeZone")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["version"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("version")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["thresholdFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("thresholdFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["thresholdQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("thresholdQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["discountFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("discountFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["discountBasisPoints"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("discountBasisPoints")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["fixedPriceFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("fixedPriceFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["store"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_store"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("stores", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"store_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StoreQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["targets"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_targets"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_promotion_targets", ctx)+" "+_alias+" ON "+_alias+"."+"promotion_id"+" = "+alias+".id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePromotionTargetQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type StorePromotionTargetQueryFilter struct {
	Query *string
}

func (qf *StorePromotionTargetQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("store_promotion_targets", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *StorePromotionTargetQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["requiredQuantity"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("requiredQuantity")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["promotion"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_promotion"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_promotions", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"promotion_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePromotionQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["offer"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_offer"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("store_package_offers", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"offer_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := StorePackageOfferQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type CustomerCouponGrantQueryFilter struct {
	Query *string
}

func (qf *CustomerCouponGrantQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("customer_coupon_grants", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *CustomerCouponGrantQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["amountFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("amountFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["minSpendFen"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("minSpendFen")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["daysAfterActivation"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("daysAfterActivation")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["issuedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("issuedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["activatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("activatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["expiresAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("expiresAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["revokedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("revokedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["requestKey"]; ok {

		column := alias + "." + SnakeString("requestKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["issuerRequestDigest"]; ok {

		column := alias + "." + SnakeString("issuerRequestDigest")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["member"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_member"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"member_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerMemberQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["template"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_template"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"template_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type CustomerCouponDistributionJobQueryFilter struct {
	Query *string
}

func (qf *CustomerCouponDistributionJobQueryFilter) Apply(ctx context.Context, db *gorm.DB, selectionSet *ast.SelectionSet, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}

	fields := []*ast.Field{}
	if selectionSet != nil {
		for _, s := range *selectionSet {
			if f, ok := s.(*ast.Field); ok {
				fields = append(fields, f)
			}
		}
	} else {
		return fmt.Errorf("Cannot query with 'q' attribute without items field.")
	}

	if *qf.Query != "" {
		queryParts := strings.Split(*qf.Query, " ")
		for _, part := range queryParts {
			ors := []string{}
			if err := qf.applyQueryWithFields(db, fields, part, TableName("customer_coupon_distribution_jobs", ctx), &ors, values, joins); err != nil {
				return err
			}
			*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
		}
	}
	return nil
}

func (qf *CustomerCouponDistributionJobQueryFilter) applyQueryWithFields(db *gorm.DB, fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}

	fieldsMap := map[string][]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = append(fieldsMap[f.Name], f)
	}

	if _, ok := fieldsMap["requestKey"]; ok {

		column := alias + "." + SnakeString("requestKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["availableAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("availableAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["leaseExpiresAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("leaseExpiresAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["leaseToken"]; ok {

		column := alias + "." + SnakeString("leaseToken")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["cursorCreatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("cursorCreatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["cursorKey"]; ok {

		column := alias + "." + SnakeString("cursorKey")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["attempts"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("attempts")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["lastErrorCode"]; ok {

		column := alias + "." + SnakeString("lastErrorCode")

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["isDelete"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("isDelete")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["weight"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("weight")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["state"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("state")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["deletedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("deletedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["updatedAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("updatedAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	if _, ok := fieldsMap["createdAt"]; ok {

		cast := "TEXT"
		if db.Name() == "mysql" {
			cast = "CHAR"
		}
		column := fmt.Sprintf("CAST(%s"+SnakeString("createdAt")+" AS %s)", alias+".", cast)

		*ors = append(*ors, fmt.Sprintf("%[1]s LIKE ? OR %[1]s LIKE ?", column))
		*values = append(*values, query+"%", "%"+query+"%")
	}

	// if fs, ok := fieldsMap["template"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_template"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_coupon_templates", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"template_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerCouponTemplateQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if fs, ok := fieldsMap["member"]; ok {
	// 	_fields := []*ast.Field{}
	// 	_alias := alias + "_member"
	// 	*joins = append(*joins,"LEFT JOIN "+TableName("customer_members", ctx)+" "+_alias+" ON "+_alias+".id = "+alias+"."+"member_id")

	// 	for _, f := range fs {
	// 		for _, s := range f.SelectionSet {
	// 			if f, ok := s.(*ast.Field); ok {
	// 				_fields = append(_fields, f)
	// 			}
	// 		}
	// 	}
	// 	q := CustomerMemberQueryFilter{qf.Query}
	// 	err := q.applyQueryWithFields(db, _fields, query, _alias, ors, values, joins)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
