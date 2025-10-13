package pg

import (
	"context"

	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/WithSoull/platform_common/pkg/client/db/prettier"
	"github.com/WithSoull/platform_common/pkg/contextx/txctx"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
}

type pg struct {
	dbc *pgxpool.Pool
	l   Logger
}

func NewDB(dbc *pgxpool.Pool, logger Logger) db.DB {
	return &pg{
		dbc: dbc,
		l:   logger,
	}
}

func (p *pg) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...any) (pgconn.CommandTag, error) {
	p.logQuery(ctx, q, args...)

	tx, ok := txctx.ExtractTx(ctx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...any) (pgx.Rows, error) {
	p.logQuery(ctx, q, args...)

	tx, ok := txctx.ExtractTx(ctx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...any) pgx.Row {
	p.logQuery(ctx, q, args...)

	tx, ok := txctx.ExtractTx(ctx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.QueryRow(ctx, "SELECT 1").Scan(new(int))
}

func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

func (p *pg) Close() {
	p.dbc.Close()
}

func (p *pg) logQuery(ctx context.Context, q db.Query, args ...any) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	p.l.Debug(
		ctx,
		"PG Querry",
		zap.String("sql", q.Name),
		zap.String("query", prettyQuery),
	)
}
