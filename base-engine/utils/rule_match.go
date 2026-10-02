package utils

import "regexp"

// MatchesRule reports whether value matches the named validation rule.
func MatchesRule(ruleName, value string) bool {
	rule, ok := Rule[ruleName]
	if !ok {
		return false
	}

	pattern, ok := rule["rgx"].(string)
	if !ok || pattern == "" {
		return false
	}

	matcher, err := regexp.Compile(pattern)
	return err == nil && matcher.MatchString(value)
}
