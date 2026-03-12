package system

import (
	"fmt"
	"os/exec"
)

// CARPConfig represents a Common Address Redundancy Protocol setup
type CARPConfig struct {
	Interface string `json:"interface"`
	VHID      int    `json:"vhid"`
	Password  string `json:"password"`
	VIP       string `json:"vip"`      // Virtual IP
	AdvSkew   int    `json:"advskew"`  // 0 for Master, >0 for Backup
}

// ConfigureCARP initializes the redundant interface
func ConfigureCARP(config CARPConfig) error {
	cmd := exec.Command("ifconfig", config.Interface, "alias", "vhid", fmt.Sprintf("%d", config.VHID), 
		"pass", config.Password, "alias", config.VIP, "advskew", fmt.Sprintf("%d", config.AdvSkew))
	
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to configure CARP: %s", string(out))
	}
	
	return nil
}

// DemoteCARP lowers priority to trigger failover
func DemoteCARP() error {
	cmd := exec.Command("sysctl", "net.inet.carp.demotion=100")
	return cmd.Run()
}
