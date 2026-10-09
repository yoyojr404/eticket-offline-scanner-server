package preflight

import (
	"context"

	"github.com/yoyojr404/internal/preflight/service/seed"
)

type Runner struct {
	seedService seed.SeedService
}

func NewRunner(seedService seed.SeedService) *Runner {
	return &Runner{seedService: seedService}
}

func (r *Runner) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// this will orchestrate some preflight service
	if err := r.seedService.Seed(ctx); err != nil {
		return err
	}

	return nil
}
