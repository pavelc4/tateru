package toolchain

import (
	"fmt"
	"os/exec"
	"strings"
)

func DetectClang(override string) (*Env, error) {
	bin := override
	if bin == "" {
		bin = "clang"
	}

	path, err := exec.LookPath(bin)
	if err != nil {
		return nil, fmt.Errorf("clang not found (%s): %w", bin, err)
	}

	ver, err := ClangVersion(path)
	if err != nil {
		return nil, err
	}

	return &Env{ClangBin: path, Version: ver}, nil
}

func ClangVersion(bin string) (string, error) {
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("clang --version: %w", err)
	}
	line := strings.SplitN(string(out), "\n", 2)[0]
	return strings.TrimSpace(line), nil
}
