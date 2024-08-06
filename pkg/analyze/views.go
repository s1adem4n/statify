package analyze

import (
	"errors"
	"sort"
	"time"

	"github.com/pocketbase/pocketbase/models"
	// dataframes gota
)

type DeviceType int

// custom json marshalling
func (d DeviceType) MarshalJSON() ([]byte, error) {
	switch d {
	case DeviceDesktop:
		return []byte(`"desktop"`), nil
	case DeviceMobile:
		return []byte(`"mobile"`), nil
	default:
		return nil, errors.New("invalid device type")
	}
}

const (
	DeviceDesktop DeviceType = iota
	DeviceMobile
)

type View struct {
	Created   time.Time
	Domain    string
	Path      string
	IP        string
	UserAgent string
	Session   string
	Device    DeviceType
}

func (v *View) FromRecord(record *models.Record) {
	v.Created = record.Created.Time()
	v.Domain = record.GetString("domain")
	v.Path = record.GetString("path")
	v.IP = record.GetString("ip")
	v.UserAgent = record.GetString("user_agent")
	v.Session = record.GetString("session")

	if record.GetString("device") == "desktop" {
		v.Device = DeviceDesktop
	} else {
		v.Device = DeviceMobile
	}
}

func RemoveDuplicates(views []View) []View {
	unique := make(map[string]bool)
	var deduped []View

	for _, view := range views {
		key := view.Path + view.Session
		if _, ok := unique[key]; !ok {
			unique[key] = true
			deduped = append(deduped, view)
		}
	}

	return deduped
}

func FilterViewsByRange(views []View, start, end time.Time) []View {
	var filtered []View

	for _, view := range views {
		if view.Created.After(start) && view.Created.Before(end) {
			filtered = append(filtered, view)
		}
	}

	return filtered
}

func CountViewsOverTime(views []View, interval time.Duration) []int {
	if len(views) == 0 {
		return nil
	}

	sort.Slice(views, func(i, j int) bool {
		return views[i].Created.Before(views[j].Created)
	})

	start := views[0].Created
	end := views[len(views)-1].Created

	numIntervals := int(end.Sub(start) / interval)
	counts := make([]int, numIntervals+1)

	for _, view := range views {
		index := int(view.Created.Sub(start) / interval)
		counts[index]++
	}

	return counts
}

func CountViewsByPath(views []View) map[string]int {
	counts := make(map[string]int)

	for _, view := range views {
		counts[view.Path]++
	}

	return counts
}

func CountViewsByDevice(views []View) map[DeviceType]int {
	counts := make(map[DeviceType]int)

	for _, view := range views {
		counts[view.Device]++
	}

	return counts
}

func CountViewsBySession(views []View) map[string]int {
	counts := make(map[string]int)

	for _, view := range views {
		counts[view.Session]++
	}

	return counts
}

func CountViewsByDomain(views []View) map[string]int {
	counts := make(map[string]int)

	for _, view := range views {
		counts[view.Domain]++
	}

	return counts
}
