/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Code 是供各客户端稳定判断的业务错误码。
type Code string

const (
	CodeAuthRequired                    Code = "AUTH_REQUIRED"
	CodeInvalidCredentials              Code = "INVALID_CREDENTIALS"
	CodePhoneNotFound                   Code = "PHONE_NOT_FOUND"
	CodePasswordIncorrect               Code = "PASSWORD_INCORRECT"
	CodeAccountInactive                 Code = "ACCOUNT_INACTIVE"
	CodeAccountUnavailable              Code = "ACCOUNT_UNAVAILABLE"
	CodePasswordWeak                    Code = "PASSWORD_WEAK"
	CodePasswordConfirmMismatch         Code = "PASSWORD_CONFIRM_MISMATCH"
	CodeRateLimited                     Code = "RATE_LIMITED"
	CodeTempPasswordChangeRequired      Code = "TEMP_PASSWORD_CHANGE_REQUIRED"
	CodeAccountLocked                   Code = "ACCOUNT_LOCKED"
	CodeWorkspaceForbidden              Code = "WORKSPACE_FORBIDDEN"
	CodeOrganizationSuspended           Code = "ORGANIZATION_SUSPENDED"
	CodeMembershipInactive              Code = "MEMBERSHIP_INACTIVE"
	CodePermissionDenied                Code = "PERMISSION_DENIED"
	CodeStoreScopeDenied                Code = "STORE_SCOPE_DENIED"
	CodeStoreNotActive                  Code = "STORE_NOT_ACTIVE"
	CodeInvitationPending               Code = "INVITATION_PENDING"
	CodeInvitationExpired               Code = "INVITATION_EXPIRED"
	CodeLastOwnerRequired               Code = "LAST_OWNER_REQUIRED"
	CodePermissionDelegationDenied      Code = "PERMISSION_DELEGATION_DENIED"
	CodeLastHQSuperAdminRequired        Code = "LAST_HQ_SUPER_ADMIN_REQUIRED"
	CodeSelfMembershipChangeDenied      Code = "SELF_MEMBERSHIP_CHANGE_DENIED"
	CodeRoleInUse                       Code = "ROLE_IN_USE"
	CodeSessionRevoked                  Code = "SESSION_REVOKED"
	CodeCredentialsChanged              Code = "CREDENTIALS_CHANGED"
	CodeConflict                        Code = "CONFLICT"
	CodeStocktakeReconciliationRequired Code = "STOCKTAKE_RECONCILIATION_REQUIRED"
	CodeCustomerPointsMismatch          Code = "CUSTOMER_POINTS_MISMATCH"
	CodeOpeningRecordNumberConflict     Code = "OPENING_RECORD_NUMBER_CONFLICT"
	CodeValidationFailed                Code = "VALIDATION_FAILED"
	CodeInternalError                   Code = "INTERNAL_ERROR"
	CodeHQAlreadyBootstrapped           Code = "HQ_ALREADY_BOOTSTRAPPED"
	CodeGenerationRequired              Code = "GENERATION_REQUIRED"
	CodeSessionPubSubRequired           Code = "SESSION_PUBSUB_REQUIRED"
)

var ErrAuthRequired = NewError(CodeAuthRequired)

// NewError 创建不泄露内部细节的 GraphQL 错误。
func NewError(code Code) *gqlerror.Error {
	return &gqlerror.Error{Message: string(code), Extensions: map[string]any{"code": string(code)}}
}

// PresentError 保留已知业务码，并隐藏未分类的内部错误细节。
func PresentError(ctx context.Context, err error) *gqlerror.Error {
	presented := graphql.DefaultErrorPresenter(ctx, err)
	if presented == nil {
		return nil
	}
	if code, ok := presented.Extensions["code"].(string); ok && code != "" {
		presented.Message = code
		return presented
	}
	return &gqlerror.Error{
		Message: string(CodeInternalError), Path: presented.Path, Locations: presented.Locations,
		Extensions: map[string]any{"code": string(CodeInternalError)},
	}
}

// ErrorCode 从错误链提取稳定业务码。
func ErrorCode(err error) Code {
	var target *gqlerror.Error
	if !errors.As(err, &target) {
		return ""
	}
	code, _ := target.Extensions["code"].(string)
	return Code(code)
}
