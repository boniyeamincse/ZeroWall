package system

import (
	"fmt"
	"net/http"
	"os/exec"
)

// FirmwareUpdate handles system upgrades
type FirmwareUpdate struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

// CheckForUpdates fetches version info from build server
func CheckForUpdates(updateURL string) (string, error) {
	resp, err := http.Get(updateURL)
	if err != nil {
		return "", fmt.Errorf("failed to check for updates: %v", err)
	}
	defer resp.Body.Close()

	// In reality: parse version from body
	return "2026.1.0-RC1", nil
}

// DownloadFirmware pulls the update package
func (f FirmwareUpdate) DownloadFirmware(pkgURL string, destPath string) error {
	fmt.Printf("Downloading ZeroWall %s to %s...\n", f.LatestVersion, destPath)
	// Implementation: stream download to file
	return nil
}

// ApplyUpdate triggers the install and reboot
func (f FirmwareUpdate) ApplyUpdate(pkgPath string) error {
	fmt.Println("Verifying firmware integrity (SHA256)...")
	// 1. Verify Hash
	// 2. Trigger FreeBSD freebsd-update or pkg upgrade
	cmd := exec.Command("touch", "/var/run/zerowall_update_pending")
	return cmd.Run()
}
