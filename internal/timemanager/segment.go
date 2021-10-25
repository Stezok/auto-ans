package timemanager

import (
	"fmt"
	"time"
)

type TimeSegment struct {
	Start    time.Time
	End      time.Time
	Template string
}

func (ts *TimeSegment) Test(d time.Duration) bool {
	t := time.Unix(0, int64(d))
	return t.After(ts.Start) == t.Before(ts.End)
}

func (ts *TimeSegment) Message(args ...interface{}) string {
	return fmt.Sprintf(ts.Template, args...)
}

func NewTimeSegment(start, end time.Duration, template string) TimeSegment {
	return TimeSegment{
		Start:    time.Unix(0, int64(start)),
		End:      time.Unix(0, int64(end)),
		Template: template,
	}
}
