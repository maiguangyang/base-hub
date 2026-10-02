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
