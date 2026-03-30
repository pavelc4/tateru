package stages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pavelc4/tateru/pkg/config"
)

type BaseStage struct {
	name string
	deps []string
	cfg  *config.BuildConfig
	env  []string
}

func (b *BaseStage) Name() string   { return b.name }
func (b *BaseStage) Deps() []string { return b.deps }

func (b *BaseStage) make(targets ...string) error {
	args := []string{
		fmt.Sprintf("-j%d", numCPU()),
		"O=out",
		"ARCH=arm64",
		"LLVM=1",
		"LLVM_IAS=1",
	}
	args = append(args, targets...)

	cmd := exec.Command("make", args...)
	cmd.Dir = b.cfg.KernelSrc
	cmd.Env = append(os.Environ(), b.env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func numCPU() int {
	out, err := exec.Command("nproc", "--all").Output()
	if err != nil {
		return 4
	}
	n := 4
	fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &n)
	return n
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
		strings.Join(s.cfg.Kernel.Defconfig.Fragments, " - "),
	)

	if err := s.make(base); err != nil {
		return fmt.Errorf("base defconfig: %w", err)
	}

	frags := s.cfg.Kernel.Defconfig.Fragments
	if len(frags) == 0 {
		return nil
	}

	kernelSrcAbs, err := filepath.Abs(s.cfg.KernelSrc)
	if err != nil {
		return fmt.Errorf("resolving KernelSrc: %w", err)
	}

	mergeScriptAbs := filepath.Join(kernelSrcAbs, "scripts/kconfig/merge_config.sh")
	if _, err := os.Stat(mergeScriptAbs); err != nil {
		return fmt.Errorf("merge_config.sh not found at %s: %w", mergeScriptAbs, err)
	}

	args := []string{
		mergeScriptAbs,
		"-O", filepath.Join(kernelSrcAbs, "out"),
		filepath.Join(kernelSrcAbs, "out", ".config"),
	}
	for _, f := range frags {
		args = append(args, filepath.Join(kernelSrcAbs, f))
	}

	cmd := exec.Command("bash", args...)
	cmd.Dir = kernelSrcAbs
	cmd.Env = append(os.Environ(), s.env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *DefconfigStage) CacheInputs() []string {
	inputs := []string{}
	for _, frag := range s.cfg.Kernel.Defconfig.Fragments {
		inputs = append(inputs, filepath.Join(s.cfg.KernelSrc, frag))
	}
	return inputs
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
