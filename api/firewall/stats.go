package firewall

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

type FirewallStats struct {
	States     StateStats     `json:"states"`
	Counters   CounterStats   `json:"counters"`
	Interfaces []IfaceStats   `json:"interfaces"`
	Timestamp  time.Time      `json:"timestamp"`
}

type StateStats struct {
	Current    int `json:"current"`
	Peak       int `json:"peak"`
	Max        int `json:"max"`
	Inserted   int `json:"inserted"`
	Removed    int `json:"removed"`
	Expires    int `json:"expires"`
}

type CounterStats struct {
	PacketsIn  uint64 `json:"packets_in"`
	PacketsOut uint64 `json:"packets_out"`
	BytesIn    uint64 `json:"bytes_in"`
	BytesOut   uint64 `json:"bytes_out"`
	Passed     uint64 `json:"passed"`
	Blocked    uint64 `json:"blocked"`
}

type IfaceStats struct {
	Name         string `json:"name"`
	PacketsIn    uint64 `json:"packets_in"`
	BytesIn      uint64 `json:"bytes_in"`
	PacketsOut   uint64 `json:"packets_out"`
	BytesOut     uint64 `json:"bytes_out"`
	ErrorsIn     uint64 `json:"errors_in"`
	ErrorsOut    uint64 `json:"errors_out"`
	Collisions   uint64 `json:"collisions"`
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.collectStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to collect stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *StatsHandler) collectStats() (*FirewallStats, error) {
	stats := &FirewallStats{
		Timestamp: time.Now(),
	}

	if err := h.getPFStats(&stats.States, &stats.Counters); err != nil {
		return stats, err
	}

	if err := h.getIfaceStats(&stats.Interfaces); err != nil {
		return stats, err
	}

	return stats, nil
}

func (h *StatsHandler) getPFStats(states *StateStats, counters *CounterStats) error {
	cmd := exec.Command("pfctl", "-si")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pfctl failed: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.Contains(line, "current entries") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				states.Current, _ = strconv.Atoi(parts[2])
			}
		}
		if strings.Contains(line, "peak entries") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				states.Peak, _ = strconv.Atoi(parts[2])
			}
		}
		if strings.Contains(line, "max entries") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				states.Max, _ = strconv.Atoi(parts[2])
			}
		}
		if strings.Contains(line, "inserted") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				states.Inserted, _ = strconv.Atoi(parts[2])
			}
		}
		if strings.Contains(line, "removed") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				states.Removed, _ = strconv.Atoi(parts[2])
			}
		}
	}

	cmd = exec.Command("pfctl", "-sm")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	lines = strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "states") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				states.Max, _ = strconv.Atoi(strings.TrimSuffix(parts[1], "("))
			}
		}
	}

	return nil
}

func (h *StatsHandler) getIfaceStats(interfaces *[]IfaceStats) error {
	cmd := exec.Command("netstat", "-ib")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	lines := strings.Split(string(output), "\n")
	*interfaces = make([]IfaceStats, 0)

	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		var iface IfaceStats
		iface.Name = strings.TrimSuffix(fields[0], ":")

		iface.PacketsIn, _ = strconv.ParseUint(fields[1], 10, 64)
		iface.BytesIn, _ = strconv.ParseUint(fields[4], 10, 64)
		iface.PacketsOut, _ = strconv.ParseUint(fields[6], 10, 64)
		iface.BytesOut, _ = strconv.ParseUint(fields[9], 10, 64)

		*interfaces = append(*interfaces, iface)
	}

	return nil
}

func (h *StatsHandler) GetStateList(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("pfctl", "-ss")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "failed to get state list", http.StatusInternalServerError)
		return
	}

	type StateEntry struct {
		ID        string `json:"id"`
		Protocol  string `json:"protocol"`
		Source    string `json:"source"`
		Destination string `json:"destination"`
		State     string `json:"state"`
		Age       string `json:"age"`
		Expires   string `json:"expires"`
	}

	lines := strings.Split(string(output), "\n")
	states := make([]StateEntry, 0)

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			state := StateEntry{
				ID:          fmt.Sprintf("%d", i),
				Protocol:    fields[0],
				Source:      fields[2],
				Destination: fields[4],
				State:       fields[5],
			}
			states = append(states, state)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    states,
	})
}

func (h *StatsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	cmd := exec.Command("tcpdump", "-i", "pflog0", "-c", strconv.Itoa(limit), "-n", "-tt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "failed to get logs", http.StatusInternalServerError)
		return
	}

	type LogEntry struct {
		Timestamp string `json:"timestamp"`
		Action    string `json:"action"`
		Interface string `json:"interface"`
		Protocol  string `json:"protocol"`
		Source    string `json:"source"`
		Destination string `json:"destination"`
	}

	lines := strings.Split(string(output), "\n")
	logs := make([]LogEntry, 0)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			entry := LogEntry{
				Timestamp:  fields[0],
				Action:     fields[1],
				Interface:  fields[2],
				Protocol:   fields[3],
				Source:     fields[4],
				Destination: fields[5],
			}
			logs = append(logs, entry)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    logs,
	})
}

func (h *StatsHandler) FlushStates(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("pfctl", "-Fs")
	if err := cmd.Run(); err != nil {
		http.Error(w, "failed to flush states", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "states flushed successfully",
	})
}
