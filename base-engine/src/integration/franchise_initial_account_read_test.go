package integration_test

import (
	"strings"
	"testing"

	"base-engine/gen"
)

func TestFranchiseListSurvivesDeletedInitialAccount(t *testing.T) {
	fixture := newSecurityFixture(t)
	fixture.db.Model(&gen.Organization{}).Where("id = ?", "org-a").Update("initial_account_id", "account-staff")
	fixture.db.Model(&gen.Account{}).Where("id = ?", "account-staff").Update("is_delete", 2)
	response := fixture.execute("session-hq", `query { organizations(filter: {type: FRANCHISE}) { data { id initialAccountId initialAccount { id } } } }`)
	if len(response.Errors) != 0 || !strings.Contains(response.Body, `"initialAccountId":"account-staff"`) || !strings.Contains(response.Body, `"initialAccount":null`) {
		t.Fatalf("deleted account broke franchise list: %s", response.Body)
	}
}
