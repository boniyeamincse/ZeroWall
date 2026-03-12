package firewall

import (
	"fmt"
)

// PolicyRoute defines a source-based routing rule
type PolicyRoute struct {
	Source      string
	Destination string
	Gateway     string
}

// GeneratePFPolicyRules creates pf route-to markings
func GeneratePFPolicyRules(routes []PolicyRoute) string {
	rules := "# Policy Based Routing Rules\n"
	for _, r := range routes {
		rules += fmt.Sprintf("pass in quick on lan from %s to %s route-to %s\n", 
			r.Source, r.Destination, r.Gateway)
	}
	return rules
}
