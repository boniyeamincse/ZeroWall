package routing

import "fmt"

// OSPFConfig represents an Open Shortest Path First setup
type OSPFConfig struct {
	RouterID string   `json:"router_id"`
	Areas    []OSPFArea `json:"areas"`
}

type OSPFArea struct {
	ID       string   `json:"id"` // 0.0.0.0
	Networks []string `json:"networks"`
}

// GenerateOSPFFragment creates the ospfd configuration block
func (o OSPFConfig) GenerateOSPFFragment() string {
	fragment := "router ospf\n"
	fragment += fmt.Sprintf(" ospf router-id %s\n", o.RouterID)
	
	for _, area := range o.Areas {
		for _, net := range area.Networks {
			fragment += fmt.Sprintf(" network %s area %s\n", net, area.ID)
		}
	}
	
	return fragment + "!\n"
}
