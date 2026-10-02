package authorization

import (
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"os"
	"testing"
)

func TestProductSpecificationSchemaRelationships(t *testing.T) {
	source, err := os.ReadFile("../../../model/model.graphql")
	if err != nil {
		t.Fatal(err)
	}
	document, err := parser.ParseSchema(&ast.Source{Input: string(source)})
	if err != nil {
		t.Fatal(err)
	}
	pairs := []struct{ left, field, right, inverse string }{
		{"Organization", "specificationDefinitions", "SpecificationDefinition", "organization"},
		{"Product", "specificationChoices", "ProductSpecificationChoice", "product"},
		{"Product", "defaultPackageTemplate", "ProductPackageTemplate", "defaultProducts"},
		{"SpecificationDefinition", "values", "SpecificationValue", "specification"},
		{"SpecificationValue", "productChoices", "ProductSpecificationChoice", "value"},
		{"SpecificationValue", "skuValues", "ProductSkuSpecificationValue", "value"},
		{"ProductSku", "specificationValues", "ProductSkuSpecificationValue", "sku"},
	}
	for _, pair := range pairs {
		assertTypedInverse(t, document.Definitions, pair.left, pair.field, pair.right, pair.inverse)
		assertTypedInverse(t, document.Definitions, pair.right, pair.inverse, pair.left, pair.field)
	}
}
