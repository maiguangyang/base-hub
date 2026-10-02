package customer

import (
	"context"
	"errors"
	"log"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

const couponDistributionBatchSize = 100

type CouponDistributionBatchResult struct {
	Claimed, Completed bool
	JobID, Kind        string
	Issued, Skipped    int
}

type couponDistributionCandidate struct {
	ID        string
	CreatedAt int64
}

func (s *Service) ProcessNextCouponDistributionBatch(ctx context.Context) (CouponDistributionBatchResult, error) {
	started := time.Now()
	nowMillis := s.now().UnixMilli()
	job, err := s.claimNextCouponDistributionJob(ctx, nowMillis)
	if err != nil || job == nil {
		return CouponDistributionBatchResult{}, err
	}
	result := CouponDistributionBatchResult{Claimed: true, JobID: job.ID, Kind: string(job.Kind)}
	candidates, complete, err := s.couponDistributionCandidates(ctx, job, nowMillis)
	if err == nil {
		result.Issued, result.Skipped, err = s.processCouponDistributionCandidates(ctx, job, candidates)
	}
	nowMillis = s.now().UnixMilli()
	if err == nil && !complete {
		complete, err = s.couponDistributionTargetComplete(ctx, job, nowMillis)
	}
	if err == nil {
		err = s.finishCouponDistributionBatch(ctx, job, complete, nowMillis)
		result.Completed = complete
	}
	if err != nil && !errors.Is(err, errCouponDistributionLeaseRelinquished) {
		s.releaseFailedCouponDistributionJob(ctx, job, nowMillis, "COUPON_DISTRIBUTION_FAILED")
	}
	log.Printf("COUPON_DISTRIBUTION_BATCH job_id=%s kind=%s issued=%d skipped=%d completed=%t attempts=%d duration_ms=%d error_code=%s",
		job.ID, job.Kind, result.Issued, result.Skipped, result.Completed, job.Attempts, time.Since(started).Milliseconds(), couponDistributionErrorCode(err))
	return result, err
}

func couponDistributionErrorCode(err error) string {
	if err == nil {
		return "NONE"
	}
	if errors.Is(err, errCouponDistributionLeaseRelinquished) {
		return "LEASE_RELINQUISHED"
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return "CONTEXT_ENDED"
	}
	return "COUPON_DISTRIBUTION_FAILED"
}

func (s *Service) processCouponDistributionCandidates(ctx context.Context, job *gen.CustomerCouponDistributionJob, candidates []couponDistributionCandidate) (int, int, error) {
	issued, skipped := 0, 0
	for _, candidate := range candidates {
		created, err := s.processCouponDistributionCandidate(ctx, job, candidate, s.now().UnixMilli())
		if err != nil {
			return issued, skipped, err
		}
		if created {
			issued++
		} else {
			skipped++
		}
	}
	return issued, skipped, nil
}

func (s *Service) couponDistributionTargetComplete(ctx context.Context, job *gen.CustomerCouponDistributionJob, nowMillis int64) (bool, error) {
	if job.Kind == gen.CustomerCouponDistributionJobKindTemplateFanout {
		_, complete, err := s.fanoutTemplate(ctx, job, nowMillis)
		return complete, err
	}
	_, complete, err := s.catchupMember(ctx, job)
	return complete, err
}

func (s *Service) couponDistributionCandidates(ctx context.Context, job *gen.CustomerCouponDistributionJob, nowMillis int64) ([]couponDistributionCandidate, bool, error) {
	query := s.db.WithContext(ctx)
	if job.CursorCreatedAt != nil && job.CursorKey != nil {
		query = query.Where("created_at > ? OR (created_at = ? AND id > ?)", *job.CursorCreatedAt, *job.CursorCreatedAt, *job.CursorKey)
	}
	var candidates []couponDistributionCandidate
	var err error
	switch job.Kind {
	case gen.CustomerCouponDistributionJobKindTemplateFanout:
		template, complete, targetErr := s.fanoutTemplate(ctx, job, nowMillis)
		if targetErr != nil || complete {
			return nil, complete, targetErr
		}
		err = query.Model(&gen.CustomerMember{}).
			Where("organization_id = ? AND status = ?", template.OrganizationID, gen.CustomerMemberStatusActive).
			Order("created_at, id").Limit(couponDistributionBatchSize).Find(&candidates).Error
	case gen.CustomerCouponDistributionJobKindMemberCatchup:
		member, complete, targetErr := s.catchupMember(ctx, job)
		if targetErr != nil || complete {
			return nil, complete, targetErr
		}
		err = query.Model(&gen.CustomerCouponTemplate{}).
			Where("organization_id = ? AND enabled = ? AND effective_at <= ?", member.OrganizationID, true, nowMillis).
			Where("distribution_ends_at IS NULL OR distribution_ends_at > ?", nowMillis).
			Where("total_issue_limit = 0 OR issued_count < total_issue_limit").
			Order("created_at, id").Limit(couponDistributionBatchSize).Find(&candidates).Error
	default:
		return nil, false, auth.NewError(auth.CodeInternalError)
	}
	return candidates, len(candidates) < couponDistributionBatchSize, err
}

func (s *Service) fanoutTemplate(ctx context.Context, job *gen.CustomerCouponDistributionJob, nowMillis int64) (*gen.CustomerCouponTemplate, bool, error) {
	if job.TemplateID == nil {
		return nil, false, auth.NewError(auth.CodeInternalError)
	}
	var template gen.CustomerCouponTemplate
	if err := s.db.WithContext(ctx).Where("id = ?", *job.TemplateID).First(&template).Error; err != nil {
		return nil, false, err
	}
	complete := !couponTemplateDistributable(&template, nowMillis) ||
		template.TotalIssueLimit > 0 && template.IssuedCount >= template.TotalIssueLimit
	return &template, complete, nil
}

func (s *Service) catchupMember(ctx context.Context, job *gen.CustomerCouponDistributionJob) (*gen.CustomerMember, bool, error) {
	if job.MemberID == nil {
		return nil, false, auth.NewError(auth.CodeInternalError)
	}
	var member gen.CustomerMember
	if err := s.db.WithContext(ctx).Where("id = ?", *job.MemberID).First(&member).Error; err != nil {
		return nil, false, err
	}
	return &member, member.Status != gen.CustomerMemberStatusActive, nil
}

func (s *Service) processCouponDistributionCandidate(ctx context.Context, job *gen.CustomerCouponDistributionJob, candidate couponDistributionCandidate, nowMillis int64) (bool, error) {
	created := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		input, err := couponDistributionGrantInput(tx, job, candidate.ID, nowMillis)
		if err != nil {
			return err
		}
		_, created, err = s.createCouponGrantInTransaction(tx, input)
		if auth.ErrorCode(err) == auth.CodeConflict {
			created, err = false, nil
		}
		if err != nil {
			return err
		}
		return renewCouponDistributionCursor(tx, job, candidate, s.now().UnixMilli())
	})
	if mysqlErrorNumber(err) == 1062 {
		input, inputErr := couponDistributionGrantInput(s.db.WithContext(ctx), job, candidate.ID, nowMillis)
		if inputErr != nil {
			return false, inputErr
		}
		if _, _, duplicateErr := s.resolveAutomaticCouponGrantDuplicate(ctx, err, input.TemplateID, input.MemberID); duplicateErr != nil {
			return false, duplicateErr
		}
		return false, s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return renewCouponDistributionCursor(tx, job, candidate, s.now().UnixMilli())
		})
	}
	return created, err
}

