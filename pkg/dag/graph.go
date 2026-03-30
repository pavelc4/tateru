package dag

import "fmt"

type Stage interface {
	Name() string
	Deps() []string
	Run() error
}

type Graph struct {
	stages map[string]Stage
	order  []string
}

func New() *Graph {
	return &Graph{stages: make(map[string]Stage)}
}

func (g *Graph) Add(s Stage) {
	g.stages[s.Name()] = s
}

func (g *Graph) Build() error {
	inDegree := make(map[string]int, len(g.stages))
	for name := range g.stages {
		inDegree[name] = 0
	}

	for _, s := range g.stages {
		for _, dep := range s.Deps() {
			if _, ok := g.stages[dep]; !ok {
				return fmt.Errorf("stage %q: unknown dep %q", s.Name(), dep)
			}
			inDegree[s.Name()]++
		}
	}

	queue := []string{}
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	order := []string{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)

		for _, s := range g.stages {
			for _, dep := range s.Deps() {
				if dep == cur {
					inDegree[s.Name()]--
					if inDegree[s.Name()] == 0 {
						queue = append(queue, s.Name())
					}
				}
			}
		}
	}

	if len(order) != len(g.stages) {
		return fmt.Errorf("cycle detected in stage graph")
	}

	g.order = order
	return nil
}

func (g *Graph) Order() []string {
	return g.order
}

func (g *Graph) Get(name string) (Stage, bool) {
	s, ok := g.stages[name]
	return s, ok
}
