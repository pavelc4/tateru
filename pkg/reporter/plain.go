package reporter

import (
	"fmt"
	"time"

	"github.com/pavelc4/tateru/pkg/dag"
)

type Plain struct{}

func NewPlain() *Plain { return &Plain{} }

var statusIcon = map[dag.StageStatus]string{
	dag.StatusDone:    "✓",
	dag.StatusFailed:  "✗",
	dag.StatusSkipped: "–",
	dag.StatusRunning: "…",
	dag.StatusPending: "·",
}

func (p *Plain) StageStart(name string) {
	fmt.Printf("  → %-24s running\n", name)
}

func (p *Plain) StageEnd(r *dag.StageResult) {
	icon := statusIcon[r.Status]
	cache := ""
	if r.CacheHit {
		cache = " [cached]"
	}
	errStr := ""
	if r.Err != nil {
		errStr = fmt.Sprintf(" (%s)", r.Err)
	}
	fmt.Printf("  %s %-24s %s%s%s\n",
		icon, r.Name,
		r.Duration.Round(time.Millisecond),
		cache, errStr,
	)
}

func (p *Plain) Summary(results map[string]*dag.StageResult, elapsed string) {
	total, done, failed, skipped, cached := 0, 0, 0, 0, 0
	for _, r := range results {
		total++
		switch r.Status {
		case dag.StatusDone:
			done++
			if r.CacheHit {
				cached++
			}
		case dag.StatusFailed:
			failed++
		case dag.StatusSkipped:
			skipped++
		}
	}

	fmt.Printf("\n╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌\n")
	fmt.Printf("  total: %d  done: %d  failed: %d  skipped: %d  cached: %d\n",
		total, done, failed, skipped, cached)
	fmt.Printf("  elapsed: %s\n", elapsed)

	if failed > 0 {
		fmt.Println("  BUILD FAILED")
	} else {
		fmt.Println("  BUILD OK")
	}
}
