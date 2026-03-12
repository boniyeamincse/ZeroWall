package system

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// Interface represents a physical or logical network interface in ZeroWall
type Interface struct {
	Name        string   `json:"name"`
	HardwareAdd string   `json:"mac"`
	IPAddresses []string `json:"ips"`
	Status      string   `json:"status"` // up, down
	MTU         int      `json:"mtu"`
	Flags       []string `json:"flags"`
}

// GetInterfaces retrieves all available network interfaces on the system
func GetInterfaces() ([]Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}

	var result []Interface
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		var ipStrs []string
		for _, addr := range addrs {
			ipStrs = append(ipStrs, addr.String())
		}

		status := "down"
		if iface.Flags&net.FlagUp != 0 {
			status = "up"
		}

		var flags []string
		if iface.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "loopback")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "broadcast")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "multicast")
		}

		result = append(result, Interface{
			Name:        iface.Name,
			HardwareAdd: iface.HardwareAddr.String(),
			IPAddresses: ipStrs,
			Status:      status,
			MTU:         iface.MTU,
			Flags:       flags,
		})
	}

	return result, nil
}

// ConfigureInterface applies basic settings to an interface using ifconfig
func ConfigureInterface(name string, up bool, mtu int) error {
	state := "down"
	if up {
		state = "up"
	}

	cmd := exec.Command("ifconfig", name, state, "mtu", fmt.Sprintf("%d", mtu))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ifconfig failed: %s", string(output))
	}

	return nil
}

// GetDetailedStatus returns driver information for a specific interface (FreeBSD specific)
func GetDetailedStatus(name string) (string, error) {
	cmd := exec.Command("ifconfig", name)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parsing would go here to extract media types, duplex settings, etc.
	return strings.TrimSpace(string(output)), nil
}
