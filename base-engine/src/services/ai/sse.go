package ai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

const maxSSEFrameBytes = 64 * 1024

type AIEvent struct {
	Name string
	Data any
}

type eventSink struct {
	ctx    context.Context
	events chan AIEvent
	mu     sync.Mutex
	closed bool
}

func newEventSink(ctx context.Context) *eventSink {
	return &eventSink{ctx: ctx, events: make(chan AIEvent, 64)}
}

func (sink *eventSink) Emit(name string, data any) error {
	if sink == nil || !validEventName(name) {
		return errors.New("INVALID_AI_EVENT")
	}
	sink.mu.Lock()
	defer sink.mu.Unlock()
	if sink.closed {
		return errors.New("AI_EVENT_STREAM_CLOSED")
	}
	select {
	case <-sink.ctx.Done():
		return sink.ctx.Err()
	case sink.events <- AIEvent{Name: name, Data: data}:
		return nil
	}
}

func (sink *eventSink) Close() {
	if sink == nil {
		return
	}
	sink.mu.Lock()
	defer sink.mu.Unlock()
	if sink.closed {
		return
	}
	sink.closed = true
	close(sink.events)
}

func validEventName(name string) bool {
	switch name {
	case "run_started", "text_delta", "tool_started", "tool_finished", "preview_ready", "secret", "run_finished", "error":
		return true
	default:
		return false
	}
}

func writeSSE(ctx context.Context, writer http.ResponseWriter, events <-chan AIEvent, heartbeat time.Duration) error {
	flusher, ok := writer.(http.Flusher)
	if !ok || heartbeat <= 0 {
		return errors.New("SSE_UNAVAILABLE")
	}
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-store")
	writer.Header().Set("X-Accel-Buffering", "no")
	writer.WriteHeader(http.StatusOK)
	flusher.Flush()
	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, open := <-events:
			if !open {
				return nil
			}
			if err := writeSSEEvent(writer, event); err != nil {
				return err
			}
			flusher.Flush()
		case <-ticker.C:
			if _, err := io.WriteString(writer, ": keepalive\n\n"); err != nil {
				return err
			}
			flusher.Flush()
		}
	}
}

func writeSSEEvent(writer io.Writer, event AIEvent) error {
	if !validEventName(event.Name) {
		return errors.New("INVALID_AI_EVENT")
	}
	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	frame := []byte("event: " + event.Name + "\ndata: " + string(data) + "\n\n")
	if len(frame) > maxSSEFrameBytes {
		return errors.New("SSE_FRAME_TOO_LARGE")
	}
	_, err = writer.Write(frame)
	return err
}
