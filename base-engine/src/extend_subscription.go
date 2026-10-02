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

// SessionEvents 只订阅当前权威 Session 的终态事件。
func (r *SubscriptionResolver) SessionEvents(ctx context.Context) (<-chan *gen.SessionEvent, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if r.Services.Publisher == nil {
		return nil, generationRequiredError()
	}
	return r.Services.Publisher.Subscribe(ctx, principal.SessionID)
}

func generationRequiredError() error {
	return auth.NewError(auth.CodeSessionPubSubRequired)
}
