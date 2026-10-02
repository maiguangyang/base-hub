package tools

import (
	"errors"

	"github.com/google/jsonschema-go/jsonschema"
	"base-engine/src/services/ai"
	"google.golang.org/adk/v2/tool"
)

func Specs() []ai.ToolSpec { return append(hqFranchiseSpecs(), remainingSpecs()...) }

func PreparedSpecs() ([]ai.ToolSpec, error) {
	specs := Specs()
	for index := range specs {
		var schema *jsonschema.Schema
		var err error
		if specs[index].Document != "" {
			schema, err = reviewedInputSchema(specs[index])
		} else if specs[index].Path == "/api/franchise-initial-account" {
			schema = initialAccountInputSchema(specs[index])
		} else {
			schema = paymentInputSchema(specs[index])
		}
		if err != nil {
			return nil, err
		}
		if schema == nil {
			return nil, errors.New("UNREVIEWED_TOOL_INPUT_SCHEMA")
		}
		specs[index].InputSchema, err = schema.Resolve(nil)
		if err != nil {
			return nil, err
		}
	}
	return specs, nil
}

func Build(runtime ai.FixedToolRuntime) ([]tool.Tool, error) {
	if runtime.Call == nil || runtime.DeliverSecret == nil {
		return nil, errors.New("INVALID_TOOL_RUNTIME")
	}
	items, err := buildHQFranchiseTools(runtime, hqFranchiseSpecs())
	if err != nil {
		return nil, err
	}
	for _, spec := range remainingSpecs() {
		var item tool.Tool
		if spec.Path != "" {
			item, err = buildPaymentTool(spec, runtime)
		} else {
			item, err = buildReviewedTool(spec, runtime)
		}
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
