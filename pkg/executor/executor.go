package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/pavelc4/tateru/pkg/dag"
)

type Options struct {
	DryRun   bool
	Jobs     int
	CacheDir string
}

type Executor struct {
	graph   *dag.Graph
	opts    Options
	cache   *Cache
	results map[string]*dag.StageResult
	mu      sync.Mutex
}

func New(g *dag.Graph, opts Options) *Executor {
	cacheDir := opts.CacheDir
	if cacheDir == "" {
		cacheDir = ".tateru-cache"
	}
	return &Executor{
		graph:   g,
		opts:    opts,
		cache:   NewCache(cacheDir),
		results: make(map[string]*dag.StageResult),
	}
}

func (e *Executor) Run(ctx context.Context) error {
	if err := e.graph.Build(); err != nil {
		return fmt.Errorf("build graph: %w", err)
	}

	order := e.graph.Order()
	jobs := e.opts.Jobs
	if jobs <= 0 {
		jobs = 4
	}

	sem := make(chan struct{}, jobs)
	eg, ctx := errgroup.WithContext(ctx)

	done := make(map[string]chan struct{}, len(order))
	for _, name := range order {
		done[name] = make(chan struct{})
	}

	for _, name := range order {
		name := name
		stage, _ := e.graph.Get(name)

		eg.Go(func() error {
			for _, dep := range stage.Deps() {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-done[dep]:
					e.mu.Lock()
					res := e.results[dep]
					e.mu.Unlock()
					if res != nil && res.Status == dag.StatusFailed {
						close(done[name])
						e.setResult(name, dag.StatusSkipped, 0, nil, false)
						return nil
					}
				}
			}

			sem <- struct{}{}
			defer func() { <-sem }()

			result := e.runStage(ctx, stage)
			e.setResult(name, result.Status, result.Duration, result.Err, result.CacheHit)
			close(done[name])

			if result.Err != nil {
				return fmt.Errorf("stage %s: %w", name, result.Err)
			}
			return nil
		})
	}

	return eg.Wait()
}

func (e *Executor) runStage(ctx context.Context, s dag.Stage) dag.StageResult {
	start := time.Now()

	if e.opts.DryRun {
		fmt.Printf("  [dry-run] %s\n", s.Name())
		return dag.StageResult{
			Name:   s.Name(),
			Status: dag.StatusDone,
		}
	}

	if err := s.Run(); err != nil {
		return dag.StageResult{
			Name:     s.Name(),
			Status:   dag.StatusFailed,
			Duration: time.Since(start),
			Err:      err,
		}
	}

	return dag.StageResult{
		Name:     s.Name(),
		Status:   dag.StatusDone,
		Duration: time.Since(start),
	}
}

func (e *Executor) setResult(name string, status dag.StageStatus, dur time.Duration, err error, hit bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.results[name] = &dag.StageResult{
		Name:     name,
		Status:   status,
		Duration: dur,
		Err:      err,
		CacheHit: hit,
	}
}

func (e *Executor) Results() map[string]*dag.StageResult {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make(map[string]*dag.StageResult, len(e.results))
	for k, v := range e.results {
		out[k] = v
	}
	return out
}
