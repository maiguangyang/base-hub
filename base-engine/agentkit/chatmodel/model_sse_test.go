package chatmodel

import (
	"strings"
	"testing"
)

func TestParseModelSSEPreservesMultilineData(t *testing.T) {
	var frames []string
	frame := "data: first\n:" + strings.Repeat("x", 4090) + "\ndata: second\n\ndata: [DONE]\n\n"
	err := parseModelSSE(strings.NewReader(frame), func(data []byte) error {
		frames = append(frames, string(data))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 1 || frames[0] != "first\nsecond" {
		t.Fatalf("frames = %q", frames)
	}
}
