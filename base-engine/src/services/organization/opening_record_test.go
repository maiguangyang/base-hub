package organization

import (
	"context"
	"testing"

	"base-engine/gen"
)

func TestProvisionFranchiseWritesOpeningRecordForSharedAccount(t *testing.T) {
	service, db := newOrganizationFixture(t)
	principal := headquartersPrincipal()
	for _, code := range []string{"F001", "F002"} {
		result, err := service.ProvisionFranchise(context.Background(), principal, ProvisionInput{Code: code, Name: code, OwnerPhone: "13800000000", OwnerDisplayName: "Owner"})
		if err != nil {
			t.Fatal(err)
		}
		var record gen.FranchiseOpeningRecord
		if err := db.Where("organization_id = ?", result.Organization.ID).First(&record).Error; err != nil {
			t.Fatal(err)
		}
		if record.RecordNumber == "" || record.Source != gen.FranchiseOpeningSourceSystemProvision || record.InitialAccountID != result.Membership.AccountID || record.RecordedByAccountID != principal.AccountID {
			t.Fatalf("incorrect opening record: %#v", record)
		}
	}
}
