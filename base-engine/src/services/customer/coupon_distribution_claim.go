package customer

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/gen"
	"gorm.io/gorm"
)

const couponDistributionLeaseDuration = 60 * time.Second

var errCouponDistributionLeaseRelinquished = errors.New("coupon distribution lease relinquished")

func (s *Service) claimNextCouponDistributionJob(ctx context.Context, nowMillis int64) (*gen.CustomerCouponDistributionJob, error) {
	for attempt := 0; attempt < 3; attempt++ {
		var job gen.CustomerCouponDistributionJob
		eligible := "(status = ? AND available_at <= ?) OR (status = ? AND lease_expires_at IS NOT NULL AND lease_expires_at <= ?)"
		err := s.db.WithContext(ctx).Where(eligible, gen.CustomerCouponDistributionJobStatusPending, nowMillis,
			gen.CustomerCouponDistributionJobStatusRunning, nowMillis).Order("available_at, created_at, id").First(&job).Error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		token := uuid.Must(uuid.NewV4()).String()
		leaseExpiresAt := nowMillis + couponDistributionLeaseDuration.Milliseconds()
		update := s.db.WithContext(ctx).Model(&gen.CustomerCouponDistributionJob{}).Where("id = ?", job.ID).
			Where(eligible, gen.CustomerCouponDistributionJobStatusPending, nowMillis, gen.CustomerCouponDistributionJobStatusRunning, nowMillis).
			Updates(map[string]any{"status": gen.CustomerCouponDistributionJobStatusRunning, "lease_token": token, "lease_expires_at": leaseExpiresAt})
		if update.Error != nil {
			return nil, update.Error
		}
		if update.RowsAffected == 1 {
			job.Status, job.LeaseToken, job.LeaseExpiresAt = gen.CustomerCouponDistributionJobStatusRunning, &token, &leaseExpiresAt
			return &job, nil
		}
	}
	return nil, nil
}

func renewCouponDistributionCursor(tx *gorm.DB, job *gen.CustomerCouponDistributionJob, candidate couponDistributionCandidate, nowMillis int64) error {
	leaseExpiresAt := nowMillis + couponDistributionLeaseDuration.Milliseconds()
	update := tx.Model(&gen.CustomerCouponDistributionJob{}).
		Where("id = ? AND status = ? AND lease_token = ?", job.ID, gen.CustomerCouponDistributionJobStatusRunning, job.LeaseToken).
		Updates(map[string]any{"cursor_created_at": candidate.CreatedAt, "cursor_key": candidate.ID, "lease_expires_at": leaseExpiresAt})
	if update.Error != nil {
		return update.Error
	}
	if update.RowsAffected != 1 {
		return errCouponDistributionLeaseRelinquished
	}
	job.CursorCreatedAt, job.CursorKey, job.LeaseExpiresAt = &candidate.CreatedAt, &candidate.ID, &leaseExpiresAt
	return nil
}

func (s *Service) finishCouponDistributionBatch(ctx context.Context, job *gen.CustomerCouponDistributionJob, complete bool, nowMillis int64) error {
	status := gen.CustomerCouponDistributionJobStatusPending
	if complete {
		status = gen.CustomerCouponDistributionJobStatusCompleted
	}
	update := s.db.WithContext(ctx).Model(&gen.CustomerCouponDistributionJob{}).
		Where("id = ? AND status = ? AND lease_token = ?", job.ID, gen.CustomerCouponDistributionJobStatusRunning, job.LeaseToken).
		Updates(map[string]any{"status": status, "available_at": nowMillis, "lease_token": nil, "lease_expires_at": nil, "last_error_code": nil})
	if update.Error != nil {
		return update.Error
	}
	if update.RowsAffected != 1 {
		return errCouponDistributionLeaseRelinquished
	}
	return nil
}

func (s *Service) releaseFailedCouponDistributionJob(ctx context.Context, job *gen.CustomerCouponDistributionJob, nowMillis int64, code string) {
	attempts := job.Attempts + 1
	delay := time.Second << min(attempts-1, 9)
	if delay > 10*time.Minute {
		delay = 10 * time.Minute
	}
	s.db.WithContext(ctx).Model(&gen.CustomerCouponDistributionJob{}).
		Where("id = ? AND status = ? AND lease_token = ?", job.ID, gen.CustomerCouponDistributionJobStatusRunning, job.LeaseToken).
		Updates(map[string]any{"status": gen.CustomerCouponDistributionJobStatusPending, "attempts": attempts,
			"last_error_code": code, "available_at": nowMillis + delay.Milliseconds(), "lease_token": nil, "lease_expires_at": nil})
}
