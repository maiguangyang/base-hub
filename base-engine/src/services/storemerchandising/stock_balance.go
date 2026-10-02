package storemerchandising

import (
	"math"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func getBalance(tx *gorm.DB, batchID, packageID string) (*gen.StoreStockBalance, error) {
	var balance gen.StoreStockBalance
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("batch_id = ? AND package_id = ?", batchID, packageID).First(&balance).Error
	if err == nil {
		return &balance, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return nil, nil
}

func addBalance(tx *gorm.DB, batchID, packageID string, quantity int64) error {
	if quantity <= 0 || quantity > math.MaxInt32 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	balance, err := getBalance(tx, batchID, packageID)
	if err != nil {
		return err
	}
	if balance == nil {
		return tx.Create(&gen.StoreStockBalance{ID: uuid.Must(uuid.NewV4()).String(), BatchID: batchID, PackageID: packageID, Quantity: quantity, Version: 1}).Error
	}
	if balance.Quantity > math.MaxInt32-quantity {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return tx.Model(balance).Updates(map[string]any{
		"quantity": gorm.Expr("quantity + ?", quantity),
		"version":  gorm.Expr("version + 1"),
	}).Error
}

func subtractBalance(tx *gorm.DB, batchID, packageID string, quantity int64) error {
	if quantity <= 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	result := tx.Model(&gen.StoreStockBalance{}).Where("batch_id = ? AND package_id = ? AND quantity >= ?", batchID, packageID, quantity).
		Updates(map[string]any{"quantity": gorm.Expr("quantity - ?", quantity), "version": gorm.Expr("version + 1")})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}
