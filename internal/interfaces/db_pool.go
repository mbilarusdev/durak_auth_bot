package interfaces

import (
	"context"
)

type DBPool interface {
	Acquire(ctx context.Context) (c DBConn, err error)
}
