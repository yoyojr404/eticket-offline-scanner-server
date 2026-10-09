package seed

import (
	"context"
)

type SeedService interface {
	Seed(ctx context.Context) error
}
