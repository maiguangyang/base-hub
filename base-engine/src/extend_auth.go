/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

// Login 登录并创建 DISCOVERY 权威会话。
func (r *MutationResolver) Login(ctx context.Context, input gen.LoginInput) (*gen.LoginPayload, error) {
	result, err := r.Services.Authentication.Login(ctx, input.Phone, input.Password, RemoteIPFromContext(ctx), time.Now())
	if err != nil {
		return nil, err
	}
	if writer := ResponseWriterFromContext(ctx); writer != nil {
		auth.SetSessionCookie(writer, result.Token, r.Services.SecurityConfig)
	}
	principal := &auth.WorkspacePrincipal{
		AccountID: result.Account.ID, SessionID: result.Session.ID,
		WorkspaceType: auth.WorkspaceTypeDiscovery,
		Permissions:   make(map[string]struct{}), StoreIDs: make(map[string]struct{}),
	}
	viewer, err := buildViewer(ctx, r.Services, principal)
	if err != nil {
		return nil, err
	}
	return &gen.LoginPayload{Viewer: viewer, RequiresPasswordChange: result.Account.MustChangePassword}, nil
}

// Logout 撤销当前会话并清除 Cookie。
func (r *MutationResolver) Logout(ctx context.Context) (bool, error) {
	principal, principalErr := auth.RequirePrincipal(ctx)
	if principalErr != nil {
		return false, principalErr
	}
	now := time.Now()
	err := r.Services.DB.Transaction(func(tx *gorm.DB) error {
		return r.Services.Sessions.Revoke(ctx, tx, principal.SessionID, "SESSION_REVOKED", now)
	})
	if err != nil {
		return false, err
	}
	if r.Services.Publisher != nil {
		r.Services.Publisher.PublishSession(principal.SessionID, &gen.SessionEvent{Code: gen.SessionEventCodeSessionRevoked, SessionID: principal.SessionID, OrganizationID: principal.OrganizationID, OccurredAt: now})
	}
	if writer := ResponseWriterFromContext(ctx); writer != nil {
		auth.ClearSessionCookie(writer, r.Services.SecurityConfig)
	}
	return true, nil
}

// ChangeTemporaryPassword 更换临时密码并轮换会话。
func (r *MutationResolver) ChangeTemporaryPassword(ctx context.Context, input gen.ChangePasswordInput) (*gen.Viewer, error) {
	principal, principalErr := auth.RequirePrincipal(ctx)
	if principalErr != nil {
		return nil, principalErr
	}
	result, err := r.Services.Authentication.ChangePassword(ctx, principal, input.CurrentPassword, input.NewPassword, time.Now())
	if err != nil {
		return nil, err
	}
	if writer := ResponseWriterFromContext(ctx); writer != nil {
		auth.SetSessionCookie(writer, result.Token, r.Services.SecurityConfig)
	}
	discovery := &auth.WorkspacePrincipal{AccountID: result.Account.ID, SessionID: result.Session.ID, WorkspaceType: auth.WorkspaceTypeDiscovery, Permissions: map[string]struct{}{}, StoreIDs: map[string]struct{}{}}
	return buildViewer(ctx, r.Services, discovery)
}

// SelectWorkspace 切换权威 Session 的工作台并轮换签名 Cookie。
func (r *MutationResolver) SelectWorkspace(ctx context.Context, input gen.SelectWorkspaceInput) (*gen.Viewer, error) {
	principal, principalErr := auth.RequirePrincipal(ctx)
	if principalErr != nil {
		return nil, principalErr
	}
	account, session, err := selectWorkspaceRecords(ctx, r.Services, principal, input)
	if err != nil {
		return nil, err
	}
	claims := auth.SessionClaims{SessionID: session.ID, AccountID: account.ID, WorkspaceType: auth.WorkspaceType(session.WorkspaceType), OrganizationID: session.OrganizationID, CredentialVersion: account.CredentialVersion}
	token, err := auth.SignSessionClaims(r.Services.SecurityConfig, claims)
	if err != nil {
		return nil, err
	}
	if writer := ResponseWriterFromContext(ctx); writer != nil {
		auth.SetSessionCookie(writer, token, r.Services.SecurityConfig)
	}
	resolved, err := r.Services.Principal.Resolve(ctx, token, time.Now())
	if err != nil {
		return nil, err
	}
	return buildViewer(ctx, r.Services, resolved)
}
