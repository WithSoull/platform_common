package txctx

import (
	"context"

	"github.com/WithSoull/AuthService/internal/contextx"
	"github.com/jackc/pgx/v4"
)

const TxKey contextx.CtxKey = "tx"

func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func ExtractTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return nil, false
	}
	return tx, true
}
