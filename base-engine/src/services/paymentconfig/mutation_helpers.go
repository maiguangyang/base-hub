package paymentconfig

import (
	"encoding/json"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func paymentDBError(err error) error {
	var duplicate *mysql.MySQLError
	if errors.As(err, &duplicate) && duplicate.Number == 1062 {
		return ErrConflict
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrConflict
	}
	return err
}

func value(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func encodeCredential(value any) ([]byte, error) { return json.Marshal(value) }
