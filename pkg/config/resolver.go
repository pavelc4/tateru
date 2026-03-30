package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type rawNode = map[string]interface{}

func Load(cfgPath string) (*BuildConfig, error) {
	merged, err := resolveChain(cfgPath, nil)
	if err != nil {
		return nil, err
	}

	raw, err := yaml.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("re-encode: %w", err)
	}

	var cfg BuildConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &cfg, nil
}

func resolveChain(cfgPath string, visited map[string]bool) (rawNode, error) {
	abs, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil, err
	}
	if visited == nil {
		visited = make(map[string]bool)
	}
	if visited[abs] {
		return nil, fmt.Errorf("extends cycle detected at %s", abs)
	}
	visited[abs] = true

	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", abs, err)
	}

	var node rawNode
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("parse %s: %w", abs, err)
	}
	if node == nil {
		node = make(rawNode)
	}

	extendsRaw, ok := node["extends"]
	if !ok {
		return node, nil
	}
	extendsPath, ok := extendsRaw.(string)
	if !ok {
		return nil, fmt.Errorf("%s: extends must be a string", abs)
	}

	parentPath := filepath.Join(filepath.Dir(abs), extendsPath)
	parent, err := resolveChain(parentPath, visited)
	if err != nil {
		return nil, fmt.Errorf("extends of %s: %w", abs, err)
	}

	delete(node, "extends")
	return deepMerge(parent, node), nil
}

func deepMerge(dst, src rawNode) rawNode {
	result := make(rawNode, len(dst))
	for k, v := range dst {
		result[k] = v
	}
	for k, sv := range src {
		if dv, exists := result[k]; exists {
			if dstMap, ok1 := dv.(rawNode); ok1 {
				if srcMap, ok2 := sv.(rawNode); ok2 {
					result[k] = deepMerge(dstMap, srcMap)
					continue
				}
			}
		}
		result[k] = sv
	}
	return result
}