func couponDistributionGrantInput(tx *gorm.DB, job *gen.CustomerCouponDistributionJob, candidateID string, nowMillis int64) (couponGrantInput, error) {
	templateID, memberID := candidateID, candidateID
	if job.Kind == gen.CustomerCouponDistributionJobKindTemplateFanout {
		if job.TemplateID == nil {
			return couponGrantInput{}, auth.NewError(auth.CodeInternalError)
		}
		templateID = *job.TemplateID
	} else if job.MemberID != nil {
		memberID = *job.MemberID
	} else {
		return couponGrantInput{}, auth.NewError(auth.CodeInternalError)
	}
	var template gen.CustomerCouponTemplate
	if err := tx.Where("id = ?", templateID).First(&template).Error; err != nil {
		return couponGrantInput{}, err
	}
	input := couponGrantInput{TemplateID: templateID, MemberID: memberID, RequestKey: autoCouponRequestKey(templateID, memberID),
		ExpectedOrganizationID: template.OrganizationID, Automatic: true,
		EnforceMemberPriority: job.Kind == gen.CustomerCouponDistributionJobKindMemberCatchup, NowMillis: nowMillis}
	if template.IssuerScope == gen.CouponIssuerScopeStore {
		input.ExpectedStoreID = template.ApplicableStoreID
	}
	return input, nil
}
