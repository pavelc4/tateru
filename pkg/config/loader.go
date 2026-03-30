package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const workspaceMarker = ".tateruroot"

func FindWorkspaceRoot(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, workspaceMarker)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("workspace root not found (no %s in %s or any parent)", workspaceMarker, start)
		}
		dir = parent
	}
}

func ResolveTarget(wsRoot, target string) (string, BuildMode, error) {
	// GKI target: configs/gki/gki-5.10.yml
	gkiPath := filepath.Join(wsRoot, "configs", "gki", target+".yml")
	if _, err := os.Stat(gkiPath); err == nil {
		return gkiPath, BuildModeGKI, nil
	}

	// Device target: configs/platforms/*/<device>/build.yml
	pattern := filepath.Join(wsRoot, "configs", "platforms", "*", target, "build.yml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", 0, err
	}
	if len(matches) == 1 {
		return matches[0], BuildModeDevice, nil
	}
	if len(matches) > 1 {
		return "", 0, fmt.Errorf("ambiguous target %q: matches %v", target, matches)
	}

	return "", 0, fmt.Errorf("target %q not found", target)
}
