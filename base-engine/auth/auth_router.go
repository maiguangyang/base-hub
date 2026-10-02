/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/iancoleman/strcase"
)

// GetMethodName 返回当前 GraphQL 顶层字段名。
func GetMethodName(ctx context.Context) (*string, error) {
	resolver := graphql.GetFieldContext(ctx)
	if resolver == nil {
		return nil, fmt.Errorf("GRAPHQL_FIELD_REQUIRED")
	}
	path := strings.Split(resolver.Path().String(), ".")
	if len(path) == 0 || path[0] == "" {
		return nil, fmt.Errorf("GRAPHQL_FIELD_REQUIRED")
	}
	name := strcase.ToCamel(path[0])
	return &name, nil
}

// CheckAuthorization 要求请求已建立数据库权威身份。
func CheckAuthorization(ctx context.Context, _ string) error {
	_, err := requirePrincipal(ctx)
	return err
}

// UserTokenVerify 保留 generated 指令接口并要求加盟工作台。
func UserTokenVerify(ctx context.Context, _ string) error {
	principal, err := requirePrincipal(ctx)
	if err != nil {
		return err
	}
	if principal.WorkspaceType != WorkspaceTypeFranchise {
		return NewError(CodeWorkspaceForbidden)
	}
	return nil
}

// AdminTokenVerify 保留 generated 指令接口并要求总部工作台。
func AdminTokenVerify(ctx context.Context, _ string) error {
	principal, err := requirePrincipal(ctx)
	if err != nil {
		return err
	}
	if principal.WorkspaceType != WorkspaceTypeHeadquarters {
		return NewError(CodeWorkspaceForbidden)
	}
	return nil
}

func requirePrincipal(ctx context.Context) (*WorkspacePrincipal, error) {
	return RequirePrincipal(ctx)
}
