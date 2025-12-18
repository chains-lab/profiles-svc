package pgdb

import (
	"context"
	"database/sql"

	"github.com/chains-lab/profiles-svc/internal/domain/entity"
)

type txKeyType struct{}

var TxKey = txKeyType{}

func TxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	return tx, ok
}

func (p Profile) ToEntity() entity.Profile {
	profile := entity.Profile{
		AccountID:   p.AccountID,
		Username:    p.Username,
		Official:    p.Official,
		Pseudonym:   p.Pseudonym,
		Description: p.Description,
		Avatar:      p.Avatar,

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	return profile
}
