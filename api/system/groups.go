package system

import (
	"fmt"
	"os/exec"
)

// InterfaceGroup represents a logical group of physical or virtual interfaces
type InterfaceGroup struct {
	Name        string   `json:"name"`
	Members     []string `json:"members"`
	Description string   `json:"description"`
}

// ApplyGroup uses interface groups in pf by generating macros or using ifconfig groups
func ApplyGroup(group InterfaceGroup) error {
	for _, member := range group.Members {
		cmd := exec.Command("ifconfig", member, "group", group.Name)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add %s to group %s: %s", member, group.Name, string(output))
		}
	}
	return nil
}

// RemoveFromGroup removes an interface from a group
func RemoveFromGroup(member string, group string) error {
	cmd := exec.Command("ifconfig", member, "-group", group)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove %s from group %s: %s", member, group, string(output))
	}
	return nil
}
