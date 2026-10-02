package customer

import (
	"context"
	"testing"
	"time"

	"base-engine/gen"
)

func TestCouponDistributionSchedulesFutureTemplateAtEffectiveTime(t *testing.T) {
	service, db, principal, _ := couponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	future := couponDistributionTestNow + 60_000
	input := validCouponInput()
	input.EffectiveAt = &future
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJob(t, db, templateCouponDistributionKey(template.ID), gen.CustomerCouponDistributionJobStatusPending, future)
}

func TestCouponDistributionStoreEnableSchedulesStableJob(t *testing.T) {
	service, db, _, principal, _ := storeCouponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	input := validCouponInput()
	input.Enabled = false
	template, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJobCount(t, db, templateCouponDistributionKey(template.ID), 0)
	if _, err := service.SetStoreCouponTemplateEnabled(context.Background(), principal, "store", template.ID, true); err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJob(t, db, templateCouponDistributionKey(template.ID), gen.CustomerCouponDistributionJobStatusPending, couponDistributionTestNow)
}

func TestCouponDistributionExhaustedTemplateEnableCompletesJob(t *testing.T) {
	service, db, principal, _ := couponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	input := validCouponInput()
	input.Enabled, input.TotalIssueLimit = false, 1
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(template).Update("issued_count", 1).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetCouponTemplateEnabled(context.Background(), principal, template.ID, true); err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJob(t, db, templateCouponDistributionKey(template.ID), gen.CustomerCouponDistributionJobStatusCompleted, couponDistributionTestNow)
}

func TestCouponDistributionMemberScheduleRollsBackWithAudit(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	if err := db.Exec(`CREATE TRIGGER fail_coupon_member_audit BEFORE INSERT ON audit_logs WHEN NEW.resource_type = 'customerMember' BEGIN SELECT RAISE(ABORT, 'forced audit failure'); END`).Error; err != nil {
		t.Fatal(err)
	}
	input := customerCreateInput("rollback-member")
	if _, err := service.CreateMember(context.Background(), principal, input); err == nil {
		t.Fatal("expected forced audit failure")
	}
	var members, jobs int64
	if err := db.Model(&gen.CustomerMember{}).Where("request_key = ?", input.RequestKey).Count(&members).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Count(&jobs).Error; err != nil {
		t.Fatal(err)
	}
	if members != 0 || jobs != 0 {
		t.Fatalf("rollback left members=%d jobs=%d", members, jobs)
	}
}

func TestCouponDistributionQueriesRemainReadOnly(t *testing.T) {
	service, db, principal, _ := couponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	if _, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput()); err != nil {
		t.Fatal(err)
	}
	if _, err := service.ListCouponTemplates(context.Background(), principal, 1, 20); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SearchMembers(context.Background(), principal, MemberSearch{Page: 1, PerPage: 20}); err != nil {
		t.Fatal(err)
	}
	var grants int64
	if err := db.Model(&gen.CustomerCouponGrant{}).Count(&grants).Error; err != nil || grants != 0 {
		t.Fatalf("queries created grants=%d err=%v", grants, err)
	}
	var pending int64
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("status = ?", gen.CustomerCouponDistributionJobStatusPending).Count(&pending).Error; err != nil || pending != 1 {
		t.Fatalf("queries mutated pending jobs=%d err=%v", pending, err)
	}
}
