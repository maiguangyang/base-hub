/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

// New 构造 gqlgen 配置并注册全部授权 handler 与指令。
func New(db *gen.DB, ec *gen.EventController, dependencies ...Dependencies) gen.Config {
	services := Dependencies{DB: db.Query()}
	if len(dependencies) > 0 {
		services = dependencies[0]
	}
	resolver := NewResolver(db, ec, services)
	resolver.Handlers.OnEvent = func(context.Context, *gen.GeneratedResolver, *gen.Event) error {
		return nil
	}
	resolver.Handlers.WebSocket = func(context.Context, *gen.GeneratedResolver) (<-chan any, error) {
		return nil, auth.NewError(auth.CodeSessionPubSubRequired)
	}
	configuration := gen.Config{Resolvers: resolver}
	configureDirectives(&configuration, db)
	return configuration
}
