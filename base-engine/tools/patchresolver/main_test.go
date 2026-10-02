package main

import (
	"strings"
	"testing"
)

func TestPatchInitialAccountProjectionTargetsOrganization(t *testing.T) {
	source := `type OrganizationResolver interface {
	Memberships(ctx context.Context, obj *Organization)
}
func (ec *executionContext) _FranchiseOpeningRecord_initialAccountId() {
	return obj.InitialAccountID, nil
}
func (ec *executionContext) _Organization_initialAccountId() {
	return ec.fieldContext_Organization_initialAccountId(ctx, field)
		},
		func(ctx context.Context) (any, error) {
			return obj.InitialAccountID, nil
	}
}`
	patched, err := patchInitialAccountProjection(source)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(patched, "return ec.Resolvers.Organization().InitialAccountID(ctx, obj)") != 1 ||
		strings.Count(patched, "return obj.InitialAccountID, nil") != 1 {
		t.Fatalf("projection patch targeted wrong resolver: %s", patched)
	}
}

func TestPatchSpecificationOrderDocsUsesNonNullIDs(t *testing.T) {
	source := `mutation hqReorderSpecificationValues($specificationId: ID!, $orderedIds: [ID]!) {
  hqReorderSpecificationValues(specificationId: $specificationId, orderedIds: $orderedIds)
}`
	patched, err := patchSpecificationOrderDocs(source)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(patched, `$orderedIds: [ID!]!`) {
		t.Fatalf("generated example still has nullable IDs: %s", patched)
	}
	again, err := patchSpecificationOrderDocs(patched)
	if err != nil || again != patched {
		t.Fatalf("documentation patch is not repeatable: %v", err)
	}
}

func TestPatchCustomerCouponDistributionJobIDProjectionsDeniesBeforeLoaders(t *testing.T) {
	source := `func (r *GeneratedCustomerMemberResolver) CouponDistributionJobsIds(ctx context.Context, obj *CustomerMember) (ids []string, err error) {
	loaders := ctx.Value(KeyLoaders)
}
func (r *GeneratedCustomerCouponTemplateResolver) DistributionJobsIds(ctx context.Context, obj *CustomerCouponTemplate) (ids []string, err error) {
	loaders := ctx.Value(KeyLoaders)
}`
	patched, err := patchCustomerCouponDistributionJobIDProjections(source)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(patched, "denied := auth.NewError(auth.CodePermissionDenied)") != 2 {
		t.Fatalf("missing projection guards: %s", patched)
	}
	if strings.Index(patched, "auth.CodePermissionDenied") > strings.Index(patched, "ctx.Value(KeyLoaders)") {
		t.Fatalf("member projection reaches loaders before denial: %s", patched)
	}
	again, err := patchCustomerCouponDistributionJobIDProjections(patched)
	if err != nil || again != patched {
		t.Fatalf("projection patch is not repeatable: %v", err)
	}
}
