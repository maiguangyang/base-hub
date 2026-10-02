package utils

import "testing"

func TestMatchesRule(t *testing.T) {
	tests := []struct {
		name     string
		ruleName string
		value    string
		want     bool
	}{
		{name: "valid phone", ruleName: "phone", value: "13800000000", want: true},
		{name: "invalid phone", ruleName: "phone", value: "123", want: false},
		{name: "valid email", ruleName: "email", value: "staff@example.com", want: true},
		{name: "invalid email", ruleName: "email", value: "invalid", want: false},
		{name: "unknown rule", ruleName: "missing", value: "anything", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := MatchesRule(test.ruleName, test.value); got != test.want {
				t.Fatalf("MatchesRule(%q, %q) = %v, want %v", test.ruleName, test.value, got, test.want)
			}
		})
	}
}
