package background

import (
	"context"
	"errors"
	"sync"
)

type Worker interface {
	Shutdown(ctx context.Context) error
	Run(ctx context.Context) error
}

type RunnerGroup struct {
	mu      sync.RWMutex
	runners map[Worker]bool
}

func NewRunnerGroup() *RunnerGroup {
	return &RunnerGroup{
		runners: make(map[Worker]bool),
		mu:      sync.RWMutex{},
	}
}

func (rg *RunnerGroup) Add(runner Worker) {
	rg.mu.Lock()
	defer rg.mu.Unlock()
	rg.runners[runner] = true
}

func (rg *RunnerGroup) Remove(runner Worker) {
	rg.mu.Lock()
	defer rg.mu.Unlock()
	delete(rg.runners, runner)
}

func (rg *RunnerGroup) RunAll(ctx context.Context) error {
	for runner := range rg.runners {
		go runner.Run(ctx)
	}
	return nil
}

func (rg *RunnerGroup) ShutdownAll(ctx context.Context) error {
	var err error

	for runner := range rg.runners {
		err = errors.Join(err, runner.Shutdown(ctx))
	}

	return err
}

func (rg *RunnerGroup) Shutdown(ctx context.Context, runner Worker) error {
	err := runner.Shutdown(ctx)
	if err != nil {
		return err
	}

	rg.Remove(runner)
	return nil
}
