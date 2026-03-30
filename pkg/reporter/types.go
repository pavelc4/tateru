package reporter

import "github.com/pavelc4/tateru/pkg/dag"

type Reporter interface {
	Start(stageName string)
	Done(result *dag.StageResult)
	Summary(results map[string]*dag.StageResult)
}
