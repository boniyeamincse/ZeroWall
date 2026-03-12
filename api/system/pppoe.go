package system

import (
	"fmt"
	"os/exec"
)

// PPPoEConfig represents PPP over Ethernet settings
type PPPoEConfig struct {
	Interface string `json:"interface"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Service   string `json:"service_name"`
}

// ConfigurePPPoE sets up mpd5 or ppp for WAN connection
func ConfigurePPPoE(config PPPoEConfig) error {
	// Simplified ifconfig-based approach (FreeBSD's pppoe(4))
	cmd := exec.Command("ifconfig", "pppoe0", "create")
	exec.Command("ifconfig", "pppoe0", "pppoe", "device", config.Interface).Run()
	exec.Command("ifconfig", "pppoe0", "pppoe", "auth", "proto", "pap", "user", config.Username, "password", config.Password).Run()
	exec.Command("ifconfig", "pppoe0", "up").Run()

	return nil
}
