package dberr

import (
	"bookmark-management/internal/domain"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTranslate(t *testing.T) {
	t.Parallel()
	otherErr := errors.New("connection refused")

	testcases := []struct {
		name        string
		err         error
		expectedErr error
	}{
		{name: "nil", err: nil, expectedErr: nil},
		{name: "duplicated key", err: gorm.ErrDuplicatedKey, expectedErr: domain.ErrAlreadyExists},
		{name: "check constraint", err: gorm.ErrCheckConstraintViolated, expectedErr: domain.ErrInvalidInput},
		{name: "not null violation", err: &pgconn.PgError{Code: "23502"}, expectedErr: domain.ErrInvalidInput},
		{name: "string too long", err: &pgconn.PgError{Code: "22001"}, expectedErr: domain.ErrInvalidInput},
		{name: "other pg error is unchanged", err: &pgconn.PgError{Code: "42P01"}, expectedErr: &pgconn.PgError{Code: "42P01"}},
		{name: "unknown error is unchanged", err: otherErr, expectedErr: otherErr},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := Translate(tc.err)
			if tc.expectedErr == nil {
				assert.NoError(t, actual)
				return
			}
			if pgErr, ok := errors.AsType[*pgconn.PgError](tc.expectedErr); ok {
				var got *pgconn.PgError
				assert.ErrorAs(t, actual, &got)
				assert.Equal(t, pgErr.Code, got.Code)
				assert.NotErrorIs(t, actual, domain.ErrInvalidInput)
				return
			}
			assert.ErrorIs(t, actual, tc.expectedErr)
		})
	}
}
