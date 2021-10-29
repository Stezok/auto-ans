package timemanager

import (
	"log"
	"time"
)

type Manager struct {
	TimeSegments []TimeSegment
	loc          *time.Location
}

func (m *Manager) Messages(args ...interface{}) (res []string) {
	hour, min, _ := time.Now().In(m.loc).Clock()
	log.Print(time.Now().In(m.loc).Clock())
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
	return &Manager{
		loc: time.FixedZone("", 3600*2),
	}
}
