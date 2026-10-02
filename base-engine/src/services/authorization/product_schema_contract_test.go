package authorization

import (
	"os"
	"strings"
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func TestProductSchemaRelationships(t *testing.T) {
	source, err := os.ReadFile("../../../model/model.graphql")
	if err != nil {
		t.Fatal(err)
	}
	document, err := parser.ParseSchema(&ast.Source{Input: string(source)})
	if err != nil {
		t.Fatal(err)
	}
	definitions := document.Definitions
	pairs := []struct{ left, field, right, inverse string }{
		{"Organization", "productCategories", "ProductCategory", "organization"},
		{"Organization", "products", "Product", "organization"},
		{"Organization", "productBrands", "ProductBrand", "organization"},
		{"Organization", "productPackageTemplates", "ProductPackageTemplate", "organization"},
		{"ProductBrand", "products", "Product", "brand"},
		{"ProductCategory", "parent", "ProductCategory", "children"},
		{"ProductCategory", "products", "Product", "category"},
		{"Product", "skus", "ProductSku", "product"},
		{"ProductSku", "packages", "ProductPackage", "sku"},
		{"ProductPackageTemplate", "createdPackages", "ProductPackage", "template"},
		{"ProductPackageTemplate", "containsPackage", "ProductPackageTemplate", "containedByPackages"},
		{"ProductSku", "listings", "StoreListing", "sku"},
		{"ProductPackage", "containsPackage", "ProductPackage", "containedByPackages"},
		{"ProductPackage", "offers", "StorePackageOffer", "package"},
		{"ProductPackage", "balances", "StoreStockBalance", "package"},
		{"ProductPackage", "movementSources", "StoreStockMovement", "sourcePackage"},
		{"ProductPackage", "movementTargets", "StoreStockMovement", "targetPackage"},
		{"Store", "productListings", "StoreListing", "store"},
		{"Store", "promotions", "StorePromotion", "store"},
		{"StoreListing", "offers", "StorePackageOffer", "listing"},
		{"StoreListing", "batches", "StoreInventoryBatch", "listing"},
		{"StorePackageOffer", "priceRevisions", "StorePriceRevision", "offer"},
		{"StorePackageOffer", "promotionTargets", "StorePromotionTarget", "offer"},
		{"StoreInventoryBatch", "balances", "StoreStockBalance", "batch"},
		{"StoreInventoryBatch", "movements", "StoreStockMovement", "batch"},
		{"StorePromotion", "targets", "StorePromotionTarget", "promotion"},
		{"Store", "couponTemplates", "CustomerCouponTemplate", "applicableStore"},
	}
	for _, pair := range pairs {
		t.Run(pair.left+"."+pair.field, func(t *testing.T) {
			assertTypedInverse(t, definitions, pair.left, pair.field, pair.right, pair.inverse)
			assertTypedInverse(t, definitions, pair.right, pair.inverse, pair.left, pair.field)
		})
	}
}

func TestStocktakeSchemaRelationships(t *testing.T) {
	source, err := os.ReadFile("../../../model/model.graphql")
	if err != nil {
		t.Fatal(err)
	}
	document, err := parser.ParseSchema(&ast.Source{Input: string(source)})
	if err != nil {
		t.Fatal(err)
	}
	pairs := []struct{ left, field, right, inverse string }{
		{"Store", "stocktakes", "StoreStocktake", "store"},
		{"Account", "createdStocktakes", "StoreStocktake", "initiatedByAccount"},
		{"Account", "postedStocktakes", "StoreStocktake", "postedBy"},
		{"StoreStocktake", "lines", "StoreStocktakeLine", "stocktake"},
		{"StoreInventoryBatch", "stocktakeLines", "StoreStocktakeLine", "batch"},
		{"ProductPackage", "stocktakeLines", "StoreStocktakeLine", "package"},
		{"StoreStocktakeLine", "movements", "StoreStockMovement", "stocktakeLine"},
	}
	for _, pair := range pairs {
		assertTypedInverse(t, document.Definitions, pair.left, pair.field, pair.right, pair.inverse)
		assertTypedInverse(t, document.Definitions, pair.right, pair.inverse, pair.left, pair.field)
	}
	if document.Definitions.ForName("StoreStockBalance").Fields.ForName("version") == nil {
		t.Fatal("stock balance version missing")
	}
	stocktake := document.Definitions.ForName("StoreStocktake")
	for _, generated := range []string{"createdAt", "createdBy"} {
		if stocktake.Fields.ForName(generated) != nil {
			t.Fatalf("handwritten stocktake field %s collides with generated entity metadata", generated)
		}
	}
	if stocktake.Fields.ForName("startedAt") == nil {
		t.Fatal("stocktake startedAt missing")
	}
}

func assertTypedInverse(t *testing.T, definitions ast.DefinitionList, owner, fieldName, target, inverse string) {
	t.Helper()
	definition := definitions.ForName(owner)
	if definition == nil || definition.Directives.ForName("entity") == nil {
		t.Fatalf("missing entity %s", owner)
	}
	field := definition.Fields.ForName(fieldName)
	if field == nil {
		t.Fatalf("missing %s.%s", owner, fieldName)
	}
	if field.Type.Name() != target {
		t.Fatalf("%s.%s must target %s, got %s", owner, fieldName, target, field.Type.Name())
	}
	relation := field.Directives.ForName("relationship")
	if relation == nil || relation.Arguments.ForName("inverse") == nil ||
		relation.Arguments.ForName("inverse").Value.Raw != inverse {
		t.Fatalf("%s.%s must declare inverse %s", owner, fieldName, inverse)
	}
}

func TestProductEntitiesHaveNoScalarRelationshipIDs(t *testing.T) {
	source, err := os.ReadFile("../../../model/model.graphql")
	if err != nil {
		t.Fatal(err)
	}
	document, err := parser.ParseSchema(&ast.Source{Input: string(source)})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"ProductCategory", "ProductBrand", "Product", "ProductSku", "ProductPackage", "SpecificationDefinition", "SpecificationValue", "ProductSpecificationChoice", "ProductSkuSpecificationValue", "ProductPackageTemplate", "StoreListing",
		"StorePackageOffer", "StorePriceRevision", "StoreInventoryBatch", "StoreStockBalance",
		"StoreStockMovement", "StoreStocktake", "StoreStocktakeLine", "StorePromotion", "StorePromotionTarget",
	} {
		definition := document.Definitions.ForName(name)
		if definition == nil {
			t.Fatalf("missing entity %s", name)
		}
		for _, field := range definition.Fields {
			if strings.HasSuffix(field.Name, "Id") || strings.HasSuffix(field.Name, "Ids") {
				t.Errorf("%s.%s is a handwritten scalar relationship candidate", name, field.Name)
			}
		}
	}
}
