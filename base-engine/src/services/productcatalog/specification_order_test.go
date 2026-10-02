package productcatalog

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestSpecificationValuesKeepSavedOrderAndAppendNewValues(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	spec := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "口味"))
	first := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "特辣"))
	second := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "微辣"))
	assertSpecificationOrder(t, service, principal, spec.ID, []string{first.ID, second.ID})
	if err := service.ReorderSpecificationValues(ctx, principal, spec.ID, []string{second.ID, first.ID}); err != nil {
		t.Fatal(err)
	}
	third := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "中辣"))
	assertSpecificationOrder(t, service, principal, spec.ID, []string{second.ID, first.ID, third.ID})
}

func TestReorderSpecificationValuesRejectsIncompleteDuplicateAndForeignIDs(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	spec := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "口味"))
	other := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "产地"))
	first := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "微辣"))
	second := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "特辣"))
	foreign := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, other.ID, "国产"))
	for _, ids := range [][]string{{first.ID}, {first.ID, first.ID}, {first.ID, foreign.ID}} {
		if err := service.ReorderSpecificationValues(ctx, principal, spec.ID, ids); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("ids %v: expected validation failure, got %v", ids, err)
		}
		assertSpecificationOrder(t, service, principal, spec.ID, []string{first.ID, second.ID})
	}
	principal.WorkspaceType = auth.WorkspaceTypeFranchise
	if err := service.ReorderSpecificationValues(ctx, principal, spec.ID, []string{second.ID, first.ID}); auth.ErrorCode(err) != auth.CodeWorkspaceForbidden {
		t.Fatalf("franchise reorder = %v", err)
	}
}

func assertSpecificationOrder(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, specificationID string, want []string) {
	t.Helper()
	page, err := service.SpecificationValues(context.Background(), principal, &specificationID, CatalogFilter{}, 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Data) != len(want) {
		t.Fatalf("values count = %d, want %d", len(page.Data), len(want))
	}
	for index, value := range page.Data {
		if value.ID != want[index] {
			t.Fatalf("value[%d] = %q, want %q", index, value.ID, want[index])
		}
	}
}
