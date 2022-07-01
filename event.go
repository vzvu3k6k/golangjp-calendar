package calendar

import "time"

type event struct {
	start    time.Time
	end      time.Time
	location string
	title    string
	url      string
}
