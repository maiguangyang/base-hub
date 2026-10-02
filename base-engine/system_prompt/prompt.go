package system_prompt

import (
	_ "embed"
	"encoding/json"
	"sync"
)

//go:embed admin_agent.md
var AdminAgent string

//go:embed tool_descriptions.json
var toolDescriptionsJSON []byte

var descriptionsOnce sync.Once
var descriptions map[string]string

func ToolDescription(id string) (string, bool) {
	descriptionsOnce.Do(func() { _ = json.Unmarshal(toolDescriptionsJSON, &descriptions) })
	value, ok := descriptions[id]
	return value, ok && value != ""
}
