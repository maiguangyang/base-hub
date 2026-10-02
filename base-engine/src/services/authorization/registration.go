/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"fmt"
	"reflect"

	"base-engine/gen"
	sessionservice "base-engine/src/services/session"
)

// HandlerDependencies supplies handwritten side effects to generated write handlers.
type HandlerDependencies struct {
	Sessions *sessionservice.Service
}

// RegisterHandlers 用显式授权处理器覆盖全部 generated 数据入口。
func RegisterHandlers(handlers gen.ResolutionHandlers, dependencies ...HandlerDependencies) gen.ResolutionHandlers {
	var services HandlerDependencies
	if len(dependencies) > 0 {
		services = dependencies[0]
	}
	registerIdentityHandlers(&handlers)
	registerOrganizationHandlers(&handlers, services)
	registerRoleHandlers(&handlers, services)
	registerStoreHandlers(&handlers)
	registerAuditHandlers(&handlers)
	registerOpeningRecordHandlers(&handlers)
	registerPaymentConfigHandlers(&handlers)
	registerCustomerHandlers(&handlers)
	registerProductHandlers(&handlers)
	// Verify the handwritten handlers before wrapping query functions, so the
	// wrapper cannot mask an accidental generated fallback.
	if err := ValidateHandlerCoverage(handlers, gen.DefaultResolutionHandlers()); err != nil {
		panic(err)
	}
	rejectGeneratedStocktakeProjections(&handlers)
	return handlers
}

// ValidateHandlerCoverage 检测任何回退到 generated 默认实现的入口。
func ValidateHandlerCoverage(handlers, defaults gen.ResolutionHandlers) error {
	registeredValue := reflect.ValueOf(handlers)
	defaultValue := reflect.ValueOf(defaults)
	handlerType := registeredValue.Type()
	for index := 0; index < registeredValue.NumField(); index++ {
		name := handlerType.Field(index).Name
		if name == "OnEvent" || name == "WebSocket" {
			continue
		}
		registered := registeredValue.Field(index)
		fallback := defaultValue.Field(index)
		if registered.IsNil() || registered.Pointer() == fallback.Pointer() {
			return fmt.Errorf("HANDLER_ISOLATION_REQUIRED:%s", name)
		}
	}
	return nil
}
