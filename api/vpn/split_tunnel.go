package vpn

import "fmt"

// SplitTunnelConfig manages selective routing
type SplitTunnelConfig struct {
	TunnelInterface string   `json:"interface"`
	AllowedNetworks []string `json:"allowed_networks"`
	BypassNetworks  []string `json:"bypass_networks"`
}

// GeneratePFPolicy creates pf rules for selective routing (policy routing)
func (s SplitTunnelConfig) GeneratePFPolicy() string {
	rules := "# ZeroWall Split Tunneling Policy\n"
	
	// Force allowed networks through the tunnel
	for _, net := range s.AllowedNetworks {
		rules += fmt.Sprintf("pass out quick on ! %s to %s route-to %s\n", 
			s.TunnelInterface, net, s.TunnelInterface)
	}

	// Ensure local traffic bypasses the tunnel
	for _, net := range s.BypassNetworks {
		rules += fmt.Sprintf("pass out quick to %s\n", net)
	}

	return rules
}
