package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// LAGGProto represents the link aggregation protocol (lacp, failover, loadbalance)
type LAGGProto string

const (
	LAGGProtocolLACP        LAGGProto = "lacp"
	LAGGProtocolFailover    LAGGProto = "failover"
	LAGGProtocolLoadBalance LAGGProto = "loadbalance"
)

// LAGG represents a link aggregation group
type LAGG struct {
	Name            string      `json:"name"`
	Members         []string    `json:"members"`
	Protocol        LAGGProto   `json:"protocol"`
}

// CreateLAGG initializes a new link aggregation interface
func CreateLAGG(lagg LAGG) (string, error) {
	// Create cloning interface (e.g., lagg0)
	cmd := exec.Command("ifconfig", "lagg", "create")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create lagg device: %v", err)
	}
	laggName := strings.TrimSpace(string(output))

	// Set protocol
	cmd = exec.Command("ifconfig", laggName, "laggproto", string(lagg.Protocol))
	if out, err := cmd.CombinedOutput(); err != nil {
		exec.Command("ifconfig", laggName, "destroy").Run()
		return "", fmt.Errorf("failed to set laggproto: %s", string(out))
	}

	// Add members
	for _, member := range lagg.Members {
		cmd = exec.Command("ifconfig", laggName, "laggport", member)
		if out, err := cmd.CombinedOutput(); err != nil {
			exec.Command("ifconfig", laggName, "destroy").Run()
			return "", fmt.Errorf("failed to add member %s to lagg: %s", member, string(out))
		}
	}

	return laggName, nil
}

// DestroyLAGG tears down a lagg interface
func DestroyLAGG(name string) error {
	cmd := exec.Command("ifconfig", name, "destroy")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to destroy lagg: %s", string(out))
	}
	return nil
}
