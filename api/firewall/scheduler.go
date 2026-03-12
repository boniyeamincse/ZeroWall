package firewall

import (
	"time"
)

// Schedule represents the active time range for a rule
type Schedule struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Weekdays  []int     `json:"weekdays"` // 0=Sunday, 1=Monday...
}

// IsActive checks if the schedule is currently active
func (s Schedule) IsActive() bool {
	now := time.Now()

	// Check time range
	if now.Before(s.StartTime) || now.After(s.EndTime) {
		return false
	}

	// Check weekday
	currentWeekday := int(now.Weekday())
	found := false
	for _, wd := range s.Weekdays {
		if wd == currentWeekday {
			found = true
			break
		}
	}

	return found
}
