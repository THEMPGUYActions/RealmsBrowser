package license

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type LicenseResponse struct {
	Status string `json:"status"`
	Owner  string `json:"owner"`
	Reason string `json:"reason"`
}

type LicenseClient struct{}

//go:embed realms-license-validator.cmd
var validatorScript []byte

func (c *LicenseClient) Validate(key string) (*LicenseResponse, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("no license key has been configured")
	}

	validator := os.Getenv("REALMS_LICENSE_VALIDATOR")
	if validator == "" {
		return c.validateWithEmbeddedAdapter(key)
	}

	return validateWithCommand(validator, key)
}

func (c *LicenseClient) validateWithEmbeddedAdapter(key string) (*LicenseResponse, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("the bundled RealmsLauncher license adapter requires Windows")
	}

	dir, err := os.MkdirTemp("", "realmsbrowser-license-validator-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	adapter := filepath.Join(dir, "realms-license-validator.cmd")
	if err := os.WriteFile(adapter, validatorScript, 0o600); err != nil {
		return nil, err
	}
	return validateWithCommand(adapter, key)
}

func validateWithCommand(command string, key string) (*LicenseResponse, error) {
	args := []string{"check", key, DeviceFingerprint(), DeviceName()}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" && strings.HasSuffix(strings.ToLower(command), ".cmd") {
		cmd = exec.Command("cmd.exe", "/d", "/c", command, args[0], args[1], args[2], args[3])
	} else {
		cmd = exec.Command(command, args...)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not reach the license validator: %w", err)
	}

	var response LicenseResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("the license validator returned an invalid response: %w", err)
	}
	if response.Status != "OK" {
		if response.Reason == "" {
			response.Reason = "The license was rejected by the license server."
		}
		return &response, fmt.Errorf("%s", response.Reason)
	}
	return &response, nil
}
