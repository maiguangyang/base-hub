/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"context"
	"errors"
	"testing"
)

// TestPresentErrorPreservesCodesAndRedactsInternals 验证统一错误呈现契约。
func TestPresentErrorPreservesCodesAndRedactsInternals(t *testing.T) {
	known := PresentError(context.Background(), NewError(CodePermissionDenied))
	if known.Message != string(CodePermissionDenied) || known.Extensions["code"] != string(CodePermissionDenied) {
		t.Fatalf("known error changed: %#v", known)
	}

	internal := PresentError(context.Background(), errors.New("database password leaked"))
	if internal.Message != string(CodeInternalError) || internal.Extensions["code"] != string(CodeInternalError) {
		t.Fatalf("internal error not redacted: %#v", internal)
	}
}
