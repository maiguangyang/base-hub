/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import "context"

// UserTokenToMap 为尚未迁移的调用方提供稳定标识，不解析客户端权限。
func UserTokenToMap(ctx context.Context) (map[string]interface{}, error) {
	principal, err := requirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"accountId":     principal.AccountID,
		"sessionId":     principal.SessionID,
		"workspaceType": principal.WorkspaceType,
	}, nil
}
