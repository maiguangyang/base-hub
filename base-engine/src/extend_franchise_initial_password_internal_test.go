package src

import (
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestFranchiseInitialAccountTargetPreservesDatabaseErrors(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	assertDatabaseError := func(stage string) {
		t.Helper()
		_, lookupErr := franchiseInitialAccountTarget(db, "franchise")
		if lookupErr == nil || auth.ErrorCode(lookupErr) != "" || !strings.Contains(lookupErr.Error(), "no such table") {
			t.Fatalf("%s query error = %v", stage, lookupErr)
		}
	}
	assertDatabaseError("organization")
	if err := db.AutoMigrate(&gen.Organization{}); err != nil {
		t.Fatal(err)
	}
	dropInitialAccountTestIndexes(t, db)
	accountID := "owner"
	if err := db.Create(&gen.Organization{ID: "franchise", Code: "F001", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive, InitialAccountID: &accountID}).Error; err != nil {
		t.Fatal(err)
	}
	assertDatabaseError("opening record")
	if err := db.AutoMigrate(&gen.FranchiseOpeningRecord{}); err != nil {
		t.Fatal(err)
	}
	dropInitialAccountTestIndexes(t, db)
	assertDatabaseError("account")
	if err := db.AutoMigrate(&gen.Account{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.Account{ID: accountID, Phone: "13800000000", DisplayName: "Owner", Status: gen.AccountStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	assertDatabaseError("membership")
}

func dropInitialAccountTestIndexes(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
		if err := db.Exec("DROP INDEX IF EXISTS `" + index + "`").Error; err != nil {
			t.Fatal(err)
		}
	}
}

func TestResetTargetPolicySupportsTablePrefix(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true, NamingStrategy: schema.NamingStrategy{TablePrefix: "prefix_"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Organization{}, &gen.OperatorMembership{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		dropInitialAccountTestIndexes(t, db)
	}
	if err := db.Create(&gen.Organization{ID: "org-a", Code: "F001", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.OperatorMembership{ID: "member-a", AccountID: "account-a", OrganizationID: "org-a", Status: gen.MembershipStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	if err := verifyResetTargetPolicy(db, "org-a", "account-a"); err != nil {
		t.Fatalf("prefixed reset target was rejected: %v", err)
	}
	if err := db.Create(&gen.Organization{ID: "org-hq", Code: "HQ01", Name: "Headquarters", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.OperatorMembership{ID: "member-hq", AccountID: "account-a", OrganizationID: "org-hq", Status: gen.MembershipStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	if code := auth.ErrorCode(verifyResetTargetPolicy(db, "org-a", "account-a")); code != auth.CodePermissionDenied {
		t.Fatalf("HQ account reset target was allowed: %s", code)
	}
}
