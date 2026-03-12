package firewall

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Rule represents a basic pf firewall rule
type Rule struct {
	Action      string `json:"action"`    // pass, block, scrub
	Direction   string `json:"direction"` // in, out
	Interface   string `json:"interface"` // em0, lan, any
	Protocol    string `json:"protocol"`  // tcp, udp, icmp, any
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Port        string `json:"port"`
	Description string `json:"description"`
}

// PFEngine handles the generation and application of pf.conf
type PFEngine struct {
	ConfigPath string
}

func NewPFEngine() *PFEngine {
	return &PFEngine{ConfigPath: "/etc/pf.conf"}
}

// GenerateConfig creates the string content for pf.conf based on a list of rules
func (e *PFEngine) GenerateConfig(rules []Rule, natRules []string) string {
	var sb strings.Builder

	sb.WriteString("# ZeroWall - Generated pf.conf\n")
	sb.WriteString("set skip on lo\n")
	sb.WriteString("scrub in all\n\n")

	// NAT Rules
	if len(natRules) > 0 {
		sb.WriteString("# NAT & Redirects\n")
		for _, nr := range natRules {
			sb.WriteString(nr + "\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("# Filtering Rules\n")
	sb.WriteString("block all\n") // Default Deny

	for _, r := range rules {
		ruleStr := fmt.Sprintf("%s %s on %s", r.Action, r.Direction, r.Interface)
		if r.Protocol != "any" {
			ruleStr += fmt.Sprintf(" proto %s", r.Protocol)
		}
		ruleStr += fmt.Sprintf(" from %s to %s", r.Source, r.Destination)
		if r.Port != "" {
			ruleStr += fmt.Sprintf(" port %s", r.Port)
		}
		if r.Description != "" {
			ruleStr += fmt.Sprintf(" # %s", r.Description)
		}
		sb.WriteString(ruleStr + "\n")
	}

	return sb.String()
}

// ApplyConfig writes the config and reloads pf
func (e *PFEngine) ApplyConfig(content string) error {
	// 1. Write to temp file for validation
	tempPath := e.ConfigPath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write temp config: %v", err)
	}

	// 2. Validate syntax
	cmd := exec.Command("pfctl", "-nf", tempPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pf syntax check failed: %s", string(out))
	}

	// 3. Move to official path and reload
	if err := os.Rename(tempPath, e.ConfigPath); err != nil {
		return fmt.Errorf("failed to move config: %v", err)
	}

	cmd = exec.Command("pfctl", "-f", e.ConfigPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload pf: %s", string(out))
	}

	return nil
}
