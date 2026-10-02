package productcatalog

import (
	"context"
	"testing"
)

func TestProductSearchByCurrentBarcodeAndClear(t *testing.T) {
	service, _, principal := catalogFixture(t)
	sku, pack := barcodePackage(t, service, principal, "0012345678905")
	ctx := context.Background()
	q := "0012345678905"
	page, err := service.Products(ctx, principal, nil, CatalogFilter{Q: &q}, 1, 1)
	catalogNoError(t, err)
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != sku.ProductID {
		t.Fatalf("barcode search: %+v", page)
	}
	_, err = service.SetPackageBarcode(ctx, principal, pack.ID, "")
	catalogNoError(t, err)
	page, err = service.Products(ctx, principal, nil, CatalogFilter{Q: &q}, 1, 1)
	catalogNoError(t, err)
	if page.Total != 0 {
		t.Fatal("cleared barcode still matches")
	}
}
