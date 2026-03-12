package firewall

import (
	"encoding/xml"
	"fmt"
	"time"
)

type FilterRule struct {
	ID          string     `json:"id" xml:"id"`
	UUID        string     `json:"uuid" xml:"uuid"`
	Sequence    int        `json:"sequence" xml:"sequence"`
	Enabled     bool       `json:"enabled" xml:"enabled"`
	Action      string     `json:"action" xml:"action"`
	Interface   string     `json:"interface" xml:"interface"`
	Direction   string     `json:"direction" xml:"direction"`
	Protocol    string     `json:"protocol" xml:"protocol"`
	Source      *Address   `json:"source" xml:"source"`
	Destination *Address   `json:"destination" xml:"destination"`
	Log         bool       `json:"log" xml:"log"`
	Description string     `json:"description" xml:"description"`
	Schedule    string     `json:"schedule" xml:"schedule"`
	State       string     `json:"state" xml:"state"` // keep, modulate, synproxy
	Track       string     `json:"track" xml:"track"`
	Created     time.Time  `json:"created" xml:"created"`
	Modified    time.Time  `json:"modified" xml:"modified"`
}

type Address struct {
	Network string `json:"network" xml:"network"`
	Port    string `json:"port" xml:"port"`
	Not     bool   `json:"not" xml:"not"`
}

type NATRule struct {
	ID          string   `json:"id" xml:"id"`
	UUID        string   `json:"uuid" xml:"uuid"`
	Enabled     bool     `json:"enabled" xml:"enabled"`
	Type        string   `json:"type" xml:"type"` // nat, rdr
	Interface   string   `json:"interface" xml:"interface"`
	Protocol    string   `json:"protocol" xml:"protocol"`
	Source      *Address `json:"source" xml:"source"`
	Destination *Address `json:"destination" xml:"destination"`
	Target      *Address `json:"target" xml:"target"`
	Log         bool     `json:"log" xml:"log"`
	Description string   `json:"description" xml:"description"`
}

type FirewallConfig struct {
	XMLName     xml.Name     `xml:"zerowall"`
	Version     string       `xml:"version,attr"`
	Rules       []FilterRule `json:"rules" xml:"filter>rule"`
	NATRules    []NATRule    `json:"nat_rules" xml:"nat>rule"`
	Aliases     []Alias      `json:"aliases" xml:"aliases>alias"`
	Schedules   []Schedule   `json:"schedules" xml:"schedules>schedule"`
	Interfaces  []Interface  `json:"interfaces" xml:"interfaces>interface"`
}

type Interface struct {
	Name        string `json:"name" xml:"name"`
	Type        string `json:"type" xml:"type"` // wan, lan, opt
	Enable      bool   `json:"enable" xml:"enable"`
	IPAddress   string `json:"ip_address" xml:"ipaddr"`
	Subnet      int    `json:"subnet" xml:"subnet"`
	Gateway     string `json:"gateway" xml:"gateway"`
	Description string `json:"description" xml:"descr"`
	MTU         int    `json:"mtu" xml:"mtu"`
}

type Schedule struct {
	Name      string   `json:"name" xml:"name"`
	TimeRange []TimeRange `json:"time_range" xml:"timerange"`
	Weekdays  []int    `json:"weekdays" xml:"weekday"`
}

type TimeRange struct {
	Start string `json:"start" xml:"start"`
	End   string `json:"end" xml:"end"`
}

type RuleSet struct {
	Rules    []FilterRule
	NATRules []NATRule
	Aliases  []Alias
}

func NewFilterRule() *FilterRule {
	return &FilterRule{
		UUID:     generateUUID(),
		Enabled:  true,
		State:    "keep",
		Protocol: "any",
		Direction: "in",
		Created:  time.Now(),
		Modified: time.Now(),
	}
}

func (r *FilterRule) ToPFRule() string {
	rule := fmt.Sprintf("%s %s on %s", r.Action, r.Direction, r.Interface)

	if r.Protocol != "any" {
		rule += fmt.Sprintf(" proto %s", r.Protocol)
	}

	if r.Source != nil {
		src := r.Source.Network
		if r.Source.Not {
			src = "!" + src
		}
		rule += fmt.Sprintf(" from %s", src)
		if r.Source.Port != "" {
			rule += fmt.Sprintf(" port %s", r.Source.Port)
		}
	} else {
		rule += " from any"
	}

	if r.Destination != nil {
		dst := r.Destination.Network
		if r.Destination.Not {
			dst = "!" + dst
		}
		rule += fmt.Sprintf(" to %s", dst)
		if r.Destination.Port != "" {
			rule += fmt.Sprintf(" port %s", r.Destination.Port)
		}
	} else {
		rule += " to any"
	}

	if r.State != "" && r.Action == "pass" {
		rule += fmt.Sprintf(" %s state", r.State)
	}

	if r.Log {
		rule += " log"
	}

	if r.Description != "" {
		rule += fmt.Sprintf(" # %s", r.Description)
	}

	return rule
}

func (n *NATRule) ToPFNATRule() string {
	if n.Type == "rdr" {
		rule := fmt.Sprintf("rdr on %s", n.Interface)
		if n.Protocol != "any" {
			rule += fmt.Sprintf(" proto %s", n.Protocol)
		}
		if n.Source != nil {
			rule += fmt.Sprintf(" from %s", n.Source.Network)
		} else {
			rule += " from any"
		}
		if n.Destination != nil && n.Destination.Network != "" {
			rule += fmt.Sprintf(" to %s", n.Destination.Network)
		} else {
			rule += " to any"
		}
		if n.Destination != nil && n.Destination.Port != "" {
			rule += fmt.Sprintf(" port %s", n.Destination.Port)
		}
		if n.Target != nil {
			rule += fmt.Sprintf(" -> %s", n.Target.Network)
			if n.Target.Port != "" {
				rule += fmt.Sprintf(" port %s", n.Target.Port)
			}
		}
		return rule
	}

	rule := fmt.Sprintf("nat on %s", n.Interface)
	if n.Source != nil {
		rule += fmt.Sprintf(" from %s", n.Source.Network)
	} else {
		rule += " from any"
	}
	rule += " to any"
	rule += fmt.Sprintf(" -> (%s)", n.Interface)

	return rule
}

func generateUUID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
