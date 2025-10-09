package interfaces

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPool interface {
	Acquire(ctx context.Context) (c *pgxpool.Conn, err error)
}
