//go:build !windows

package license

func windowsMachineID() string {
	return ""
}
