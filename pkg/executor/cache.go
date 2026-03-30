package executor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CacheKey struct {
	Stage   string
	Inputs  []string
	EnvHash string
}

type Cache struct {
	dir string
}

func NewCache(dir string) *Cache {
	_ = os.MkdirAll(dir, 0755)
	return &Cache{dir: dir}
}

func (c *Cache) Key(ck CacheKey) (string, error) {
	h := sha256.New()
	for _, path := range ck.Inputs {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("cache hash %s: %w", path, err)
		}
		h.Write(data)
	}

	meta, _ := json.Marshal(map[string]string{
		"stage": ck.Stage,
		"env":   ck.EnvHash,
	})
	h.Write(meta)

	return hex.EncodeToString(h.Sum(nil))[:16], nil
}

func (c *Cache) Hit(key string) bool {
	_, err := os.Stat(filepath.Join(c.dir, key))
	return err == nil
}

func (c *Cache) Store(key string) error {
	f, err := os.Create(filepath.Join(c.dir, key))
	if err != nil {
		return err
	}
	return f.Close()
}
