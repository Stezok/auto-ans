package timemanager

import (
	"log"
	"time"
)

type Manager struct {
	LockDays     map[int]struct{}
	TimeSegments []TimeSegment
	loc          *time.Location
	unavalible   string
}

func (m *Manager) Messages(args ...interface{}) (res []string) {
	day := time.Now().In(m.loc).Weekday()
	if _, ok := m.LockDays[int(day)]; ok {
		return []string{m.unavalible}
	}

	hour, min, _ := time.Now().In(m.loc).Clock()
	log.Print(time.Now().In(m.loc).Clock())
	for _, timeSegment := range m.TimeSegments {
		if timeSegment.Test(time.Hour*time.Duration(hour) + time.Minute*time.Duration(min)) {
			res = append(res, timeSegment.Message(args...))
		}
	}
	return
}

func (m *Manager) AddLockDay(days ...int) {
	for _, day := range days {
		m.LockDays[day] = struct{}{}
	}
}

func (m *Manager) Push(timeSegments ...TimeSegment) {
	m.TimeSegments = append(m.TimeSegments, timeSegments...)
}

func NewManager(unav string) *Manager {
	return &Manager{
		unavalible: unav,
		loc:        time.FixedZone("", 3600*2),
		LockDays:   make(map[int]struct{}),
	}
}
