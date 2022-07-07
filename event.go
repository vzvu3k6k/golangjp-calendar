package calendar

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Event struct {
	start    time.Time
	end      time.Time
	location string
	title    string
	url      string
}

type Events []Event

func (es Events) buildCalendar() string {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for _, event := range es {
		e := cal.AddEvent(fmt.Sprintf("id@domain-%d", event.start.Unix())) // TODO: まともなIDを生成する
		e.SetStartAt(event.start)
		e.SetEndAt(event.end)
		e.SetSummary(event.title)
		e.SetLocation(event.location)
		e.SetURL(event.url)
	}
	return cal.Serialize()
}
