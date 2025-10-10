package interfaces

import (
	"context"
)

type DBConn interface {
	Release()
	QueryRow(ctx context.Context, sql string, args ...any) DBRow
}
