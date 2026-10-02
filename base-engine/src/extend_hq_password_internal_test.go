package src

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPasswordResetAuthorizesTargetBeforeGeneratingPassword(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&gen.OperatorMembership{}); err != nil {
		t.Fatal(err)
	}
	organizationID := "hq-org"
	principal := &auth.WorkspacePrincipal{AccountID: "actor", OrganizationID: &organizationID, WorkspaceType: auth.WorkspaceTypeHeadquarters}
	generated := false
	_, err = resetTemporaryPasswordWithGenerator(context.Background(), Dependencies{DB: database}, principal, "missing", func() (string, error) {
		generated = true
		return "temporary", nil
	})
	if auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("authorization error = %v", err)
	}
	if generated {
		t.Fatal("temporary password was generated before target authorization")
	}
}
