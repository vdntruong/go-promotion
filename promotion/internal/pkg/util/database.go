package util

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)


func IsDuplicateError(err error) bool {
	var perr *pgconn.PgError
	if ok := errors.As(err, &perr); ok {
		return perr.Code == "23505"
	}
	return  errors.Is(err, gorm.ErrDuplicatedKey)
}
