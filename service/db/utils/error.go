package dbUtil

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
)

var (
	ErrorRecordNotFound      = pgx.ErrNoRows
	ErrorForeignKeyViolation = pgconn.PgError{Code: foreignKeyViolation}
	ErrorUniqueKeyViolation  = pgconn.PgError{Code: uniqueViolation}
)

func GetErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func CheckErrorCode(err error, targetCode string) bool {
	if errorCode := GetErrorCode(err); errorCode == targetCode {
		return true
	}
	return false
}
