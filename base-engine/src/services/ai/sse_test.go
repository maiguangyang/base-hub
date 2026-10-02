package ai

import (
	"context"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSSESerializesEventsAndHeartbeat(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	sink := newEventSink(ctx)
	response := httptest.NewRecorder()
	done := make(chan error, 1)
	go func() { done <- writeSSE(ctx, response, sink.events, time.Millisecond) }()
	var group sync.WaitGroup
	for index := range 10 {
		group.Go(func() {
			if err := sink.Emit("tool_finished", map[string]any{"sequence": index}); err != nil {
				t.Error(err)
			}
		})
	}
	group.Wait()
	time.Sleep(3 * time.Millisecond)
	sink.Close()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	body := response.Body.String()
	if strings.Count(body, "event: tool_finished\n") != 10 {
		t.Fatalf("event count: %s", body)
	}
	if !strings.Contains(body, ": keepalive\n\n") {
		t.Fatalf("heartbeat missing: %s", body)
	}
	for _, frame := range strings.Split(body, "\n\n") {
		if strings.HasPrefix(frame, "event:") && !strings.Contains(frame, "\ndata: ") {
			t.Fatalf("broken frame: %q", frame)
		}
	}
}

func TestSSERejectsOversizedFrame(t *testing.T) {
	response := httptest.NewRecorder()
	ch := make(chan AIEvent, 1)
	ch <- AIEvent{Name: "text_delta", Data: map[string]any{"text": strings.Repeat("x", 70*1024)}}
	close(ch)
	if err := writeSSE(t.Context(), response, ch, time.Hour); err == nil {
		t.Fatal("oversized frame accepted")
	}
}
