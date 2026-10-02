package customer

import (
	"context"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func pointsFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal, *gen.CustomerMember) {
	t.Helper()
	service, db, principal := newCustomerServiceFixture(t)
	for _, action := range []string{"hqCustomerPoints:read", "hqCustomerPoints:grant", "hqCustomerPoints:reverse", "hqCustomerPoints:correct"} {
		principal.Permissions[action] = struct{}{}
	}
	member, err := service.CreateMember(context.Background(), principal, customerCreateInput("member-one"))
	if err != nil {
		t.Fatal(err)
	}
	return service, db, principal, member
}

func configurePoints(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, single, daily int64) {
	t.Helper()
	input := validCustomerPolicy()
	input.ManualGrantMaxSingle = single
	input.ManualGrantMaxDaily = daily
	if _, err := service.SavePolicy(context.Background(), principal, 0, input); err != nil {
		t.Fatal(err)
	}
}

func TestCustomerPointGrantLimits(t *testing.T) {
	service, db, principal, member := pointsFixture(t)
	ctx := context.Background()
	if _, err := service.GrantPoints(ctx, principal, member.ID, 1, gen.CustomerPointReasonCodeReward, "gift", "pre-policy"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("unconfigured gift = %v", err)
	}
	configurePoints(t, service, principal, 100, 150)
	assertInvalidPointAmounts(t, service, principal, member.ID)
	if _, err := service.GrantPoints(ctx, principal, member.ID, 100, gen.CustomerPointReasonCodeReward, "gift", "one"); err != nil {
		t.Fatal(err)
	}
	input := customerCreateInput("member-two")
	input.Phone = "13900000002"
	other, err := service.CreateMember(ctx, principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantPoints(ctx, principal, other.ID, 51, gen.CustomerPointReasonCodeReward, "gift", "two"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("shared daily cap = %v", err)
	}
	if _, err := service.GrantPoints(ctx, principal, other.ID, 50, gen.CustomerPointReasonCodeReward, "gift", "three"); err != nil {
		t.Fatal(err)
	}
	var budget gen.CustomerDailyPointGrantBudget
	if err := db.First(&budget).Error; err != nil || budget.UsedPoints != 150 {
		t.Fatalf("daily used = %d, %v", budget.UsedPoints, err)
	}
}

func assertInvalidPointAmounts(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, memberID string) {
	t.Helper()
	ctx := context.Background()
	for _, amount := range []int64{0, -1, 1 << 31} {
		if _, err := service.GrantPoints(ctx, principal, memberID, amount, gen.CustomerPointReasonCodeReward, "gift", "bad"); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("gift %d = %v", amount, err)
		}
	}
	if _, err := service.GrantPoints(ctx, principal, memberID, 101, gen.CustomerPointReasonCodeReward, "gift", "too-much"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("over single cap = %v", err)
	}
}

func TestCustomerPointIdempotencyAndReversal(t *testing.T) {
	service, db, principal, member := pointsFixture(t)
	configurePoints(t, service, principal, 100, 100)
	ctx := context.Background()
	first, err := service.GrantPoints(ctx, principal, member.ID, 60, gen.CustomerPointReasonCodeCompensation, "adjust", "grant-1")
	if err != nil {
		t.Fatal(err)
	}
	assertCreatedTime(t, first.CreatedAt)
	replayed, err := service.GrantPoints(ctx, principal, member.ID, 60, gen.CustomerPointReasonCodeCompensation, "adjust", "grant-1")
	if err != nil || replayed.ID != first.ID {
		t.Fatalf("grant replay = %+v, %v", replayed, err)
	}
	reversed, err := service.ReversePoints(ctx, principal, first.ID, gen.CustomerPointReasonCodeCorrection, "undo", "grant-1")
	if err != nil || reversed.Delta != -60 {
		t.Fatalf("reverse = %+v, %v", reversed, err)
	}
	retry, err := service.ReversePoints(ctx, principal, first.ID, gen.CustomerPointReasonCodeCorrection, "undo", "grant-1")
	if err != nil || retry.ID != reversed.ID {
		t.Fatalf("reverse replay = %+v, %v", retry, err)
	}
	if _, err := service.ReversePoints(ctx, principal, first.ID, gen.CustomerPointReasonCodeCorrection, "different", "grant-1"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed reverse intent = %v", err)
	}
	if _, err := service.ReversePoints(ctx, principal, first.ID, gen.CustomerPointReasonCodeCorrection, "undo", "another"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("duplicate reversal = %v", err)
	}
	assertBudgetNotRestored(t, db)
}

func assertBudgetNotRestored(t *testing.T, db *gorm.DB) {
	t.Helper()
	var budget gen.CustomerDailyPointGrantBudget
	if err := db.First(&budget).Error; err != nil || budget.UsedPoints != 60 {
		t.Fatalf("reversal restored budget: %d, %v", budget.UsedPoints, err)
	}
}

func TestCustomerPointFrozenMember(t *testing.T) {
	service, db, principal, member := pointsFixture(t)
	configurePoints(t, service, principal, 100, 100)
	service.now = func() time.Time { return time.Date(2026, 9, 25, 23, 59, 0, 0, time.FixedZone("CST", 8*3600)) }
	if err := db.Model(member).Updates(map[string]any{"status": gen.CustomerMemberStatusCancelPending}).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantPoints(context.Background(), principal, member.ID, 1, gen.CustomerPointReasonCodeReward, "gift", "pending"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("pending gift = %v", err)
	}
}

func TestCustomerPointDailyBudgetShanghaiBoundary(t *testing.T) {
	service, db, principal, member := pointsFixture(t)
	configurePoints(t, service, principal, 100, 100)
	service.now = func() time.Time { return time.Date(2026, 9, 25, 23, 59, 0, 0, shanghaiLocation) }
	if _, err := service.GrantPoints(context.Background(), principal, member.ID, 100, gen.CustomerPointReasonCodeReward, "gift", "day-one"); err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Date(2026, 9, 26, 0, 1, 0, 0, shanghaiLocation) }
	if _, err := service.GrantPoints(context.Background(), principal, member.ID, 100, gen.CustomerPointReasonCodeReward, "gift", "day-two"); err != nil {
		t.Fatal(err)
	}
	var budgets []gen.CustomerDailyPointGrantBudget
	if err := db.Order("business_date").Find(&budgets).Error; err != nil {
		t.Fatal(err)
	}
	if len(budgets) != 2 || budgets[0].UsedPoints != 100 || budgets[1].UsedPoints != 100 {
		t.Fatalf("daily budgets = %+v", budgets)
	}
}
