package authorization

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SetFranchiseInitialAccount records an HQ operator's one-time evidence-backed choice for an unmarked historical franchise.
func SetFranchiseInitialAccount(db *gorm.DB, principal *auth.WorkspacePrincipal, organizationID, accountID, evidenceReference string) error {
	if err := authorizeInitialAccountWriter(principal); err != nil {
		return err
	}
	if !validInitialAccountCommand(organizationID, accountID, evidenceReference) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		return setFranchiseInitialAccountTransaction(tx, principal, organizationID, accountID, evidenceReference)
	})
	if err != nil && auth.ErrorCode(err) == "" {
		if conflict := openingRecordCollision(db, organizationID, evidenceReference); conflict != nil {
			return conflict
		}
	}
	return err
}

func setFranchiseInitialAccountTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, organizationID, accountID, evidenceReference string) error {
	var organization gen.Organization
	if err := activeRecordCondition(tx.Clauses(clause.Locking{Strength: "UPDATE"})).First(&organization, "id = ? AND type = ?", organizationID, gen.OrganizationTypeFranchise).Error; err != nil {
		return concealInitialAccountWriteError(err)
	}
	if organization.InitialAccountID != nil {
		return auth.NewError(auth.CodeConflict)
	}
	if err := verifyInitialAccountCandidate(tx, organizationID, accountID); err != nil {
		return err
	}
	if err := recordHistoricalOpening(tx, principal.AccountID, organizationID, accountID, evidenceReference); err != nil {
		return err
	}
	if err := assignInitialAccount(tx, organizationID, accountID); err != nil {
		return err
	}
	return audit.NewService().Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID, OrganizationID: &organizationID,
		Action: "franchiseInitialAccount:confirm", ResourceType: "account", ResourceID: accountID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{Source: "historical_manual_attestation", EvidenceReference: evidenceReference}),
	})
}

func recordHistoricalOpening(tx *gorm.DB, actorID, organizationID, accountID, recordNumber string) error {
	if conflict := openingRecordCollision(tx, organizationID, recordNumber); conflict != nil {
		return conflict
	}
	record := gen.FranchiseOpeningRecord{
		ID: uuid.Must(uuid.NewV4()).String(), RecordNumber: recordNumber,
		Source:         gen.FranchiseOpeningSourceHistoricalAttestation,
		OrganizationID: organizationID, InitialAccountID: accountID, RecordedByAccountID: actorID,
	}
	return tx.Create(&record).Error
}

func openingRecordCollision(db *gorm.DB, organizationID, recordNumber string) error {
	for _, candidate := range []struct {
		where string
		value string
		code  auth.Code
	}{{"organization_id = ?", organizationID, auth.CodeConflict}, {"record_number = ?", recordNumber, auth.CodeOpeningRecordNumberConflict}} {
		var count int64
		if err := db.Model(&gen.FranchiseOpeningRecord{}).Where(candidate.where, candidate.value).Count(&count).Error; err != nil {
			return err
		}
		if count != 0 {
			return auth.NewError(candidate.code)
		}
	}
	return nil
}

func assignInitialAccount(tx *gorm.DB, organizationID, accountID string) error {
	result := tx.Model(&gen.Organization{}).Where("id = ? AND initial_account_id IS NULL", organizationID).Update("initial_account_id", accountID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func validInitialAccountCommand(organizationID, accountID, evidenceReference string) bool {
	for _, value := range []string{organizationID, accountID} {
		if value == "" || value != strings.TrimSpace(value) || len(value) > 128 {
			return false
		}
	}
	return evidenceReference != "" && evidenceReference == strings.TrimSpace(evidenceReference) &&
		utf8.RuneCountInString(evidenceReference) <= 128 && strings.IndexFunc(evidenceReference, unicode.IsControl) == -1
}
