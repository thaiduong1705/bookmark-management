// Package dberr translates database errors into domain errors so upper layers
// never depend on SQL or driver details.
package dberr

import (
	"bookmark-management/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// PostgreSQL error codes gorm does not translate on its own,
// see https://www.postgresql.org/docs/current/errcodes-appendix.html.
const (
	pgNotNullViolation = "23502"
	pgStringTooLong    = "22001"
)

// Translate maps a database error to a domain error. The original error is kept
// in the chain for logging. Unknown errors are returned unchanged and should be
// treated as system errors. It relies on gorm.Config{TranslateError: true}.
func Translate(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: %w", domain.ErrAlreadyExists, err)
	}
	if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		return fmt.Errorf("%w: %w", domain.ErrInvalidInput, err)
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case pgNotNullViolation, pgStringTooLong:
			return fmt.Errorf("%w: %w", domain.ErrInvalidInput, err)
		}
	}

	return err
}
