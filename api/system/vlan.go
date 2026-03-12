package system

import (
	"fmt"
	"os/exec"
)

// VLAN represents an 802.1Q tagged interface
type VLAN struct {
	ParentInterface string `json:"parent"`
	VlanTag         int    `json:"tag"`
	Description     string `json:"description"`
}

// CreateVLAN creates a new tagged interface on FreeBSD
func CreateVLAN(vlan VLAN) (string, error) {
	ifaceName := fmt.Sprintf("%s.%d", vlan.ParentInterface, vlan.VlanTag)
	
	// Create the vlan pseudo-device
	cmd := exec.Command("ifconfig", ifaceName, "create")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to create vlan interface: %s", string(output))
	}

	// Associate with parent and tag
	cmd = exec.Command("ifconfig", ifaceName, "vlan", fmt.Sprintf("%d", vlan.VlanTag), "vlandev", vlan.ParentInterface)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Cleanup if association fails
		exec.Command("ifconfig", ifaceName, "destroy").Run()
		return "", fmt.Errorf("failed to associate vlan: %s", string(output))
	}

	return ifaceName, nil
}

// DestroyVLAN removes a vlan interface
func DestroyVLAN(name string) error {
	cmd := exec.Command("ifconfig", name, "destroy")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to destroy vlan: %s", string(output))
	}
	return nil
}
