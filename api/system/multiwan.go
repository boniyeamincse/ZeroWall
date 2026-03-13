package system

import (
	"fmt"
)

// Gateway represents a WAN link
type Gateway struct {
	Name      string
	Interface string
	IP        string
	MonitorIP string
	Status    string // online, offline, pending
}

// MultiWANManager handles gateway failover and monitoring
type MultiWANManager struct {
	Gateways []Gateway
}

// UpdateStatuses checks dpinger/ping results for each gateway
func (m *MultiWANManager) UpdateStatuses() {
	for i, g := range m.Gateways {
		fmt.Printf("Multi-WAN: Monitoring Gateway %s via %s\n", g.Name, g.MonitorIP)
		// cmd := exec.Command("ping", "-c", "1", "-t", "1", g.MonitorIP)
		// err := cmd.Run()
		m.Gateways[i].Status = "online" // Mock
	}
}

// TriggerFailover switches the default gateway if primary is down
func (m *MultiWANManager) TriggerFailover() {
	fmt.Println("Multi-WAN: Checking for failover conditions...")
	// Logic to modify system routing table if needed
}
