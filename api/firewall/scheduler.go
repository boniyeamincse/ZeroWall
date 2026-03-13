package firewall

import (
	"time"
)

// IsActive checks if the schedule is currently active
func (s Schedule) IsActive() bool {
	now := time.Now()
	currentTime := now.Format("15:04")

	// Check time ranges
	timeActive := false
	if len(s.TimeRange) == 0 {
		timeActive = true
	} else {
		for _, tr := range s.TimeRange {
			if currentTime >= tr.Start && currentTime <= tr.End {
				timeActive = true
				break
			}
		}
	}

	if !timeActive {
		return false
	}

	// Check weekday
	if len(s.Weekdays) == 0 {
		return true
	}

	currentWeekday := int(now.Weekday())
	for _, wd := range s.Weekdays {
		if wd == currentWeekday {
			return true
		}
	}

	return false
}
