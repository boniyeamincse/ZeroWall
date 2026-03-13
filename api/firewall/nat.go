package firewall

import "fmt"

// PortForward generates a pf 'rdr' rule string
func (nr NATRule) PortForward() string {
	if nr.Type != "rdr" || nr.Interface == "" {
		return ""
	}

	src := "any"
	if nr.Source != nil && nr.Source.Network != "" {
		src = nr.Source.Network
	}

	dst := "any"
	if nr.Destination != nil && nr.Destination.Network != "" {
		dst = nr.Destination.Network
	}

	extPort := ""
	if nr.Destination != nil && nr.Destination.Port != "" {
		extPort = "port " + nr.Destination.Port
	}

	target := ""
	if nr.Target != nil {
		target = nr.Target.Network
		if nr.Target.Port != "" {
			target += " port " + nr.Target.Port
		}
	}

	return fmt.Sprintf("rdr on %s proto %s from %s to %s %s -> %s",
		nr.Interface, nr.Protocol, src, dst, extPort, target)
}

// OutboundNAT generates a pf 'nat' rule string
func (nr NATRule) OutboundNAT() string {
	if nr.Type != "nat" || nr.Interface == "" {
		return ""
	}

	src := "any"
	if nr.Source != nil && nr.Source.Network != "" {
		src = nr.Source.Network
	}

	return fmt.Sprintf("nat on %s from %s to any -> (%s)",
		nr.Interface, src, nr.Interface)
}
