package toolchain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DetectClang(override string) (*Env, error) {
	bin := resolveClangBin(override)

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

func resolveClangBin(cfgOverride string) string {
	if cfgOverride != "" {
		return cfgOverride
	}

	if dir := os.Getenv("TATERU_CLANG"); dir != "" {
		candidate := filepath.Join(dir, "clang")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		return dir
	}

	return "clang"
}

func ClangVersion(bin string) (string, error) {
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("clang --version: %w", err)
	}
	line := strings.SplitN(string(out), "\n", 2)[0]
	return strings.TrimSpace(line), nil
}
