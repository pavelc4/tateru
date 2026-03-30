package stages

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pavelc4/tateru/pkg/config"
)

type BaseStage struct {
	name string
	deps []string
	cfg  *config.BuildConfig
	env  []string
}

func (b *BaseStage) Name() string    { return b.name }
func (b *BaseStage) Deps() []string  { return b.deps }

func (b *BaseStage) make(targets ...string) error {
	args := append([]string{"-C", b.cfg.KernelSrc}, targets...)
	cmd := exec.Command("make", args...)
	cmd.Env = append(os.Environ(), b.env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}


type DefconfigStage struct{ BaseStage }

func NewDefconfig(cfg *config.BuildConfig, env []string) *DefconfigStage {
	return &DefconfigStage{BaseStage{
		name: "defconfig",
		deps: []string{},
		cfg:  cfg,
		env:  env,
	}}
}

func (s *DefconfigStage) Run() error {
	base := s.cfg.Kernel.Defconfig.Base
	if base == "" {
		return fmt.Errorf("kernel.defconfig.base is empty")
	}
	fmt.Printf("  [defconfig] base=%s fragments=%s\n",
		base,
		strings.Join(s.cfg.Kernel.Defconfig.Fragments, " "),
	)
	return s.make(base)
}

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
		targets = []string{s.cfg.GKI.Image}
	}
	fmt.Printf("  [kernel] targets=%s\n", strings.Join(targets, " "))
	return s.make(targets...)
}
