package firewall

// FloatingRule represents a rule applied across multiple interfaces or globally
type FloatingRule struct {
	Rule
	ApplyToInterfaces []string `json:"interfaces"`
	Quick            bool     `json:"quick"` // pf 'quick' keyword to stop processing further
}

// GeneratePFString generates the floating rule representation for pf.conf
func (fr FloatingRule) GeneratePFString() string {
	r := fr.Rule
	quickStr := ""
	if fr.Quick {
		quickStr = " quick"
	}

	// If no specific interfaces, it's global
	if len(fr.ApplyToInterfaces) == 0 {
		r.Interface = "any"
	}

	// This is a simplified generator; actual implementation would handle 
	// multiple interfaces by repeating the rule or using interface groups.
	return r.Action + quickStr + " on " + r.Interface + " proto " + r.Protocol + 
		" from " + r.Source + " to " + r.Destination + " port " + r.Port
}
