package license

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type LicenseCache struct {
	LicenseKey string `json:"license_key"`
	Owner      string `json:"owner,omitempty"`
}

func CachePath(dataPath string) string {
	return filepath.Join(dataPath, "license.json")
}

func LoadCache(dataPath string) (*LicenseCache, error) {
	path := CachePath(dataPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &LicenseCache{}, nil
		}
		return nil, err
	}
	var cache LicenseCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	cache.LicenseKey = strings.TrimSpace(cache.LicenseKey)
	return &cache, nil
}

func SaveCache(dataPath string, cache *LicenseCache) error {
	if err := os.MkdirAll(dataPath, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(CachePath(dataPath), data, 0o600)
}
