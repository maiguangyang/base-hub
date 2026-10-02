/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"strings"
	"testing"

	"base-engine/gen"
)

// TestGeneratedHandlerCoverage 锁定所有 generated 数据入口均被授权层接管。
func TestGeneratedHandlerCoverage(t *testing.T) {
	defaults := gen.DefaultResolutionHandlers()
	registered := RegisterHandlers(defaults)
	if err := ValidateHandlerCoverage(registered, defaults); err != nil {
		t.Fatal(err)
	}
	registered.StoreMembers = defaults.StoreMembers
	err := ValidateHandlerCoverage(registered, defaults)
	if err == nil || !strings.Contains(err.Error(), "StoreMembers") {
		t.Fatalf("missing StoreMembers was not detected: %v", err)
	}
}
