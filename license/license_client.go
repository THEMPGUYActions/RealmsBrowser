package license

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type LicenseResponse struct {
	Status string `json:"status"`
	Owner  string `json:"owner"`
	Reason string `json:"reason"`
}

type LicenseClient struct{}

// Validate preserves the RealmsLauncher server command and JSON response
// contract. The executable is supplied externally so authentication material
// never needs to be embedded in RealmsBrowser.
func (c *LicenseClient) Validate(key string) (*LicenseResponse, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("no license key has been configured")
	}

	validator := os.Getenv("REALMS_LICENSE_VALIDATOR")
	if validator == "" {
		return nil, fmt.Errorf("RealmsLauncher license validator is not configured")
	}

	cmd := exec.Command(validator, "check", key, DeviceFingerprint(), DeviceName())
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
