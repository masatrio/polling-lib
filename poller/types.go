package poller

import (
	"context"
)

type Poller interface {
	Run(ctx context.Context, resp chan interface{}) error
	SetTimeout(timeout int64)
	Stop()
}
