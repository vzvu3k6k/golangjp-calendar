package calendar

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

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
		baseDate, err := extractBaseDate(post.Title)
		if err != nil {
			return err
		}

		events, err := extractEvents(post.Content, baseDate)
		if err != nil {
			return err
		}
		fmt.Println(events)
	}

	return nil
}

func getEventPosts(source io.Reader) ([]*gofeed.Item, error) {
	feed, err := gofeed.NewParser().Parse(source)
	if err != nil {
		return nil, err
	}

	var items []*gofeed.Item
	for _, item := range feed.Items {
		if item.Categories != nil {
			for _, c := range item.Categories {
				if c == "Events" {
					items = append(items, item)
					break
				}
			}
		}
	}
	return items, nil
}

var eventTextRegexp = regexp.MustCompile(`(\d{1,2}/\d{1,2})\([日月火水木金土]\) (\d{1,2}:\d{2})〜(\d{1,2}:\d{2}) \[(.+?)\]`)

func extractEvents(content string, baseDate time.Time) ([]event, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	var events []event
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		var e event

		link := s.Find("a")
		e.title = link.Text()
		e.url = link.AttrOr("href", "")

		text := s.Contents().First().Text()
		matches := eventTextRegexp.FindStringSubmatch(text)
		e.location = matches[4]

		start, end, _ := parseStartAndEnd(text, baseDate)
		e.start = start
		e.end = end

		events = append(events, e)
	})

	return events, nil
}

func parseStartAndEnd(text string, baseDate time.Time) (time.Time, time.Time, error) {
	matches := eventTextRegexp.FindStringSubmatch(text)
	date, err := parsePartialDate("1/2", matches[1], baseDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start, err := parsePartialTime("15:04", matches[2], date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := parsePartialTime("15:04", matches[3], date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func parsePartialDate(layout, value string, base time.Time) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(base.Year(), t.Month(), t.Day(), 0, 0, 0, 0, base.Location()), nil
}

func parsePartialTime(layout, value string, base time.Time) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(base.Year(), base.Month(), base.Day(), t.Hour(), t.Minute(), 0, 0, base.Location()), nil
}

func extractBaseDate(postTitle string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("2006年1月のGoイベント一覧", postTitle, loc)
}
