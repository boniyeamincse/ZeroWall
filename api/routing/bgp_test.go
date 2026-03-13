package routing

import (
	"strings"
	"testing"
)

func TestBGPConfig_GenerateBGPFragment(t *testing.T) {
	config := BGPConfig{
		ASNumber: 65001,
		RouterID: "1.1.1.1",
		Neighbors: []Neighbor{
			{RemoteIP: "10.0.0.1", RemoteAS: 65002, Description: "ISP-A"},
		},
		Networks: []string{"192.168.10.0/24"},
	}

	got := config.GenerateBGPFragment()

	wants := []string{
		"router bgp 65001",
		"bgp router-id 1.1.1.1",
		"neighbor 10.0.0.1 remote-as 65002",
		"neighbor 10.0.0.1 description ISP-A",
		"network 192.168.10.0/24",
	}

	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("GenerateBGPFragment() = %q, want to contain %q", got, want)
		}
	}
}
