package dbup

import "testing"

func TestProductEntityPermissionsAreSeeded(t *testing.T) {
	byAction := make(map[string]PermissionSeed)
	for _, seed := range PermissionSeeds() {
		byAction[seed.Action] = seed
	}
	for _, resource := range []string{
		"productCategory", "product", "productSku", "productPackage",
		"storeListing", "storePackageOffer", "storePriceRevision",
		"storeInventoryBatch", "storeStockBalance", "storeStockMovement",
		"storePromotion", "storePromotionTarget",
	} {
		for _, verb := range []string{"read", "create", "update", "delete"} {
			action := resource + ":" + verb
			if _, ok := byAction[action]; !ok {
				t.Errorf("missing permission seed %s", action)
			}
		}
	}
}
