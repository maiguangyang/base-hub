package dbup

import (
	"fmt"

	"github.com/gofrs/uuid"
	"base-engine/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func EnsureCustomerCouponDistributionJobs(db *gorm.DB, nowMillis int64) error {
	if nowMillis < minimumCouponMillis || nowMillis > maximumCouponMillis {
		return fmt.Errorf("invalid coupon distribution seed time")
	}
	var templates []gen.CustomerCouponTemplate
	if err := db.Where("enabled = ?", true).Order("created_at, id").Find(&templates).Error; err != nil {
		return err
	}
	for _, template := range templates {
		availableAt := nowMillis
		if template.EffectiveAt > availableAt {
			availableAt = template.EffectiveAt
		}
		job := gen.CustomerCouponDistributionJob{
			ID: uuid.Must(uuid.NewV4()).String(), Kind: gen.CustomerCouponDistributionJobKindTemplateFanout,
			Status: gen.CustomerCouponDistributionJobStatusPending, RequestKey: "TEMPLATE_FANOUT:" + template.ID,
			AvailableAt: availableAt, TemplateID: &template.ID, CreatedAt: nowMillis,
		}
		if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "request_key"}}, DoNothing: true}).Create(&job).Error; err != nil {
			return err
		}
	}
	return nil
}
