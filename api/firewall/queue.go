package firewall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type QueueHandler struct {
	AltqConfigPath string
}

func NewQueueHandler() *QueueHandler {
	return &QueueHandler{
		AltqConfigPath: "/etc/zerowall/altq.conf",
	}
}

type Queue struct {
	Name       string   `json:"name"`
	Parent     string   `json:"parent"`
	Interface  string   `json:"interface"`
	Bandwidth  int      `json:"bandwidth"` // Mbps
	Priority   int      `json:"priority"`  // 0-7
	Scheduler  string   `json:"scheduler"` // hfsc, fairq, cbq
	Bytes      uint64   `json:"bytes"`
	Packets    uint64   `json:"packets"`
	Dropped    uint64   `json:"dropped"`
	Queued     uint64   `json:"queued"`
	Average    uint64   `json:"average"`
	Enabled    bool     `json:"enabled"`
}

type QueueConfig struct {
	Queues []Queue `json:"queues"`
}

func (h *QueueHandler) GetQueues(w http.ResponseWriter, r *http.Request) {
	queues, err := h.listQueues()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    queues,
	})
}

func (h *QueueHandler) CreateQueue(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var queue Queue
	if err := json.Unmarshal(body, &queue); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if queue.Name == "" || queue.Interface == "" {
		h.respondError(w, http.StatusBadRequest, "name and interface are required")
		return
	}

	if queue.Scheduler == "" {
		queue.Scheduler = "hfsc"
	}
	if queue.Priority == 0 {
		queue.Priority = 3
	}
	if queue.Bandwidth == 0 {
		queue.Bandwidth = 100
	}

	queues := h.loadQueueConfig()
	queues = append(queues, queue)

	if err := h.saveQueueConfig(queues); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.applyQueues(queues); err != nil {
		h.respondError(w, http.StatusInternalServerError, "queue saved but failed to apply: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "queue created successfully",
		"data":    queue,
	})
}

func (h *QueueHandler) UpdateQueue(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/firewall/queues/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "queue name required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var updateQueue Queue
	if err := json.Unmarshal(body, &updateQueue); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	queues := h.loadQueueConfig()
	found := false
	for i, q := range queues {
		if q.Name == id {
			updateQueue.Name = q.Name
			updateQueue.Parent = q.Parent
			updateQueue.Interface = q.Interface
			queues[i] = updateQueue
			found = true
			break
		}
	}

	if !found {
		h.respondError(w, http.StatusNotFound, "queue not found")
		return
	}

	if err := h.saveQueueConfig(queues); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.applyQueues(queues)

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "queue updated successfully",
	})
}

func (h *QueueHandler) DeleteQueue(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/firewall/queues/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "queue name required")
		return
	}

	queues := h.loadQueueConfig()
	newQueues := make([]Queue, 0)
	for _, q := range queues {
		if q.Name != id {
			newQueues = append(newQueues, q)
		}
	}

	if err := h.saveQueueConfig(newQueues); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.applyQueues(newQueues)

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "queue deleted successfully",
	})
}

func (h *QueueHandler) GetQueueStats(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("pfctl", "-vq")
	output, err := cmd.CombinedOutput()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	queues := h.parseQueueStats(string(output))

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    queues,
	})
}

func (h *QueueHandler) listQueues() ([]Queue, error) {
	cmd := exec.Command("pfctl", "-sq")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	queues := make([]Queue, 0)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "queue") {
			continue
		}

		var q Queue
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			q.Name = fields[1]
			q.Enabled = true
		}
		queues = append(queues, q)
	}

	return queues, nil
}

func (h *QueueHandler) parseQueueStats(output string) []Queue {
	queues := make([]Queue, 0)
	lines := strings.Split(output, "\n")

	var currentQueue *Queue
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "queue") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				currentQueue = &Queue{Name: fields[1], Enabled: true}
			}
		} else if currentQueue != nil {
			if strings.Contains(line, "bytes") {
				fmt.Sscanf(line, "%*s %d %*s %d", &currentQueue.Bytes, &currentQueue.Packets)
			}
			if strings.Contains(line, "dropped") {
				fmt.Sscanf(line, "%*s %d", &currentQueue.Dropped)
			}
			if strings.Contains(line, "queued") {
				fmt.Sscanf(line, "%*s %d", &currentQueue.Queued)
			}
		}
		if currentQueue != nil && strings.HasPrefix(line, "queue") && currentQueue.Name != "" {
			queues = append(queues, *currentQueue)
			currentQueue = nil
		}
	}

	return queues
}

func (h *QueueHandler) loadQueueConfig() []Queue {
	data, err := os.ReadFile(h.AltqConfigPath)
	if err != nil {
		return []Queue{}
	}

	var config QueueConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return []Queue{}
	}

	return config.Queues
}

func (h *QueueHandler) saveQueueConfig(queues []Queue) error {
	config := QueueConfig{Queues: queues}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(h.AltqConfigPath, data, 0600)
}

func (h *QueueHandler) applyQueues(queues []Queue) error {
	var sb strings.Builder

	sb.WriteString("# ZeroWall - Generated ALTQ config\n")

	interfaces := make(map[string]bool)
	for _, q := range queues {
		interfaces[q.Interface] = true
	}

	for iface := range interfaces {
		var ifaceQueues []Queue
		for _, q := range queues {
			if q.Interface == iface {
				ifaceQueues = append(ifaceQueues, q)
			}
		}

		if len(ifaceQueues) == 0 {
			continue
		}

		bandwidth := ifaceQueues[0].Bandwidth
		scheduler := ifaceQueues[0].Scheduler
		if scheduler == "" {
			scheduler = "hfsc"
		}

		sb.WriteString(fmt.Sprintf("altq on %s %s bandwidth %dM queue { ", iface, scheduler, bandwidth))

		qNames := make([]string, 0)
		for _, q := range ifaceQueues {
			qNames = append(qNames, q.Name)
		}
		sb.WriteString(strings.Join(qNames, ", "))
		sb.WriteString(" }\n")

		for _, q := range ifaceQueues {
			sb.WriteString(fmt.Sprintf("queue %s on %s ", q.Name, iface))
			if q.Parent != "" {
				sb.WriteString(fmt.Sprintf("parent %s ", q.Parent))
			}
			sb.WriteString(fmt.Sprintf("bandwidth %dq priority %d\n", q.Bandwidth*1000/8, q.Priority))
		}
	}

	content := sb.String()
	tempPath := "/tmp/altq.conf.tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0600); err != nil {
		return err
	}

	cmd := exec.Command("pfctl", "-f", tempPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply queues: %v", err)
	}

	return nil
}

func (h *QueueHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *QueueHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
