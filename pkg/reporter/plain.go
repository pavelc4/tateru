package reporter

import (
	"fmt"
	"sort"
	"time"

	"github.com/pavelc4/tateru/pkg/dag"
)

type Plain struct{}

func NewPlain() *Plain { return &Plain{} }

func (p *Plain) Start(stageName string) {
	fmt.Printf("  → %-20s running...\n", stageName)
}

func (p *Plain) Done(r *dag.StageResult) {
	icon := statusIcon(r.Status)
	dur := ""
	if r.Duration > 0 {
		dur = fmt.Sprintf(" (%s)", r.Duration.Round(time.Millisecond))
	}
	cache := ""
	if r.CacheHit {
		cache = " [cached]"
	}
	if r.Err != nil {
		fmt.Printf("  %s %-20s%s%s — %v\n", icon, r.Name, dur, cache, r.Err)
	} else {
		fmt.Printf("  %s %-20s%s%s\n", icon, r.Name, dur, cache)
	}
}

func (p *Plain) Summary(results map[string]*dag.StageResult) {
	names := make([]string, 0, len(results))
	for n := range results {
		names = append(names, n)
	}
	sort.Strings(names)

	var done, failed, skipped int
	var total time.Duration

	for _, n := range names {
		r := results[n]
		switch r.Status {
		case dag.StatusDone:
			done++
		case dag.StatusFailed:
			failed++
		case dag.StatusSkipped:
			skipped++
		}
		total += r.Duration
	}

	fmt.Println("\n────────────────────────────────")
	fmt.Printf("  done: %d  failed: %d  skipped: %d  total: %s\n",
		done, failed, skipped, total.Round(time.Millisecond))

	if failed > 0 {
		fmt.Println("\nfailed stages:")
		for _, n := range names {
			r := results[n]
			if r.Status == dag.StatusFailed {
				fmt.Printf("  • %s: %v\n", r.Name, r.Err)
			}
		}
	}
}

func statusIcon(s dag.StageStatus) string {
	switch s {
	case dag.StatusDone:
		return "✓"
	case dag.StatusFailed:
		return "✗"
	case dag.StatusSkipped:
		return "–"
	default:
		return "?"
	}
}
