package system

import (
	"fmt"
	"os/exec"
)

// NetflowConfig handles telemetry export via softflowd
type NetflowConfig struct {
	Interface string `json:"interface"`
	Collector string `json:"collector"` // Host:Port
	Version   int    `json:"version"`   // 5, 9, 10 (IPFIX)
}

// EnableNetflow starts the softflowd instance
func (c NetflowConfig) EnableNetflow() error {
	fmt.Printf("Enabling Netflow v%d on %s, exporting to %s\n", c.Version, c.Interface, c.Collector)
	
	// FreeBSD: softflowd -i em0 -n collector:9995 -v 9
	cmd := exec.Command("softflowd", "-i", c.Interface, "-n", c.Collector, "-v", fmt.Sprintf("%d", c.Version))
	return cmd.Run()
}

// DisableNetflow stops the telemetry export
func (c NetflowConfig) DisableNetflow() error {
	cmd := exec.Command("pkill", "-f", fmt.Sprintf("softflowd -i %s", c.Interface))
	return cmd.Run()
}
