package agentkit

import (
	"sync"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/tool"
)

type toolFailureTracker struct {
	mu    sync.Mutex
	first error
}

func (tracker *toolFailureTracker) record(err error) {
	if err == nil {
		return
	}
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	if tracker.first == nil {
		tracker.first = err
	}
}

func (tracker *toolFailureTracker) Err() error {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	return tracker.first
}

func trackToolFailures(callbacks []llmagent.AfterToolCallback, record func(error)) []llmagent.AfterToolCallback {
	return []llmagent.AfterToolCallback{func(ctx agent.Context, item tool.Tool, args, output map[string]any, callErr error) (map[string]any, error) {
		record(callErr)
		for _, callback := range callbacks {
			result, err := callback(ctx, item, args, output, callErr)
			record(err)
			if result != nil || err != nil {
				return result, err
			}
		}
		return output, callErr
	}}
}
