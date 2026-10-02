package customer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCustomerMySQLAutomaticDuplicateRequiresExactReplay(t *testing.T) {
	db := newCustomerCouponMySQLTestDB(t)
	service := NewService(db, audit.NewService())
	nowMillis := time.Now().UnixMilli()
	unrelated := gen.CustomerCouponGrant{ID: "mysql-duplicate-source", TemplateID: "other-template", MemberID: "other-member",
		RequestKey: "other-request", Status: gen.CustomerCouponGrantStatusPendingActivation, AmountFen: 1,
		DaysAfterActivation: 1, IssuedAt: nowMillis, CreatedAt: nowMillis}
	if err := db.Create(&unrelated).Error; err != nil {
		t.Fatal(err)
	}
	duplicate := unrelated
	duplicate.RequestKey = "colliding-primary-key"
	duplicateErr := db.Create(&duplicate).Error
	var mysqlErr *mysqldriver.MySQLError
	if !errors.As(duplicateErr, &mysqlErr) || mysqlErr.Number != 1062 {
		t.Fatalf("expected real MySQL 1062, got %v", duplicateErr)
	}

	const templateID, memberID = "target-template", "target-member"
	if _, _, err := service.resolveAutomaticCouponGrantDuplicate(context.Background(), duplicateErr, templateID, memberID); !errors.Is(err, duplicateErr) {
		t.Fatalf("unrelated MySQL duplicate was swallowed: %v", err)
	}
	exact := gen.CustomerCouponGrant{ID: "mysql-exact-auto", TemplateID: templateID, MemberID: memberID,
		RequestKey: autoCouponRequestKey(templateID, memberID), Status: gen.CustomerCouponGrantStatusPendingActivation,
		AmountFen: 100, DaysAfterActivation: 30, IssuedAt: nowMillis, CreatedAt: nowMillis}
	if err := db.Create(&exact).Error; err != nil {
		t.Fatal(err)
	}
	replay, created, err := service.resolveAutomaticCouponGrantDuplicate(context.Background(), duplicateErr, templateID, memberID)
	if err != nil || created || replay == nil || replay.ID != exact.ID {
		t.Fatalf("exact MySQL duplicate replay=%+v created=%t err=%v", replay, created, err)
	}
}

func newCustomerCouponMySQLTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	rootDSN := os.Getenv("KOREAN_MYSQL_TEST_DSN")
	if rootDSN == "" {
		t.Skip("set KOREAN_MYSQL_TEST_DSN to run coupon MySQL tests")
	}
	if !strings.HasSuffix(rootDSN, "/") {
		t.Fatal("KOREAN_MYSQL_TEST_DSN must end with /")
	}
	admin, err := gorm.Open(mysql.Open(rootDSN+"mysql?parseTime=true"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	database := fmt.Sprintf("korean_customer_coupon_test_%d", time.Now().UnixNano())
	if err := admin.Exec("CREATE DATABASE `" + database + "`").Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP DATABASE `" + database + "`").Error; err != nil {
			t.Errorf("drop disposable coupon database: %v", err)
		}
	})
	db, err := gorm.Open(mysql.Open(rootDSN+database+"?parseTime=true"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&gen.CustomerCouponGrant{}); err != nil {
		t.Fatal(err)
	}
	return db
}
