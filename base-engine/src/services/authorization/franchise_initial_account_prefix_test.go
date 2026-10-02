package authorization

import (
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestInitialAccountCandidateSupportsTablePrefix(t *testing.T) {
	db := prefixedInitialAccountDB(t)
	mustCreateInitialCandidate(t, db, &gen.Account{ID: "account-a", Phone: "13800000001", DisplayName: "Owner", Status: gen.AccountStatusActive})
	mustCreateInitialCandidate(t, db, &gen.Organization{ID: "org-a", Code: "F001", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive})
	mustCreateInitialCandidate(t, db, &gen.OperatorMembership{ID: "member-a", AccountID: "account-a", OrganizationID: "org-a", Status: gen.MembershipStatusActive})
	if err := verifyInitialAccountCandidate(db, "org-a", "account-a"); err != nil {
		t.Fatalf("prefixed account candidate was rejected: %v", err)
	}
	mustCreateInitialCandidate(t, db, &gen.Organization{ID: "org-hq", Code: "HQ01", Name: "Headquarters", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive})
	mustCreateInitialCandidate(t, db, &gen.OperatorMembership{ID: "member-hq", AccountID: "account-a", OrganizationID: "org-hq", Status: gen.MembershipStatusActive})
	if code := auth.ErrorCode(verifyInitialAccountCandidate(db, "org-a", "account-a")); code != auth.CodePermissionDenied {
		t.Fatalf("HQ account candidate was allowed: %s", code)
	}
}

func prefixedInitialAccountDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true, NamingStrategy: schema.NamingStrategy{TablePrefix: "prefix_"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Account{}, &gen.Organization{}, &gen.OperatorMembership{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			if err := db.Exec("DROP INDEX IF EXISTS `" + name + "`").Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	return db
}

func mustCreateInitialCandidate(t *testing.T, db *gorm.DB, record any) {
	t.Helper()
	if err := db.Create(record).Error; err != nil {
		t.Fatal(err)
	}
}
