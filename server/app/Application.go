package app

import "context"

type Application interface {
	Start(ctx context.Context) error
}
