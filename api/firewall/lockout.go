package firewall

// AntiLockoutManager ensures administrative access is never blocked
type AntiLockoutManager struct {
	ManagementPort int
	ManagementIP   string
}

// GenerateBaseRules returns the rules required to prevent admin lockout
func (m AntiLockoutManager) GenerateBaseRules() []Rule {
	return []Rule{
		{
			Action:      "pass",
			Direction:   "in",
			Interface:   "lan", // Usually enforced on LAN
			Protocol:    "tcp",
			Source:      "any",
			Destination: m.ManagementIP,
			Port:        "443", // Default Web UI
			Description: "Anti-Lockout: Allow Web UI access",
		},
		{
			Action:      "pass",
			Direction:   "in",
			Interface:   "lan",
			Protocol:    "tcp",
			Source:      "any",
			Destination: m.ManagementIP,
			Port:        "22", // Default SSH
			Description: "Anti-Lockout: Allow SSH access",
		},
	}
}
