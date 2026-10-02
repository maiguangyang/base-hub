package storemerchandising

import (
	"context"
	"testing"
	"time"

	"base-engine/auth"
)

func assertStocktakeCreatedAtFilter(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, created *StocktakeView) {
	t.Helper()
	before := created.StartedAt.Add(-time.Hour)
	after := created.StartedAt.Add(time.Hour)
	for _, filter := range []StocktakeFilter{{From: &after}, {To: &before}} {
		outside, err := service.Stocktakes(context.Background(), principal, "store", filter, 1, 20)
		merchandisingNoError(t, err)
		if outside.Total != 0 {
			t.Fatalf("creation date outside range returned %+v", outside)
		}
	}
	inside, err := service.Stocktakes(context.Background(), principal, "store", StocktakeFilter{From: &before, To: &after}, 1, 20)
	merchandisingNoError(t, err)
	if inside.Total != 1 || inside.Data[0].ID != created.ID {
		t.Fatalf("creation date inside range returned %+v", inside)
	}
}
