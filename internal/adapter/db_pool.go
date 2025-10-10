package adapter

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
)

type AdapterPool struct {
	pool *pgxpool.Pool
}

func NewAdapterPool(pool *pgxpool.Pool) interfaces.DBPool {
	return &AdapterPool{pool: pool}
}

func (p *AdapterPool) Acquire(ctx context.Context) (interfaces.DBConn, error) {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return NewAdapterConn(conn), nil
}
