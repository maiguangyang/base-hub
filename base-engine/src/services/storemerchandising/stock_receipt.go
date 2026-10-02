package storemerchandising

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockReceipt struct {
	StoreID, ListingID, PackageID, BatchNumber, RequestKey string
	Quantity                                               int64
	ProducedAt, ExpiresAt                                  *time.Time
	SourceReference                                        *string
	Barcode                                                *string
}

func (s *Service) ReceiveStock(ctx context.Context, principal *auth.WorkspacePrincipal, input StockReceipt) (*gen.StoreStockMovement, error) {
	input.BatchNumber = strings.TrimSpace(input.BatchNumber)
	if !validStockReceipt(input) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var movement gen.StoreStockMovement
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.receiveStockTx(tx, principal, input, &movement)
	})
	return &movement, err
}

func validStockReceipt(input StockReceipt) bool {
	if input.StoreID == "" || input.ListingID == "" || input.PackageID == "" {
		return false
	}
	if input.Quantity <= 0 || input.Quantity > math.MaxInt32 || !validRequestKey(input.RequestKey) {
		return false
	}
	if len(input.BatchNumber) > 128 {
		return false
	}
	return input.ProducedAt == nil || input.ExpiresAt == nil || input.ExpiresAt.After(*input.ProducedAt)
}

func (s *Service) receiveStockTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, input StockReceipt, movement *gen.StoreStockMovement) error {
	// 首次数据库读取即锁定门店：防止并发序号冲突及重复请求读取旧快照。
	store, err := storeScope(tx.Clauses(clause.Locking{Strength: "UPDATE"}), principal, input.StoreID, "franchiseStock:manage", authorization.AccessCreate)
	if err != nil {
		return err
	}
	found, err := replayMovement(tx, store.ID, gen.StockMovementKindReceive, input.RequestKey)
	if err != nil {
		return err
	}
	if found != nil {
		*movement = *found
		return validateReceiptReplay(tx, movement, input)
	}
	batch, pack, err := receiptTarget(tx, store.ID, input)
	if err != nil {
		return err
	}
	if err := assertBatchBalanced(tx, batch.ID); err != nil {
		return err
	}
	if err := addBalance(tx, batch.ID, pack.ID, input.Quantity); err != nil {
		return err
	}
	*movement = gen.StoreStockMovement{ID: uuid.Must(uuid.NewV4()).String(), StoreID: store.ID, BatchID: batch.ID,
		Kind: gen.StockMovementKindReceive, RequestKey: input.RequestKey, TargetQuantity: &input.Quantity,
		TargetPackageID: &pack.ID, PackageSetVersion: pack.PackageSetVersion, OccurredAt: time.Now()}
	if err := tx.Create(movement).Error; err != nil {
		return err
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseStock:manage",
		ResourceType: "storeStockMovement", ResourceID: movement.ID, ResultCode: "SUCCESS"})
}

func receiptTarget(tx *gorm.DB, storeID string, input StockReceipt) (*gen.StoreInventoryBatch, *gen.ProductPackage, error) {
	var listing gen.StoreListing
	if err := tx.Where("id = ? AND store_id = ? AND enabled = ?", input.ListingID, storeID, true).First(&listing).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodeValidationFailed)
	}
	if err := publishedSku(tx, listing.SkuID); err != nil {
		return nil, nil, err
	}
	var pack gen.ProductPackage
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND sku_id = ? AND enabled = ?", input.PackageID, listing.SkuID, true).First(&pack).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodeValidationFailed)
	}
	if err := validateReceiptBarcode(tx, &pack, input.Barcode); err != nil {
		return nil, nil, err
	}
	if err := validateReceiptExpiry(tx, listing.SkuID, input); err != nil {
		return nil, nil, err
	}
	batch, err := receiptBatch(tx, listing.ID, input)
	return batch, &pack, err
}

func validateReceiptExpiry(tx *gorm.DB, skuID string, input StockReceipt) error {
	var sku gen.ProductSku
	if err := tx.First(&sku, "id = ?", skuID).Error; err != nil {
		return err
	}
	if sku.ShelfLifeDays != nil && input.ExpiresAt == nil {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if input.ExpiresAt != nil && !input.ExpiresAt.After(time.Now()) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func receiptBatch(tx *gorm.DB, listingID string, input StockReceipt) (*gen.StoreInventoryBatch, error) {
	batchNumber := input.BatchNumber
	if batchNumber == "" {
		var err error
		batchNumber, err = nextStockBatchNumber(tx, input.StoreID, time.Now())
		if err != nil {
			return nil, err
		}
	}
	var batch gen.StoreInventoryBatch
	err := tx.Where("listing_id = ? AND batch_number = ?", listingID, batchNumber).First(&batch).Error
	if err == nil {
		if !sameReceiptBatchTerms(&batch, input) {
			return nil, auth.NewError(auth.CodeValidationFailed)
		}
		return &batch, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	batch = gen.StoreInventoryBatch{ID: uuid.Must(uuid.NewV4()).String(), ListingID: listingID, BatchNumber: batchNumber,
		ProducedAt: input.ProducedAt, ExpiresAt: input.ExpiresAt, SourceReference: input.SourceReference}
	return &batch, tx.Create(&batch).Error
}

func sameTime(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Equal(*b)
}
