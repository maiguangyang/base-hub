package productcatalog

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCategoryMySQLConcurrentMovesCannotCreateCycle(t *testing.T) {
	service, db, principal := categoryMySQLFixture(t)
	left, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: "Left"})
	catalogNoError(t, err)
	right, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: "Right"})
	catalogNoError(t, err)
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, move := range []func() error{
		func() error { return service.MoveCategory(context.Background(), principal, left.ID, &right.ID) },
		func() error {
			_, err := service.UpdateCategory(context.Background(), principal, right.ID, right.Name, &left.ID)
			return err
		},
	} {
		wg.Add(1)
		go func(move func() error) { defer wg.Done(); <-start; results <- move() }(move)
	}
	close(start)
	wg.Wait()
	close(results)
	assertNoCategoryCycle(t, db, left.ID, right.ID, results)
}

func categoryMySQLFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	rootDSN := os.Getenv("PRODUCT_TEST_MYSQL_DSN")
	if rootDSN == "" {
		t.Skip("PRODUCT_TEST_MYSQL_DSN is required for the disposable MySQL test")
	}
	if !strings.HasSuffix(rootDSN, "/") {
		t.Fatal("PRODUCT_TEST_MYSQL_DSN must end in / with no database selected")
	}
	admin, err := gorm.Open(mysql.Open(rootDSN+"mysql?parseTime=true"), &gorm.Config{})
	catalogNoError(t, err)
	name := fmt.Sprintf("product_category_test_%d", time.Now().UnixNano())
	catalogNoError(t, admin.Exec("CREATE DATABASE `"+name+"`").Error)
	t.Cleanup(func() { catalogNoError(t, admin.Exec("DROP DATABASE `"+name+"`").Error) })
	db, err := gorm.Open(mysql.Open(rootDSN+name+"?parseTime=true"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true})
	catalogNoError(t, err)
	for _, model := range []any{&gen.Organization{}, &gen.ProductCategory{}, &gen.AuditLog{}} {
		catalogNoError(t, db.AutoMigrate(model))
	}
	catalogNoError(t, db.Create(&gen.Organization{ID: "hq", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}).Error)
	hqID := "hq"
	principal := &auth.WorkspacePrincipal{AccountID: "admin", SessionID: "session", WorkspaceType: auth.WorkspaceTypeHeadquarters,
		OrganizationID: &hqID, Permissions: map[string]struct{}{"hqProductCatalog:manage": {}}}
	return NewService(db, audit.NewService()), db, principal
}

func assertNoCategoryCycle(t *testing.T, db *gorm.DB, leftID, rightID string, results <-chan error) {
	t.Helper()
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("successful opposing moves = %d, want 1", successes)
	}
	var categories []gen.ProductCategory
	catalogNoError(t, db.Where("id IN ?", []string{leftID, rightID}).Find(&categories).Error)
	if len(categories) != 2 || categories[0].ParentID != nil && categories[1].ParentID != nil {
		t.Fatalf("category cycle: %+v", categories)
	}
}
