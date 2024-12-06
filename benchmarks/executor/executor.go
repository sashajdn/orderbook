package executor

import (
	"context"
)

type Executor interface {
	RunIteration(ctx context.Context) error
	Name() string
}
