package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// Bridge represents a Layer-2 bridge interface
type Bridge struct {
	Name        string   `json:"name"`
	Members     []string `json:"members"`
	Description string   `json:"description"`
}

// CreateBridge initializes a new bridge interface
func CreateBridge(bridge Bridge) (string, error) {
	// Create cloning interface
	cmd := exec.Command("ifconfig", "bridge", "create")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create bridge device: %v", err)
	}
	bridgeName := strings.TrimSpace(string(output))

	// Add members
	for _, member := range bridge.Members {
		cmd = exec.Command("ifconfig", bridgeName, "addm", member)
		if out, err := cmd.CombinedOutput(); err != nil {
			exec.Command("ifconfig", bridgeName, "destroy").Run()
			return "", fmt.Errorf("failed to add member %s to bridge: %s", member, string(out))
		}
	}

	// Set UP
	exec.Command("ifconfig", bridgeName, "up").Run()

	return bridgeName, nil
}

// DestroyBridge removes a bridge interface
func DestroyBridge(name string) error {
	cmd := exec.Command("ifconfig", name, "destroy")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to destroy bridge: %s", string(out))
	}
	return nil
}
