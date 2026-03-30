package dag

import "time"

type StageStatus int

const (
	StatusPending StageStatus = iota
	StatusRunning
	StatusDone
	StatusFailed
	StatusSkipped
)

func (s StageStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusDone:
		return "done"
	case StatusFailed:
		return "failed"
	case StatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

type StageResult struct {
	Name     string
	Status   StageStatus
	Duration time.Duration
	Err      error
	CacheHit bool
}
