package stages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pavelc4/tateru/pkg/config"
)

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
