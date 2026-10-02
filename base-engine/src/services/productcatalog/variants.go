package productcatalog

import (
	"sort"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

type SpecificationSelection struct {
	SpecificationID string
	ValueIDs        []string
}

type SkuOverride struct {
	ValueIDs              []string
	Enabled               bool
	PackageTemplateID     *string
	DisableDefaultPackage bool
}

type variantChoice struct {
	values     []gen.SpecificationValue
	key        string
	name       string
	restricted bool
}

func variantKey(ids []string) string {
	copyIDs := append([]string(nil), ids...)
	sort.Strings(copyIDs)
	return strings.Join(copyIDs, ",")
}

func resolveVariants(tx *gorm.DB, hqID string, input ProductInput) ([]variantChoice, []string, error) {
	variants := []variantChoice{{values: []gen.SpecificationValue{}}}
	seenDefinitions := make(map[string]struct{}, len(input.Selections))
	selected := make([]string, 0)
	for _, selection := range input.Selections {
		if len(variants)*len(selection.ValueIDs) > 100 {
			return nil, nil, auth.NewError(auth.CodeValidationFailed)
		}
		values, definitionEnabled, err := loadSelectionValues(tx, hqID, selection, seenDefinitions)
		if err != nil {
			return nil, nil, err
		}
		for _, value := range values {
			selected = append(selected, value.ID)
		}
		variants = expandVariants(variants, values, !definitionEnabled)
	}
	for i := range variants {
		labelVariant(&variants[i])
	}
	return variants, selected, nil
}

func loadSelectionValues(tx *gorm.DB, hqID string, selection SpecificationSelection, seen map[string]struct{}) ([]gen.SpecificationValue, bool, error) {
	if selection.SpecificationID == "" || len(selection.ValueIDs) == 0 {
		return nil, false, auth.NewError(auth.CodeValidationFailed)
	}
	if _, exists := seen[selection.SpecificationID]; exists {
		return nil, false, auth.NewError(auth.CodeValidationFailed)
	}
	seen[selection.SpecificationID] = struct{}{}
	var definition gen.SpecificationDefinition
	if err := tx.Where("id = ? AND organization_id = ?", selection.SpecificationID, hqID).First(&definition).Error; err != nil {
		return nil, false, auth.NewError(auth.CodePermissionDenied)
	}
	values := make([]gen.SpecificationValue, 0, len(selection.ValueIDs))
	seenValues := make(map[string]struct{}, len(selection.ValueIDs))
	for _, id := range selection.ValueIDs {
		if id == "" {
			return nil, false, auth.NewError(auth.CodeValidationFailed)
		}
		if _, exists := seenValues[id]; exists {
			return nil, false, auth.NewError(auth.CodeValidationFailed)
		}
		seenValues[id] = struct{}{}
		var value gen.SpecificationValue
		if err := tx.Where("id = ? AND specification_id = ?", id, definition.ID).First(&value).Error; err != nil {
			return nil, false, auth.NewError(auth.CodePermissionDenied)
		}
		values = append(values, value)
	}
	return values, definition.Enabled, nil
}

func expandVariants(variants []variantChoice, values []gen.SpecificationValue, definitionRestricted bool) []variantChoice {
	next := make([]variantChoice, 0, len(variants)*len(values))
	for _, variant := range variants {
		for _, value := range values {
			copyValues := append(append([]gen.SpecificationValue(nil), variant.values...), value)
			next = append(next, variantChoice{values: copyValues, restricted: variant.restricted || definitionRestricted || !value.Enabled})
		}
	}
	return next
}

func labelVariant(variant *variantChoice) {
	ids := make([]string, 0, len(variant.values))
	names := make([]string, 0, len(variant.values))
	for _, value := range variant.values {
		ids = append(ids, value.ID)
		names = append(names, value.Name)
	}
	variant.key = variantKey(ids)
	label := strings.Join(names, " / ")
	if label == "" {
		label = "默认规格"
	}
	name := []rune(label)
	if len(name) > 128 {
		name = name[:128]
	}
	variant.name = string(name)
}
