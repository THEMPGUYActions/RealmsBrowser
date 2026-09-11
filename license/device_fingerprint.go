package license

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"runtime"
	"strings"
)

func DeviceFingerprint() string {
	id := persistedMachineID()
	if id != "" {
		return id
	}
	fallback := strings.Join([]string{runtime.GOOS, runtime.GOARCH, DeviceName()}, "|")
	sum := sha256.Sum256([]byte(fallback))
	return hex.EncodeToString(sum[:])[:32]
}

func persistedMachineID() string {
	var id string
	switch runtime.GOOS {
	case "windows":
		id = windowsMachineID()
	case "linux":
		if data, err := os.ReadFile("/etc/machine-id"); err == nil {
			id = strings.TrimSpace(string(data))
		}
	}
	return strings.TrimSpace(id)
}

func DeviceName() string {
	name, _ := os.Hostname()
	return strings.TrimSpace(name)
}
