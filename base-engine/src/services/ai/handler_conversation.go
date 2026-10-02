package ai

import "strings"

func validChatHistory(history []chatTurn) bool {
	if len(history) > 12 {
		return false
	}
	total := 0
	for _, turn := range history {
		if turn.Role != "user" && turn.Role != "assistant" {
			return false
		}
		if strings.TrimSpace(turn.Text) == "" || len(turn.Text) > 2000 {
			return false
		}
		total += len(turn.Text)
	}
	return total <= 6000
}
