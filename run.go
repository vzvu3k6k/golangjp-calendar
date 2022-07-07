package calendar

import (
	"fmt"
	"io"
	"net/http"
)

func Run(out io.Writer, args []string) error {
	feedURL := args[0]

	resp, err := http.Get(feedURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	feed, err := NewFeed(resp.Body)
	if err != nil {
		return err
	}

	var events Events
	for _, post := range feed.GetEventPosts()[:1] {
		e, err := post.GetEvents()
		if err != nil {
			return err
		}
		events = append(events, e...)
	}

	calendar := events.buildCalendar()
	fmt.Fprintln(out, calendar)

	return nil
}
