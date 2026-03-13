package firewall

import (
	"strings"
	"testing"
)

func TestFloatingRule_GeneratePFString(t *testing.T) {
	tests := []struct {
		name     string
		fr       FloatingRule
		contains []string
	}{
		{
			name: "Quick global pass rule",
			fr: FloatingRule{
				Rule: Rule{
					Action:      "pass",
					Protocol:    "tcp",
					Source:      "any",
					Destination: "any",
					Port:        "80",
				},
				Quick: true,
			},
			contains: []string{"pass quick", "on any", "proto tcp", "port 80"},
		},
		{
			name: "Non-quick interface-specific block rule",
			fr: FloatingRule{
				Rule: Rule{
					Action:      "block",
					Interface:   "em0",
					Protocol:    "udp",
					Source:      "192.168.1.0/24",
					Destination: "any",
					Port:        "53",
				},
				ApplyToInterfaces: []string{"em0"},
				Quick:            false,
			},
			contains: []string{"block", "on em0", "proto udp", "from 192.168.1.0/24", "port 53"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fr.GeneratePFString()
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("GeneratePFString() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}
