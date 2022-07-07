package calendar

import (
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

type Post struct {
	*gofeed.Item
}

func (p Post) GetEvents() ([]Event, error) {
	baseDate, err := getBaseDate(p.Title)
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(p.Content))
	if err != nil {
		return nil, err
	}

	nodes, err := htmlquery.QueryAll(doc, "//div/ul/li")
	if err != nil {
		return nil, err
	}

	var events []Event
	for _, node := range nodes {
		e, err := parseEvent(node, baseDate)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

var eventTextRegexp = regexp.MustCompile(`(\d{1,2}/\d{1,2})\([日月火水木金土]\)\p{Zs}*(\d{1,2}:\d{2})〜(\d{1,2}:\d{2})\p{Zs}*\[(.+?)\]`)

func parseEvent(node *html.Node, baseDate time.Time) (Event, error) {
	var e Event

	link := htmlquery.FindOne(node, "./a")
	title, url, err := parseTitleAndURL(link)
	if err != nil {
		return Event{}, err
	}
	e.title = title
	e.url = url

	text := htmlquery.InnerText(node)
	matches := eventTextRegexp.FindStringSubmatch(text)
	e.location = matches[4]

	start, end, err := parseStartAndEnd(text, baseDate)
	if err != nil {
		return Event{}, err
	}
	e.start = start
	e.end = end

	return e, nil
}

func parseTitleAndURL(link *html.Node) (string, string, error) {
	title := htmlquery.InnerText(link)
	url := htmlquery.SelectAttr(link, "href")
	return title, url, nil
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

func getBaseDate(postTitle string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("2006年1月のGoイベント一覧", postTitle, loc)
}
