package ai

import (
	"context"
	"net/http"
	"time"

	"base-engine/auth"
)

type aiRunIDKey struct{}

func runIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(aiRunIDKey{}).(string)
	return value
}

func (s *Service) servePhase(w http.ResponseWriter, r *http.Request, principal *auth.WorkspacePrincipal, token string, phase Phase, release func(), run func(context.Context, *eventSink) error) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Minute)
	defer cancel()
	runID := randomPlanID()
	ctx = context.WithValue(ctx, aiRunIDKey{}, runID)
	sink := newEventSink(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer release()
		defer sink.Close()
		_ = sink.Emit("run_started", map[string]any{"runId": runID, "phase": phase})
		err := run(ctx, sink)
		if err != nil {
			_ = sink.Emit("error", map[string]string{"code": safeRunErrorCode(err)})
		}
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		_ = sink.Emit("run_finished", map[string]string{"status": status})
	}()
	if err := writeSSE(ctx, w, sink.events, 15*time.Second); err != nil {
		cancel()
	}
	<-done
}
