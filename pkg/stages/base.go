package stages

import (
	"fmt"
	"io"
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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
