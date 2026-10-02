/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import "context"

type principalContextKey struct{}
type authenticationErrorContextKey struct{}
type principalResolverContextKey struct{}

// PrincipalResolverFunc 在长连接上的每次受保护操作中重新解析权威身份。
type PrincipalResolverFunc func(context.Context) (*WorkspacePrincipal, error)

// WorkspacePrincipal 是每次请求从数据库重新解析的权威身份。
type WorkspacePrincipal struct {
	AccountID      string
	SessionID      string
	RequestID      string
	WorkspaceType  WorkspaceType
	OrganizationID *string
	MembershipID   *string
	Permissions    map[string]struct{}
	StoreIDs       map[string]struct{}
	AllStores      bool
}

// Has 判断权威权限集合是否包含指定动作。
func (p *WorkspacePrincipal) Has(action string) bool {
	if p == nil {
		return false
	}
	_, ok := p.Permissions[action]
	return ok
}

// HasStore 判断门店是否位于解析后的活跃门店范围。
func (p *WorkspacePrincipal) HasStore(storeID string) bool {
	if p == nil {
		return false
	}
	_, ok := p.StoreIDs[storeID]
	return ok
}

// WithPrincipal 把权威身份写入请求上下文。
func WithPrincipal(ctx context.Context, principal *WorkspacePrincipal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

// PrincipalFromContext 读取请求级权威身份。
func PrincipalFromContext(ctx context.Context) (*WorkspacePrincipal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(*WorkspacePrincipal)
	return principal, ok && principal != nil
}

// WithAuthenticationError 保存权威 Session 解析出的结构化终态错误。
func WithAuthenticationError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, authenticationErrorContextKey{}, err)
}

// WithPrincipalResolver 为长连接保存按操作执行的权威身份解析器。
func WithPrincipalResolver(ctx context.Context, resolve PrincipalResolverFunc) context.Context {
	return context.WithValue(ctx, principalResolverContextKey{}, resolve)
}

// RequirePrincipal 返回权威身份，并保留可安全暴露的会话终态错误码。
func RequirePrincipal(ctx context.Context) (*WorkspacePrincipal, error) {
	if resolve, ok := ctx.Value(principalResolverContextKey{}).(PrincipalResolverFunc); ok && resolve != nil {
		principal, err := resolve(ctx)
		if err != nil {
			return nil, safeAuthenticationError(err)
		}
		if principal == nil {
			return nil, NewError(CodeAuthRequired)
		}
		return principal, nil
	}
	if principal, ok := PrincipalFromContext(ctx); ok {
		return principal, nil
	}
	if err, ok := ctx.Value(authenticationErrorContextKey{}).(error); ok {
		return nil, safeAuthenticationError(err)
	}
	return nil, NewError(CodeAuthRequired)
}

func safeAuthenticationError(err error) error {
	if code := ErrorCode(err); code != "" {
		return NewError(code)
	}
	return NewError(CodeAuthRequired)
}
