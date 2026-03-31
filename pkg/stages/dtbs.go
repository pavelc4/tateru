package stages

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pavelc4/tateru/pkg/config"
)

type DTBsStage struct{ BaseStage }

func NewDTBs(cfg *config.BuildConfig, env []string) *DTBsStage {
	return &DTBsStage{BaseStage{
		name: "dtbs",
		deps: []string{"defconfig"},
		cfg:  cfg,
		env:  env,
	}}
}

func (s *DTBsStage) Run() error {
	fmt.Printf("  [dtbs] wildcard=%s dtbo=%s\n",
		s.cfg.DTB.Wildcard,
		s.cfg.DTB.DTBOWildcard,
	)

	if err := s.make("dtbs"); err != nil {
		return fmt.Errorf("make dtbs: %w", err)
	}

	kernelSrcAbs, err := filepath.Abs(s.cfg.KernelSrc)
	if err != nil {
		return fmt.Errorf("resolving KernelSrc: %w", err)
	}

	distDir := s.cfg.Output.Dist
	if !filepath.IsAbs(distDir) {
		wsRoot := filepath.Join(kernelSrcAbs, "..", "..")
		distDir = filepath.Join(wsRoot, distDir)
	}
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return fmt.Errorf("create dist dir: %w", err)
	}

	dtbSrcDir := filepath.Join(kernelSrcAbs, "out", "arch", "arm64", "boot", "dts")

	if s.cfg.DTB.Wildcard != "" {
		if err := copyByPattern(dtbSrcDir, distDir, s.cfg.DTB.Wildcard, ".dtb"); err != nil {
			return fmt.Errorf("copy dtbs: %w", err)
		}
	}

	if s.cfg.DTB.DTBOWildcard != "" {
		if err := copyByPattern(dtbSrcDir, distDir, s.cfg.DTB.DTBOWildcard, ".dtbo"); err != nil {
			return fmt.Errorf("copy dtbos: %w", err)
		}
	}

	return nil
}

func (s *DTBsStage) CacheInputs() []string {
	kernelSrcAbs, _ := filepath.Abs(s.cfg.KernelSrc)
	return []string{
		filepath.Join(kernelSrcAbs, "out", ".config"),
	}
}

func copyByPattern(srcDir, dstDir, pattern, ext string) error {
	copied := 0
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		name := d.Name()
		if !strings.HasSuffix(name, ext) {
			return nil
		}
		base := strings.TrimSuffix(name, ext)
		matched, err := filepath.Match(pattern, base)
		if err != nil || !matched {
			return err
		}
		dst := filepath.Join(dstDir, name)
		if err := copyFile(path, dst); err != nil {
			return fmt.Errorf("copy %s: %w", name, err)
		}
		fmt.Printf("  [dtbs] %s → dist\n", name)
		copied++
		return nil
	})
	if err != nil {
		return err
	}
	if copied == 0 {
		fmt.Printf("  [dtbs] warning: no files matched pattern=%s ext=%s\n", pattern, ext)
	}
	return nil
}
