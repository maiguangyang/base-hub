package tools

import (
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"strings"
	"testing"
)

func TestProductReadToolIncludesFieldsRequiredToPreserveVariantsOnUpdate(t *testing.T) {
	for _, spec := range Specs() {
		if spec.ID != "HqProducts" {
			continue
		}
		selection := strings.SplitN(spec.Document, "data{", 2)
		if len(selection) != 2 {
			t.Fatal("HqProducts query has no product selection")
		}
		for _, field := range []string{"brandId", "description", "selectedValueIds", "defaultPackageTemplateId"} {
			if !strings.Contains(selection[1], field) {
				t.Errorf("HqProducts query omits %s needed for safe update", field)
			}
		}
		return
	}
	t.Fatal("HqProducts tool missing")
}
func TestFranchiseStockBatchesToolSupportsStoreWideBatchRows(t *testing.T) {
	for _, spec := range Specs() {
		if spec.ID != "FranchiseStockBatches" {
			continue
		}
		document, err := parser.ParseQuery(&ast.Source{Input: spec.Document})
		if err != nil {
			t.Fatal(err)
		}
		operation := document.Operations.ForName("FranchiseStockBatches")
		if operation == nil {
			t.Fatal("FranchiseStockBatches operation missing")
		}
		listing := operation.VariableDefinitions.ForName("listingId")
		if listing == nil || listing.Type.NamedType != "ID" || listing.Type.NonNull {
			t.Fatal("listingId must be optional ID, not ID!")
		}
		root := stockBatchToolField(t, operation.SelectionSet, "franchiseStockBatches")
		rows := stockBatchToolField(t, root.SelectionSet, "data")
		for _, name := range []string{"skuId", "productName", "skuName"} {
			stockBatchToolField(t, rows.SelectionSet, name)
		}
		packages := stockBatchToolField(t, rows.SelectionSet, "packages")
		for _, name := range []string{"id", "name", "enabled", "containsPackageId", "containsQuantity", "packageSetVersion"} {
			stockBatchToolField(t, packages.SelectionSet, name)
		}
		return
	}
	t.Fatal("FranchiseStockBatches tool missing")
}

func stockBatchToolField(t *testing.T, selection ast.SelectionSet, name string) *ast.Field {
	t.Helper()
	for _, item := range selection {
		if field, ok := item.(*ast.Field); ok && field.Name == name {
			return field
		}
	}
	t.Fatalf("missing batch tool field %s", name)
	return nil
}
