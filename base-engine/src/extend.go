/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"time"

	"base-engine/config"
	"base-engine/src/services/audit"
	"base-engine/src/services/authentication"
	"base-engine/src/services/authorization"
	"base-engine/src/services/customer"
	"base-engine/src/services/membership"
	"base-engine/src/services/organization"
	"base-engine/src/services/productcatalog"
	sessionservice "base-engine/src/services/session"
	storeservice "base-engine/src/services/store"
	"base-engine/src/services/storemerchandising"
	"gorm.io/gorm"
)

// Dependencies 汇总自定义 GraphQL resolver 的业务服务。
type Dependencies struct {
	DB             *gorm.DB
	SecurityConfig config.SecurityConfig
	Audit          *audit.Service
	Authentication *authentication.Service
	Sessions       *sessionservice.Service
	Principal      *authorization.PrincipalResolver
	Organizations  *organization.Service
	Memberships    *membership.Service
	Roles          *membership.RoleService
	Stores         *storeservice.Service
	Customers      *customer.Service
	ProductCatalog *productcatalog.Service
	Merchandising  *storemerchandising.Service
	Publisher      sessionservice.Publisher
}

// NewDependencies 构造共享同一数据库、配置与发布器的完整依赖图。
func NewDependencies(db *gorm.DB, cfg config.SecurityConfig, publisher sessionservice.Publisher) Dependencies {
	auditService := audit.NewService()
	sessions := sessionservice.NewService(cfg.SessionDuration)
	limit := cfg.LoginIPAttemptsPerMinute
	if limit < 1 {
		limit = 20
	}
	limiter := authentication.NewLoginLimiter(limit, time.Minute)
	products := productcatalog.NewService(db, auditService)
	stores := storeservice.NewService(db, auditService)
	stores.SetDocumentRoot(products.ImageRoot())
	return Dependencies{
		DB: db, SecurityConfig: cfg,
		Audit:          auditService,
		Authentication: authentication.NewService(db, cfg, sessions, limiter, auditService, publisher),
		Sessions:       sessions, Principal: authorization.NewPrincipalResolver(db, cfg),
		Organizations: organization.NewService(db, auditService, publisher),
		Memberships: membership.NewService(db, auditService, membership.ServiceDependencies{
			Sessions: sessions, Publisher: publisher,
		}),
		Roles:     membership.NewRoleService(db, auditService),
		Stores:    stores,
		Customers: customer.NewService(db, auditService), ProductCatalog: products,
		Merchandising: storemerchandising.NewService(db, auditService), Publisher: publisher,
	}
}
