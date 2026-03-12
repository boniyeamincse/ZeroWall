package firewall

import "fmt"

// NATRule represents a Network Address Translation or Port Forwarding rule
type NATRule struct {
	Type        string `json:"type"` // nat, rdr (port forward)
	Interface   string `json:"interface"`
	Protocol    string `json:"protocol"`
	Source      string `json:"source"`
	ExtPort     string `json:"external_port"`
	InternalIP  string `json:"internal_ip"`
	IntPort     string `json:"internal_port"`
}

// PortForward generates a pf 'rdr' rule string
func (nr NATRule) PortForward() string {
	if nr.Type != "rdr" {
		return ""
	}
	return fmt.Sprintf("rdr on %s proto %s from any to any port %s -> %s port %s",
		nr.Interface, nr.Protocol, nr.ExtPort, nr.InternalIP, nr.IntPort)
}

// OutboundNAT generates a pf 'nat' rule string
func (nr NATRule) OutboundNAT() string {
	if nr.Type != "nat" {
		return ""
	}
	return fmt.Sprintf("nat on %s from %s to any -> (%s)",
		nr.Interface, nr.Source, nr.Interface)
}
