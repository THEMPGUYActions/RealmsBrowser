package license

import (
	"fmt"
	"os"
	"path/filepath"
)

// RealmsBrowser uses the exact same SSH authentication credential as
// RealmsLauncher. The credential is supplied externally and is never stored
// in this repository.
func writeAuthFile() (string, error) {
	source := os.Getenv("REALMS_LICENSE_SSH_AUTH_FILE")
	if source == "" {
		return "", fmt.Errorf("RealmsLauncher license SSH authentication file is not configured")
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("RealmsLauncher license SSH authentication file is empty")
	}

	dir, err := os.MkdirTemp("", "realmsbrowser-license-")
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "license_auth")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		_ = os.RemoveAll(dir)
		return "", err
	}
	return path, nil
}

func removeAuthFile(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path)
	_ = os.Remove(filepath.Dir(path))
}
