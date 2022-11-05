package calendar

import (
	"crypto/sha256"
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
		e := cal.AddEvent(event.uid())
		e.SetStartAt(event.start)
		e.SetEndAt(event.end)
		e.SetSummary(event.title)
		e.SetLocation(event.location)
		e.SetURL(event.url)
	}
	return cal.Serialize()
}

func (e Event) uid() string {
	h := sha256.New()
	h.Write([]byte(e.url))
	return fmt.Sprintf("%x-golangjp-calendar@vzvu3k6k.github.io", h.Sum(nil))
}
