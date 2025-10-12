package pg

import (
	"context"

	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type pgClient struct {
	masterDBC db.DB
}

func NewPGClient(ctx context.Context, dsn string, logger Logger) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err.Error())
	}

	return &pgClient{
		masterDBC: NewDB(dbc, logger),
	}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
