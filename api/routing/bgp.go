package routing

import "fmt"

// BGPConfig represents a Border Gateway Protocol setup
type BGPConfig struct {
	ASNumber    int      `json:"as"`
	RouterID    string   `json:"router_id"`
	Neighbors   []Neighbor `json:"neighbors"`
	Networks    []string   `json:"networks"`
}

type Neighbor struct {
	RemoteIP string `json:"ip"`
	RemoteAS int    `json:"remote_as"`
	Description string `json:"desc"`
}

// GenerateBGPFragment creates the bgpd configuration block
func (b BGPConfig) GenerateBGPFragment() string {
	fragment := fmt.Sprintf("router bgp %d\n", b.ASNumber)
	fragment += fmt.Sprintf(" bgp router-id %s\n", b.RouterID)
	
	for _, n := range b.Neighbors {
		fragment += fmt.Sprintf(" neighbor %s remote-as %d\n", n.RemoteIP, n.RemoteAS)
		if n.Description != "" {
			fragment += fmt.Sprintf(" neighbor %s description %s\n", n.RemoteIP, n.Description)
		}
	}
	
	for _, net := range b.Networks {
		fragment += fmt.Sprintf(" network %s\n", net)
	}
	
	return fragment + "!\n"
}
