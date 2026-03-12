package firewall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type FirewallHandler struct {
	ConfigPath string
	Engine     *PFEngine
}

func NewFirewallHandler() *FirewallHandler {
	return &FirewallHandler{
		ConfigPath: "/etc/zerowall/config.xml",
		Engine:     NewPFEngine(),
	}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *FirewallHandler) GetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    rules,
	})
}

func (h *FirewallHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/firewall/rules/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "rule id required")
		return
	}

	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, rule := range rules {
		if rule.ID == id || rule.UUID == id {
			h.respondJSON(w, http.StatusOK, APIResponse{
				Success: true,
				Data:    rule,
			})
			return
		}
	}

	h.respondError(w, http.StatusNotFound, "rule not found")
}

func (h *FirewallHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var newRule FilterRule
	if err := json.Unmarshal(body, &newRule); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	validator := NewValidator()
	if err := validator.validateRule(&newRule); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newRule.UUID = generateUUID()
	newRule.ID = fmt.Sprintf("%d", len(rules)+1000)
	newRule.Sequence = len(rules) + 1
	newRule.Created = time.Now()
	newRule.Modified = time.Now()

	rules = append(rules, newRule)

	if err := h.saveRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.applyRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, "rule saved but failed to apply: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "rule created successfully",
		Data:    newRule,
	})
}

func (h *FirewallHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/firewall/rules/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "rule id required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var updateRule FilterRule
	if err := json.Unmarshal(body, &updateRule); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	found := false
	for i, rule := range rules {
		if rule.ID == id || rule.UUID == id {
			updateRule.ID = rule.ID
			updateRule.UUID = rule.UUID
			updateRule.Created = rule.Created
			updateRule.Modified = time.Now()
			rules[i] = updateRule
			found = true
			break
		}
	}

	if !found {
		h.respondError(w, http.StatusNotFound, "rule not found")
		return
	}

	if err := h.saveRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.applyRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, "rule updated but failed to apply: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "rule updated successfully",
		Data:    updateRule,
	})
}

func (h *FirewallHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/firewall/rules/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "rule id required")
		return
	}

	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	found := false
	newRules := make([]FilterRule, 0)
	for _, rule := range rules {
		if rule.ID != id && rule.UUID != id {
			newRules = append(newRules, rule)
		} else {
			found = true
		}
	}

	if !found {
		h.respondError(w, http.StatusNotFound, "rule not found")
		return
	}

	for i := range newRules {
		newRules[i].Sequence = i + 1
	}

	if err := h.saveRules(newRules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.applyRules(newRules); err != nil {
		h.respondError(w, http.StatusInternalServerError, "rule deleted but failed to apply: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "rule deleted successfully",
	})
}

func (h *FirewallHandler) ReorderRules(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var reorderRequest struct {
		Rules []struct {
			ID       string `json:"id"`
			Sequence int    `json:"sequence"`
		} `json:"rules"`
	}

	if err := json.Unmarshal(body, &reorderRequest); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ruleMap := make(map[string]*FilterRule)
	for i := range rules {
		ruleMap[rules[i].ID] = &rules[i]
		ruleMap[rules[i].UUID] = &rules[i]
	}

	for _, item := range reorderRequest.Rules {
		if rule, ok := ruleMap[item.ID]; ok {
			rule.Sequence = item.Sequence
		}
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Sequence < rules[j].Sequence
	})

	if err := h.saveRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.applyRules(rules)

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "rules reordered successfully",
	})
}

func (h *FirewallHandler) ToggleRule(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/firewall/rules/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "rule id required")
		return
	}

	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	found := false
	for i, rule := range rules {
		if rule.ID == id || rule.UUID == id {
			rules[i].Enabled = !rules[i].Enabled
			rules[i].Modified = time.Now()
			found = true
			break
		}
	}

	if !found {
		h.respondError(w, http.StatusNotFound, "rule not found")
		return
	}

	if err := h.saveRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.applyRules(rules)

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "rule toggled successfully",
	})
}

func (h *FirewallHandler) GetNATRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.loadNATRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    rules,
	})
}

func (h *FirewallHandler) CreateNATRule(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var newRule NATRule
	if err := json.Unmarshal(body, &newRule); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	rules, err := h.loadNATRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newRule.UUID = generateUUID()
	newRule.ID = fmt.Sprintf("%d", len(rules)+2000)
	rules = append(rules, newRule)

	if err := h.saveNATRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.applyNATRules(rules)

	h.respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "NAT rule created",
		Data:    newRule,
	})
}

func (h *FirewallHandler) DeleteNATRule(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/firewall/nat/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "rule id required")
		return
	}

	rules, err := h.loadNATRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newRules := make([]NATRule, 0)
	for _, rule := range rules {
		if rule.ID != id && rule.UUID != id {
			newRules = append(newRules, rule)
		}
	}

	if err := h.saveNATRules(newRules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.applyNATRules(newRules)

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "NAT rule deleted",
	})
}

