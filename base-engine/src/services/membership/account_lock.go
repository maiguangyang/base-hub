package membership

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func inviteAccountQuery(tx *gorm.DB, phone string) *gorm.DB {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("phone = ?", phone)
}
