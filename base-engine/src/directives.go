/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"base-engine/utils"
)

func configureDirectives(configuration *gen.Config, db *gen.DB) {
	configuration.Directives.HasPermission = hasPermissionDirective
	configuration.Directives.HasRole = hasRoleDirective
	configuration.Directives.Validator = validatorDirective(db)
}

func hasPermissionDirective(ctx context.Context, _ any, next graphql.Resolver, action string) (any, error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	workspaceAction, err := authorization.WorkspaceAction(principal, action)
	if err != nil {
		return nil, err
	}
	if !principal.Has(workspaceAction) {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return next(ctx)
}

func hasRoleDirective(ctx context.Context, _ any, next graphql.Resolver, role gen.Role) (any, error) {
	if _, err := auth.RequirePrincipal(ctx); err != nil {
		return nil, err
	}
	if role != gen.RoleAll {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return next(ctx)
}

func validatorDirective(db *gen.DB) func(context.Context, any, graphql.Resolver, *string, *string, *string, *int, *int, *int, *int, *string, *string) (any, error) {
	return func(ctx context.Context, obj any, next graphql.Resolver, required, immutable, typeArg *string, minLength, maxLength, minValue, maxValue *int, unique, uniqueScope *string) (any, error) {
		value, err := next(ctx)
		if err != nil {
			return nil, err
		}
		fieldName := utils.GetFieldName(obj, value)
		ctx = uniqueValidationContext(ctx, db, unique)
		if err := utils.ValidateField(ctx, obj, fieldName, value, required, immutable, typeArg, minLength, maxLength, minValue, maxValue, unique, uniqueScope); err != nil {
			return nil, auth.NewError(auth.CodeValidationFailed)
		}
		return value, nil
	}
}

func uniqueValidationContext(ctx context.Context, db *gen.DB, unique *string) context.Context {
	if unique == nil || *unique != "true" {
		return ctx
	}
	field := graphql.GetFieldContext(ctx)
	if field == nil {
		return ctx
	}
	name := mutationEntityName(field.Field.Name)
	if name == "" {
		return ctx
	}
	for _, model := range gen.TableMap {
		if fmt.Sprintf("%T", model) == "gen."+name {
			ctx = context.WithValue(ctx, "db", db.Query())
			return context.WithValue(ctx, "modelStruct", model)
		}
	}
	return ctx
}

func mutationEntityName(fieldName string) string {
	if strings.HasPrefix(fieldName, "create") {
		return strings.TrimPrefix(fieldName, "create")
	}
	if strings.HasPrefix(fieldName, "update") {
		return strings.TrimPrefix(fieldName, "update")
	}
	return ""
}
