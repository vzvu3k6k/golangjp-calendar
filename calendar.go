package calendar

import (
	"fmt"
	"io"
	"net/http"

	ical "github.com/arran4/golang-ical"
	"github.com/mmcdole/gofeed"
)

const eventCategory = "Events"

func getEventPosts(source io.Reader) ([]*gofeed.Item, error) {
	feed, err := gofeed.NewParser().Parse(source)
	if err != nil {
		return nil, err
	}

	var items []*gofeed.Item
	for _, item := range feed.Items {
		if item.Categories != nil {
			for _, c := range item.Categories {
				if c == eventCategory {
					items = append(items, item)
					break
				}
			}
		}
	}
	return items, nil
}

func extractEvents(content string) []*ical.VEvent {
	return nil
}

func Run(args []string) error {
	feedURL := args[0]

	resp, err := http.Get(feedURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	posts, err := getEventPosts(resp.Body)
	if err != nil {
		return err
	}

	for _, post := range posts {
		events := extractEvents(post.Content)
		fmt.Println(events)
	}

	return nil
}
