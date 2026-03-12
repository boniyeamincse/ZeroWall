package vpn

import (
	"fmt"
	"os/exec"
	"strings"
)

// WireGuardInterface represents a wg(4) device
type WireGuardInterface struct {
	Name       string `json:"name"`
	ListenPort int    `json:"port"`
	PrivateKey string `json:"private_key"`
	Address    string `json:"address"`
}

// CreateWGInterface initializes the kernel module and device
func CreateWGInterface(config WireGuardInterface) (string, error) {
	// 1. Create interface
	cmd := exec.Command("ifconfig", "wg", "create", "name", config.Name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to create wg: %s", string(out))
	}

	// 2. Set configuration
	cmd = exec.Command("wg", "set", config.Name, "listen-port", fmt.Sprintf("%d", config.ListenPort), "private-key", "/etc/wireguard/"+config.Name+".key")
	// Note: In real app, write private key to file first
	
	// 3. Set IP address
	exec.Command("ifconfig", config.Name, "inet", config.Address, "alias").Run()
	exec.Command("ifconfig", config.Name, "up").Run()

	return config.Name, nil
}
