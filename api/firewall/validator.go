package firewall

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidAction      = errors.New("invalid action: must be 'pass', 'block', or 'scrub'")
	ErrInvalidDirection   = errors.New("invalid direction: must be 'in' or 'out'")
	ErrInvalidProtocol    = errors.New("invalid protocol: must be 'tcp', 'udp', 'icmp', or 'any'")
	ErrInvalidInterface   = errors.New("invalid interface name")
	ErrInvalidPort        = errors.New("invalid port number or range")
	ErrInvalidIPAddress   = errors.New("invalid IP address")
	ErrInvalidDescription = errors.New("invalid description: contains forbidden characters")
)

var (
	validInterface = regexp.MustCompile(`^[a-z][a-z0-9]{0,14}$`)
	validIP        = regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(\/\d{1,2})?)|([0-9a-f:]+(\/[0-9]+)?)$`)
	allowedChars    = regexp.MustCompile(`^[a-zA-Z0-9\s\-_\.,;:\(\)]+$`)
)

var validActions = map[string]bool{
	"pass":  true,
	"block": true,
	"scrub": true,
	"nat":   true,
	"rdr":   true,
}

var validDirections = map[string]bool{
	"in": true,
	"out": true,
}

var validProtocols = map[string]bool{
	"tcp":  true,
	"udp":  true,
	"icmp": true,
	"any":  true,
	"esp":  true,
	"ah":   true,
}

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateRule(rule Rule) error {
	if err := v.ValidateAction(rule.Action); err != nil {
		return err
	}
	if err := v.ValidateDirection(rule.Direction); err != nil {
		return err
	}
	if err := v.ValidateProtocol(rule.Protocol); err != nil {
		return err
	}
	if err := v.ValidateInterface(rule.Interface); err != nil {
		return err
	}
	if rule.Port != "" {
		if err := v.ValidatePort(rule.Port); err != nil {
			return err
		}
	}
	if rule.Source != "" && rule.Source != "any" {
		if err := v.ValidateIP(rule.Source); err != nil {
			return err
		}
	}
	if rule.Destination != "" && rule.Destination != "any" {
		if err := v.ValidateIP(rule.Destination); err != nil {
			return err
		}
	}
	if rule.Description != "" {
		if err := v.ValidateDescription(rule.Description); err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) ValidateAction(action string) error {
	action = strings.ToLower(action)
	if !validActions[action] {
		return ErrInvalidAction
	}
	return nil
}

func (v *Validator) ValidateDirection(direction string) error {
	direction = strings.ToLower(direction)
	if !validDirections[direction] {
		return ErrInvalidDirection
	}
	return nil
}

func (v *Validator) ValidateProtocol(protocol string) error {
	protocol = strings.ToLower(protocol)
	if !validProtocols[protocol] {
		return ErrInvalidProtocol
	}
	return nil
}

func (v *Validator) ValidateInterface(iface string) error {
	iface = strings.ToLower(iface)
	if iface == "any" || iface == "lo0" {
		return nil
	}
	if !validInterface.MatchString(iface) {
		return ErrInvalidInterface
	}
	return nil
}

func (v *Validator) ValidatePort(port string) error {
	if port == "any" {
		return nil
	}

	if strings.Contains(port, "-") {
		parts := strings.Split(port, "-")
		if len(parts) != 2 {
			return ErrInvalidPort
		}
		for _, p := range parts {
			if err := v.validatePortNumber(p); err != nil {
				return err
			}
		}
		return nil
	}

	return v.validatePortNumber(port)
}

func (v *Validator) validatePortNumber(port string) error {
	p, err := strconv.Atoi(port)
	if err != nil {
		return ErrInvalidPort
	}
	if p < 1 || p > 65535 {
		return ErrInvalidPort
	}
	return nil
}

func (v *Validator) ValidateIP(ip string) error {
	ip = strings.ToLower(ip)

	if strings.Contains(ip, "/") {
		parts := strings.Split(ip, "/")
		ip = parts[0]
	}

	if net.ParseIP(ip) == nil {
		if !validIP.MatchString(ip) {
			return ErrInvalidIPAddress
		}
	}
	return nil
}

func (v *Validator) ValidateDescription(desc string) error {
	if len(desc) > 255 {
		return ErrInvalidDescription
	}
	if !allowedChars.MatchString(desc) {
		return ErrInvalidDescription
	}
	return nil
}

func (v *Validator) SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = regexp.MustCompile(`[<>;&|$` + "`" + `]`).ReplaceAllString(input, "")
	return input
}

func ValidateRules(rules []Rule) []error {
	validator := NewValidator()
	var errors []error

	for i, rule := range rules {
		if err := validator.ValidateRule(rule); err != nil {
			errors = append(errors, fmt.Errorf("rule %d: %w", i+1, err))
		}
	}

	return errors
}
