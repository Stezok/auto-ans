package timemanager

import "time"

type Manager struct {
	TimeSegments []TimeSegment
}

func (m *Manager) Messages(args ...interface{}) (res []string) {
	hour, min, _ := time.Now().Clock()
	for _, timeSegment := range m.TimeSegments {
		if timeSegment.Test(time.Hour*time.Duration(hour) + time.Minute*time.Duration(min)) {
			res = append(res, timeSegment.Message(args...))
		}
	}
	return
}

func (m *Manager) Push(timeSegments ...TimeSegment) {
	m.TimeSegments = append(m.TimeSegments, timeSegments...)
}

func NewManager() *Manager {
	return &Manager{}
}
