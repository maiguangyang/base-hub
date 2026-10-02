package ai

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadOnlyReplyPreviewUsesArrayFields(t *testing.T) {
	state := &runState{}
	var frame bytes.Buffer
	if err := writeSSEEvent(&frame, AIEvent{Name: "preview_ready", Data: state.previewView("你好")}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(frame.String(), `"steps":[]`) || !strings.Contains(frame.String(), `"requiredInputs":[]`) {
		t.Fatalf("read-only preview is not a valid array-shaped SSE payload: %s", frame.String())
	}
}
