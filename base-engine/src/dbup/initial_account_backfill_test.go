package dbup

import (
	"testing"

	"base-engine/gen"
	"gorm.io/gorm"
)

func TestReportUnresolvedFranchiseInitialAccountsDoesNotInferFromOrder(t *testing.T) {
	db := seedBackfillFixture(t)
	unresolved, err := ReportUnresolvedFranchiseInitialAccounts(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(unresolved) != 4 || unresolved[0] != "no-audit" || unresolved[1] != "pre-audit" || unresolved[2] != "provable" || unresolved[3] != "tied" {
		t.Fatalf("unresolved = %v", unresolved)
	}
	assertUnmarkedOrganizations(t, db)
	again, err := ReportUnresolvedFranchiseInitialAccounts(db)
	if err != nil || len(again) != 4 {
		t.Fatalf("repeat unresolved=%v err=%v", again, err)
	}
}

func TestReportUnresolvedFranchiseInitialAccountsExcludesDeletedOrganizations(t *testing.T) {
	db := seedBackfillFixture(t)
	deleted := int64(2)
	if err := db.Model(&gen.Organization{}).Where("id = ?", "pre-audit").Update("is_delete", deleted).Error; err != nil {
		t.Fatal(err)
	}
	unresolved, err := ReportUnresolvedFranchiseInitialAccounts(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(unresolved) != 3 || unresolved[0] != "no-audit" || unresolved[1] != "provable" || unresolved[2] != "tied" {
		t.Fatalf("unresolved = %v", unresolved)
	}
}

func seedBackfillFixture(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := migrateGeneratedForSQLite(db, &gen.Account{}, &gen.Organization{}, &gen.OperatorMembership{}, &gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	deleted := int64(2)
	accounts := []gen.Account{
		{ID: "first", Phone: "13800000101", DisplayName: "First", Status: gen.AccountStatusActive},
		{ID: "later", Phone: "13800000102", DisplayName: "Later", Status: gen.AccountStatusActive},
	}
	organizations := []gen.Organization{
		{ID: "provable", Code: "F001", Name: "Provable", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive},
		{ID: "tied", Code: "F002", Name: "Tied", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive},
		{ID: "no-audit", Code: "F003", Name: "No audit", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive},
		{ID: "pre-audit", Code: "F004", Name: "Pre audit", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive},
	}
	memberships := []gen.OperatorMembership{
		{ID: "earliest", AccountID: "first", OrganizationID: "provable", Status: gen.MembershipStatusActive, IsDelete: &deleted, CreatedAt: 100},
		{ID: "subsequent", AccountID: "later", OrganizationID: "provable", Status: gen.MembershipStatusActive, CreatedAt: 200},
		{ID: "tie-1", AccountID: "first", OrganizationID: "tied", Status: gen.MembershipStatusActive, CreatedAt: 100},
		{ID: "tie-2", AccountID: "later", OrganizationID: "tied", Status: gen.MembershipStatusActive, CreatedAt: 100},
		{ID: "unaudited", AccountID: "first", OrganizationID: "no-audit", Status: gen.MembershipStatusActive, CreatedAt: 100},
		{ID: "pre-1", AccountID: "first", OrganizationID: "pre-audit", Status: gen.MembershipStatusActive, CreatedAt: 100},
		{ID: "pre-2", AccountID: "later", OrganizationID: "pre-audit", Status: gen.MembershipStatusActive, CreatedAt: 120},
	}
	audits := []gen.AuditLog{
		{ID: "provision-1", Action: "franchise:provision", OrganizationID: pointer("provable"), ResourceID: pointer("provable"), ResourceType: "organization", ResultCode: "SUCCESS", CreatedAt: 150},
		{ID: "provision-2", Action: "franchise:provision", OrganizationID: pointer("tied"), ResourceID: pointer("tied"), ResourceType: "organization", ResultCode: "SUCCESS", CreatedAt: 150},
		{ID: "provision-3", Action: "franchise:provision", OrganizationID: pointer("pre-audit"), ResourceID: pointer("pre-audit"), ResourceType: "organization", ResultCode: "SUCCESS", CreatedAt: 150},
	}
	for _, records := range []any{accounts, organizations, memberships, audits} {
		if err := db.Create(records).Error; err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func assertUnmarkedOrganizations(t *testing.T, db *gorm.DB) {
	t.Helper()
	var provable, tied, noAudit, preAudit gen.Organization
	db.First(&provable, "id = ?", "provable")
	db.First(&tied, "id = ?", "tied")
	db.First(&noAudit, "id = ?", "no-audit")
	db.First(&preAudit, "id = ?", "pre-audit")
	if provable.InitialAccountID != nil || tied.InitialAccountID != nil || noAudit.InitialAccountID != nil || preAudit.InitialAccountID != nil {
		t.Fatalf("unsafe backfill result: %#v %#v %#v %#v", provable.InitialAccountID, tied.InitialAccountID, noAudit.InitialAccountID, preAudit.InitialAccountID)
	}
}

func pointer(value string) *string { return &value }
