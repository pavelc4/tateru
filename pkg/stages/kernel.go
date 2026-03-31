package stages

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pavelc4/tateru/pkg/config"
)

type KernelStage struct{ BaseStage }

func NewKernel(cfg *config.BuildConfig, env []string) *KernelStage {
	return &KernelStage{BaseStage{
		name: "kernel",
		deps: []string{"defconfig"},
		cfg:  cfg,
		env:  env,
	}}
}

func (s *KernelStage) Run() error {
	targets := s.cfg.GKI.MakeTargets
	if len(targets) == 0 {
		targets = []string{s.cfg.GKI.Image, "modules", "dtbs"}
	}
	fmt.Printf("  [kernel] targets=%s\n", strings.Join(targets, " "))
	return s.make(targets...)
}

func (s *KernelStage) CacheInputs() []string {
	return []string{
		filepath.Join(s.cfg.KernelSrc, "Makefile"),
	}
}