func (h *FirewallHandler) ApplyFirewall(w http.ResponseWriter, r *http.Request) {
	rules, err := h.loadRules()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.applyRules(rules); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "firewall rules applied successfully",
	})
}

func (h *FirewallHandler) GetAliases(w http.ResponseWriter, r *http.Request) {
	aliases, err := h.loadAliases()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    aliases,
	})
}

func (h *FirewallHandler) loadRules() ([]FilterRule, error) {
	data, err := os.ReadFile(h.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []FilterRule{}, nil
		}
		return nil, err
	}

	var config FirewallConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Rules, nil
}

func (h *FirewallHandler) saveRules(rules []FilterRule) error {
	data, err := os.ReadFile(h.ConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var config FirewallConfig
	if len(data) > 0 {
		if err := json.Unmarshal(data, &config); err != nil {
			config = FirewallConfig{Version: "1.0"}
		}
	} else {
		config = FirewallConfig{Version: "1.0"}
	}

	config.Rules = rules

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(h.ConfigPath, output, 0600)
}

func (h *FirewallHandler) loadNATRules() ([]NATRule, error) {
	data, err := os.ReadFile(h.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []NATRule{}, nil
		}
		return nil, err
	}

	var config FirewallConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.NATRules, nil
}

func (h *FirewallHandler) saveNATRules(rules []NATRule) error {
	data, err := os.ReadFile(h.ConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var config FirewallConfig
	if len(data) > 0 {
		if err := json.Unmarshal(data, &config); err != nil {
			config = FirewallConfig{Version: "1.0"}
		}
	} else {
		config = FirewallConfig{Version: "1.0"}
	}

	config.NATRules = rules

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(h.ConfigPath, output, 0600)
}

func (h *FirewallHandler) loadAliases() ([]Alias, error) {
	data, err := os.ReadFile(h.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Alias{}, nil
		}
		return nil, err
	}

	var config FirewallConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Aliases, nil
}

func (h *FirewallHandler) applyRules(rules []FilterRule) error {
	var natRules []string
	nats, _ := h.loadNATRules()
	for _, nat := range nats {
		if nat.Enabled {
			natRules = append(natRules, nat.ToPFNATRule())
		}
	}

	filterRules := make([]Rule, 0)
	for _, r := range rules {
		if r.Enabled {
			filterRules = append(filterRules, Rule{
				Action:      r.Action,
				Direction:   r.Direction,
				Interface:   r.Interface,
				Protocol:    r.Protocol,
				Source:      "",
				Destination: "",
				Description: r.Description,
			})
		}
	}

	config := h.Engine.GenerateConfig(filterRules, natRules)
	return h.Engine.ApplyConfig(config)
}

func (h *FirewallHandler) applyNATRules(rules []NATRule) error {
	var natRuleStrs []string
	for _, nat := range rules {
		if nat.Enabled {
			natRuleStrs = append(natRuleStrs, nat.ToPFNATRule())
		}
	}

	filterRules, _ := h.loadRules()
	filterRuleObjs := make([]Rule, 0)
	for _, r := range filterRules {
		if r.Enabled {
			filterRuleObjs = append(filterRuleObjs, Rule{
				Action:      r.Action,
				Direction:   r.Direction,
				Interface:   r.Interface,
				Protocol:    r.Protocol,
				Source:      "",
				Destination: "",
				Description: r.Description,
			})
		}
	}

	config := h.Engine.GenerateConfig(filterRuleObjs, natRuleStrs)
	return h.Engine.ApplyConfig(config)
}

func (h *FirewallHandler) respondJSON(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func (h *FirewallHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}

func (v *Validator) validateRule(rule *FilterRule) error {
	validActions := map[string]bool{"pass": true, "block": true, "reject": true}
	validDirections := map[string]bool{"in": true, "out": true}
	validProtocols := map[string]bool{"tcp": true, "udp": true, "icmp": true, "any": true, "esp": true, "ah": true}
	validStates := map[string]bool{"": true, "keep": true, "modulate": true, "synproxy": true}

	if !validActions[strings.ToLower(rule.Action)] {
		return fmt.Errorf("invalid action: %s", rule.Action)
	}
	if !validDirections[strings.ToLower(rule.Direction)] {
		return fmt.Errorf("invalid direction: %s", rule.Direction)
	}
	if !validProtocols[strings.ToLower(rule.Protocol)] {
		return fmt.Errorf("invalid protocol: %s", rule.Protocol)
	}
	if !validStates[strings.ToLower(rule.State)] {
		return fmt.Errorf("invalid state: %s", rule.State)
	}

	return nil
}
