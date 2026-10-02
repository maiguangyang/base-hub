package customer

import (
	"context"
	"errors"

	mysql "github.com/go-sql-driver/mysql"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func mysqlErrorNumber(err error) uint16 {
	var dbError *mysql.MySQLError
	if errors.As(err, &dbError) {
		return dbError.Number
	}
	return 0
}

type PointMutationResult struct {
	*gen.CustomerPointEntry
	CurrentBalance int64
}

func pointMutationResult(tx *gorm.DB, entry *gen.CustomerPointEntry) (*PointMutationResult, error) {
	var member gen.CustomerMember
	err := tx.Select("points_balance").Where("id = ? AND organization_id = ?", entry.MemberID, entry.SourceOrganizationID).First(&member).Error
	if err != nil {
		return nil, err
	}
	if member.PointsBalance < 0 || member.PointsBalance > 2_147_483_647 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	return &PointMutationResult{CustomerPointEntry: entry, CurrentBalance: member.PointsBalance}, nil
}

func (s *Service) runPointWrite(ctx context.Context, write func(*gorm.DB) (*gen.CustomerPointEntry, error), replay func() (*gen.CustomerPointEntry, error)) (*PointMutationResult, error) {
	for attempt := 0; attempt < 3; attempt++ {
		var result *PointMutationResult
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			entry, err := write(tx)
			if err != nil {
				return err
			}
			result, err = pointMutationResult(tx, entry)
			return err
		})
		if err == nil {
			return result, nil
		}
		if mysqlErrorNumber(err) == 1062 || auth.ErrorCode(err) == auth.CodeConflict {
			entry, replayErr := replay()
			if replayErr != nil {
				return nil, replayErr
			}
			return pointMutationResult(s.db.WithContext(ctx), entry)
		}
		if mysqlErrorNumber(err) != 1213 && mysqlErrorNumber(err) != 1205 {
			return nil, err
		}
	}
	return nil, auth.NewError(auth.CodeConflict)
}
