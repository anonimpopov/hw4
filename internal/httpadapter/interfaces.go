package httpadapter

import (
	"context"
)

type Adapter interface {
	Serve() error
	Shutdown(ctx context.Context)
}
