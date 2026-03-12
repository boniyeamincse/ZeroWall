package vpn

import (
	"fmt"
	"os/exec"
)

// SiteToSiteRoute handles routing for VPN tunnels
type SiteToSiteRoute struct {
	Network   string `json:"network"`
	Gateway   string `json:"gateway"`
	Interface string `json:"interface"`
}

// AddRoute injects a static route into the system routing table
func (r SiteToSiteRoute) AddRoute() error {
	cmd := exec.Command("route", "add", "-net", r.Network, r.Gateway)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add route: %s", string(out))
	}
	return nil
}

// DeleteRoute removes a static route
func (r SiteToSiteRoute) DeleteRoute() error {
	cmd := exec.Command("route", "delete", "-net", r.Network, r.Gateway)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete route: %s", string(out))
	}
	return nil
}

// SetupVTI initializes a Virtual Tunnel Interface (FreeBSD if_ipsec)
func SetupVTI(name string, local string, remote string) error {
	exec.Command("ifconfig", name, "create").Run()
	cmd := exec.Command("ifconfig", name, "tunnel", local, remote)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to setup vti tunnel: %s", string(out))
	}
	return nil
}
