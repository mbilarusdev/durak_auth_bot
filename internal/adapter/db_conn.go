package adapter

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
)

type AdapterConn struct {
	conn *pgxpool.Conn
}

func (a *AdapterConn) QueryRow(
	ctx context.Context,
	query string,
	args ...interface{},
) interfaces.DBRow {
	return a.conn.QueryRow(ctx, query, args...)
}

func NewAdapterConn(conn *pgxpool.Conn) interfaces.DBConn {
	return &AdapterConn{conn: conn}
}

func (p *AdapterConn) Release() {
	p.conn.Release()
}
