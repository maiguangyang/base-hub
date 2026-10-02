package integration_test

import (
	"encoding/json"
	"errors"
	"testing"

	"base-engine/gen"
	"gorm.io/gorm"
)

func TestCustomerCancellationRollsBackWhenRightsCountFails(t *testing.T) {
	fixture := newCustomerFixture(t)
	create := fixture.execute("session-hq", `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"create-1"}) { id } }`)
	assertCustomerSuccess(t, create)
	var member struct{ ID string }
	if err := json.Unmarshal(create.Data["hqCreateCustomerMember"], &member); err != nil {
		t.Fatal(err)
	}
	if err := fixture.db.Callback().Query().Before("gorm:query").Register("fail_customer_rights_count", func(tx *gorm.DB) {
		if tx.Statement.Table == "customer_coupon_grants" {
			tx.AddError(errors.New("injected rights count failure"))
		}
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = fixture.db.Callback().Query().Remove("fail_customer_rights_count") })
	response := fixture.execute("session-hq", `mutation { hqRequestCustomerCancellation(id:"`+member.ID+`", identityEvidence:"proof", basisCode:"REQUEST") { status pendingCouponCount } }`)
	if len(response.Errors) == 0 {
		t.Fatalf("expected rights count failure: %s", response.Body)
	}
	var persisted gen.CustomerMember
	if err := fixture.db.First(&persisted, "id = ?", member.ID).Error; err != nil || persisted.Status != gen.CustomerMemberStatusActive {
		t.Fatalf("cancellation committed despite failed response: %+v, %v", persisted, err)
	}
}

func TestCustomerPointGrantRollsBackWhenResponseBalanceFails(t *testing.T) {
	fixture := newCustomerFixture(t)
	fixture.createAll([]gen.CustomerBenefitPolicy{{
		ID: "policy-hq", OrganizationID: "org-hq", Version: 1,
		ManualGrantMaxSingle: 100, ManualGrantMaxDaily: 100,
	}})
	create := fixture.execute("session-hq", `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"create-1"}) { id } }`)
	assertCustomerSuccess(t, create)
	var member struct{ ID string }
	if err := json.Unmarshal(create.Data["hqCreateCustomerMember"], &member); err != nil {
		t.Fatal(err)
	}
	removeFailure := installPointBalanceFailure(t, fixture)
	response := fixture.execute("session-hq", `mutation { hqGrantCustomerPoints(memberId:"`+member.ID+`", points:50, reasonCode:REWARD, note:"gift", requestKey:"gift-1") { currentBalance } }`)
	removeFailure()
	if len(response.Errors) == 0 {
		t.Fatalf("expected balance read failure: %s", response.Body)
	}
	var persisted gen.CustomerMember
	if err := fixture.db.First(&persisted, "id = ?", member.ID).Error; err != nil || persisted.PointsBalance != 0 {
		t.Fatalf("point grant committed despite failed response: %+v, %v", persisted, err)
	}
	var entries int64
	if err := fixture.db.Model(&gen.CustomerPointEntry{}).Count(&entries).Error; err != nil || entries != 0 {
		t.Fatalf("point entries committed despite failed response: %d, %v", entries, err)
	}
}

func installPointBalanceFailure(t *testing.T, fixture *securityFixture) func() {
	t.Helper()
	pointCreated := false
	if err := fixture.db.Callback().Create().After("gorm:create").Register("mark_customer_point_created", func(tx *gorm.DB) {
		if tx.Statement.Table == "customer_point_entries" {
			pointCreated = true
		}
	}); err != nil {
		t.Fatal(err)
	}
	if err := fixture.db.Callback().Query().Before("gorm:query").Register("fail_response_balance", func(tx *gorm.DB) {
		if pointCreated && tx.Statement.Table == "customer_members" {
			tx.AddError(errors.New("injected response balance failure"))
		}
	}); err != nil {
		t.Fatal(err)
	}
	return func() {
		_ = fixture.db.Callback().Create().Remove("mark_customer_point_created")
		_ = fixture.db.Callback().Query().Remove("fail_response_balance")
	}
}
