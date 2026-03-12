package firewall

import "strings"

// Alias represents a named group of hosts or ports
type Alias struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"` // host, network, port
	Values      []string `json:"values"`
	Description string   `json:"description"`
}

// ToPFTable converts an Alias to a pf 'table' definition
func (a Alias) ToPFTable() string {
	if a.Type == "port" {
		return "" // Ports are usually handled as macros or lists in PF
	}
	return "table <" + a.Name + "> { " + strings.Join(a.Values, ", ") + " }"
}

// ToPFMacro converts a port Alias to a pf macro
func (a Alias) ToPFMacro() string {
	if a.Type != "port" {
		return ""
	}
	return a.Name + " = \"{" + strings.Join(a.Values, ", ") + "}\""
}
