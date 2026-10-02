/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package src

import (
	"context"

	"base-engine/gen"
	sessionservice "base-engine/src/services/session"
)

// UpdateOperatorMembership publishes authority-change events only after commit.
func (r *MutationResolver) UpdateOperatorMembership(ctx context.Context, id string, input map[string]interface{}) (*gen.OperatorMembership, error) {
	ctx = sessionservice.WithPendingEvents(ctx)
	ctx = gen.EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err := r.Handlers.UpdateOperatorMembership(ctx, r.GeneratedResolver, id, input)
	if err = finishGeneratedMutation(ctx, r.GeneratedResolver, r.Services.Publisher, err); err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteOperatorMemberships publishes authority-change events only after commit.
func (r *MutationResolver) DeleteOperatorMemberships(ctx context.Context, ids []string, unscoped *bool) (bool, error) {
	ctx = sessionservice.WithPendingEvents(ctx)
	ctx = gen.EnrichContextWithMutations(ctx, r.GeneratedResolver)
	done, err := r.Handlers.DeleteOperatorMemberships(ctx, r.GeneratedResolver, ids, unscoped)
	err = finishGeneratedMutation(ctx, r.GeneratedResolver, r.Services.Publisher, err)
	return done && err == nil, err
}

// UpdateOperatorRole publishes member authority-change events only after commit.
func (r *MutationResolver) UpdateOperatorRole(ctx context.Context, id string, input map[string]interface{}) (*gen.OperatorRole, error) {
	ctx = sessionservice.WithPendingEvents(ctx)
	ctx = gen.EnrichContextWithMutations(ctx, r.GeneratedResolver)
	item, err := r.Handlers.UpdateOperatorRole(ctx, r.GeneratedResolver, id, input)
	if err = finishGeneratedMutation(ctx, r.GeneratedResolver, r.Services.Publisher, err); err != nil {
		return nil, err
	}
	return item, nil
}

func finishGeneratedMutation(ctx context.Context, resolver *gen.GeneratedResolver, publisher sessionservice.Publisher, handlerError error) error {
	if handlerError != nil {
		_ = gen.RollbackMutationContext(ctx, resolver)
		return handlerError
	}
	if err := gen.FinishMutationContext(ctx, resolver); err != nil {
		_ = gen.RollbackMutationContext(ctx, resolver)
		return err
	}
	sessionservice.PublishPending(ctx, publisher)
	return nil
}
