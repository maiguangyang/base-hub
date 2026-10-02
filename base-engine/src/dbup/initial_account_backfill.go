package dbup

import (
	"base-engine/gen"
	"gorm.io/gorm"
)

// ReportUnresolvedFranchiseInitialAccounts reports legacy rows needing verified provenance.
// Historical provision audits do not identify the account chosen at creation time.
func ReportUnresolvedFranchiseInitialAccounts(db *gorm.DB) ([]string, error) {
	var ids []string
	err := db.Model(&gen.Organization{}).
		Where("type = ? AND initial_account_id IS NULL AND (is_delete IS NULL OR is_delete = 1)", gen.OrganizationTypeFranchise).
		Order("id").Pluck("id", &ids).Error
	return ids, err
}
